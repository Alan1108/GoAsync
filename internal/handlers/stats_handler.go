package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/alan.bermudez/goasync/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// StatsHandler maneja las peticiones HTTP relacionadas con estadísticas
type StatsHandler struct {
	statsService *services.StatsService
	logger       *logrus.Logger
}

// NewStatsHandler crea una nueva instancia del handler de estadísticas
func NewStatsHandler(statsService *services.StatsService, logger *logrus.Logger) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		logger:       logger,
	}
}

// GetDatabaseStats obtiene estadísticas generales de la base de datos
func (h *StatsHandler) GetDatabaseStats(c *gin.Context) {
	stats, err := h.statsService.GetDatabaseStats()
	if err != nil {
		h.logger.Errorf("Error obteniendo estadísticas de la base de datos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GetActivityLogs obtiene logs de actividad con filtros
func (h *StatsHandler) GetActivityLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 10000 {
		perPage = 10
	}

	filter := models.ActivityLogFilter{
		Page:    page,
		PerPage: perPage,
	}

	// Aplicar filtros opcionales
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			filter.UserID = userID
		}
	}

	if action := c.Query("action"); action != "" {
		filter.Action = action
	}

	if resourceType := c.Query("resource_type"); resourceType != "" {
		filter.ResourceType = resourceType
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = endDate
		}
	}

	logs, err := h.statsService.GetActivityLogs(filter)
	if err != nil {
		h.logger.Errorf("Error obteniendo logs de actividad: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"filters": gin.H{
			"page":          filter.Page,
			"per_page":      filter.PerPage,
			"user_id":       filter.UserID,
			"action":        filter.Action,
			"resource_type": filter.ResourceType,
			"start_date":    filter.StartDate,
			"end_date":      filter.EndDate,
		},
	})
}

// GetRecentActivity obtiene actividad reciente
func (h *StatsHandler) GetRecentActivity(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 10000 {
		limit = 10
	}

	activity, err := h.statsService.GetRecentActivity(limit)
	if err != nil {
		h.logger.Errorf("Error obteniendo actividad reciente: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activity": activity,
	})
}

// GetUserActivity obtiene actividad de un usuario específico
func (h *StatsHandler) GetUserActivity(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 10000 {
		limit = 10
	}

	activity, err := h.statsService.GetUserActivity(userID, limit)
	if err != nil {
		h.logger.Errorf("Error obteniendo actividad del usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"activity": activity,
	})
}

// GetPostStats obtiene estadísticas de posts
func (h *StatsHandler) GetPostStats(c *gin.Context) {
	stats, err := h.statsService.GetPostStats()
	if err != nil {
		h.logger.Errorf("Error obteniendo estadísticas de posts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post_stats": stats,
	})
}

// GetDailyStats obtiene estadísticas diarias
func (h *StatsHandler) GetDailyStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days < 1 || days > 365 {
		days = 7
	}

	stats, err := h.statsService.GetDailyStats(days)
	if err != nil {
		h.logger.Errorf("Error obteniendo estadísticas diarias: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"daily_stats": stats,
	})
}
