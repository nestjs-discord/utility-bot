package google

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch_SearchResultLength(t *testing.T) {
	search := NewGoogle()
	keyword := "nestjs"
	results, err := search.Search(keyword)

	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Greater(t, len(results), 1, "Google result should have at least one item")
}
