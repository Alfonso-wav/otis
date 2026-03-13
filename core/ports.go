package core

type PokemonFetcher interface {
	FetchPokemon(name string) (Pokemon, error)
	FetchPokemonList(offset int, limit int) (PokemonListResponse, error)
	FetchTypeList() (TypeListResponse, error)
	FetchType(name string) (PokemonTypeDetail, error)
	FetchRegions() ([]Region, error)
	FetchRegion(name string) (Region, error)
	FetchMove(name string) (Move, error)
	FetchAbility(name string) (Ability, error)
	FetchEvolutionChain(id int) (EvolutionChain, error)
}
