package algolia_test

import (
	"github.com/nestjs-discord/utility-bot/internal/algolia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	res, err := algolia.Search(algolia.NestJS, "test")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, len(res), 20)
}

func TestGetObject(t *testing.T) {
	res, err := algolia.GetObject(algolia.NestJS, "1-https://docs.nestjs.com/standalone-applications")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "https://docs.nestjs.com/standalone-applications#getting-started", res.URL)
}
