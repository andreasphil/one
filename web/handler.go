package web

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/andreasphil/one/util"
)

type httpStatusError struct {
	status int
	err    error
}

func (e *httpStatusError) Error() string {
	return e.err.Error()
}

func (e *httpStatusError) Unwrap() error {
	return e.err
}

func httpStatusErrorf(status int, format string, v ...any) error {
	return &httpStatusError{status: status, err: fmt.Errorf(format, v...)}
}

type handler func(w http.ResponseWriter, r *http.Request) error

func handle(errw io.Writer, h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.Debugf(errw, "%v %v %v", r.Method, r.URL.EscapedPath(), util.FormatNestedMap(r.URL.Query()))

		err := h(w, r)
		if err == nil {
			return
		}

		status := http.StatusInternalServerError
		if statusErr, ok := errors.AsType[*httpStatusError](err); ok {
			status = statusErr.status
		}

		util.Errorf(errw, "[%v] %v", status, err.Error())
		http.Error(w, fmt.Sprintf("%v", err.Error()), status)
	}
}
