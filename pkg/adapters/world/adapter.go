package world

import (
	"github.com/c-kuroki/alien_invasion/pkg/model"
)

//go:generate mockery --name Adapter

// World Adapter interface to manage world state
type Adapter interface {
	GetNumCities() int
	GetWidth() int
	GetHeight() int
	GetAllCities() []*model.City
	GetAliens() map[int]*model.Alien
	GetCityByID(ID int) (*model.City, error)
	GetCityByName(name string) (*model.City, error)
	GetExits(cityID int) ([]int, error)
	GetAlienByID(ID int) (*model.Alien, error)
	GetAliensByCity(cityID int) (map[int]*model.Alien, error)
	GetAllAliensByCity() map[int]map[int]*model.Alien
	AddAlien(*model.Alien) error
	MoveAlien(alienID, toCityID int) error
	AddCity(args ...string) error
	RemoveCity(CityID int) error
	Load() error
	Save(string) error
}
