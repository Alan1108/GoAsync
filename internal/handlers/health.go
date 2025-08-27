package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HealthHandler maneja las peticiones de health check
type HealthHandler struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewHealthHandler crea una nueva instancia del handler de health check
func NewHealthHandler(db *sql.DB, logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// HealthCheck verifica el estado de la aplicación
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Verificar conexión a la base de datos
	var dbStatus string
	var dbError error

	if h.db != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		err := h.db.PingContext(ctx)
		if err != nil {
			dbStatus = "error"
			dbError = err
		} else {
			dbStatus = "ok"
		}
	} else {
		dbStatus = "not_configured"
	}

	// Construir respuesta
	response := gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
		"services": gin.H{
			"database": gin.H{
				"status": dbStatus,
			},
		},
	}

	// Agregar error de base de datos si existe
	if dbError != nil {
		response["services"].(gin.H)["database"].(gin.H)["error"] = dbError.Error()
		h.logger.Errorf("Health check failed - Database error: %v", dbError)
	}

	// Determinar código de estado HTTP
	statusCode := http.StatusOK
	if dbStatus == "error" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
