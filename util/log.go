package util

import (
	"fmt"
	"log"
	"net/http"
)

const (
	reset  = "\033[0m"
	blue   = "\033[34m"
	yellow = "\033[33m"
	red    = "\033[31m"
	green  = "\033[32m"
	bgRed  = "\033[41m"
	black  = "\033[30m"
)

func style(color string, text string) string {
	return color + text + reset
}

// Infof logs an informational message, formatted according to format, and
// prefixed with a blue arrow.
func Infof(format string, v ...any) {
	log.Printf(style(blue, "→")+" "+format, v...)
}

// Warnf logs a warning message, formatted according to format, and prefixed
// with a yellow triangle.
func Warnf(format string, v ...any) {
	log.Printf(style(yellow, "△")+" "+format, v...)
}

// Errorf logs an error message, formatted according to format, and prefixed
// with a red cross.
func Errorf(format string, v ...any) {
	log.Printf(style(red, "✗")+" "+format, v...)
}

// Okf logs a success message, formatted according to format, and prefixed with
// a green checkmark.
func Okf(format string, v ...any) {
	log.Printf(style(green, "✓")+" "+format, v...)
}

// HttpErrorf logs an error message formatted according to format, and writes it
// as an HTTP error response with the given status code.
func HttpErrorf(w http.ResponseWriter, status int, format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Errorf("%v %v", status, message)
	http.Error(w, message, status)
}
