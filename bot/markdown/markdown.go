package markdown

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
	"os"
	"strings"
)

type Markdown struct {
	logger   *slog.Logger
	commands yaml.Commands
}

func NewMarkdown() *Markdown {
	return &Markdown{
		logger:   logger.NewWithSubsystem("bot", "markdown"),
		commands: make(yaml.Commands),
	}
}

// CacheCommands will cache Markdown content from the disk onto the memory
func (m *Markdown) CacheCommands(commands yaml.Commands) error {
	charLimit := 2000

	for cmdName, c := range commands { // TODO: concurrent loop
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

		c.Content = string(data)
		m.commands[cmdName] = c

		m.logger.Debug("cached content",
			slog.Int("char-len", len(c.Content)),
			slog.String("name", cmdName),
		)
	}

	return nil
}
