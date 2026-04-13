package shell

import (
	"testing"
)

// TestToDomainPokemon_ParsesAbilitiesWithIsHidden verifies that the
// shell→core mapping preserves the is_hidden flag from PokéAPI.
func TestToDomainPokemon_ParsesAbilitiesWithIsHidden(t *testing.T) {
	raw := apiPokemon{
		ID:   1,
		Name: "bulbasaur",
	}
	raw.Abilities = []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
	}{
		{Ability: struct {
			Name string `json:"name"`
		}{Name: "overgrow"}, IsHidden: false},
		{Ability: struct {
			Name string `json:"name"`
		}{Name: "chlorophyll"}, IsHidden: true},
	}
	p := toDomainPokemon(raw)
	if len(p.Abilities) != 2 {
		t.Fatalf("abilities count = %d, want 2", len(p.Abilities))
	}
	if p.Abilities[0].Name != "overgrow" || p.Abilities[0].IsHidden {
		t.Errorf("slot 0: want overgrow hidden=false, got %+v", p.Abilities[0])
	}
	if p.Abilities[1].Name != "chlorophyll" || !p.Abilities[1].IsHidden {
		t.Errorf("slot 1: want chlorophyll hidden=true, got %+v", p.Abilities[1])
	}
}
