package models

import (
	"time"

	"github.com/google/uuid"
)

// Tag representa una etiqueta para posts
type Tag struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	// Relaciones
	Posts []Post `json:"posts,omitempty"`
}

// TagCreateRequest representa la solicitud para crear una etiqueta
type TagCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Slug        string `json:"slug" validate:"omitempty,min=1,max=50"`
	Description string `json:"description"`
}

// TagUpdateRequest representa la solicitud para actualizar una etiqueta
type TagUpdateRequest struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=50"`
	Slug        string `json:"slug" validate:"omitempty,min=1,max=50"`
	Description string `json:"description"`
}
