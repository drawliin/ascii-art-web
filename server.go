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
type AnError struct {
	Code    int
	Message string
}

const port = "8080"

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ascii-art", ascciHandler)

	// file server for /static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/static/" {
			// prevent spying
			errorPage(w, http.StatusNotFound)
			return
		}

		path := "." + r.URL.Path
		_, err := os.Stat(path)
		if err != nil {
			errorPage(w, http.StatusNotFound)
			return
		} else {
			http.StripPrefix("/static/", fs).ServeHTTP(w, r)
		}
	}))

	fmt.Printf("Server Starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorPage(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		errorPage(w, http.StatusMethodNotAllowed)
		return
	}

	// Render the main HTML template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		errorPage(w, http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, PageData{
		UserInput: "",
		Font:      "standard",
		Art:       "",
	})
	if err != nil {
		errorPage(w, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func ascciHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorPage(w, http.StatusMethodNotAllowed)
		return
	}

	// parse data
	err := r.ParseForm()
	if err != nil {
		errorPage(w, http.StatusInternalServerError)
		return
	}
	data := PageData{}
	data.UserInput = strings.ReplaceAll(r.FormValue("input"), "\r", "")
	if data.UserInput == "" {
		errorPage(w, http.StatusBadRequest)
		return
	}

	if data.Font = r.FormValue("banner"); !validFont(data.Font) {
		errorPage(w, http.StatusBadRequest)
		return
	}

	bytesF, err := os.ReadFile(fmt.Sprintf("%s.txt", data.Font))
	if err != nil {
		errorPage(w, http.StatusInternalServerError)
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
		errorPage(w, http.StatusInternalServerError)
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
		errorPage(w, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func errorPage(w http.ResponseWriter, status int) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// error in the error page *o*
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	w.WriteHeader(status)
	tmpl.Execute(w, AnError{
		Code:    404,
		Message: "Not found",
	})
}

func validFont(s string) bool {
	switch s {
	case "standard", "shadow", "thinkertoy":
		return true
	default:
		return false
	}
}
