package cli_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreasphil/one/cli"
	"github.com/andreasphil/one/lib/note"
)

func writeNotes(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func newProvider(t *testing.T, content string) (*cli.FileNotesProvider, string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "one.md")
	writeNotes(t, path, content)

	provider, err := cli.NewFileNotesProvider(path, io.Discard)
	if err != nil {
		t.Fatalf("NewFileNotesProvider() error = %v", err)
	}

	return provider, path
}

func titles(notes []note.Note) []string {
	var got []string
	for _, n := range note.Flat(notes) {
		got = append(got, n.Title)
	}

	return got
}

func TestFileNotesProviderParsesOnStartup(t *testing.T) {
	provider, _ := newProvider(t, "# First\n\nhello\n")

	if got := titles(provider.Notes()); len(got) != 1 || got[0] != "First" {
		t.Errorf("Notes() titles = %v, want [First]", got)
	}
}

func TestFileNotesProviderReparsesChangedFile(t *testing.T) {
	provider, path := newProvider(t, "# First\n\nhello\n")

	writeNotes(t, path, "# First\n\nhello\n\n# Second\n\nworld, and then some\n")

	got := titles(provider.Notes())
	if len(got) != 2 || got[1] != "Second" {
		t.Errorf("Notes() titles = %v, want [First Second]", got)
	}
}

func TestFileNotesProviderKeepsNotesOnParseError(t *testing.T) {
	provider, path := newProvider(t, "# First\n\nhello\n")

	writeNotes(t, path, "no heading here, so parsing fails\n")

	if got := titles(provider.Notes()); len(got) != 1 || got[0] != "First" {
		t.Errorf("Notes() titles = %v, want the previously parsed [First]", got)
	}
}

func TestFileNotesProviderKeepsNotesOnMissingFile(t *testing.T) {
	provider, path := newProvider(t, "# First\n\nhello\n")

	if err := os.Remove(path); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	if got := titles(provider.Notes()); len(got) != 1 || got[0] != "First" {
		t.Errorf("Notes() titles = %v, want the previously parsed [First]", got)
	}
}

func TestFileNotesProviderRecoversAfterParseError(t *testing.T) {
	provider, path := newProvider(t, "# First\n\nhello\n")

	writeNotes(t, path, "no heading here, so parsing fails\n")
	provider.Notes()

	writeNotes(t, path, "# Fixed again\n\nhello\n")

	if got := titles(provider.Notes()); len(got) != 1 || got[0] != "Fixed again" {
		t.Errorf("Notes() titles = %v, want [Fixed again]", got)
	}
}
