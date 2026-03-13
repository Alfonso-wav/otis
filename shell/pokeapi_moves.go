package shell

import (
	"fmt"

	"github.com/alfon/pokemon-app/core"
)

// --- Lista de movimientos ---

type apiMoveList struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *PokeAPIClient) FetchMoveList(offset int, limit int) (core.MoveListResponse, error) {
	url := fmt.Sprintf("%s/move?offset=%d&limit=%d", c.baseURL, offset, limit)
	var raw apiMoveList
	if err := fetchOne(c, url, &raw); err != nil {
		return core.MoveListResponse{}, fmt.Errorf("fetching move list: %w", err)
	}

	results := make([]core.NamedResource, len(raw.Results))
	for i, r := range raw.Results {
		results[i] = core.NamedResource{Name: r.Name, URL: r.URL}
	}
	return core.MoveListResponse{Count: raw.Count, Results: results}, nil
}

// --- Clases de daño ---

type apiMoveDamageClass struct {
	Name         string `json:"name"`
	Descriptions []struct {
		Description string `json:"description"`
		Language    struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"descriptions"`
	Moves []struct {
		Name string `json:"name"`
	} `json:"moves"`
}

func (c *PokeAPIClient) FetchMoveDamageClass(name string) (core.MoveDamageClass, error) {
	url := fmt.Sprintf("%s/move-damage-class/%s", c.baseURL, name)
	var raw apiMoveDamageClass
	if err := fetchOne(c, url, &raw); err != nil {
		return core.MoveDamageClass{}, fmt.Errorf("fetching move damage class %q: %w", name, err)
	}

	desc := firstEnglishDesc(raw.Descriptions)

	moves := make([]string, len(raw.Moves))
	for i, m := range raw.Moves {
		moves[i] = m.Name
	}
	return core.MoveDamageClass{Name: raw.Name, Description: desc, Moves: moves}, nil
}

// --- Alteraciones de estado ---

type apiMoveAilment struct {
	Name  string `json:"name"`
	Moves []struct {
		Name string `json:"name"`
	} `json:"moves"`
}

func (c *PokeAPIClient) FetchMoveAilment(name string) (core.MoveAilment, error) {
	url := fmt.Sprintf("%s/move-ailment/%s", c.baseURL, name)
	var raw apiMoveAilment
	if err := fetchOne(c, url, &raw); err != nil {
		return core.MoveAilment{}, fmt.Errorf("fetching move ailment %q: %w", name, err)
	}

	moves := make([]string, len(raw.Moves))
	for i, m := range raw.Moves {
		moves[i] = m.Name
	}
	return core.MoveAilment{Name: raw.Name, Moves: moves}, nil
}

// --- Objetivos de movimiento ---

type apiMoveTarget struct {
	Name         string `json:"name"`
	Descriptions []struct {
		Description string `json:"description"`
		Language    struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"descriptions"`
}

func (c *PokeAPIClient) FetchMoveTarget(name string) (core.MoveTarget, error) {
	url := fmt.Sprintf("%s/move-target/%s", c.baseURL, name)
	var raw apiMoveTarget
	if err := fetchOne(c, url, &raw); err != nil {
		return core.MoveTarget{}, fmt.Errorf("fetching move target %q: %w", name, err)
	}

	desc := firstEnglishDesc(raw.Descriptions)
	return core.MoveTarget{Name: raw.Name, Description: desc}, nil
}

// --- Máquinas (TM/HM) ---

type apiMachine struct {
	ID   int `json:"id"`
	Item struct {
		Name string `json:"name"`
	} `json:"item"`
	Move struct {
		Name string `json:"name"`
	} `json:"move"`
	VersionGroup struct {
		Name string `json:"name"`
	} `json:"version_group"`
}

func (c *PokeAPIClient) FetchMachine(id int) (core.Machine, error) {
	url := fmt.Sprintf("%s/machine/%d", c.baseURL, id)
	var raw apiMachine
	if err := fetchOne(c, url, &raw); err != nil {
		return core.Machine{}, fmt.Errorf("fetching machine %d: %w", id, err)
	}

	return core.Machine{
		ID:           raw.ID,
		Move:         raw.Move.Name,
		Item:         raw.Item.Name,
		VersionGroup: raw.VersionGroup.Name,
	}, nil
}
