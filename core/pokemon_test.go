package core

import "testing"

func TestNormalizeName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Pikachu", "pikachu"},
		{"  BULBASAUR  ", "bulbasaur"},
		{"mr-mime", "mr-mime"},
		{"", ""},
	}
	for _, c := range cases {
		got := NormalizeName(c.input)
		if got != c.want {
			t.Errorf("NormalizeName(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestFilterByType(t *testing.T) {
	pokemons := []Pokemon{
		{Name: "pikachu", Types: []PokemonType{{Name: "electric"}}},
		{Name: "charmander", Types: []PokemonType{{Name: "fire"}}},
		{Name: "charizard", Types: []PokemonType{{Name: "fire"}, {Name: "flying"}}},
		{Name: "gyarados", Types: []PokemonType{{Name: "water"}, {Name: "flying"}}},
	}

	got := FilterByType(pokemons, "fire")
	if len(got) != 2 {
		t.Fatalf("FilterByType(fire) len = %d, want 2", len(got))
	}

	got = FilterByType(pokemons, "flying")
	if len(got) != 2 {
		t.Fatalf("FilterByType(flying) len = %d, want 2", len(got))
	}

	got = FilterByType(pokemons, "psychic")
	if len(got) != 0 {
		t.Fatalf("FilterByType(psychic) len = %d, want 0", len(got))
	}
}

func TestFilterByTypeCaseInsensitive(t *testing.T) {
	pokemons := []Pokemon{
		{Name: "pikachu", Types: []PokemonType{{Name: "Electric"}}},
	}
	got := FilterByType(pokemons, "electric")
	if len(got) != 1 {
		t.Fatalf("FilterByType case-insensitive: len = %d, want 1", len(got))
	}
}

func TestFilterMovesByType(t *testing.T) {
	moves := []Move{
		{Name: "flamethrower", Type: "fire", Category: "special"},
		{Name: "ember", Type: "fire", Category: "special"},
		{Name: "surf", Type: "water", Category: "special"},
		{Name: "earthquake", Type: "ground", Category: "physical"},
	}

	got := FilterMovesByType(moves, "fire")
	if len(got) != 2 {
		t.Fatalf("FilterMovesByType(fire) len = %d, want 2", len(got))
	}
	got = FilterMovesByType(moves, "normal")
	if len(got) != 0 {
		t.Fatalf("FilterMovesByType(normal) len = %d, want 0", len(got))
	}
}

func TestFilterMovesByCategory(t *testing.T) {
	moves := []Move{
		{Name: "tackle", Category: "physical"},
		{Name: "flamethrower", Category: "special"},
		{Name: "toxic", Category: "status"},
		{Name: "earthquake", Category: "physical"},
	}

	got := FilterMovesByCategory(moves, "physical")
	if len(got) != 2 {
		t.Fatalf("FilterMovesByCategory(physical) len = %d, want 2", len(got))
	}
	got = FilterMovesByCategory(moves, "status")
	if len(got) != 1 {
		t.Fatalf("FilterMovesByCategory(status) len = %d, want 1", len(got))
	}
}

func TestSearchMoves(t *testing.T) {
	moves := []Move{
		{Name: "flamethrower"},
		{Name: "flame-wheel"},
		{Name: "surf"},
	}

	got := SearchMoves(moves, "flame")
	if len(got) != 2 {
		t.Fatalf("SearchMoves(flame) len = %d, want 2", len(got))
	}
	got = SearchMoves(moves, "")
	if len(got) != 3 {
		t.Fatalf("SearchMoves('') len = %d, want 3", len(got))
	}
}

func TestSearchAbilities(t *testing.T) {
	abilities := []Ability{
		{Name: "overgrow"},
		{Name: "over-coat"},
		{Name: "blaze"},
	}

	got := SearchAbilities(abilities, "over")
	if len(got) != 2 {
		t.Fatalf("SearchAbilities(over) len = %d, want 2", len(got))
	}
	got = SearchAbilities(abilities, "")
	if len(got) != 3 {
		t.Fatalf("SearchAbilities('') len = %d, want 3", len(got))
	}
}

func TestGetStat(t *testing.T) {
	p := Pokemon{
		Stats: []Stat{
			{Name: "hp", BaseStat: 45},
			{Name: "attack", BaseStat: 49},
			{Name: "speed", BaseStat: 90},
		},
	}

	if got := GetStat(p, "hp"); got != 45 {
		t.Errorf("GetStat(hp) = %d, want 45", got)
	}
	if got := GetStat(p, "Speed"); got != 90 {
		t.Errorf("GetStat(Speed) = %d, want 90", got)
	}
	if got := GetStat(p, "defense"); got != -1 {
		t.Errorf("GetStat(defense) = %d, want -1", got)
	}
}
