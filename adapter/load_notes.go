package adapter

import "github.com/andreasphil/one/lib"

type NotesLoader interface {
	LoadNotes() []lib.Note
}
