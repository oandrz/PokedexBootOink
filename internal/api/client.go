package api

import (
	"net/http"
	"pokedex_go/internal/pokeCache"
	"time"
)

type Client struct {
	cache      *pokeCache.Cache
	httpClient http.Client
	Pokedex    map[string]Pokemon
}

// NewClient -
func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokeCache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
		Pokedex: make(map[string]Pokemon),
	}
}
