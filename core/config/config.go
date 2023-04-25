package config

var c config

func GetConfig() *config {
	return &c
}

type config struct {
	GuildID  string              `mapstructure:"guild_id" validate:"required,min=1"`
	Commands map[string]*Command `mapstructure:"commands" validate:"required,min=1,dive"`
}

type Command struct {
	Description string `mapstructure:"description" validate:"required,min=1"`
	Content     string `mapstructure:"content" validate:"required,min=1"`
}
