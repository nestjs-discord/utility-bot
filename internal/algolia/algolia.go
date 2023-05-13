// Package algolia provides utilities related to the Algolia search service.
// Documentation: https://www.algolia.com/doc/api-reference/api-methods/
package algolia

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Hit represents a search result hit.
type Hit struct {
	ObjectID  string    `json:"objectID"`
	URL       string    `json:"url"`
	Anchor    string    `json:"anchor"`
	Content   string    `json:"content"`
	Hierarchy Hierarchy `json:"hierarchy"`
}

// Hierarchy represents the hierarchical structure of a search result hit.
type Hierarchy struct {
	Lvl0 string `json:"lvl0"`
	Lvl1 string `json:"lvl1"`
	Lvl2 string `json:"lvl2"`
	Lvl3 string `json:"lvl3"`
	Lvl4 string `json:"lvl4"`
	Lvl5 string `json:"lvl5"`
	Lvl6 string `json:"lvl6"`
}

type queryResponse struct {
	Hits []Hit `json:"hits"`
}

// Search performs a search query on the Algolia search service for the specified app and query.
// It returns a list of search result hits and any error encountered.
func Search(app App, query string) (hits []Hit, error error) {
	query = strings.TrimSpace(query)
	if len(query) > 500 {
		return nil, errors.New("There's a hard limit of 500 characters per query")
	}

	credential := credentials[app]
	baseURL := getBaseURL(credential)

	jsonPayload, err := generateQueryJSONPayload(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate JSON payload")
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"query", bytes.NewBuffer(jsonPayload))

	setHeaders(req, credential)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "expected status code OK but received %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read all response body failed")
	}

	var data queryResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal failed")
	}

	return data.Hits, nil
}

func generateQueryJSONPayload(query string) ([]byte, error) {
	payload := map[string]string{
		"query":          query,
		"hitsPerPage":    "24", // Because Discord commands option choices is limited to 25
		"queryLanguages": "en",
	}

	return json.Marshal(payload)
}

// GetObject retrieves a specific object from the Algolia search service based on the objectID.
// It returns the hit representing the object and any error encountered.
func GetObject(app App, objectID string) (*Hit, error) {
	credential := credentials[app]

	baseURL := getBaseURL(credential)

	req, err := http.NewRequest(http.MethodGet, baseURL+url.QueryEscape(objectID), nil)

	setHeaders(req, credential)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http request failed")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "expected status code OK but received %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read all response body failed")
	}

	var hit Hit
	err = json.Unmarshal(body, &hit)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal failed")
	}

	return &hit, nil
}
