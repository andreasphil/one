package cli

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

type fileNotesProvider struct {
	path string
	logw io.Writer

	mu      sync.Mutex
	notes   []note.Note
	modTime time.Time
	size    int64
}

func newFileNotesProvider(path string, logw io.Writer) (*fileNotesProvider, error) {
	notes, err := note.ParseFile(path)
	if err != nil {
		return nil, err
	}

	p := &fileNotesProvider{path: path, logw: logw, notes: notes}
	if info, err := os.Stat(path); err == nil {
		p.modTime, p.size = info.ModTime(), info.Size()
	}

	util.Infof(logw, "parsed %v notes", note.Count(notes))

	return p, nil
}

// Notes returns the parsed notes, reparsing the file first if it changed since
// the last time it was read.
func (p *fileNotesProvider) Notes() []note.Note {
	p.mu.Lock()
	defer p.mu.Unlock()

	info, err := os.Stat(p.path)
	if err != nil {
		util.Warnf(p.logw, "failed to check %v for changes, serving previously parsed notes: %v", p.path, err)
		return p.notes
	}

	if info.ModTime().Equal(p.modTime) && info.Size() == p.size {
		util.Debugf(p.logw, "%v has not changed", p.path)
		return p.notes
	}

	notes, err := note.ParseFile(p.path)
	if err != nil {
		util.Warnf(p.logw, "failed to reparse %v, serving previously parsed notes: %v", p.path, err)
		return p.notes
	}

	p.notes, p.modTime, p.size = notes, info.ModTime(), info.Size()
	util.Infof(p.logw, "%v changed, reparsed %v notes", p.path, note.Count(notes))

	return p.notes
}
