// brain-fitness/musicalnotes/musicalnotes.go
package musicalnotes

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
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

// TestUser will be replaced by API logic in main.go
func (n *MusicalNote) TestUser() bool {
	// This function will largely be obsolete, replaced by multiple API calls and frontend logic
	return n.CheckNext() && n.CheckPrevious() && n.CheckPosition() && n.CheckSound()
}

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

// RunNote will NOT be called from the Go backend.
// The frontend (Vue.js) will handle audio playback using HTML5 Audio API.
// This function is kept here for reference or if you later decide to stream audio differently.
func RunNote(audioPath string) error {
	fmt.Printf("Attempting to play sound (backend only) %s\n", audioPath)
	f, err := os.Open(audioPath)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	// Initialize oto context for playing sound
	// This part is for *backend* sound playing, which we are moving away from.
	// We're keeping it for compilation, but it won't be called for the game logic.
	context, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()

	if _, err := io.Copy(player, d); err != nil {
		return err
	}
	return nil
}

// DisplayStaff is now purely a frontend concern.
// This function will be removed from actual API usage.
func DisplayStaff() {
	fmt.Println("\n(DisplayStaff is now handled by frontend. This is a placeholder output.)\n")
	fmt.Println("C6                          ---                     29")
	fmt.Println("                                                    28")
	fmt.Println("                            ---                     27")
	fmt.Println("                                                    26")
	fmt.Println("            ----------------------------        25")
	fmt.Println("                                                    24")
	fmt.Println("            ----------------------------        23")
	fmt.Println("C5                                                22")
	fmt.Println("            ----------------------------        21")
	fmt.Println("                                                    20")
	fmt.Println("G4          ----------------------------        19")
	fmt.Println("                                                    18")
	fmt.Println("            ----------------------------        17")
	fmt.Println("                                                    16")
	fmt.Println("C4                          ---                     15")
	fmt.Println("                                                    14")
	fmt.Println("            ----------------------------        13")
	fmt.Println("                                                    12")
	fmt.Println("F3          ----------------------------        11")
	fmt.Println("                                                    10")
	fmt.Println("            ----------------------------        9")
	fmt.Println("C3                                                8")
	fmt.Println("            ----------------------------        7")
	fmt.Println("                                                    6")
	fmt.Println("            ----------------------------        5")
	fmt.Println("                                                    4")
	fmt.Println("                            ---                     3")
	fmt.Println("                                                    2")
	fmt.Println("C2                          ---                     1")
	fmt.Println("")
}
