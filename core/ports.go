package core

// PokemonDBScraper define la interfaz para extraer datos de pokemondb.net.
type PokemonDBScraper interface {
	FetchPokedex() ([]PokedexDBEntry, error)
}

type PokemonFetcher interface {
	// --- Pokémon base ---
	FetchPokemon(name string) (Pokemon, error)
	FetchPokemonList(offset int, limit int) (PokemonListResponse, error)

	// --- Grupo A: Pokémon extendido ---
	FetchPokemonSpecies(name string) (PokemonSpecies, error)
	FetchPokemonForm(name string) (PokemonForm, error)
	FetchPokemonColors() ([]NamedResource, error)
	FetchPokemonShapes() ([]NamedResource, error)
	FetchPokemonHabitats() ([]NamedResource, error)

	// --- Tipos ---
	FetchTypeList() (TypeListResponse, error)
	FetchType(name string) (PokemonTypeDetail, error)

	// --- Regiones ---
	FetchRegions() ([]Region, error)
	FetchRegion(name string) (Region, error)

	// --- Grupo B: Naturalezas y cría ---
	FetchNatureList() ([]NamedResource, error)
	FetchNatureDetail(name string) (NatureDetail, error)
	FetchEggGroup(name string) (EggGroup, error)
	FetchGender(name string) (Gender, error)
	FetchGrowthRate(name string) (GrowthRate, error)

	// --- Grupo C: Movimientos ---
	FetchMoveList(offset int, limit int) (MoveListResponse, error)
	FetchMove(name string) (Move, error)
	FetchMoveDamageClass(name string) (MoveDamageClass, error)
	FetchMoveAilment(name string) (MoveAilment, error)
	FetchMoveTarget(name string) (MoveTarget, error)
	FetchMachine(id int) (Machine, error)

	// --- Grupo D: Habilidades ---
	FetchAbilityList(offset int, limit int) (AbilityListResponse, error)
	FetchAbility(name string) (Ability, error)

	// --- Evolución ---
	FetchEvolutionChain(id int) (EvolutionChain, error)

	// --- Grupo F: Ubicaciones ---
	FetchLocation(name string) (LocationDetail, error)
	FetchLocationArea(name string) (LocationArea, error)

	// --- Grupo G: Stats y generaciones ---
	FetchStat(name string) (StatDetail, error)
	FetchGenerations() ([]NamedResource, error)
	FetchGeneration(name string) (Generation, error)
	FetchPokedexList() ([]NamedResource, error)
	FetchPokedex(name string) (Pokedex, error)
	FetchVersionGroups() ([]NamedResource, error)
	FetchVersionGroup(name string) (VersionGroup, error)
}
