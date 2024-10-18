package yaml

import "errors"

type SolvedCommand struct {
	Description string            `yaml:"description"`
	Response    string            `yaml:"response"`
	SolvedTags  map[string]string `json:"solvedTags"`
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

	return nil
}

func NewSolvedCommand(config *Config) SolvedCommand {
	return config.SolvedCommand
}
