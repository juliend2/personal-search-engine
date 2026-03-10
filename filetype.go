package main

import (
	"os"
	// "fmt"
	"strings"
	"regexp"
	"net/http"
)

func IsBinary(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)

	contentType := http.DetectContentType(buf[:n])
	return !strings.HasPrefix(contentType, "text/"), nil
}

func IsWordDoc(path string) bool {
	match, _ := regexp.MatchString("\\.docx$", path)
	return match
}

