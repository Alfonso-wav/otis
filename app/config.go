package app

import "os"

type Config struct {
	PokeAPIBaseURL string
}

func LoadConfig() Config {
	baseURL := os.Getenv("POKEAPI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://pokeapi.co/api/v2"
	}

	return Config{
		PokeAPIBaseURL: baseURL,
	}
}
