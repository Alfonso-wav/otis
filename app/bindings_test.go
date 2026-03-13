package app

import (
	"errors"
	"testing"

	"github.com/alfon/pokemon-app/core"
)

// mockFetcher implementa core.PokemonFetcher para tests.
type mockFetcher struct {
	pokemon     core.Pokemon
	pokemonErr  error
	list        core.PokemonListResponse
	listErr     error
}

func (m *mockFetcher) FetchPokemon(name string) (core.Pokemon, error) {
	return m.pokemon, m.pokemonErr
}

func (m *mockFetcher) FetchPokemonList(offset int, limit int) (core.PokemonListResponse, error) {
	return m.list, m.listErr
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
