package search

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch_SearchResultLength(t *testing.T) {
	search := NewSearch()
	keyword := "nestjs"
	results, err := search.Search(keyword)

	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Greater(t, len(results), 1, "Search result should have at least one item")
}
