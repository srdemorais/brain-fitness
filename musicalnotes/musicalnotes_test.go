// Unit tests for the musicalnotes package
// The following command tells Go to run all tests in the musicalnotes directory and its subdirectories. You should see output indicating all tests passed.
// $ go test ./musicalnotes/...
package musicalnotes

import (
	"strings"
	"testing"
)

func TestGetNoteNameByIdx(t *testing.T) {
	tests := []struct {
		name     string
		idx      int
		expected string
	}{
		{"Valid Index 0", 0, "C2"},
		{"Valid Index Middle", 24, "C4"},
		{"Valid Index Last", len(NotesArray) - 1, "C6"},
		{"Invalid Index Negative", -1, ""}, // Should return empty string or error
		{"Invalid Index Out of Bounds", len(NotesArray), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetNoteNameByIdx(tt.idx)
			if actual != tt.expected {
				t.Errorf("GetNoteNameByIdx(%d): expected %s, got %s", tt.idx, tt.expected, actual)
			}
		})
	}
}

func TestGetNotePositionByIdx(t *testing.T) {
	tests := []struct {
		name     string
		idx      int
		expected int
	}{
		{"Valid Index 0", 0, 1},
		{"Valid Index Middle", 24, 15},
		{"Valid Index Last", len(NotesPosArray) - 1, 29},
		{"Invalid Index Negative", -1, 0}, // Should return 0 or handle error
		{"Invalid Index Out of Bounds", len(NotesPosArray), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetNotePositionByIdx(tt.idx)
			if actual != tt.expected {
				t.Errorf("GetNotePositionByIdx(%d): expected %d, got %d", tt.idx, tt.expected, actual)
			}
		})
	}
}

func TestCheckPosition(t *testing.T) {
	note := MusicalNote{Note: "C4", Idx: 24} // Position is 15
	tests := []struct {
		name    string
		input   int
		correct bool
	}{
		{"Correct Guess", 15, true},
		{"Incorrect Guess", 14, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := note.CheckPosition(tt.input)
			if actual != tt.correct {
				t.Errorf("CheckPosition with input %d: expected %v, got %v", tt.input, tt.correct, actual)
			}
		})
	}
}

func TestInit(t *testing.T) {
	note := Init()
	if note.Idx < 0 || note.Idx >= len(NotesArray) {
		t.Errorf("Init: returned invalid index %d", note.Idx)
	}
	if note.Note == "" {
		t.Errorf("Init: returned empty note string")
	}
	if !strings.HasPrefix(note.AudioPath, "/audio/") || !strings.HasSuffix(note.AudioPath, ".mp3") {
		t.Errorf("Init: AudioPath '%s' has incorrect format", note.AudioPath)
	}
	if note.Position == 0 {
		t.Errorf("Init: Position not correctly populated: %+v", note)
	}
}

func TestGetGuessNotes(t *testing.T) {
	testNote := MusicalNote{Note: "C4", Idx: 24} // The note we want to ensure is in the guess list
	guessNotes, correctPos := testNote.GetGuessNotes()

	if len(guessNotes) != 6 {
		t.Errorf("GetGuessNotes: expected 6 notes, got %d", len(guessNotes))
	}

	found := false
	for i, n := range guessNotes {
		if n.Note == testNote.Note {
			if i != correctPos {
				t.Errorf("GetGuessNotes: correct note '%s' found at index %d, but correctPos indicates %d", testNote.Note, i, correctPos)
			}
			found = true
		}
		// Basic validation of each guess note
		if n.Note == "" || n.AudioPath == "" || n.Position == 0 {
			t.Errorf("GetGuessNotes: found malformed guess note at index %d: %+v", i, n)
		}
	}

	if !found {
		t.Errorf("GetGuessNotes: original note '%s' not found in guessNotes", testNote.Note)
	}
}
