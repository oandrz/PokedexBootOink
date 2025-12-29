package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex_go/internal/api"
	"strings"
)

type config struct {
	pokemonClient api.Client
	nextUrl       *string
	prevUrl       *string
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	supportedCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a map of the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays a map of the Pokemon world in reverse order",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific city in the Pokemon world",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a specific Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a specific Pokemon",
			callback:    commandInspect,
		},
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := strings.ToLower(scanner.Text())
		cleaned := cleanInput(input)

		if cleaned[0] == "" {
			continue
		}

		command := cleaned[0]
		args := []string{}
		if len(cleaned) > 1 {
			args = cleaned[1:]
		}

		if supportedCommands[command].callback != nil {
			if err := supportedCommands[command].callback(cfg, args...); err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
		}
	}
}
