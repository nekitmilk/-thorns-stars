package handlers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// GetHostByID возвращает информацию о конкретном хосте
// @Summary Get host by ID
// @Description Get detailed information about a specific host
// @Tags hosts
// @Produce json
// @Param id path string true "Host ID"
// @Success 200 {object} models.Host
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hosts/{id} [get]
func (h *HostHandler) GetHostByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid host ID format",
		})
		return
	}

	ctx := c.Request.Context()
	host, err := h.hostRepo.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch host",
		})
		return
	}

	if host == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Host not found",
		})
		return
	}

	c.JSON(http.StatusOK, host)
}

// UpdateHost обновляет информацию о хосте
// @Summary Update host
// @Description Update information about a specific host
// @Tags hosts
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Param request body models.CreateHostRequest true "Host data"
// @Success 200 {object} models.Host
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hosts/{id} [put]
func (h *HostHandler) UpdateHost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid host ID format",
		})
		return
	}

	var req models.CreateHostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Проверяем существование хоста
	existingHost, err := h.hostRepo.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check host existence",
		})
		return
	}
	if existingHost == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Host not found",
		})
		return
	}

	// Проверяем уникальность имени (исключая текущий хост)
	if exists, err := h.hostRepo.IsNameExistsExcluding(ctx, req.Name, id); err != nil {
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

	// Проверяем уникальность IP (исключая текущий хост)
	if exists, err := h.hostRepo.IsIPExistsExcluding(ctx, req.IP, id); err != nil {
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

	// Обновляем данные
	existingHost.Name = req.Name
	existingHost.IP = req.IP
	existingHost.Priority = req.Priority

	if err := h.hostRepo.Update(ctx, existingHost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update host",
		})
		return
	}

	c.JSON(http.StatusOK, existingHost)
}

// DeleteHost удаляет хост
// @Summary Delete host
// @Description Delete a specific host from monitoring
// @Tags hosts
// @Produce json
// @Param id path string true "Host ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/hosts/{id} [delete]
func (h *HostHandler) DeleteHost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid host ID format",
		})
		return
	}

	ctx := c.Request.Context()

	// Проверяем существование хоста
	existingHost, err := h.hostRepo.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check host existence",
		})
		return
	}
	if existingHost == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Host not found",
		})
		return
	}

	if err := h.hostRepo.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete host",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetMasterHost возвращает текущий мастер-хост
// @Summary Get master host
// @Description Get the current master host (host with highest priority among online hosts)
// @Tags hosts
// @Produce json
// @Success 200 {object} models.Host
// @Success 204 "No master host available"
// @Failure 500 {object} map[string]string
// @Router /api/hosts/master [get]
func (h *HostHandler) GetMasterHost(c *gin.Context) {
	ctx := c.Request.Context()

	masterHost, err := h.hostRepo.FindMasterHost(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to find master host",
			"details": err.Error(),
		})
		return
	}

	if masterHost == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, masterHost)
}
