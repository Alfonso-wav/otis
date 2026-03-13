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

// GetNatures retorna la lista de todas las naturalezas disponibles.
func (a *App) GetNatures() []core.Nature {
	natures := make([]core.Nature, 0, len(core.Natures))
	for _, n := range core.Natures {
		natures = append(natures, n)
	}
	return natures
}

// CalculateEVs calcula los EVs estimados a partir de los stats actuales.
func (a *App) CalculateEVs(input core.EVCalculatorInput) (core.EVCalculatorResult, error) {
	pokemon, err := a.fetcher.FetchPokemon(core.NormalizeName(input.PokemonName))
	if err != nil {
		return core.EVCalculatorResult{}, err
	}

	baseStats := core.PokemonToBaseStats(pokemon)
	nature, ok := core.Natures[input.NatureName]
	if !ok {
		nature = core.Natures["Hardy"] // Neutral por defecto
	}

	ivs := core.DefaultIVs()
	if input.KnownIVs != nil {
		ivs = *input.KnownIVs
	}

	// Calcular rangos de EVs para cada stat
	evRanges := make(map[string]core.StatRange)
	evRanges["hp"] = core.EstimateEVRangeFromHP(input.CurrentStats.HP, baseStats.HP, ivs.HP, input.Level)
	evRanges["attack"] = core.EstimateEVRangeFromStat(input.CurrentStats.Attack, baseStats.Attack, ivs.Attack, input.Level, core.GetNatureModifier(nature, "attack"))
	evRanges["defense"] = core.EstimateEVRangeFromStat(input.CurrentStats.Defense, baseStats.Defense, ivs.Defense, input.Level, core.GetNatureModifier(nature, "defense"))
	evRanges["spAttack"] = core.EstimateEVRangeFromStat(input.CurrentStats.SpAttack, baseStats.SpAttack, ivs.SpAttack, input.Level, core.GetNatureModifier(nature, "spAttack"))
	evRanges["spDefense"] = core.EstimateEVRangeFromStat(input.CurrentStats.SpDefense, baseStats.SpDefense, ivs.SpDefense, input.Level, core.GetNatureModifier(nature, "spDefense"))
	evRanges["speed"] = core.EstimateEVRangeFromStat(input.CurrentStats.Speed, baseStats.Speed, ivs.Speed, input.Level, core.GetNatureModifier(nature, "speed"))

	// Usar el mínimo del rango como estimación
	estimatedEVs := core.Stats{
		HP:        evRanges["hp"].Min,
		Attack:    evRanges["attack"].Min,
		Defense:   evRanges["defense"].Min,
		SpAttack:  evRanges["spAttack"].Min,
		SpDefense: evRanges["spDefense"].Min,
		Speed:     evRanges["speed"].Min,
	}

	totalUsed := core.TotalEVs(estimatedEVs)

	// Calcular stats máximos posibles (252 EVs, 31 IVs)
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
