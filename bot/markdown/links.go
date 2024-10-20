package markdown

import (
	"net/url"
	"regexp"
)

// ConvertLinksToHyperlinks converts URLs in the input string to markdown-style hyperlinks.
// Example:
//
//	https://foo.bar/x/y -> [foo.bar](<https://foo.bar/x/y>)
func (m *Markdown) ConvertLinksToHyperlinks(input string) string {
	// Regular expression to match URLs
	urlRegex := `https?://[^\s]+`
	re := regexp.MustCompile(urlRegex)

	// Replace all URLs with markdown-style hyperlinks
	return re.ReplaceAllStringFunc(input, func(urlStr string) string {
		// Parse the URL
		parsedUrl, err := url.Parse(urlStr)
		if err != nil {
			// If parsing fails, just return the original URL (fallback behavior)
			return urlStr
		}

		// Extract the host (domain) from the parsed URL
		domain := parsedUrl.Host

		// Convert to markdown-style hyperlink
		return "[" + domain + "](" + urlStr + ")"
	})
}
