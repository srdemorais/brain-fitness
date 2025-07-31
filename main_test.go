// Integration Tests for main.go API Endpoints
// The following command tells Go to run all tests in the current directory (main_test.go)
// $ go test .
// You can also run all tests (unit and integration) with the following command:
// $ go test ./...

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/srdemorais/brain-fitness/musicalnotes" // Import our musicalnotes package
	"github.com/stretchr/testify/assert"               // We'll use testify for cleaner assertions
)

// Helper function to setup the Gin router for testing
func setupRouter() *gin.Engine {
	router := gin.Default()
	// Re-add CORS for testing if needed, or disable it
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Register your API endpoints as they are in main.go
	router.Static("/audio", "./mp3") // Ensure static files are also served for tests if needed

	router.POST("/api/note/new", func(c *gin.Context) {
		note, _ := musicalnotes.Init()
		c.JSON(http.StatusOK, note)
	})

	return router
}

// Install testify: go get github.com/stretchr/testify
func TestNewNoteEndpoint(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/note/new", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var note musicalnotes.MusicalNote
	err := json.Unmarshal(w.Body.Bytes(), &note)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, note.Idx, 0)
	assert.Less(t, note.Idx, len(musicalnotes.NotesCodeArray)) // Assuming you expose notes array for test or mock it
	assert.NotEmpty(t, note.Code)
	assert.Contains(t, note.AudioPath, "/audio/")
	assert.Greater(t, note.Position, 0)
}
