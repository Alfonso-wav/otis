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
	Moves []struct {
		Move struct {
			Name string `json:"name"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
			} `json:"move_learn_method"`
		} `json:"version_group_details"`
	} `json:"moves"`
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

func toDomainTypeList(raw apiTypeList) core.TypeListResponse {
	results := make([]core.PokemonListItem, len(raw.Results))
	for i, r := range raw.Results {
		results[i] = core.PokemonListItem{Name: r.Name, URL: r.URL}
	}
	return core.TypeListResponse{Count: raw.Count, Results: results}
}

func toDomainTypeDetail(raw apiTypeDetail) core.PokemonTypeDetail {
	pokemon := make([]core.TypePokemonEntry, len(raw.Pokemon))
	for i, p := range raw.Pokemon {
		pokemon[i] = core.TypePokemonEntry{Name: p.Pokemon.Name, URL: p.Pokemon.URL}
	}
	return core.PokemonTypeDetail{Name: raw.Name, Pokemon: pokemon}
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

	moves := make([]core.PokemonMoveEntry, 0, len(raw.Moves))
	for _, m := range raw.Moves {
		entry := core.PokemonMoveEntry{Name: m.Move.Name}
		if len(m.VersionGroupDetails) > 0 {
			last := m.VersionGroupDetails[len(m.VersionGroupDetails)-1]
			entry.Method = last.MoveLearnMethod.Name
			entry.Level = last.LevelLearnedAt
		}
		moves = append(moves, entry)
	}

	return core.Pokemon{
		ID:    raw.ID,
		Name:  raw.Name,
		Types: types,
		Stats: stats,
		Sprites: core.Sprites{
			FrontDefault: raw.Sprites.FrontDefault,
			FrontShiny:   raw.Sprites.FrontShiny,
		},
		Height: raw.Height,
		Weight: raw.Weight,
		Moves:  moves,
	}
}

// apiTypeList es la estructura raw que devuelve PokéAPI para /type.
type apiTypeList struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// apiTypeDetail es la estructura raw que devuelve PokéAPI para /type/{name}.
type apiTypeDetail struct {
	Name    string `json:"name"`
	Pokemon []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon"`
}

func (c *PokeAPIClient) FetchTypeList() (core.TypeListResponse, error) {
	url := fmt.Sprintf("%s/type?limit=100", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return core.TypeListResponse{}, fmt.Errorf("fetching type list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return core.TypeListResponse{}, fmt.Errorf("pokeapi returned status %d for type list", resp.StatusCode)
	}

	var raw apiTypeList
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return core.TypeListResponse{}, fmt.Errorf("decoding type list: %w", err)
	}

	return toDomainTypeList(raw), nil
}

func (c *PokeAPIClient) FetchType(name string) (core.PokemonTypeDetail, error) {
	url := fmt.Sprintf("%s/type/%s", c.baseURL, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return core.PokemonTypeDetail{}, fmt.Errorf("fetching type %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.PokemonTypeDetail{}, fmt.Errorf("type %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return core.PokemonTypeDetail{}, fmt.Errorf("pokeapi returned status %d for type %q", resp.StatusCode, name)
	}

	var raw apiTypeDetail
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return core.PokemonTypeDetail{}, fmt.Errorf("decoding type %q: %w", name, err)
	}

	return toDomainTypeDetail(raw), nil
}

// --- Region ---

type apiRegionList struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type apiRegion struct {
	Name      string `json:"name"`
	Locations []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"locations"`
}

func (c *PokeAPIClient) FetchRegions() ([]core.Region, error) {
	url := fmt.Sprintf("%s/region?limit=100", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching regions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pokeapi returned status %d for regions", resp.StatusCode)
	}

	var raw apiRegionList
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding regions: %w", err)
	}

	regions := make([]core.Region, len(raw.Results))
	for i, r := range raw.Results {
		regions[i] = core.Region{Name: r.Name}
	}
	return regions, nil
}

func (c *PokeAPIClient) FetchRegion(name string) (core.Region, error) {
	url := fmt.Sprintf("%s/region/%s", c.baseURL, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return core.Region{}, fmt.Errorf("fetching region %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.Region{}, fmt.Errorf("region %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return core.Region{}, fmt.Errorf("pokeapi returned status %d for region %q", resp.StatusCode, name)
	}

	var raw apiRegion
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return core.Region{}, fmt.Errorf("decoding region %q: %w", name, err)
	}

	locations := make([]core.Location, len(raw.Locations))
	for i, l := range raw.Locations {
		locations[i] = core.Location{Name: l.Name, Region: raw.Name}
	}
	return core.Region{Name: raw.Name, Locations: locations}, nil
}

// --- Move ---

type apiMove struct {
	Name     string `json:"name"`
	Type     struct{ Name string `json:"name"` } `json:"type"`
	Power    *int   `json:"power"`
	Accuracy *int   `json:"accuracy"`
	PP       int    `json:"pp"`
	DamageClass struct{ Name string `json:"name"` } `json:"damage_class"`
	FlavorTextEntries []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct{ Name string `json:"name"` } `json:"language"`
	} `json:"flavor_text_entries"`
}

func (c *PokeAPIClient) FetchMove(name string) (core.Move, error) {
	url := fmt.Sprintf("%s/move/%s", c.baseURL, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return core.Move{}, fmt.Errorf("fetching move %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.Move{}, fmt.Errorf("move %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return core.Move{}, fmt.Errorf("pokeapi returned status %d for move %q", resp.StatusCode, name)
	}

	var raw apiMove
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return core.Move{}, fmt.Errorf("decoding move %q: %w", name, err)
	}

	power := 0
	if raw.Power != nil {
		power = *raw.Power
	}
	accuracy := 0
	if raw.Accuracy != nil {
		accuracy = *raw.Accuracy
	}

	desc := ""
	for _, fe := range raw.FlavorTextEntries {
		if fe.Language.Name == "en" {
			desc = fe.FlavorText
			break
		}
	}

	return core.Move{
		Name:        raw.Name,
		Type:        raw.Type.Name,
		Power:       power,
		Accuracy:    accuracy,
		PP:          raw.PP,
		Category:    raw.DamageClass.Name,
		Description: desc,
	}, nil
}

// --- Ability ---

type apiAbility struct {
	Name string `json:"name"`
	FlavorTextEntries []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct{ Name string `json:"name"` } `json:"language"`
	} `json:"flavor_text_entries"`
	Pokemon []struct {
		Pokemon struct{ Name string `json:"name"` } `json:"pokemon"`
	} `json:"pokemon"`
}

func (c *PokeAPIClient) FetchAbility(name string) (core.Ability, error) {
	url := fmt.Sprintf("%s/ability/%s", c.baseURL, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return core.Ability{}, fmt.Errorf("fetching ability %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.Ability{}, fmt.Errorf("ability %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return core.Ability{}, fmt.Errorf("pokeapi returned status %d for ability %q", resp.StatusCode, name)
	}

	var raw apiAbility
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return core.Ability{}, fmt.Errorf("decoding ability %q: %w", name, err)
	}

	desc := ""
	for _, fe := range raw.FlavorTextEntries {
		if fe.Language.Name == "en" {
			desc = fe.FlavorText
			break
		}
	}

	pokemon := make([]string, len(raw.Pokemon))
	for i, p := range raw.Pokemon {
		pokemon[i] = p.Pokemon.Name
	}

	return core.Ability{Name: raw.Name, Description: desc, Pokemon: pokemon}, nil
}

// --- EvolutionChain ---

type apiEvolutionChain struct {
	ID   int            `json:"id"`
	Chain apiChainLink  `json:"chain"`
}

type apiChainLink struct {
	Species struct{ Name string `json:"name"` } `json:"species"`
	EvolutionDetails []struct {
		MinLevel    *int   `json:"min_level"`
		Trigger     struct{ Name string `json:"name"` } `json:"trigger"`
	} `json:"evolution_details"`
	EvolvesTo []apiChainLink `json:"evolves_to"`
}

func toEvolutionStage(link apiChainLink) core.EvolutionStage {
	minLevel := 0
	trigger := ""
	if len(link.EvolutionDetails) > 0 {
		d := link.EvolutionDetails[0]
		if d.MinLevel != nil {
			minLevel = *d.MinLevel
		}
		trigger = d.Trigger.Name
	}

	children := make([]core.EvolutionStage, len(link.EvolvesTo))
	for i, child := range link.EvolvesTo {
		children[i] = toEvolutionStage(child)
	}

	return core.EvolutionStage{
		Name:        link.Species.Name,
		MinLevel:    minLevel,
		TriggerName: trigger,
		EvolvesTo:   children,
	}
}

func (c *PokeAPIClient) FetchEvolutionChain(id int) (core.EvolutionChain, error) {
	url := fmt.Sprintf("%s/evolution-chain/%d", c.baseURL, id)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return core.EvolutionChain{}, fmt.Errorf("fetching evolution chain %d: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return core.EvolutionChain{}, fmt.Errorf("evolution chain %d not found", id)
	}
	if resp.StatusCode != http.StatusOK {
		return core.EvolutionChain{}, fmt.Errorf("pokeapi returned status %d for evolution chain %d", resp.StatusCode, id)
	}

	var raw apiEvolutionChain
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return core.EvolutionChain{}, fmt.Errorf("decoding evolution chain %d: %w", id, err)
	}

	return core.EvolutionChain{ID: raw.ID, Chain: toEvolutionStage(raw.Chain)}, nil
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
