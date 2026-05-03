package web

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

type data[T any] struct {
	Title         string
	CurrentUrl    string
	CurrentSearch string
	CurrentFilter string
	CurrentQuery  string
	Data          T
}

type nilPage data[struct{}]

type renderFunc[T any] func(http.ResponseWriter, data[T]) error

func newRenderFunc[T any](name string) renderFunc[T] {
	helpers := template.FuncMap{
		"hasPrefix": strings.HasPrefix,

		"withQuery": func(path string, query string) template.URL {
			if query == "" {
				return template.URL(path)
			}
			u := url.URL{Path: path, RawQuery: query}
			return template.URL(u.String())
		},
	}

	t := template.Must(template.New("").Funcs(helpers).ParseGlob("./web/templates/shared/*.html"))
	template.Must(t.ParseGlob("./web/templates/components/*.html"))
	template.Must(t.ParseGlob("./web/templates/icons/*.svg"))

	fullName := fmt.Sprintf("./web/templates/%v", name)
	template.Must(t.ParseFiles(fullName))

	return func(w http.ResponseWriter, data data[T]) error {
		return t.ExecuteTemplate(w, name, data)
	}
}
