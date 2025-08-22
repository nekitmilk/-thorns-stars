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

type CreateHostRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=255"`
	IP       string `json:"ip" binding:"required,ip"`
	Priority int    `json:"priority" binding:"required,min=1,max=100"`
}
