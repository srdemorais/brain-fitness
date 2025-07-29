// musical-ear-trainer/main.go
package main

import (
	"fmt"
	"log"
	"net/http" // Required for HTTP status codes

	"github.com/gin-gonic/gin" // Gin framework
	// "github.com/srdemorais/brain-fitness/musicalnotes"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// --- CORS (Cross-Origin Resource Sharing) Setup ---
	// Since the frontend (Vue.js) will be running on a different port (e.g., 8081)
	// than the backend (8080), we need to enable CORS for the browser to allow
	// requests from the frontend to the backend.
	router.Use(func(c *gin.Context) {
		// Allow requests from all origins during development.
		// In production, we would replace "*" with our frontend's exact domain.
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	// --- End CORS Setup ---

	// Define a simple API endpoint for testing purposes
	router.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Go API! This is your new project.",
		})
	})

	// Add a handler for serving MP3 files
	// This will serve files from the 'mp3' directory
	router.Static("/audio", "./musicalnotes/mp3")

	// Start the server
	port := ":8080" // Backend will run on port 8080
	fmt.Printf("Go API server listening on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
