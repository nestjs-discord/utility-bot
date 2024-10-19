package cron

import (
	"fmt"
	gc "github.com/go-co-op/gocron/v2"
	"github.com/nestjs-discord/utility-bot/bot/auto_mod"
	"github.com/nestjs-discord/utility-bot/bot/status"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

type Option struct {
	AutoMod *auto_mod.AutoMod
	Status  *status.Status
}

type Cron struct {
	opts      Option
	logger    *slog.Logger
	scheduler gc.Scheduler
}

func NewCron(opts Option) (*Cron, func(), error) {
	s, err := gc.NewScheduler()
	if err != nil {
		return nil, nil, fmt.Errorf("gocron.NewScheduler failed: %v", err)
	}

	c := &Cron{
		opts:      opts,
		logger:    logger.NewWithSubsystem("bot", "cron"),
		scheduler: s,
	}

	err = c.registerJobs()
	if err != nil {
		return nil, nil, fmt.Errorf("register jobs failed: %v", err)
	}

	cleanup := func() {
		if err = c.Shutdown(); err != nil {
			c.logger.Error("shutdown failed",
				slog.Any("err", err),
			)
		}
	}

	return c, cleanup, nil
}

func (c *Cron) Start() {
	c.scheduler.Start()
}

func (c *Cron) Shutdown() error {
	return c.scheduler.Shutdown()
}
