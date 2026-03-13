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

// Stats representa los 6 stats de un Pokémon
type Stats struct {
	HP        int `json:"hp"`
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	SpAttack  int `json:"spAttack"`
	SpDefense int `json:"spDefense"`
	Speed     int `json:"speed"`
}

// Nature representa una naturaleza y sus modificadores
type Nature struct {
	Name         string  `json:"name"`
	IncreasedStat string `json:"increasedStat"` // "" si neutral
	DecreasedStat string `json:"decreasedStat"` // "" si neutral
}

// EVCalculatorInput es el input del usuario para calcular EVs
type EVCalculatorInput struct {
	PokemonName  string `json:"pokemonName"`
	Level        int    `json:"level"`
	NatureName   string `json:"natureName"`
	CurrentStats Stats  `json:"currentStats"`
	KnownIVs     *Stats `json:"knownIVs,omitempty"` // nil = asumir 31
}

// StatRange representa un rango min-max para un stat
type StatRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// EVCalculatorResult es el resultado del cálculo de EVs
type EVCalculatorResult struct {
	Pokemon         string              `json:"pokemon"`
	Level           int                 `json:"level"`
	Nature          string              `json:"nature"`
	BaseStats       Stats               `json:"baseStats"`
	EstimatedEVs    Stats               `json:"estimatedEVs"`
	EVRanges        map[string]StatRange `json:"evRanges"`
	TotalEVsUsed    int                 `json:"totalEVsUsed"`
	EVsRemaining    int                 `json:"evsRemaining"`
	MaxPossibleStats Stats              `json:"maxPossibleStats"`
	UsedIVs         Stats               `json:"usedIVs"`
}

// StatCalculatorInput para calcular stats finales
type StatCalculatorInput struct {
	PokemonName string `json:"pokemonName"`
	Level       int    `json:"level"`
	NatureName  string `json:"natureName"`
	IVs         Stats  `json:"ivs"`
	EVs         Stats  `json:"evs"`
}
