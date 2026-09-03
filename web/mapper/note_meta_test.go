package mapper_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/web/mapper"
	"github.com/google/go-cmp/cmp"
)

func TestNewNoteMeta(t *testing.T) {
	type testcase struct {
		name     string
		note     note.Note
		expected mapper.NoteMeta
	}

	testcases := []testcase{
		{
			name:     "maps title and slug",
			note:     note.Note{Title: "Hello World"},
			expected: mapper.NoteMeta{Title: "Hello World", Slug: "hello-world"},
		},
		{
			name:     "maps daily note",
			note:     note.Note{Title: "01.02.2026", Date: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
			expected: mapper.NoteMeta{Title: "01.02.2026", Slug: "2026-02-01"},
		},
		{
			name:     "maps empty note",
			note:     note.Note{},
			expected: mapper.NoteMeta{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := mapper.NewNoteMeta(tc.note)

			if !cmp.Equal(result, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestToNoteMeta(t *testing.T) {
	date := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	notes := []note.Note{
		{Title: "Root 1"},
		{
			Title: "01.02.2026",
			Date:  date,
			Children: []note.Note{
				{Title: "Child 1", Date: date},
				{
					Title: "Child 2",
					Date:  date,
					Children: []note.Note{
						{Title: "Grandchild 1", Date: date},
					},
				},
			},
		},
		{Title: "Root 3"},
	}

	t.Run("flattens notes and children in depth-first, pre-order", func(t *testing.T) {
		result := mapper.ToNoteMeta(notes)

		expected := []mapper.NoteMeta{
			{Title: "Root 1", Slug: "root-1"},
			{Title: "01.02.2026", Slug: "2026-02-01"},
			{Title: "Child 1", Slug: "2026-02-01-child-1"},
			{Title: "Child 2", Slug: "2026-02-01-child-2"},
			{Title: "Grandchild 1", Slug: "2026-02-01-grandchild-1"},
			{Title: "Root 3", Slug: "root-3"},
		}

		if !cmp.Equal(result, expected) {
			t.Errorf("expected %+v, got %+v", expected, result)
		}
	})

	t.Run("returns an empty slice for no notes", func(t *testing.T) {
		for _, notes := range [][]note.Note{nil, {}} {
			result := mapper.ToNoteMeta(notes)

			if result == nil {
				t.Fatalf("expected a non-nil slice, got nil")
			}

			if len(result) != 0 {
				t.Errorf("expected an empty slice, got %+v", result)
			}
		}
	})

	t.Run("serializes an empty slice to an empty JSON array", func(t *testing.T) {
		result, err := json.Marshal(mapper.ToNoteMeta(nil))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if string(result) != "[]" {
			t.Errorf("expected %q, got %q", "[]", string(result))
		}
	})
}
