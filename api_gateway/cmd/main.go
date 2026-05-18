package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Adiilkwz/music-grpc-go/auth"
	"github.com/Adiilkwz/music-grpc-go/catalog"
	"github.com/Adiilkwz/music-grpc-go/streaming"

	"api_gateway/internal/delivery/http"
	"api_gateway/internal/middleware"
)

func main() {
	log.Println("Initializing API Gateway...")

	authConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Auth Service Error: %v", err)
	}
	defer authConn.Close()
	authClient := auth.NewAuthServiceClient(authConn)

	streamingConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Streaming Service Error: %v", err)
	}
	defer streamingConn.Close()
	streamingClient := streaming.NewStreamingServiceClient(streamingConn)

	catalogConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Catalog Service Error: %v", err)
	}
	defer catalogConn.Close()
	catalogClient := catalog.NewCatalogServiceClient(catalogConn)

	r := gin.Default()

	api := r.Group("/api")
	{
		http.RegisterPublicAuthRoutes(api, authClient)
		http.RegisterPublicCatalogRoutes(api, catalogClient)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authClient))
		{
			http.RegisterStreamRoutes(protected, streamingClient, catalogClient)
			http.RegisterProtectedAuthRoutes(protected, authClient)
			http.RegisterProtectedCatalogRoutes(protected, catalogClient)
		}
	}

	log.Println("API Gateway listening on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
