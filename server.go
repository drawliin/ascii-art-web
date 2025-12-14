package main

import (
	"bytes"
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
			http.Error(w, "Error: 405 Not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Render the main HTML template
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Error: 500 InternalServerError", http.StatusInternalServerError)
			return
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, PageData{
			UserInput: "",
			Font:      "standard",
			Art:       "",
		})
		if err != nil {
			http.Error(w, "Error: 500 InternalServerError", http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)
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

	// parse data
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

	bytesF, err := os.ReadFile(fmt.Sprintf("%s.txt", data.Font))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: Internal Server Error"))
		return
	}

	// Split into 2d Slice
	fontTxt := strings.ReplaceAll(string(bytesF), "\r", "")
	arr := helpers.Split2D(fontTxt)

	// Split the input by NewLine
	lines := strings.Split(data.UserInput, "\n")

	// check trailing empty string
	if len(lines) > 1 && helpers.ContainOnlyNewLines(lines) {
		lines = lines[:len(lines)-1]
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error: 500 InternalServerError", http.StatusInternalServerError)
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
					data.ErrorMsg = fmt.Sprintf("Error: 400 unsupported character: %q\n", c)
					w.WriteHeader(http.StatusBadRequest)
					tmpl.Execute(w, data)
					return
				}
			}
			res.WriteRune('\n')
		}
	}
	// filling result with the output to print it in root "/"
	data.Art = res.String()

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		http.Error(w, "Error: 500 InternalServerError", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func validFont(s string) bool {
	switch s {
	case "standard", "shadow", "thinkertoy":
		return true
	default:
		return false
	}
}
