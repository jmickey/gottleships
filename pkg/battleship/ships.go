package battleship

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

type ship struct {
	length   int
	health   int
	position []string // e.g. []string{"A1, A2, A3"}
}

type ships struct {
	ships     map[string]*ship
	remaining int16
	allCells  []string
}

const (
	cClass = "Canberra"
	hClass = "Hobart"
	lClass = "Leeuwin"
	aClass = "Armidale"
	cols   = "ABCDEFGHI"
	rows   = "123456789"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Details of algorithm
func (s *ships) setPositions() {
	s.remaining = 4
	s.ships = map[string]*ship{
		cClass: &ship{length: 5, health: 5},
		hClass: &ship{length: 4, health: 4},
		lClass: &ship{length: 3, health: 3},
		aClass: &ship{length: 2, health: 2},
	}

	for _, sh := range s.ships {
		var retry bool
		// Acts as a do...while loop
		for cont := true; cont; cont = retry {
			retry = true
			var orientations []byte
			startRow := rand.Intn(9) + 1
			startCol := rand.Intn(9)

			if startRow-sh.length >= 0 {
				orientations = append(orientations, 'N')
			}
			if startRow+sh.length <= 10 {
				orientations = append(orientations, 'S')
			}
			if (startCol+sh.length)+1 <= 10 {
				orientations = append(orientations, 'E')
			}
			if (startCol-sh.length)+1 >= 0 {
				orientations = append(orientations, 'W')
			}

			var failed bool
			for _cont := true; _cont; _cont = failed {
				orientation := orientations[rand.Intn(len(orientations))]
				failed = false
				pos := []string{}

				switch orientation {
				case 'N':
					for i := 0; i < sh.length; i++ {
						p := fmt.Sprintf("%v%v", string(cols[startCol]), startRow-i)
						if contains(s.allCells, p) {
							in := bytes.IndexByte(orientations, orientation)
							orientations = append(orientations[:in], orientations[in+1:]...)
							failed = true
							break
						}
						pos = append(pos, p)
					}
				case 'S':
					for i := 0; i < sh.length; i++ {
						p := fmt.Sprintf("%v%v", string(cols[startCol]), startRow+i)
						if contains(s.allCells, p) {
							in := bytes.IndexByte(orientations, orientation)
							orientations = append(orientations[:in], orientations[in+1:]...)
							failed = true
							break
						}
						pos = append(pos, p)
					}
				case 'E':
					for i := 0; i < sh.length; i++ {
						p := fmt.Sprintf("%v%v", string(cols[startCol+i]), startRow)
						if contains(s.allCells, p) {
							in := bytes.IndexByte(orientations, orientation)
							orientations = append(orientations[:in], orientations[in+1:]...)
							failed = true
							break
						}
						pos = append(pos, p)
					}
				case 'W':
					for i := 0; i < sh.length; i++ {
						p := fmt.Sprintf("%v%v", string(cols[startCol-i]), startRow)
						if contains(s.allCells, p) {
							in := bytes.IndexByte(orientations, orientation)
							orientations = append(orientations[:in], orientations[in+1:]...)
							failed = true
							break
						}
						pos = append(pos, p)
					}
				}

				if !failed {
					sh.position = append(sh.position, pos...)
					s.allCells = append(s.allCells, sh.position...)
					retry = false
				}
				if len(orientations) < 1 {
					// Setting failed = false will allow us to
					// jump out of current loop and retry
					failed = false
				}
			}
		}
	}
}

func (s *ships) getShip(class string) (*ship, error) {
	if ship := s.ships[class]; ship != nil {
		return ship, nil
	}
	return nil, fmt.Errorf("no ship of class '%s' found", class)
}

func (s *ships) getPositionByClass(class string) ([]string, error) {
	if ship, ok := s.ships[class]; ok {
		return ship.position, nil
	}
	return []string{}, fmt.Errorf("class '%s' does not exist", class)
}

func (s *ships) checkCell(cell string) (bool, string) {
	for k, s := range s.ships {
		for _, p := range s.position {
			if cell == p {
				return true, k
			}
		}
	}
	return false, ""
}

func (s *ships) registerHit(cell string, class string) {
	s.ships[class].health--
	if s.ships[class].health == 0 {
		s.remaining--
	}
}

func contains(s []string, v string) bool {
	for _, p := range s {
		if p == v {
			return true
		}
	}
	return false
}
