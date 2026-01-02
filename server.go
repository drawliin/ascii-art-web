package main

import (
	"fmt"
	"net/http"
	"os"
	"web/helpers"
)

const port = "8080"

func main() {
	http.HandleFunc("/", helpers.RootHandler)
	http.HandleFunc("/ascii-art", helpers.AsciiHandler)

	// file server for /static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/static/" {
			// prevent spying
			helpers.ErrorPage(w, http.StatusNotFound)
			return
		}

		path := "." + r.URL.Path
		_, err := os.Stat(path)
		if err != nil {
			helpers.ErrorPage(w, http.StatusNotFound)
			return
		}
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	}))

	http.HandleFunc("/download", helpers.DownloadHandler)

	fmt.Printf("Server Starting on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
