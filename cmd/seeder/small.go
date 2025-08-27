package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// runSmallSeeder ejecuta un seeder con solo datos básicos para desarrollo
func runSmallSeeder(db *sql.DB) error {
	// Limpiar datos existentes
	if err := cleanDatabase(db); err != nil {
		return fmt.Errorf("error limpiando base de datos: %w", err)
	}

	// Insertar categorías básicas
	if err := seedBasicCategories(db); err != nil {
		return fmt.Errorf("error insertando categorías básicas: %w", err)
	}

	// Insertar usuarios básicos
	if err := seedBasicUsers(db); err != nil {
		return fmt.Errorf("error insertando usuarios básicos: %w", err)
	}

	// Insertar tags básicos
	if err := seedBasicTags(db); err != nil {
		return fmt.Errorf("error insertando tags básicos: %w", err)
	}

	// Insertar posts básicos
	if err := seedBasicPosts(db); err != nil {
		return fmt.Errorf("error insertando posts básicos: %w", err)
	}

	// Insertar comentarios básicos
	if err := seedBasicComments(db); err != nil {
		return fmt.Errorf("error insertando comentarios básicos: %w", err)
	}

	// Insertar relaciones post-tag básicas
	if err := seedBasicPostTags(db); err != nil {
		return fmt.Errorf("error insertando relaciones post-tag básicas: %w", err)
	}

	return nil
}

