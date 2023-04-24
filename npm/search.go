package npm

import (
	"encoding/json"
	"errors"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"strings"
	"time"
)

// SearchOptions represents options for searching NPM packages.
type SearchOptions struct {
	Text        string  `url:"text"`        // Text does full-text search
	Size        int     `url:"size"`        // Size sets how many results should be returned (default 20, max 250)
	From        *int    `url:"from"`        // From offset to return results from
	Quality     float64 `url:"quality"`     // Quality how much of an effect should quality have
	Popularity  float64 `url:"popularity"`  // Popularity how much of an effect should popularity have
	Maintenance float64 `url:"maintenance"` // Maintenance how much of an effect should maintenance have
}

type SearchResponse struct {
	Total   int `json:"total"` // Total is the total number of packages matching the search
	Objects []struct {
		Package struct {
			Name        string    `json:"name"`        // Name is the name of the NPM package.
			Scope       string    `json:"scope"`       // Scope is the scope of the NPM package.
			Version     string    `json:"version"`     // Version is the version of the NPM package.
			Description string    `json:"description"` // Description is the description of the NPM package.
			Date        time.Time `json:"date"`        // Date is the publishing date of the NPM package.
			Links       struct {
				NPM        string `json:"npm"`        // NPM is the NPM link of the NPM package.
				Homepage   string `json:"homepage"`   // Homepage is the homepage link of the NPM package.
				Repository string `json:"repository"` // Repository is the repository link of the NPM package.
				Bugs       string `json:"bugs"`       // Bugs is the link to the bugs' section of the NPM package.
			} `json:"links"`
		} `json:"package"`
	} `json:"objects"`
}

// Search searches the NPM registry for packages matching the given options.
func Search(options *SearchOptions) (*SearchResponse, error) {
	// Validation
	options.Text = strings.TrimSpace(options.Text)
	if options.Text == "" {
		return nil, errors.New("package name is required")
	}
	if len(options.Text) > 214 {
		return nil, errors.New("package name length too long")
	}

	if options.Size <= 0 {
		options.Size = 20
	}

	// Generate the querystring
	v, err := query.Values(options)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	url := "https://registry.npmjs.org/-/v1/search?" + v.Encode()
	resp, err := http.Get(url) // TODO: replace with client with timeout context
	if err != nil {
		return nil, err // TODO: wrap error
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, err // TODO: wrap error
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	var data SearchResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	return &data, nil
}
