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
