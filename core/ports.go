package core

type PokemonFetcher interface {
	FetchPokemon(name string) (Pokemon, error)
	FetchPokemonList(offset int, limit int) (PokemonListResponse, error)
}
