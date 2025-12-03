package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"web/helpers"
)

func main() {
	http.HandleFunc("/ascii", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Route", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	//////////////
	input := r.Form["input"][0]
	fileName := r.Form["banner"][0]

	bytes, err := os.ReadFile(fmt.Sprintf("%s.txt", fileName))
	if err != nil {
		fmt.Fprintf(w,"error: %v\n", err)
		return
	}

	// Split into 2d Slice
	arr := helpers.Split2D(string(bytes))

	// Split the input by NewLine
	lines := strings.Split(input, "\\n")

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
	fmt.Fprint(w, res.String())
	//////////////
	fmt.Fprintf(w, "\n\nYour Data %q", r.PostForm)
}
