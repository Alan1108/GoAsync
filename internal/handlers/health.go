package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/alan.bermudez/goasync/pkg/logger"
)

// HealthResponse estructura de respuesta para el endpoint de salud
type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// HealthCheck verifica el estado de salud de la API
func HealthCheck(c *gin.Context) {
	logger.Info("Health check request received")
	
	response := HealthResponse{
		Status:    "ok",
		Message:   "API funcionando correctamente",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}
	
	c.JSON(http.StatusOK, response)
}

// DetailedHealthCheck proporciona información detallada del estado de salud
func DetailedHealthCheck(c *gin.Context) {
	logger.Info("Detailed health check request received")
	
	response := gin.H{
		"status":    "ok",
		"message":   "API funcionando correctamente",
		"timestamp": time.Now(),
		"version":   "1.0.0",
		"services": gin.H{
			"api": gin.H{
				"status": "healthy",
				"uptime": "running",
			},
			"database": gin.H{
				"status": "not_configured",
				"message": "Base de datos no configurada aún",
			},
		},
	}
	
	c.JSON(http.StatusOK, response)
}
