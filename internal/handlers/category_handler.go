package handlers

import (
	"net/http"

	"github.com/alan.bermudez/goasync/internal/models"
	"github.com/alan.bermudez/goasync/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CategoryHandler maneja las peticiones HTTP relacionadas con categorías
type CategoryHandler struct {
	categoryService *services.CategoryService
	statsService    *services.StatsService
	logger          *logrus.Logger
}

// NewCategoryHandler crea una nueva instancia del handler de categorías
func NewCategoryHandler(categoryService *services.CategoryService, statsService *services.StatsService, logger *logrus.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		statsService:    statsService,
		logger:          logger,
	}
}

// GetCategories obtiene todas las categorías
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		h.logger.Errorf("Error obteniendo categorías: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// GetCategory obtiene una categoría por su ID
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de categoría inválido",
		})
		return
	}

	category, err := h.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		if err.Error() == "categoría no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Categoría no encontrada",
			})
			return
		}
		h.logger.Errorf("Error obteniendo categoría: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

// GetCategoryBySlug obtiene una categoría por su slug
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Slug de categoría requerido",
		})
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(slug)
	if err != nil {
		if err.Error() == "categoría no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Categoría no encontrada",
			})
			return
		}
		h.logger.Errorf("Error obteniendo categoría por slug: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

// GetCategoryWithPosts obtiene una categoría con sus posts
func (h *CategoryHandler) GetCategoryWithPosts(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de categoría inválido",
		})
		return
	}

	category, err := h.categoryService.GetCategoryWithPosts(categoryID)
	if err != nil {
		if err.Error() == "categoría no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Categoría no encontrada",
			})
			return
		}
		h.logger.Errorf("Error obteniendo categoría con posts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
	})
}

// CreateCategory crea una nueva categoría
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		if err.Error() == "el slug de categoría ya existe" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Errorf("Error creando categoría: %v", err)
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
		"category_created",
		"category",
		&category.ID,
		map[string]interface{}{
			"name": category.Name,
			"slug": category.Slug,
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusCreated, gin.H{
		"category": category,
		"message":  "Categoría creada exitosamente",
	})
}

// UpdateCategory actualiza una categoría existente
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de categoría inválido",
		})
		return
	}

	var req models.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	category, err := h.categoryService.UpdateCategory(categoryID, req)
	if err != nil {
		if err.Error() == "categoría no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Categoría no encontrada",
			})
			return
		}
		h.logger.Errorf("Error actualizando categoría: %v", err)
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
		"category_updated",
		"category",
		&categoryID,
		map[string]interface{}{
			"name": category.Name,
			"slug": category.Slug,
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"message":  "Categoría actualizada exitosamente",
	})
}

// DeleteCategory elimina una categoría
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de categoría inválido",
		})
		return
	}

	err = h.categoryService.DeleteCategory(categoryID)
	if err != nil {
		if err.Error() == "categoría no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Categoría no encontrada",
			})
			return
		}
		if err.Error() == "no se puede eliminar la categoría porque tiene posts asociados" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Errorf("Error eliminando categoría: %v", err)
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
		"category_deleted",
		"category",
		&categoryID,
		map[string]interface{}{
			"category_id": categoryID.String(),
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Categoría eliminada exitosamente",
	})
}
