package pqa

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nestjs-discord/utility-bot/bot/status"
	"github.com/nestjs-discord/utility-bot/cmd/status/tidy"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const scrapedDataFilename = "./pqa/data.json"

type scrapedDataType map[string]bool

type response []struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

func Scrape() error {
	scrapedData := readAndParseScrapedData()

	u := "https://programming-quotesapi.vercel.app/api/bulk"
	res, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("http get failed: %s", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var resBody response
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return fmt.Errorf("json decode failed: %s", err)
	}
	if len(resBody) == 0 {
		return errors.New("no data found")
	}

	for _, item := range resBody {
		key := item.Quote

		if strings.ToLower(item.Author) != "anonymous" {
			key += " - " + item.Author
		}

		if len(key) > status.MaxCustomStatusLength {
			continue
		}

		scrapedData[key] = true
	}

	encodedScrapedData, err := json.MarshalIndent(scrapedData, "", " ")
	if err != nil {
		return fmt.Errorf("json marshal failed: %s", err)
	}

	err = os.WriteFile(scrapedDataFilename, encodedScrapedData, 0644)
	if err != nil {
		return fmt.Errorf("write file failed: %s", err)
	}

	return nil
}

func readAndParseScrapedData() scrapedDataType {
	fileContent, err := os.ReadFile(scrapedDataFilename)
	if err != nil {
		log.Fatal(err)
	}
	var scrapedData scrapedDataType
	err = json.Unmarshal(fileContent, &scrapedData)
	if err != nil {
		log.Fatal(err)
	}

	return scrapedData
}

func Fetch() []string {
	scrapedData := readAndParseScrapedData()
	texts := make([]string, 0, len(scrapedData))
	for k := range scrapedData {
		k = tidy.Text(k)
		if len(k) > status.MaxCustomStatusLength {
			continue
		}

		texts = append(texts, k)
	}
	return texts
}
