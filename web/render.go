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
	CurrentUrl string
	Data       T
}

type nilPage data[struct{}]

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
	}

	t := template.Must(template.New("").Funcs(helpers).ParseFS(templatesFS, "templates/shared/*.html"))
	// template.Must(t.ParseFS(templatesFS, "templates/components/*.html"))
	// template.Must(t.ParseFS(templatesFS, "templates/icons/*.svg"))

	template.Must(t.ParseFS(templatesFS, fmt.Sprintf("templates/%v", name)))

	return func(w http.ResponseWriter, data data[T]) error {
		return t.ExecuteTemplate(w, name, data)
	}
}
