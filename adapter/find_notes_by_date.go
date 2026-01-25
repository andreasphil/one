package adapter

import (
	"time"

	"github.com/andreasphil/one/lib"
)

type NotesByDateFinder interface {
	FindNotesByDate(from time.Time, to time.Time) []lib.Note
}
