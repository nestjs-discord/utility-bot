package yaml

import "errors"

type Moderators []string

func (m Moderators) validate() error {
	if len(m) == 0 {
		return errors.New("at least 1 moderator is required")
	}

	return nil
}

func NewModerators(config *Config) Moderators {
	return config.Moderators
}
