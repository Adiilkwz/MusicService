package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"streaming_service/config"
	grpc_delivery "streaming_service/internal/delivery/grpc"
	gwhttp "streaming_service/internal/delivery/http"
	"streaming_service/internal/domain"
	"streaming_service/internal/metrics"
	"streaming_service/internal/repository/nats"
	"streaming_service/internal/repository/postgres"
	redis_repo "streaming_service/internal/repository/redis"
	"streaming_service/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, reading from system env")
	}

	cfg := config.Load()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	db, err := postgres.NewDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		sugar.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	var cacheRepoDomain interface{}
	var redisClient *redis.Client
	if cfg.RedisAddr != "" {
		redisClient = redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: cfg.RedisPassword, DB: 0})
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			sugar.Warnf("Redis not available: %v", err)
			redisClient = nil
		} else {
			cacheRepoDomain = redis_repo.NewCacheRepository(redisClient)
			sugar.Info("Connected to Redis")
		}
	}

	eventPublisher, err := nats.NewEventPublisher(cfg.NATSUrl)
	if err != nil {
		sugar.Warnf("Warning: Failed to connect to NATS, events won't be published: %v", err)
	}
	if eventPublisher != nil {
		defer eventPublisher.Close()
	}

	historyRepo := postgres.NewHistoryRepository(db)
	playlistRepo := postgres.NewPlaylistRepository(db)
	likeRepo := postgres.NewLikeRepository(db)
	trendingRepo := postgres.NewTrendingRepository(db)
	audioRepo := postgres.NewAudioRepository(cfg.AudioDir)

	metrics.Register()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		sugar.Infof("metrics endpoint listening on :9090")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			sugar.Errorf("metrics server error: %v", err)
		}
	}()

	streamingUsecase := usecase.NewStreamingUsecase(historyRepo, trendingRepo, audioRepo, eventPublisher)
	playlistUsecase := usecase.NewPlaylistUsecase(playlistRepo)

	var likeUsecase domain.LikeUsecase
	if cacheRepoDomain != nil {
		likeUsecase = usecase.NewLikeUsecase(likeRepo, cacheRepoDomain.(domain.CacheRepository))
	} else {
		likeUsecase = usecase.NewLikeUsecase(likeRepo, nil)
	}

	handler := grpc_delivery.NewHandler(streamingUsecase, playlistUsecase, likeUsecase)

	go func() {
		gw, err := gwhttp.NewGateway("localhost:" + cfg.ServerPort)
		if err != nil {
			sugar.Errorf("failed to create gateway: %v", err)
			return
		}
		if err := gw.Serve(":8081"); err != nil {
			sugar.Errorf("gateway stopped: %v", err)
		}
	}()

	server := grpc.NewServer()

	grpc_delivery.Register(server, handler)

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		sugar.Fatalf("Failed to listen: %v", err)
	}

	sugar.Infof("Streaming service starting on port %s", cfg.ServerPort)
	if err := server.Serve(lis); err != nil {
		sugar.Fatalf("Failed to serve: %v", err)
	}
}
