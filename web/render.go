package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed templates
var templatesFS embed.FS

type data[T any] struct {
	Title      string
	CurrentURL string
	Data       T
}

type renderFunc[T any] func(http.ResponseWriter, data[T]) error

func newRenderFunc[T any](name string) renderFunc[T] {
	helpers := template.FuncMap{
		// "hasPrefix": strings.HasPrefix,

		// "withQuery": func(path string, query string) template.URL {
		// 	if query == "" {
		// 		return template.URL(path)
		// 	}
		// 	u := url.URL{Path: path, RawQuery: query}
		// 	return template.URL(u.String())
		// },

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

	return func(w http.ResponseWriter, data data[T]) error {
		if err := t.ExecuteTemplate(w, name, data); err != nil {
			return fmt.Errorf("failed to render page template: %w", err)
		}

		return nil
	}
}
