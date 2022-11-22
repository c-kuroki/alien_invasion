package model

import "fmt"

type Coord string

func NewCoord(x, y int) *Coord {
	coord := Coord(fmt.Sprintf("(%d,%d)", x, y))
	return &coord
}

func (c Coord) String() string {
	return string(c)
}
