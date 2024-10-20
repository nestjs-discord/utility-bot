package status

import (
	"errors"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

const (
	MaxCustomStatusLength = 60
)

type Options struct {
	Session *dgo.Session
}

type Status struct {
	opts   Options
	logger *slog.Logger
}

func NewStatus(opts Options) *Status {
	return &Status{
		opts:   opts,
		logger: logger.NewWithSubsystem("bot", "status"),
	}
}

func (s *Status) setCustomActivity(text string) error {
	if len(text) > MaxCustomStatusLength {
		return errors.New("text too long")
	}

	activities := []*dgo.Activity{
		{
			//Name: "NestJS 😎🍿",
			//Type: dgo.ActivityTypeStreaming,
			//// URL: "https://www.youtube.com/watch?v=0M8AYU_hPas",
			//// URL: "https://www.twitch.tv/directory/all/tags/nestjs",
			//
			//// Revealing framework fundamentals: NestJS behind the curtain
			//URL: "https://www.youtube.com/watch?v=jo-1EUxMmxc",
			Type:  dgo.ActivityTypeCustom,
			Name:  "custom",
			State: text,
		},
	}

	return s.opts.Session.UpdateStatusComplex(dgo.UpdateStatusData{
		Activities: activities,
	})
}
