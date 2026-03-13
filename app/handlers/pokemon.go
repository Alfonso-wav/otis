package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/alfon/pokemon-app/core"
)

type PokemonHandler struct {
	fetcher core.PokemonFetcher
}

func NewPokemonHandler(fetcher core.PokemonFetcher) *PokemonHandler {
	return &PokemonHandler{fetcher: fetcher}
}

func (h *PokemonHandler) ListPokemon(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset := 0
	limit := 20

	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil && v >= 0 {
			offset = v
		}
	}
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	list, err := h.fetcher.FetchPokemonList(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, list)
}

func (h *PokemonHandler) GetPokemon(w http.ResponseWriter, r *http.Request) {
	// Extrae el nombre del path: /api/pokemon/{name}
	name := strings.TrimPrefix(r.URL.Path, "/api/pokemon/")
	name = strings.TrimSpace(name)
	if name == "" {
		http.Error(w, "pokemon name is required", http.StatusBadRequest)
		return
	}

	pokemon, err := h.fetcher.FetchPokemon(core.NormalizeName(name))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, pokemon)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
