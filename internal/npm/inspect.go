package npm

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type InspectOptions struct {
	Name    string
	Version string
}

// InspectResponse package metadata https://github.com/npm/registry/blob/master/docs/responses/package-metadata.md
type InspectResponse struct {
	// Name is the NPM package name
	Name string `json:"name"`
	// Version for this version
	Version string `json:"version"`
	// Description a short description of the package
	Description string `json:"description"`
	// License the SPDX identifier of the package's license
	License string `json:"license"`
	// Homepage url
	Homepage string `json:"homepage"`

	// Repository as given in package.json, for the given version
	Repository struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"repository"`

	Engines          map[string]string `json:"engines"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	// Types            string            `json:"types"`
	GitHead string `json:"gitHead"`
	Bugs    struct {
		Url string `json:"url"`
	} `json:"bugs"`
	Dist struct {
		Integrity string `json:"integrity"`
		// Shasum       string `json:"shasum"`
		// Tarball is the url of the tarball containing the payload for this package
		Tarball string `json:"tarball"`
		// FileCount the number of files in the tarball, folder excluded
		FileCount int64 `json:"fileCount"`
		// UnpackedSize is the total byte of the unpacked files in the tarball
		UnpackedSize uint64 `json:"unpackedSize"`
	} `json:"dist"`

	Keywords []string `json:"keywords"`

	// Donation
	Funding struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"funding"`
}

func Inspect(options *InspectOptions) (*InspectResponse, error) {
	// Validation
	options.Name = strings.TrimSpace(options.Name)
	if options.Name == "" {
		return nil, errors.New("name field is required")
	}
	if len(options.Name) > 214 {
		return nil, errors.New("name length too long")
	}

	if options.Version == "" {
		options.Version = "latest"
	}

	url := fmt.Sprintf("https://registry.npmjs.org/%v/%v", options.Name, options.Version)
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
		return nil, fmt.Errorf("NPM registry API returned unacceptable status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read NPM registry response body")
	}

	var data InspectResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to JSON parse NPM registry response")
	}

	return &data, nil
}
