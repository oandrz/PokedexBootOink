package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedex_go/internal"
	"time"
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

var cache *internal.Cache

func getPokemonMapLocation(url string, commandConfig *config) ([]MapLocation, error) {
	if cache == nil {
		cache = internal.NewCache(10 * time.Second)
	}

	if cached, ok := cache.Get(url); ok {
		fmt.Println("Using cached data")
		pokemonMap, err := decodePokemonMap(cached)
		if err != nil {
			return nil, err
		}
		return pokemonMap.Results, nil
	}

	fmt.Println("Using Remote data")

	return fetchRemotePokemonMapLocation(url, commandConfig)
}

func fetchRemotePokemonMapLocation(url string, commandConfig *config) ([]MapLocation, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, data)

	pokemonMap, err := decodePokemonMap(data)
	if err != nil {
		return nil, err
	}

	commandConfig.prevUrl = commandConfig.nextUrl
	commandConfig.nextUrl = pokemonMap.Next

	return pokemonMap.Results, nil
}

func decodePokemonMap(data []byte) (PokemonMapResponse, error) {
	var pokemonMap PokemonMapResponse
	err := json.Unmarshal(data, &pokemonMap)
	return pokemonMap, err
}
