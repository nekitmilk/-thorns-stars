package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client *mongo.Client
}

func NewMongoStorage(connString string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoStorage{client: client}, nil
}

func (s *MongoStorage) Ping(ctx context.Context) error {
	return s.client.Ping(ctx, nil)
}

func (s *MongoStorage) Close() error {
	return s.client.Disconnect(context.Background())
}
