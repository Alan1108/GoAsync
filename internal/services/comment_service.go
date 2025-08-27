package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CommentService maneja la lógica de negocio para comentarios
type CommentService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewCommentService crea una nueva instancia del servicio de comentarios
func NewCommentService(db *sql.DB, logger *logrus.Logger) *CommentService {
	return &CommentService{
		db:     db,
		logger: logger,
	}
}

// GetCommentsByPostID obtiene comentarios de un post con paginación
func (s *CommentService) GetCommentsByPostID(postID uuid.UUID, page, perPage int) (*models.CommentListResponse, error) {
	offset := (page - 1) * perPage

	// Obtener total de comentarios
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = $1 AND is_approved = true", postID).Scan(&total)
	if err != nil {
		s.logger.Errorf("Error contando comentarios: %v", err)
		return nil, err
	}

	// Obtener comentarios principales (sin parent_id)
	query := `
		SELECT c.id, c.post_id, c.author_id, c.parent_id, c.content, c.is_approved, 
		       c.created_at, c.updated_at,
		       u.username as author_username, u.first_name as author_first_name, u.last_name as author_last_name
		FROM comments c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.post_id = $1 AND c.parent_id IS NULL AND c.is_approved = true
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.Query(query, postID, perPage, offset)
	if err != nil {
		s.logger.Errorf("Error obteniendo comentarios: %v", err)
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var authorUsername, authorFirstName, authorLastName sql.NullString

		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.AuthorID, &comment.ParentID,
			&comment.Content, &comment.IsApproved, &comment.CreatedAt, &comment.UpdatedAt,
			&authorUsername, &authorFirstName, &authorLastName,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando comentario: %v", err)
			continue
		}

		// Construir autor
		if authorUsername.Valid {
			comment.Author = &models.User{
				ID:        comment.AuthorID,
				Username:  authorUsername.String,
				FirstName: authorFirstName.String,
				LastName:  authorLastName.String,
			}
		}

		// Obtener respuestas del comentario
		replies, err := s.getCommentReplies(comment.ID)
		if err != nil {
			s.logger.Errorf("Error obteniendo respuestas del comentario: %v", err)
		} else {
			comment.Replies = replies
		}

		comments = append(comments, comment)
	}

	// Calcular total de páginas
	totalPages := (total + perPage - 1) / perPage

	return &models.CommentListResponse{
		Comments:   comments,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetCommentByID obtiene un comentario por su ID
func (s *CommentService) GetCommentByID(id uuid.UUID) (*models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.author_id, c.parent_id, c.content, c.is_approved, 
		       c.created_at, c.updated_at,
		       u.username as author_username, u.first_name as author_first_name, u.last_name as author_last_name
		FROM comments c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.id = $1
	`

	var comment models.Comment
	var authorUsername, authorFirstName, authorLastName sql.NullString

	err := s.db.QueryRow(query, id).Scan(
		&comment.ID, &comment.PostID, &comment.AuthorID, &comment.ParentID,
		&comment.Content, &comment.IsApproved, &comment.CreatedAt, &comment.UpdatedAt,
		&authorUsername, &authorFirstName, &authorLastName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comentario no encontrado")
		}
		s.logger.Errorf("Error obteniendo comentario por ID: %v", err)
		return nil, err
	}

	// Construir autor
	if authorUsername.Valid {
		comment.Author = &models.User{
			ID:        comment.AuthorID,
			Username:  authorUsername.String,
			FirstName: authorFirstName.String,
			LastName:  authorLastName.String,
		}
	}

	return &comment, nil
}

