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
