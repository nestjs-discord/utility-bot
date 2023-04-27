package npm

import (
	"encoding/json"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// SearchOptions represents options for searching NPM packages.
type SearchOptions struct {
	Text        string  `url:"text"`                  // Text does full-text search
	Size        int64   `url:"size,omitempty"`        // Size sets how many results should be returned (default 20, max 250)
	From        int64   `url:"from,omitempty"`        // From offset to return results from
	Quality     float64 `url:"quality,omitempty"`     // Quality how much of an effect should quality have
	Popularity  float64 `url:"popularity,omitempty"`  // Popularity how much of an effect should popularity have
	Maintenance float64 `url:"maintenance,omitempty"` // Maintenance how much of an effect should maintenance have
}

type SearchResponse struct {
	Total   int64 `json:"total"` // Total is the total number of packages matching the search
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
		return nil, errors.New("text field is required")
	}
	if len(options.Text) > 214 {
		return nil, errors.New("text length too long")
	}

	if options.Size <= 0 {
		options.Size = 20
	}

	// Generate the querystring
	v, err := query.Values(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate querystring")
	}

	url := "https://registry.npmjs.org/-/v1/search?" + v.Encode()
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to interact with NPM registry API")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(err, "NPM registry API returned unacceptable status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read NPM registry response body")
	}

	var data SearchResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to JSON parse NPM registry response")
	}

	return &data, nil
}
