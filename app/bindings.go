package app

import (
	"context"

	"github.com/alfon/pokemon-app/core"
)

// App es el struct bindeado a Wails. Sus métodos exportados son invocables
// desde el frontend via IPC (window.go.app.App.Method()).
type App struct {
	ctx     context.Context
	fetcher core.PokemonFetcher
}

// NewApp crea una instancia de App con el fetcher inyectado desde main.
func NewApp(fetcher core.PokemonFetcher) *App {
	return &App{fetcher: fetcher}
}

// Startup es llamado por Wails al arrancar la ventana.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// ListPokemon retorna la lista paginada de Pokémon.
func (a *App) ListPokemon(offset int, limit int) (core.PokemonListResponse, error) {
	return a.fetcher.FetchPokemonList(offset, limit)
}

// GetPokemon retorna el detalle de un Pokémon por nombre.
func (a *App) GetPokemon(name string) (core.Pokemon, error) {
	return a.fetcher.FetchPokemon(core.NormalizeName(name))
}

// ListTypes retorna la lista de todos los tipos de Pokémon.
func (a *App) ListTypes() (core.TypeListResponse, error) {
	return a.fetcher.FetchTypeList()
}

// GetType retorna el detalle de un tipo con sus Pokémon.
func (a *App) GetType(name string) (core.PokemonTypeDetail, error) {
	return a.fetcher.FetchType(core.NormalizeName(name))
}

// ListRegions retorna todas las regiones del mundo Pokémon.
func (a *App) ListRegions() ([]core.Region, error) {
	return a.fetcher.FetchRegions()
}

// GetRegion retorna el detalle de una región con sus locations.
func (a *App) GetRegion(name string) (core.Region, error) {
	return a.fetcher.FetchRegion(core.NormalizeName(name))
}

// GetMove retorna el detalle de un movimiento.
func (a *App) GetMove(name string) (core.Move, error) {
	return a.fetcher.FetchMove(core.NormalizeName(name))
}

// GetAbility retorna el detalle de una habilidad.
func (a *App) GetAbility(name string) (core.Ability, error) {
	return a.fetcher.FetchAbility(core.NormalizeName(name))
}

// GetEvolutionChain retorna la cadena evolutiva por ID.
func (a *App) GetEvolutionChain(id int) (core.EvolutionChain, error) {
	return a.fetcher.FetchEvolutionChain(id)
}
