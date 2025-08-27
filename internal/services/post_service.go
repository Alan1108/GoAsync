package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// PostService maneja la lógica de negocio para posts
type PostService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewPostService crea una nueva instancia del servicio de posts
func NewPostService(db *sql.DB, logger *logrus.Logger) *PostService {
	return &PostService{
		db:     db,
		logger: logger,
	}
}

// GetAllPosts obtiene todos los posts con paginación y filtros
func (s *PostService) GetAllPosts(filter models.PostFilter) (*models.PostListResponse, error) {
	offset := (filter.Page - 1) * filter.PerPage

	// Construir query base
	baseQuery := `
		SELECT p.id, p.title, p.slug, p.content, p.excerpt, p.author_id, p.category_id,
		       p.status, p.published_at, p.created_at, p.updated_at,
		       u.username as author_username, u.first_name as author_first_name, u.last_name as author_last_name,
		       c.name as category_name, c.slug as category_slug
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		LEFT JOIN categories c ON p.category_id = c.id
	`

	whereConditions := []string{}
	args := []interface{}{}
	argCount := 0

	// Aplicar filtros
	if filter.Status != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("p.status = $%d", argCount))
		args = append(args, filter.Status)
	}

	if filter.CategoryID != uuid.Nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("p.category_id = $%d", argCount))
		args = append(args, filter.CategoryID)
	}

	if filter.AuthorID != uuid.Nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("p.author_id = $%d", argCount))
		args = append(args, filter.AuthorID)
	}

	if filter.Search != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("(p.title ILIKE $%d OR p.content ILIKE $%d OR p.excerpt ILIKE $%d)", argCount, argCount, argCount))
		args = append(args, "%"+filter.Search+"%")
	}

	// Construir WHERE clause
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Query para contar total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM posts p %s", whereClause)
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		s.logger.Errorf("Error contando posts: %v", err)
		return nil, err
	}

	// Query para obtener posts
	argCount++
	limitArg := fmt.Sprintf("$%d", argCount)
	argCount++
	offsetArg := fmt.Sprintf("$%d", argCount)
	args = append(args, filter.PerPage, offset)

	query := fmt.Sprintf(`
		%s
		%s
		ORDER BY p.published_at DESC NULLS LAST, p.created_at DESC
		LIMIT %s OFFSET %s
	`, baseQuery, whereClause, limitArg, offsetArg)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		s.logger.Errorf("Error obteniendo posts: %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var authorUsername, authorFirstName, authorLastName sql.NullString
		var categoryName, categorySlug sql.NullString

		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Content, &post.Excerpt,
			&post.AuthorID, &post.CategoryID, &post.Status, &post.PublishedAt,
			&post.CreatedAt, &post.UpdatedAt,
			&authorUsername, &authorFirstName, &authorLastName,
			&categoryName, &categorySlug,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando post: %v", err)
			continue
		}

		// Construir relaciones
		if authorUsername.Valid {
			post.Author = &models.User{
				ID:        post.AuthorID,
				Username:  authorUsername.String,
				FirstName: authorFirstName.String,
				LastName:  authorLastName.String,
			}
		}

		if categoryName.Valid {
			post.Category = &models.Category{
				ID:   post.CategoryID,
				Name: categoryName.String,
				Slug: categorySlug.String,
			}
		}

		posts = append(posts, post)
	}

	// Calcular total de páginas
	totalPages := (total + filter.PerPage - 1) / filter.PerPage

	return &models.PostListResponse{
		Posts:      posts,
		Total:      total,
		Page:       filter.Page,
		PerPage:    filter.PerPage,
		TotalPages: totalPages,
	}, nil
}

