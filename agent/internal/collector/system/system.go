package system

import (
	"github.com/nekitmilk/agent/internal/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemCollector struct{}

func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

func (c *SystemCollector) Collect() ([]models.Metric, error) {
	var metrics []models.Metric

	// Сбор CPU метрик
	cpuMetrics, err := c.collectCPUMetrics()
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, cpuMetrics...)

	// Сбор RAM метрик
	ramMetrics, err := c.collectRAMMetrics()
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, ramMetrics...)

	// Сбор Disk метрик
	diskMetrics, err := c.collectDiskMetrics()
	if err != nil {
		return nil, err
	}
	metrics = append(metrics, diskMetrics...)

	return metrics, nil
}

func (c *SystemCollector) collectCPUMetrics() ([]models.Metric, error) {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	return []models.Metric{
		{
			Type:  models.MetricCPU,
			Value: percent[0],
			Data: models.CPUData{
				UsagePercent: percent[0],
				Cores:        len(info),
			},
		},
	}, nil
}

func (c *SystemCollector) collectRAMMetrics() ([]models.Metric, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return []models.Metric{
		{
			Type:  models.MetricRAM,
			Value: memory.UsedPercent,
			Data: models.RAMData{
				Total:        memory.Total,
				Used:         memory.Used,
				UsagePercent: memory.UsedPercent,
			},
		},
	}, nil
}

func (c *SystemCollector) collectDiskMetrics() ([]models.Metric, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var metrics []models.Metric
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // Пропускаем проблемные разделы
		}

		metrics = append(metrics, models.Metric{
			Type:  models.MetricDisk,
			Value: usage.UsedPercent,
			Data: models.DiskData{
				MountPoint:   partition.Mountpoint,
				Total:        usage.Total,
				Used:         usage.Used,
				Free:         usage.Free,
				UsagePercent: usage.UsedPercent,
			},
		})
	}

	return metrics, nil
}
