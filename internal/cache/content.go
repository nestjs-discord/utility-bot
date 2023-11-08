package cache

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

// Content will cache local Markdown content on memory
func Content() error {
	charLimit := 2000

	for _, c := range config.GetConfig().Commands {
		// Ignore non-markdown files
		if !strings.HasSuffix(c.Content, ".md") {
			return fmt.Errorf("expected '%v' file, to have '.md' extension", c.Content)
		}

		p := c.Content
		data, err := os.ReadFile(p)
		if err != nil {
			return errors.Wrapf(err, "failed to read file '%v'", p)
		}

		c.Content = string(data)

		// Slash commands can have a maximum of 4000 characters for combined name, description,
		// and value properties for each command, its options (including subcommands and groups), and choices.
		if len(c.Content) > charLimit {
			return fmt.Errorf("file '%v' contains too many characters, expected maximum of %v but received %v", p, charLimit, len(c.Content))
		}

		log.Debug().Str("path", p).Msg("cached file content")
	}

	return nil
}