// GetPostByID obtiene un post por su ID
func (s *PostService) GetPostByID(id uuid.UUID) (*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.slug, p.content, p.excerpt, p.author_id, p.category_id,
		       p.status, p.published_at, p.created_at, p.updated_at,
		       u.username as author_username, u.first_name as author_first_name, u.last_name as author_last_name,
		       c.name as category_name, c.slug as category_slug
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var post models.Post
	var authorUsername, authorFirstName, authorLastName sql.NullString
	var categoryName, categorySlug sql.NullString

	err := s.db.QueryRow(query, id).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content, &post.Excerpt,
		&post.AuthorID, &post.CategoryID, &post.Status, &post.PublishedAt,
		&post.CreatedAt, &post.UpdatedAt,
		&authorUsername, &authorFirstName, &authorLastName,
		&categoryName, &categorySlug,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post no encontrado")
		}
		s.logger.Errorf("Error obteniendo post por ID: %v", err)
		return nil, err
	}

	// Construir relaciones
	if authorUsername.Valid {
		post.Author = &models.User{
			ID:        post.AuthorID,
			Username:  authorUsername.String,
			FirstName: authorFirstName.String,
			LastName:  authorLastName.String,
		}
	}

	if categoryName.Valid {
		post.Category = &models.Category{
			ID:   post.CategoryID,
			Name: categoryName.String,
			Slug: categorySlug.String,
		}
	}

	return &post, nil
}

