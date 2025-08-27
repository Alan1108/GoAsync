package handlers

import (
	"net/http"
	"strconv"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/alan.bermudez/goasync/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// UserHandler maneja las peticiones HTTP relacionadas con usuarios
type UserHandler struct {
	userService  *services.UserService
	statsService *services.StatsService
	logger       *logrus.Logger
}

// NewUserHandler crea una nueva instancia del handler de usuarios
func NewUserHandler(userService *services.UserService, statsService *services.StatsService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService:  userService,
		statsService: statsService,
		logger:       logger,
	}
}

// GetUsers obtiene todos los usuarios con paginación
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Obtener parámetros de paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// Validar parámetros
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 10000 {
		perPage = 10
	}

	users, total, err := h.userService.GetAllUsers(page, perPage)
	if err != nil {
		h.logger.Errorf("Error obteniendo usuarios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// Calcular total de páginas
	totalPages := (total + perPage - 1) / perPage

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetUser obtiene un usuario por su ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuario no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetUserWithProfile obtiene un usuario con su perfil
func (h *UserHandler) GetUserWithProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	user, err := h.userService.GetUserWithProfile(userID)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuario no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo usuario con perfil: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetUserStats obtiene estadísticas de un usuario
func (h *UserHandler) GetUserStats(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	stats, err := h.userService.GetUserStats(userID)
	if err != nil {
		if err.Error() == "estadísticas de usuario no encontradas" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Estadísticas de usuario no encontradas",
			})
			return
		}
		h.logger.Errorf("Error obteniendo estadísticas de usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// GetAllUserStats obtiene estadísticas de todos los usuarios
func (h *UserHandler) GetAllUserStats(c *gin.Context) {
	stats, err := h.userService.GetAllUserStats()
	if err != nil {
		h.logger.Errorf("Error obteniendo estadísticas de usuarios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// CreateUser crea un nuevo usuario
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		if err.Error() == "el nombre de usuario ya existe" || err.Error() == "el email ya existe" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Errorf("Error creando usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&user.ID,
		"user_created",
		"user",
		&user.ID,
		map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"message": "Usuario creado exitosamente",
	})
}

// UpdateUser actualiza un usuario existente
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	user, err := h.userService.UpdateUser(userID, req)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuario no encontrado",
			})
			return
		}
		h.logger.Errorf("Error actualizando usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&userID,
		"user_updated",
		"user",
		&userID,
		map[string]interface{}{
			"username": user.Username,
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"message": "Usuario actualizado exitosamente",
	})
}

// DeleteUser elimina un usuario
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	err = h.userService.DeleteUser(userID)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuario no encontrado",
			})
			return
		}
		h.logger.Errorf("Error eliminando usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		nil,
		"user_deleted",
		"user",
		&userID,
		map[string]interface{}{
			"user_id": userID.String(),
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario eliminado exitosamente",
	})
}

// GetUserActivity obtiene la actividad de un usuario
func (h *UserHandler) GetUserActivity(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de usuario inválido",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 10000 {
		limit = 10
	}

	activity, err := h.statsService.GetUserActivity(userID, limit)
	if err != nil {
		h.logger.Errorf("Error obteniendo actividad del usuario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activity": activity,
	})
}
