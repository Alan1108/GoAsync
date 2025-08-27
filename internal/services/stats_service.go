package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// StatsService maneja la lógica de negocio para estadísticas
type StatsService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewStatsService crea una nueva instancia del servicio de estadísticas
func NewStatsService(db *sql.DB, logger *logrus.Logger) *StatsService {
	return &StatsService{
		db:     db,
		logger: logger,
	}
}

// GetDatabaseStats obtiene estadísticas generales de la base de datos
func (s *StatsService) GetDatabaseStats() (*models.DatabaseStats, error) {
	query := "SELECT * FROM get_database_stats()"

	var stats models.DatabaseStats
	err := s.db.QueryRow(query).Scan(
		&stats.TotalUsers, &stats.TotalPosts, &stats.TotalComments,
		&stats.TotalCategories, &stats.TotalTags,
	)

	if err != nil {
		s.logger.Errorf("Error obteniendo estadísticas de la base de datos: %v", err)
		return nil, err
	}

	return &stats, nil
}

// GetActivityLogs obtiene logs de actividad con filtros
func (s *StatsService) GetActivityLogs(filter models.ActivityLogFilter) ([]models.ActivityLog, error) {
	// Construir query base
	baseQuery := `
		SELECT al.id, al.user_id, al.action, al.resource_type, al.resource_id, 
		       al.details, al.ip_address, al.user_agent, al.created_at,
		       u.username as user_username, u.first_name as user_first_name, u.last_name as user_last_name
		FROM activity_logs al
		LEFT JOIN users u ON al.user_id = u.id
	`

	whereConditions := []string{}
	args := []interface{}{}
	argCount := 0

	// Aplicar filtros
	if filter.UserID != uuid.Nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("al.user_id = $%d", argCount))
		args = append(args, filter.UserID)
	}

	if filter.Action != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("al.action = $%d", argCount))
		args = append(args, filter.Action)
	}

	if filter.ResourceType != "" {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("al.resource_type = $%d", argCount))
		args = append(args, filter.ResourceType)
	}

	if !filter.StartDate.IsZero() {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("al.created_at >= $%d", argCount))
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("al.created_at <= $%d", argCount))
		args = append(args, filter.EndDate)
	}

	// Construir WHERE clause
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + whereConditions[0]
		for i := 1; i < len(whereConditions); i++ {
			whereClause += " AND " + whereConditions[i]
		}
	}

	// Agregar paginación
	argCount++
	limitArg := fmt.Sprintf("$%d", argCount)
	argCount++
	offsetArg := fmt.Sprintf("$%d", argCount)
	args = append(args, filter.PerPage, (filter.Page-1)*filter.PerPage)

	query := fmt.Sprintf(`
		%s
		%s
		ORDER BY al.created_at DESC
		LIMIT %s OFFSET %s
	`, baseQuery, whereClause, limitArg, offsetArg)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		s.logger.Errorf("Error obteniendo logs de actividad: %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var log models.ActivityLog
		var userUsername, userFirstName, userLastName sql.NullString

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
			&log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
			&userUsername, &userFirstName, &userLastName,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando log de actividad: %v", err)
			continue
		}

		// Construir usuario si existe
		if userUsername.Valid {
			log.User = &models.User{
				ID:        *log.UserID,
				Username:  userUsername.String,
				FirstName: userFirstName.String,
				LastName:  userLastName.String,
			}
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// CreateActivityLog crea un nuevo log de actividad
func (s *StatsService) CreateActivityLog(userID *uuid.UUID, action, resourceType string, resourceID *uuid.UUID, details map[string]interface{}, ipAddress, userAgent string) error {
	query := `
		INSERT INTO activity_logs (user_id, action, resource_type, resource_id, details, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.Exec(query, userID, action, resourceType, resourceID, details, ipAddress, userAgent)
	if err != nil {
		s.logger.Errorf("Error creando log de actividad: %v", err)
		return err
	}

	return nil
}

// GetPostStats obtiene estadísticas de posts
func (s *StatsService) GetPostStats() ([]models.PostStats, error) {
	query := `
		SELECT p.id as post_id, p.title, COUNT(c.id) as comment_count, 
		       0 as view_count, p.published_at
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id AND c.is_approved = true
		WHERE p.status = 'published'
		GROUP BY p.id, p.title, p.published_at
		ORDER BY p.published_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Errorf("Error obteniendo estadísticas de posts: %v", err)
		return nil, err
	}
	defer rows.Close()

	var stats []models.PostStats
	for rows.Next() {
		var stat models.PostStats
		err := rows.Scan(&stat.PostID, &stat.Title, &stat.CommentCount, &stat.ViewCount, &stat.PublishedAt)
		if err != nil {
			s.logger.Errorf("Error escaneando estadísticas de post: %v", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// GetRecentActivity obtiene actividad reciente
func (s *StatsService) GetRecentActivity(limit int) ([]models.ActivityLog, error) {
	query := `
		SELECT al.id, al.user_id, al.action, al.resource_type, al.resource_id, 
		       al.details, al.ip_address, al.user_agent, al.created_at,
		       u.username as user_username, u.first_name as user_first_name, u.last_name as user_last_name
		FROM activity_logs al
		LEFT JOIN users u ON al.user_id = u.id
		ORDER BY al.created_at DESC
		LIMIT $1
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		s.logger.Errorf("Error obteniendo actividad reciente: %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var log models.ActivityLog
		var userUsername, userFirstName, userLastName sql.NullString

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
			&log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
			&userUsername, &userFirstName, &userLastName,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando log de actividad: %v", err)
			continue
		}

		// Construir usuario si existe
		if userUsername.Valid {
			log.User = &models.User{
				ID:        *log.UserID,
				Username:  userUsername.String,
				FirstName: userFirstName.String,
				LastName:  userLastName.String,
			}
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// GetUserActivity obtiene actividad de un usuario específico
func (s *StatsService) GetUserActivity(userID uuid.UUID, limit int) ([]models.ActivityLog, error) {
	query := `
		SELECT al.id, al.user_id, al.action, al.resource_type, al.resource_id, 
		       al.details, al.ip_address, al.user_agent, al.created_at,
		       u.username as user_username, u.first_name as user_first_name, u.last_name as user_last_name
		FROM activity_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.user_id = $1
		ORDER BY al.created_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(query, userID, limit)
	if err != nil {
		s.logger.Errorf("Error obteniendo actividad del usuario: %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var log models.ActivityLog
		var userUsername, userFirstName, userLastName sql.NullString

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.ResourceType, &log.ResourceID,
			&log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
			&userUsername, &userFirstName, &userLastName,
		)
		if err != nil {
			s.logger.Errorf("Error escaneando log de actividad: %v", err)
			continue
		}

		// Construir usuario si existe
		if userUsername.Valid {
			log.User = &models.User{
				ID:        *log.UserID,
				Username:  userUsername.String,
				FirstName: userFirstName.String,
				LastName:  userLastName.String,
			}
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// GetDailyStats obtiene estadísticas diarias
func (s *StatsService) GetDailyStats(days int) (map[string]interface{}, error) {
	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as total_activities,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(CASE WHEN action = 'user_login' THEN 1 END) as logins,
			COUNT(CASE WHEN action = 'post_created' THEN 1 END) as posts_created,
			COUNT(CASE WHEN action = 'comment_added' THEN 1 END) as comments_added
		FROM activity_logs
		WHERE created_at >= CURRENT_DATE - INTERVAL '$1 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`

	rows, err := s.db.Query(query, days)
	if err != nil {
		s.logger.Errorf("Error obteniendo estadísticas diarias: %v", err)
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]interface{})
	var dailyStats []map[string]interface{}

	for rows.Next() {
		var date time.Time
		var totalActivities, uniqueUsers, logins, postsCreated, commentsAdded int

		err := rows.Scan(&date, &totalActivities, &uniqueUsers, &logins, &postsCreated, &commentsAdded)
		if err != nil {
			s.logger.Errorf("Error escaneando estadísticas diarias: %v", err)
			continue
		}

		dailyStats = append(dailyStats, map[string]interface{}{
			"date":             date.Format("2006-01-02"),
			"total_activities": totalActivities,
			"unique_users":     uniqueUsers,
			"logins":           logins,
			"posts_created":    postsCreated,
			"comments_added":   commentsAdded,
		})
	}

	stats["daily_stats"] = dailyStats
	stats["period_days"] = days

	return stats, nil
}
