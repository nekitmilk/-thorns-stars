package handlers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nekitmilk/monitoring-center/internal/models"
	"github.com/nekitmilk/monitoring-center/internal/storage/postgres"
)

type HostHandler struct {
	hostRepo *postgres.HostRepository
}

func NewHostHandler(hostRepo *postgres.HostRepository) *HostHandler {
	return &HostHandler{hostRepo: hostRepo}
}

// CreateHost создает новый хост
// @Summary Create new host
// @Description Add a new host to monitoring system
// @Tags hosts
// @Accept json
// @Produce json
// @Param request body models.CreateHostRequest true "Host data"
// @Success 201 {object} models.Host  // Успех: 201 Created, вернет объект Host
// @Failure 400 {object} map[string]string  // Ошибка клиента
// @Failure 409 {object} map[string]string  // Конфликт (дубликат)
// @Failure 500 {object} map[string]string  // Ошибка сервера
// @Router /api/hosts [post]  // Путь и метод
func (h *HostHandler) CreateHost(c *gin.Context) {
	var req models.CreateHostRequest

	// Валидация входных данных
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	if exists, err := h.hostRepo.IsNameExists(ctx, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check host name",
		})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Host with this name already exists",
		})
		return
	}

	if exists, err := h.hostRepo.IsIPExists(ctx, req.IP); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check host IP",
		})
		return
	} else if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Host with this IP already exists",
		})
		return
	}

	host := &models.Host{
		Name:     req.Name,
		IP:       req.IP,
		Priority: req.Priority,
	}

	if err := h.hostRepo.Create(ctx, host); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create host",
		})
		return
	}

	c.JSON(http.StatusCreated, host)
}

// GetHosts возвращает список хостов с пагинацией и фильтрацией
// @Summary Get all hosts
// @Description Get list of all monitored hosts with pagination and filtering
// @Tags hosts
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Number of items per page" default(20) minimum(1) maximum(100)
// @Param status query string false "Filter by status" Enums(online, offline, unknown)
// @Param priority query int false "Filter by priority" minimum(1) maximum(100)
// @Param search query string false "Search by name or IP"
// @Success 200 {object} models.HostsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hosts [get]
func (h *HostHandler) GetHosts(c *gin.Context) {
	var query models.HostsQuery

	// Биндим параметры запроса
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Устанавливаем значения по умолчанию
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	}

	ctx := c.Request.Context()
	hosts, total, err := h.hostRepo.FindAll(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch hosts",
			"details": err.Error(),
		})
		return
	}

	// Вычисляем пагинацию
	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))
	hasNext := query.Page < totalPages
	hasPrevious := query.Page > 1

	// Формируем ответ
	response := models.HostsResponse{
		Hosts:       hosts,
		Total:       total,
		Page:        query.Page,
		Limit:       query.Limit,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrevious: hasPrevious,
	}

	c.JSON(http.StatusOK, response)
}
