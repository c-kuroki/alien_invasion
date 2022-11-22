package world

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/c-kuroki/alien_invasion/pkg/model"
)

// check that interface is implemented
var _ Adapter = (*InMemoryState)(nil)

var (
	notFoundErr         = errors.New("not found")
	invalidCityErr      = errors.New("invalid city")
	invalidMapErr       = errors.New("invalid map")
	invalidDirectionErr = errors.New("invalid direction")
	duplicatedErr       = errors.New("duplicated city")
	dupConnErr          = errors.New("duplicated connection")
)

// InMemoryState loads and save worlds from/to files storing state in memory
type InMemoryState struct {
	filename     string
	nextID       int
	citiesByName map[string]*model.City
	citiesByID   map[int]*model.City
	aliensByCity map[int]map[int]*model.Alien
	aliensByID   map[int]*model.Alien
	mapHeight    int
	mapWidth     int
}

func NewInMemoryState(filename string) *InMemoryState {
	return &InMemoryState{
		filename:     filename,
		citiesByName: make(map[string]*model.City),
		citiesByID:   make(map[int]*model.City),
		aliensByCity: make(map[int]map[int]*model.Alien),
		aliensByID:   make(map[int]*model.Alien),
	}
}

// Load loads a world map from a file (max line size of 64K)
func (st *InMemoryState) Load() error {
	file, err := os.Open(st.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// scan line by line
	var lineNum uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		fields, err := parseLine(lineNum, scanner.Text())
		if err != nil {
			return err
		}
		if err := st.AddCity(fields...); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	err = st.validateCities()
	if err != nil {
		return err
	}
	err = st.setCoordinates()
	if err != nil {
		return err
	}
	return nil
}

func (st *InMemoryState) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	cities := st.GetAllCities()
	// sort result
	sort.Slice(cities, func(i, j int) bool {
		return cities[i].ID < cities[j].ID
	})
	for _, city := range cities {
		_, err = file.WriteString(city.Name)
		if err != nil {
			return err
		}
		if city.North != "" {
			_, err = file.WriteString(fmt.Sprintf(" north=%s", city.North))
			if err != nil {
				return err
			}
		}
		if city.East != "" {
			_, err = file.WriteString(fmt.Sprintf(" east=%s", city.East))
			if err != nil {
				return err
			}
		}
		if city.South != "" {
			_, err = file.WriteString(fmt.Sprintf(" south=%s", city.South))
			if err != nil {
				return err
			}
		}
		if city.West != "" {
			_, err = file.WriteString(fmt.Sprintf(" west=%s", city.West))
			if err != nil {
				return err
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}

	}
	_ = file.Sync()
	return nil
}

func (st *InMemoryState) GetNumCities() int {
	return len(st.citiesByID)
}

func (st *InMemoryState) GetWidth() int {
	return st.mapWidth
}

func (st *InMemoryState) GetHeight() int {
	return st.mapHeight
}

func (st *InMemoryState) GetAllCities() []*model.City {
	var cities []*model.City
	for _, city := range st.citiesByID {
		cities = append(cities, city)
	}
	return cities
}

func (st *InMemoryState) GetAliens() map[int]*model.Alien {
	return st.aliensByID
}

func (st *InMemoryState) GetAllAliensByCity() map[int]map[int]*model.Alien {
	return st.aliensByCity
}

func (st *InMemoryState) GetCityByID(cityID int) (*model.City, error) {
	city, ok := st.citiesByID[cityID]
	if !ok {
		return city, notFoundErr
	}
	return city, nil
}

func (st *InMemoryState) GetCityByName(name string) (*model.City, error) {
	city, ok := st.citiesByName[name]
	if !ok {
		return city, notFoundErr
	}
	return city, nil
}

func (st *InMemoryState) GetAlienByID(alienID int) (*model.Alien, error) {
	alien, ok := st.aliensByID[alienID]
	if !ok {
		return alien, notFoundErr
	}
	return alien, nil
}

func (st *InMemoryState) GetAliensByCity(cityID int) (map[int]*model.Alien, error) {
	alien, ok := st.aliensByCity[cityID]
	if !ok {
		return alien, notFoundErr
	}
	return alien, nil
}

func (st *InMemoryState) AddAlien(alien *model.Alien) error {
	st.aliensByID[alien.ID] = alien
	return st.addAlienToCity(alien, alien.City)
}

func (st *InMemoryState) addAlienToCity(alien *model.Alien, cityID int) error {
	if aliens, ok := st.aliensByCity[alien.City]; ok {
		aliens[alien.ID] = alien
	} else {
		alienMap := make(map[int]*model.Alien, 1)
		alienMap[alien.ID] = alien
		st.aliensByCity[alien.City] = alienMap
	}
	st.aliensByID[alien.ID] = alien
	return nil
}

func (st *InMemoryState) MoveAlien(alienID, cityID int) error {
	alien, ok := st.aliensByID[alienID]
	if !ok {
		return notFoundErr
	}
	if source, ok := st.aliensByCity[alien.City]; ok {
		delete(source, alienID)
	}
	alien.City = cityID
	return st.addAlienToCity(alien, cityID)
}

func (st *InMemoryState) AddCity(args ...string) error {
	if len(args) != 5 {
		return invalidCityErr
	}
	name := args[0]
	north := args[1]
	east := args[2]
	south := args[3]
	west := args[4]

	// check if city already exists
	_, err := st.GetCityByName(name)
	if err == nil {
		return duplicatedErr
	}
	city := &model.City{
		ID:    st.nextID,
		Name:  name,
		North: north,
		East:  east,
		South: south,
		West:  west,
	}
	st.nextID++
	st.citiesByName[name] = city
	st.citiesByID[city.ID] = city
	return nil
}

func (st *InMemoryState) validateCities() error {
	cities := st.GetAllCities()
	for _, city := range cities {
		// check north
		if city.North != "" {
			n, err := st.GetCityByName(city.North)
			if err != nil {
				return fmt.Errorf("%s city : error %s", city.North, err.Error())
			}
			if n.South != city.Name {
				return fmt.Errorf("invalid connection: %s north connection (%s) doesnt match %s south connection (%s)", city.Name, city.North, n.Name, n.South)
			}
		}
		// check east
		if city.East != "" {
			e, err := st.GetCityByName(city.East)
			if err != nil {
				return fmt.Errorf("%s city : error %s", city.East, err.Error())
			}
			if e.West != city.Name {
				return fmt.Errorf("invalid connection: %s east connection (%s) doesnt match %s west connection (%s)", city.Name, city.East, e.Name, e.West)
			}
		}
		// check south
		if city.South != "" {
			s, err := st.GetCityByName(city.South)
			if err != nil {
				return fmt.Errorf("%s city : error %s", city.South, err.Error())
			}
			if s.North != city.Name {
				return fmt.Errorf("invalid connection: %s south connection (%s) doesnt match %s north connection (%s)", city.Name, city.South, s.Name, s.North)
			}
		}
		// check west
		if city.West != "" {
			w, err := st.GetCityByName(city.West)
			if err != nil {
				return fmt.Errorf("%s city : error %s", city.West, err.Error())
			}
			if w.East != city.Name {
				return fmt.Errorf("invalid connection: %s west connection (%s) doesnt match %s east connection (%s)", city.Name, city.West, w.Name, w.East)
			}
		}
	}
	return nil
}

func (st *InMemoryState) setCoordinates() error {
	// first traverse world and set relative coordinates , using (0,0) from first city
	visited := make(map[int]bool)
	var x, y int
	y = 1 // initialize y coord to 1 because first call is direction north
	numCities := st.GetNumCities()
	err := st.Traverse(numCities, visited, 0, 0, func(c *model.City, direction, action int) error {
		// if the action is just rotation and no move, not modify position
		if action == 4 {
			return nil
		}
		// tracks coordinates while traversing
		switch direction {
		case 0:
			y = y - 1
		case 1:
			x = x + 1
		case 2:
			y = y + 1
		case 3:
			x = x - 1
		}
		c.X = x
		c.Y = y
		return nil
	}, 0)
	if err != nil {
		return err
	}

	cities := st.GetAllCities()
	// if not all cities were visited exists one ore more unreachable (isolated) cities on map, that is invalid
	if len(visited) != len(cities) {
		return invalidMapErr
	}
	// reindex with absolute values
	// get min and max X and Y
	var minX, maxX, minY, maxY int
	for _, city := range cities {
		if city.X < minX {
			minX = city.X
		}
		if city.X > maxX {
			maxX = city.X
		}
		if city.Y < minY {
			minY = city.Y
		}
		if city.Y > maxY {
			maxY = city.Y
		}
	}
	// set height and width
	st.mapWidth = maxX - minX + 1
	st.mapHeight = maxY - minY + 1
	// now replace with absolute values, and add to coordinate indexed map
	for _, city := range cities {
		city.X = city.X + (minX * -1)
		city.Y = city.Y + (minY * -1)
		//	st.citiesByCoord[model.NewCoord(city.X, city.Y)] = city
	}
	return nil
}

// traverse recursively the world calling func passed as parameter on each new city visited
// uses right-hand wallet follower algorithm
// params:
// visited = keeps already visited cities ids
// direction = next direction to follow (facing), coded as 0:N, 1:E, 2:S, 3:W
// current = current city ID
// onEach = optional function to be called when first visiting a city
// action = previous action type id
func (st *InMemoryState) Traverse(numCities int, visited map[int]bool, direction int, current int, onEach func(*model.City, int, int) error, action int) error {
	// base case is reached when we fall again in starting state ( on first indexed city facing north)
	if len(visited) == numCities || (len(visited) > 0 && current == 0 && direction == 0) {
		return nil
	}
	// get current city
	city, err := st.GetCityByID(current)
	if err != nil {
		return err
	}
	// append to visited list
	visited[current] = true
	// call on each function if defined
	if onEach != nil {
		err := onEach(city, direction, action)
		if err != nil {
			return err
		}
	}
	right, err := rightConnectionExists(city, direction)
	if err != nil {
		return err
	}
	front, err := frontConnectionExists(city, direction)
	if err != nil {
		return err
	}
	// action 1 -  no wall at right - turn 90 clockwise and move forward
	if right {
		direction = rotate(direction, true)
		nextCityID, err := st.forward(city, direction)
		if err != nil {
			return err
		}
		err = st.Traverse(numCities, visited, direction, nextCityID, onEach, 1)
		if err != nil {
			return err
		}
		return nil
	}
	// action 2 - wall at right, no wall at front - move forward
	if !right && front {
		nextCityID, err := st.forward(city, direction)
		if err != nil {
			return err
		}
		err = st.Traverse(numCities, visited, direction, nextCityID, onEach, 2)
		if err != nil {
			return err
		}
		return nil
	}
	// action 3 - no wall at right - turn 90 clockwise and move forward
	if right {
		direction = rotate(direction, true)
		nextCityID, err := st.forward(city, direction)
		if err != nil {
			return err
		}
		err = st.Traverse(numCities, visited, direction, nextCityID, onEach, 3)
		if err != nil {
			return err
		}
		return nil
	}
	// action 4 - wall at right and wall at front - turn 90 counter-clockwise
	if !right && !front {
		direction = rotate(direction, false)
		err = st.Traverse(numCities, visited, direction, current, onEach, 4)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// forward moves forward in passed direction, returning new city id
func (st *InMemoryState) forward(city *model.City, direction int) (int, error) {
	var name string
	switch direction {
	case 0: // north
		name = city.North
	case 1: // east
		name = city.East
	case 2: // south
		name = city.South
	case 3: // west
		name = city.West
	}
	newCity, err := st.GetCityByName(name)
	if err != nil {
		return 0, err
	}
	return newCity.ID, nil
}

// GetExits returns an array of cities IDs connected with the input city
func (st *InMemoryState) GetExits(cityID int) ([]int, error) {
	var exitNames []string
	// retrieve city
	city, err := st.GetCityByID(cityID)
	if err != nil {
		return nil, err
	}
	// check connections with other cities
	if city.North != "" {
		exitNames = append(exitNames, city.North)
	}
	if city.East != "" {
		exitNames = append(exitNames, city.East)
	}
	if city.South != "" {
		exitNames = append(exitNames, city.South)
	}
	if city.West != "" {
		exitNames = append(exitNames, city.West)
	}

	// build exits ids
	exits := make([]int, len(exitNames))
	for ix := range exitNames {
		exit, err := st.GetCityByName(exitNames[ix])
		if err != nil {
			return nil, err
		}
		exits[ix] = exit.ID
	}
	return exits, nil
}

func (st *InMemoryState) RemoveCity(cityID int) error {
	// retrieve city
	city, err := st.GetCityByID(cityID)
	if err != nil {
		return err
	}
	// remove connections from other cities
	if city.North != "" {
		err = st.removeConnection(city.North, model.South)
		if err != nil {
			return err
		}
	}
	if city.East != "" {
		err = st.removeConnection(city.East, model.West)
		if err != nil {
			return err
		}
	}
	if city.South != "" {
		err = st.removeConnection(city.South, model.North)
		if err != nil {
			return err
		}
	}
	if city.West != "" {
		err = st.removeConnection(city.West, model.East)
		if err != nil {
			return err
		}
	}
	// remove aliens at the city
	aliens, err := st.GetAliensByCity(cityID)
	if err == nil {
		for alienID := range aliens {
			delete(st.aliensByID, alienID)
		}
		delete(st.aliensByCity, cityID)
	}
	// remove city
	delete(st.citiesByName, city.Name)
	delete(st.citiesByID, cityID)
	return nil
}

func (st *InMemoryState) removeConnection(cityName, connection string) error {
	// retrieve city
	city, err := st.GetCityByName(cityName)
	if err != nil {
		return err
	}
	switch connection {
	case model.North:
		city.North = ""
	case model.East:
		city.East = ""
	case model.South:
		city.South = ""
	case model.West:
		city.West = ""
	}
	return nil
}
