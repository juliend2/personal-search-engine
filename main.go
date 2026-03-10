package main

import (
	"fmt"
	"os"
	"runtime"

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
	numberOfSupportedArguments := 3
	if len(os.Args) < numberOfSupportedArguments {
		usage := `
Usage: %s folder query

folder is the path from that we'll search into.
query is a string that you search for among those files.
`
		panic(fmt.Sprintf(usage, os.Args[0]))
	}

	sourcePath := os.Args[1]
	searchQuery := os.Args[2]

	fs_crawler := crawler.NewFSCrawler(sourcePath)
	err := fs_crawler.Crawl()
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}

	// Create or Load the index
	index, err := bleve.Open(INDEX_FILE)
	if err == bleve.ErrorIndexPathDoesNotExist {
		fmt.Println("Creating new index...")
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(INDEX_FILE, mapping)
		if err != nil {
			panic(err)
		}
	}

	const BATCH_SIZE = 100
	
	batch := index.NewBatch()
	batchCount := 0


	for _, path := range fs_crawler.Files {
		fmt.Println(path)

		if IsSkippable(path) {
			fmt.Printf("Skipping %s \n", path)
			continue
		}
		binary, _ := IsBinary(path)
			
		content := ""
		if binary {
			switch {
			case IsPDF(path):
				content, _ = parser.GetTextFromPdf(path)
			case IsWordDoc(path):
				content, _ = parser.GetTextFromWordDoc(path)
			default:
				content = ""
			}
		} else {
			content = parser.GetTextContent(path)
		}
		if content == "" {
			continue // don't bother persisting to index if this is empty
		}
		batch.Index(path, content)
    batchCount++

		if batchCount >= BATCH_SIZE {
			if err := index.Batch(batch); err != nil {
				fmt.Println("batch error:", err)
			}
			batch = index.NewBatch()  // reset
			batchCount = 0
			runtime.GC()
		}

		// datum := &CrawledDocument{
		// 	ID: path,
		// 	Content: content,
		// }

		// index.Index(datum.ID, datum.Content)
	}

	// flush remaining
	if batchCount > 0 {
			index.Batch(batch)
	}

	count, _ := index.DocCount()
	fmt.Println(count)

	// search for some text
	// keep this only for dev's feedback loop.
	query := bleve.NewMatchQuery(searchQuery)
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}

