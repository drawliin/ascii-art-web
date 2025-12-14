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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			errtmpl, tmperr := template.ParseFiles("templates/error.html")
			if tmperr != nil {
				w.WriteHeader(500)
				w.Write([]byte("500 Internal Server Error"))
				return
			}
			w.WriteHeader(404)
			errtmpl.Execute(w, AnError{
				Code:    404,
				Message: "Not found",
			})
			return
		}

		if r.Method != http.MethodGet {
			errtmpl, tmperr := template.ParseFiles("templates/error.html")
			if tmperr != nil {
				w.WriteHeader(500)
				w.Write([]byte("500 Internal Server Error"))
				return
			}
			w.WriteHeader(405)
			errtmpl.Execute(w, AnError{
				Code:    405,
				Message: "Not allowed",
			})
			return
		}

		// Render the main HTML template
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			w.WriteHeader(500)
			errtmpl, tmperr := template.ParseFiles("templates/error.html")
			if tmperr != nil {
				w.Write([]byte("500 Internal Server Error"))
				return
			}
			errtmpl.Execute(w, AnError{
				Code:    500,
				Message: "Internal Server Error",
			})
			return
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, PageData{
			UserInput: "",
			Font:      "standard",
			Art:       "",
		})
		if err != nil {
			w.WriteHeader(500)
			errtmpl, tmperr := template.ParseFiles("templates/error.html")
			if tmperr != nil {
				w.Write([]byte("500 Internal Server Error"))
				return
			}
			errtmpl.Execute(w, AnError{
				Code:    500,
				Message: "Internal Server Error",
			})
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
		errtmpl, tmperr := template.ParseFiles("templates/error.html")
		if tmperr != nil {
			w.WriteHeader(500)
			w.Write([]byte("500 Internal Server Error"))
			return
		}
		w.WriteHeader(405)
		errtmpl.Execute(w, AnError{
			Code:    405,
			Message: "Not Allowed",
		})
		return
	}

	// parse data
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		errtmpl, tmperr := template.ParseFiles("templates/error.html")
		if tmperr != nil {
			w.Write([]byte("500 Internal Server Error"))
			return
		}
		errtmpl.Execute(w, AnError{
			Code:    500,
			Message: "Internal Server Error",
		})
		return
	}
	data := PageData{}
	data.UserInput = strings.ReplaceAll(r.FormValue("input"), "\r", "")
	if data.UserInput == "" {
		w.WriteHeader(500)
		errtmpl, tmperr := template.ParseFiles("templates/error.html")
		if tmperr != nil {
			w.Write([]byte("500 Internal Server Error"))
			return
		}
		errtmpl.Execute(w, AnError{
			Code:    500,
			Message: "Internal Server Error",
		})
		return
	}

	if data.Font = r.FormValue("banner"); !validFont(data.Font) {
		errtmpl, tmperr := template.ParseFiles("templates/error.html")
		if tmperr != nil {
			w.WriteHeader(500)
			w.Write([]byte("500 Internal Server Error"))
			return
		}
		w.WriteHeader(400)
		errtmpl.Execute(w, AnError{
			Code:    400,
			Message: "Bad Request",
		})
		return
	}

	bytesF, err := os.ReadFile(fmt.Sprintf("%s.txt", data.Font))
	if err != nil {
		w.WriteHeader(500)
		errtmpl, tmperr := template.ParseFiles("templates/error.html")
		if tmperr != nil {
			w.Write([]byte("500 Internal Server Error"))
			return
		}
		errtmpl.Execute(w, AnError{
			Code:    500,
			Message: "Internal Server Error",
		})
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
		w.WriteHeader(500)
		errtmpl, tmperr := template.ParseFiles("templates/error.html")
		if tmperr != nil {
			w.Write([]byte("500 Internal Server Error"))
			return
		}
		errtmpl.Execute(w, AnError{
			Code:    500,
			Message: "Internal Server Error",
		})
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
		w.WriteHeader(500)
		errtmpl, tmperr := template.ParseFiles("templates/error.html")
		if tmperr != nil {
			w.Write([]byte("500 Internal Server Error"))
			return
		}
		errtmpl.Execute(w, AnError{
			Code:    500,
			Message: "Internal Server Error",
		})
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
