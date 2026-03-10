package parser

import (
	"os"
	"encoding/xml"
	"fmt"
	"strings"
	"github.com/gomutex/godocx"
)

func GetTextContent(filePath string) string {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading the file: '%s' \n", filePath)
		return ""
	}
	return string(bytes)
}

func GetTextFromWordDoc(filePath string) string {
	rootDoc, err := godocx.OpenDocument(filePath)
	if err != nil {
		fmt.Printf("Error reading the docx %s \n", err)
	}
	xmlBytes, err := xml.Marshal(rootDoc)
	decoder := xml.NewDecoder(strings.NewReader(string(xmlBytes)))
	var result strings.Builder
	for {
		tok, err := decoder.Token()
		if err != nil {
			break // io.EOF when done
		}
		if charData, ok := tok.(xml.CharData); ok {
			result.Write(charData)
		}

	}
	return result.String()
}

