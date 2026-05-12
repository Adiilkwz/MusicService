package main

import (
	"database/sql"
	"log"
	"net"

	"auth_service/internal/config"
	"auth_service/internal/delivery/grpc"
	"auth_service/internal/infrastructure/email"
	"auth_service/internal/repository/postgres"
	"auth_service/internal/usecase"

	"github.com/Adiilkwz/music-grpc-go/auth"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	emailSender := email.NewSMTPSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)

	userRepo := postgres.NewUserRepository(db)

	accessUC := usecase.NewAccessUsecase(userRepo, cfg.JWTSecret, emailSender)
	profileUC := usecase.NewProfileUsecase(userRepo)
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
