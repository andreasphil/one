package adapter

import "github.com/andreasphil/one/lib"

type NotesProvider interface {
	Notes() []lib.Note
}
