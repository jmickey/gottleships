package battleship

// Game represents an instance of a game
type Game struct {
	sh *ships
}

// NewGame creates a new instance of a game,
// and sets the ships positions
func NewGame() *Game {
	sh := &ships{}
	sh.setPositions()
	return &Game{sh: sh}
}

// Fire checks the position povided and returns a bool
func (g *Game) Fire(cell string) bool {
	if hit, class := g.sh.checkCell(cell); hit {
		g.sh.registerHit(cell, class)
		return true
	}
	return false
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
