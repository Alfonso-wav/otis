package shell

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alfon/pokemon-app/core"
)

type PokeAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPokeAPIClient(baseURL string) *PokeAPIClient {
	return &PokeAPIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// apiPokemon es la estructura raw que devuelve PokéAPI para /pokemon/{name}.
type apiPokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
		FrontShiny   string `json:"front_shiny"`
	} `json:"sprites"`
}

// apiList es la estructura raw que devuelve PokéAPI para /pokemon?offset=&limit=.
type apiList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *PokeAPIClient) FetchPokemon(name string) (core.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return core.Pokemon{}, fmt.Errorf("fetching pokemon %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.Pokemon{}, fmt.Errorf("pokemon %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return core.Pokemon{}, fmt.Errorf("pokeapi returned status %d for pokemon %q", resp.StatusCode, name)
	}

	var raw apiPokemon
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return core.Pokemon{}, fmt.Errorf("decoding pokemon %q: %w", name, err)
	}

	return toDomainPokemon(raw), nil
}

func (c *PokeAPIClient) FetchPokemonList(offset int, limit int) (core.PokemonListResponse, error) {
	url := fmt.Sprintf("%s/pokemon?offset=%d&limit=%d", c.baseURL, offset, limit)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return core.PokemonListResponse{}, fmt.Errorf("fetching pokemon list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return core.PokemonListResponse{}, fmt.Errorf("pokeapi returned status %d for pokemon list", resp.StatusCode)
	}

	var raw apiList
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return core.PokemonListResponse{}, fmt.Errorf("decoding pokemon list: %w", err)
	}

	return toDomainList(raw), nil
}

func toDomainPokemon(raw apiPokemon) core.Pokemon {
	types := make([]core.PokemonType, len(raw.Types))
	for i, t := range raw.Types {
		types[i] = core.PokemonType{Name: t.Type.Name}
	}

	stats := make([]core.Stat, len(raw.Stats))
	for i, s := range raw.Stats {
		stats[i] = core.Stat{Name: s.Stat.Name, BaseStat: s.BaseStat}
	}

	return core.Pokemon{
		ID:   raw.ID,
		Name: raw.Name,
		Types: types,
		Stats: stats,
		Sprites: core.Sprites{
			FrontDefault: raw.Sprites.FrontDefault,
			FrontShiny:   raw.Sprites.FrontShiny,
		},
		Height: raw.Height,
		Weight: raw.Weight,
	}
}

func toDomainList(raw apiList) core.PokemonListResponse {
	results := make([]core.PokemonListItem, len(raw.Results))
	for i, r := range raw.Results {
		results[i] = core.PokemonListItem{Name: r.Name, URL: r.URL}
	}
	return core.PokemonListResponse{
		Count:    raw.Count,
		Next:     raw.Next,
		Previous: raw.Previous,
		Results:  results,
	}
}
