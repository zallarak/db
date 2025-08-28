package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Test database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/dbxyz?sslmode=disable"
	}

	log.Printf("Connecting to database: %s", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✅ Database connection successful")

	// Setup simple HTTP server
	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		log.Println("Health check requested")
		c.JSON(http.StatusOK, gin.H{"status": "ok", "database": "connected"})
	})

	r.POST("/test", func(c *gin.Context) {
		log.Println("Test endpoint requested")
		c.JSON(http.StatusOK, gin.H{"message": "test successful"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Starting server on port %s", port)
	log.Fatal(r.Run(":" + port))
}