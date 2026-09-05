package mapper_test

import (
	"encoding/json"
	"testing"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
	"github.com/andreasphil/one/web/mapper"
	"github.com/google/go-cmp/cmp"
)

func TestToTags(t *testing.T) {
	notes := []note.Note{
		{Title: "Root 1", Tags: util.NewSetFrom([]note.Tag{"#work", "#Idea"})},
		{
			Title: "Root 2",
			Tags:  util.NewSetFrom([]note.Tag{"#work"}),
			Children: []note.Note{
				{Title: "Child 1", Tags: util.NewSetFrom([]note.Tag{"#recipe"})},
			},
		},
	}

	t.Run("collects unique tags of notes and children without the leading #", func(t *testing.T) {
		result := mapper.ToTags(notes)

		expected := []string{"Idea", "recipe", "work"}

		if !cmp.Equal(result, expected) {
			t.Errorf("expected %+v, got %+v", expected, result)
		}
	})

	t.Run("returns an empty slice for no notes", func(t *testing.T) {
		for _, notes := range [][]note.Note{nil, {}} {
			result := mapper.ToTags(notes)

			if result == nil {
				t.Fatalf("expected a non-nil slice, got nil")
			}

			if len(result) != 0 {
				t.Errorf("expected an empty slice, got %+v", result)
			}
		}
	})

	t.Run("serializes an empty slice to an empty JSON array", func(t *testing.T) {
		result, err := json.Marshal(mapper.ToTags(nil))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if string(result) != "[]" {
			t.Errorf("expected %q, got %q", "[]", string(result))
		}
	})
}
