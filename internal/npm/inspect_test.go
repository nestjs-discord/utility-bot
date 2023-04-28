package npm_test

import (
	"github.com/nestjs-discord/utility-bot/internal/npm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInspectNestJSCore(t *testing.T) {
	res, err := npm.Inspect(&npm.InspectOptions{
		Name: "@nestjs/core",
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Name, "@nestjs/core")
	// fmt.Printf("%+v\n", res.Objects)
}
