package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

	// "path/filepath"
	"strings"

	"desrosiers.org/pse/crawler"

	"github.com/blevesearch/bleve"
	"github.com/gomutex/godocx"
)

type CrawledDocument struct {
	ID string
	Content string
}

func isBinary(path string) (bool, error) {
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

func getTextContent(filePath string) string {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading the file: '%s' \n", filePath)
		return ""
	}
	return string(bytes)
}

func getTextFromWordDoc(filePath string) string {
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

func main() {

	fs_crawler := crawler.NewFSCrawler("./fixtures/filesystem")
	err := fs_crawler.Crawl()
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}

	// open a new index
	index, err := bleve.Open("example.bleve")
	if err == bleve.ErrorIndexPathDoesNotExist {
			mapping := bleve.NewIndexMapping()
			index, err = bleve.New("example.bleve", mapping)
			if err != nil {
					panic(err)
			}
	}

	// index, err := bleve.New("example.bleve", mapping)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, path := range fs_crawler.Files {
		binary, _ := isBinary(path)
			
		content := ""
		if binary {
			fmt.Println(path)
			// TODO: more precise mime/ext matching for word documents:
			content = getTextFromWordDoc(path)
		} else {
			content = getTextContent(path)
		}
		datum := &CrawledDocument{
			ID: path,
			Content: content,
		}

		index.Index(datum.ID, datum.Content)
	}

	count, _ := index.DocCount()
	fmt.Println(count)

	// search for some text
	query := bleve.NewMatchQuery("Check")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}

