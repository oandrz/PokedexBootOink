package main

import (
	"fmt"
	"os"
)

const pokemonLocationUrl = "https://pokeapi.co/api/v2/location-area/"

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func commandExit(commandConfig *config, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commandConfig *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")

	return nil
}

func commandMap(commandConfig *config, args ...string) error {
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

func commandMapB(commandConfig *config, args ...string) error {
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

func commandExplore(commandConfig *config, args ...string) error {
	cityToExplore := args[0]
	url := pokemonLocationUrl + cityToExplore
	fmt.Println("Exploring " + cityToExplore + "...")
	result, err := commandConfig.pokemonClient.GetPokemonMapLocation(url)
	if err != nil {
		return err
	}

	for _, pokemonEncounter := range result.PokemonsEncounter {
		fmt.Println(pokemonEncounter.PokemonEncounter.Name)
	}

	return nil
}

func commandCatch(commandConfig *config, args ...string) error {
	pokemonName := args[0]
	fmt.Println("Throwing a Pokeball at " + pokemonName + "...")
	isSuccess, err := commandConfig.pokemonClient.CatchPokemon(pokemonName)
	if err != nil {
		return err
	}

	if isSuccess {
		fmt.Println("Success Caught " + pokemonName + "...")
		return nil
	}

	fmt.Println("Failed to catch " + pokemonName + "...")
	return nil
}

func commandInspect(commandConfig *config, args ...string) error {
	pokemonName := args[0]
	pokedex := commandConfig.pokemonClient.Pokedex

	pokemon, ok := pokedex[pokemonName]
	if !ok {
		fmt.Println("Pokemon not found in pokedex")
		return nil
	}

	pokemon.Print()
	return nil
}
