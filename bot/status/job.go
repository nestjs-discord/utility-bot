package status

import (
	"log/slog"
)

func (s *Status) ExecuteBackgroundJob() error {
	s.logger.Debug("executing background job")

	text := randomText()
	err := s.setCustomActivity(text)
	if err != nil {
		s.logger.Error("set custom activity failed",
			slog.Any("err", err),
		)
		return err
	}
	s.logger.Debug("updated",
		slog.String("text", text),
	)

	return nil
}
