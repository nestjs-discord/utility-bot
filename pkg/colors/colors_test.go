package colors_test

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"scan-hub/cmd/colors"
	"testing"
)

func TestColors(t *testing.T) {
	t.Run("Test Purple Color", func(t *testing.T) {
		assert.Equal(t, lipgloss.Color("5"), colors.Purple.GetForeground())
	})

	t.Run("Test Red Color", func(t *testing.T) {
		assert.Equal(t, lipgloss.Color("9"), colors.Red.GetForeground())
	})

	t.Run("Test Yellow Color", func(t *testing.T) {
		assert.Equal(t, lipgloss.Color("11"), colors.Yellow.GetForeground())
	})

	t.Run("Test Blue Color", func(t *testing.T) {
		assert.Equal(t, lipgloss.Color("12"), colors.Blue.GetForeground())
	})
}
