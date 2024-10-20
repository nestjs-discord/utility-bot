package algolia

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	cacheSize := 100

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ClientSessionCache: tls.NewLRUClientSessionCache(cacheSize),
	}

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     tlsConfig,
	}

	client = &http.Client{
		Timeout:   8 * time.Second,
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
	req.Header.Set("X-Algolia-Api-Key", credential.apiKey)
	req.Header.Set("X-Algolia-Application-Id", credential.appId)
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
}
