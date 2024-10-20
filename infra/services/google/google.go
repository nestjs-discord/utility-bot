package google

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Google struct {
	client *http.Client
}

func NewGoogle() *Google {
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,               // Maximum idle connections to keep open.
			IdleConnTimeout:     20 * time.Second, // Maximum time an idle (keep-alive) connection will be kept.
			TLSHandshakeTimeout: 2 * time.Second,  // Timeout for TLS handshake.
		},
	}
	return &Google{
		client: client,
	}
}

type TopLevel struct {
	CompleteSuggestions []CompleteSuggestion `xml:"CompleteSuggestion"`
}

type CompleteSuggestion struct {
	Suggestion Suggestion `xml:"suggestion"`
}

type Suggestion struct {
	Data string `xml:"data,attr"`
}

func (g *Google) Search(keyword string) ([]string, error) {
	keyword = strings.TrimSpace(keyword) // normalize
	//apiUrl := "https://google.com/complete/search?client=gws-wiz&xssi=t&hl=en-US&authuser=0&dpr=1&q=" +
	//	url.QueryEscape(keyword)
	apiUrl := "https://suggestqueries.google.com/complete/search?output=toolbar&hl=en&q=" + url.QueryEscape(keyword)

	resp, err := g.client.Get(apiUrl)
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

	var topLevel TopLevel

	// Unmarshal the XML into the TopLevel struct
	err = xml.Unmarshal(body, &topLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	var res []string

	// Print the parsed suggestions
	for _, completeSuggestion := range topLevel.CompleteSuggestions {
		res = append(res, completeSuggestion.Suggestion.Data)
	}

	return res, nil

	//jsonStr := strings.Replace(string(body), ")]}'", "", 1)
	//
	//r := gjson.Get(strings.Replace(string(body), ")]}'", "", 1), "0.#.0")
	//
	//var res []string
	//for i, v := range r.Array() {
	//	if i == 8 { // limit the result to eight elements
	//		break
	//	}
	//
	//	text := strings.ReplaceAll(v.String(), "<b>", " ")
	//	text = strings.ReplaceAll(text, "</b>", "")
	//	text = strings.ReplaceAll(text, "  ", " ")
	//	text = strings.TrimSpace(text)
	//
	//	if len(text) == 0 { // to avoid adding empty text to the response array
	//		continue
	//	}
	//
	//	res = append(res, text)
	//}
}
