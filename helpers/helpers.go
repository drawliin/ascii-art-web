package helpers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"html/template"
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

var Cache = make(map[string]*template.Template)

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
	tmpl := Cache["error.html"]

	w.WriteHeader(status)

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, AnError{
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

func GenerateArt(input, font string) (string, error) {
	bytesF, err := os.ReadFile(fmt.Sprintf("%s.txt", font))
	if err != nil {
		return "", err
	}

	fontTxt := strings.ReplaceAll(string(bytesF), "\r", "")
	arr := Split2D(fontTxt)
	lines := strings.Split(input, "\n")

	if len(lines) > 1 && ContainOnlyNewLines(lines) {
		lines = lines[:len(lines)-1]
	}

	var res strings.Builder
	for _, line := range lines {
		if line == "" {
			res.WriteRune('\n')
			continue
		}
		for j := range 8 {
			for _, c := range line {
				if c < ' ' || c > '~' {
					return "", fmt.Errorf("unsupported character: %q", c)
				}
				res.WriteString(arr[c-' '][j])
			}
			res.WriteRune('\n')
		}
	}
	return res.String(), nil
}

func ContainsUnsupportedChars(err error) bool {
	return strings.Contains(err.Error(), "unsupported character")
}
