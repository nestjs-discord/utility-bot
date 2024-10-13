package config

type Yaml struct {
	Moderators []string            `mapstructure:"moderators" validate:"required,min=1,dive,min=1"`
	RateLimit  RateLimit           `mapstructure:"rateLimit" validate:"required"`
	AutoMod    AutoMod             `mapstructure:"autoMod" validate:"required"`
	Forms      map[string]Form     `mapstructure:"forms" validate:"required,min=1,dive"`
	Commands   map[string]*Command `mapstructure:"commands" validate:"required,max-one-space-allowed,min=1,max=85,dive"`
}

type Form struct {
	ButtonLabel     string      `mapstructure:"buttonLabel" validate:"required,min=5"`
	Title           string      `mapstructure:"modalTitle" validate:"required,min=10"`
	ChannelId       string      `mapstructure:"channelId" validate:"required,min=5"`
	ModChannelId    string      `mapstructure:"modChannelId" validate:"required,min=5"`
	ModSkipApproval bool        `mapstructure:"modSkipApproval"`
	Color           int         `mapstructure:"color" validate:"required"`
	Footer          string      `mapstructure:"footer" validate:"required"`
	Inputs          []FormInput `mapstructure:"inputs" validate:"required,min=1,max=10,dive"`
}

type FormInput struct {
	Id          string `mapstructure:"id" validate:"required,min=5"`
	Placeholder string `mapstructure:"placeholder" validate:"required,min=1,max=100"`
	Multiline   bool   `mapstructure:"multiline"`
	Min         int    `mapstructure:"min" validate:"min=0"`
	Max         int    `mapstructure:"max" validate:"min=0,max=1000"`
	Required    bool   `mapstructure:"required"`
}

type RateLimit struct {
	TTL     int    `mapstructure:"ttl" validate:"required,min=1"`
	Usage   int    `mapstructure:"usage" validate:"required,min=2"`
	Message string `mapstructure:"message" validate:"required,min=3"`
}

type AutoMod struct {
	Enabled                 bool   `mapstructure:"enabled" validate:"boolean"`
	ModeratorsBypass        bool   `mapstructure:"moderatorsBypass" validate:"boolean"`
	LogChannelId            string `mapstructure:"logChannelId" validate:"required,min=1"`
	LogMentionRoleId        string `mapstructure:"logMentionRoleId"`
	MessageTTL              int    `mapstructure:"messageTTL" validate:"required,min=1"`
	MaxChannelsLimitPerUser int    `mapstructure:"maxChannelsLimitPerUser" validate:"required,min=1"`
	DenyTTL                 int    `mapstructure:"denyTTL" validate:"required,min=1"`
}

type Command struct {
	Description string      `mapstructure:"description" validate:"required,min=1,max=100"`
	Content     string      `mapstructure:"content" validate:"required,min=1"`
	Protected   bool        `mapstructure:"protected" validate:"boolean"`
	Buttons     [][]*Button `mapstructure:"buttons" validate:"min=0,max=8,dive,min=1,max=4,dive"`
}

type Button struct {
	Label string `mapstructure:"label" validate:"required,min=3,max=40"`
	URL   string `mapstructure:"url" validate:"required,url,min=3"`
	Emoji string `mapstructure:"emoji" validate:"regexp=^[\p{Emoji}]$"`
}
