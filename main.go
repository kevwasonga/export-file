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

	url := ":8080"

	http.HandleFunc("/", pathHandler)
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

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		tmpl, _ := template.ParseFiles("./static/errors.html")
		tmpl.ExecuteTemplate(w, "errors.html", Error{405, "Method Not allowed"})
		return
	}

	var request struct {
		Banner     string `json:"banner"`
		Input      string `json:"input"`
		SaveToFile bool   `json:"saveToFile"` // New field to indicate whether to save to file
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}

	fileName := asciiArt.BannerFile(request.Banner)

	bannerMap, err := asciiArt.LoadBannerMap(w, fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, "/error", http.StatusInternalServerError)
		return
	}

	response, httpErr := generateASCIIArt(w, request.Input, bannerMap)
	if httpErr != nil {
		return
	}

	if request.SaveToFile {
		// Save the response to a file
		err := saveToFile(response, "output.txt")
		if err != nil {
			log.Println("Failed to save ASCII art to file:", err)
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

func saveToFile(data string, filename string) error {
	return os.WriteFile(filename, []byte(data), 0o644)
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
