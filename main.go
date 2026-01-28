package main

import (
	_ "embed"
	"os"
	"snakehem/game"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05.999Z07:00"})
	game.Run()
}
