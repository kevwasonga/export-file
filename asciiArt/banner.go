package asciiArt

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

type Error struct {
	ErrNo   int
	ErrText string
}

// LoadBannerMap reads an ASCII banner file and maps each banner to a unique key.

func LoadBannerMap(w http.ResponseWriter, fileName string) (map[int][]string, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// checksum results, to ascertain non of the banner files is tampered with

	// r *http.Request
	results := VerifyBanners()

	for filePath, isValid := range results {
		if isValid {
			//	fmt.Printf("Checksum for %s is valid.\n", filePath)
		} else {
			return nil, fmt.Errorf("banner file  %s is corrupted\n", filePath)
		}
	}

	// Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Map to store the banner strings
	bannerMap := make(map[int][]string)
	// Starting ASCII code for mapping
	key := 32
	lineCount := 0
	chunk := []string{}

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			chunk = append(chunk, line)
			lineCount++
		}

		// Each banner is 8 lines tall
		if lineCount == 8 {
			bannerMap[key] = chunk
			key++
			chunk = []string{}
			lineCount = 0
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return bannerMap, nil
}
