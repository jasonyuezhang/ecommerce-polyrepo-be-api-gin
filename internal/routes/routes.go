package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ecommerce/be-api-gin/internal/config"
	"github.com/ecommerce/be-api-gin/internal/handlers"
	"github.com/ecommerce/be-api-gin/internal/middleware"
	grpcclient "github.com/ecommerce/be-api-gin/pkg/grpc"
)

// Setup configures all routes and returns the router
func Setup(cfg *config.Config, grpcClients *grpcclient.Clients) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware(cfg))
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RequestIDMiddleware())

	// Health check endpoints
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck(grpcClients))

	// Initialize handlers
	productHandler := handlers.NewProductHandler(grpcClients)
	orderHandler := handlers.NewOrderHandler(grpcClients)

	// Setup product and order routes function
	setupAPIRoutes := func(apiGroup *gin.RouterGroup) {
		// Product routes
		products := apiGroup.Group("/products")
		{
			// Public routes
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)

			// Protected routes
			products.POST("", middleware.AuthMiddleware(cfg), productHandler.CreateProduct)
			products.PUT("/:id", middleware.AuthMiddleware(cfg), productHandler.UpdateProduct)
			products.DELETE("/:id", middleware.AuthMiddleware(cfg), productHandler.DeleteProduct)
			products.PUT("/:id/inventory", middleware.AuthMiddleware(cfg), productHandler.UpdateInventory)
		}

		// Order routes (all protected)
		orders := apiGroup.Group("/orders")
		orders.Use(middleware.AuthMiddleware(cfg))
		{
			orders.GET("", orderHandler.ListOrders)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.POST("", orderHandler.CreateOrder)
			orders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			orders.DELETE("/:id", orderHandler.CancelOrder)
		}
	}

	// API routes without version (for backward compatibility)
	api := router.Group("/api")
	setupAPIRoutes(api)

	// API v1 routes (versioned)
	v1 := router.Group("/api/v1")
	setupAPIRoutes(v1)

	// Handle 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "The requested resource was not found",
		})
	})

	// Handle 405
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":   "Method Not Allowed",
			"message": "The requested method is not allowed for this resource",
		})
	})

	return router
}

// healthCheck returns the health status of the service
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "api-gateway",
	})
}

// readinessCheck checks if all dependencies are ready
func readinessCheck(grpcClients *grpcclient.Clients) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check gRPC connections
		status := grpcClients.HealthCheck(c.Request.Context())

		allHealthy := true
		for _, healthy := range status {
			if !healthy {
				allHealthy = false
				break
			}
		}

		if allHealthy {
			c.JSON(http.StatusOK, gin.H{
				"status":   "ready",
				"services": status,
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "not ready",
				"services": status,
			})
		}
	}
}
