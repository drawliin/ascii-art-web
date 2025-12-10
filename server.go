package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"web/helpers"
)

type PageData struct {
	UserInput string
	Font      string
	Art       string
	ErrorMsg  string
}

const port = "8080"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Error: 404 Not found", http.StatusNotFound)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Error: Not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Render the main HTML template
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Error: 404 Not found", http.StatusNotFound)
			return
		}

		tmpl.Execute(w, PageData{
			UserInput: "",
			Font:      "standard",
			Art:       "",
		})
	})

	// serve static files like (css || js) so the html can access them if needed
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/ascii-art", handler)

	fmt.Printf("Server Starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error: Not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error: Internal Server error", http.StatusInternalServerError)
		return
	}
	data := PageData{}
	data.UserInput = strings.ReplaceAll(r.FormValue("input"), "\r", "")
	if data.UserInput == "" {
		http.Error(w, "Error: Bad Request", http.StatusBadRequest)
		return
	}

	if data.Font = r.FormValue("banner"); !validFont(data.Font) {
		http.Error(w, "Error: Bad Request", http.StatusBadRequest)
		return
	}

	bytes, err := os.ReadFile(fmt.Sprintf("%s.txt", data.Font))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: Internal Server Error"))
		return
	}

	// Split into 2d Slice
	arr := helpers.Split2D(string(bytes))

	// Split the input by NewLine
	lines := strings.Split(data.UserInput, "\n")

	// check trailing empty string
	if len(lines) > 1 && helpers.ContainOnlyNewLines(lines) {
		lines = lines[:len(lines)-1]
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error: 404 Not found", http.StatusNotFound)
		return
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
				if c >= ' ' && c <= '~' {
					res.WriteString(arr[c-' '][j])
				} else {
					data.ErrorMsg = fmt.Sprintf("error: unsupported character: %q\n", c)
					tmpl.Execute(w, data)
					return
				}
			}
			res.WriteRune('\n')
		}
	}
	// filling result with the output to print it in root "/"
	data.Art = res.String()
	tmpl.Execute(w, data)
}

func validFont(s string) bool {
	switch s {
	case "standard", "shadow", "thinkertoy":
		return true
	default:
		return false
	}
}
