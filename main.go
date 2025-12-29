package main

import (
	"pokedex_go/internal/api"
	"time"
)

func main() {
	client := api.NewClient(
		5*time.Second,
		10*time.Minute,
	)
	commandConfig := config{
		pokemonClient: client,
	}

	startRepl(&commandConfig)
}
