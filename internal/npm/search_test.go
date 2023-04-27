package npm_test

import (
	"github.com/nestjs-discord/utility-bot/internal/npm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchNestJS(t *testing.T) {
	res, err := npm.Search(&npm.SearchOptions{
		Text: "scope:nestjs",
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, res.Total, int64(5), "should return at least 5 packages")
	// fmt.Printf("%+v\n", res.Objects)
}

func TestSearchNestJSWithOptions(t *testing.T) {
	res, err := npm.Search(&npm.SearchOptions{
		Text:        "scope:nestjs",
		Popularity:  1,
		Maintenance: 1,
		Quality:     1,
		Size:        10,
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 10, len(res.Objects), "should return 10 packages only")
	assert.GreaterOrEqual(t, res.Total, int64(10), "should return at least 10 packages")
	// fmt.Printf("%+v\n", res.Objects)
}
