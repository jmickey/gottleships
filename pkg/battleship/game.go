package battleship

import (
	"fmt"
	"sort"
)

// Game represents an instance of a game
type Game struct {
	sh        *ships
	shots     int
	attempted []string
}

// NewGame initiates a new game instance and returns it
func NewGame() *Game {
	sh := &ships{}
	sh.setPositions()
	return &Game{sh: sh}
}

// Fire checks the position provided and returns a bool
// It will return an error if the received position has already been attempted
func (g *Game) Fire(cell string) (bool, error) {
	if !g.isValid(cell) {
		return false, fmt.Errorf("cell '%s' has already been attempted", cell)
	}
	g.shots++
	g.attempted = append(g.attempted, cell)
	if hit, class := g.sh.checkCell(cell); hit {
		g.sh.registerHit(cell, class)
		return true, nil
	}
	return false, nil
}

// Checks if the cell has already been attempted
func (g *Game) isValid(cell string) bool {
	sort.Strings(g.attempted)
	i := sort.SearchStrings(g.attempted, cell)
	if i < len(g.attempted) && g.attempted[i] == cell {
		return false
	}
	return true
}

// GetRemaining returns the number of remaining ships
func (g *Game) GetRemaining() int16 {
	return g.sh.remaining
}

// IsGameOver checks if any ships remain and returns a bool
func (g *Game) IsGameOver() bool {
	if g.sh.remaining == 0 {
		return true
	}
	return false
}

// Shots returns the number of shots taken since the start of the game
func (g *Game) Shots() int {
	return g.shots
}
