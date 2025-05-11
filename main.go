package main

import (
	_ "embed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"snakehem/game"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	game.Run()
}
