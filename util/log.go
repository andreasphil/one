package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

func Infof(format string, v ...any) {
	log.Printf(style(blue, "INFO")+" "+format, v...)
}

func Warnf(format string, v ...any) {
	log.Printf(style(yellow, "WARN")+" "+format, v...)
}

func Errorf(format string, v ...any) {
	log.Printf(style(red, "ERROR")+" "+format, v...)
}

func Okf(format string, v ...any) {
	log.Printf(style(green, "OK")+" "+format, v...)
}

func Fatalf(format string, v ...any) {
	log.Printf(style(bgRed+black, " FATAL ")+" "+format, v...)
	os.Exit(1)
}

func HttpErrorf(w http.ResponseWriter, status int, format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	Errorf("%v %v", status, message)
	http.Error(w, message, status)
}
