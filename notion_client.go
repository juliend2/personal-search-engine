package main

import (
	"net/http"
	"io"
	"os"
	"fmt"
	"github.com/joho/godotenv"
)

func GetPage() {
	
	godotenv.Load() // loads .env file automatically

	apiKey := os.Getenv("NOTION_SECRET_KEY")


	url := "https://api.notion.com/v1/pages/312379235fe6803babe3f8cd2dd8cce7"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Notion-Version", "2025-09-03")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}
