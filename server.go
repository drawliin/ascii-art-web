package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"web/helpers"
)

var userInput string
var font string = "standard"
var art string = ""

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Error: 404 Not found", http.StatusNotFound)
			return
		}

		// Render the main HTML template and inject the generated ASCII art.
		tmpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			http.Error(w, "Error: 404 Not found", http.StatusNotFound)
			return
		}

		// We pass "Art" into the template so index.html can display the result.
		tmpl.Execute(w, map[string]string{
			"Art":       art,
			"PrevInput": userInput,
			"Font":      font,
		})
	})

	// serve static files like (css || js) so the html can access them if needed
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/ascii-art", handler)

	fmt.Printf("%s\n", "Server Starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error: Bad request", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	userInput = strings.ReplaceAll(r.Form["input"][0], "\r", "")
	if userInput == "" {
		http.Error(w, "Error: Bad request", http.StatusBadRequest)
		return
	}

	if font = r.Form["banner"][0]; !validFont(font) {
		http.Error(w, "Error: Not Found ", http.StatusNotFound)
		return
	}

	bytes, err := os.ReadFile(fmt.Sprintf("%s.txt", font))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: Internal Server error"))
		return
	}

	// Split into 2d Slice
	arr := helpers.Split2D(string(bytes))

	// Split the input by NewLine
	lines := strings.Split(userInput, "\n")

	// check trailing empty string
	if len(lines) > 1 && helpers.ContainOnlyNewLines(lines) {
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
				// check if valid and printable ascii
				if c >= 32 && c <= 126 {
					res.WriteString(arr[c-32][j])
				} else {
					art = fmt.Sprintf("error: unsupported character: %q\n", c)
					http.Error(w, "Error: Internal error", http.StatusInternalServerError)
					return
				}
			}
			res.WriteRune('\n')
		}
	}
	// filling result with the output to print it in root "/"
	art = res.String()

	// return to root "/" and show data
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validFont(s string) bool {
	switch s {
	case "standard", "shadow", "thinkertoy":
		return true
	default:
		return false
	}
}
