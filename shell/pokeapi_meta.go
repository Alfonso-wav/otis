package shell

import (
	"fmt"

	"github.com/alfon/pokemon-app/core"
)

// --- Grupo D: Lista de habilidades ---

type apiAbilityList struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *PokeAPIClient) FetchAbilityList(offset int, limit int) (core.AbilityListResponse, error) {
	url := fmt.Sprintf("%s/ability?offset=%d&limit=%d", c.baseURL, offset, limit)
	var raw apiAbilityList
	if err := fetchOne(c, url, &raw); err != nil {
		return core.AbilityListResponse{}, fmt.Errorf("fetching ability list: %w", err)
	}

	results := make([]core.NamedResource, len(raw.Results))
	for i, r := range raw.Results {
		results[i] = core.NamedResource{Name: r.Name, URL: r.URL}
	}
	return core.AbilityListResponse{Count: raw.Count, Results: results}, nil
}

// --- Grupo G: Stats ---

type apiStatDetail struct {
	Name         string `json:"name"`
	IsBattleOnly bool   `json:"is_battle_only"`
	AffectingMoves struct {
		Increase []struct {
			Move struct{ Name string `json:"name"` } `json:"move"`
		} `json:"increase"`
		Decrease []struct {
			Move struct{ Name string `json:"name"` } `json:"move"`
		} `json:"decrease"`
	} `json:"affecting_moves"`
	AffectingNatures struct {
		Increase []struct {
			Name string `json:"name"`
		} `json:"increase"`
		Decrease []struct {
			Name string `json:"name"`
		} `json:"decrease"`
	} `json:"affecting_natures"`
}

func (c *PokeAPIClient) FetchStat(name string) (core.StatDetail, error) {
	url := fmt.Sprintf("%s/stat/%s", c.baseURL, name)
	var raw apiStatDetail
	if err := fetchOne(c, url, &raw); err != nil {
		return core.StatDetail{}, fmt.Errorf("fetching stat %q: %w", name, err)
	}

	buff := make([]string, len(raw.AffectingMoves.Increase))
	for i, m := range raw.AffectingMoves.Increase {
		buff[i] = m.Move.Name
	}
	nerf := make([]string, len(raw.AffectingMoves.Decrease))
	for i, m := range raw.AffectingMoves.Decrease {
		nerf[i] = m.Move.Name
	}
	natureBuff := make([]string, len(raw.AffectingNatures.Increase))
	for i, n := range raw.AffectingNatures.Increase {
		natureBuff[i] = n.Name
	}
	natureNerf := make([]string, len(raw.AffectingNatures.Decrease))
	for i, n := range raw.AffectingNatures.Decrease {
		natureNerf[i] = n.Name
	}

	return core.StatDetail{
		Name:                 raw.Name,
		IsBattleOnly:         raw.IsBattleOnly,
		AffectingMovesBuff:   buff,
		AffectingMovesNerf:   nerf,
		AffectingNaturesBuff: natureBuff,
		AffectingNaturesNerf: natureNerf,
	}, nil
}

// --- Generaciones ---

type apiGeneration struct {
	Name       string `json:"name"`
	MainRegion struct {
		Name string `json:"name"`
	} `json:"main_region"`
	VersionGroups []struct {
		Name string `json:"name"`
	} `json:"version_groups"`
	PokemonSpecies []struct {
		Name string `json:"name"`
	} `json:"pokemon_species"`
	Types []struct {
		Name string `json:"name"`
	} `json:"types"`
	Moves []struct {
		Name string `json:"name"`
	} `json:"moves"`
	Abilities []struct {
		Name string `json:"name"`
	} `json:"abilities"`
}

func toDomainGeneration(raw apiGeneration) core.Generation {
	games := make([]string, len(raw.VersionGroups))
	for i, vg := range raw.VersionGroups {
		games[i] = vg.Name
	}
	species := make([]string, len(raw.PokemonSpecies))
	for i, p := range raw.PokemonSpecies {
		species[i] = p.Name
	}
	types := make([]string, len(raw.Types))
	for i, t := range raw.Types {
		types[i] = t.Name
	}
	moves := make([]string, len(raw.Moves))
	for i, m := range raw.Moves {
		moves[i] = m.Name
	}
	abilities := make([]string, len(raw.Abilities))
	for i, a := range raw.Abilities {
		abilities[i] = a.Name
	}
	return core.Generation{
		Name:           raw.Name,
		MainRegion:     raw.MainRegion.Name,
		Games:          games,
		PokemonSpecies: species,
		Types:          types,
		Moves:          moves,
		Abilities:      abilities,
	}
}

