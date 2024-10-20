package env

import (
	"errors"
	"os"
)

type Stage string

const (
	StageProd Stage = "prod"
	StageDev  Stage = "dev"
)

func NewStageConfig() (Stage, error) {
	stage := os.Getenv("STAGE")
	if stage == "" {
		return "", errors.New("STAGE environment variable not set")
	}

	return Stage(stage), nil
}
