package cache

import (
	"github.com/nestjs-discord/utility-bot/config/yaml"
	"log/slog"

	"fmt"
	"os"
	"strings"
)

// DynamicContent is a map of command key to the markdown content
var DynamicContent = make(map[string]string) // TODO: use this in a private struct + where the dynamic command is being handled

// MarkdownContent will cache Markdown content from the disk onto the memory
func MarkdownContent(commands yaml.Commands) error {
	charLimit := 2000

	for _, c := range commands {
		// Ignore non-markdown files
		if !strings.HasSuffix(c.Content, ".md") {
			return fmt.Errorf("expected '%v' file, to have '.md' extension", c.Content)
		}

		p := c.Content
		data, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("failed to read file '%v': %s", p, err)
		}

		// Slash commands can have a maximum of 4000 characters for combined name, description,
		// and value properties for each command, its options (including subcommands and groups), and choices.
		if len(c.Content) > charLimit {
			return fmt.Errorf("file '%v' contains too many characters, expected maximum of %v but received %v", p, charLimit, len(c.Content))
		}

		strData := string(data)

		DynamicContent[c.Content] = strData

		slog.Debug("cached file content",
			slog.Int("char-len", len(strData)),
			slog.String("path", p),
		)
	}

	return nil
}
