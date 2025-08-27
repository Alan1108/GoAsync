package models

import (
	"time"

	"github.com/google/uuid"
)

// Post representa un artículo o post en el sistema
type Post struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Slug        string     `json:"slug" db:"slug"`
	Content     string     `json:"content" db:"content"`
	Excerpt     string     `json:"excerpt" db:"excerpt"`
	AuthorID    uuid.UUID  `json:"author_id" db:"author_id"`
	CategoryID  uuid.UUID  `json:"category_id" db:"category_id"`
	Status      string     `json:"status" db:"status"`
	PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Relaciones
	Author   *User     `json:"author,omitempty"`
	Category *Category `json:"category,omitempty"`
	Tags     []Tag     `json:"tags,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
}

// PostCreateRequest representa la solicitud para crear un post
type PostCreateRequest struct {
	Title      string      `json:"title" validate:"required,min=1,max=255"`
	Content    string      `json:"content" validate:"required"`
	Excerpt    string      `json:"excerpt"`
	CategoryID uuid.UUID   `json:"category_id" validate:"required"`
	Status     string      `json:"status" validate:"omitempty,oneof=draft published archived"`
	TagIDs     []uuid.UUID `json:"tag_ids"`
}

// PostUpdateRequest representa la solicitud para actualizar un post
type PostUpdateRequest struct {
	Title      string      `json:"title" validate:"omitempty,min=1,max=255"`
	Content    string      `json:"content"`
	Excerpt    string      `json:"excerpt"`
	CategoryID *uuid.UUID  `json:"category_id"`
	Status     string      `json:"status" validate:"omitempty,oneof=draft published archived"`
	TagIDs     []uuid.UUID `json:"tag_ids"`
}

// PostListResponse representa la respuesta paginada de posts
type PostListResponse struct {
	Posts      []Post `json:"posts"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
}

// PostFilter representa los filtros para listar posts
type PostFilter struct {
	Status     string    `json:"status"`
	CategoryID uuid.UUID `json:"category_id"`
	AuthorID   uuid.UUID `json:"author_id"`
	TagID      uuid.UUID `json:"tag_id"`
	Search     string    `json:"search"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
}
