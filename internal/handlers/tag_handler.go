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

// TagHandler maneja las peticiones HTTP relacionadas con tags
type TagHandler struct {
	tagService   *services.TagService
	statsService *services.StatsService
	logger       *logrus.Logger
}

// NewTagHandler crea una nueva instancia del handler de tags
func NewTagHandler(tagService *services.TagService, statsService *services.StatsService, logger *logrus.Logger) *TagHandler {
	return &TagHandler{
		tagService:   tagService,
		statsService: statsService,
		logger:       logger,
	}
}

// GetTags obtiene todos los tags
func (h *TagHandler) GetTags(c *gin.Context) {
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		h.logger.Errorf("Error obteniendo tags: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}

// GetTag obtiene un tag por su ID
func (h *TagHandler) GetTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de tag inválido",
		})
		return
	}

	tag, err := h.tagService.GetTagByID(tagID)
	if err != nil {
		if err.Error() == "tag no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tag no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo tag: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag": tag,
	})
}

// GetTagBySlug obtiene un tag por su slug
func (h *TagHandler) GetTagBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Slug de tag requerido",
		})
		return
	}

	tag, err := h.tagService.GetTagBySlug(slug)
	if err != nil {
		if err.Error() == "tag no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tag no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo tag por slug: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag": tag,
	})
}

// GetTagWithPosts obtiene un tag con sus posts
func (h *TagHandler) GetTagWithPosts(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de tag inválido",
		})
		return
	}

	tag, err := h.tagService.GetTagWithPosts(tagID)
	if err != nil {
		if err.Error() == "tag no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tag no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo tag con posts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag": tag,
	})
}

// GetPopularTags obtiene los tags más populares
func (h *TagHandler) GetPopularTags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 10000 {
		limit = 10
	}

	tags, err := h.tagService.GetPopularTags(limit)
	if err != nil {
		h.logger.Errorf("Error obteniendo tags populares: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}

// CreateTag crea un nuevo tag
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req models.TagCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	tag, err := h.tagService.CreateTag(req)
	if err != nil {
		if err.Error() == "el slug de tag ya existe" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Errorf("Error creando tag: %v", err)
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
		"tag_created",
		"tag",
		&tag.ID,
		map[string]interface{}{
			"name": tag.Name,
			"slug": tag.Slug,
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusCreated, gin.H{
		"tag":     tag,
		"message": "Tag creado exitosamente",
	})
}

// UpdateTag actualiza un tag existente
func (h *TagHandler) UpdateTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de tag inválido",
		})
		return
	}

	var req models.TagUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	tag, err := h.tagService.UpdateTag(tagID, req)
	if err != nil {
		if err.Error() == "tag no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tag no encontrado",
			})
			return
		}
		h.logger.Errorf("Error actualizando tag: %v", err)
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
		"tag_updated",
		"tag",
		&tagID,
		map[string]interface{}{
			"name": tag.Name,
			"slug": tag.Slug,
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"tag":     tag,
		"message": "Tag actualizado exitosamente",
	})
}

// DeleteTag elimina un tag
func (h *TagHandler) DeleteTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de tag inválido",
		})
		return
	}

	err = h.tagService.DeleteTag(tagID)
	if err != nil {
		if err.Error() == "tag no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tag no encontrado",
			})
			return
		}
		if err.Error() == "no se puede eliminar el tag porque tiene posts asociados" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Errorf("Error eliminando tag: %v", err)
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
		"tag_deleted",
		"tag",
		&tagID,
		map[string]interface{}{
			"tag_id": tagID.String(),
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Tag eliminado exitosamente",
	})
}
