package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetricType string

const (
	MetricCPU       MetricType = "cpu"
	MetricRAM       MetricType = "ram"
	MetricDisk      MetricType = "disk"
	MetricProcess   MetricType = "process"
	MetricPort      MetricType = "port"
	MetricContainer MetricType = "container"
)

type Metric struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HostID    string             `bson:"host_id" json:"host_id"`
	Type      MetricType         `bson:"type" json:"type"`
	Value     float64            `bson:"value" json:"value"`
	Data      any                `bson:"data" json:"data"` // Детальные данные
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type CPUData struct {
	UsagePercent float64 `bson:"usage_percent" json:"usage_percent"`
	Cores        int     `bson:"cores" json:"cores"`
}

type RAMData struct {
	Total        uint64  `bson:"total" json:"total"`
	Used         uint64  `bson:"used" json:"used"`
	UsagePercent float64 `bson:"usage_percent" json:"usage_percent"`
}

type DiskData struct {
	MountPoint   string  `bson:"mount_point" json:"mount_point"`
	Total        uint64  `bson:"total" json:"total"`
	Used         uint64  `bson:"used" json:"used"`
	Free         uint64  `bson:"free" json:"free"`
	UsagePercent float64 `bson:"usage_percent" json:"usage_percent"`
}

type ProcessData struct {
	Name     string  `bson:"name" json:"name"`
	PID      int     `bson:"pid" json:"pid"`
	Status   string  `bson:"status" json:"status"`
	CPUUsage float64 `bson:"cpu_usage" json:"cpu_usage"`
	RAMUsage uint64  `bson:"ram_usage" json:"ram_usage"`
}

type PortData struct {
	Port     int    `bson:"port" json:"port"`
	Protocol string `bson:"protocol" json:"protocol"`
	Status   string `bson:"status" json:"status"` // "open", "closed", "filtered"
	Service  string `bson:"service,omitempty" json:"service,omitempty"`
}

type ContainerData struct {
	ID     string `bson:"id" json:"id"`
	Name   string `bson:"name" json:"name"`
	Image  string `bson:"image" json:"image"`
	Status string `bson:"status" json:"status"`
	State  string `bson:"state" json:"state"` // "running", "exited", etc.
}

// Запрос от агента
type MetricsRequest struct {
	HostID    string    `json:"host_id" binding:"required"`
	Metrics   []Metric  `json:"metrics" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
}
