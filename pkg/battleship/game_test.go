package battleship

import (
	"fmt"
	"math/rand"
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

func TestHit(t *testing.T) {
	game := NewGame()
	c := game.battleships.allCells[0]
	h, _ := game.Fire(c)
	assert.True(t, h, "game.Fire() should be true")
}

func TestMiss(t *testing.T) {
	game := NewGame()
	for {
		c := string(rand.Intn(len(cols))) + string(rand.Intn(len(rows)))
		if !contains(game.battleships.allCells, c) {
			h, _ := game.Fire(c)
			assert.False(t, h, "game.Fire() should be false")
			return
		}
	}
}

func TestGameOver(t *testing.T) {
	game := NewGame()
	for _, p := range game.battleships.allCells {
		h, err := game.Fire(p)
		if err != nil {
			fmt.Printf("Err")
		}
		if h == false {
			fmt.Printf("Hello")
		}
		assert.True(t, h, "should be a hit")
	}
	assert.True(t, game.IsGameOver(), "game should be over")
}
