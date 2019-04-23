package battleship

import (
	"bytes"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShipPositionsValid(t *testing.T) {
	sh := &ships{}
	sh.setPositions()

	for k, s := range sh.ships {
		for _, p := range s.position {
			match, _ := regexp.MatchString("^[A-Z][1-9]$", p)
			assert.Truef(t, match, "position %s in %s doesn't match regex", p, k)
		}
	}
}

func TestShipPlacement(t *testing.T) {
	sh := &ships{}
	sh.setPositions()
	assert.Equal(t, int16(4), sh.remaining, "should be 4 ships remaining")

	pos1 := sh.ships["Canberra"].position[0]
	pos2 := sh.ships["Canberra"].position[1]
	assert.NotEqual(t, pos1, pos2, "cells should not be the same")

	switch {
	case pos1[0] == pos2[0]:
		i, _ := strconv.Atoi(string(pos1[1]))
		j, _ := strconv.Atoi(string(pos2[1]))
		rows := []int{}
		if i > 0 {
			rows = append(rows, i-1)
		}
		if i < 8 {
			rows = append(rows, i+1)
		}
		assert.Contains(t, rows, j, "row should be +- 1")

	case pos1[1] == pos2[1]:
		col := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I'}
		index := bytes.IndexByte(col, pos1[0])
		cols := []byte{}
		if index > 0 {
			cols = append(cols, col[index-1])
		}
		if index < 8 {
			cols = append(cols, col[index+1])
		}
		assert.Contains(t, cols, col[bytes.IndexByte(col, pos2[0])], "should be adjacent column")
	}
}

func TestCheckCell(t *testing.T) {
	sh := &ships{}
	sh.setPositions()
	hit, class := sh.checkCell(sh.ships[hClass].position[1])
	assert.True(t, hit, "checkCell should return true")
	assert.Equal(t, hClass, class, "cells should be equal")
}

func TestRegisterHit(t *testing.T) {
	sh := &ships{}
	sh.setPositions()
	sh.registerHit(sh.ships[lClass].position[2], lClass)
	assert.Equal(t, 2, sh.ships[lClass].health, "registerHit should return true")
}
