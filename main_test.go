// Integration Tests for main.go API Endpoints
// The following command tells Go to run all tests in the current directory (main_test.go)
// $ go test .
// You can also run all tests (unit and integration) with the following command:
// $ go test ./...

package main

import (
	"bytes"
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
		note := musicalnotes.Init()
		c.JSON(http.StatusOK, note)
	})

	router.POST("/api/note/check_text_position", func(c *gin.Context) {
		var req CheckTextPositionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		currentNote := musicalnotes.MusicalNote{
			Idx:  req.NoteIdx,
			Note: musicalnotes.GetNoteNameByIdx(req.NoteIdx),
		}

		positionCorrect := currentNote.CheckPosition(req.UserPosition)

		resp := CheckTextPositionResponse{
			PositionCorrect: positionCorrect,
		}
		if !positionCorrect {
			resp.CorrectPosition = musicalnotes.GetNotePositionByIdx(req.NoteIdx)
		}
		c.JSON(http.StatusOK, resp)
	})

	router.POST("/api/note/prepare_sound_test", func(c *gin.Context) {
		var req PrepareSoundTestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		dummyNote := musicalnotes.MusicalNote{Idx: req.NoteIdx}
		guessNotes, correctPos := dummyNote.GetGuessNotes()
		resp := PrepareSoundTestResponse{
			GuessNotes:      guessNotes,
			CorrectGuessPos: correctPos,
		}
		c.JSON(http.StatusOK, resp)
	})

	router.POST("/api/note/check_sound_guess", func(c *gin.Context) {
		var req CheckSoundGuessRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		soundCorrect := (req.CorrectGuessPos == req.UserSoundGuessPos)
		resp := CheckSoundGuessResponse{
			SoundCorrect: soundCorrect,
		}
		c.JSON(http.StatusOK, resp)
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
	assert.Less(t, note.Idx, len(musicalnotes.NotesArray)) // Assuming you expose notes array for test or mock it
	assert.NotEmpty(t, note.Note)
	assert.Contains(t, note.AudioPath, "/audio/")
	assert.Greater(t, note.Position, 0)
}

func TestCheckTextPositionEndpoint(t *testing.T) {
	router := setupRouter()

	// --- Test Case 1: All correct answers ---
	correctNoteIdx := 24 // C4 (Next: Db4, Previous: B3, Pos: 15)
	correctReqBody, _ := json.Marshal(CheckTextPositionRequest{
		NoteIdx:      correctNoteIdx,
		UserPosition: 15,
	})
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/note/check_text_position", bytes.NewBuffer(correctReqBody))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	var resp1 CheckTextPositionResponse
	json.Unmarshal(w1.Body.Bytes(), &resp1)
	assert.True(t, resp1.PositionCorrect)
	assert.Equal(t, 0, resp1.CorrectPosition) // For int, omitempty usually means 0

	// --- Test Case 2: One incorrect answer (UserNext is wrong) ---
	incorrectNextReqBody, _ := json.Marshal(CheckTextPositionRequest{
		NoteIdx:      correctNoteIdx,
		UserPosition: 15,
	})
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/note/check_text_position", bytes.NewBuffer(incorrectNextReqBody))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp2 CheckTextPositionResponse
	json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.True(t, resp2.PositionCorrect)
	assert.Equal(t, 0, resp2.CorrectPosition) // For int, omitempty usually means 0
}

func TestCheckSoundGuessEndpoint(t *testing.T) {
	router := setupRouter()

	// --- Test Case 1: Correct guess ---
	correctGuessReqBody, _ := json.Marshal(CheckSoundGuessRequest{
		CorrectGuessPos:   3,
		UserSoundGuessPos: 3,
	})
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/note/check_sound_guess", bytes.NewBuffer(correctGuessReqBody))
	req1.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	var resp1 CheckSoundGuessResponse
	json.Unmarshal(w1.Body.Bytes(), &resp1)
	assert.True(t, resp1.SoundCorrect)

	// --- Test Case 2: Incorrect guess ---
	incorrectGuessReqBody, _ := json.Marshal(CheckSoundGuessRequest{
		CorrectGuessPos:   3,
		UserSoundGuessPos: 2, // Incorrect
	})
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/note/check_sound_guess", bytes.NewBuffer(incorrectGuessReqBody))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var resp2 CheckSoundGuessResponse
	json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.False(t, resp2.SoundCorrect)
}
