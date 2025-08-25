package postgres

import (
	"context"
	"fmt"
	"strings"
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

// FindAll возвращает хосты с пагинацией и фильтрацией
func (r *HostRepository) FindAll(ctx context.Context, query models.HostsQuery) ([]models.Host, int, error) {
	// Базовый запрос
	baseQuery := `
        SELECT id, name, ip, priority, status, created_at, updated_at 
        FROM hosts 
        WHERE 1=1
    `

	// Запрос для подсчета общего количества
	countQuery := `SELECT COUNT(*) FROM hosts WHERE 1=1`

	// Параметры для запросов
	// var params []interface{}
	var params []any
	var conditions []string

	// Добавляем условия фильтрации
	if query.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(params)+1))
		params = append(params, query.Status)
	}

	if query.Priority > 0 {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", len(params)+1))
		params = append(params, query.Priority)
	}

	if query.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR ip ILIKE $%d)", len(params)+1, len(params)+1))
		params = append(params, "%"+query.Search+"%")
	}

	// Добавляем условия к запросам
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Добавляем сортировку и пагинацию
	baseQuery += " ORDER BY created_at DESC"
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(params)+1, len(params)+2)

	// Вычисляем offset
	offset := (query.Page - 1) * query.Limit
	params = append(params, query.Limit, offset)

	// Получаем общее количество записей
	var total int
	err := r.pool.QueryRow(ctx, countQuery, params[:len(params)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count hosts: %w", err)
	}

	// Получаем данные
	rows, err := r.pool.Query(ctx, baseQuery, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query hosts: %w", err)
	}
	defer rows.Close()

	var hosts []models.Host
	for rows.Next() {
		var host models.Host
		err := rows.Scan(
			&host.ID,
			&host.Name,
			&host.IP,
			&host.Priority,
			&host.Status,
			&host.CreatedAt,
			&host.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan host: %w", err)
		}
		hosts = append(hosts, host)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating hosts: %w", err)
	}

	return hosts, total, nil
}

// func (r *HostRepository) GetMasterHost(ctx context.Context) (*models.Host, error) {
// 	query := `
// 	    SELECT id, name, ip, priority, status, created_at, updated_at
//         FROM hosts
//         ORDER BY priority DESC
// 		LIMIT 1
// 	`

// 	row := r.pool.QueryRow(ctx, query)
// 	var host models.Host
// 	err := row.Scan(
// 		&host.ID,
// 		&host.Name,
// 		&host.IP,
// 		&host.Priority,
// 		&host.Status,
// 		&host.CreatedAt,
// 		&host.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to scan master host: %w", err)
// 	}

// 	return &host, nil

// }
