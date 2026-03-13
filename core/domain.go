package core

type Region struct {
	Name      string
	Locations []Location
}

type Location struct {
	Name   string
	Region string
}

type Move struct {
	Name        string
	Type        string
	Power       int
	Accuracy    int
	PP          int
	Category    string
	Description string
}

type Ability struct {
	Name        string
	Description string
	Pokemon     []string
}

type EvolutionStage struct {
	Name            string
	MinLevel        int
	TriggerName     string
	EvolvesTo       []EvolutionStage
}

type EvolutionChain struct {
	ID    int
	Chain EvolutionStage
}

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

type TypePokemonEntry struct {
	Name string
	URL  string
}

type PokemonTypeDetail struct {
	Name    string
	Pokemon []TypePokemonEntry
}

type TypeListResponse struct {
	Count   int
	Results []PokemonListItem
}
