package search

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 2 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        10,               // Maximum idle connections to keep open.
		IdleConnTimeout:     20 * time.Second, // Maximum time an idle (keep-alive) connection will be kept.
		TLSHandshakeTimeout: 2 * time.Second,  // Timeout for TLS handshake.
	},
}

type Search struct{}

func NewSearch() *Search {
	return &Search{}
}

func (s *Search) Search(keyword string) ([]string, error) {
	keyword = strings.TrimSpace(keyword) // normalize
	apiUrl := "https://google.com/complete/search?client=gws-wiz&xssi=t&hl=en-US&authuser=0&dpr=1&q=" +
		url.QueryEscape(keyword)

	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google returned non-ok status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	r := gjson.Get(strings.Replace(string(body), ")]}'", "", 1), "0.#.0")

	var res []string
	for i, v := range r.Array() {
		if i == 8 { // limit the result to eight elements
			break
		}

		text := strings.ReplaceAll(v.String(), "<b>", " ")
		text = strings.ReplaceAll(text, "</b>", "")
		text = strings.ReplaceAll(text, "  ", " ")
		text = strings.TrimSpace(text)
		// url := "https://google.com/search?q=" + url.QueryEscape(text)

		res = append(res, text)
	}

	return res, nil
}
