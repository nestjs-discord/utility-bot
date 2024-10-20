package cf

import (
	"encoding/json"
	"fmt"
	"github.com/nestjs-discord/utility-bot/bot/status"
	"github.com/nestjs-discord/utility-bot/cmd/status/tidy"
	"io"
	"log"
	"net/http"
)

type response struct {
	Data []struct {
		Fact string `json:"fact"`
	} `json:"data"`
}

func Fetch() []string {
	u := fmt.Sprintf("https://catfact.ninja/facts?limit=1000&max_length=%d", status.MaxCustomStatusLength)
	res, err := http.Get(u)
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

	var bodyData response
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		log.Fatal(err)
	}

	uniqueTexts := make(map[string]bool)
	for _, v := range bodyData.Data {
		if v.Fact == "" {
			continue
		}

		text := tidy.Text(v.Fact)
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
