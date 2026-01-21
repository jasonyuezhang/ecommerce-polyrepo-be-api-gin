package main

import (
	"log"
	"os"

	"github.com/ecommerce/be-api-gin/internal/config"
	"github.com/ecommerce/be-api-gin/internal/routes"
	grpcclient "github.com/ecommerce/be-api-gin/pkg/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting API Gateway on port %s", cfg.Port)

	// Initialize gRPC clients
	grpcClients, err := grpcclient.NewClients(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer grpcClients.Close()

	// Setup routes
	router := routes.Setup(cfg, grpcClients)

	// Start server
	port := cfg.Port
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	log.Printf("API Gateway listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
