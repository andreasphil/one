package web

import (
	"fmt"
	"html/template"
	"net/http"
)

type data[T any] struct {
	Title      string
	CurrentUrl string
	Data       T
}

type nilPage data[struct{}]

type renderFunc[T any] func(http.ResponseWriter, data[T]) error

func newRenderFunc[T any](name string) renderFunc[T] {
	t := template.Must(template.ParseGlob("./web/templates/shared/*.html"))
	template.Must(t.ParseGlob("./web/templates/icons/*.svg"))

	fullName := fmt.Sprintf("./web/templates/%v", name)
	template.Must(t.ParseFiles(fullName))

	return func(w http.ResponseWriter, data data[T]) error {
		return t.ExecuteTemplate(w, name, data)
	}
}
