package handlers

import (
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
