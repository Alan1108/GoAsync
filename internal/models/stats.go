package models

import (
	"time"

	"github.com/google/uuid"
)

// DatabaseStats representa estadísticas generales de la base de datos
type DatabaseStats struct {
	TotalUsers      int64 `json:"total_users"`
	TotalPosts      int64 `json:"total_posts"`
	TotalComments   int64 `json:"total_comments"`
	TotalCategories int64 `json:"total_categories"`
	TotalTags       int64 `json:"total_tags"`
}

// ActivityLog representa un log de actividad del sistema
type ActivityLog struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	UserID       *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	Action       string                 `json:"action" db:"action"`
	ResourceType string                 `json:"resource_type" db:"resource_type"`
	ResourceID   *uuid.UUID             `json:"resource_id,omitempty" db:"resource_id"`
	Details      map[string]interface{} `json:"details" db:"details"`
	IPAddress    string                 `json:"ip_address" db:"ip_address"`
	UserAgent    string                 `json:"user_agent" db:"user_agent"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`

	// Relaciones
	User *User `json:"user,omitempty"`
}

// ActivityLogFilter representa los filtros para listar logs de actividad
type ActivityLogFilter struct {
	UserID       uuid.UUID `json:"user_id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Page         int       `json:"page"`
	PerPage      int       `json:"per_page"`
}

// PostStats representa estadísticas de un post específico
type PostStats struct {
	PostID       uuid.UUID `json:"post_id" db:"post_id"`
	Title        string    `json:"title" db:"title"`
	CommentCount int       `json:"comment_count" db:"comment_count"`
	ViewCount    int       `json:"view_count" db:"view_count"`
	PublishedAt  time.Time `json:"published_at" db:"published_at"`
}
