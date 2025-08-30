package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/nekitmilk/monitoring-center/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MetricRepository struct {
	collection *mongo.Collection
}

func NewMetricRepository(client *mongo.Client, dbname string) *MetricRepository {
	collection := client.Database(dbname).Collection("metrics")
	return &MetricRepository{collection: collection}
}

// CreateIndexes создает индексы для оптимизации запросов
func (r *MetricRepository) CreateIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "host_id", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "timestamp", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

// SaveMetrics сохраняет метрики от агента
func (r *MetricRepository) SaveMetrics(ctx context.Context, req models.MetricsRequest) error {
	if len(req.Metrics) == 0 {
		return nil
	}

	var documents []interface{}
	for _, metric := range req.Metrics {
		metric.HostID = req.HostID
		if req.Timestamp.IsZero() {
			metric.Timestamp = time.Now()
		} else {
			metric.Timestamp = req.Timestamp
		}
		documents = append(documents, metric)
	}

	_, err := r.collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert metrics: %w", err)
	}

	return nil
}

// GetHostMetrics возвращает метрики для конкретного хоста
func (r *MetricRepository) GetHostMetrics(ctx context.Context, hostID string, metricType models.MetricType, from, to time.Time, limit int64) ([]models.Metric, error) {
	filter := bson.M{
		"host_id": hostID,
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}

	if metricType != "" {
		filter["type"] = metricType
	}

	opts := options.Find().SetSort(bson.M{"timestamp": -1})
	if limit > 0 {
		opts.SetLimit(limit)
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find metrics: %w", err)
	}
	defer cursor.Close(ctx)

	var metrics []models.Metric
	if err := cursor.All(ctx, &metrics); err != nil {
		return nil, fmt.Errorf("failed to decode metrics: %w", err)
	}

	return metrics, nil
}

// GetLatestMetrics возвращает последние метрики для хоста
func (r *MetricRepository) GetLatestMetrics(ctx context.Context, hostID string) (map[models.MetricType]models.Metric, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"host_id": hostID}}},
		bson.D{{Key: "$sort", Value: bson.M{"timestamp": -1}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":    "$type",
			"latest": bson.M{"$first": "$$ROOT"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate metrics: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID     models.MetricType `bson:"_id"`
		Latest models.Metric     `bson:"latest"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregation results: %w", err)
	}

	latestMetrics := make(map[models.MetricType]models.Metric)
	for _, result := range results {
		latestMetrics[result.ID] = result.Latest
	}

	return latestMetrics, nil
}

// CleanupOldMetrics удаляет старые метрики
func (r *MetricRepository) CleanupOldMetrics(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"timestamp": bson.M{"$lt": cutoff},
	})

	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old metrics: %w", err)
	}

	return result.DeletedCount, nil
}
