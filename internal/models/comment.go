package models

import (
	"time"

	"github.com/google/uuid"
)

// Comment representa un comentario en un post
type Comment struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	PostID     uuid.UUID  `json:"post_id" db:"post_id"`
	AuthorID   uuid.UUID  `json:"author_id" db:"author_id"`
	ParentID   *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	Content    string     `json:"content" db:"content"`
	IsApproved bool       `json:"is_approved" db:"is_approved"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`

	// Relaciones
	Author  *User     `json:"author,omitempty"`
	Post    *Post     `json:"post,omitempty"`
	Replies []Comment `json:"replies,omitempty"`
}

// CommentCreateRequest representa la solicitud para crear un comentario
type CommentCreateRequest struct {
	PostID   uuid.UUID  `json:"post_id" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
	Content  string     `json:"content" validate:"required,min=1"`
}

// CommentUpdateRequest representa la solicitud para actualizar un comentario
type CommentUpdateRequest struct {
	Content    string `json:"content" validate:"required,min=1"`
	IsApproved *bool  `json:"is_approved"`
}

// CommentListResponse representa la respuesta paginada de comentarios
type CommentListResponse struct {
	Comments   []Comment `json:"comments"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}
