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
	u := "https://raw.githubusercontent.com/15Dkatz/official_joke_api/refs/heads/master/jokes/index.json"
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
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

		text := strings.Join([]string{
			tidy.Text(v.Setup),
			tidy.Text(v.Punchline),
		}, " ")
		text = tidy.Text(text)
		uniqueTexts[text] = true
	}
	texts := make([]string, 0, len(uniqueTexts))
	for k := range uniqueTexts {
		if len(k) > status.MaxCustomStatusLength {
			continue
		}

		texts = append(texts, k)
	}

	return texts
}
