package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"catalog_service/config"
	delivery_grpc "catalog_service/internal/delivery/grpc"
	"catalog_service/internal/repository/postgres"
	"catalog_service/internal/usecase"

	"github.com/Adiilkwz/music-grpc-go/catalog"
	"github.com/jackc/pgx/v5/pgxpool"
	google_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"catalog_service/internal/worker"

	"github.com/nats-io/nats.go"
)

func main() {
	cfg := config.Load()

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS: %v", err)
	}
	defer nc.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbPool, err := pgxpool.New(ctx, cfg.GetDBConnString())
	if err != nil {
		log.Fatalf("Не удалось создать пул соединений с БД: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL")

	artistRepo := postgres.NewArtistRepository(dbPool)
	albumRepo := postgres.NewAlbumRepository(dbPool)
	songRepo := postgres.NewSongRepository(dbPool)

	artistUC := usecase.NewArtistUsecase(artistRepo, albumRepo)
	albumUC := usecase.NewAlbumUsecase(albumRepo, songRepo)
	songUC := usecase.NewSongUsecase(songRepo)
	searchUC := usecase.NewSearchUsecase(artistRepo, albumRepo, songRepo)

	natsWorker := worker.NewNatsWorker(nc, songUC)
	go natsWorker.Start(ctx)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Не удалось запустить прослушивание порта %s: %v", cfg.GRPCPort, err)
	}

	grpcServer := google_grpc.NewServer()

	serviceHandler := delivery_grpc.NewServer(artistUC, albumUC, songUC, searchUC)

	catalog.RegisterCatalogServiceServer(grpcServer, serviceHandler)

	reflection.Register(grpcServer)

	log.Printf("=== Music Catalog Service запущен на порту %s ===", cfg.GRPCPort)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Ошибка при работе gRPC сервера: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("\nПолучен сигнал остановки. Завершаем работу...")

	grpcServer.GracefulStop()
	log.Println("Сервер успешно остановлен")
}
