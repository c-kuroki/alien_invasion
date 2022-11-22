package model

import (
	"fmt"
	"math/rand"
	"time"
)

type Alien struct {
	ID   int
	Name string
	City int
}

func NewAlien(ID, cityID int) *Alien {
	return &Alien{
		ID:   ID,
		Name: randomAlienName(ID),
		City: cityID,
	}
}

func randomAlienName(ID int) string {
	const letters = "zaxorukigmle"

	randomizer := make([]byte, 5)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range randomizer {
		randomizer[i] = letters[r.Intn(len(letters))]
	}
	return fmt.Sprintf("%s%d", string(randomizer), ID)
}
