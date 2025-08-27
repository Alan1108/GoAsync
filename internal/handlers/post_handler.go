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

// PostHandler maneja las peticiones HTTP relacionadas con posts
type PostHandler struct {
	postService  *services.PostService
	statsService *services.StatsService
	logger       *logrus.Logger
}

// NewPostHandler crea una nueva instancia del handler de posts
func NewPostHandler(postService *services.PostService, statsService *services.StatsService, logger *logrus.Logger) *PostHandler {
	return &PostHandler{
		postService:  postService,
		statsService: statsService,
		logger:       logger,
	}
}

// GetPosts obtiene todos los posts con filtros y paginación
func (h *PostHandler) GetPosts(c *gin.Context) {
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

	// Construir filtros
	filter := models.PostFilter{
		Page:    page,
		PerPage: perPage,
	}

	// Aplicar filtros opcionales
	if status := c.Query("status"); status != "" {
		filter.Status = status
	}

	if search := c.Query("search"); search != "" {
		filter.Search = search
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			filter.CategoryID = categoryID
		}
	}

	if authorIDStr := c.Query("author_id"); authorIDStr != "" {
		if authorID, err := uuid.Parse(authorIDStr); err == nil {
			filter.AuthorID = authorID
		}
	}

	if tagIDStr := c.Query("tag_id"); tagIDStr != "" {
		if tagID, err := uuid.Parse(tagIDStr); err == nil {
			filter.TagID = tagID
		}
	}

	response, err := h.postService.GetAllPosts(filter)
	if err != nil {
		h.logger.Errorf("Error obteniendo posts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": response.Posts,
		"pagination": gin.H{
			"page":        response.Page,
			"per_page":    response.PerPage,
			"total":       response.Total,
			"total_pages": response.TotalPages,
		},
	})
}

// GetPublishedPosts obtiene solo posts publicados
func (h *PostHandler) GetPublishedPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 10000 {
		perPage = 10
	}

	response, err := h.postService.GetPublishedPosts(page, perPage)
	if err != nil {
		h.logger.Errorf("Error obteniendo posts publicados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": response.Posts,
		"pagination": gin.H{
			"page":        response.Page,
			"per_page":    response.PerPage,
			"total":       response.Total,
			"total_pages": response.TotalPages,
		},
	})
}

// GetPost obtiene un post por su ID
func (h *PostHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de post inválido",
		})
		return
	}

	post, err := h.postService.GetPostByID(postID)
	if err != nil {
		if err.Error() == "post no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// GetPostBySlug obtiene un post por su slug
func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Slug de post requerido",
		})
		return
	}

	post, err := h.postService.GetPostBySlug(slug)
	if err != nil {
		if err.Error() == "post no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo post por slug: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// GetPostWithTags obtiene un post con sus tags
func (h *PostHandler) GetPostWithTags(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de post inválido",
		})
		return
	}

	post, err := h.postService.GetPostWithTags(postID)
	if err != nil {
		if err.Error() == "post no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post no encontrado",
			})
			return
		}
		h.logger.Errorf("Error obteniendo post con tags: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// CreatePost crea un nuevo post
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req models.PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	// TODO: Obtener el ID del usuario autenticado
	// Por ahora, usamos un ID falso para pruebas
	authorID := uuid.New() // En producción, esto vendría del contexto de autenticación

	post, err := h.postService.CreatePost(req, authorID)
	if err != nil {
		h.logger.Errorf("Error creando post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&authorID,
		"post_created",
		"post",
		&post.ID,
		map[string]interface{}{
			"title":  post.Title,
			"slug":   post.Slug,
			"status": post.Status,
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusCreated, gin.H{
		"post":    post,
		"message": "Post creado exitosamente",
	})
}

// UpdatePost actualiza un post existente
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de post inválido",
		})
		return
	}

	var req models.PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	post, err := h.postService.UpdatePost(postID, req)
	if err != nil {
		if err.Error() == "post no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post no encontrado",
			})
			return
		}
		h.logger.Errorf("Error actualizando post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// TODO: Obtener el ID del usuario autenticado
	authorID := uuid.New() // En producción, esto vendría del contexto de autenticación

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&authorID,
		"post_updated",
		"post",
		&postID,
		map[string]interface{}{
			"title":  post.Title,
			"slug":   post.Slug,
			"status": post.Status,
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"post":    post,
		"message": "Post actualizado exitosamente",
	})
}

// DeletePost elimina un post
func (h *PostHandler) DeletePost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de post inválido",
		})
		return
	}

	err = h.postService.DeletePost(postID)
	if err != nil {
		if err.Error() == "post no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post no encontrado",
			})
			return
		}
		h.logger.Errorf("Error eliminando post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error interno del servidor",
		})
		return
	}

	// TODO: Obtener el ID del usuario autenticado
	authorID := uuid.New() // En producción, esto vendría del contexto de autenticación

	// Crear log de actividad
	h.statsService.CreateActivityLog(
		&authorID,
		"post_deleted",
		"post",
		&postID,
		map[string]interface{}{
			"post_id": postID.String(),
		},
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Post eliminado exitosamente",
	})
}
