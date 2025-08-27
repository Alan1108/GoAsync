package middleware

import (
	"time"

	"github.com/alan.bermudez/goasync/pkg/logger"
	"github.com/gin-gonic/gin"
)

// RequestLogger middleware para logging de requests HTTP
func RequestLogger() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Tiempo de inicio
		start := time.Now()

		// Procesar request
		c.Next()

		// Tiempo de duración
		duration := time.Since(start)

		// Obtener información del request
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Log del request
		logger.WithFields(map[string]interface{}{
			"method":     method,
			"path":       path,
			"status":     status,
			"duration":   duration.String(),
			"client_ip":  clientIP,
			"user_agent": userAgent,
		}).Info("HTTP Request")

		// Log de errores si el status es 4xx o 5xx
		if status >= 400 {
			logger.WithFields(map[string]interface{}{
				"method":     method,
				"path":       path,
				"status":     status,
				"duration":   duration.String(),
				"client_ip":  clientIP,
				"user_agent": userAgent,
			}).Error("HTTP Request Error")
		}
	})
}
