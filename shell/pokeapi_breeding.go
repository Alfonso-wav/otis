package shell

import (
	"fmt"

	"github.com/alfon/pokemon-app/core"
)

// --- Naturalezas ---

type apiNatureDetail struct {
	Name          string `json:"name"`
	IncreasedStat *struct {
		Name string `json:"name"`
	} `json:"increased_stat"`
	DecreasedStat *struct {
		Name string `json:"name"`
	} `json:"decreased_stat"`
	LikesFlavor *struct {
		Name string `json:"name"`
	} `json:"likes_flavor"`
	HatesFlavor *struct {
		Name string `json:"name"`
	} `json:"hates_flavor"`
}

func toDomainNatureDetail(raw apiNatureDetail) core.NatureDetail {
	increased, decreased, likes, hates := "", "", "", ""
	if raw.IncreasedStat != nil {
		increased = raw.IncreasedStat.Name
	}
	if raw.DecreasedStat != nil {
		decreased = raw.DecreasedStat.Name
	}
	if raw.LikesFlavor != nil {
		likes = raw.LikesFlavor.Name
	}
	if raw.HatesFlavor != nil {
		hates = raw.HatesFlavor.Name
	}
	return core.NatureDetail{
		Name:          raw.Name,
		IncreasedStat: increased,
		DecreasedStat: decreased,
		LikesFlavor:   likes,
		HatesFlavor:   hates,
	}
}

func (c *PokeAPIClient) FetchNatureList() ([]core.NamedResource, error) {
	return fetchNamedResourceList(c, "nature")
}

func (c *PokeAPIClient) FetchNatureDetail(name string) (core.NatureDetail, error) {
	url := fmt.Sprintf("%s/nature/%s", c.baseURL, name)
	var raw apiNatureDetail
	if err := fetchOne(c, url, &raw); err != nil {
		return core.NatureDetail{}, fmt.Errorf("fetching nature %q: %w", name, err)
	}
	return toDomainNatureDetail(raw), nil
}

// --- Grupos de huevo ---

type apiEggGroup struct {
	Name           string `json:"name"`
	PokemonSpecies []struct {
		Name string `json:"name"`
	} `json:"pokemon_species"`
}

func (c *PokeAPIClient) FetchEggGroup(name string) (core.EggGroup, error) {
	url := fmt.Sprintf("%s/egg-group/%s", c.baseURL, name)
	var raw apiEggGroup
	if err := fetchOne(c, url, &raw); err != nil {
		return core.EggGroup{}, fmt.Errorf("fetching egg group %q: %w", name, err)
	}

	pokemon := make([]string, len(raw.PokemonSpecies))
	for i, p := range raw.PokemonSpecies {
		pokemon[i] = p.Name
	}
	return core.EggGroup{Name: raw.Name, Pokemon: pokemon}, nil
}

// --- Género ---

type apiGender struct {
	Name                  string `json:"name"`
	PokemonSpeciesDetails []struct {
		PokemonSpecies struct {
			Name string `json:"name"`
		} `json:"pokemon_species"`
	} `json:"pokemon_species_details"`
}

func (c *PokeAPIClient) FetchGender(name string) (core.Gender, error) {
	url := fmt.Sprintf("%s/gender/%s", c.baseURL, name)
	var raw apiGender
	if err := fetchOne(c, url, &raw); err != nil {
		return core.Gender{}, fmt.Errorf("fetching gender %q: %w", name, err)
	}

	pokemon := make([]string, len(raw.PokemonSpeciesDetails))
	for i, p := range raw.PokemonSpeciesDetails {
		pokemon[i] = p.PokemonSpecies.Name
	}
	return core.Gender{Name: raw.Name, Pokemon: pokemon}, nil
}

// --- Tasas de crecimiento ---

type apiGrowthRate struct {
	Name    string `json:"name"`
	Formula string `json:"formula"`
	Levels  []struct {
		Level      int `json:"level"`
		Experience int `json:"experience"`
	} `json:"levels"`
	PokemonSpecies []struct {
		Name string `json:"name"`
	} `json:"pokemon_species"`
}

func (c *PokeAPIClient) FetchGrowthRate(name string) (core.GrowthRate, error) {
	url := fmt.Sprintf("%s/growth-rate/%s", c.baseURL, name)
	var raw apiGrowthRate
	if err := fetchOne(c, url, &raw); err != nil {
		return core.GrowthRate{}, fmt.Errorf("fetching growth rate %q: %w", name, err)
	}

	levels := make([]core.GrowthRateLevel, len(raw.Levels))
	for i, l := range raw.Levels {
		levels[i] = core.GrowthRateLevel{Level: l.Level, Experience: l.Experience}
	}

	pokemon := make([]string, len(raw.PokemonSpecies))
	for i, p := range raw.PokemonSpecies {
		pokemon[i] = p.Name
	}

	return core.GrowthRate{
		Name:    raw.Name,
		Formula: raw.Formula,
		Levels:  levels,
		Pokemon: pokemon,
	}, nil
}
