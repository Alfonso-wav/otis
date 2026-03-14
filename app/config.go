package app

import "os"

type Config struct {
	PokeAPIBaseURL    string
	PokemonDBBaseURL  string
}

func LoadConfig() Config {
	baseURL := os.Getenv("POKEAPI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://pokeapi.co/api/v2"
	}

	pokemonDBURL := os.Getenv("POKEMONDB_BASE_URL")
	if pokemonDBURL == "" {
		pokemonDBURL = "https://pokemondb.net"
	}

	return Config{
		PokeAPIBaseURL:   baseURL,
		PokemonDBBaseURL: pokemonDBURL,
	}
}
