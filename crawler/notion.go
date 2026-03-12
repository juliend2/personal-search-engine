package crawler

import (
	"errors"
	"encoding/json"
	"net/http"
	"io"
	"os"
	"fmt"
)

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


func GetNotionPage() (Page, error) {

	apiKey := os.Getenv("NOTION_SECRET_KEY")

	url := "https://api.notion.com/v1/pages/2d5379235fe6804983a4e8b552ea211c"

	req, _ := http.NewRequest("GET", url, nil)

	if apiKey == "" {
		return Page{}, errors.New("API Key not provided.")
	}
	fmt.Printf("API KEY: %s \n", apiKey)
	req.Header.Add("Notion-Version", "2025-09-03")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var block Page
	errUnmarshal := json.Unmarshal(body, &block)

	if errUnmarshal != nil {
		fmt.Println(errUnmarshal)
	}

	return block, nil
}

func GetMarkdown(pageId string) (string, error) {
	apiKey := os.Getenv("NOTION_SECRET_KEY")

	url := "https://api.notion.com/v1/pages/"+pageId+"/markdown"

	req, _ := http.NewRequest("GET", url, nil)

	if apiKey == "" {
		return "", errors.New("API Key not provided.")
	}
	fmt.Printf("API KEY: %s \n", apiKey)
	req.Header.Add("Notion-Version", "2026-03-11")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var jsn map[string]any
	errUnmarshal := json.Unmarshal(body, &jsn)

	if errUnmarshal != nil {
		fmt.Println(errUnmarshal)
	}

	return jsn["markdown"].(string), nil
}
