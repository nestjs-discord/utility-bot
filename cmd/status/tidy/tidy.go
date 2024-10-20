package tidy

import "strings"

func Text(text string) string {
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.ReplaceAll(text, `"`, `\"`)
	text = strings.TrimSpace(text)
	return text
}
