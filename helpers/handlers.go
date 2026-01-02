package helpers

import (
	"bytes"
	"net/http"
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
	if err := r.ParseForm(); err != nil {
		ErrorPage(w, http.StatusBadRequest)
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

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	// filling result with the output to print it in root "/"
	data.Art, err = GenerateArt(data.UserInput, data.Font)

	if err != nil {
		if ContainsUnsupportedChars(err) {
			data.ErrorMsg = err.Error()
		} else {
			ErrorPage(w, http.StatusInternalServerError)
			return
		}
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		ErrorPage(w, http.StatusInternalServerError)
		return
	}
	if data.ErrorMsg != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	buf.WriteTo(w)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		ErrorPage(w, http.StatusBadRequest)
		return
	}

	input := strings.ReplaceAll(r.FormValue("input"), "\r", "")
	banner := r.FormValue("banner")

	if input == "" || len(input) > 2000 {
		ErrorPage(w, http.StatusBadRequest)
		return
	}

	if !ValidFont(banner) {
		ErrorPage(w, http.StatusBadRequest)
		return
	}

	art, err := GenerateArt(input, banner)
	if err != nil {
		if ContainsUnsupportedChars(err) {
			ErrorPage(w, http.StatusBadRequest)
			return
		}
		ErrorPage(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.txt"`)
	w.Write([]byte(art))
}
