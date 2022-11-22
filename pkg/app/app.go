package app

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/c-kuroki/alien_invasion/pkg/adapters/renderer"
	"github.com/c-kuroki/alien_invasion/pkg/adapters/world"
	"github.com/c-kuroki/alien_invasion/pkg/logger"
	"github.com/c-kuroki/alien_invasion/pkg/model"
)

type AlienInvasionApp struct {
	state    world.Adapter
	renderer renderer.Adapter
	cfg      *model.Config
	log      logger.Logger
}

func NewAlienInvasionApp(cfg *model.Config, state world.Adapter, renderer renderer.Adapter, log logger.Logger) *AlienInvasionApp {
	return &AlienInvasionApp{
		cfg:      cfg,
		state:    state,
		renderer: renderer,
		log:      log,
	}
}

func (app *AlienInvasionApp) RenderMap(ctx context.Context, w io.Writer) error {
	cities := app.state.GetAllCities()
	aliens := app.state.GetAllAliensByCity()
	return app.renderer.Render(ctx, cities, aliens, w)
}

func (app *AlienInvasionApp) Start() {
	app.log.Infow("starting invasion app")
	// load map
	err := app.state.Load()
	if err != nil {
		app.log.Errorw("error loading map", "error", err.Error())
		return
	}
	cities := app.state.GetAllCities()
	// add aliens
	max := len(cities) - 1
	for i := 0; i < app.cfg.NumAliens; i++ {
		cityID := getRandomInRange(0, max)
		alien := model.NewAlien(i, cityID)
		err = app.state.AddAlien(alien)
		if err != nil {
			app.log.Warnw("adding alien", "error", err.Error())
			continue
		}
		app.log.Infow("added alien", "id", fmt.Sprintf("%d", alien.ID), "name", alien.Name)
	}
	app.MainLoop()
}

// main loop
func (app *AlienInvasionApp) MainLoop() {
	var moves int
	tickerChan := time.NewTicker(time.Duration(int64(time.Millisecond) * int64(app.cfg.TickInterval))).C
	for {
		<-tickerChan
		moves++
		err := app.makeMove()
		if err != nil {
			app.log.Warnw("making move", "move #", fmt.Sprint(moves), "error", err.Error())
		}
		aliens := app.state.GetAliens()
		if moves > app.cfg.MaxMoves || len(aliens) == 0 {
			// save final map and exit
			finalMapFile := fmt.Sprintf("%s.map", time.Now().Format(time.RFC3339))
			err = app.state.Save(finalMapFile)
			if err != nil {
				app.log.Errorw("writing final map", "filename", finalMapFile, "error", err.Error())
			}
			return
		}
	}
}

func (app *AlienInvasionApp) makeMove() error {
	// move aliens
	aliens := app.state.GetAliens()
	for _, alien := range aliens {
		exits, err := app.state.GetExits(alien.City)
		if err != nil {
			app.log.Warnw("getting exits ", "cityID", alien.City, "error", err.Error())
			continue
		}
		numExits := len(exits)
		if numExits > 0 {
			// get next move
			moveIndex := getRandomInRange(0, numExits)
			// if random index is exactly numExits will not move this turn
			if moveIndex != numExits {
				err := app.state.MoveAlien(alien.ID, exits[moveIndex])
				if err != nil {
					app.log.Warnw("moving alien", "error", err.Error())
					continue
				}
			}
		}
	}

	// check fights
	aliensByCity := app.state.GetAllAliensByCity()
	for cityID, aliensMap := range aliensByCity {
		if len(aliensMap) > 1 {
			fightCity, _ := app.state.GetCityByID(cityID)
			app.log.Infow("Fight !!", "city", fightCity.Name, "aliens", len(aliensMap))
			err := app.state.RemoveCity(cityID)
			if err != nil {
				app.log.Warnw("removing city", "error", err.Error())
			}
		}
	}
	return nil
}
