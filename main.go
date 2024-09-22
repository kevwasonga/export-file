package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"ascii-art-web/asciiArt"
)

type Error struct {
	ErrNo   int
	ErrText string
}

// Resizing
// file deletion
// error pages
func main() {
	if len(os.Args) != 1 {
		fmt.Printf("usage: go run main.go" + "\n")
		return
	}

	url := ":8084"

	http.HandleFunc("/", pathHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)
	http.HandleFunc("/error", errorHandler)

	fmt.Printf("  Server listening on http://localhost%s...", url)
	log.Fatal(http.ListenAndServe(url, nil))
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	code, err := strconv.Atoi(r.URL.Query().Get("error"))
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	msg := "Page not found"
	if code == 400 {
		msg = "Bad request"
		w.WriteHeader(http.StatusBadRequest)
	} else if code == 500 {
		msg = "Internal server error"
		w.WriteHeader(http.StatusInternalServerError)

	} else {
		code = 404
		w.WriteHeader(http.StatusNotFound)
	}
	tmpl, _ := template.ParseFiles("./static/errors.html")
	tmpl.ExecuteTemplate(w, "errors.html", Error{code, msg})
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path

	if filePath == "/" {
		http.ServeFile(w, r, "./static")
	} else if strings.HasPrefix(filePath, "/static/") {
		http.ServeFile(w, r, "."+filePath)
	} else {
		w.WriteHeader(http.StatusNotFound)
		tmpl, _ := template.ParseFiles("./static/errors.html")
		tmpl.ExecuteTemplate(w, "errors.html", Error{404, "Page Not Found"})
	}
}

func generateASCIIArtResponse(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method == http.MethodPost {
		// Struct for decoding JSON input
		var request struct {
			Banner string `json:"banner"`
			Input  string `json:"input"`
		}

		// Decode the request body
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return "", err
		}

		return processASCIIArt(w, request.Banner, request.Input)
	} else if r.Method == http.MethodGet {
		// Handle GET request: extract query parameters
		banner := r.URL.Query().Get("banner")
		input := r.URL.Query().Get("input")

		if banner == "" || input == "" {
			return "", fmt.Errorf("missing banner or input parameters")
		}

		return processASCIIArt(w, banner, input)
	}
	return "", fmt.Errorf("Method Not Allowed")
}

func processASCIIArt(w http.ResponseWriter, banner, input string) (string, error) {
	// Load the appropriate banner file
	fileName := asciiArt.BannerFile(banner)
	bannerMap, err := asciiArt.LoadBannerMap(w, fileName)
	if err != nil {
		return "", err
	}

	// Generate the ASCII art
	return generateASCIIArt(w, input, bannerMap)
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// Call the function to generate ASCII art
	response, err := generateASCIIArtResponse(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	// Set content type and write response
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Allow the download to be handled with GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		tmpl, _ := template.ParseFiles("./static/errors.html")
		tmpl.ExecuteTemplate(w, "errors.html", Error{405, "Method Not Allowed"})
		return
	}

	// Generate ASCII art using a helper function
	asciiArt, err := generateASCIIArtResponse(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpl, _ := template.ParseFiles("./static/errors.html")
		tmpl.ExecuteTemplate(w, "errors.html", Error{500, "Internal Server Error"})
		return
	}

	// Set headers for file download
	fileName := "ascii_art.txt"
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/plain")

	// Write ASCII art to the response
	w.Write([]byte(asciiArt))
}

// generateASCIIArt generates the ASCII art from input and bannerMap.
// Logs an error using http.Error if input contains characters outside ASCII limits.
func generateASCIIArt(w http.ResponseWriter, input string, bannerMap map[int][]string) (string, error) {
	var str strings.Builder
	lines := make([]string, 8)
	input = strings.ReplaceAll(input, "\r", "\n")
	arr := strings.Split(input, "\n")
	for _, line := range arr {
		for _, char := range line {

			banner, exists := bannerMap[int(char)]
			if !exists {
				w.WriteHeader(http.StatusBadRequest)
				// http.Error(w, fmt.Sprintf("400 - Character '%c' not found in banner map", char), http.StatusBadRequest)
				tmpl, _ := template.ParseFiles("./static/errors.html")
				tmpl.ExecuteTemplate(w, "errors.html", Error{400, "Bad request"})
				return "", fmt.Errorf("character '%c' not found in banner map", char)
			}
			for i := 0; i < 8; i++ {
				lines[i] += banner[i]
			}
		}
		str.WriteString(strings.Join(lines, "\n"))
		str.WriteString("\n")
		lines = make([]string, 8)
	}
	return str.String(), nil
}