// GetAllComments obtiene todos los comentarios con filtros
func (s *CommentService) GetAllComments(page, perPage int, approvedOnly bool) (*models.CommentListResponse, error) {
	offset := (page - 1) * perPage

	// Construir query base
	baseQuery := `
		SELECT c.id, c.post_id, c.author_id, c.parent_id, c.content, c.is_approved, 
		       c.created_at, c.updated_at,
		       u.username as author_username, u.first_name as author_first_name, u.last_name as author_last_name,
		       p.title as post_title, p.slug as post_slug
		FROM comments c
		LEFT JOIN users u ON c.author_id = u.id
		LEFT JOIN posts p ON c.post_id = p.id
	`

	whereClause := ""
	if approvedOnly {
		whereClause = "WHERE c.is_approved = true"
	}

	// Query para contar total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM comments c %s", whereClause)
	var total int
	err := s.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		s.logger.Errorf("Error contando comentarios: %v", err)
		return nil, err
	}

	// Query para obtener comentarios
	query := fmt.Sprintf(`
		%s
		%s
		ORDER BY c.created_at DESC
		LIMIT $1 OFFSET $2
	`, baseQuery, whereClause)

	rows, err := s.db.Query(query, perPage, offset)
	if err != nil {
		s.logger.Errorf("Error obteniendo comentarios: %v", err)
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var authorUsername, authorFirstName, authorLastName sql.NullString
		var postTitle, postSlug sql.NullString

		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.AuthorID, &comment.ParentID,
			&comment.Content, &comment.IsApproved, &comment.CreatedAt, &comment.UpdatedAt,
			&authorUsername, &authorFirstName, &authorLastName,
			&postTitle, &postSlug,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando comentario: %v", err)
			continue
		}

		// Construir autor
		if authorUsername.Valid {
			comment.Author = &models.User{
				ID:        comment.AuthorID,
				Username:  authorUsername.String,
				FirstName: authorFirstName.String,
				LastName:  authorLastName.String,
			}
		}

		// Construir post
		if postTitle.Valid {
			comment.Post = &models.Post{
				ID:    comment.PostID,
				Title: postTitle.String,
				Slug:  postSlug.String,
			}
		}

		comments = append(comments, comment)
	}

	// Calcular total de páginas
	totalPages := (total + perPage - 1) / perPage

	return &models.CommentListResponse{
		Comments:   comments,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// CreateComment crea un nuevo comentario
func (s *CommentService) CreateComment(req models.CommentCreateRequest, authorID uuid.UUID) (*models.Comment, error) {
	// Verificar que el post existe y está publicado
	var postStatus string
	err := s.db.QueryRow("SELECT status FROM posts WHERE id = $1", req.PostID).Scan(&postStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post no encontrado")
		}
		s.logger.Errorf("Error verificando post: %v", err)
		return nil, err
	}

	if postStatus != "published" {
		return nil, fmt.Errorf("no se pueden agregar comentarios a posts no publicados")
	}

	// Verificar parent_id si se proporciona
	if req.ParentID != nil {
		var parentExists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1 AND post_id = $2)",
			req.ParentID, req.PostID).Scan(&parentExists)
		if err != nil {
			s.logger.Errorf("Error verificando comentario padre: %v", err)
			return nil, err
		}
		if !parentExists {
			return nil, fmt.Errorf("comentario padre no encontrado")
		}
	}

	query := `
		INSERT INTO comments (post_id, author_id, parent_id, content, is_approved)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, post_id, author_id, parent_id, content, is_approved, created_at, updated_at
	`

	// Por defecto, los comentarios no están aprobados
	isApproved := false

	var comment models.Comment
	err = s.db.QueryRow(query, req.PostID, authorID, req.ParentID, req.Content, isApproved).Scan(
		&comment.ID, &comment.PostID, &comment.AuthorID, &comment.ParentID,
		&comment.Content, &comment.IsApproved, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error creando comentario: %v", err)
		return nil, err
	}

	return &comment, nil
}

// UpdateComment actualiza un comentario existente
func (s *CommentService) UpdateComment(id uuid.UUID, req models.CommentUpdateRequest) (*models.Comment, error) {
	// Verificar que el comentario existe
	existingComment, err := s.GetCommentByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Content != "" {
		existingComment.Content = req.Content
	}
	if req.IsApproved != nil {
		existingComment.IsApproved = *req.IsApproved
	}

	query := `
		UPDATE comments 
		SET content = $1, is_approved = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, post_id, author_id, parent_id, content, is_approved, created_at, updated_at
	`

	var comment models.Comment
	err = s.db.QueryRow(query, existingComment.Content, existingComment.IsApproved, time.Now(), id).Scan(
		&comment.ID, &comment.PostID, &comment.AuthorID, &comment.ParentID,
		&comment.Content, &comment.IsApproved, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error actualizando comentario: %v", err)
		return nil, err
	}

	return &comment, nil
}

// DeleteComment elimina un comentario
func (s *CommentService) DeleteComment(id uuid.UUID) error {
	query := "DELETE FROM comments WHERE id = $1"

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Errorf("Error eliminando comentario: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comentario no encontrado")
	}

	return nil
}

// ApproveComment aprueba un comentario
func (s *CommentService) ApproveComment(id uuid.UUID) error {
	query := "UPDATE comments SET is_approved = true, updated_at = $1 WHERE id = $2"

	result, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		s.logger.Errorf("Error aprobando comentario: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comentario no encontrado")
	}

	return nil
}

// getCommentReplies obtiene las respuestas de un comentario
func (s *CommentService) getCommentReplies(commentID uuid.UUID) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.author_id, c.parent_id, c.content, c.is_approved, 
		       c.created_at, c.updated_at,
		       u.username as author_username, u.first_name as author_first_name, u.last_name as author_last_name
		FROM comments c
		LEFT JOIN users u ON c.author_id = u.id
		WHERE c.parent_id = $1 AND c.is_approved = true
		ORDER BY c.created_at ASC
	`

	rows, err := s.db.Query(query, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []models.Comment
	for rows.Next() {
		var reply models.Comment
		var authorUsername, authorFirstName, authorLastName sql.NullString

		err := rows.Scan(
			&reply.ID, &reply.PostID, &reply.AuthorID, &reply.ParentID,
			&reply.Content, &reply.IsApproved, &reply.CreatedAt, &reply.UpdatedAt,
			&authorUsername, &authorFirstName, &authorLastName,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando respuesta: %v", err)
			continue
		}

		// Construir autor
		if authorUsername.Valid {
			reply.Author = &models.User{
				ID:        reply.AuthorID,
				Username:  authorUsername.String,
				FirstName: authorFirstName.String,
				LastName:  authorLastName.String,
			}
		}

		replies = append(replies, reply)
	}

	return replies, nil
}
