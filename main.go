package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/alan.bermudez/goasync/internal/config"
	"github.com/alan.bermudez/goasync/internal/handlers"
	"github.com/alan.bermudez/goasync/pkg/logger"
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
	log := logger.GetLogger()

	// Conectar a la base de datos
	db, err := sql.Open("postgres", cfg.Database.URL())
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer db.Close()

	// Verificar conexión a la base de datos
	if err := db.Ping(); err != nil {
		log.Fatal("Error verificando conexión a la base de datos:", err)
	}

	log.Info("Conexión a la base de datos establecida exitosamente")

	// Configurar el modo de Gin
	if cfg.Server.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Crear el router de Gin
	router := gin.New()

	// Middleware de recuperación
	router.Use(gin.Recovery())

	// Configurar rutas
	handlers.SetupRoutes(router, db, log)

	// Iniciar el servidor
	log.Printf("Servidor iniciando en el puerto %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
