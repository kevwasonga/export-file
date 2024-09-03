package asciiArt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Map to store the known checksums for each banner file
var knownChecksums = map[string]string{
	"banners/standard.txt":   "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf",
	"banners/shadow.txt":     "26b94d0b134b77e9fd23e0360bfd81740f80fb7f6541d1d8c5d85e73ee550f73",
	"banners/thinkertoy.txt": "64285e4960d199f4819323c4dc6319ba34f1f0dd9da14d07111345f5d76c3fa3",
}

func VerifyBanners() map[string]bool {
	results := make(map[string]bool)
	for filePath, expectedChecksum := range knownChecksums {
		isValid, err := VerifyChecksum(filePath, expectedChecksum)
		if err != nil {
			fmt.Printf("Error verifying %s: %v\n", filePath, err)
			results[filePath] = false
			continue
		}

		results[filePath] = isValid
	}
	return results
}

func VerifyChecksum(filePath, expectedChecksum string) (bool, error) {
	computedChecksum, err := ComputeSHA256(filePath)
	if err != nil {
		return false, err
	}

	return computedChecksum == expectedChecksum, nil
}

func ComputeSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
