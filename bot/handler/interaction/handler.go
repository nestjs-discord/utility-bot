package interaction

import (
	"github.com/nestjs-discord/utility-bot/bot/commands/archive"
	"github.com/nestjs-discord/utility-bot/bot/commands/credits"
	"github.com/nestjs-discord/utility-bot/bot/commands/dont_ping_mods"
	"github.com/nestjs-discord/utility-bot/bot/commands/google_it"
	"github.com/nestjs-discord/utility-bot/bot/commands/reference"
	"github.com/nestjs-discord/utility-bot/bot/commands/solved"
	"github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/bot/markdown"
	"github.com/nestjs-discord/utility-bot/bot/rate_limit"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"

	"github.com/nestjs-discord/utility-bot/bot/moderators"
)

type Options struct {
	Forms        *forms.Forms
	Moderators   *moderators.Moderators
	RateLimit    *rate_limit.RateLimit
	Markdown     *markdown.Markdown
	Archive      *archive.Archive
	Credits      *credits.Credits
	GoogleIt     *google_it.GoogleIt
	Reference    *reference.Reference
	Solved       *solved.Solved
	DontPingMods *dont_ping_mods.DontPingMods
}

type Handler struct {
	logger *slog.Logger
	opts   Options
}

func NewInteractionHandler(opts Options) *Handler {
	return &Handler{
		logger: logger.NewWithSubsystem("bot", "handler", "interaction"),
		opts:   opts,
	}
}
