package main

import (
	"log"
	"os"

	"github.com/zallarak/db/api/internal/auth"
	"github.com/zallarak/db/api/internal/db"
	"github.com/zallarak/db/api/internal/handlers"
	"github.com/zallarak/db/api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	database, err := db.Init()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Create auth service
	authService := auth.NewService(database)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(database)
	orgHandler := handlers.NewOrgHandler(database)

	// Setup router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve OpenAPI documentation
	r.Static("/docs", "./openapi.yaml")
	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("./openapi.yaml")
	})

	// API v1 routes
	v1 := r.Group("/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired(authService))
		{
			// User routes
			protected.GET("/users/me", userHandler.GetCurrentUser)

			// Org routes
			orgs := protected.Group("/orgs")
			{
				orgs.GET("", orgHandler.ListOrgs)
				orgs.POST("", orgHandler.CreateOrg)
				orgs.GET("/:orgId", orgHandler.GetOrg)
				orgs.PATCH("/:orgId", orgHandler.UpdateOrg)
				orgs.DELETE("/:orgId", orgHandler.DeleteOrg)
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}