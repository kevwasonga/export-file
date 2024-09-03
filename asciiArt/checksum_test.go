package asciiArt_test

import (
	"path/filepath"
	"testing"
	"ascii-art-web/asciiArt"
)

func TestVerifyBanners(t *testing.T) {
	// Get the absolute path to the root of the project
	baseDir, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("Failed to determine the base directory: %v", err)
	}

	// Construct the absolute path to the banners directory
	bannersDir := filepath.Join(baseDir, "banners")

	expectedResults := map[string]bool{
		filepath.Join(bannersDir, "standard.txt"):   true,
		filepath.Join(bannersDir, "shadow.txt"):     true,
		filepath.Join(bannersDir, "thinkertoy.txt"): true,
	}

	results := asciiArt.VerifyBanners()

	for filePath, expected := range expectedResults {
		if got, ok := results[filePath]; !ok {
			t.Errorf("File %s was not checked", filePath)
		} else if got != expected {
			t.Errorf("File %s integrity check failed. Expected %v, got %v", filePath, expected, got)
		}
	}
}
