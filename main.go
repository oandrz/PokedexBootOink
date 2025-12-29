package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const pokemonLocationUrl = "https://pokeapi.co/api/v2/location-area/"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commandConfig := config{
		nextUrl: pokemonLocationUrl,
		prevUrl: "",
	}
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
			if err := supportedCommands[cleaned[0]].callback(&commandConfig); err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
		}
	}
}
