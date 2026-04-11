package core

import (
	"sort"
	"testing"
)

func TestAggregateEncounters_Dedup(t *testing.T) {
	areas := []LocationArea{
		{
			Name: "area-1",
			PokemonEncounters: []PokemonEncounter{
				{PokemonName: "pikachu", MaxChance: 10},
				{PokemonName: "bulbasaur", MaxChance: 5},
			},
		},
		{
			Name: "area-2",
			PokemonEncounters: []PokemonEncounter{
				{PokemonName: "pikachu", MaxChance: 30},
				{PokemonName: "charmander", MaxChance: 20},
			},
		},
	}

	result := AggregateEncounters(areas)
	sort.Slice(result, func(i, j int) bool { return result[i].PokemonName < result[j].PokemonName })

	if len(result) != 3 {
		t.Fatalf("expected 3 encounters, got %d", len(result))
	}

	expected := map[string]int{
		"bulbasaur":  5,
		"charmander": 20,
		"pikachu":    30,
	}
	for _, enc := range result {
		want, ok := expected[enc.PokemonName]
		if !ok {
			t.Errorf("unexpected pokemon %q", enc.PokemonName)
			continue
		}
		if enc.MaxChance != want {
			t.Errorf("%s: expected MaxChance %d, got %d", enc.PokemonName, want, enc.MaxChance)
		}
	}
}

func TestAggregateEncounters_CapsAt100(t *testing.T) {
	areas := []LocationArea{
		{
			Name: "area-1",
			PokemonEncounters: []PokemonEncounter{
				{PokemonName: "pikachu", MaxChance: 120},
				{PokemonName: "bulbasaur", MaxChance: 80},
				{PokemonName: "charmander", MaxChance: 100},
			},
		},
		{
			Name: "area-2",
			PokemonEncounters: []PokemonEncounter{
				{PokemonName: "pikachu", MaxChance: 150},
			},
		},
	}

	result := AggregateEncounters(areas)
	sort.Slice(result, func(i, j int) bool { return result[i].PokemonName < result[j].PokemonName })

	expected := map[string]int{
		"bulbasaur":  80,
		"charmander": 100,
		"pikachu":    100,
	}
	for _, enc := range result {
		want, ok := expected[enc.PokemonName]
		if !ok {
			t.Errorf("unexpected pokemon %q", enc.PokemonName)
			continue
		}
		if enc.MaxChance != want {
			t.Errorf("%s: expected MaxChance %d, got %d", enc.PokemonName, want, enc.MaxChance)
		}
	}
}

func TestAggregateEncounters_Empty(t *testing.T) {
	result := AggregateEncounters(nil)
	if len(result) != 0 {
		t.Fatalf("expected 0 encounters for nil input, got %d", len(result))
	}

	result = AggregateEncounters([]LocationArea{})
	if len(result) != 0 {
		t.Fatalf("expected 0 encounters for empty input, got %d", len(result))
	}
}
