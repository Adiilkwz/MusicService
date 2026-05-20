package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"auth_service/internal/config"
	"auth_service/internal/delivery/grpc"
	"auth_service/internal/infrastructure/email"
	nats_infra "auth_service/internal/infrastructure/nats"
	"auth_service/internal/repository/postgres"
	"auth_service/internal/repository/redis"
	"auth_service/internal/usecase"

	redis_client "github.com/redis/go-redis/v9"

	"github.com/Adiilkwz/music-grpc-go/auth"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	grpc_lib "google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, reading from system env")
	}

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL!")

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsConn.Close()
	log.Println("Successfully connected to NATS!")

	redisClient := redis_client.NewClient(&redis_client.Options{
		Addr: "redis:6379",
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis is unreachable: %v", err)
	}
	log.Println("Successfully connected to Redis!")

	userRepo := postgres.NewUserRepository(db)
	userCacheRepo := redis.NewUserCache(redisClient, 15*time.Minute)
	emailSender := email.NewSMTPSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
	eventPublisher := nats_infra.NewNatsPublisher(natsConn)

	nats_infra.StartRegistrationEmailWorker(natsConn, emailSender)

	accessUC := usecase.NewAccessUsecase(userRepo, eventPublisher, cfg.JWTSecret, emailSender)

	profileUC := usecase.NewProfileUsecase(userRepo, userCacheRepo)

	adminUC := usecase.NewAdminUsecase(userRepo)

	authServer := grpc.NewAuthServer(accessUC, profileUC, adminUC)
	grpcServer := grpc_lib.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, authServer)

	port := ":" + cfg.GRPCPort
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("Auth Service gRPC server is running on port %s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
