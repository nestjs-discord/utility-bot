package algolia

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBaseURL(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		credential     credential
		expectedResult string
	}{
		{
			name: "Test case 1",
			credential: credential{
				appId:  "APP_ID",
				apiKey: "API_KEY",
				index:  "INDEX_NAME",
			},
			expectedResult: "https://APP_ID.algolia.net/1/indexes/INDEX_NAME/",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getBaseURL(tc.credential)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestSetHeaders(t *testing.T) {
	// Create a fake HTTP request
	req, _ := http.NewRequest("GET", "https://example.com", nil)

	// Test cases
	testCases := []struct {
		name           string
		credential     credential
		expectedResult http.Header
	}{
		{
			name: "Test case 1",
			credential: credential{
				appId:  "APP_ID",
				apiKey: "API_KEY",
				index:  "INDEX_NAME",
			},
			expectedResult: http.Header{
				"Content-Type":             []string{"application/json"},
				"X-Algolia-Api-Key":        []string{"API_KEY"},
				"X-Algolia-Application-Id": []string{"APP_ID"},
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setHeaders(req, tc.credential)
			assert.Equal(t, tc.expectedResult, req.Header)
		})
	}
}
