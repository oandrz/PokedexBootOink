package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokemonMapResponse struct {
	Next              *string       `json:"next"`
	Previous          *string       `json:"previous"`
	Results           []MapLocation `json:"results"`
	PokemonsEncounter []struct {
		PokemonEncounter struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type MapLocation struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (c *Client) GetPokemonMapLocation(url string) (PokemonMapResponse, error) {
	if cached, ok := c.cache.Get(url); ok {
		fmt.Println("Using cached data")
		pokemonMap, err := decodePokemonMap(cached)
		if err != nil {
			return PokemonMapResponse{}, err
		}
		return pokemonMap, nil
	}

	fmt.Println("Using Remote data")

	return c.fetchRemotePokemonMapLocation(url)
}

func (c *Client) fetchRemotePokemonMapLocation(url string) (PokemonMapResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return PokemonMapResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PokemonMapResponse{}, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonMapResponse{}, err
	}
	c.cache.Add(url, data)

	pokemonMap, err := decodePokemonMap(data)
	if err != nil {
		return PokemonMapResponse{}, err
	}

	return pokemonMap, nil
}

func decodePokemonMap(data []byte) (PokemonMapResponse, error) {
	var pokemonMap PokemonMapResponse
	err := json.Unmarshal(data, &pokemonMap)
	return pokemonMap, err
}
