package world

import (
	"fmt"
	"strings"

	"github.com/c-kuroki/alien_invasion/pkg/model"
)

// Parse a line and return 5 fields, city name, north, east , south and west
func parseLine(num uint64, line string) ([]string, error) {
	var name, north, east, south, west string
	fields := strings.Fields(line)
	done := make(map[string]bool)
	for ix, f := range fields {
		// first field is city name
		if ix == 0 {
			name = f
			continue
		}
		// rest of the fields should be a cardinal point
		card := strings.Split(f, "=")
		if len(card) != 2 {
			return nil, invalidMapErr
		}
		if !model.IsValidCard(card[0]) {
			return nil, fmt.Errorf("invalid card [%s] (should be north,east,south or east)", card[0])
		}
		switch card[0] {
		case model.North:
			if done[model.North] {
				return nil, dupConnErr
			}
			north = card[1]
			done[model.North] = true
		case model.East:
			if done[model.East] {
				return nil, dupConnErr
			}
			east = card[1]
			done[model.East] = true
		case model.South:
			if done[model.South] {
				return nil, dupConnErr
			}
			south = card[1]
			done[model.South] = true
		case model.West:
			if done[model.West] {
				return nil, dupConnErr
			}
			west = card[1]
			done[model.West] = true
		}
	}
	return []string{name, north, east, south, west}, nil
}

// rightConnectionExists returns true if there is connection to other city to the right side of current facing position
func rightConnectionExists(city *model.City, direction int) (bool, error) {
	switch direction {
	case 0: // north
		return city.East != "", nil
	case 1: // east
		return city.South != "", nil
	case 2: // south
		return city.West != "", nil
	case 3: // west
		return city.North != "", nil
	}
	return false, invalidDirectionErr
}

// frontConnectionExists returns true if there is connection to other city to the front side of current facing position
func frontConnectionExists(city *model.City, direction int) (bool, error) {
	switch direction {
	case 0: // north
		return city.North != "", nil
	case 1: // east
		return city.East != "", nil
	case 2: // south
		return city.South != "", nil
	case 3: // west
		return city.West != "", nil
	}
	return false, invalidDirectionErr
}

// rotate direction 90 degrees
func rotate(direction int, clockwise bool) int {
	// clockwise
	if clockwise {
		direction++
		if direction > 3 {
			return 0
		}
		return direction
	}
	// anti clockwise
	direction--
	if direction < 0 {
		return 3
	}
	return direction
}
