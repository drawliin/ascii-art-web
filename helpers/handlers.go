package helpers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorPage(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		ErrorPage(w, http.StatusMethodNotAllowed)
		return
	}

	// Render the main HTML template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, PageData{
		UserInput: "",
		Font:      "standard",
		Art:       "",
	})
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func AsciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(w, http.StatusMethodNotAllowed)
		return
	}

	// parse data
	err := r.ParseForm()
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	data := PageData{}
	data.UserInput = strings.ReplaceAll(r.FormValue("input"), "\r", "")
	if data.UserInput == "" || len(data.UserInput) > 2000 {
		ErrorPage(w, http.StatusBadRequest)
		return
	}

	if data.Font = r.FormValue("banner"); !ValidFont(data.Font) {
		ErrorPage(w, http.StatusBadRequest)
		return
	}

	bytesF, err := os.ReadFile(fmt.Sprintf("%s.txt", data.Font))
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	// Split into 2d Slice
	fontTxt := strings.ReplaceAll(string(bytesF), "\r", "")
	arr := Split2D(fontTxt)

	// Split the input by NewLine
	lines := strings.Split(data.UserInput, "\n")

	// check trailing empty string
	if len(lines) > 1 && ContainOnlyNewLines(lines) {
		lines = lines[:len(lines)-1]
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
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
					data.ErrorMsg = fmt.Sprintf("Unsupported character: %q\n", c)
					w.WriteHeader(http.StatusBadRequest)

					var buf bytes.Buffer
					err = tmpl.Execute(&buf, data)
					if err != nil {
						ErrorPage(w, http.StatusInternalServerError)
						return
					}
					buf.WriteTo(w)

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
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(w, http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	art := r.FormValue("art")
	if art == "" {
		ErrorPage(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/txt")
	w.Header().Set("Content-Disposition", "attachment; filename=\"file.txt\"")
	w.Write([]byte(art))
}
