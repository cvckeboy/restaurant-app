package main

import (
	"context"
	"github.com/cvckeboy/restaurant-app/database"
	"github.com/cvckeboy/restaurant-app/restaurant/handlers"
	"github.com/cvckeboy/restaurant-app/restaurant/services"
	"github.com/cvckeboy/restaurant-app/restaurant/storage"
	"github.com/cvckeboy/restaurant-app/utils"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

const defaultTimeout = 5 * time.Second

func main() {
	// Load configuration
	config, err := utils.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger(config)

	// Initialize database connection pool
	dbUrl := config.Database.ConnectionUrl
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	pool, err := database.NewDatabasePool(ctx, dbUrl, logger)
	if err != nil {
		logger.Error("failed to connect to database", err)
		return
	}
	defer pool.Close()

	// Initialize storages
	productStorage := storage.NewProductStorage(pool)
	//categoryStorage := storage.NewCategoryStorage(pool)

	// Initialize services
	productService := services.NewProductService(productStorage)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService, logger)

	// Set up router and register handlers
	router := gin.Default()
	productHandler.Register(router)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		logger.Error("failed to run server", err)
	}
}
