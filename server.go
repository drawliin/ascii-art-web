package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"web/helpers"
)

var result string

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Render the main HTML template and inject the generated ASCII art.
		tmpl := template.Must(template.ParseFiles("template/index.html"))
		
		// We pass "Art" into the template so index.html can display the result.
		tmpl.Execute(w, map[string]string{
			"Art": result,
		})
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/ascii", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Route", http.StatusMethodNotAllowed)
		return
	}

	// tmpl := template.Must(template.ParseFiles("template/ascii-art-web.html"))

	r.ParseForm()
	input := strings.ReplaceAll(r.Form["input"][0], "\r", "")
	fileName := r.Form["banner"][0]

	bytes, err := os.ReadFile(fmt.Sprintf("%s.txt", fileName))
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		return
	}

	// Split into 2d Slice
	arr := helpers.Split2D(string(bytes))

	// Split the input by NewLine
	lines := strings.Split(input, "\n")

	// check trailing empty string
	if len([]rune(input)) > 1 && helpers.ContainOnlyNewLines(input) {
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
					fmt.Fprintf(w, "error: unsupported character: %q\n", c)
					return
				}
			}
			res.WriteRune('\n')
		}
	}
	// filling result with the output to print it in root "/"
	result = res.String()

	// return to root "/" and show data
	http.Redirect(w, r, "/", http.StatusSeeOther)

	// tmpl.Execute(w, data)
	// fmt.Fprintf(w, "%s\n", res.String())
	// fmt.Fprintf(w, "\n\nYour Data %q", r.PostForm)
}
