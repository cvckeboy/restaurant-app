package main

import (
	"context"
	"errors"
	"github.com/cvckeboy/restaurant-app/database"
	"github.com/cvckeboy/restaurant-app/middleware"
	"github.com/cvckeboy/restaurant-app/restaurant/handlers"
	"github.com/cvckeboy/restaurant-app/restaurant/services"
	"github.com/cvckeboy/restaurant-app/restaurant/storage"
	"github.com/cvckeboy/restaurant-app/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	userStorage := storage.NewUserStorage(pool)
	userService := services.NewUserService(userStorage)
	userHandler := handlers.NewUserHandler(userService, logger)

	logger.Info("Set up router and register handlers")
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorHandler(logger))

	productHandler.Register(router)
	userHandler.Register(router)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server failed", "error", err)
	}
}
