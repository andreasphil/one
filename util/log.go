package util

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	blue   = "\033[34m"
	yellow = "\033[33m"
	red    = "\033[31m"
	green  = "\033[32m"
	gray   = "\033[90m"

	// 256 color palette colors
	plum = "\033[38;5;176m"
)

const face = "(っᵔᴗᵔ)っ"

func style(color string, text string) string {
	return color + text + reset
}

// Banner writes text to w in a violet rounded frame.
func Banner(w io.Writer, text string) {
	rule := strings.Repeat("─", utf8.RuneCountInString(text)+4)

	fmt.Fprintf(w, "\n %v\n", style(plum, "╭"+rule+"╮"))
	fmt.Fprintf(w, " %v  %v  %v  %v\n", style(plum, "│"), style(bold, text), style(plum, "│"), style(gray, face))
	fmt.Fprintf(w, " %v\n\n", style(plum, "╰"+rule+"╯"))
}

func logf(w io.Writer, marker string, format string, v ...any) {
	fmt.Fprintf(w, " "+marker+" "+format+"\n", v...)
}

// Infof writes an informational message to w, formatted according to format,
// and prefixed with a gray arrow.
func Infof(w io.Writer, format string, v ...any) {
	logf(w, style(gray, "→"), format, v...)
}

// Warnf writes a warning message to w, formatted according to format, and
// prefixed with a yellow triangle.
func Warnf(w io.Writer, format string, v ...any) {
	logf(w, style(yellow, "△"), format, v...)
}

// Errorf writes an error message to w, formatted according to format, and
// prefixed with a red cross.
func Errorf(w io.Writer, format string, v ...any) {
	logf(w, style(red, "✗"), format, v...)
}

// Successf writes a success message to w, formatted according to format, and
// prefixed with a green checkmark.
func Successf(w io.Writer, format string, v ...any) {
	logf(w, style(green, "✓"), format, v...)
}

// HTTPErrorf writes an error message, formatted according to format, to errw,
// and sends it as an HTTP error response with the given status code.
func HTTPErrorf(errw io.Writer, w http.ResponseWriter, status int, format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Errorf(errw, "%v %v", status, message)
	http.Error(w, message, status)
}
