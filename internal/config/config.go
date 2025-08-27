package config

import (
	"os"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port    string
	GinMode string
}

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// LogConfig configuración de logs
type LogConfig struct {
	Level string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "goasync"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
