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

// CommentHandler maneja las peticiones HTTP relacionadas con comentarios
type CommentHandler struct {
	commentService *services.CommentService
	statsService   *services.StatsService
	logger         *logrus.Logger
}

// NewCommentHandler crea una nueva instancia del handler de comentarios
func NewCommentHandler(commentService *services.CommentService, statsService *services.StatsService, logger *logrus.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		statsService:   statsService,
		logger:         logger,
	}
}

// GetComments obtiene comentarios de un post
func (h *CommentHandler) GetComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de post inválido",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 10000 {
		perPage = 10
	}

	response, err := h.commentService.GetCommentsByPostID(postID, page, perPage)
	if err != nil {
		h.logger.Errorf("Error obteniendo comentarios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": response.Comments,
		"pagination": gin.H{
			"page":        response.Page,
			"per_page":    response.PerPage,
			"total":       response.Total,
			"total_pages": response.TotalPages,
		},
	})
}

// GetAllComments obtiene todos los comentarios
func (h *CommentHandler) GetAllComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	approvedOnly := c.DefaultQuery("approved_only", "true") == "true"

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 10000 {
		perPage = 10
	}

	response, err := h.commentService.GetAllComments(page, perPage, approvedOnly)
	if err != nil {
		h.logger.Errorf("Error obteniendo comentarios: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": response.Comments,
		"pagination": gin.H{
			"page":        response.Page,
			"per_page":    response.PerPage,
			"total":       response.Total,
			"total_pages": response.TotalPages,
		},
	})
}

// GetComment obtiene un comentario por su ID
func (h *CommentHandler) GetComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de comentario inválido",
		})
		return
	}

	comment, err := h.commentService.GetCommentByID(commentID)
	if err != nil {
		if err.Error() == "comentario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Comentario no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo comentario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

// CreateComment crea un nuevo comentario
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req models.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	// TODO: Obtener el ID del usuario autenticado
	authorID := uuid.New() // En producción, esto vendría del contexto de autenticación

	comment, err := h.commentService.CreateComment(req, authorID)
	if err != nil {
		if err.Error() == "post no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post no encontrado",
			})
			return
		}
		if err.Error() == "no se pueden agregar comentarios a posts no publicados" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err.Error() == "comentario padre no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Errorf("Error creando comentario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&authorID,
		"comment_added",
		"comment",
		&comment.ID,
		map[string]interface{}{
			"post_id": comment.PostID.String(),
			"content": comment.Content[:50] + "...", // Solo los primeros 50 caracteres
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusCreated, gin.H{
		"comment": comment,
		"message": "Comentario creado exitosamente",
	})
}

// UpdateComment actualiza un comentario existente
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de comentario inválido",
		})
		return
	}

	var req models.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	comment, err := h.commentService.UpdateComment(commentID, req)
	if err != nil {
		if err.Error() == "comentario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Comentario no encontrado",
			})
			return
		}
		h.logger.Errorf("Error actualizando comentario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// TODO: Obtener el ID del usuario autenticado
	userID := uuid.New() // En producción, esto vendría del contexto de autenticación

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&userID,
		"comment_updated",
		"comment",
		&commentID,
		map[string]interface{}{
			"post_id": comment.PostID.String(),
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"comment": comment,
		"message": "Comentario actualizado exitosamente",
	})
}

// DeleteComment elimina un comentario
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de comentario inválido",
		})
		return
	}

	err = h.commentService.DeleteComment(commentID)
	if err != nil {
		if err.Error() == "comentario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Comentario no encontrado",
			})
			return
		}
		h.logger.Errorf("Error eliminando comentario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// TODO: Obtener el ID del usuario autenticado
	userID := uuid.New() // En producción, esto vendría del contexto de autenticación

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&userID,
		"comment_deleted",
		"comment",
		&commentID,
		map[string]interface{}{
			"comment_id": commentID.String(),
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comentario eliminado exitosamente",
	})
}

// ApproveComment aprueba un comentario
func (h *CommentHandler) ApproveComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de comentario inválido",
		})
		return
	}

	err = h.commentService.ApproveComment(commentID)
	if err != nil {
		if err.Error() == "comentario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Comentario no encontrado",
			})
			return
		}
		h.logger.Errorf("Error aprobando comentario: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// TODO: Obtener el ID del usuario autenticado
	userID := uuid.New() // En producción, esto vendría del contexto de autenticación

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&userID,
		"comment_approved",
		"comment",
		&commentID,
		map[string]interface{}{
			"comment_id": commentID.String(),
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comentario aprobado exitosamente",
	})
}
