package app

import (
	"math/rand"
	"time"
)

func getRandomInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
