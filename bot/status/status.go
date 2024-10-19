package status

import (
	"errors"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

const (
	maxCustomStatusLength = 60
)

type Status struct {
	logger  *slog.Logger
	session *dgo.Session
}

func NewStatus(session *dgo.Session) *Status {
	return &Status{
		logger:  logger.NewWithSubsystem("bot", "status"),
		session: session,
	}
}

func (s *Status) setCustomActivity(text string) error {
	if len(text) > maxCustomStatusLength {
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

	return s.session.UpdateStatusComplex(dgo.UpdateStatusData{
		Activities: activities,
	})
}
