package service

import (
	"github.com/andreasphil/one/lib"
)

type NotesService struct {
	index []lib.Note
}

func NewNotesService(path string) (*NotesService, error) {
	notes, err := lib.ParseFile(path)
	if err != nil {
		return nil, err
	}

	service := NotesService{
		index: notes,
	}

	return &service, nil
}

func (n NotesService) LoadNotes() []lib.Note {
	return n.index
}

func (n NotesService) LoadNote(slug string) (lib.Note, bool) {
	return lib.FindRecursive(n.index, slug)
}
