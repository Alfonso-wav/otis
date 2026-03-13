package core

type Pokemon struct {
	ID      int
	Name    string
	Types   []PokemonType
	Stats   []Stat
	Sprites Sprites
	Height  int
	Weight  int
}

type PokemonType struct {
	Name string
}

type Stat struct {
	Name     string
	BaseStat int
}

type Sprites struct {
	FrontDefault string
	FrontShiny   string
}

type PokemonListItem struct {
	Name string
	URL  string
}

type PokemonListResponse struct {
	Count    int
	Next     string
	Previous string
	Results  []PokemonListItem
}
