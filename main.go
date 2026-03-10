package main

import (
	"fmt"

	"desrosiers.org/pse/crawler"
	"desrosiers.org/pse/parser"

	"github.com/blevesearch/bleve"
)

const INDEX_FILE = "pse.bleve"

type CrawledDocument struct {
	ID string
	Content string
}

func main() {

	fs_crawler := crawler.NewFSCrawler("./fixtures/filesystem")
	err := fs_crawler.Crawl()
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}

	// open a new index
	index, err := bleve.Open(INDEX_FILE)
	if err == bleve.ErrorIndexPathDoesNotExist {
			mapping := bleve.NewIndexMapping()
			index, err = bleve.New(INDEX_FILE, mapping)
			if err != nil {
					panic(err)
			}
	}

	for _, path := range fs_crawler.Files {
		binary, _ := IsBinary(path)
			
		content := ""
		if binary {
			if IsWordDoc(path) {
				content = parser.GetTextFromWordDoc(path)
			} else {
				content = "" // Just so it doesn't create an error
			}
		} else {
			content = parser.GetTextContent(path)
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

