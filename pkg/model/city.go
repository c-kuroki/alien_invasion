package model

import (
	"encoding/json"
	"fmt"
)

const (
	North = "north"
	East  = "east"
	South = "south"
	West  = "west"
)

type City struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	North string `json:"north"`
	East  string `json:"east"`
	South string `json:"south"`
	West  string `json:"west"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

func (c *City) String() string {
	b, _ := json.Marshal(c)
	return fmt.Sprintf("%s\n", string(b))
}

var validCard map[string]bool = map[string]bool{
	North: true,
	East:  true,
	South: true,
	West:  true,
}

func IsValidCard(card string) bool {
	return validCard[card]
}
