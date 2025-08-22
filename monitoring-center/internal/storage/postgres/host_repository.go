package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nekitmilk/monitoring-center/internal/models"
	// "go.mongodb.org/mongo-driver/internal/uuid"
)

// Этот класс для работы с конкретной таблицей - Hosts
// Это паттерн проектирования - репозиторий

// Инкапсулируем доступ к БД
type HostRepository struct {
	pool *pgxpool.Pool
}

// Конструктор
func NewHostRepository(pool *pgxpool.Pool) *HostRepository {
	return &HostRepository{pool: pool}
}

// Метод, который добавляет нового хоста в БД
func (r *HostRepository) Create(ctx context.Context, host *models.Host) error {
	query := `INSERT INTO hosts (id, name, ip, priority, status, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	now := time.Now()
	host.ID = uuid.New()
	host.CreatedAt = now
	host.UpdatedAt = now
	host.Status = models.StatusUnknown

	_, err := r.pool.Exec(ctx, query, host.ID, host.Name, host.IP, host.Priority, host.Status, host.CreatedAt, host.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create host: %w", err)
	}

	return nil
}

func (r *HostRepository) IsNameExists(ctx context.Context, name string) (bool, error) {
	query := `SELECT COUNT(*) FROM hosts WHERE name = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, name).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *HostRepository) IsIPExists(ctx context.Context, ip string) (bool, error) {
	query := `SELECT COUNT(*) FROM hosts WHERE ip = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, ip).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
