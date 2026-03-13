package app

import (
	"errors"
	"testing"

	"github.com/alfon/pokemon-app/core"
)

// mockFetcher implementa core.PokemonFetcher para tests.
type mockFetcher struct {
	pokemon    core.Pokemon
	pokemonErr error
	list       core.PokemonListResponse
	listErr    error
	typeList   core.TypeListResponse
	typeListErr error
	typeDetail  core.PokemonTypeDetail
	typeDetailErr error
}

func (m *mockFetcher) FetchPokemon(name string) (core.Pokemon, error) {
	return m.pokemon, m.pokemonErr
}

func (m *mockFetcher) FetchPokemonList(offset int, limit int) (core.PokemonListResponse, error) {
	return m.list, m.listErr
}

func (m *mockFetcher) FetchTypeList() (core.TypeListResponse, error) {
	return m.typeList, m.typeListErr
}

func (m *mockFetcher) FetchType(name string) (core.PokemonTypeDetail, error) {
	return m.typeDetail, m.typeDetailErr
}

func (m *mockFetcher) FetchRegions() ([]core.Region, error) {
	return nil, nil
}

func (m *mockFetcher) FetchRegion(name string) (core.Region, error) {
	return core.Region{}, nil
}

func (m *mockFetcher) FetchMove(name string) (core.Move, error) {
	return core.Move{}, nil
}

func (m *mockFetcher) FetchAbility(name string) (core.Ability, error) {
	return core.Ability{}, nil
}

func (m *mockFetcher) FetchEvolutionChain(id int) (core.EvolutionChain, error) {
	return core.EvolutionChain{}, nil
}

// --- Grupo A ---
func (m *mockFetcher) FetchPokemonSpecies(name string) (core.PokemonSpecies, error) {
	return core.PokemonSpecies{}, nil
}
func (m *mockFetcher) FetchPokemonForm(name string) (core.PokemonForm, error) {
	return core.PokemonForm{}, nil
}
func (m *mockFetcher) FetchPokemonColors() ([]core.NamedResource, error)   { return nil, nil }
func (m *mockFetcher) FetchPokemonShapes() ([]core.NamedResource, error)   { return nil, nil }
func (m *mockFetcher) FetchPokemonHabitats() ([]core.NamedResource, error) { return nil, nil }

// --- Grupo B ---
func (m *mockFetcher) FetchNatureList() ([]core.NamedResource, error) { return nil, nil }
func (m *mockFetcher) FetchNatureDetail(name string) (core.NatureDetail, error) {
	return core.NatureDetail{}, nil
}
func (m *mockFetcher) FetchEggGroup(name string) (core.EggGroup, error)   { return core.EggGroup{}, nil }
func (m *mockFetcher) FetchGender(name string) (core.Gender, error)       { return core.Gender{}, nil }
func (m *mockFetcher) FetchGrowthRate(name string) (core.GrowthRate, error) {
	return core.GrowthRate{}, nil
}

// --- Grupo C ---
func (m *mockFetcher) FetchMoveList(offset int, limit int) (core.MoveListResponse, error) {
	return core.MoveListResponse{}, nil
}
func (m *mockFetcher) FetchMoveDamageClass(name string) (core.MoveDamageClass, error) {
	return core.MoveDamageClass{}, nil
}
func (m *mockFetcher) FetchMoveAilment(name string) (core.MoveAilment, error) {
	return core.MoveAilment{}, nil
}
func (m *mockFetcher) FetchMoveTarget(name string) (core.MoveTarget, error) {
	return core.MoveTarget{}, nil
}
func (m *mockFetcher) FetchMachine(id int) (core.Machine, error) { return core.Machine{}, nil }

// --- Grupo D ---
func (m *mockFetcher) FetchAbilityList(offset int, limit int) (core.AbilityListResponse, error) {
	return core.AbilityListResponse{}, nil
}

