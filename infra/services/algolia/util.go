package algolia

import (
	"strings"
)

// GetFormattedHierarchy formats a hierarchy from the provided Hit
// struct into a human-readable string representation.
func GetFormattedHierarchy(hit Hit) string {
	var builder strings.Builder

	builder.WriteString(hit.Hierarchy.Lvl0)
	if hit.Hierarchy.Lvl1 != "" {
		builder.WriteString(" - ")
		builder.WriteString(hit.Hierarchy.Lvl1)
	}

	if hit.Hierarchy.Lvl2 != "" {
		builder.WriteString(": ")
		builder.WriteString(hit.Hierarchy.Lvl2)
	}

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
