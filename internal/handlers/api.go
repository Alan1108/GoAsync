package handlers

import (
	"database/sql"

	"github.com/alan.bermudez/goasync/internal/services"
	"github.com/alan.bermudez/goasync/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(r *gin.Engine, db *sql.DB, logger *logrus.Logger) {
	// Crear servicios
	userService := services.NewUserService(db, logger)
	postService := services.NewPostService(db, logger)
	categoryService := services.NewCategoryService(db, logger)
	tagService := services.NewTagService(db, logger)
	commentService := services.NewCommentService(db, logger)
	statsService := services.NewStatsService(db, logger)

	// Crear handlers
	userHandler := NewUserHandler(userService, statsService, logger)
	postHandler := NewPostHandler(postService, statsService, logger)
	categoryHandler := NewCategoryHandler(categoryService, statsService, logger)
	tagHandler := NewTagHandler(tagService, statsService, logger)
	commentHandler := NewCommentHandler(commentService, statsService, logger)
	statsHandler := NewStatsHandler(statsService, logger)
	healthHandler := NewHealthHandler(db, logger)

	// Middleware global
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(logger))

	// Grupo de rutas de la API
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthHandler.HealthCheck)

		// Rutas de usuarios
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/stats", userHandler.GetAllUserStats)
			users.GET("/:id", userHandler.GetUser)
			users.GET("/:id/profile", userHandler.GetUserWithProfile)
			users.GET("/:id/stats", userHandler.GetUserStats)
			users.GET("/:id/activity", userHandler.GetUserActivity)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Rutas de posts
		posts := api.Group("/posts")
		{
			posts.GET("", postHandler.GetPosts)
			posts.GET("/published", postHandler.GetPublishedPosts)
			posts.GET("/:id", postHandler.GetPost)
			posts.GET("/slug/:slug", postHandler.GetPostBySlug)
			posts.GET("/:id/with-tags", postHandler.GetPostWithTags)
			posts.POST("", postHandler.CreatePost)
			posts.PUT("/:id", postHandler.UpdatePost)
			posts.DELETE("/:id", postHandler.DeletePost)
			posts.GET("/:id/comments", commentHandler.GetComments)
		}

		// Rutas de categorías
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
			categories.GET("/:id/with-posts", categoryHandler.GetCategoryWithPosts)
			categories.POST("", categoryHandler.CreateCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// Rutas de tags
		tags := api.Group("/tags")
		{
			tags.GET("", tagHandler.GetTags)
			tags.GET("/popular", tagHandler.GetPopularTags)
			tags.GET("/:id", tagHandler.GetTag)
			tags.GET("/slug/:slug", tagHandler.GetTagBySlug)
			tags.GET("/:id/with-posts", tagHandler.GetTagWithPosts)
			tags.POST("", tagHandler.CreateTag)
			tags.PUT("/:id", tagHandler.UpdateTag)
			tags.DELETE("/:id", tagHandler.DeleteTag)
		}

		// Rutas de comentarios
		comments := api.Group("/comments")
		{
			comments.GET("", commentHandler.GetAllComments)
			comments.GET("/:id", commentHandler.GetComment)
			comments.POST("", commentHandler.CreateComment)
			comments.PUT("/:id", commentHandler.UpdateComment)
			comments.DELETE("/:id", commentHandler.DeleteComment)
			comments.PATCH("/:id/approve", commentHandler.ApproveComment)
		}

		// Rutas de estadísticas
		stats := api.Group("/stats")
		{
			stats.GET("/database", statsHandler.GetDatabaseStats)
			stats.GET("/activity", statsHandler.GetActivityLogs)
			stats.GET("/activity/recent", statsHandler.GetRecentActivity)
			stats.GET("/activity/user/:user_id", statsHandler.GetUserActivity)
			stats.GET("/posts", statsHandler.GetPostStats)
			stats.GET("/daily", statsHandler.GetDailyStats)
		}
	}

	// Ruta raíz
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Bienvenido a la API de GoAsync",
			"version": "1.0.0",
			"docs":    "/api/v1/health",
		})
	})
}
