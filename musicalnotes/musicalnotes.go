package musicalnotes

import (
	"errors"
	"math/rand"
)

// NotesArray and NotesPosArray are the musical notes and their positions in the staff
var NotesCodeArray = [...]string{"C2", "Db2", "D2", "Eb2", "E2", "F2", "Gb2", "G2", "Ab2", "A2", "Bb2", "B2", "C3", "Db3", "D3", "Eb3", "E3", "F3", "Gb3", "G3", "Ab3", "A3", "Bb3", "B3", "C4", "Db4", "D4", "Eb4", "E4", "F4", "Gb4", "G4", "Ab4", "A4", "Bb4", "B4", "C5", "Db5", "D5", "Eb5", "E5", "F5", "Gb5", "G5", "Ab5", "A5", "Bb5", "B5", "C6"}
var NotesPosArray = [...]int{1, 2, 2, 3, 3, 4, 5, 5, 6, 6, 7, 7, 8, 9, 9, 10, 10, 11, 12, 12, 13, 13, 14, 14, 15, 16, 16, 17, 17, 18, 19, 19, 20, 20, 21, 21, 22, 23, 23, 24, 24, 25, 26, 26, 27, 27, 28, 28, 29}

type MusicalNote struct {
	Idx       int    `json:"idx"`       // Index of the note in the array
	Code      string `json:"code"`      // Note code (e.g., "C4", "Db4")
	AudioPath string `json:"audioPath"` // Path to the audio file (e.g., "/audio/C4.mp3")
	Position  int    `json:"position"`  // Position in the staff (1-29)
}

// Init() initializes a new MusicalNote with a random index and populates its fields
func Init() (MusicalNote, error) {

	// Note: since Go 1.20 the global math/rand source is now automatically seeded at program startup
	// with a non-deterministic value. So we can use rand.Intn directly without seeding.
	idx := rand.Intn(len(NotesCodeArray)) // Randomly select an index from the NotesCodeArray

	var oMusicalNote MusicalNote
	if idx < 0 || idx >= len(NotesCodeArray) {
		return oMusicalNote, errors.New("index out of bounds for notes array (0-28)")
	}

	// Populate the MusicalNote fields
	oMusicalNote.Idx = idx
	oMusicalNote.Code = NotesCodeArray[idx]
	oMusicalNote.AudioPath = "/audio/" + NotesCodeArray[idx] + ".mp3" // Path matches Gin static server
	oMusicalNote.Position = NotesPosArray[idx]

	return oMusicalNote, nil
}
