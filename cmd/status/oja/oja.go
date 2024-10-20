package oja

import (
	"encoding/json"
	"github.com/nestjs-discord/utility-bot/bot/status"
	"github.com/nestjs-discord/utility-bot/cmd/status/tidy"
	"io"
	"log"
	"net/http"
	"strings"
)

type ojaResponse []struct {
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func Fetch() []string {
	res, err := http.Get("https://raw.githubusercontent.com/15Dkatz/official_joke_api/refs/heads/master/jokes/index.json")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ojaRes ojaResponse
	err = json.Unmarshal(body, &ojaRes)
	if err != nil {
		log.Fatal(err)
	}

	uniqueTexts := make(map[string]bool)

	for _, v := range ojaRes {
		if v.Type == "" || v.Setup == "" || v.Punchline == "" {
			continue
		}

		text := strings.Join([]string{v.Setup, v.Punchline}, " ")
		text = tidy.Text(text)
		if len(text) > status.MaxCustomStatusLength {
			continue
		}

		uniqueTexts[text] = true
	}
	texts := make([]string, 0, len(uniqueTexts))
	for k := range uniqueTexts {
		texts = append(texts, k)
	}

	return texts

}
