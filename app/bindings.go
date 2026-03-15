package app

import (
	"context"
	"fmt"
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

// GetRegionPokemonByType retorna los nombres de Pokémon de un tipo dado dentro de una región.
func (a *App) GetRegionPokemonByType(region, typeName string) ([]string, error) {
	pokedexMap := map[string]string{
		"kanto":  "kanto",
		"johto":  "original-johto",
		"hoenn":  "hoenn",
		"sinnoh": "original-sinnoh",
		"unova":  "original-unova",
		"kalos":  "kalos-central",
	}
	pokedexName, ok := pokedexMap[core.NormalizeName(region)]
	if !ok {
		pokedexName = core.NormalizeName(region)
	}

	pokedex, err := a.fetcher.FetchPokedex(pokedexName)
	if err != nil {
		return nil, err
	}

	typeDetail, err := a.fetcher.FetchType(core.NormalizeName(typeName))
	if err != nil {
		return nil, err
	}

	pokedexNames := make([]string, len(pokedex.PokemonEntries))
	for i, e := range pokedex.PokemonEntries {
		pokedexNames[i] = e.Pokemon
	}

	typePokemonNames := make([]string, len(typeDetail.Pokemon))
	for i, p := range typeDetail.Pokemon {
		typePokemonNames[i] = p.Name
	}

	return core.FilterPokedexByType(pokedexNames, typePokemonNames), nil
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

// SimulateTeamBattle runs a full team-vs-team battle simulation.
func (a *App) SimulateTeamBattle(team1Name, team2Name string) (core.TeamBattleState, error) {
	members1, err := a.resolveTeamBattleMembers(team1Name)
	if err != nil {
		return core.TeamBattleState{}, fmt.Errorf("team1: %w", err)
	}
	members2, err := a.resolveTeamBattleMembers(team2Name)
	if err != nil {
		return core.TeamBattleState{}, fmt.Errorf("team2: %w", err)
	}
	input := core.TeamBattleInput{
		Team1Name: team1Name, Team1Members: members1,
		Team2Name: team2Name, Team2Members: members2,
	}
	return core.SimulateTeamBattle(input, func(n int) int { return rand.Intn(n) }), nil
}

// SimulateMultipleTeamBattles runs N team battle simulations and returns aggregated statistics.
func (a *App) SimulateMultipleTeamBattles(team1Name, team2Name string, n int) (core.TeamBattleReport, error) {
	if n < 1 || n > 10000 {
		return core.TeamBattleReport{}, fmt.Errorf("n must be between 1 and 10000, got %d", n)
	}
	members1, err := a.resolveTeamBattleMembers(team1Name)
	if err != nil {
		return core.TeamBattleReport{}, fmt.Errorf("team1: %w", err)
	}
	members2, err := a.resolveTeamBattleMembers(team2Name)
	if err != nil {
		return core.TeamBattleReport{}, fmt.Errorf("team2: %w", err)
	}
	input := core.TeamBattleInput{
		Team1Name: team1Name, Team1Members: members1,
		Team2Name: team2Name, Team2Members: members2,
	}
	return core.SimulateMultipleTeamBattles(input, n, func(v int) int { return rand.Intn(v) }), nil
}

// resolveTeamBattleMembers loads a team and resolves each member's stats, types, and moves.
func (a *App) resolveTeamBattleMembers(teamName string) ([]core.TeamBattleMember, error) {
	team, err := a.teams.GetTeam(teamName)
	if err != nil {
		return nil, err
	}
	members := make([]core.TeamBattleMember, 0, len(team.Members))
	for _, m := range team.Members {
		pokemon, err := a.fetcher.FetchPokemon(core.NormalizeName(m.PokemonName))
		if err != nil {
			return nil, fmt.Errorf("fetch %s: %w", m.PokemonName, err)
		}
		baseStats := core.PokemonToBaseStats(pokemon)
		nature, ok := core.Natures[m.Nature]
		if !ok {
			nature = core.Natures["Hardy"]
		}
		stats := core.CalculateAllStats(baseStats, m.IVs, m.EVs, m.Level, nature)

		var moves []core.Move
		for _, moveName := range m.Moves {
			move, err := a.fetcher.FetchMove(core.NormalizeName(moveName))
			if err == nil {
				moves = append(moves, move)
			}
		}
		if len(moves) == 0 {
			moves = []core.Move{{Name: "struggle", Type: "normal", Power: 50, Category: "physical", Accuracy: 100}}
		}

		members = append(members, core.TeamBattleMember{
			PokemonName: m.PokemonName,
			Stats:       stats,
			Types:       pokemon.Types,
			Moves:       moves,
			Level:       m.Level,
		})
	}
	return members, nil
}

// SimulateFullBattle runs a complete automatic battle and returns the final state.
func (a *App) SimulateFullBattle(input core.FullBattleInput) core.BattleState {
	return core.SimulateFullBattle(input, func(n int) int {
		return rand.Intn(n)
	})
}

// SimulateMultipleBattles runs N full battle simulations and returns aggregated statistics.
func (a *App) SimulateMultipleBattles(input core.FullBattleInput, n int) (core.BattleReport, error) {
	if n < 1 || n > 10000 {
		return core.BattleReport{}, fmt.Errorf("n must be between 1 and 10000, got %d", n)
	}
	return core.SimulateMultipleBattles(input, n, func(v int) int {
		return rand.Intn(v)
	}), nil
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

// GetAllMoves retorna todos los movimientos de PokeAPI (~920) con cache en memoria.
func (a *App) GetAllMoves() ([]core.Move, error) {
	return a.fetcher.FetchAllMoves()
}

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

// GetAllAbilities retorna todas las habilidades de PokeAPI (~300) con cache en memoria.
func (a *App) GetAllAbilities() ([]core.Ability, error) {
	return a.fetcher.FetchAllAbilities()
}

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

// FillTeamRandom fills empty slots of a team with random Pokemon.
func (a *App) FillTeamRandom(teamName string) (core.Team, error) {
	team, err := a.teams.GetTeam(teamName)
	if err != nil {
		return core.Team{}, err
	}
	if len(team.Members) >= core.MaxTeamMembers {
		return team, nil
	}

	slotsNeeded := core.MaxTeamMembers - len(team.Members)
	list, err := a.fetcher.FetchPokemonList(0, 151)
	if err != nil {
		return core.Team{}, err
	}

	var pokemon []core.Pokemon
	for _, item := range list.Results {
		if len(pokemon) >= slotsNeeded+len(team.Members) {
			break
		}
		p, ferr := a.fetcher.FetchPokemon(item.Name)
		if ferr == nil {
			pokemon = append(pokemon, p)
		}
	}

	filled := core.FillTeamRandom(team, pokemon, func(n int) int {
		return rand.Intn(n)
	})
	if err := a.teams.SaveTeam(filled); err != nil {
		return core.Team{}, err
	}
	return filled, nil
}

// CreateTeam creates a new empty team with the given name.
func (a *App) CreateTeam(name string) error {
	team := core.Team{Name: name, Members: []core.TeamMember{}}
	if err := core.ValidateTeam(team); err != nil {
		return err
	}
	if _, err := a.teams.GetTeam(name); err == nil {
		return fmt.Errorf("team %q already exists", name)
	}
	return a.teams.SaveTeam(team)
}

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

// UpdateTeamMember updates a member of a team at the given index.
func (a *App) UpdateTeamMember(teamName string, memberIndex int, member core.TeamMember) error {
	team, err := a.teams.GetTeam(teamName)
	if err != nil {
		return err
	}
	team, err = core.UpdateTeamMember(team, memberIndex, member)
	if err != nil {
		return err
	}
	return a.teams.SaveTeam(team)
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
