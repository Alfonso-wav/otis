package shell

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alfon/pokemon-app/core"
)

// apiEncounterEntry is the raw structure returned by PokeAPI for /pokemon/{name}/encounters.
type apiEncounterEntry struct {
	LocationArea struct {
		Name string `json:"name"`
	} `json:"location_area"`
	VersionDetails []struct {
		Version struct {
			Name string `json:"name"`
		} `json:"version"`
		MaxChance        int `json:"max_chance"`
		EncounterDetails []struct {
			Method struct {
				Name string `json:"name"`
			} `json:"method"`
			Chance         int `json:"chance"`
			MinLevel       int `json:"min_level"`
			MaxLevel       int `json:"max_level"`
			ConditionValues []struct {
				Name string `json:"name"`
			} `json:"condition_values"`
		} `json:"encounter_details"`
	} `json:"version_details"`
}

func (c *PokeAPIClient) FetchPokemonEncounters(name string) ([]core.PokemonLocationEncounter, error) {
	url := fmt.Sprintf("%s/pokemon/%s/encounters", c.baseURL, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching encounters for %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("pokemon %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pokeapi returned status %d for encounters of %q", resp.StatusCode, name)
	}

	var raw []apiEncounterEntry
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding encounters for %q: %w", name, err)
	}

	return toDomainEncounters(raw), nil
}

func toDomainEncounters(raw []apiEncounterEntry) []core.PokemonLocationEncounter {
	encounters := make([]core.PokemonLocationEncounter, len(raw))
	for i, entry := range raw {
		versions := make([]core.VersionEncounter, len(entry.VersionDetails))
		for j, vd := range entry.VersionDetails {
			details := make([]core.EncounterMethodDetail, len(vd.EncounterDetails))
			for k, ed := range vd.EncounterDetails {
				conditions := make([]core.EncounterCondition, len(ed.ConditionValues))
				for l, cv := range ed.ConditionValues {
					conditions[l] = core.EncounterCondition{Name: cv.Name}
				}
				details[k] = core.EncounterMethodDetail{
					Method:     ed.Method.Name,
					Chance:     ed.Chance,
					MinLevel:   ed.MinLevel,
					MaxLevel:   ed.MaxLevel,
					Conditions: conditions,
				}
			}
			versions[j] = core.VersionEncounter{
				Version:   vd.Version.Name,
				MaxChance: vd.MaxChance,
				Details:   details,
			}
		}
		encounters[i] = core.PokemonLocationEncounter{
			LocationArea: entry.LocationArea.Name,
			Versions:     versions,
		}
	}
	return encounters
}
