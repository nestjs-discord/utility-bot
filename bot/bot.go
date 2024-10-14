package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler"
	"github.com/nestjs-discord/utility-bot/internal/logger"
	"log/slog"
)

type OptFunc func(*Options)

type Options struct {
	token   string
	appId   string
	guildId string
	handler *handler.Handler
}

func defaultOpts() Options {
	return Options{}
}

func WithToken(token string) OptFunc {
	return func(o *Options) {
		o.token = token
	}
}

func WithAppId(appId string) OptFunc {
	return func(o *Options) {
		o.appId = appId
	}
}

func WithGuildId(guildId string) OptFunc {
	return func(o *Options) {
		o.guildId = guildId
	}
}

func WithHandler(handler *handler.Handler) OptFunc {
	return func(o *Options) {
		o.handler = handler
	}
}

type Bot struct {
	opts    Options
	session *discordgo.Session
	logger  *slog.Logger
}

func NewBot(opts ...OptFunc) (*Bot, error) {
	o := defaultOpts()
	for _, opt := range opts {
		opt(&o)
	}

	bot := &Bot{
		opts:   o,
		logger: logger.NewWithSubsystem("bot"),
	}

	err := bot.newSession()
	if err != nil {
		return nil, err
	}

	bot.addHandlers()

	return bot, nil
}

func (b *Bot) addHandlers() {
	b.session.AddHandler(b.opts.handler.Ready)
	b.session.AddHandler(b.opts.handler.InteractionCreate)
	b.session.AddHandler(b.opts.handler.MessageCreate)
}
