package battleship

// Board represents the state of a battleships board
type Board struct {
	layout map[byte][]string
}

// NewBoard populates a board and returns it
func NewBoard() *Board {
	board := &Board{
		layout: populateBoard(),
	}
	return board
}

func populateBoard() map[byte][]string {
	b := make(map[byte][]string)
	cols := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'}
	for i := 0; i < 9; i++ {
		b[cols[i]] = []string{"●", "●", "●", "●", "●", "●", "●", "●", "●"}
	}
	return b
}
