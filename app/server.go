package app

import (
	"net/http"

	"github.com/alfon/pokemon-app/app/handlers"
	"github.com/alfon/pokemon-app/core"
)

func NewServer(cfg Config, fetcher core.PokemonFetcher) *http.ServeMux {
	mux := http.NewServeMux()
	h := handlers.NewPokemonHandler(fetcher)

	// API endpoints
	mux.HandleFunc("/api/pokemon", h.ListPokemon)
	mux.HandleFunc("/api/pokemon/", h.GetPokemon)

	// Frontend: sirve archivos estaticos desde ./frontend/
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	return mux
}