func (c *PokeAPIClient) FetchGenerations() ([]core.NamedResource, error) {
	return fetchNamedResourceList(c, "generation")
}

func (c *PokeAPIClient) FetchGeneration(name string) (core.Generation, error) {
	url := fmt.Sprintf("%s/generation/%s", c.baseURL, name)
	var raw apiGeneration
	if err := fetchOne(c, url, &raw); err != nil {
		return core.Generation{}, fmt.Errorf("fetching generation %q: %w", name, err)
	}
	return toDomainGeneration(raw), nil
}

// --- Pokédex ---

type apiPokedex struct {
	Name         string `json:"name"`
	IsMainSeries bool   `json:"is_main_series"`
	Region       *struct {
		Name string `json:"name"`
	} `json:"region"`
	PokemonEntries []struct {
		EntryNumber    int `json:"entry_number"`
		PokemonSpecies struct {
			Name string `json:"name"`
		} `json:"pokemon_species"`
	} `json:"pokemon_entries"`
}

func (c *PokeAPIClient) FetchPokedexList() ([]core.NamedResource, error) {
	return fetchNamedResourceList(c, "pokedex")
}

func (c *PokeAPIClient) FetchPokedex(name string) (core.Pokedex, error) {
	url := fmt.Sprintf("%s/pokedex/%s", c.baseURL, name)
	var raw apiPokedex
	if err := fetchOne(c, url, &raw); err != nil {
		return core.Pokedex{}, fmt.Errorf("fetching pokedex %q: %w", name, err)
	}

	region := ""
	if raw.Region != nil {
		region = raw.Region.Name
	}

	entries := make([]core.PokedexEntry, len(raw.PokemonEntries))
	for i, e := range raw.PokemonEntries {
		entries[i] = core.PokedexEntry{EntryNumber: e.EntryNumber, Pokemon: e.PokemonSpecies.Name}
	}

	return core.Pokedex{
		Name:           raw.Name,
		IsMainSeries:   raw.IsMainSeries,
		Region:         region,
		PokemonEntries: entries,
	}, nil
}

// --- Grupos de versión ---

type apiVersionGroup struct {
	Name       string `json:"name"`
	Order      int    `json:"order"`
	Generation struct {
		Name string `json:"name"`
	} `json:"generation"`
	Versions []struct {
		Name string `json:"name"`
	} `json:"versions"`
	Pokedexes []struct {
		Name string `json:"name"`
	} `json:"pokedexes"`
	Regions []struct {
		Name string `json:"name"`
	} `json:"regions"`
}

func (c *PokeAPIClient) FetchVersionGroups() ([]core.NamedResource, error) {
	return fetchNamedResourceList(c, "version-group")
}

func (c *PokeAPIClient) FetchVersionGroup(name string) (core.VersionGroup, error) {
	url := fmt.Sprintf("%s/version-group/%s", c.baseURL, name)
	var raw apiVersionGroup
	if err := fetchOne(c, url, &raw); err != nil {
		return core.VersionGroup{}, fmt.Errorf("fetching version group %q: %w", name, err)
	}

	versions := make([]string, len(raw.Versions))
	for i, v := range raw.Versions {
		versions[i] = v.Name
	}
	pokedexes := make([]string, len(raw.Pokedexes))
	for i, p := range raw.Pokedexes {
		pokedexes[i] = p.Name
	}
	regions := make([]string, len(raw.Regions))
	for i, r := range raw.Regions {
		regions[i] = r.Name
	}

	return core.VersionGroup{
		Name:       raw.Name,
		Order:      raw.Order,
		Generation: raw.Generation.Name,
		Versions:   versions,
		Pokedexes:  pokedexes,
		Regions:    regions,
	}, nil
}
