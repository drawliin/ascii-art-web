package helpers

import (
	"bytes"
	"net/http"
	"text/template"
)

type PageData struct {
	UserInput string
	Font      string
	Art       string
	ErrorMsg  string
}

type AnError struct {
	Code    int
	Message string
}

// Splits banner file's content into a slice of ascii representations
// Returns a slice of ascii represented symbols (each of them as a slice of 8 lines/strings)
func Split2D(s string) [][]string {
	arr := [][]string{}
	symbol := []string{}
	wordStart := 0
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if s[wordStart:i] == "" {
				if count == 8 {
					arr = append(arr, symbol)
					count = 0
					symbol = []string{}
				}
			} else {
				symbol = append(symbol, s[wordStart:i])
				count++
			}
			wordStart = i + 1
		}
	}
	if len(symbol) > 0 {
		arr = append(arr, symbol)
	}
	return arr
}

func ContainOnlyNewLines(arr []string) bool {
	for _, c := range arr {
		if c != "" {
			return false
		}
	}
	return true
}

func ErrorPage(w http.ResponseWriter, status int) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// error in the error page *o*
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	w.WriteHeader(status)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, AnError{
		Code:    status,
		Message: http.StatusText(status),
	})
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func ValidFont(s string) bool {
	switch s {
	case "standard", "shadow", "thinkertoy":
		return true
	default:
		return false
	}
}
