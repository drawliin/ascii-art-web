
# Ascii Art Web

This project is a web version of the ascii-art program.
It provides a graphical web interface that allows users to enter text, choose a banner style (standard, shadow, or thinkertoy), and generate ASCII art directly in the browser.

The application is built with a Go HTTP server and HTML templates.
It processes user input, converts it into ASCII art using predefined banner files, and displays the result on the web page.


## Authors

- [@Houssam Eddine Hamouich](https://www.github.com/drawliin)
- [@Mohammed belhoussine drissi](https://www.github.com/DissonantVoid)
## Run Locally

Clone the project

```bash
  git clone https://github.com/drawliin/ascii-art-web
```

Go to the project directory

```bash
  cd ascii-art-web
```

Start the server

```bash
  go run server.go
```

Open your browser and go to:

```bash
  http://localhost:8080
```

## Implementation details (Algorithm)

1. The server starts and listens on port 8080.

2. When a client sends a GET / request:

- The server loads the HTML template.

- It renders the main page and displays any previous result.

3. When the client sends a POST /ascii-art request:

- The server validates the HTTP method.

- It reads the form data (text and selected banner).

- It removes carriage return characters (\r) from the input text and validates that the text is not empty.

- It loads the corresponding banner file (standard.txt, shadow.txt, or thinkertoy.txt).

- It splits the banner file into a 2D structure representing characters.

- It processes each character of the input text and builds the ASCII art line by line.

- It stores the result and sends it back to be displayed (either by redirecting or rendering directly).

4. If an error occurs:

- 400 Bad Request is returned for invalid requests.

- 404 Not Found is returned if a banner or template is missing.

- 500 Internal Server Error is returned for unexpected failures.