package app

import "os"

type Config struct {
	Port          string
	PokeAPIBaseURL string
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseURL := os.Getenv("POKEAPI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://pokeapi.co/api/v2"
	}

	return Config{
		Port:          port,
		PokeAPIBaseURL: baseURL,
	}
}
