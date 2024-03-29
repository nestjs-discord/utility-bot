package algolia_test

import (
	algolia2 "github.com/nestjs-discord/utility-bot/pkg/algolia"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFormattedHierarchy_CategoryAndSubcategoryProvided(t *testing.T) {
	hit := algolia2.Hit{
		Hierarchy: algolia2.Hierarchy{
			Lvl0: "Books",
			Lvl1: "Fiction",
			Lvl2: "Mystery",
			Lvl3: "Thriller",
		},
	}
	expectedResult := "Books - Fiction: Mystery - Thriller"

	result := algolia2.GetFormattedHierarchy(hit)
	assert.Equal(t, expectedResult, result)
}

func TestGetFormattedHierarchy_OnlyCategoryProvided(t *testing.T) {
	hit := algolia2.Hit{
		Hierarchy: algolia2.Hierarchy{
			Lvl0: "Books",
			Lvl1: "Fiction",
			Lvl2: "",
			Lvl3: "Drama",
		},
	}
	expectedResult := "Books - Fiction - Drama"

	result := algolia2.GetFormattedHierarchy(hit)
	assert.Equal(t, expectedResult, result)
}

func TestGetFormattedHierarchy_OnlyLvl0Provided(t *testing.T) {
	hit := algolia2.Hit{
		Hierarchy: algolia2.Hierarchy{
			Lvl0: "Books",
			Lvl1: "",
			Lvl2: "",
			Lvl3: "",
		},
	}
	expectedResult := "Books"

	result := algolia2.GetFormattedHierarchy(hit)
	assert.Equal(t, expectedResult, result)
}

func TestTruncate_WhenTextIsShorterOrEqual_ReturnsOriginalText(t *testing.T) {
	// Text length <= length
	text := "Hello, world!"
	truncated := algolia2.Truncate(text, 15)
	assert.Equal(t, text, truncated)
}

func TestTruncate_WhenTextExceedsLength_TruncatesAtWordBoundary(t *testing.T) {
	// Text length > length, truncate at word boundary
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
	truncated := algolia2.Truncate(text, 20)
	expected := "Lorem ipsum dolor..."
	assert.Equal(t, expected, truncated)
}

func TestTruncate_WhenTextDoesNotContainSpaces_TruncatesAtCharacterLimit(t *testing.T) {
	// Text length > length, truncate at character limit
	text := "ThisIsAVeryLongWordWithoutSpaces"
	truncated := algolia2.Truncate(text, 10)
	expected := "ThisIsAVer..."
	assert.Equal(t, expected, truncated)
}

func TestTruncate_WhenTextContainsLeadingOrTrailingWhitespaces_TrimsWhitespaces(t *testing.T) {
	// Text length > length, trim leading and trailing whitespaces
	text := "   Trim leading and trailing   "
	truncated := algolia2.Truncate(text, 20)
	expected := "Trim leading and..."
	assert.Equal(t, expected, truncated)
}
