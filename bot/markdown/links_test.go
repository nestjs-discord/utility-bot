package markdown

import (
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertLinksToHyperlinks(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Visit https://foo.bar/x/y for more info.",
			expected: "Visit [foo.bar](https://foo.bar/x/y) for more info.",
		},
		{
			input:    "Check out http://example.com/a/b.",
			expected: "Check out [example.com](http://example.com/a/b.)",
		},
		{
			input:    "Multiple links https://abc.com/path1 and http://xyz.org/path2.",
			expected: "Multiple links [abc.com](https://abc.com/path1) and [xyz.org](http://xyz.org/path2.)",
		},
		{
			input:    "This has no links at all.",
			expected: "This has no links at all.",
		},
		{
			input:    "Edge case https://example.com.",
			expected: "Edge case [example.com.](https://example.com.)",
		},
	}

	opts := Options{
		Commands: yaml.Commands{},
	}

	m := NewMarkdown(opts)

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := m.ConvertLinksToHyperlinks(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
