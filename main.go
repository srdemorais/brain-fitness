// $ curl -X POST http://localhost:8080/api/note/new
// {"idx":3,"code":"Eb2","audioPath":"/audio/Eb2.mp3","position":3}

package main

import (
	"fmt"
	"log"
	"net/http" // Required for HTTP status codes

	"github.com/gin-gonic/gin"                         // Gin framework
	"github.com/srdemorais/brain-fitness/musicalnotes" // Our local musicalnotes package
)

// --- Main Server Setup ---

func main() {
	router := gin.Default()

	// --- CORS (Cross-Origin Resource Sharing) Setup ---
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Adjust in production
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// --- End CORS Setup ---

	// Serve static MP3 files from the /mp3 directory under the /audio path
	router.Static("/audio", "./musicalnotes/mp3")

	// --- API Endpoints ---

	// POST /api/note/new: Get a new random note for a round
	router.POST("/api/note/new", func(c *gin.Context) {
		note, err := musicalnotes.Init() // Generates a random note
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, note) // Returns the MusicalNote struct as JSON
	})

	// Start the server
	port := ":8080"
	fmt.Printf("Go API server listening on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