// --- Grupo F ---
func (m *mockFetcher) FetchLocation(name string) (core.LocationDetail, error) {
	return core.LocationDetail{}, nil
}
func (m *mockFetcher) FetchLocationArea(name string) (core.LocationArea, error) {
	return core.LocationArea{}, nil
}

// --- Grupo G ---
func (m *mockFetcher) FetchStat(name string) (core.StatDetail, error) { return core.StatDetail{}, nil }
func (m *mockFetcher) FetchGenerations() ([]core.NamedResource, error) { return nil, nil }
func (m *mockFetcher) FetchGeneration(name string) (core.Generation, error) {
	return core.Generation{}, nil
}
func (m *mockFetcher) FetchPokedexList() ([]core.NamedResource, error) { return nil, nil }
func (m *mockFetcher) FetchPokedex(name string) (core.Pokedex, error)  { return core.Pokedex{}, nil }
func (m *mockFetcher) FetchVersionGroups() ([]core.NamedResource, error) { return nil, nil }
func (m *mockFetcher) FetchVersionGroup(name string) (core.VersionGroup, error) {
	return core.VersionGroup{}, nil
}

func TestGetPokemon(t *testing.T) {
	expected := core.Pokemon{ID: 25, Name: "pikachu"}
	a := NewApp(&mockFetcher{pokemon: expected})

	got, err := a.GetPokemon("Pikachu")
	if err != nil {
		t.Fatalf("GetPokemon error: %v", err)
	}
	if got.ID != expected.ID || got.Name != expected.Name {
		t.Errorf("got %+v, want %+v", got, expected)
	}
}

func TestGetPokemonError(t *testing.T) {
	a := NewApp(&mockFetcher{pokemonErr: errors.New("not found")})

	_, err := a.GetPokemon("unknown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListPokemon(t *testing.T) {
	expected := core.PokemonListResponse{
		Count:   2,
		Results: []core.PokemonListItem{{Name: "bulbasaur"}, {Name: "ivysaur"}},
	}
	a := NewApp(&mockFetcher{list: expected})

	got, err := a.ListPokemon(0, 2)
	if err != nil {
		t.Fatalf("ListPokemon error: %v", err)
	}
	if got.Count != expected.Count || len(got.Results) != len(expected.Results) {
		t.Errorf("got %+v, want %+v", got, expected)
	}
}

func TestListPokemonError(t *testing.T) {
	a := NewApp(&mockFetcher{listErr: errors.New("api error")})

	_, err := a.ListPokemon(0, 20)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListTypes(t *testing.T) {
	expected := core.TypeListResponse{
		Count:   2,
		Results: []core.PokemonListItem{{Name: "fire"}, {Name: "water"}},
	}
	a := NewApp(&mockFetcher{typeList: expected})

	got, err := a.ListTypes()
	if err != nil {
		t.Fatalf("ListTypes error: %v", err)
	}
	if got.Count != expected.Count || len(got.Results) != len(expected.Results) {
		t.Errorf("got %+v, want %+v", got, expected)
	}
}

func TestListTypesError(t *testing.T) {
	a := NewApp(&mockFetcher{typeListErr: errors.New("api error")})

	_, err := a.ListTypes()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetType(t *testing.T) {
	expected := core.PokemonTypeDetail{
		Name:    "fire",
		Pokemon: []core.TypePokemonEntry{{Name: "charmander"}},
	}
	a := NewApp(&mockFetcher{typeDetail: expected})

	got, err := a.GetType("Fire")
	if err != nil {
		t.Fatalf("GetType error: %v", err)
	}
	if got.Name != expected.Name || len(got.Pokemon) != len(expected.Pokemon) {
		t.Errorf("got %+v, want %+v", got, expected)
	}
}

func TestGetTypeError(t *testing.T) {
	a := NewApp(&mockFetcher{typeDetailErr: errors.New("not found")})

	_, err := a.GetType("unknown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
