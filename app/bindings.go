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

// GetNatures retorna la lista de todas las naturalezas disponibles.
func (a *App) GetNatures() []core.Nature {
	natures := make([]core.Nature, 0, len(core.Natures))
	for _, n := range core.Natures {
		natures = append(natures, n)
	}
	return natures
}

// CalculateEVs calcula los EVs estimados a partir de los stats actuales del Pokémon.
func (a *App) CalculateEVs(input core.EVCalculatorInput) (core.EVCalculatorResult, error) {
	pokemon, err := a.fetcher.FetchPokemon(core.NormalizeName(input.PokemonName))
	if err != nil {
		return core.EVCalculatorResult{}, err
	}

	baseStats := core.PokemonToBaseStats(pokemon)
	nature, ok := core.Natures[input.NatureName]
	if !ok {
		nature = core.Natures["Hardy"]
	}

	ivs := core.DefaultIVs()
	if input.KnownIVs != nil {
		ivs = *input.KnownIVs
	}

	evRanges := map[string]core.StatRange{
		"hp":        core.EstimateEVRangeFromHP(input.CurrentStats.HP, baseStats.HP, ivs.HP, input.Level),
		"attack":    core.EstimateEVRangeFromStat(input.CurrentStats.Attack, baseStats.Attack, ivs.Attack, input.Level, core.GetNatureModifier(nature, "attack")),
		"defense":   core.EstimateEVRangeFromStat(input.CurrentStats.Defense, baseStats.Defense, ivs.Defense, input.Level, core.GetNatureModifier(nature, "defense")),
		"spAttack":  core.EstimateEVRangeFromStat(input.CurrentStats.SpAttack, baseStats.SpAttack, ivs.SpAttack, input.Level, core.GetNatureModifier(nature, "spAttack")),
		"spDefense": core.EstimateEVRangeFromStat(input.CurrentStats.SpDefense, baseStats.SpDefense, ivs.SpDefense, input.Level, core.GetNatureModifier(nature, "spDefense")),
		"speed":     core.EstimateEVRangeFromStat(input.CurrentStats.Speed, baseStats.Speed, ivs.Speed, input.Level, core.GetNatureModifier(nature, "speed")),
	}

	estimatedEVs := core.Stats{
		HP:        evRanges["hp"].Min,
		Attack:    evRanges["attack"].Min,
		Defense:   evRanges["defense"].Min,
		SpAttack:  evRanges["spAttack"].Min,
		SpDefense: evRanges["spDefense"].Min,
		Speed:     evRanges["speed"].Min,
	}

	totalUsed := core.TotalEVs(estimatedEVs)
	maxEVs := core.Stats{HP: 252, Attack: 252, Defense: 252, SpAttack: 252, SpDefense: 252, Speed: 252}
	maxStats := core.CalculateAllStats(baseStats, ivs, maxEVs, input.Level, nature)

	return core.EVCalculatorResult{
		Pokemon:          pokemon.Name,
		Level:            input.Level,
		Nature:           nature.Name,
		BaseStats:        baseStats,
		EstimatedEVs:     estimatedEVs,
		EVRanges:         evRanges,
		TotalEVsUsed:     totalUsed,
		EVsRemaining:     510 - totalUsed,
		MaxPossibleStats: maxStats,
		UsedIVs:          ivs,
	}, nil
}

// CalculateStats calcula los stats finales dados IVs, EVs, nivel y naturaleza.
func (a *App) CalculateStats(input core.StatCalculatorInput) (core.Stats, error) {
	pokemon, err := a.fetcher.FetchPokemon(core.NormalizeName(input.PokemonName))
	if err != nil {
		return core.Stats{}, err
	}

	baseStats := core.PokemonToBaseStats(pokemon)
	nature, ok := core.Natures[input.NatureName]
	if !ok {
		nature = core.Natures["Hardy"]
	}

	return core.CalculateAllStats(baseStats, input.IVs, input.EVs, input.Level, nature), nil
}
