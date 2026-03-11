all: main search

crawl: main
	./main notion

search:
	go run cmd/search/cli.go crainte

main: clean
	go build -o main main.go filetype.go notion_client.go

clean:
	rm -f ./main
