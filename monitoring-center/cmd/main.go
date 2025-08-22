package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nekitmilk/monitoring-center/internal/config"
	"github.com/nekitmilk/monitoring-center/internal/storage/mongo"
	"github.com/nekitmilk/monitoring-center/internal/storage/postgres"
	"github.com/nekitmilk/monitoring-center/internal/transport/http/handlers"
)

func main() {
	cfg := config.Load()

	// Проверка подключения к PostgreSQL
	pgStorage, err := postgres.NewPostgresStorage(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Postgres connection failed: %v", err)
	}
	defer pgStorage.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pgStorage.Ping(ctx); err != nil {
		log.Fatalf("Postgres ping failed: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	// Проверка подключения к MongoDB
	mongoStorage, err := mongo.NewMongoStorage(cfg.MongoURL)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer mongoStorage.Close()

	if err := mongoStorage.Ping(ctx); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	log.Println("Successfully connected to MongoDB")

	// Инициализация репозитория
	hostRepo := postgres.NewHostRepository(pgStorage.GetPool())

	// Инициализация обработчиков
	hostHandler := handlers.NewHostHandler(hostRepo)

	// Настройка роутинга
	router := gin.Default()

	api := router.Group("/api")
	{
		hosts := api.Group("/hosts")
		{
			hosts.POST("", hostHandler.CreateHost)
		}
	}

	// Здесь будет запуск HTTP-сервера
	// log.Printf("Starting server on %s", cfg.ServerAddress)

	// Запуск HTTP сервера с graceful shutdown
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
