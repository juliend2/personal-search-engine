package main

import (
	"fmt"
	"os"
	"net/http"
	"strings"

	"desrosiers.org/pse/crawler"

	"github.com/blevesearch/bleve"
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
		fmt.Printf("Error reading the file: '%s'", filePath)
		return ""
	}
	return string(bytes)
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
		// fmt.Printf("%v", binary)

		if binary {
			fmt.Println("binary format not handled yet...")
		}
		// fmt.Println(path)
		// fmt.Println(getTextContent(path))
		datum := &CrawledDocument{
			ID: path,
			Content: getTextContent(path),
		}

		index.Index(datum.ID, datum.Content)
	}

	count, _ := index.DocCount()
	fmt.Println(count)

	// search for some text
	query := bleve.NewMatchQuery("français")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}

