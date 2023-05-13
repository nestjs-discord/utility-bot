package algolia

import (
	"strings"
)

// GetFormattedHierarchy formats a hierarchy from the provided Hit
// struct into a human-readable string representation.
//
// The function follows the following logic:
//   - If the "lvl1" field of the Hit's Hierarchy is not empty, it is used as the category.
//     Otherwise, if the "lvl0" field is not empty, it is used as the category.
//   - If the "lvl2" field of the Hit's Hierarchy is not empty, it is used as the subcategory.
//     Otherwise, if the "lvl1" field is not empty, it is used as the subcategory.
//   - If the category is the same as the subcategory or the subcategory is empty,
//     the subcategory is set to "Introduction".
//   - If the "lvl3" field of the Hit's Hierarchy is not empty,
//     it is appended to the result string separated by a dash.
//
// The resulting formatted hierarchy string is returned.
func GetFormattedHierarchy(hit Hit) string {
	var builder strings.Builder

	category := hit.Hierarchy.Lvl1
	if category == "" {
		category = hit.Hierarchy.Lvl0
	}

	subcategory := hit.Hierarchy.Lvl2
	if subcategory == "" {
		subcategory = hit.Hierarchy.Lvl1
	}

	if category == subcategory || subcategory == "" {
		subcategory = "Introduction"
	}

	builder.WriteString(category)
	builder.WriteString(": ")
	builder.WriteString(subcategory)

	if hit.Hierarchy.Lvl3 != "" {
		builder.WriteString(" - ")
		builder.WriteString(hit.Hierarchy.Lvl3)
	}

	return builder.String()
}

// Truncate truncates a given text to a specified length by cutting at word boundaries.
// If the length of the text is less than or equal to the specified length, it returns the original text.
// If truncation occurs, it appends "..." to the truncated text.
func Truncate(text string, length int) string {
	if len(text) <= length {
		return text
	}

	// Trim leading and trailing whitespaces
	text = strings.TrimSpace(text)

	if len(text) > length {
		// Find the last complete word within the character limit
		lastSpaceIndex := strings.LastIndex(text[:length], " ")
		if lastSpaceIndex > 0 {
			text = text[:lastSpaceIndex]
		} else {
			text = text[:length]
		}
	}

	return text + "..."
}
