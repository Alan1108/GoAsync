package services

import (
	"database/sql"
	"fmt"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CategoryService maneja la lógica de negocio para categorías
type CategoryService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewCategoryService crea una nueva instancia del servicio de categorías
func NewCategoryService(db *sql.DB, logger *logrus.Logger) *CategoryService {
	return &CategoryService{
		db:     db,
		logger: logger,
	}
}

// GetAllCategories obtiene todas las categorías
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	query := `
		SELECT id, name, description, slug, is_active, created_at, updated_at
		FROM categories 
		WHERE is_active = true
		ORDER BY name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Errorf("Error obteniendo categorías: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID, &category.Name, &category.Description, &category.Slug,
			&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando categoría: %v", err)
			continue
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategoryByID obtiene una categoría por su ID
func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	query := `
		SELECT id, name, description, slug, is_active, created_at, updated_at
		FROM categories 
		WHERE id = $1
	`

	var category models.Category
	err := s.db.QueryRow(query, id).Scan(
		&category.ID, &category.Name, &category.Description, &category.Slug,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("categoría no encontrada")
		}
		s.logger.Errorf("Error obteniendo categoría por ID: %v", err)
		return nil, err
	}

	return &category, nil
}

// GetCategoryBySlug obtiene una categoría por su slug
func (s *CategoryService) GetCategoryBySlug(slug string) (*models.Category, error) {
	query := `
		SELECT id, name, description, slug, is_active, created_at, updated_at
		FROM categories 
		WHERE slug = $1
	`

	var category models.Category
	err := s.db.QueryRow(query, slug).Scan(
		&category.ID, &category.Name, &category.Description, &category.Slug,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("categoría no encontrada")
		}
		s.logger.Errorf("Error obteniendo categoría por slug: %v", err)
		return nil, err
	}

	return &category, nil
}

// GetCategoryWithPosts obtiene una categoría con sus posts
func (s *CategoryService) GetCategoryWithPosts(id uuid.UUID) (*models.Category, error) {
	category, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	// Obtener posts de la categoría
	postsQuery := `
		SELECT p.id, p.title, p.slug, p.content, p.excerpt, p.author_id, p.category_id,
		       p.status, p.published_at, p.created_at, p.updated_at
		FROM posts p
		WHERE p.category_id = $1 AND p.status = 'published'
		ORDER BY p.published_at DESC
	`

	rows, err := s.db.Query(postsQuery, id)
	if err != nil {
		s.logger.Errorf("Error obteniendo posts de la categoría: %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Content, &post.Excerpt,
			&post.AuthorID, &post.CategoryID, &post.Status, &post.PublishedAt,
			&post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando post: %v", err)
			continue
		}
		posts = append(posts, post)
	}

	category.Posts = posts
	return category, nil
}

// CreateCategory crea una nueva categoría
func (s *CategoryService) CreateCategory(req models.CategoryCreateRequest) (*models.Category, error) {
	// Verificar si el nombre ya existe
	existingCategory, _ := s.GetCategoryBySlug(req.Slug)
	if existingCategory != nil {
		return nil, fmt.Errorf("el slug de categoría ya existe")
	}

	// Generar slug si no se proporciona
	slug := req.Slug
	if slug == "" {
		slug = req.Name
	}

	query := `
		INSERT INTO categories (name, description, slug)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, slug, is_active, created_at, updated_at
	`

	var category models.Category
	err := s.db.QueryRow(query, req.Name, req.Description, slug).Scan(
		&category.ID, &category.Name, &category.Description, &category.Slug,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error creando categoría: %v", err)
		return nil, err
	}

	return &category, nil
}

// UpdateCategory actualiza una categoría existente
func (s *CategoryService) UpdateCategory(id uuid.UUID, req models.CategoryUpdateRequest) (*models.Category, error) {
	// Verificar que la categoría existe
	existingCategory, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Name != "" {
		existingCategory.Name = req.Name
	}
	if req.Description != "" {
		existingCategory.Description = req.Description
	}
	if req.Slug != "" {
		existingCategory.Slug = req.Slug
	}
	if req.IsActive != nil {
		existingCategory.IsActive = *req.IsActive
	}

	query := `
		UPDATE categories 
		SET name = $1, description = $2, slug = $3, is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING id, name, description, slug, is_active, created_at, updated_at
	`

	var category models.Category
	err = s.db.QueryRow(query, existingCategory.Name, existingCategory.Description,
		existingCategory.Slug, existingCategory.IsActive, id).Scan(
		&category.ID, &category.Name, &category.Description, &category.Slug,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error actualizando categoría: %v", err)
		return nil, err
	}

	return &category, nil
}

// DeleteCategory elimina una categoría
func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	// Verificar si hay posts asociados
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM posts WHERE category_id = $1", id).Scan(&count)
	if err != nil {
		s.logger.Errorf("Error verificando posts de categoría: %v", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("no se puede eliminar la categoría porque tiene posts asociados")
	}

	query := "DELETE FROM categories WHERE id = $1"

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Errorf("Error eliminando categoría: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("categoría no encontrada")
	}

	return nil
}