// GetPostBySlug obtiene un post por su slug
func (s *PostService) GetPostBySlug(slug string) (*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.slug, p.content, p.excerpt, p.author_id, p.category_id,
		       p.status, p.published_at, p.created_at, p.updated_at,
		       u.username as author_username, u.first_name as author_first_name, u.last_name as author_last_name,
		       c.name as category_name, c.slug as category_slug
		FROM posts p
		LEFT JOIN users u ON p.author_id = u.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.slug = $1
	`

	var post models.Post
	var authorUsername, authorFirstName, authorLastName sql.NullString
	var categoryName, categorySlug sql.NullString

	err := s.db.QueryRow(query, slug).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content, &post.Excerpt,
		&post.AuthorID, &post.CategoryID, &post.Status, &post.PublishedAt,
		&post.CreatedAt, &post.UpdatedAt,
		&authorUsername, &authorFirstName, &authorLastName,
		&categoryName, &categorySlug,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post no encontrado")
		}
		s.logger.Errorf("Error obteniendo post por slug: %v", err)
		return nil, err
	}

	// Construir relaciones
	if authorUsername.Valid {
		post.Author = &models.User{
			ID:        post.AuthorID,
			Username:  authorUsername.String,
			FirstName: authorFirstName.String,
			LastName:  authorLastName.String,
		}
	}

	if categoryName.Valid {
		post.Category = &models.Category{
			ID:   post.CategoryID,
			Name: categoryName.String,
			Slug: categorySlug.String,
		}
	}

	return &post, nil
}

// GetPostWithTags obtiene un post con sus tags
func (s *PostService) GetPostWithTags(id uuid.UUID) (*models.Post, error) {
	post, err := s.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	// Obtener tags del post
	tagsQuery := `
		SELECT t.id, t.name, t.slug, t.description, t.created_at
		FROM tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		WHERE pt.post_id = $1
		ORDER BY t.name
	`

	rows, err := s.db.Query(tagsQuery, id)
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

	post.Tags = tags
	return post, nil
}

// GetPublishedPosts obtiene solo posts publicados
func (s *PostService) GetPublishedPosts(page, perPage int) (*models.PostListResponse, error) {
	filter := models.PostFilter{
		Status:  "published",
		Page:    page,
		PerPage: perPage,
	}
	return s.GetAllPosts(filter)
}

// CreatePost crea un nuevo post
func (s *PostService) CreatePost(req models.PostCreateRequest, authorID uuid.UUID) (*models.Post, error) {
	// Generar slug si no se proporciona
	slug := req.Title
	if slug == "" {
		slug = "untitled"
	}

	// Convertir a slug básico
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Verificar si el slug ya existe
	existingPost, _ := s.GetPostBySlug(slug)
	if existingPost != nil {
		// Agregar timestamp al slug para hacerlo único
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	}

	// Determinar published_at
	var publishedAt *time.Time
	if req.Status == "published" {
		now := time.Now()
		publishedAt = &now
	}

	query := `
		INSERT INTO posts (title, slug, content, excerpt, author_id, category_id, status, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, slug, content, excerpt, author_id, category_id, status, published_at, created_at, updated_at
	`

	var post models.Post
	err := s.db.QueryRow(query, req.Title, slug, req.Content, req.Excerpt,
		authorID, req.CategoryID, req.Status, publishedAt).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content, &post.Excerpt,
		&post.AuthorID, &post.CategoryID, &post.Status, &post.PublishedAt,
		&post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error creando post: %v", err)
		return nil, err
	}

	// Asociar tags si se proporcionan
	if len(req.TagIDs) > 0 {
		err = s.associateTags(post.ID, req.TagIDs)
		if err != nil {
			s.logger.Errorf("Error asociando tags al post: %v", err)
		}
	}

	return &post, nil
}

// UpdatePost actualiza un post existente
func (s *PostService) UpdatePost(id uuid.UUID, req models.PostUpdateRequest) (*models.Post, error) {
	// Verificar que el post existe
	existingPost, err := s.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Title != "" {
		existingPost.Title = req.Title
	}
	if req.Content != "" {
		existingPost.Content = req.Content
	}
	if req.Excerpt != "" {
		existingPost.Excerpt = req.Excerpt
	}
	if req.CategoryID != nil {
		existingPost.CategoryID = *req.CategoryID
	}
	if req.Status != "" {
		existingPost.Status = req.Status
		// Actualizar published_at si se publica
		if req.Status == "published" && existingPost.PublishedAt == nil {
			now := time.Now()
			existingPost.PublishedAt = &now
		}
	}

	query := `
		UPDATE posts 
		SET title = $1, content = $2, excerpt = $3, category_id = $4, status = $5, published_at = $6, updated_at = $7
		WHERE id = $8
		RETURNING id, title, slug, content, excerpt, author_id, category_id, status, published_at, created_at, updated_at
	`

	var post models.Post
	err = s.db.QueryRow(query, existingPost.Title, existingPost.Content, existingPost.Excerpt,
		existingPost.CategoryID, existingPost.Status, existingPost.PublishedAt, time.Now(), id).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content, &post.Excerpt,
		&post.AuthorID, &post.CategoryID, &post.Status, &post.PublishedAt,
		&post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error actualizando post: %v", err)
		return nil, err
	}

	// Actualizar tags si se proporcionan
	if req.TagIDs != nil {
		// Eliminar tags existentes
		err = s.removeAllTags(id)
		if err != nil {
			s.logger.Errorf("Error eliminando tags del post: %v", err)
		}

		// Agregar nuevos tags
		if len(req.TagIDs) > 0 {
			err = s.associateTags(id, req.TagIDs)
			if err != nil {
				s.logger.Errorf("Error asociando tags al post: %v", err)
			}
		}
	}

	return &post, nil
}

// DeletePost elimina un post
func (s *PostService) DeletePost(id uuid.UUID) error {
	query := "DELETE FROM posts WHERE id = $1"

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Errorf("Error eliminando post: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post no encontrado")
	}

	return nil
}

// associateTags asocia tags a un post
func (s *PostService) associateTags(postID uuid.UUID, tagIDs []uuid.UUID) error {
	query := "INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2)"

	for _, tagID := range tagIDs {
		_, err := s.db.Exec(query, postID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

// removeAllTags elimina todas las asociaciones de tags de un post
func (s *PostService) removeAllTags(postID uuid.UUID) error {
	query := "DELETE FROM post_tags WHERE post_id = $1"
	_, err := s.db.Exec(query, postID)
	return err
}
