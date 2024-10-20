package yaml

import (
	"errors"
	"fmt"
)

type Forms map[string]Form

func (f Forms) validate() error {
	for key, value := range f {
		// key
		if len(key) > 2 {
			return fmt.Errorf("key '%s' must not exceed 2 characters", key)
		}

		// value
		err := value.validate()
		if err != nil {
			return fmt.Errorf("key '%s' has error: %s", key, err)
		}
	}
	return nil
}

type Form struct {
	ChannelId             string               `yaml:"channelId"`
	ModChannelId          string               `yaml:"modChannelId"`
	ModalTitle            string               `yaml:"modalTitle"`
	ModSkipApproval       bool                 `yaml:"modSkipApproval"`
	Color                 int                  `yaml:"color"`
	MinimumAccountAgeDays int                  `yaml:"minimumAccountAgeDays"`
	MinimumServerJoinDays int                  `yaml:"minimumServerJoinDays"`
	OpenModalMessage      FormOpenModalMessage `yaml:"openModalMessage"`
	Footer                string               `yaml:"footer"`
	Inputs                []FormInput          `yaml:"inputs"`
}

func (f Form) validate() error {
	if len(f.ChannelId) < 5 {
		return errors.New("channelId is too short")
	}

	if len(f.ModChannelId) < 5 {
		return errors.New("modChannelId is too short")
	}

	if len(f.ModalTitle) < 10 {
		return errors.New("modalTitle is too short")
	}
	if len(f.ModalTitle) > 40 {
		return errors.New("modalTitle is too long")
	}

	if f.Color == 0 {
		return errors.New("color must be greater than zero")
	}

	if f.MinimumAccountAgeDays < 3 {
		return errors.New("minimumAccountAgeDays must be greater than 3")
	}
	if f.MinimumServerJoinDays < 3 {
		return errors.New("minimumServerJoinDays must be greater than 3")
	}

	err := f.OpenModalMessage.validate()
	if err != nil {
		return fmt.Errorf("openModalMessage has error: %s", err)
	}

	if f.Footer == "" {
		return errors.New("footer is required")
	}

	if len(f.Inputs) == 0 {
		return errors.New("inputs is required")
	}

	if len(f.Inputs) > 5 {
		return errors.New("inputs must not exceed 5")
	}

	var totalInputsLength = 0

	for index, input := range f.Inputs {
		err = input.validate()
		if err != nil {
			return fmt.Errorf("inputs[%d] is invalid: %s", index, err)
		}
		totalInputsLength += input.Max
	}

	if totalInputsLength > 3000 {
		return errors.New("inputs must not exceed 3000 characters")
	}

	return nil
}

type FormOpenModalMessage struct {
	ButtonLabel      string `yaml:"buttonLabel"`
	EmbedColor       int    `yaml:"embedColor"`
	EmbedTitle       string `yaml:"embedTitle"`
	EmbedDescription string `yaml:"embedDescription"`
}

func (f FormOpenModalMessage) validate() error {
	if f.ButtonLabel == "" {
		return errors.New("buttonLabel is required")
	}
	if len(f.ButtonLabel) > 70 {
		return errors.New("buttonLabel is too long")
	}

	if f.EmbedColor == 0 {
		return errors.New("embedColor must be greater than zero")
	}

	if f.EmbedTitle == "" {
		return errors.New("embedTitle is required")
	}

	if f.EmbedDescription == "" {
		return errors.New("embedDescription is required")
	}

	if len(f.EmbedDescription) > 2000 {
		return errors.New("embedDescription is too long")
	}

	return nil
}

type FormInput struct {
	Id          string `yaml:"id"`
	Label       string `yaml:"label"`
	Placeholder string `yaml:"placeholder"`
	Multiline   bool   `yaml:"multiline"`
	Min         int    `yaml:"min"`
	Max         int    `yaml:"max"`
	Required    bool   `yaml:"required"`
}

func (f FormInput) validate() error {
	if len(f.Id) < 5 {
		return errors.New("id is too short")
	}
	if len(f.Id) > 40 {
		return errors.New("id is too long")
	}

	if len(f.Label) < 5 {
		return errors.New("label is too short")
	}
	if len(f.Label) > 40 {
		return errors.New("label is too long")
	}

	if f.Placeholder == "" {
		return errors.New("placeholder is required")
	}
	if len(f.Placeholder) < 5 {
		return errors.New("placeholder is too short")
	}
	if len(f.Placeholder) > 100 {
		return errors.New("placeholder is too long")
	}

	if f.Min < 0 {
		return errors.New("min must be greater or equal to 0")
	}
	if f.Max < 0 {
		return errors.New("max must be greater or equal to 0")
	}
	if f.Max > 1000 {
		return errors.New("max must be less than 1000")
	}
	if f.Max < f.Min {
		return errors.New("max must be greater or equal to min")
	}

	return nil
}

func NewForms(config *Config) Forms {
	return config.Forms
}
