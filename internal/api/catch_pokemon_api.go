package api

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

const basePokemonUrl = "https://pokeapi.co/api/v2/pokemon/"

type PokemonDetailResponse struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

/*
* CatchPokemon
Core formula: catchChance = 50 / baseExperience

	This creates an inverse relationship — higher base experience → lower catch chance.

	Example Calculations

	| Pokemon   | Base Exp | Calculation   | Catch Chance   |
	|-----------|----------|---------------|----------------|
	| Caterpie  | 39       | 50/39 = 1.28  | 100% (clamped) |
	| Pikachu   | 112      | 50/112 = 0.44 | 44%            |
	| Charizard | 267      | 50/267 = 0.18 | 18%            |
	| Mewtwo    | 340      | 50/340 = 0.14 | 14%            |

	Why clamp to 1.0?

	if catchChance > 1.0 {
	    catchChance = 1.0
	}

	When baseExperience < 50, the math gives us values greater than 1.0 (like 1.28 for Caterpie). But a probability can't exceed 100%, so we cap it.

	The "50" is your difficulty knob

	- Higher value (e.g., 80) → easier catches overall
	- Lower value (e.g., 30) → harder catches overall

	Think of it as: "You need roughly this much 'luck' to catch any Pokemon, scaled by their strength."
*/
func (c *Client) CatchPokemon(pokemonName string) (bool, error) {
	pokemonDetail, err := c.fetchPokemonDetailRemotely(pokemonName)
	if err != nil {
		return false, err
	}

	catchChance := float64(50) / float64(pokemonDetail.BaseExperience)
	if catchChance > 1.0 {
		catchChance = 1.0
	}

	roll := rand.Float64()
	isSuccess := roll <= catchChance

	if isSuccess {
		c.pokedex[pokemonName] = Pokemon{Name: pokemonName}
	}

	return isSuccess, nil
}

func (c *Client) fetchPokemonDetailRemotely(pokemonName string) (PokemonDetailResponse, error) {
	url := basePokemonUrl + pokemonName
	resp, err := http.Get(url)
	if err != nil {
		return PokemonDetailResponse{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonDetailResponse{}, err
	}

	pokemonDetail, err := decodePokemonDetail(data)
	if err != nil {
		return PokemonDetailResponse{}, err
	}

	return pokemonDetail, nil
}

func decodePokemonDetail(data []byte) (PokemonDetailResponse, error) {
	var pokemonDetail PokemonDetailResponse
	err := json.Unmarshal(data, &pokemonDetail)
	return pokemonDetail, err
}
