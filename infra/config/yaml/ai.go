package yaml

import "errors"

type AI struct {
	Enabled bool     `yaml:"enabled"`
	Forums  []string `yaml:"forums"`
}

func (a AI) validate() error {
	if !a.Enabled {
		return nil
	}

	if len(a.Forums) == 0 {
		return errors.New("at least one forum channel id must be specified when AI is enabled")
	}

	for _, forum := range a.Forums {
		if forum == "" {
			return errors.New("forum channel id cannot be empty")
		}
	}

	return nil
}

func NewAI(config *Config) AI {
	return config.AI
}
