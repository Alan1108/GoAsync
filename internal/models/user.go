package models

import (
	"time"

	"github.com/google/uuid"
)

// User representa un usuario en el sistema
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Relaciones
	Profile *UserProfile `json:"profile,omitempty"`
	Posts   []Post       `json:"posts,omitempty"`
}

// UserProfile representa el perfil de un usuario
type UserProfile struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Bio         string     `json:"bio" db:"bio"`
	AvatarURL   string     `json:"avatar_url" db:"avatar_url"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Phone       string     `json:"phone" db:"phone"`
	Address     string     `json:"address" db:"address"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// UserStats representa estadísticas de un usuario
type UserStats struct {
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Username      string    `json:"username" db:"username"`
	Email         string    `json:"email" db:"email"`
	PostsCount    int       `json:"posts_count" db:"posts_count"`
	CommentsCount int       `json:"comments_count" db:"comments_count"`
	JoinedAt      time.Time `json:"joined_at" db:"joined_at"`
}

// UserCreateRequest representa la solicitud para crear un usuario
type UserCreateRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
}

// UserUpdateRequest representa la solicitud para actualizar un usuario
type UserUpdateRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=1,max=100"`
	LastName  string `json:"last_name" validate:"omitempty,min=1,max=100"`
	IsActive  *bool  `json:"is_active"`
}

// UserProfileUpdateRequest representa la solicitud para actualizar un perfil
type UserProfileUpdateRequest struct {
	Bio         string     `json:"bio"`
	AvatarURL   string     `json:"avatar_url"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Phone       string     `json:"phone"`
	Address     string     `json:"address"`
}
