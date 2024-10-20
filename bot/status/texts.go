package status

import (
	"golang.org/x/exp/rand"
)

var Texts []string

func randomText() string {
	return Texts[rand.Intn(len(Texts))]
}
