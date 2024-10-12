package forms

import (
	"net/url"
	"regexp"
)

const FormButtonIdPrefix = "form-"
const FormModalIdPrefix = "modal-"

type UserInput struct {
	InputId string
	Value   string
}

// ConvertLinksToHyperlinks converts URLs in the input string to markdown-style hyperlinks.
// Example:
//
//	https://foo.bar/x/y -> [foo.bar](<https://example.com/x/y>)
//func ConvertLinksToHyperlinks(input string) string {
//	// Regular expression to match URLs
//	urlRegex := `https?://[^\s]+`
//	re := regexp.MustCompile(urlRegex)
//
//	// Replace all URLs with markdown-style hyperlinks
//	return re.ReplaceAllStringFunc(input, func(url string) string {
//		// Extract the domain name from the URL for display
//		domain := url
//		if idx := strings.Index(url, "//"); idx != -1 {
//			domain = url[idx+2:]
//		}
//		if idx := strings.Index(domain, "/"); idx != -1 {
//			domain = domain[:idx]
//		}
//
//		// Convert to markdown-style hyperlink
//		return "[" + domain + "](" + url + ")"
//	})
//}

// ConvertLinksToHyperlinks converts URLs in the input string to markdown-style hyperlinks.
// Example:
//	https://foo.bar/x/y -> [foo.bar](<https://example.com/x/y>)
//func ConvertLinksToHyperlinks(input string) string {
//	// Regular expression to match URLs
//	urlRegex := `https?://[^\s]+`
//	re := regexp.MustCompile(urlRegex)
//
//	// Replace all URLs with markdown-style hyperlinks
//	return re.ReplaceAllStringFunc(input, func(url string) string {
//		// Extract the domain name from the URL for display
//		domain := url
//		if idx := strings.Index(url, "//"); idx != -1 {
//			domain = url[idx+2:]
//		}
//		if idx := strings.Index(domain, "/"); idx != -1 {
//			domain = domain[:idx]
//		}
//
//		// Convert to markdown-style hyperlink
//		return "[" + domain + "](" + url + ")"
//	})
//}

// ConvertLinksToHyperlinks converts URLs in the input string to markdown-style hyperlinks.
// Example:
//
//	https://foo.bar/x/y -> [foo.bar](<https://foo.bar/x/y>)
func ConvertLinksToHyperlinks(input string) string {
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
