package world

import (
	"fmt"
	"os"
	"testing"

	"github.com/c-kuroki/alien_invasion/pkg/model"
	"github.com/stretchr/testify/suite"
)

type testCase struct {
	name           string
	fileContent    string
	expectedCities []*model.City
	expectedError  string
}

type InMemoryStateTestSuite struct {
	suite.Suite
	st        *InMemoryState
	testCases []testCase
}

const (
	exampleMapFile = "../../../examples/world.map"
	testMapFile    = "testSave.map"
)

func (suite *InMemoryStateTestSuite) SetupTest() {
	suite.st = NewInMemoryState(exampleMapFile)
	suite.testCases = getTestCases()
}

func (suite *InMemoryStateTestSuite) TearDownSuite() {
	// remove generated files
	for ix := range suite.testCases {
		fileName := fmt.Sprintf("map_%d.map", ix)
		_ = os.Remove(fileName)
	}
	_ = os.Remove(testMapFile)
}

func (suite *InMemoryStateTestSuite) TestLoad() {
	for ix, tc := range suite.testCases {
		// run tests in parallel
		suite.T().Run(tc.name, func(t *testing.T) {

			fileName := fmt.Sprintf("map_%d.map", ix)
			// create file
			f, err := os.Create(fileName)
			suite.Require().NoError(err)
			defer f.Close()

			// write test content
			_, err = f.Write([]byte(tc.fileContent))
			suite.Require().NoError(err)

			// test map parse and load
			state := NewInMemoryState(fileName)
			err = state.Load()
			if tc.expectedError != "" {
				suite.Require().Error(err)
				suite.Assert().Contains(err.Error(), tc.expectedError)
			} else {
				suite.Require().NoError(err)
				cities := state.GetAllCities()
				suite.Assert().Equal(len(tc.expectedCities), len(cities))
				for _, expectedCity := range tc.expectedCities {
					city, err := state.GetCityByID(expectedCity.ID)
					suite.Require().NoError(err)
					suite.Assert().Equal(expectedCity, city)
				}
			}
		})
	}
}

func (suite *InMemoryStateTestSuite) TestSave() {
	err := suite.st.Load()
	suite.Require().NoError(err)

	err = suite.st.Save(testMapFile)
	suite.Require().NoError(err)

	expected, err := os.ReadFile(exampleMapFile)
	suite.Require().NoError(err)
	result, err := os.ReadFile(testMapFile)
	suite.Require().NoError(err)
	suite.Assert().Equal(expected, result)
}

func (suite *InMemoryStateTestSuite) TestRemove() {
	state := NewInMemoryState(exampleMapFile)
	err := state.Load()
	suite.Require().NoError(err)

	// remove first city
	city, err := state.GetCityByName("Foo")
	suite.Require().NoError(err)
	err = state.RemoveCity(city.ID)
	suite.Require().NoError(err)

	// check that city was removed
	_, err = state.GetCityByName("Foo")
	suite.Require().Error(err)
	suite.Assert().Equal(err, notFoundErr)

	// check that city was unconnected
	city, err = state.GetCityByName("Qu-ux")
	suite.Require().NoError(err)
	suite.Assert().Equal(city.North, "")
}

func (suite *InMemoryStateTestSuite) TestAliens() {
	testAliens := getAliens()
	// add test aliens
	for _, alien := range testAliens {
		err := suite.st.AddAlien(alien)
		suite.Require().NoError(err)
	}
	// retrieve and check aliens
	aliens := suite.st.GetAliens()
	suite.Assert().Equal(3, len(aliens))
	for _, testAlien := range testAliens {
		alien, err := suite.st.GetAlienByID(testAlien.ID)
		suite.Require().NoError(err)
		suite.Assert().Equal(testAlien, alien)
	}

	// check aliens per city
	aliensAtCity2, err := suite.st.GetAliensByCity(2)
	suite.Require().NoError(err)
	suite.Assert().Equal(2, len(aliensAtCity2))
	aliensAtCity0, err := suite.st.GetAliensByCity(0)
	suite.Require().NoError(err)
	suite.Assert().Equal(1, len(aliensAtCity0))
	_, err = suite.st.GetAliensByCity(1)
	suite.Require().Error(err)
	suite.Assert().Equal("not found", err.Error())

	// move an alien check if aliens per city was updated
	err = suite.st.MoveAlien(2, 2)
	suite.Require().NoError(err)
	aliensAtCity2, err = suite.st.GetAliensByCity(2)
	suite.Require().NoError(err)
	suite.Assert().Equal(3, len(aliensAtCity2))
	aliensAtCity0, err = suite.st.GetAliensByCity(0)
	suite.Require().NoError(err)
	suite.Assert().Equal(0, len(aliensAtCity0))

}

// TestInMemoryState is the entry point of this test suite
func TestInMemoryState(t *testing.T) {
	suite.Run(t, new(InMemoryStateTestSuite))
}

func getTestCases() []testCase {
	return []testCase{
		{
			name: "valid 5 cities",
			fileContent: `Foo north=Bar west=Baz south=Qu-ux 
Bar south=Foo west=Bee
Qu-ux north=Foo
Baz east=Foo
Bee east=Bar`,
			expectedCities: []*model.City{
				{ID: 0, Name: "Foo", North: "Bar", South: "Qu-ux", West: "Baz", X: 1, Y: 1},
				{ID: 1, Name: "Bar", South: "Foo", West: "Bee", X: 1, Y: 0},
				{ID: 2, Name: "Qu-ux", North: "Foo", X: 1, Y: 2},
				{ID: 3, Name: "Baz", East: "Foo", X: 0, Y: 1},
				{ID: 4, Name: "Bee", East: "Bar", X: 0, Y: 0},
			},
		},
		{
			name: "invalid connection",
			fileContent: `Foo south=Qu-ux west=Bar
Qu-ux north=Bar
Bar east=Foo`,
			expectedError: "invalid connection",
		},
		{
			name: "invalid duplicated connection",
			fileContent: `Foo south=Bar south=Qu-ux
Qu-ux north=Foo`,
			expectedError: "duplicated connection",
		},
		{
			name: "invalid loop",
			fileContent: `Foo south=Qu-ux 
Qu-ux north=Foo
Bar east=Bee
Bee west=Bar
`,
			expectedError: "invalid map",
		},
	}
}

func getAliens() []*model.Alien {
	return []*model.Alien{
		{
			ID:   1,
			Name: "Zork",
			City: 2,
		},
		{
			ID:   2,
			Name: "Mork",
			City: 0,
		},
		{
			ID:   3,
			Name: "Gork",
			City: 2,
		},
	}
}
