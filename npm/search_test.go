package npm_test

import (
	"github.com/erosdesire/discord-nestjs-utility-bot/npm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchNestJS(t *testing.T) {
	res, err := npm.Search(&npm.SearchOptions{
		Text: "@nestjs/",
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, res.Total, 10)
	// fmt.Printf("%+v\n", res.Objects)
}

func TestSearchNestJSWithOptions(t *testing.T) {
	res, err := npm.Search(&npm.SearchOptions{
		Text:        "@nestjs/",
		Popularity:  1,
		Maintenance: 1,
		Quality:     1,
		Size:        -50,
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, res.Total, 10)
	// fmt.Printf("%+v\n", res.Objects)
}
