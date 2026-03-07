package main

import (
	"encoding/json"
	"net/http"
	"io"
	"os"
	"fmt"
)

const API_VERSION = "2025-09-03"

type TitleItem struct {
    PlainText   string      `json:"plain_text"`
    Type        string      `json:"type"`
}

type Page struct {
	Id 					string `json:id`
	Url 				string `json:url`
	CreatedTime string `json:created_time`
	Properties struct {
        Title struct {
            ID    string      `json:"id"`
            Type  string      `json:"type"`
            Title []TitleItem `json:"title"` // Note: this is an array
        } `json:"title"`
    } `json:"properties"`
}
func (p *Page) GetTitle() string {
    items := p.Properties.Title.Title
    if len(items) == 0 {
        return ""
    }
    return items[0].PlainText
}


func GetSubBlocks(blockId string) {

	apiKey := os.Getenv("NOTION_SECRET_KEY")

	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", blockId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Notion-Version", API_VERSION)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
	
	// return page
}

func GetNotionPage() Page {

	apiKey := os.Getenv("NOTION_SECRET_KEY")

	url := "https://api.notion.com/v1/pages/312379235fe6803babe3f8cd2dd8cce7"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Notion-Version", API_VERSION)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var block Page
	err := json.Unmarshal(body, &block)

	if err != nil {
		fmt.Println(err)
	}

	return block
}
