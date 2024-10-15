package markdown

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/config/yaml"
	"github.com/nestjs-discord/utility-bot/logger"
	"log/slog"
	"os"
	"strings"
)

type Markdown struct {
	logger *slog.Logger
	data   map[string]string
}

func NewMarkdown() *Markdown {
	return &Markdown{
		logger: logger.NewWithSubsystem("bot", "markdown"),
		data:   make(map[string]string),
	}
}

// CacheCommands will cache Markdown content from the disk onto the memory
func (m *Markdown) CacheCommands(commands yaml.Commands) error {
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
			return fmt.Errorf("the '%v' file contains too many characters, expected maximum of %d but received %d", p, charLimit, len(c.Content))
		}

		m.data[c.Content] = string(data)

		m.logger.Debug("cached content",
			slog.Int("char-len", len(m.data[c.Content])),
			slog.String("path", p),
		)
	}

	return nil
}
