package main

import (
	"fmt"

	"github.com/blevesearch/bleve"
)

func main() {
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

	data := struct {
		Name string
	}{
		Name: "text",
	}

	// index some data
	index.Index("id", data)

	data1 := struct {
		Name string
	}{
		Name: "julien",
	}

	// index some data
	index.Index("dowit", data1)

	count, _ := index.DocCount()
	fmt.Println(count)



	// search for some text
	query := bleve.NewMatchQuery("julien")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}

