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

type Handler struct {
	logger       *slog.Logger
	forms        *forms.Forms
	moderators   *moderators.Moderators
	rateLimit    *rate_limit.RateLimit
	markdown     *markdown.Markdown
	archive      *archive.Archive
	credits      *credits.Credits
	googleIt     *google_it.GoogleIt
	reference    *reference.Reference
	solved       *solved.Solved
	dontPingMods *dont_ping_mods.DontPingMods
}

func NewHandler(
	forms *forms.Forms,
	moderators *moderators.Moderators,
	rateLimit *rate_limit.RateLimit,
	markdown *markdown.Markdown,
	archive *archive.Archive,
	credits *credits.Credits,
	googleIt *google_it.GoogleIt,
	reference *reference.Reference,
	solved *solved.Solved,
	dontPingMods *dont_ping_mods.DontPingMods,
) *Handler {
	return &Handler{
		logger:       logger.NewWithSubsystem("bot", "handler", "interaction"),
		forms:        forms,
		moderators:   moderators,
		rateLimit:    rateLimit,
		markdown:     markdown,
		archive:      archive,
		credits:      credits,
		googleIt:     googleIt,
		reference:    reference,
		solved:       solved,
		dontPingMods: dontPingMods,
	}
}
