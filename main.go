// musical-ear-trainer/main.go
package main

import (
	"fmt"
	"log"
	"net/http" // Required for HTTP status codes

	"github.com/gin-gonic/gin"                         // Gin framework
	"github.com/srdemorais/brain-fitness/musicalnotes" // Our local musicalnotes package
)

// --- Request and Response Structs for API Endpoints ---

// CheckTextPositionRequest represents the JSON payload for checking text/position answers
type CheckTextPositionRequest struct {
	NoteIdx      int    `json:"noteIdx"` // Index of the original note
	UserNext     string `json:"userNext"`
	UserPrevious string `json:"userPrevious"`
	UserPosition int    `json:"userPosition"`
}

// CheckTextPositionResponse represents the JSON response for text/position checks
type CheckTextPositionResponse struct {
	NextCorrect     bool   `json:"nextCorrect"`
	CorrectNext     string `json:"correctNext,omitempty"` // omitempty means it won't be included if empty (i.e., if correct)
	PreviousCorrect bool   `json:"previousCorrect"`
	CorrectPrevious string `json:"correctPrevious,omitempty"`
	PositionCorrect bool   `json:"positionCorrect"`
	CorrectPosition int    `json:"correctPosition,omitempty"`
}

// PrepareSoundTestRequest represents the JSON payload for preparing the sound test
type PrepareSoundTestRequest struct {
	NoteIdx int `json:"noteIdx"` // Index of the original note
}

// PrepareSoundTestResponse represents the JSON response for preparing the sound test
type PrepareSoundTestResponse struct {
	GuessNotes      [6]musicalnotes.MusicalNote `json:"guessNotes"`      // The 6 notes to play
	CorrectGuessPos int                         `json:"correctGuessPos"` // 0-indexed position of the correct note within GuessNotes
}

// CheckSoundGuessRequest represents the JSON payload for checking the sound guess
type CheckSoundGuessRequest struct {
	CorrectGuessPos   int `json:"correctGuessPos"`   // The original correct 0-indexed position from prepare_sound_test
	UserSoundGuessPos int `json:"userSoundGuessPos"` // The user's 0-indexed guess
}

// CheckSoundGuessResponse represents the JSON response for the sound guess
type CheckSoundGuessResponse struct {
	SoundCorrect bool `json:"soundCorrect"`
}

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
	router.Static("/audio", "./mp3")

	// --- API Endpoints ---

	// POST /api/note/new: Get a new random note for a round
	router.POST("/api/note/new", func(c *gin.Context) {
		note := musicalnotes.Init() // Generates a random note
		c.JSON(http.StatusOK, note) // Returns the MusicalNote struct as JSON
	})

	// POST /api/note/check_text_position: Check user's text and position answers
	router.POST("/api/note/check_text_position", func(c *gin.Context) {
		var req CheckTextPositionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Re-initialize a MusicalNote based on the received index for checking
		// This ensures we're checking against the correct, untampered backend data
		currentNote := musicalnotes.MusicalNote{
			Idx: req.NoteIdx,
			// Other fields are not strictly needed for Check methods but good to have a complete object
			Note: musicalnotes.GetNoteNameByIdx(req.NoteIdx), // Helper needed: see note below
		}

		// Perform checks
		nextCorrect := currentNote.CheckNext(req.UserNext)
		previousCorrect := currentNote.CheckPrevious(req.UserPrevious)
		positionCorrect := currentNote.CheckPosition(req.UserPosition)

		resp := CheckTextPositionResponse{
			NextCorrect:     nextCorrect,
			PreviousCorrect: previousCorrect,
			PositionCorrect: positionCorrect,
		}

		// If incorrect, provide the correct answer
		if !nextCorrect {
			resp.CorrectNext = currentNote.GetNext()
		}
		if !previousCorrect {
			resp.CorrectPrevious = currentNote.GetPrevious()
		}
		if !positionCorrect {
			resp.CorrectPosition = musicalnotes.GetNotePositionByIdx(req.NoteIdx) // Helper needed: see note below
		}

		c.JSON(http.StatusOK, resp)
	})

	// POST /api/note/prepare_sound_test: Get notes for the sound test
	router.POST("/api/note/prepare_sound_test", func(c *gin.Context) {
		var req PrepareSoundTestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create a dummy note just to use its GetGuessNotes method
		dummyNote := musicalnotes.MusicalNote{Idx: req.NoteIdx}
		guessNotes, correctPos := dummyNote.GetGuessNotes()

		resp := PrepareSoundTestResponse{
			GuessNotes:      guessNotes,
			CorrectGuessPos: correctPos,
		}

		c.JSON(http.StatusOK, resp)
	})

	// POST /api/note/check_sound_guess: Check user's guess for the sound test
	router.POST("/api/note/check_sound_guess", func(c *gin.Context) {
		var req CheckSoundGuessRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// The CheckSoundGuess logic is essentially comparing two integers
		soundCorrect := (req.CorrectGuessPos == req.UserSoundGuessPos)

		resp := CheckSoundGuessResponse{
			SoundCorrect: soundCorrect,
		}

		c.JSON(http.StatusOK, resp)
	})

	// Start the server
	port := ":8080"
	fmt.Printf("Go API server listening on port %s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
