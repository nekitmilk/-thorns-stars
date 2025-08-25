package models

import (
	"time"

	"github.com/google/uuid"
	// uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
)

type HostStatus string

const (
	StatusOnline  HostStatus = "online"
	StatusOffline HostStatus = "offline"
	StatusUnknown HostStatus = "unknown"
)

type Host struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	IP        string     `json:"ip" db:"ip"`
	Priority  int        `json:"priority" db:"priority"`
	Status    HostStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateHostRequest параметры запроса для создания нового хоста
type CreateHostRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=255"`
	IP       string `json:"ip" binding:"required,ip"`
	Priority int    `json:"priority" binding:"required,min=1,max=100"`
}

// HostsQuery параметры запроса для получения хостов
type HostsQuery struct {
	Page     int        `form:"page" json:"page" binding:"omitempty,min=1"`
	Limit    int        `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
	Status   HostStatus `form:"status" json:"status"`
	Priority int        `form:"priority" json:"priority" binding:"omitempty,min=1,max=100"`
	Search   string     `form:"search" json:"search"`
}

// HostsResponse ответ с пагинацией
type HostsResponse struct {
	Hosts       []Host `json:"hosts"`
	Total       int    `json:"total"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	TotalPages  int    `json:"total_pages"`
	HasNext     bool   `json:"has_next"`
	HasPrevious bool   `json:"has_previous"`
}
