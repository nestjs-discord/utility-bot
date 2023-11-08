package algolia_test

import (
	algolia2 "github.com/nestjs-discord/utility-bot/pkg/algolia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchNestJS(t *testing.T) {
	res, err := algolia2.Search(algolia2.NestJS, "test")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, len(res), 20)
}

func TestGetObject(t *testing.T) {
	res, err := algolia2.GetObject(algolia2.NestJS, "1-https://docs.nestjs.com/standalone-applications")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "https://docs.nestjs.com/standalone-applications#getting-started", res.URL)
}
