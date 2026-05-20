package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"catalog_service/config"
	delivery_grpc "catalog_service/internal/delivery/grpc"
	gwhttp "catalog_service/internal/delivery/http"
	"catalog_service/internal/repository/postgres"
	redis_repo "catalog_service/internal/repository/redis"
	"catalog_service/internal/usecase"
	"catalog_service/internal/worker"

	"github.com/Adiilkwz/music-grpc-go/catalog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	google_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		sugar.Fatalf("Ошибка подключения к NATS: %v", err)
	}
	defer func() {
		_ = nc.Drain()
		nc.Close()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbPool, err := pgxpool.New(ctx, cfg.GetDBConnString())
	if err != nil {
		sugar.Fatalf("Не удалось создать пул соединений с БД: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		sugar.Fatalf("База данных недоступна: %v", err)
	}
	sugar.Info("Успешное подключение к PostgreSQL")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		sugar.Fatalf("Redis недоступен: %v", err)
	}
	sugar.Info("Успешное подключение к Redis")

	artistRepo := postgres.NewArtistRepository(dbPool)
	albumRepo := postgres.NewAlbumRepository(dbPool)
	songRepo := postgres.NewSongRepository(dbPool)
	cacheRepo := redis_repo.NewCacheRepository(redisClient)

	artistUC := usecase.NewArtistUsecase(artistRepo, albumRepo)
	albumUC := usecase.NewAlbumUsecase(albumRepo, songRepo)
	songUC := usecase.NewSongUsecase(songRepo)
	searchUC := usecase.NewSearchUsecase(artistRepo, albumRepo, songRepo, cacheRepo)

	natsWorker := worker.NewNatsWorker(nc, songUC)
	go natsWorker.Start(ctx)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		sugar.Fatalf("Не удалось запустить прослушивание порта %s: %v", cfg.GRPCPort, err)
	}

	grpcCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests",
	}, []string{"method"})
	prometheus.MustRegister(grpcCounter)

	grpcServer := google_grpc.NewServer(
		google_grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *google_grpc.UnaryServerInfo, handler google_grpc.UnaryHandler) (interface{}, error) {
			grpcCounter.WithLabelValues(info.FullMethod).Inc()
			return handler(ctx, req)
		}),
	)

	serviceHandler := delivery_grpc.NewServer(artistUC, albumUC, songUC, searchUC)

	catalog.RegisterCatalogServiceServer(grpcServer, serviceHandler)

	reflection.Register(grpcServer)

	sugar.Infof("=== Music Catalog Service запущен на порту %s ===", cfg.GRPCPort)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := ":9090"
		sugar.Infof("metrics endpoint listening on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			sugar.Errorf("metrics server error: %v", err)
		}
	}()

	go func() {
		gw, err := gwhttp.NewGateway("localhost:" + cfg.GRPCPort)
		if err != nil {
			sugar.Errorf("failed to create gateway: %v", err)
			return
		}
		if err := gw.Serve(":8080"); err != nil {
			sugar.Errorf("gateway stopped: %v", err)
		}
	}()

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			sugar.Fatalf("Ошибка при работе gRPC сервера: %v", err)
		}
	}()

	<-ctx.Done()
	sugar.Info("Получен сигнал остановки. Завершаем работу...")

	grpcServer.GracefulStop()
	sugar.Info("Сервер успешно остановлен")
}
