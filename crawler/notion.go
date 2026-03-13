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

const CACHE_FILE_MASK = "cache/%s.json"

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

func GetBlockMarkdown(pageID string) ([]byte, error) {
	cacheKey := pageID + "-markdown"
	data, err := getCachedBlock(cacheKey)
	if err == nil {
		fmt.Println("prend la cache de markdown")
		return data, nil
	}

	apiKey := os.Getenv("NOTION_SECRET_KEY")
	url := "https://api.notion.com/v1/pages/"+pageID+"/markdown"
	req, _ := http.NewRequest("GET", url, nil)

	if apiKey == "" {
		return []byte{}, errors.New("API Key not provided.")
	}
	req.Header.Add("Notion-Version", "2026-03-11")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	cacheBlock(cacheKey, body)
	return body, nil
}

func GetMarkdown(pageID string) (string, error) {
	body, err := GetBlockMarkdown(pageID)
	if err != nil {
		panic(err)
	}

	var jsn map[string]any
	errUnmarshal := json.Unmarshal(body, &jsn)

	if errUnmarshal != nil {
		fmt.Println(errUnmarshal)
	}

	return jsn["markdown"].(string), nil
}

func getCachedBlock(blockID string) ([]byte, error) {
	data, err := os.ReadFile(fmt.Sprintf(CACHE_FILE_MASK, blockID))
	if err != nil {
		return []byte{}, errors.New("Error while opening the file. Maybe it does not exist.")
	}
	fmt.Printf("getCachedBlock %s \n", blockID)
	return data, nil
}

func cacheBlock(blockID string, data []byte) error {
	err := os.WriteFile(fmt.Sprintf(CACHE_FILE_MASK, blockID), data, 0666)
	return err
}

func GetBlockJSON(pageID string) ([]byte, error) {
	// Maybe it's cached?
	cachedBytes, err := getCachedBlock(pageID)
	if err == nil && len(cachedBytes) > 0 {
		return cachedBytes, nil
	}

	// Not cached; get it:
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", pageID)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Notion-Version", "2026-03-11")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NOTION_SECRET_KEY")))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	cacheBlock(pageID, body)

	return body, nil
}

func GetChildPageIds(parentPageID string) ([]string, error) {
	var pageIds []string
	apiKey := os.Getenv("NOTION_SECRET_KEY")
	if apiKey == "" {
		return pageIds, errors.New("API Key not provided.")
	}

	body, _ := GetBlockJSON(parentPageID)
	var resp map[string]any
	errUnmarshal := json.Unmarshal(body, &resp)
	if errUnmarshal != nil {
		fmt.Println(errUnmarshal)
	}

	results, ok := resp["results"].([]any)
	if !ok {
		return pageIds, errors.New("No result in JSON response")
	}

	// get page ids
	for _, r := range results {
		result := r.(map[string]any)
		if result["type"] == "child_page" {
			pageIds = append(pageIds, result["id"].(string))
		}
	}

	return pageIds, nil
}
