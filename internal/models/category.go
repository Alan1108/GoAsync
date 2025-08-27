package models

import (
	"time"

	"github.com/google/uuid"
)

// Category representa una categoría de posts
type Category struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Slug        string    `json:"slug" db:"slug"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Relaciones
	Posts []Post `json:"posts,omitempty"`
}

// CategoryCreateRequest representa la solicitud para crear una categoría
type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
	Slug        string `json:"slug" validate:"omitempty,min=1,max=100"`
}

// CategoryUpdateRequest representa la solicitud para actualizar una categoría
type CategoryUpdateRequest struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=100"`
	Description string `json:"description"`
	Slug        string `json:"slug" validate:"omitempty,min=1,max=100"`
	IsActive    *bool  `json:"is_active"`
}
