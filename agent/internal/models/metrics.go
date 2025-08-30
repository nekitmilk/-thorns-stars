package models

import (
	"time"
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
	Type      MetricType `json:"type"`
	Value     float64    `json:"value"`
	Data      any        `json:"data,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

type MetricsRequest struct {
	HostID    string    `json:"host_id"`
	Metrics   []Metric  `json:"metrics"`
	Timestamp time.Time `json:"timestamp"`
}

// Детальные структуры данных (аналогичные ЦМ)
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
