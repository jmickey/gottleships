package battleship

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionsAreAdjacent(t *testing.T) {
	sh := &ships{}
	sh.setPositions()

	pos, err := sh.getPositionByClass("Hobart")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	f, s := pos[0][0], pos[1][0]
	if f == s {
		f, s = pos[0][1], pos[1][1]
	}

	assert.Truef(t,
		(f-1 == s || f+1 == s),
		"row/col %s is not adjacent to %s", string(s), string(f))
}

func TestShipPositionsValid(t *testing.T) {
	s := &ships{}
	s.setPositions()
	assert.True(t, len(s.allCells) == 14, "there should be exactly 14 positions")
	for _, p := range s.allCells {
		match, _ := regexp.MatchString("^[A-Z][1-9]$", p)
		assert.Truef(t, match, "position %s doesn't match regex", p)
	}
}

func TestCheckCell(t *testing.T) {
	s := &ships{}
	s.setPositions()
	hit, class := s.checkCell(s.ships[hClass].position[1])
	assert.True(t, hit, "checkCell should return true")
	assert.Equal(t, hClass, class, "cells should be equal")
}

func TestRegisterHit(t *testing.T) {
	s := &ships{}
	s.setPositions()
	s.registerHit(s.ships[lClass].position[2], lClass)
	assert.Equal(t, 2, s.ships[lClass].health, "registerHit should return true")
}

func TestGetShip(t *testing.T) {
	s := &ships{}
	s.setPositions()
	_, err := s.getShip("Hobart")
	assert.Nil(t, err, "error was not nil")
	_, err = s.getShip("NonExistant")
	assert.Error(t, err, "s.getShip(\"NonExistant\") should return an error")
}
