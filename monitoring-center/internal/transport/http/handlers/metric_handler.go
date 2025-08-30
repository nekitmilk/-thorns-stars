package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nekitmilk/monitoring-center/internal/models"
	"github.com/nekitmilk/monitoring-center/internal/storage/mongo"
	"github.com/nekitmilk/monitoring-center/internal/storage/postgres"
)

type MetricHandler struct {
	metricRepo *mongo.MetricRepository
	hostRepo   *postgres.HostRepository
}

func NewMetricHandler(metricRepo *mongo.MetricRepository, hostRepo *postgres.HostRepository) *MetricHandler {
	return &MetricHandler{
		metricRepo: metricRepo,
		hostRepo:   hostRepo,
	}
}

// ReceiveMetrics принимает метрики от агента
// @Summary Receive metrics from agent
// @Description Receive monitoring metrics from agent
// @Tags metrics
// @Accept json
// @Produce json
// @Param request body models.MetricsRequest true "Metrics data"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/metrics [post]
func (h *MetricHandler) ReceiveMetrics(c *gin.Context) {
	var req models.MetricsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Проверяем существование хоста
	ctx := c.Request.Context()
	hostID, err := uuid.Parse(req.HostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid host ID format",
		})
		return
	}

	host, err := h.hostRepo.FindByID(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify host",
		})
		return
	}

	if host == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Host not found",
		})
		return
	}

	// Сохраняем метрики
	if err := h.metricRepo.SaveMetrics(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save metrics",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Metrics received successfully",
		"count":   len(req.Metrics),
	})
}

// GetHostMetrics возвращает метрики хоста
// @Summary Get host metrics
// @Description Get monitoring metrics for specific host
// @Tags metrics
// @Produce json
// @Param host_id path string true "Host ID"
// @Param type query string false "Metric type" Enums(cpu, ram, disk, process, port, container)
// @Param from query string false "Start time (RFC3339)"
// @Param to query string false "End time (RFC3339)"
// @Param limit query int false "Limit results" default(100) minimum(1) maximum(1000)
// @Success 200 {array} models.Metric
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hosts/{host_id}/metrics [get]
func (h *MetricHandler) GetHostMetrics(c *gin.Context) {
	hostID := c.Param("id")
	metricType := models.MetricType(c.Query("type"))

	var from, to time.Time
	var err error

	if fromStr := c.Query("from"); fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid from date format",
			})
			return
		}
	} else {
		from = time.Now().Add(-24 * time.Hour) // default: last 24 hours
	}

	if toStr := c.Query("to"); toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid to date format",
			})
			return
		}
	} else {
		to = time.Now()
	}

	limit := int64(100)
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err = strconv.ParseInt(limitStr, 10, 64); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid limit",
			})
			return
		}
	}

	ctx := c.Request.Context()
	metrics, err := h.metricRepo.GetHostMetrics(ctx, hostID, metricType, from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch metrics",
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetLatestHostMetrics возвращает последние метрики хоста
// @Summary Get latest host metrics
// @Description Get latest monitoring metrics for specific host
// @Tags metrics
// @Produce json
// @Param host_id path string true "Host ID"
// @Success 200 {object} map[string]models.Metric
// @Failure 500 {object} map[string]string
// @Router /api/hosts/{host_id}/metrics/latest [get]
func (h *MetricHandler) GetLatestHostMetrics(c *gin.Context) {
	hostID := c.Param("id")

	ctx := c.Request.Context()
	metrics, err := h.metricRepo.GetLatestMetrics(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch latest metrics",
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
