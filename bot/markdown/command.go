package markdown

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"log"
	"log/slog"
	"os"
	"strings"
)

// CacheCommands will cache Markdown content from the disk into memory.
func (m *Markdown) CacheCommands(commands yaml.Commands) {
	// Loop through the commands and process each concurrently.
	for cmdName, c := range commands {
		m.processCommand(cmdName, c)
	}
}

// processCommand handles the caching of a single command.
// It checks file extensions, reads file content, validates size, and stores it.
func (m *Markdown) processCommand(cmdName string, c yaml.Command) {
	const charLimit = 2000

	// Ignore non-markdown files
	if err := m.validateFileExtension(c.Content); err != nil {
		log.Fatal(err)
	}

	// Read the file content
	data, err := os.ReadFile(c.Content)
	if err != nil {
		log.Fatalf("failed to read file '%v': %s", c.Content, err)
	}

	// Ensure content doesn't exceed the character limit
	if len(data) > charLimit {
		log.Fatalf("the '%s' file contains too many characters, expected %d, received %d", c.Content, charLimit, len(data))
	}

	// Store the command content in memory
	c.Content = string(data)
	m.commands[cmdName] = c

	// Log the caching result
	m.logger.Debug("cached content",
		slog.Int("char-len", len(c.Content)),
		slog.String("name", cmdName),
	)
}

// validateFileExtension checks if the file has a '.md' extension.
func (m *Markdown) validateFileExtension(filepath string) error {
	if !strings.HasSuffix(filepath, ".md") {
		return fmt.Errorf("expected '%v' file, to have '.md' extension", filepath)
	}
	return nil
}
