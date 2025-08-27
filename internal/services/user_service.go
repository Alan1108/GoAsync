package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// UserService maneja la lógica de negocio para usuarios
type UserService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(db *sql.DB, logger *logrus.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

// GetAllUsers obtiene todos los usuarios con paginación
func (s *UserService) GetAllUsers(page, perPage int) ([]models.User, int, error) {
	offset := (page - 1) * perPage

	// Obtener total de usuarios
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		s.logger.Errorf("Error contando usuarios: %v", err)
		return nil, 0, err
	}

	// Obtener usuarios
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, 
		       is_active, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, perPage, offset)
	if err != nil {
		s.logger.Errorf("Error obteniendo usuarios: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.FirstName, &user.LastName, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando usuario: %v", err)
			continue
		}
		users = append(users, user)
	}

	return users, total, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, 
		       is_active, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user models.User
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		s.logger.Errorf("Error obteniendo usuario por ID: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername obtiene un usuario por su nombre de usuario
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, 
		       is_active, created_at, updated_at
		FROM users 
		WHERE username = $1
	`

	var user models.User
	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		s.logger.Errorf("Error obteniendo usuario por username: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail obtiene un usuario por su email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, 
		       is_active, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	var user models.User
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		s.logger.Errorf("Error obteniendo usuario por email: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetUserWithProfile obtiene un usuario con su perfil
func (s *UserService) GetUserWithProfile(id uuid.UUID) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Obtener perfil
	profileQuery := `
		SELECT id, user_id, bio, avatar_url, date_of_birth, phone, address, 
		       created_at, updated_at
		FROM user_profiles 
		WHERE user_id = $1
	`

	var profile models.UserProfile
	err = s.db.QueryRow(profileQuery, id).Scan(
		&profile.ID, &profile.UserID, &profile.Bio, &profile.AvatarURL,
		&profile.DateOfBirth, &profile.Phone, &profile.Address,
		&profile.CreatedAt, &profile.UpdatedAt,
	)

	if err != nil && err != sql.ErrNoRows {
		s.logger.Errorf("Error obteniendo perfil de usuario: %v", err)
		return nil, err
	}

	if err != sql.ErrNoRows {
		user.Profile = &profile
	}

	return user, nil
}

// GetUserStats obtiene estadísticas de un usuario
func (s *UserService) GetUserStats(id uuid.UUID) (*models.UserStats, error) {
	query := `
		SELECT user_id, username, email, posts_count, comments_count, joined_at
		FROM user_stats 
		WHERE user_id = $1
	`

	var stats models.UserStats
	err := s.db.QueryRow(query, id).Scan(
		&stats.UserID, &stats.Username, &stats.Email,
		&stats.PostsCount, &stats.CommentsCount, &stats.JoinedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("estadísticas de usuario no encontradas")
		}
		s.logger.Errorf("Error obteniendo estadísticas de usuario: %v", err)
		return nil, err
	}

	return &stats, nil
}

// GetAllUserStats obtiene estadísticas de todos los usuarios
func (s *UserService) GetAllUserStats() ([]models.UserStats, error) {
	query := `
		SELECT user_id, username, email, posts_count, comments_count, joined_at
		FROM user_stats 
		ORDER BY posts_count DESC, comments_count DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Errorf("Error obteniendo estadísticas de usuarios: %v", err)
		return nil, err
	}
	defer rows.Close()

	var stats []models.UserStats
	for rows.Next() {
		var stat models.UserStats
		err := rows.Scan(
			&stat.UserID, &stat.Username, &stat.Email,
			&stat.PostsCount, &stat.CommentsCount, &stat.JoinedAt,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando estadísticas de usuario: %v", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(req models.UserCreateRequest) (*models.User, error) {
	// Verificar si el username ya existe
	existingUser, _ := s.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("el nombre de usuario ya existe")
	}

	// Verificar si el email ya existe
	existingUser, _ = s.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("el email ya existe")
	}

	// Hash de la contraseña (aquí deberías usar bcrypt)
	passwordHash := req.Password // En producción, usar bcrypt

	query := `
		INSERT INTO users (username, email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, password_hash, first_name, last_name, 
		          is_active, created_at, updated_at
	`

	var user models.User
	err := s.db.QueryRow(query, req.Username, req.Email, passwordHash, req.FirstName, req.LastName).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error creando usuario: %v", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUser actualiza un usuario existente
func (s *UserService) UpdateUser(id uuid.UUID, req models.UserUpdateRequest) (*models.User, error) {
	// Verificar que el usuario existe
	existingUser, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.FirstName != "" {
		existingUser.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existingUser.LastName = req.LastName
	}
	if req.IsActive != nil {
		existingUser.IsActive = *req.IsActive
	}

	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, is_active = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, username, email, password_hash, first_name, last_name, 
		          is_active, created_at, updated_at
	`

	var user models.User
	err = s.db.QueryRow(query, existingUser.FirstName, existingUser.LastName,
		existingUser.IsActive, time.Now(), id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		s.logger.Errorf("Error actualizando usuario: %v", err)
		return nil, err
	}

	return &user, nil
}

// DeleteUser elimina un usuario
func (s *UserService) DeleteUser(id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Errorf("Error eliminando usuario: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuario no encontrado")
	}

	return nil
}
