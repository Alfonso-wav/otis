package app

import (
	"context"
	"math/rand"

	"github.com/alfon/pokemon-app/core"
)

// App es el struct bindeado a Wails. Sus métodos exportados son invocables
// desde el frontend via IPC (window.go.app.App.Method()).
type App struct {
	ctx     context.Context
	fetcher core.PokemonFetcher
	scraper core.PokemonDBScraper
	teams   core.TeamStorage
	sprites core.SpriteDownloader
}

// NewApp crea una instancia de App con el fetcher, scraper, team storage y sprite downloader inyectados desde main.
func NewApp(fetcher core.PokemonFetcher, scraper core.PokemonDBScraper, teams core.TeamStorage, sprites core.SpriteDownloader) *App {
	return &App{fetcher: fetcher, scraper: scraper, teams: teams, sprites: sprites}
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

// SimulateDamage calcula el daño de un movimiento dado atacante y defensor.
func (a *App) SimulateDamage(input core.DamageInput) (core.DamageResult, error) {
	return core.CalculateDamage(input), nil
}

// ComparePokemons fetcha ambos Pokémon y devuelve la comparación de sus stats base.
func (a *App) ComparePokemons(nameA, nameB string) (core.PokemonComparison, error) {
	pkA, err := a.fetcher.FetchPokemon(core.NormalizeName(nameA))
	if err != nil {
		return core.PokemonComparison{}, err
	}
	pkB, err := a.fetcher.FetchPokemon(core.NormalizeName(nameB))
	if err != nil {
		return core.PokemonComparison{}, err
	}
	return core.ComparePokemons(pkA, pkB), nil
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

// InitBattle initializes a fresh battle state with the given max HP values.
func (a *App) InitBattle(attackerMaxHP, defenderMaxHP int) core.BattleState {
	return core.InitBattle(attackerMaxHP, defenderMaxHP)
}

// ExecuteTurn applies one turn of battle and returns the new state and damage result.
func (a *App) ExecuteTurn(input core.TurnInput) core.TurnResult {
	return core.ExecuteTurn(input, func(n int) int {
		return rand.Intn(n)
	})
}

// SimulateFullBattle runs a complete automatic battle and returns the final state.
func (a *App) SimulateFullBattle(input core.FullBattleInput) core.BattleState {
	return core.SimulateFullBattle(input, func(n int) int {
		return rand.Intn(n)
	})
}

// --- Grupo A: Pokémon extendido ---

func (a *App) GetPokemonSpecies(name string) (core.PokemonSpecies, error) {
	return a.fetcher.FetchPokemonSpecies(core.NormalizeName(name))
}

func (a *App) GetPokemonForm(name string) (core.PokemonForm, error) {
	return a.fetcher.FetchPokemonForm(core.NormalizeName(name))
}

func (a *App) ListPokemonColors() ([]core.NamedResource, error) {
	return a.fetcher.FetchPokemonColors()
}

func (a *App) ListPokemonShapes() ([]core.NamedResource, error) {
	return a.fetcher.FetchPokemonShapes()
}

func (a *App) ListPokemonHabitats() ([]core.NamedResource, error) {
	return a.fetcher.FetchPokemonHabitats()
}

// --- Grupo B: Naturalezas y cría ---

func (a *App) ListNatureNames() ([]core.NamedResource, error) {
	return a.fetcher.FetchNatureList()
}

func (a *App) GetNatureDetail(name string) (core.NatureDetail, error) {
	return a.fetcher.FetchNatureDetail(core.NormalizeName(name))
}

func (a *App) GetEggGroup(name string) (core.EggGroup, error) {
	return a.fetcher.FetchEggGroup(core.NormalizeName(name))
}

func (a *App) GetGender(name string) (core.Gender, error) {
	return a.fetcher.FetchGender(core.NormalizeName(name))
}

func (a *App) GetGrowthRate(name string) (core.GrowthRate, error) {
	return a.fetcher.FetchGrowthRate(core.NormalizeName(name))
}

// --- Grupo C: Movimientos completos ---

func (a *App) ListMoves(offset int, limit int) (core.MoveListResponse, error) {
	return a.fetcher.FetchMoveList(offset, limit)
}

func (a *App) GetMoveDamageClass(name string) (core.MoveDamageClass, error) {
	return a.fetcher.FetchMoveDamageClass(core.NormalizeName(name))
}

func (a *App) GetMoveAilment(name string) (core.MoveAilment, error) {
	return a.fetcher.FetchMoveAilment(core.NormalizeName(name))
}

func (a *App) GetMoveTarget(name string) (core.MoveTarget, error) {
	return a.fetcher.FetchMoveTarget(core.NormalizeName(name))
}

func (a *App) GetMachine(id int) (core.Machine, error) {
	return a.fetcher.FetchMachine(id)
}

// --- Grupo D: Habilidades completas ---

func (a *App) ListAbilities(offset int, limit int) (core.AbilityListResponse, error) {
	return a.fetcher.FetchAbilityList(offset, limit)
}

// --- Grupo F: Ubicaciones ---

func (a *App) GetLocation(name string) (core.LocationDetail, error) {
	return a.fetcher.FetchLocation(core.NormalizeName(name))
}

func (a *App) GetLocationArea(name string) (core.LocationArea, error) {
	return a.fetcher.FetchLocationArea(core.NormalizeName(name))
}

// --- Grupo G: Stats y generaciones ---

func (a *App) GetStatDetail(name string) (core.StatDetail, error) {
	return a.fetcher.FetchStat(core.NormalizeName(name))
}

func (a *App) ListGenerations() ([]core.NamedResource, error) {
	return a.fetcher.FetchGenerations()
}

func (a *App) GetGeneration(name string) (core.Generation, error) {
	return a.fetcher.FetchGeneration(core.NormalizeName(name))
}

func (a *App) ListPokedexes() ([]core.NamedResource, error) {
	return a.fetcher.FetchPokedexList()
}

func (a *App) GetPokedex(name string) (core.Pokedex, error) {
	return a.fetcher.FetchPokedex(core.NormalizeName(name))
}

func (a *App) ListVersionGroups() ([]core.NamedResource, error) {
	return a.fetcher.FetchVersionGroups()
}

func (a *App) GetVersionGroup(name string) (core.VersionGroup, error) {
	return a.fetcher.FetchVersionGroup(core.NormalizeName(name))
}

// --- PokemonDB Scraper ---

// ScrapePokedex extrae la tabla completa de Pokémon desde pokemondb.net.
func (a *App) ScrapePokedex() ([]core.PokedexDBEntry, error) {
	return a.scraper.FetchPokedex()
}

// --- Sprites ---

// DownloadSprites descarga sprites de pokemondb.net a disco local.
func (a *App) DownloadSprites(categories []core.SpriteCategory) (core.SpriteDownloadResult, error) {
	return a.sprites.DownloadAllSprites("assets/sprites", categories)
}

// --- Teams ---

// SaveToTeam agrega un miembro a un equipo existente o crea uno nuevo.
func (a *App) SaveToTeam(teamName string, member core.TeamMember) error {
	team, err := a.teams.GetTeam(teamName)
	if err != nil {
		team = core.Team{Name: teamName}
	}
	team, err = core.AddMemberToTeam(team, member)
	if err != nil {
		return err
	}
	return a.teams.SaveTeam(team)
}

// ListTeams retorna todos los equipos guardados.
func (a *App) ListTeams() ([]core.Team, error) {
	return a.teams.ListTeams()
}

// GetTeam retorna un equipo por nombre.
func (a *App) GetTeam(name string) (core.Team, error) {
	return a.teams.GetTeam(name)
}

// DeleteTeam elimina un equipo completo.
func (a *App) DeleteTeam(name string) error {
	return a.teams.DeleteTeam(name)
}

// DeleteTeamMember elimina un miembro de un equipo por indice.
func (a *App) DeleteTeamMember(teamName string, memberIndex int) error {
	team, err := a.teams.GetTeam(teamName)
	if err != nil {
		return err
	}
	team, err = core.RemoveMemberFromTeam(team, memberIndex)
	if err != nil {
		return err
	}
	return a.teams.SaveTeam(team)
}
