package adapter

import "github.com/andreasphil/one/lib/note"

type NotesProvider interface {
	Notes() []note.Note
}
