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
	metricRepo := mongo.NewMetricRepository(mongoStorage.GetClient(), "monitoring")

	// Инициализация обработчиков
	hostHandler := handlers.NewHostHandler(hostRepo)
	metricHandler := handlers.NewMetricHandler(metricRepo, hostRepo)

	// Создание индексов MongoDB
	indexCtx, indexCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer indexCancel()

	if err := metricRepo.CreateIndexes(indexCtx); err != nil {
		log.Printf("Warning: failed to create MongoDB indexes: %v", err)
	}

	// Настройка роутинга
	// router := gin.Default()
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	api := router.Group("/api")
	{
		hosts := api.Group("/hosts")
		{
			hosts.GET("", hostHandler.GetHosts)             // GET /api/hosts
			hosts.POST("", hostHandler.CreateHost)          // POST /api/hosts
			hosts.GET("/:id", hostHandler.GetHostByID)      // GET /api/hosts/{id}
			hosts.PUT("/:id", hostHandler.UpdateHost)       // PUT /api/hosts/{id}
			hosts.DELETE("/:id", hostHandler.DeleteHost)    // DELETE /api/hosts/{id}
			hosts.GET("/master", hostHandler.GetMasterHost) // GET /api/hosts/master

			// Метрики хоста
			hosts.GET("/:id/metrics", metricHandler.GetHostMetrics)
			hosts.GET("/:id/metrics/latest", metricHandler.GetLatestHostMetrics)
		}

		// Эндпоинт для приема метрик от агентов
		api.POST("/metrics", metricHandler.ReceiveMetrics)
	}

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
