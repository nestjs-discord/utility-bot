package algolia

import (
	"fmt"
	"strings"
)

// ResolveHitToName resolves the name of a hit based on its hierarchy.
// It takes a Hit object and returns the resolved name as a string.
// If the subcategory is not specified or is the same as the category, it defaults to "Introduction".
// If the Lvl3 of the hierarchy is not empty, it appends it to the result with a dash.
func ResolveHitToName(hit Hit) string {
	category := hit.Hierarchy.Lvl1
	if category == "" {
		category = hit.Hierarchy.Lvl0
	}

	subcategory := hit.Hierarchy.Lvl2
	if subcategory == "" || category == subcategory {
		subcategory = "Introduction"
	}

	result := fmt.Sprintf("%s: %s", category, subcategory)
	if hit.Hierarchy.Lvl3 != "" {
		result += fmt.Sprintf(" - %s", hit.Hierarchy.Lvl3)
	}

	return result
}

// Truncate truncates a given text to a specified length while preserving word boundaries.
// It takes the text to truncate, the desired length, and the split character as input.
// If the length of the text is less than or equal to the specified length, it returns the original text.
// It splits the text into words using the split character and gradually adds words until the length is reached.
// If truncation occurs, it appends "..." to the truncated text.
func Truncate(text string, length int, splitChar string) string {
	if len(text) <= length {
		return text
	}

	words := strings.Split(text, splitChar)
	var res []string

	for _, word := range words {
		full := strings.Join(res, splitChar)
		if len(full)+len(word)+1 <= length-3 {
			res = append(res, word)
		}
	}

	resText := strings.Join(res, splitChar)
	if len(resText) == len(text) {
		return resText
	}

	return fmt.Sprintf("%s...", strings.TrimSpace(resText))
}
