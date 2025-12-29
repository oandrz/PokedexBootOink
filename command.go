package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	nextUrl string
	prevUrl string
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

	result, err := getPokemonMapLocation(commandConfig.nextUrl, commandConfig)
	if err != nil {
		return err
	}

	for _, location := range result {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(commandConfig *config) error {
	if commandConfig == nil {
		return fmt.Errorf("command config cannot be nil")
	}

	result, err := getPokemonMapLocation(commandConfig.prevUrl, commandConfig)
	if err != nil {
		return err
	}

	for _, location := range result {
		fmt.Println(location.Name)
	}

	return nil
}
