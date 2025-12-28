package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := strings.ToLower(scanner.Text())
		cleaned := cleanInput(input)

		if supportedCommands[cleaned[0]].callback != nil {
			if err := supportedCommands[cleaned[0]].callback(); err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
		}
	}
}
