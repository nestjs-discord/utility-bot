package ja

import (
	"encoding/json"
	"github.com/nestjs-discord/utility-bot/bot/status"
	"github.com/nestjs-discord/utility-bot/cmd/status/tidy"
	"io"
	"log"
	"net/http"
)

type response struct {
	Jokes []struct {
		Category string `json:"category"`
		Type     string `json:"type"`
		Joke     string `json:"joke,omitempty"`
		Flags    struct {
			Nsfw      bool `json:"nsfw"`
			Religious bool `json:"religious"`
			Political bool `json:"political"`
			Racist    bool `json:"racist"`
			Sexist    bool `json:"sexist"`
			Explicit  bool `json:"explicit"`
		} `json:"flags"`
		Id       int    `json:"id"`
		Safe     bool   `json:"safe"`
		Setup    string `json:"setup,omitempty"`
		Delivery string `json:"delivery,omitempty"`
	} `json:"jokes"`
}

func Fetch() []string {
	u := "https://raw.githubusercontent.com/Sv443-Network/JokeAPI/refs/heads/main/data/jokes/regular/jokes-en.json"
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

	var data response
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	uniqueTexts := make(map[string]bool)

	for _, joke := range data.Jokes {
		if !joke.Safe {
			continue
		}
		f := joke.Flags
		if f.Nsfw || f.Religious || f.Political || f.Racist || f.Sexist || f.Explicit {
			continue
		}

		switch joke.Type {
		case "single":
			key := tidy.Text(joke.Joke)
			uniqueTexts[key] = true
		case "twopart":
			key := tidy.Text(joke.Setup + " " + joke.Delivery)
			uniqueTexts[key] = true
		}
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
