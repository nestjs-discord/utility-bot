package google

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearch_SearchResultLength(t *testing.T) {
	search := NewGoogle()
	keyword := "nestjs"

	results, err := search.Search(keyword)
	require.Nil(t, err)
	require.NotNil(t, results)
	require.Greater(t, len(results), 1, "Google result should have at least one item")
}
