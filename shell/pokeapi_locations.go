package shell

import (
	"fmt"

	"github.com/alfon/pokemon-app/core"
)

// --- Ubicaciones ---

type apiLocationDetail struct {
	Name   string `json:"name"`
	Region struct {
		Name string `json:"name"`
	} `json:"region"`
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
}

func (c *PokeAPIClient) FetchLocation(name string) (core.LocationDetail, error) {
	url := fmt.Sprintf("%s/location/%s", c.baseURL, name)
	var raw apiLocationDetail
	if err := fetchOne(c, url, &raw); err != nil {
		return core.LocationDetail{}, fmt.Errorf("fetching location %q: %w", name, err)
	}

	areas := make([]string, len(raw.Areas))
	for i, a := range raw.Areas {
		areas[i] = a.Name
	}
	return core.LocationDetail{
		Name:   raw.Name,
		Region: raw.Region.Name,
		Areas:  areas,
	}, nil
}

// --- Áreas de ubicación ---

type apiLocationArea struct {
	Name     string `json:"name"`
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
		VersionDetails []struct {
			MaxChance int `json:"max_chance"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (c *PokeAPIClient) FetchLocationArea(name string) (core.LocationArea, error) {
	url := fmt.Sprintf("%s/location-area/%s", c.baseURL, name)
	var raw apiLocationArea
	if err := fetchOne(c, url, &raw); err != nil {
		return core.LocationArea{}, fmt.Errorf("fetching location area %q: %w", name, err)
	}

	encounters := make([]core.PokemonEncounter, len(raw.PokemonEncounters))
	for i, pe := range raw.PokemonEncounters {
		maxChance := 0
		for _, vd := range pe.VersionDetails {
			if vd.MaxChance > maxChance {
				maxChance = vd.MaxChance
			}
		}
		encounters[i] = core.PokemonEncounter{
			PokemonName: pe.Pokemon.Name,
			MaxChance:   maxChance,
		}
	}

	return core.LocationArea{
		Name:              raw.Name,
		Location:          raw.Location.Name,
		PokemonEncounters: encounters,
	}, nil
}
