package yaml

// TODO: remove all the "validate" tags and perform manual validation

type Config struct {
	Moderators Moderators `yaml:"moderators" validate:"required,min=1,dive,min=1"`
	RateLimit  RateLimit  `yaml:"rateLimit" validate:"required"`
	Antispam   Antispam   `yaml:"antispam" validate:"required"`
	Forms      Forms      `yaml:"forms" validate:"required,min=1,dive"`
	Commands   Commands   `yaml:"commands" validate:"required,max-one-space-allowed,min=1,max=85,dive"`
}
type Path string
type Moderators []string
type Forms map[string]Form
type Form struct {
	ChannelId        string               `yaml:"channelId" validate:"required,min=5"`
	ModChannelId     string               `yaml:"modChannelId" validate:"required,min=5"`
	ModalTitle       string               `yaml:"modalTitle" validate:"required,min=10,max=45"` // TODO: validate title of the popup modal, max 45 characters
	ModSkipApproval  bool                 `yaml:"modSkipApproval"`
	Color            int                  `yaml:"color" validate:"required"`
	OpenModalMessage FormOpenModalMessage `yaml:"openModalMessage"`
	Footer           string               `yaml:"footer" validate:"required"`
	Inputs           []FormInput          `yaml:"inputs" validate:"required,min=1,max=5,dive"` // TODO: Between 1 and 5 (inclusive) components that make up the modal
}
type FormOpenModalMessage struct {
	ButtonLabel      string `yaml:"buttonLabel"`
	EmbedColor       int    `yaml:"embedColor"`
	EmbedTitle       string `yaml:"embedTitle"`
	EmbedDescription string `yaml:"embedDescription"`
}
type FormInput struct {
	Id          string `yaml:"id" validate:"required,min=5"`
	Placeholder string `yaml:"placeholder" validate:"required,min=1,max=100"`
	Multiline   bool   `yaml:"multiline"`
	Min         int    `yaml:"min" validate:"min=0"`
	Max         int    `yaml:"max" validate:"min=0,max=1000"`
	Required    bool   `yaml:"required"`
}
type RateLimit struct {
	TTLSec   int    `yaml:"ttlSec" validate:"required,min=1"`
	MaxUsage int    `yaml:"maxUsage" validate:"required,min=2"`
	Message  string `yaml:"message" validate:"required,min=3"`
}
type Antispam struct {
	Enabled            bool     `yaml:"enabled" validate:"boolean"`
	ModeratorsBypass   bool     `yaml:"moderatorsBypass" validate:"boolean"`
	LogChannelId       string   `yaml:"logChannelId" validate:"required,min=1"`
	MessageTTLSec      int      `yaml:"messageTTLSec" validate:"required,min=1"`
	MaxChannelsPerUser int      `yaml:"maxChannelsPerUser" validate:"required,min=1"`
	DenyTTLSec         int      `yaml:"denyTTLSec" validate:"required,min=1"`
	TrackedChannelIds  []string `yaml:"trackedChannelIds"`
}
type Commands map[string]Command
type Command struct {
	Description string         `yaml:"description" validate:"required,min=1,max=100"`
	Content     string         `yaml:"content" validate:"required,min=1"`
	Protected   bool           `yaml:"protected" validate:"boolean"`
	Buttons     CommandButtons `yaml:"buttons" validate:"min=0,max=8,dive,min=1,max=4,dive"`
}
type CommandButtons [][]*CommandButton
type CommandButton struct {
	Label string `yaml:"label" validate:"required,min=3,max=40"`
	URL   string `yaml:"url" validate:"required,url,min=3"`
	Emoji string `yaml:"emoji"` // TODO: validate
}
