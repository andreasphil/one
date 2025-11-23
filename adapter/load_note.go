package adapter

import "github.com/andreasphil/one/lib"

type NoteLoader interface {
	LoadNote(slug string) (lib.Note, bool)
}
