package yaml

import (
	"errors"
	"fmt"
	"github.com/forPelevin/gomoji"
	"net/url"
	"strings"
)

type Commands map[string]Command

func (c Commands) validate() error {
	for key, value := range c {
		if len(key) < 6 {
			return fmt.Errorf("'%s' length must be at least 6 characters", key)
		}

		keyParts := strings.Split(key, " ")

		// max one space character is allowed
		if len(keyParts) > 2 {
			return fmt.Errorf(
				"max depth of allowed sub-command for the '%s' key is 1, remove extra space characters",
				key,
			)
		}

		for _, keyPart := range keyParts {
			if len(keyPart) > 30 {
				return fmt.Errorf("'%s' length must be less than 30 characters", keyPart)
			}
		}

		err := value.validate()
		if err != nil {
			return fmt.Errorf("'%s' is invalid: %w", key, err)
		}
	}

	return nil
}

type Command struct {
	Description string         `yaml:"description"`
	Content     string         `yaml:"content"`
	Protected   bool           `yaml:"protected"`
	Buttons     CommandButtons `yaml:"buttons"`
}

func (c Command) validate() error {
	if len(c.Description) < 20 {
		return fmt.Errorf("description is too short")
	}
	if len(c.Description) > 100 {
		return fmt.Errorf("description is too long")
	}

	if len(c.Content) < 30 {
		return fmt.Errorf("content is too short")
	}
	if len(c.Content) > 2000 {
		return fmt.Errorf("content is too long")
	}

	if len(c.Buttons) == 0 {
		return nil // skip further validation
	}
	if len(c.Buttons) > 8 {
		return fmt.Errorf("buttons can have maximum of 8 rows")
	}
	for rowIndex, row := range c.Buttons {
		if len(row) == 0 {
			return fmt.Errorf("buttons[%d] is empty", rowIndex)
		}

		if len(row) > 4 {
			return fmt.Errorf("buttons[%d] have too many items, maximum allowed is 4", rowIndex)
		}

		for buttonIndex, button := range row {
			err := button.validate()
			if err != nil {
				return fmt.Errorf("buttons[%d][%d] is invalid: %w", rowIndex, buttonIndex, err)
			}
		}
	}

	return nil
}

type CommandButtons [][]CommandButton

type CommandButton struct {
	Label string `yaml:"label"`
	URL   string `yaml:"url"`
	Emoji string `yaml:"emoji"`
}

func (c CommandButton) validate() error {
	if c.Label == "" {
		return errors.New("label is required")
	}
	if len(c.Label) < 3 {
		return errors.New("label is too short")
	}
	if len(c.Label) > 40 {
		return errors.New("label is too long")
	}

	if c.URL == "" {
		return errors.New("url is required")
	}
	_, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("url is invalid: %s", err)
	}

	if c.Emoji == "" {
		return errors.New("emoji is required")
	}
	if !gomoji.ContainsEmoji(c.Emoji) {
		return errors.New("invalid emoji")
	}
	if len(gomoji.FindAll(c.Emoji)) != 1 {
		return errors.New("only 1 emoji is allowed")
	}

	return nil
}

func NewCommands(config *Config) Commands {
	return config.Commands
}
