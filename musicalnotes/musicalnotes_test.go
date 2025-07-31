// Unit tests for the musicalnotes package
// The following command tells Go to run all tests in the musicalnotes directory and its subdirectories. You should see output indicating all tests passed.
// $ go test ./musicalnotes/...
package musicalnotes

import (
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	note, _ := Init()
	if note.Idx < 0 || note.Idx >= len(NotesCodeArray) {
		t.Errorf("Init: returned invalid index %d", note.Idx)
	}
	if note.Code == "" {
		t.Errorf("Init: returned empty note string")
	}
	if !strings.HasPrefix(note.AudioPath, "/audio/") || !strings.HasSuffix(note.AudioPath, ".mp3") {
		t.Errorf("Init: AudioPath '%s' has incorrect format", note.AudioPath)
	}
	if note.Position == 0 {
		t.Errorf("Init: Position not correctly populated: %+v", note)
	}
}
