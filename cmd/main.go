package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"product-service/internal/config"
	"product-service/internal/controllers"
	"product-service/internal/repositories"
	"product-service/internal/services"
	"product-service/pkg/queue"

	v1 "product-service/api/v1"

	"github.com/gin-gonic/gin"
)

func main() {
	// load the config
	cfg := config.LoadConfig()

	logger := log.New(os.Stdout, "[MAIN] ", log.LstdFlags)
	logger.Printf("Starting application with %d workers, queue size %d", cfg.Workers, cfg.QueueSize)

	// initialize the dependencies
	productRepo := repositories.NewInMemoryProductRepository()
	eventQueue := queue.NewInMemoryEventQueue(cfg.QueueSize)
	productService := services.NewProductService(productRepo, eventQueue, cfg.Workers)

	// initialize the controllers
	productController := controllers.NewProductController(productService)
	healthController := controllers.NewHealthController()

	// setup the gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// setup the routes
	v1.SetupRoutes(router, productController, healthController)

	// start the product service
	productService.Start()

	// setup the graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Println("Received shutdown signal")
		productService.Stop()
		os.Exit(0)
	}()

	// Start HTTP server
	logger.Printf("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
