package battleship

import (
	"fmt"
	"math/rand"
)

type ship struct {
	length   int
	health   int
	position []string // e.g. []string{"A1, A2, A3"}
}

type ships struct {
	ships     map[string]*ship
	remaining int16
}

const (
	cClass = "Canberra"
	hClass = "Hobart"
	lClass = "Leeuwin"
	aClass = "Armidale"
)

func (ss *ships) setPositions() {
	ss.remaining = 4
	ss.ships = map[string]*ship{
		cClass: &ship{length: 5, health: 5},
		hClass: &ship{length: 4, health: 4},
		lClass: &ship{length: 3, health: 3},
		aClass: &ship{length: 2, health: 2},
	}

	for _, s := range ss.ships {
		var orientations []byte
		startRow := rand.Intn(9)

		switch {
		case startRow+s.length > 8:
			orientations = append(orientations, 'N')
		case startRow-s.length < 0:
			orientations = append(orientations, 'S')
		default:
			orientations = append(orientations, 'N', 'S')
		}

		startCol := rand.Intn(9)

		switch {
		case startCol+s.length > 8:
			orientations = append(orientations, 'W')
		case startCol-s.length < 0:
			orientations = append(orientations, 'E')
		default:
			orientations = append(orientations, 'E', 'W')
		}

		orientation := orientations[rand.Intn(len(orientations)-1)]
		col := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'}

		switch orientation {
		case 'N':
			for i := 0; i < s.length; i++ {
				s.position = append(s.position, fmt.Sprintf("%v%v", col[startCol], startRow+i))
			}
		case 'S':
			for i := 0; i < s.length; i++ {
				s.position = append(s.position, fmt.Sprintf("%v%v", col[startCol], startRow-i))
			}
		case 'E':
			for i := 0; i < s.length; i++ {
				s.position = append(s.position, fmt.Sprintf("%v%v", col[startCol-i], startRow))
			}
		case 'W':
			for i := 0; i < s.length; i++ {
				s.position = append(s.position, fmt.Sprintf("%v%v", col[startCol+i], startRow+i))
			}
		}
	}
}

func (ss *ships) checkCell(cell string) (bool, string) {
	for k, s := range ss.ships {
		for _, p := range s.position {
			if cell == p {
				return true, k
			}
		}
	}
	return false, ""
}

func (ss *ships) registerHit(cell string, class string) {
	ss.ships[class].health--
	if ss.ships[class].health == 0 {
		ss.remaining--
	}
}
