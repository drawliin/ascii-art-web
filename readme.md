
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

## Implementation Details (Algorithm)

### 1. Server Startup
- The server starts and listens on port **8080**.

### 2. Handling `GET /`
When a client sends a `GET /` request:
- The server loads the HTML template.
- It renders the main page and displays any previously generated result.

### 3. Handling `POST /ascii-art`
When a client sends a `POST /ascii-art` request, the following steps are executed:

1. **Request Validation**
   - The server validates that the HTTP method is correct.

2. **Form Data Processing**
   - Reads the submitted form data:
     - Input text
     - Selected banner style
   - Removes all carriage return characters (`\r`) from the input.
   - Verifies that the input text is not empty.

3. **Banner Loading**
   - Loads the corresponding banner file:
     - `standard.txt`
     - `shadow.txt`
     - `thinkertoy.txt`

4. **Banner Parsing**
   - Splits the banner file into a **2D data structure** representing character patterns.

5. **ASCII Art Generation**
   - Iterates over each character of the input text.
   - Constructs the ASCII art **line by line** using the parsed banner data.

6. **Response Handling**
   - Stores the generated result.
   - Sends the result back to the client by:
     - Redirecting, **or**
     - Directly rendering the page.

### 4. Error Handling

The server responds with appropriate HTTP status codes in case of failure:

- `200 Okay` – No errors
- `400 Bad Request` – Invalid or malformed requests
- `404 Not Found` – Missing banner or template files
- `500 Internal Server Error` – Unexpected server-side errors
- `405 Method Not Allowed` – Inappropriate method used