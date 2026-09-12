package web

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/andreasphil/one/lib/note"
)

//go:embed templates
var templatesFS embed.FS

type data[T any] struct {
	CurrentURL string
	NotesMeta  []noteMeta
	Notes      []note.Note
	Tags       []string

	Title string
	Data  T
}

type renderFunc[T any] func(http.ResponseWriter, *http.Request, data[T]) error

func newRenderFunc[T any](provider NotesProvider, name string) renderFunc[T] {
	helpers := template.FuncMap{
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict: odd number of arguments")
			}

			d := make(map[string]any, len(values)/2)

			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict: keys must be strings")
				}

				d[key] = values[i+1]
			}

			return d, nil
		},
	}

	t := template.Must(template.New("").Funcs(helpers).ParseFS(templatesFS, "templates/shared/*.html"))
	template.Must(t.ParseFS(templatesFS, "templates/components/*.html"))
	template.Must(t.ParseFS(templatesFS, "templates/icons/*.svg"))
	template.Must(t.ParseFS(templatesFS, fmt.Sprintf("templates/%v", name)))

	return func(w http.ResponseWriter, r *http.Request, data data[T]) error {
		notes := provider.Notes()
		flat := note.Flat(notes)

		data.CurrentURL = r.URL.Path
		data.NotesMeta = mapToNoteMeta(flat)
		data.Notes = flat
		data.Tags = mapToTags(notes)

		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, name, data); err != nil {
			return fmt.Errorf("failed to render page template: %w", err)
		}
		_, err := buf.WriteTo(w)

		return err
	}
}
