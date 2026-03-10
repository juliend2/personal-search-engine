run:
	go run .

clean:
	rm ./main

main: clean
	go build -o main main.go filetype.go notion_client.go
