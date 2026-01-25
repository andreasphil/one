package service

import (
	"time"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
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

func (n NotesService) FindNotesByDate(from time.Time, to time.Time) []lib.Note {
	return util.Filter(n.index, func(n lib.Note) bool {
		return (n.Date.After(from) || n.Date.Equal(from)) && (n.Date.Before(to) || n.Date.Equal(to))
	})
}
