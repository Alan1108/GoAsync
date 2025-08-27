package handlers

import (
	"net/http"

	"github.com/alan.bermudez/goasync/pkg/logger"
	"github.com/gin-gonic/gin"
)

// APIInfoResponse estructura de respuesta para información de la API
type APIInfoResponse struct {
	Message   string `json:"message"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// GetAPIInfo retorna información general de la API
func GetAPIInfo(c *gin.Context) {
	logger.Info("API info request received")

	response := APIInfoResponse{
		Message:   "Bienvenido a la API GoAsync",
		Version:   "1.0.0",
		Status:    "active",
		Timestamp: "2024-01-01T00:00:00Z",
	}

	c.JSON(http.StatusOK, response)
}

// GetAPIStatus retorna el estado actual de la API
func GetAPIStatus(c *gin.Context) {
	logger.Info("API status request received")

	response := gin.H{
		"status":  "operational",
		"uptime":  "running",
		"version": "1.0.0",
		"endpoints": []string{
			"GET /health",
			"GET /health/detailed",
			"GET /api/v1/",
			"GET /api/v1/status",
		},
		"features": []string{
			"REST API",
			"Health checks",
			"Structured logging",
			"Environment configuration",
		},
	}

	c.JSON(http.StatusOK, response)
}
