package api

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

const basePokemonUrl = "https://pokeapi.co/api/v2/pokemon/"

type PokemonDetailResponse struct {
	Name           string                `json:"name"`
	BaseExperience int                   `json:"base_experience"`
	Weight         int                   `json:"weight"`
	Height         int                   `json:"height"`
	Types          []PokemonTypeResponse `json:"types"`
	Stats          []PokemonStatResponse `json:"stats"`
}

type PokemonTypeResponse struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type PokemonStatResponse struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
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
	if err != nil || pokemonDetail.Name == "" {
		return false, err
	}

	catchChance := float64(50) / float64(pokemonDetail.BaseExperience)
	if catchChance > 1.0 {
		catchChance = 1.0
	}

	roll := rand.Float64()
	isSuccess := roll <= catchChance

	if isSuccess {
		c.Pokedex[pokemonName] = Pokemon{
			Name:   pokemonName,
			Weight: pokemonDetail.Weight,
			Height: pokemonDetail.Height,
			Types: Map(pokemonDetail.Types, func(p PokemonTypeResponse) PokemonType {
				return PokemonType{
					Name: p.Type.Name,
				}
			}),
			Stats: Map(pokemonDetail.Stats, func(p PokemonStatResponse) PokemonStat {
				return PokemonStat{
					Name:  p.Stat.Name,
					Value: p.BaseStat,
				}
			}),
		}
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

	if resp.StatusCode != http.StatusOK {
		return PokemonDetailResponse{}, err
	}

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

func Map[T any, U any](slice []T, transform func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}
