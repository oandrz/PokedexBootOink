package api

import (
	"net/http"
	"pokedex_go/internal/pokeCache"
	"time"
)

type Client struct {
	cache      *pokeCache.Cache
	httpClient http.Client
}

// NewClient -
func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokeCache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
