package main

import (
	"database/sql"
	"log"
	"net"
	"os"

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
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}
	dbURL := os.Getenv("DATABASE_URL")

	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL!")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	emailSender := email.NewSMTPSender(smtpHost, smtpPort, smtpUser, smtpPass)

	userRepo := postgres.NewUserRepository(db)

	accessUC := usecase.NewAccessUsecase(userRepo, jwtSecret, emailSender)
	profileUC := usecase.NewProfileUsecase(userRepo)
	adminUC := usecase.NewAdminUsecase(userRepo)

	authServer := grpc.NewAuthServer(accessUC, profileUC, adminUC)

	grpcServer := grpc_lib.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, authServer)

	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("Auth Service gRPC server is running on port %s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
