package yaml

import (
	"errors"
	"fmt"
)

type SolvedCommand struct {
	Description      string           `yaml:"description"`
	Response         string           `yaml:"response"`
	ChannelSolvedTag ChannelSolvedTag `yaml:"channelSolvedTag"`
}

func (s SolvedCommand) validate() error {
	if s.Description == "" {
		return errors.New("description is required")
	}
	if len(s.Description) > 60 {
		return errors.New("description is too long")
	}

	if s.Response == "" {
		return errors.New("response is required")
	}
	if len(s.Response) > 1000 {
		return errors.New("response is too long")
	}

	if len(s.ChannelSolvedTag) == 0 {
		return errors.New("ChannelSolvedTag is required")
	}

	err := s.ChannelSolvedTag.validate()
	if err != nil {
		return fmt.Errorf("channelSolvedTag is invalid: %s", err)
	}

	return nil
}

type ChannelSolvedTag map[string]string

func (ch ChannelSolvedTag) validate() error {
	for key, value := range ch {
		if len(key) > 100 {
			return fmt.Errorf("the key '%s' is too long", key)
		}

		if value == "" {
			return fmt.Errorf("the '%s' key does not contain any value", key)
		}
		if len(value) > 100 {
			return fmt.Errorf("the '%s' key's value is too long", key)
		}
	}

	return nil
}

func NewSolvedCommand(config *Config) SolvedCommand {
	return config.SolvedCommand
}
