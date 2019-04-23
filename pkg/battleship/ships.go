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
}

const (
	cClass = "Canberra"
	hClass = "Hobart"
	lClass = "Leeuwin"
	aClass = "Armidale"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (ss *ships) setPositions() {
	var taken []string
	ss.remaining = 4
	ss.ships = map[string]*ship{
		cClass: &ship{length: 5, health: 5},
		hClass: &ship{length: 4, health: 4},
		lClass: &ship{length: 3, health: 3},
		aClass: &ship{length: 2, health: 2},
	}

	for _, s := range ss.ships {
		var retry bool
		for cont := true; cont; cont = retry {
			retry = true
			var orientations []byte
			cols := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'}
			startRow := rand.Intn(9) + 1
			startCol := rand.Intn(9)

			if startRow-s.length >= 0 {
				orientations = append(orientations, 'N')
			}
			if startRow+s.length <= 10 {
				orientations = append(orientations, 'S')
			}
			if (startCol+s.length)+1 <= 10 {
				orientations = append(orientations, 'E')
			}
			if (startCol-s.length)+1 >= 0 {
				orientations = append(orientations, 'W')
			}

			var failed bool
			for _cont := true; _cont; _cont = failed {
				orientation := orientations[rand.Intn(len(orientations))]
				failed = false
				pos := []string{}

				switch orientation {
				case 'N':
					for i := 0; i < s.length; i++ {
						p := fmt.Sprintf("%v%v", string(cols[startCol]), startRow-i)
						if contains(taken, p) {
							in := bytes.IndexByte(orientations, orientation)
							orientations = append(orientations[:in], orientations[in+1:]...)
							failed = true
							break
						}
						pos = append(pos, p)
					}
				case 'S':
					for i := 0; i < s.length; i++ {
						p := fmt.Sprintf("%v%v", string(cols[startCol]), startRow+i)
						if contains(taken, p) {
							in := bytes.IndexByte(orientations, orientation)
							orientations = append(orientations[:in], orientations[in+1:]...)
							failed = true
							break
						}
						pos = append(pos, p)
					}
				case 'E':
					for i := 0; i < s.length; i++ {
						p := fmt.Sprintf("%v%v", string(cols[startCol+i]), startRow)
						if contains(taken, p) {
							in := bytes.IndexByte(orientations, orientation)
							orientations = append(orientations[:in], orientations[in+1:]...)
							failed = true
							break
						}
						pos = append(pos, p)
					}
				case 'W':
					for i := 0; i < s.length; i++ {
						p := fmt.Sprintf("%v%v", string(cols[startCol-i]), startRow)
						if contains(taken, p) {
							in := bytes.IndexByte(orientations, orientation)
							orientations = append(orientations[:in], orientations[in+1:]...)
							failed = true
							break
						}
						pos = append(pos, p)
					}
				}

				if !failed {
					s.position = append(s.position, pos...)
					taken = append(taken, s.position...)
					retry = false
				}
				if len(orientations) < 1 {
					failed = false
				}
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

func contains(s []string, v string) bool {
	for _, p := range s {
		if p == v {
			return true
		}
	}
	return false
}
