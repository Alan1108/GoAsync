package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// User representa un usuario en el sistema
type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Category representa una categoría de posts
type Category struct {
	ID          string
	Name        string
	Description string
	Slug        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Post representa un post o artículo
type Post struct {
	ID          string
	Title       string
	Slug        string
	Content     string
	Excerpt     string
	AuthorID    string
	CategoryID  string
	Status      string
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Tag representa una etiqueta
type Tag struct {
	ID          string
	Name        string
	Slug        string
	Description string
	CreatedAt   time.Time
}

func main() {
	// Inicializar generador de números aleatorios
	rand.Seed(time.Now().UnixNano())

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables del sistema")
	}

	// Configuración de la base de datos
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "goasync")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// String de conexión
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Conectar a la base de datos
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer db.Close()

	// Verificar conexión
	if err := db.Ping(); err != nil {
		log.Fatal("Error haciendo ping a la base de datos:", err)
	}

	log.Println("✅ Conectado exitosamente a la base de datos")

	// Determinar modo de seeder basado en argumentos de línea de comandos
	seederMode := getSeederMode()

	// Ejecutar seeder según el modo
	switch seederMode {
	case "massive":
		log.Println("🚀 Ejecutando seeder MASIVO...")
		if err := runMassiveSeeder(db); err != nil {
			log.Fatal("Error ejecutando seeder masivo:", err)
		}
	case "small":
		log.Println("📝 Ejecutando seeder PEQUEÑO...")
		if err := runSmallSeeder(db); err != nil {
			log.Fatal("Error ejecutando seeder pequeño:", err)
		}
	default:
		log.Println("⚡ Ejecutando seeder por defecto...")
		if err := runSmallSeeder(db); err != nil {
			log.Fatal("Error ejecutando seeder por defecto:", err)
		}
	}

	log.Printf("🎉 Seeder %s completado exitosamente", seederMode)
}

// getSeederMode determina el modo del seeder basado en argumentos de línea de comandos
func getSeederMode() string {
	args := os.Args[1:]

	for i, arg := range args {
		if arg == "--massive" || arg == "-m" {
			return "massive"
		}
		if arg == "--small" || arg == "-s" {
			return "small"
		}
		// Verificar si el siguiente argumento es el valor
		if (arg == "--mode" || arg == "-mode") && i+1 < len(args) {
			return args[i+1]
		}
	}

	return "default"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func cleanDatabase(db *sql.DB) error {
	log.Println("🧹 Limpiando base de datos...")

	tables := []string{
		"post_tags",
		"comments",
		"posts",
		"tags",
		"user_profiles",
		"users",
		"categories",
		"activity_logs",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s", table)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("error limpiando tabla %s: %w", table, err)
		}
		log.Printf("  - Tabla %s limpiada", table)
	}

	return nil
}
