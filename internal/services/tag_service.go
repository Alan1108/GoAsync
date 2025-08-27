package services

import (
	"database/sql"
	"fmt"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TagService maneja la lógica de negocio para tags
type TagService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewTagService crea una nueva instancia del servicio de tags
func NewTagService(db *sql.DB, logger *logrus.Logger) *TagService {
	return &TagService{
		db:     db,
		logger: logger,
	}
}

// GetAllTags obtiene todos los tags
func (s *TagService) GetAllTags() ([]models.Tag, error) {
	query := `
		SELECT id, name, slug, description, created_at
		FROM tags 
		ORDER BY name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Errorf("Error obteniendo tags: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.CreatedAt)
		if err != nil {
			s.logger.Errorf("Error escaneando tag: %v", err)
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetTagByID obtiene un tag por su ID
func (s *TagService) GetTagByID(id uuid.UUID) (*models.Tag, error) {
	query := `
		SELECT id, name, slug, description, created_at
		FROM tags 
		WHERE id = $1
	`

	var tag models.Tag
	err := s.db.QueryRow(query, id).Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag no encontrado")
		}
		s.logger.Errorf("Error obteniendo tag por ID: %v", err)
		return nil, err
	}

	return &tag, nil
}

// GetTagBySlug obtiene un tag por su slug
func (s *TagService) GetTagBySlug(slug string) (*models.Tag, error) {
	query := `
		SELECT id, name, slug, description, created_at
		FROM tags 
		WHERE slug = $1
	`

	var tag models.Tag
	err := s.db.QueryRow(query, slug).Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag no encontrado")
		}
		s.logger.Errorf("Error obteniendo tag por slug: %v", err)
		return nil, err
	}

	return &tag, nil
}

// GetTagWithPosts obtiene un tag con sus posts
func (s *TagService) GetTagWithPosts(id uuid.UUID) (*models.Tag, error) {
	tag, err := s.GetTagByID(id)
	if err != nil {
		return nil, err
	}

	// Obtener posts del tag
	postsQuery := `
		SELECT p.id, p.title, p.slug, p.content, p.excerpt, p.author_id, p.category_id,
		       p.status, p.published_at, p.created_at, p.updated_at
		FROM posts p
		JOIN post_tags pt ON p.id = pt.post_id
		WHERE pt.tag_id = $1 AND p.status = 'published'
		ORDER BY p.published_at DESC
	`

	rows, err := s.db.Query(postsQuery, id)
	if err != nil {
		s.logger.Errorf("Error obteniendo posts del tag: %v", err)
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

	tag.Posts = posts
	return tag, nil
}

// GetTagsByPostID obtiene todos los tags de un post
func (s *TagService) GetTagsByPostID(postID uuid.UUID) ([]models.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug, t.description, t.created_at
		FROM tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		WHERE pt.post_id = $1
		ORDER BY t.name
	`

	rows, err := s.db.Query(query, postID)
	if err != nil {
		s.logger.Errorf("Error obteniendo tags del post: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.CreatedAt)
		if err != nil {
			s.logger.Errorf("Error escaneando tag: %v", err)
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// CreateTag crea un nuevo tag
func (s *TagService) CreateTag(req models.TagCreateRequest) (*models.Tag, error) {
	// Verificar si el slug ya existe
	existingTag, _ := s.GetTagBySlug(req.Slug)
	if existingTag != nil {
		return nil, fmt.Errorf("el slug de tag ya existe")
	}

	// Generar slug si no se proporciona
	slug := req.Slug
	if slug == "" {
		slug = req.Name
	}

	query := `
		INSERT INTO tags (name, slug, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, slug, description, created_at
	`

	var tag models.Tag
	err := s.db.QueryRow(query, req.Name, slug, req.Description).Scan(
		&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.CreatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error creando tag: %v", err)
		return nil, err
	}

	return &tag, nil
}

// UpdateTag actualiza un tag existente
func (s *TagService) UpdateTag(id uuid.UUID, req models.TagUpdateRequest) (*models.Tag, error) {
	// Verificar que el tag existe
	existingTag, err := s.GetTagByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Name != "" {
		existingTag.Name = req.Name
	}
	if req.Description != "" {
		existingTag.Description = req.Description
	}
	if req.Slug != "" {
		existingTag.Slug = req.Slug
	}

	query := `
		UPDATE tags 
		SET name = $1, slug = $2, description = $3
		WHERE id = $4
		RETURNING id, name, slug, description, created_at
	`

	var tag models.Tag
	err = s.db.QueryRow(query, existingTag.Name, existingTag.Slug, existingTag.Description, id).Scan(
		&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.CreatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error actualizando tag: %v", err)
		return nil, err
	}

	return &tag, nil
}

// DeleteTag elimina un tag
func (s *TagService) DeleteTag(id uuid.UUID) error {
	// Verificar si hay posts asociados
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM post_tags WHERE tag_id = $1", id).Scan(&count)
	if err != nil {
		s.logger.Errorf("Error verificando posts del tag: %v", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("no se puede eliminar el tag porque tiene posts asociados")
	}

	query := "DELETE FROM tags WHERE id = $1"

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Errorf("Error eliminando tag: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag no encontrado")
	}

	return nil
}

// GetPopularTags obtiene los tags más populares
func (s *TagService) GetPopularTags(limit int) ([]models.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug, t.description, t.created_at, COUNT(pt.post_id) as post_count
		FROM tags t
		LEFT JOIN post_tags pt ON t.id = pt.tag_id
		LEFT JOIN posts p ON pt.post_id = p.id AND p.status = 'published'
		GROUP BY t.id, t.name, t.slug, t.description, t.created_at
		ORDER BY post_count DESC, t.name
		LIMIT $1
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		s.logger.Errorf("Error obteniendo tags populares: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		var postCount int
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.CreatedAt, &postCount)
		if err != nil {
			s.logger.Errorf("Error escaneando tag popular: %v", err)
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
