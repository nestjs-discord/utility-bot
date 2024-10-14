package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var YamlFile = "./config.yml"

type YamlConfig struct {
	Moderators YamlModerators `yaml:"moderators" validate:"required,min=1,dive,min=1"`
	RateLimit  YamlRateLimit  `yaml:"rateLimit" validate:"required"`
	AutoMod    YamlAutoMod    `yaml:"autoMod" validate:"required"`
	Forms      YamlForms      `yaml:"forms" validate:"required,min=1,dive"`
	Commands   YamlCommands   `yaml:"commands" validate:"required,max-one-space-allowed,min=1,max=85,dive"`
}

func NewYamlConfig(path string) (*YamlConfig, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read the yaml config: %s", err)
	}

	var data YamlConfig
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal the yaml config: %s", err)
	}

	// TODO: manual validation

	return &data, nil
}

type YamlModerators []string
type YamlForms map[string]YamlForm

type YamlForm struct {
	ButtonLabel     string          `yaml:"buttonLabel" validate:"required,min=5"`
	Title           string          `yaml:"modalTitle" validate:"required,min=10"`
	ChannelId       string          `yaml:"channelId" validate:"required,min=5"`
	ModChannelId    string          `yaml:"modChannelId" validate:"required,min=5"`
	ModSkipApproval bool            `yaml:"modSkipApproval"`
	Color           int             `yaml:"color" validate:"required"`
	Footer          string          `yaml:"footer" validate:"required"`
	Inputs          []YamlFormInput `yaml:"inputs" validate:"required,min=1,max=10,dive"`
}

type YamlFormInput struct {
	Id          string `yaml:"id" validate:"required,min=5"`
	Placeholder string `yaml:"placeholder" validate:"required,min=1,max=100"`
	Multiline   bool   `yaml:"multiline"`
	Min         int    `yaml:"min" validate:"min=0"`
	Max         int    `yaml:"max" validate:"min=0,max=1000"`
	Required    bool   `yaml:"required"`
}

type YamlRateLimit struct {
	TTL     int    `yaml:"ttl" validate:"required,min=1"`
	Usage   int    `yaml:"usage" validate:"required,min=2"`
	Message string `yaml:"message" validate:"required,min=3"`
}

type YamlAutoMod struct {
	Enabled                 bool   `yaml:"enabled" validate:"boolean"`
	ModeratorsBypass        bool   `yaml:"moderatorsBypass" validate:"boolean"`
	LogChannelId            string `yaml:"logChannelId" validate:"required,min=1"`
	LogMentionRoleId        string `yaml:"logMentionRoleId"`
	MessageTTL              int    `yaml:"messageTTL" validate:"required,min=1"`
	MaxChannelsLimitPerUser int    `yaml:"maxChannelsLimitPerUser" validate:"required,min=1"`
	DenyTTL                 int    `yaml:"denyTTL" validate:"required,min=1"`
}

type YamlCommands map[string]YamlCommand

type YamlCommand struct {
	Description string                 `yaml:"description" validate:"required,min=1,max=100"`
	Content     string                 `yaml:"content" validate:"required,min=1"`
	Protected   bool                   `yaml:"protected" validate:"boolean"`
	Buttons     [][]*YamlCommandButton `yaml:"buttons" validate:"min=0,max=8,dive,min=1,max=4,dive"`
}

type YamlCommandButton struct {
	Label string `yaml:"label" validate:"required,min=3,max=40"`
	URL   string `yaml:"url" validate:"required,url,min=3"`
	Emoji string `yaml:"emoji" validate:"regexp=^[\p{Emoji}]$"`
}
