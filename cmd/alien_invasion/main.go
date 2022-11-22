package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/c-kuroki/alien_invasion/pkg/adapters/renderer"
	"github.com/c-kuroki/alien_invasion/pkg/adapters/world"
	"github.com/c-kuroki/alien_invasion/pkg/app"
	"github.com/c-kuroki/alien_invasion/pkg/model"
	"github.com/c-kuroki/alien_invasion/pkg/ports/http"
)

func usage() {
	fmt.Println(`Usage: alien_invasion [OPTIONS] <num aliens>

OPTIONS
-------`)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	filename := flag.String("f", "./examples/big.map", "map filename")
	tickInterval := flag.Int("t", 1000, "tick interval")
	maxMoves := flag.Int("m", 10000, "max number of moves")
	httpServiceAddress := flag.String("a", ":8080", "http service address (-1 to disable http service)")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		usage()
	}

	numAliens, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("invalid number of aliens")
		usage()
	}
	if *tickInterval < 1 || *maxMoves < 1 || numAliens < 1 {
		fmt.Println("invalid parameters : tick interval, max moves and num aliens should be greater than 0")
		usage()
	}
	cfg := &model.Config{
		MapFilename:  *filename,
		TickInterval: *tickInterval,
		MaxMoves:     *maxMoves,
		NumAliens:    numAliens,
	}

	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	state := world.NewInMemoryState(cfg.MapFilename)
	rnd := renderer.NewSVGRenderer()

	invasion := app.NewAlienInvasionApp(cfg, state, rnd, logger.Sugar())
	if *httpServiceAddress != "-1" {
		srv := http.NewHTTPService(invasion, *httpServiceAddress)
		go srv.Start()
	}
	invasion.Start()
}
