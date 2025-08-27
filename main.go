package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/alan.bermudez/goasync/internal/config"
	"github.com/alan.bermudez/goasync/internal/handlers"
	"github.com/alan.bermudez/goasync/pkg/logger"
	"github.com/alan.bermudez/goasync/pkg/middleware"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env")
	}

	// Cargar configuración
	cfg := config.Load()

	// Inicializar logger
	logger.Init(cfg.Log.Level)

	// Configurar el modo de Gin
	if cfg.Server.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Crear el router de Gin
	router := gin.New()

	// Middleware personalizado
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())

	// Rutas básicas
	setupRoutes(router)

	// Iniciar el servidor
	logger.GetLogger().Printf("Servidor iniciando en el puerto %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Error al iniciar el servidor:", err)
	}
}

func setupRoutes(router *gin.Engine) {
	// Ruta de salud
	router.GET("/health", handlers.HealthCheck)
	router.GET("/health/detailed", handlers.DetailedHealthCheck)

	// Grupo de API v1
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", handlers.GetAPIInfo)
		v1.GET("/status", handlers.GetAPIStatus)

		// Aquí puedes agregar más rutas
		// v1.GET("/users", handlers.GetUsers)
		// v1.POST("/users", handlers.CreateUser)
	}
}
