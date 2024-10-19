package cron

import (
	"fmt"
	gc "github.com/go-co-op/gocron/v2"
	"time"
)

func (c *Cron) registerJobs() error {
	_, err := c.scheduler.NewJob(
		gc.DurationJob(1*time.Hour),
		gc.NewTask(c.opts.AutoMod.ExecuteBackgroundJob),
		gc.WithSingletonMode(gc.LimitModeReschedule),
		gc.WithStartAt(gc.WithStartImmediately()),
		gc.WithTags("autoMod"),
	)
	if err != nil {
		return fmt.Errorf("autoMod: %v", err)
	}

	_, err = c.scheduler.NewJob(
		gc.DurationJob(10*time.Minute),
		gc.NewTask(c.opts.Status.ExecuteBackgroundJob),
		gc.WithSingletonMode(gc.LimitModeReschedule),
		gc.WithStartAt(gc.WithStartImmediately()),
		gc.WithTags("status"),
	)
	if err != nil {
		return fmt.Errorf("status: %v", err)
	}

	return nil
}
