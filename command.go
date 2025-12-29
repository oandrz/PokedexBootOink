package main

import (
	"fmt"
	"os"
)

const pokemonLocationUrl = "https://pokeapi.co/api/v2/location-area/"

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandExit(commandConfig *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commandConfig *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")

	return nil
}

func commandMap(commandConfig *config) error {
	if commandConfig == nil {
		return fmt.Errorf("command config cannot be nil")
	}

	var url string
	if commandConfig.nextUrl == nil {
		url = pokemonLocationUrl
	} else {
		url = *commandConfig.nextUrl
	}

	result, err := commandConfig.pokemonClient.GetPokemonMapLocation(url)
	if err != nil {
		return err
	}

	commandConfig.nextUrl = result.Next
	commandConfig.prevUrl = result.Previous

	for _, location := range result.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(commandConfig *config) error {
	if commandConfig == nil {
		return fmt.Errorf("command config cannot be nil")
	}

	if commandConfig.prevUrl == nil {
		fmt.Println("you're on the first page, no previous locations")
		return nil
	}

	result, err := commandConfig.pokemonClient.GetPokemonMapLocation(*commandConfig.prevUrl)
	if err != nil {
		return err
	}

	for _, location := range result.Results {
		fmt.Println(location.Name)
	}

	commandConfig.nextUrl = result.Next
	commandConfig.prevUrl = result.Previous

	return nil
}
