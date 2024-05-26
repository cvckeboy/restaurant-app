package main

import (
	"context"
	"github.com/cvckeboy/restaurant-app/database"
	"github.com/cvckeboy/restaurant-app/restaurant/handlers"
	"github.com/cvckeboy/restaurant-app/restaurant/services"
	"github.com/cvckeboy/restaurant-app/restaurant/storage"
	"github.com/cvckeboy/restaurant-app/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
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

	dbUrl := config.Database.ConnectionUrl
	logger.Info("Initialize database connection pool")
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	pool, err := database.NewDatabasePool(ctx, dbUrl, logger)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		return
	}

	logger.Info("Connection pool initialized", "max connections", pool.Config().MaxConns)
	defer pool.Close()

	logger.Info("Initialize storages")
	productStorage := storage.NewProductStorage(pool, logger)

	logger.Info("Initialize services")
	productService := services.NewProductService(productStorage, logger)

	logger.Info("Initialize handlers")
	productHandler := handlers.NewProductHandler(productService, logger)

	logger.Info("Set up router and register handlers")
	router := gin.Default()

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},                   // Разрешенные источники (ваш фронтенд)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Разрешенные методы
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"}, // Разрешенные заголовки
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	}
	corsMiddleware := cors.New(corsOptions)

	// Логгирование настроек CORS
	logger.Info("Initializing CORS settings", "allowed origins", corsOptions.AllowedOrigins)

	// Применение CORS middleware
	router.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	})

	productHandler.Register(router)

	logger.Info("Start the server")
	if err := router.Run(":8080"); err != nil {
		logger.Error("failed to run server", err)
	}
}
