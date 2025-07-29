// brain-fitness/musicalnotes/musicalnotes.go
package musicalnotes

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var notes = [...]string{"C2", "Db2", "D2", "Eb2", "E2", "F2", "Gb2", "G2", "Ab2", "A2", "Bb2", "B2", "C3", "Db3", "D3", "Eb3", "E3", "F3", "Gb3", "G3", "Ab3", "A3", "Bb3", "B3", "C4", "Db4", "D4", "Eb4", "E4", "F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4", "C5", "Db5", "D5", "Eb5", "E5", "F5", "Gb5", "G5", "Ab5", "A5", "Bb5", "B5", "C6"}
var notesPos = [...]int{1, 2, 2, 3, 3, 4, 5, 5, 6, 6, 7, 7, 8, 9, 9, 10, 10, 11, 12, 12, 13, 13, 14, 14, 15, 16, 16, 17, 17, 18, 19, 19, 20, 20, 21, 21, 22, 23, 23, 24, 24, 25, 26, 26, 27, 27, 28, 28, 29}

type MusicalNote struct {
	Idx          int    `json:"idx"` // Add JSON tags for easy serialization
	Note         string `json:"note"`
	AudioPath    string `json:"audioPath"`
	NextNote     string `json:"nextNote"`     // New field to pre-calculate for API
	PreviousNote string `json:"previousNote"` // New field to pre-calculate for API
	Position     int    `json:"position"`     // New field to pre-calculate for API
}

func Init() MusicalNote {
	rand.Seed(time.Now().UnixNano()) // Use UnixNano for better randomness in quick succession
	idx := rand.Intn(len(notes))

	var oMusicalNote MusicalNote
	oMusicalNote.Idx = idx
	oMusicalNote.Note = notes[idx]
	oMusicalNote.AudioPath = "/audio/" + notes[idx] + ".mp3" // Change path to match Gin static server

	// Pre-calculate next, previous, and position for the API response
	oMusicalNote.NextNote = oMusicalNote.GetNext()
	oMusicalNote.PreviousNote = oMusicalNote.GetPrevious()
	oMusicalNote.Position = notesPos[idx]

	return oMusicalNote
}

// NOTE: TestUser, CheckNext, CheckPrevious, CheckPosition, CheckSound will be heavily refactored
// to work with API requests instead of console input.
// For now, let's keep the originals for reference but they won't be used directly by main.go anymore.

// CheckNext refactored to take user input and return result
func (n *MusicalNote) CheckNext(userInput string) bool {
	next := n.GetNext()
	return strings.ToUpper(next) == strings.ToUpper(userInput)
}

func (n *MusicalNote) GetNext() string {
	var next string
	if n.Note == "C6" {
		next = "Db6" // Or potentially a custom "end of range" note
	} else {
		next = notes[n.Idx+1]
	}
	return next
}

// CheckPrevious refactored to take user input and return result
func (n *MusicalNote) CheckPrevious(userInput string) bool {
	previous := n.GetPrevious()
	return strings.ToUpper(previous) == strings.ToUpper(userInput)
}

func (n *MusicalNote) GetPrevious() string {
	var previous string
	if n.Note == "C2" {
		previous = "B1" // Or potentially a custom "start of range" note
	} else {
		previous = notes[n.Idx-1]
	}
	return previous
}

// CheckPosition refactored to take user input and return result
func (n *MusicalNote) CheckPosition(userInput int) bool {
	return notesPos[n.Idx] == userInput
}

// CheckSound will be heavily refactored. For now, this is a placeholder.
// It will need to return the list of guess notes and the correct position,
// and the frontend will handle playing and user input.
func (n *MusicalNote) CheckSound() bool {
	// This function's logic will be split between API and frontend
	fmt.Println("CheckSound is placeholder for API refactoring.")
	return false // Placeholder
}

// getGuessNotes can remain as a helper for generating sound options
// but RunNote will not be called from the backend in the final API version.
func (n *MusicalNote) GetGuessNotes() ([6]MusicalNote, int) { // Exported for API use
	var guessNotes [6]MusicalNote
	var pos int

	var randPos [6]int
	randPos[0] = n.Idx

	rand.Seed(time.Now().UnixNano())

	p := 1
	for {
		tmpIdx := rand.Intn(len(notes))
		// Ensure unique notes and don't pick the same note as the current one too many times
		isUnique := true
		for i := 0; i < p; i++ {
			if randPos[i] == tmpIdx {
				isUnique = false
				break
			}
		}
		if isUnique {
			randPos[p] = tmpIdx
			p++
		}
		if p > len(randPos)-1 {
			break
		}
	}

	// shuffle randpos
	rand.Shuffle(len(randPos), func(i, j int) { randPos[i], randPos[j] = randPos[j], randPos[i] })

	// set guessNotes
	for i := 0; i < 6; i++ {
		guessNotes[i].Idx = randPos[i]
		guessNotes[i].Note = notes[randPos[i]]
		guessNotes[i].AudioPath = "/audio/" + notes[randPos[i]] + ".mp3" // Change path for frontend
		guessNotes[i].Position = notesPos[randPos[i]]                    // Add position for frontend
		if n.Note == guessNotes[i].Note {
			pos = i
		}
	}

	return guessNotes, pos
}

// GetNoteNameByIdx returns the note string for a given index
func GetNoteNameByIdx(idx int) string {
	if idx >= 0 && idx < len(notes) {
		return notes[idx]
	}
	return "" // Or handle error appropriately
}

// GetNotePositionByIdx returns the note position for a given index
func GetNotePositionByIdx(idx int) int {
	if idx >= 0 && idx < len(notesPos) {
		return notesPos[idx]
	}
	return 0 // Or handle error appropriately
}
