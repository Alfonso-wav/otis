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
	Priority    int
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

// PokemonMoveEntry representa un movimiento que puede aprender un Pokémon
type PokemonMoveEntry struct {
	Name   string
	Method string // level-up, machine, egg, tutor
	Level  int    // nivel al que se aprende (0 si no aplica)
}

type Pokemon struct {
	ID      int
	Name    string
	Types   []PokemonType
	Stats   []Stat
	Sprites Sprites
	Height  int
	Weight  int
	Moves   []PokemonMoveEntry
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
	Name          string `json:"name"`
	IncreasedStat string `json:"increasedStat"`
	DecreasedStat string `json:"decreasedStat"`
}

// StatRange representa un rango min-max para un stat
type StatRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// EVCalculatorInput es el input del usuario para calcular EVs
type EVCalculatorInput struct {
	PokemonName  string `json:"pokemonName"`
	Level        int    `json:"level"`
	NatureName   string `json:"natureName"`
	CurrentStats Stats  `json:"currentStats"`
	KnownIVs     *Stats `json:"knownIVs,omitempty"`
}

// EVCalculatorResult es el resultado del cálculo de EVs
type EVCalculatorResult struct {
	Pokemon          string               `json:"pokemon"`
	Level            int                  `json:"level"`
	Nature           string               `json:"nature"`
	BaseStats        Stats                `json:"baseStats"`
	EstimatedEVs     Stats                `json:"estimatedEVs"`
	EVRanges         map[string]StatRange `json:"evRanges"`
	TotalEVsUsed     int                  `json:"totalEVsUsed"`
	EVsRemaining     int                  `json:"evsRemaining"`
	MaxPossibleStats Stats                `json:"maxPossibleStats"`
	UsedIVs          Stats                `json:"usedIVs"`
}

// StatCalculatorInput para calcular stats finales dados IVs, EVs y naturaleza
type StatCalculatorInput struct {
	PokemonName string `json:"pokemonName"`
	Level       int    `json:"level"`
	NatureName  string `json:"natureName"`
	IVs         Stats  `json:"ivs"`
	EVs         Stats  `json:"evs"`
}

// NamedResource es una referencia a un recurso de la API con nombre y URL.
type NamedResource struct {
	Name string
	URL  string
}

// --- Grupo A: Pokémon extendido ---

type PokemonSpecies struct {
	Name                 string
	Order                int
	GenderRate           int // -1 genderless, 0 always male, 8 always female
	CaptureRate          int
	BaseHappiness        int
	IsBaby               bool
	IsLegendary          bool
	IsMythical           bool
	HatchCounter         int
	HasGenderDifferences bool
	FormsSwitchable      bool
	Genus                string
	Color                string
	Shape                string
	Habitat              string
	EggGroups            []string
	FlavorText           string
	EvolutionChainID     int
	Varieties            []PokemonVariety
}

type PokemonVariety struct {
	IsDefault bool
	Pokemon   string
}

type PokemonForm struct {
	Name         string
	FormName     string
	IsMega       bool
	IsBattleOnly bool
	Types        []PokemonType
	Sprites      Sprites
}

// --- Grupo B: Naturalezas y cría ---

type NatureDetail struct {
	Name          string
	IncreasedStat string
	DecreasedStat string
	LikesFlavor   string
	HatesFlavor   string
}

type EggGroup struct {
	Name    string
	Pokemon []string
}

type Gender struct {
	Name    string
	Pokemon []string
}

type GrowthRateLevel struct {
	Level      int
	Experience int
}

type GrowthRate struct {
	Name    string
	Formula string
	Levels  []GrowthRateLevel
	Pokemon []string
}

// --- Grupo C: Movimientos completos ---

type MoveListResponse struct {
	Count   int
	Results []NamedResource
}

type MoveDamageClass struct {
	Name        string
	Description string
	Moves       []string
}

type MoveAilment struct {
	Name  string
	Moves []string
}

type MoveTarget struct {
	Name        string
	Description string
}

type Machine struct {
	ID           int
	Move         string
	Item         string
	VersionGroup string
}

// --- Grupo D: Habilidades ---

type AbilityListResponse struct {
	Count   int
	Results []NamedResource
}

// --- Grupo F: Ubicaciones extendidas ---

type LocationDetail struct {
	Name   string
	Region string
	Areas  []string
}

type PokemonEncounter struct {
	PokemonName string
	MaxChance   int
}

type LocationArea struct {
	Name              string
	Location          string
	PokemonEncounters []PokemonEncounter
}

// --- Grupo G: Stats y generaciones ---

type StatDetail struct {
	Name                 string
	IsBattleOnly         bool
	AffectingMovesBuff   []string
	AffectingMovesNerf   []string
	AffectingNaturesBuff []string
	AffectingNaturesNerf []string
}

type Generation struct {
	Name           string
	MainRegion     string
	Games          []string
	PokemonSpecies []PokemonListItem
	Types          []string
	Moves          []string
	Abilities      []string
}

type PokedexEntry struct {
	EntryNumber int
	Pokemon     string
}

type Pokedex struct {
	Name           string
	IsMainSeries   bool
	Region         string
	PokemonEntries []PokedexEntry
}

type VersionGroup struct {
	Name       string
	Order      int
	Generation string
	Versions   []string
	Pokedexes  []string
	Regions    []string
}

// --- PokemonDB Scraper ---

// PokedexDBEntry representa un Pokémon extraído de la tabla de pokemondb.net/pokedex/all.
type PokedexDBEntry struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Types   []string `json:"types"`
	Total   int      `json:"total"`
	HP      int      `json:"hp"`
	Attack  int      `json:"attack"`
	Defense int      `json:"defense"`
	SpAtk   int      `json:"spAtk"`
	SpDef   int      `json:"spDef"`
	Speed   int      `json:"speed"`
}
