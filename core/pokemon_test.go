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

func TestFilterPokedexByType(t *testing.T) {
	pokedex := []string{"pikachu", "charmander", "bulbasaur", "squirtle", "gastly"}
	typePokemon := []string{"charmander", "vulpix", "growlithe", "gastly"}

	got := FilterPokedexByType(pokedex, typePokemon)
	if len(got) != 2 {
		t.Fatalf("FilterPokedexByType partial: len = %d, want 2", len(got))
	}

	// Empty intersection
	got = FilterPokedexByType(pokedex, []string{"mewtwo", "mew"})
	if len(got) != 0 {
		t.Fatalf("FilterPokedexByType empty: len = %d, want 0", len(got))
	}

	// Full intersection
	got = FilterPokedexByType([]string{"pikachu"}, []string{"pikachu", "raichu"})
	if len(got) != 1 {
		t.Fatalf("FilterPokedexByType full: len = %d, want 1", len(got))
	}

	// Case insensitive
	got = FilterPokedexByType([]string{"Pikachu"}, []string{"pikachu"})
	if len(got) != 1 {
		t.Fatalf("FilterPokedexByType case: len = %d, want 1", len(got))
	}
}

func TestFilterByClassification(t *testing.T) {
	items := []PokemonListItem{
		{Name: "mewtwo", URL: ""},
		{Name: "pikachu", URL: ""},
		{Name: "mew", URL: ""},
		{Name: "rayquaza", URL: ""},
		{Name: "bulbasaur", URL: ""},
	}
	classifications := map[string]SpeciesClassification{
		"mewtwo":    {IsLegendary: true, IsMythical: false},
		"pikachu":   {IsLegendary: false, IsMythical: false},
		"mew":       {IsLegendary: false, IsMythical: true},
		"rayquaza":  {IsLegendary: true, IsMythical: false},
		"bulbasaur": {IsLegendary: false, IsMythical: false},
	}

	// Filter legendary only
	got := FilterByClassification(items, classifications, true, false)
	if len(got) != 2 {
		t.Fatalf("legendary only: len = %d, want 2", len(got))
	}
	if got[0].Name != "mewtwo" || got[1].Name != "rayquaza" {
		t.Errorf("legendary only: got %v", got)
	}

	// Filter mythical only
	got = FilterByClassification(items, classifications, false, true)
	if len(got) != 1 {
		t.Fatalf("mythical only: len = %d, want 1", len(got))
	}
	if got[0].Name != "mew" {
		t.Errorf("mythical only: got %v", got)
	}

	// Filter both legendary and mythical
	got = FilterByClassification(items, classifications, true, true)
	if len(got) != 3 {
		t.Fatalf("legendary+mythical: len = %d, want 3", len(got))
	}

	// Neither flag set — returns empty
	got = FilterByClassification(items, classifications, false, false)
	if len(got) != 0 {
		t.Fatalf("no filter: len = %d, want 0", len(got))
	}

	// Pokemon not in classifications map is skipped
	extra := []PokemonListItem{{Name: "unknown-mon", URL: ""}}
	got = FilterByClassification(extra, classifications, true, false)
	if len(got) != 0 {
		t.Fatalf("unknown pokemon: len = %d, want 0", len(got))
	}
}
