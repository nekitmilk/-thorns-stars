package main

import (
	"context"
	"log"
	"time"

	"github.com/nekitmilk/monitoring-center/internal/config"
	"github.com/nekitmilk/monitoring-center/internal/storage/mongo"
	"github.com/nekitmilk/monitoring-center/internal/storage/postgres"
)

func main() {
	cfg := config.Load()

	// Проверка подключения к PostgreSQL
	pgStorage, err := postgres.NewPostgresStorage(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Postgres connection failed: %v", err)
	}

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

	if err := mongoStorage.Ping(ctx); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	log.Println("Successfully connected to MongoDB")

	// Здесь будет запуск HTTP-сервера
	log.Printf("Starting server on %s", cfg.ServerAddress)
}
