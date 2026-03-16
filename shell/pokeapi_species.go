package shell

import (
	"fmt"
	"sync"

	"github.com/alfon/pokemon-app/core"
)

var (
	speciesClassCache map[string]core.SpeciesClassification
	speciesClassMu    sync.Mutex
)

// --- apiPokemonSpecies ---

type apiPokemonSpecies struct {
	Name                 string `json:"name"`
	Order                int    `json:"order"`
	GenderRate           int    `json:"gender_rate"`
	CaptureRate          int    `json:"capture_rate"`
	BaseHappiness        int    `json:"base_happiness"`
	IsBaby               bool   `json:"is_baby"`
	IsLegendary          bool   `json:"is_legendary"`
	IsMythical           bool   `json:"is_mythical"`
	HatchCounter         int    `json:"hatch_counter"`
	HasGenderDifferences bool   `json:"has_gender_differences"`
	FormsSwitchable      bool   `json:"forms_switchable"`
	Color                struct {
		Name string `json:"name"`
	} `json:"color"`
	Shape struct {
		Name string `json:"name"`
	} `json:"shape"`
	Habitat *struct {
		Name string `json:"name"`
	} `json:"habitat"`
	EggGroups []struct {
		Name string `json:"name"`
	} `json:"egg_groups"`
	Genera []struct {
		Genus    string `json:"genus"`
		Language struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"genera"`
	FlavorTextEntries []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"flavor_text_entries"`
	EvolutionChain struct {
		URL string `json:"url"`
	} `json:"evolution_chain"`
	Varieties []struct {
		IsDefault bool `json:"is_default"`
		Pokemon   struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"varieties"`
}

func toDomainSpecies(raw apiPokemonSpecies) core.PokemonSpecies {
	eggGroups := make([]string, len(raw.EggGroups))
	for i, eg := range raw.EggGroups {
		eggGroups[i] = eg.Name
	}

	genus := ""
	for _, g := range raw.Genera {
		if g.Language.Name == "en" {
			genus = g.Genus
			break
		}
	}

	// Tomar la última entrada en inglés (texto del juego más reciente).
	flavorText := ""
	for _, ft := range raw.FlavorTextEntries {
		if ft.Language.Name == "en" {
			flavorText = ft.FlavorText
		}
	}

	habitat := ""
	if raw.Habitat != nil {
		habitat = raw.Habitat.Name
	}

	varieties := make([]core.PokemonVariety, len(raw.Varieties))
	for i, v := range raw.Varieties {
		varieties[i] = core.PokemonVariety{IsDefault: v.IsDefault, Pokemon: v.Pokemon.Name}
	}

	return core.PokemonSpecies{
		Name:                 raw.Name,
		Order:                raw.Order,
		GenderRate:           raw.GenderRate,
		CaptureRate:          raw.CaptureRate,
		BaseHappiness:        raw.BaseHappiness,
		IsBaby:               raw.IsBaby,
		IsLegendary:          raw.IsLegendary,
		IsMythical:           raw.IsMythical,
		HatchCounter:         raw.HatchCounter,
		HasGenderDifferences: raw.HasGenderDifferences,
		FormsSwitchable:      raw.FormsSwitchable,
		Genus:                genus,
		Color:                raw.Color.Name,
		Shape:                raw.Shape.Name,
		Habitat:              habitat,
		EggGroups:            eggGroups,
		FlavorText:           flavorText,
		EvolutionChainID:     extractIDFromURL(raw.EvolutionChain.URL),
		Varieties:            varieties,
	}
}

func (c *PokeAPIClient) FetchPokemonSpecies(name string) (core.PokemonSpecies, error) {
	url := fmt.Sprintf("%s/pokemon-species/%s", c.baseURL, name)
	var raw apiPokemonSpecies
	if err := fetchOne(c, url, &raw); err != nil {
		return core.PokemonSpecies{}, fmt.Errorf("fetching pokemon species %q: %w", name, err)
	}
	return toDomainSpecies(raw), nil
}

// --- apiPokemonForm ---

type apiPokemonForm struct {
	Name         string `json:"name"`
	FormName     string `json:"form_name"`
	IsMega       bool   `json:"is_mega"`
	IsBattleOnly bool   `json:"is_battle_only"`
	Types        []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
		FrontShiny   string `json:"front_shiny"`
	} `json:"sprites"`
}

func (c *PokeAPIClient) FetchPokemonForm(name string) (core.PokemonForm, error) {
	url := fmt.Sprintf("%s/pokemon-form/%s", c.baseURL, name)
	var raw apiPokemonForm
	if err := fetchOne(c, url, &raw); err != nil {
		return core.PokemonForm{}, fmt.Errorf("fetching pokemon form %q: %w", name, err)
	}

	types := make([]core.PokemonType, len(raw.Types))
	for i, t := range raw.Types {
		types[i] = core.PokemonType{Name: t.Type.Name}
	}

	return core.PokemonForm{
		Name:         raw.Name,
		FormName:     raw.FormName,
		IsMega:       raw.IsMega,
		IsBattleOnly: raw.IsBattleOnly,
		Types:        types,
		Sprites: core.Sprites{
			FrontDefault: raw.Sprites.FrontDefault,
			FrontShiny:   raw.Sprites.FrontShiny,
		},
	}, nil
}

// --- Listas de clasificación ---

func (c *PokeAPIClient) FetchPokemonColors() ([]core.NamedResource, error) {
	return fetchNamedResourceList(c, "pokemon-color")
}

func (c *PokeAPIClient) FetchPokemonShapes() ([]core.NamedResource, error) {
	return fetchNamedResourceList(c, "pokemon-shape")
}

func (c *PokeAPIClient) FetchPokemonHabitats() ([]core.NamedResource, error) {
	return fetchNamedResourceList(c, "pokemon-habitat")
}

// FetchAllSpeciesClassifications returns a map of pokemon name → classification (legendary/mythical),
// fetching all species in parallel with a concurrency limit and caching the result in memory.
func (c *PokeAPIClient) FetchAllSpeciesClassifications() (map[string]core.SpeciesClassification, error) {
	speciesClassMu.Lock()
	if speciesClassCache != nil {
		result := make(map[string]core.SpeciesClassification, len(speciesClassCache))
		for k, v := range speciesClassCache {
			result[k] = v
		}
		speciesClassMu.Unlock()
		return result, nil
	}
	speciesClassMu.Unlock()

	// Get total pokemon count
	list, err := c.FetchPokemonList(0, 2000)
	if err != nil {
		return nil, fmt.Errorf("fetching pokemon list for classifications: %w", err)
	}

	type result struct {
		name string
		cls  core.SpeciesClassification
	}

	results := make(chan result, len(list.Results))
	sem := make(chan struct{}, 30)
	var wg sync.WaitGroup

	for _, item := range list.Results {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			species, err := c.FetchPokemonSpecies(name)
			if err != nil {
				return
			}
			results <- result{
				name: name,
				cls: core.SpeciesClassification{
					IsLegendary: species.IsLegendary,
					IsMythical:  species.IsMythical,
				},
			}
		}(item.Name)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	classMap := make(map[string]core.SpeciesClassification, len(list.Results))
	for res := range results {
		classMap[res.name] = res.cls
	}

	speciesClassMu.Lock()
	speciesClassCache = make(map[string]core.SpeciesClassification, len(classMap))
	for k, v := range classMap {
		speciesClassCache[k] = v
	}
	speciesClassMu.Unlock()

	return classMap, nil
}
