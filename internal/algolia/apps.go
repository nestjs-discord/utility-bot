package algolia

import (
	"strings"
)

// App represents an Algolia application.
type App string

// ToSlug converts the App to a slug format.
// It converts the App to lowercase and replaces spaces with hyphens.
func (a App) ToSlug() string {
	return strings.ReplaceAll(strings.ToLower(string(a)), " ", "-")
}

// Predefined App constants representing different applications.
const (
	Discord        App = "Discord"
	DiscordJSGuide App = "DiscordJS Guide"
	Express        App = "Express"
	Fastify        App = "Fastify"
	Necord         App = "NECORD"
	NestCommander  App = "Nest Commander"
	NestJS         App = "NestJS"
	Ogma           App = "Ogma"
	TypeScript     App = "TypeScript"
)

// Apps is a map that associates App slugs with their corresponding App constants.
var Apps = map[string]App{
	Discord.ToSlug():        Discord,
	DiscordJSGuide.ToSlug(): DiscordJSGuide,
	Express.ToSlug():        Express,
	Fastify.ToSlug():        Fastify,
	Necord.ToSlug():         Necord,
	NestCommander.ToSlug():  NestCommander,
	NestJS.ToSlug():         NestJS,
	Ogma.ToSlug():           Ogma,
	TypeScript.ToSlug():     TypeScript,
}
