package algolia

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// getClient returns an HTTP client configured for making requests to Algolia.
func getClient() *http.Client {
	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
}

// getBaseURL returns the base URL for the Algolia API based on the given credential.
func getBaseURL(credential credential) string {
	return fmt.Sprintf("https://%v.algolia.net/1/indexes/%v/", credential.appId, credential.index)
}

// setHeaders sets the required headers for making requests to Algolia.
func setHeaders(req *http.Request, credential credential) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Algolia-API-Key", credential.apiKey)
	req.Header.Set("X-Algolia-Application-Id", credential.appId)
}
