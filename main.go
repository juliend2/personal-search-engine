package main

import (
	"fmt"
	"os"
	"net/http"
	"strings"

	"desrosiers.org/pse/crawler"

	// "github.com/blevesearch/bleve"
)

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


func main() {

	fs_crawler := crawler.NewFSCrawler("./fixtures/filesystem")
	err := fs_crawler.Crawl()
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}

	for _, path := range fs_crawler.Files {
		b, _ := isBinary(path)
		fmt.Printf("%v", b)
		fmt.Println(path)
	}

	// open a new index
	// index, err := bleve.Open("example.bleve")
	// if err == bleve.ErrorIndexPathDoesNotExist {
	// 		mapping := bleve.NewIndexMapping()
	// 		index, err = bleve.New("example.bleve", mapping)
	// 		if err != nil {
	// 				panic(err)
	// 		}
	// }
	//
	// // index, err := bleve.New("example.bleve", mapping)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// data := struct {
	// 	Name string
	// }{
	// 	Name: "text",
	// }
	//
	// // index some data
	// index.Index("id", data)
	//
	// data1 := struct {
	// 	Name string
	// }{
	// 	Name: "julien",
	// }
	//
	// // index some data
	// index.Index("dowit", data1)
	//
	// count, _ := index.DocCount()
	// fmt.Println(count)
	//
	//
	//
	// // search for some text
	// query := bleve.NewMatchQuery("julien")
	// search := bleve.NewSearchRequest(query)
	// searchResults, err := index.Search(search)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(searchResults)
}

