package main

import (
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"streaming_service/config"
	"streaming_service/internal/clients"
	grpc_delivery "streaming_service/internal/delivery/grpc"
	"streaming_service/internal/repository/nats"
	"streaming_service/internal/repository/postgres"
	"streaming_service/internal/usecase"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, reading from system env")
	}

	cfg := config.Load()

	db, err := postgres.NewDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	eventPublisher, err := nats.NewEventPublisher(cfg.NATSUrl)
	if err != nil {
		log.Printf("Warning: Failed to connect to NATS, events won't be published: %v", err)
	}
	if eventPublisher != nil {
		defer eventPublisher.Close()
	}

	authClient, err := clients.NewAuthClient(cfg.AuthServiceUrl)
	if err != nil {
		log.Fatalf("Failed to initialize auth client: %v", err)
	}

	historyRepo := postgres.NewHistoryRepository(db)
	playlistRepo := postgres.NewPlaylistRepository(db)
	likeRepo := postgres.NewLikeRepository(db)
	trendingRepo := postgres.NewTrendingRepository(db)
	audioRepo := postgres.NewAudioRepository(cfg.AudioDir)

	streamingUsecase := usecase.NewStreamingUsecase(historyRepo, trendingRepo, audioRepo, eventPublisher)
	playlistUsecase := usecase.NewPlaylistUsecase(playlistRepo)
	likeUsecase := usecase.NewLikeUsecase(likeRepo)

	handler := grpc_delivery.NewHandler(streamingUsecase, playlistUsecase, likeUsecase)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_delivery.AuthInterceptor(authClient)),
	)

	grpc_delivery.Register(server, handler)

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Streaming service starting on port %s", cfg.ServerPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
