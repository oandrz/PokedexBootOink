package main

import (
	"encoding/json"
	"net/http"
)

type PokemonMapResponse struct {
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []MapLocation `json:"results"`
}

type MapLocation struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func fetchRemotePokemonMapLocation(url string, commandConfig *config) ([]MapLocation, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pokemonMap PokemonMapResponse
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&pokemonMap); err != nil {
		return nil, err
	}

	commandConfig.prevUrl = commandConfig.nextUrl
	commandConfig.nextUrl = pokemonMap.Next

	return pokemonMap.Results, nil
}
