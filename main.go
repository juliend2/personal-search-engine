package main

import (
	"fmt"
	"os"
	"runtime"

	"desrosiers.org/pse/crawler"
	"desrosiers.org/pse/parser"

	"github.com/blevesearch/bleve"
	"github.com/joho/godotenv"
)

const INDEX_FILE = "pse.bleve"

type CrawledDocument struct {
	ID string
	Content string
}

func FileSystemCrawl(path string) {
	fs_crawler := crawler.NewFSCrawler(path)
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

	// Loop into the crawled files
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
	}

	// flush remaining
	if batchCount > 0 {
			index.Batch(batch)
	}

	count, _ := index.DocCount()
	fmt.Println(count)

}

func main() {
	godotenv.Load()
	numberOfSupportedArguments := 2
	if len(os.Args) < numberOfSupportedArguments {
		usage := `
Usage: %s folder query

folder is the path from that we'll search into.
`
		panic(fmt.Sprintf(usage, os.Args[0]))
	}

	sourcePath := os.Args[1]

	if sourcePath == "notion" {
		page, err := crawler.GetNotionPage()

		if err != nil {
			panic(err)
		}
		fmt.Printf("%v \n", page.GetTitle())
		// GO: 2d5379235fe6804983a4e8b552ea211c
		markdown, err := crawler.GetMarkdown("2ab379235fe68009b4e9e3d00579ba1c")
		fmt.Println(markdown)

		// subPageIDs, err := crawler.GetChildPageIds("2ab379235fe68009b4e9e3d00579ba1c")
		// fmt.Printf("%v\n", subPageIDs)
		var initialPageID string = "2ab379235fe68009b4e9e3d00579ba1c"
		var pageIdAccumulator []string
		crawler.NotionPageSearch(&pageIdAccumulator, &initialPageID)
		fmt.Printf("Final list of pages: %v \n", pageIdAccumulator)
	} else {
		FileSystemCrawl(sourcePath)
	}
}

