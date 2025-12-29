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
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := strings.ToLower(scanner.Text())
		cleaned := cleanInput(input)

		if supportedCommands[cleaned[0]].callback != nil {
			if err := supportedCommands[cleaned[0]].callback(cfg); err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
		}
	}
}