func seedBasicCategories(db *sql.DB) error {
	log.Println("📂 Insertando categorías básicas...")

	categories := []Category{
		{Name: "Tecnología", Description: "Artículos sobre tecnología, programación y desarrollo", Slug: "tecnologia"},
		{Name: "Ciencia", Description: "Artículos sobre ciencia, investigación y descubrimientos", Slug: "ciencia"},
		{Name: "Salud", Description: "Artículos sobre salud, bienestar y medicina", Slug: "salud"},
		{Name: "Educación", Description: "Artículos sobre educación, aprendizaje y desarrollo personal", Slug: "educacion"},
		{Name: "Entretenimiento", Description: "Artículos sobre entretenimiento, cultura y ocio", Slug: "entretenimiento"},
		{Name: "Deportes", Description: "Artículos sobre deportes, fitness y actividades físicas", Slug: "deportes"},
		{Name: "Negocios", Description: "Artículos sobre negocios, emprendimiento y economía", Slug: "negocios"},
		{Name: "Viajes", Description: "Artículos sobre viajes, turismo y aventuras", Slug: "viajes"},
	}

	for _, cat := range categories {
		query := `
			INSERT INTO categories (name, description, slug, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

		var id string
		err := db.QueryRow(query, cat.Name, cat.Description, cat.Slug, true, time.Now(), time.Now()).Scan(&id)
		if err != nil {
			return fmt.Errorf("error insertando categoría %s: %w", cat.Name, err)
		}

		cat.ID = id
		log.Printf("  - Categoría '%s' insertada con ID: %s", cat.Name, id)
	}

	return nil
}

func seedBasicUsers(db *sql.DB) error {
	log.Println("👥 Insertando usuarios básicos...")

	users := []User{
		{Username: "admin", Email: "admin@goasync.com", PasswordHash: "$2a$10$hashedpassword", FirstName: "Admin", LastName: "User"},
		{Username: "johndoe", Email: "john.doe@example.com", PasswordHash: "$2a$10$hashedpassword", FirstName: "John", LastName: "Doe"},
		{Username: "janesmith", Email: "jane.smith@example.com", PasswordHash: "$2a$10$hashedpassword", FirstName: "Jane", LastName: "Smith"},
		{Username: "bobwilson", Email: "bob.wilson@example.com", PasswordHash: "$2a$10$hashedpassword", FirstName: "Bob", LastName: "Wilson"},
		{Username: "alicebrown", Email: "alice.brown@example.com", PasswordHash: "$2a$10$hashedpassword", FirstName: "Alice", LastName: "Brown"},
	}

	for _, user := range users {
		query := `
			INSERT INTO users (username, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`

		var id string
		err := db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName, true, time.Now(), time.Now()).Scan(&id)
		if err != nil {
			return fmt.Errorf("error insertando usuario %s: %w", user.Username, err)
		}

		user.ID = id
		log.Printf("  - Usuario '%s' insertado con ID: %s", user.Username, id)
	}

	return nil
}

func seedBasicTags(db *sql.DB) error {
	log.Println("🏷️ Insertando tags básicos...")

	tags := []Tag{
		{Name: "Go", Slug: "go", Description: "Lenguaje de programación Go"},
		{Name: "API", Slug: "api", Description: "Interfaces de programación de aplicaciones"},
		{Name: "Docker", Slug: "docker", Description: "Plataforma de contenedores"},
		{Name: "PostgreSQL", Slug: "postgresql", Description: "Base de datos relacional"},
		{Name: "Web Development", Slug: "web-development", Description: "Desarrollo web"},
		{Name: "Microservicios", Slug: "microservicios", Description: "Arquitectura de microservicios"},
		{Name: "Cloud Computing", Slug: "cloud-computing", Description: "Computación en la nube"},
		{Name: "DevOps", Slug: "devops", Description: "Prácticas de desarrollo y operaciones"},
		{Name: "Machine Learning", Slug: "machine-learning", Description: "Aprendizaje automático"},
		{Name: "Data Science", Slug: "data-science", Description: "Ciencia de datos"},
	}

	for _, tag := range tags {
		query := `
			INSERT INTO tags (name, slug, description, created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id`

		var id string
		err := db.QueryRow(query, tag.Name, tag.Slug, tag.Description, time.Now()).Scan(&id)
		if err != nil {
			return fmt.Errorf("error insertando tag %s: %w", tag.Name, err)
		}

		tag.ID = id
		log.Printf("  - Tag '%s' insertado con ID: %s", tag.Name, id)
	}

	return nil
}

func seedBasicPosts(db *sql.DB) error {
	log.Println("📝 Insertando posts básicos...")

	// Obtener IDs necesarios
	var adminID, johnID, janeID, bobID, aliceID string
	var techCategoryID string

	err := db.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID de admin: %w", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE username = 'johndoe'").Scan(&johnID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID de johndoe: %w", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE username = 'janesmith'").Scan(&janeID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID de janesmith: %w", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE username = 'bobwilson'").Scan(&bobID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID de bobwilson: %w", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE username = 'alicebrown'").Scan(&aliceID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID de alicebrown: %w", err)
	}

	err = db.QueryRow("SELECT id FROM categories WHERE slug = 'tecnologia'").Scan(&techCategoryID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID de categoría tecnología: %w", err)
	}

	posts := []Post{
		{
			Title:      "Introducción a Go: El lenguaje del futuro",
			Slug:       "introduccion-a-go-el-lenguaje-del-futuro",
			Content:    "Go es un lenguaje de programación desarrollado por Google que combina la simplicidad de Python con el rendimiento de C++. En este artículo exploraremos sus características principales, ventajas y casos de uso.",
			Excerpt:    "Go es un lenguaje de programación moderno que combina simplicidad y rendimiento.",
			AuthorID:   johnID,
			CategoryID: techCategoryID,
			Status:     "published",
		},
		{
			Title:      "Construyendo APIs RESTful con Go y Gin",
			Slug:       "construyendo-apis-restful-con-go-y-gin",
			Content:    "En este artículo aprenderemos a construir APIs RESTful robustas y escalables usando Go y el framework Gin.",
			Excerpt:    "Aprende a construir APIs RESTful robustas usando Go y Gin.",
			AuthorID:   janeID,
			CategoryID: techCategoryID,
			Status:     "published",
		},
		{
			Title:      "Docker para desarrolladores: Una guía práctica",
			Slug:       "docker-para-desarrolladores-una-guia-practica",
			Content:    "Docker ha revolucionado la forma en que desarrollamos, desplegamos y ejecutamos aplicaciones. En esta guía práctica aprenderemos los conceptos fundamentales.",
			Excerpt:    "Descubre cómo Docker puede mejorar tu flujo de trabajo de desarrollo.",
			AuthorID:   bobID,
			CategoryID: techCategoryID,
			Status:     "published",
		},
		{
			Title:      "PostgreSQL vs MySQL: ¿Cuál elegir para tu proyecto?",
			Slug:       "postgresql-vs-mysql-cual-elegir-para-tu-proyecto",
			Content:    "La elección de la base de datos es una decisión crítica en cualquier proyecto. En este artículo compararemos PostgreSQL y MySQL para ayudarte a tomar la decisión correcta.",
			Excerpt:    "Compara PostgreSQL y MySQL para elegir la base de datos correcta para tu proyecto.",
			AuthorID:   aliceID,
			CategoryID: techCategoryID,
			Status:     "published",
		},
	}

	for i, post := range posts {
		// Calcular fecha de publicación
		publishedAt := time.Now().AddDate(0, 0, -(i + 1))

		query := `
			INSERT INTO posts (title, slug, content, excerpt, author_id, category_id, status, published_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id`

		var id string
		err := db.QueryRow(query, post.Title, post.Slug, post.Content, post.Excerpt, post.AuthorID, post.CategoryID, post.Status, publishedAt, time.Now(), time.Now()).Scan(&id)
		if err != nil {
			return fmt.Errorf("error insertando post '%s': %w", post.Title, err)
		}

		post.ID = id
		log.Printf("  - Post '%s' insertado con ID: %s", post.Title, id)
	}

	return nil
}

func seedBasicComments(db *sql.DB) error {
	log.Println("💬 Insertando comentarios básicos...")

	// Obtener IDs necesarios
	var postID, janeID, bobID string

	err := db.QueryRow("SELECT id FROM posts WHERE slug = 'introduccion-a-go-el-lenguaje-del-futuro'").Scan(&postID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID del post: %w", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE username = 'janesmith'").Scan(&janeID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID de janesmith: %w", err)
	}

	err = db.QueryRow("SELECT id FROM users WHERE username = 'bobwilson'").Scan(&bobID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID de bobwilson: %w", err)
	}

	comments := []struct {
		PostID  string
		UserID  string
		Content string
	}{
		{PostID: postID, UserID: janeID, Content: "Excelente artículo! Go realmente es un lenguaje increíble para microservicios."},
		{PostID: postID, UserID: bobID, Content: "Muy útil para principiantes. ¿Podrías hacer un artículo sobre testing en Go?"},
	}

	for _, comment := range comments {
		query := `
			INSERT INTO comments (post_id, author_id, content, is_approved, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

		var id string
		err := db.QueryRow(query, comment.PostID, comment.UserID, comment.Content, true, time.Now(), time.Now()).Scan(&id)
		if err != nil {
			return fmt.Errorf("error insertando comentario: %w", err)
		}

		log.Printf("  - Comentario insertado con ID: %s", id)
	}

	return nil
}

func seedBasicPostTags(db *sql.DB) error {
	log.Println("🔗 Insertando relaciones post-tag básicas...")

	// Obtener IDs necesarios
	var postID, goTagID, webDevTagID string

	err := db.QueryRow("SELECT id FROM posts WHERE slug = 'introduccion-a-go-el-lenguaje-del-futuro'").Scan(&postID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID del post: %w", err)
	}

	err = db.QueryRow("SELECT id FROM tags WHERE slug = 'go'").Scan(&goTagID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID del tag go: %w", err)
	}

	err = db.QueryRow("SELECT id FROM tags WHERE slug = 'web-development'").Scan(&webDevTagID)
	if err != nil {
		return fmt.Errorf("error obteniendo ID del tag web-development: %w", err)
	}

	postTags := []struct {
		PostID string
		TagID  string
	}{
		{PostID: postID, TagID: goTagID},
		{PostID: postID, TagID: webDevTagID},
	}

	for _, pt := range postTags {
		query := `INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)`

		_, err := db.Exec(query, pt.PostID, pt.TagID)
		if err != nil {
			return fmt.Errorf("error insertando relación post-tag: %w", err)
		}

		log.Printf("  - Relación post-tag insertada: post=%s, tag=%s", pt.PostID, pt.TagID)
	}

	return nil
}
