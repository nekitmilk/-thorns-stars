package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	PostgresURL   string
	MongoURL      string
}

func Load() Config {
	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		PostgresURL:   getEnv("POSTGRES_URL", "postgres://monitor:securepass@postgres:5432/monitoring?sslmode=disable"),
		// PostgresURL: getEnv("POSTGRES_URL", "postgres://monitor:securepass@localhost:5432/monitoring?sslmode=disable"),
		MongoURL: getEnv("MONGO_URL", "mongodb://root:example@mongodb:27017"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
