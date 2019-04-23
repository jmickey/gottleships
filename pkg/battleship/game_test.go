package battleship

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorOnRepeatedMessages(t *testing.T) {
	game := NewGame()
	cell := "E5"
	_, err := game.Fire(cell)
	assert.Nil(t, err, "Error on game.Fire() was not nil")
	_, err = game.Fire(cell)
	assert.Error(t, err, "Error was nil when firing on same cell twice")
}
