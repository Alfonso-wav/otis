package core

import (
	"fmt"
	"sort"
)

// FullBattleInput contains everything needed to simulate a complete battle.
type FullBattleInput struct {
	AttackerStats Stats         `json:"attackerStats"`
	DefenderStats Stats         `json:"defenderStats"`
	AttackerTypes []PokemonType `json:"attackerTypes"`
	DefenderTypes []PokemonType `json:"defenderTypes"`
	AttackerLevel int           `json:"attackerLevel"`
	DefenderLevel int           `json:"defenderLevel"`
	AttackerMoves []Move        `json:"attackerMoves"`
	DefenderMoves []Move        `json:"defenderMoves"`
	AttackerName  string        `json:"attackerName"`
	DefenderName  string        `json:"defenderName"`
}

// resolveOrder determines who attacks first based on move priority, then speed.
// Returns true if attacker goes first.
func resolveOrder(atkSpeed, defSpeed int, atkMove, defMove Move, randSource func(n int) int) bool {
	if atkMove.Priority != defMove.Priority {
		return atkMove.Priority > defMove.Priority
	}
	if atkSpeed != defSpeed {
		return atkSpeed > defSpeed
	}
	return randSource(2) == 0
}

// SimulateFullBattle simulates an entire battle automatically.
// randSource is injected for testability; use rand.Intn in production.
// Each turn both sides pick a random move, then order is resolved by priority/speed.
// Max 200 turns to prevent infinite loops.
func SimulateFullBattle(input FullBattleInput, randSource func(n int) int) BattleState {
	if len(input.AttackerMoves) == 0 || len(input.DefenderMoves) == 0 {
		return BattleState{}
	}

	state := InitBattle(input.AttackerStats.HP, input.DefenderStats.HP)
	const maxTurns = 200

	for !state.IsOver && state.TurnCount < maxTurns {
		atkMove := input.AttackerMoves[randSource(len(input.AttackerMoves))]
		defMove := input.DefenderMoves[randSource(len(input.DefenderMoves))]

		attackerFirst := resolveOrder(
			input.AttackerStats.Speed, input.DefenderStats.Speed,
			atkMove, defMove, randSource,
		)

		if attackerFirst {
			state = executeAttackerTurn(state, input, atkMove, randSource)
			if !state.IsOver {
				state = executeDefenderTurn(state, input, defMove, randSource)
			}
		} else {
			state = executeDefenderTurn(state, input, defMove, randSource)
			if !state.IsOver {
				state = executeAttackerTurn(state, input, atkMove, randSource)
			}
		}
	}

	// Resolve draw by max turns: whoever has more HP wins.
	if !state.IsOver {
		state.IsOver = true
		switch {
		case state.AttackerHP > state.DefenderHP:
			state.Winner = "attacker"
		case state.DefenderHP > state.AttackerHP:
			state.Winner = "defender"
		default:
			state.Winner = "draw"
		}
	}

	return state
}

func executeAttackerTurn(state BattleState, input FullBattleInput, move Move, randSource func(n int) int) BattleState {
	result := ExecuteTurn(TurnInput{
		State:         state,
		AttackerStats: input.AttackerStats,
		DefenderStats: input.DefenderStats,
		AttackerTypes: input.AttackerTypes,
		DefenderTypes: input.DefenderTypes,
		AttackerLevel: input.AttackerLevel,
		DefenderLevel: input.DefenderLevel,
		Move:          move,
		AttackerName:  input.AttackerName,
	}, randSource)
	return result.NewState
}

func executeDefenderTurn(state BattleState, input FullBattleInput, move Move, randSource func(n int) int) BattleState {
	swapped := BattleState{
		AttackerHP:    state.DefenderHP,
		DefenderHP:    state.AttackerHP,
		AttackerMaxHP: state.DefenderMaxHP,
		DefenderMaxHP: state.AttackerMaxHP,
		TurnCount:     state.TurnCount,
		Log:           state.Log,
		IsOver:        state.IsOver,
		Winner:        state.Winner,
	}
	result := ExecuteTurn(TurnInput{
		State:         swapped,
		AttackerStats: input.DefenderStats,
		DefenderStats: input.AttackerStats,
		AttackerTypes: input.DefenderTypes,
		DefenderTypes: input.AttackerTypes,
		AttackerLevel: input.DefenderLevel,
		DefenderLevel: input.AttackerLevel,
		Move:          move,
		AttackerName:  input.DefenderName,
	}, randSource)
	ns := result.NewState
	state = BattleState{
		AttackerHP:    ns.DefenderHP,
		DefenderHP:    ns.AttackerHP,
		AttackerMaxHP: ns.DefenderMaxHP,
		DefenderMaxHP: ns.AttackerMaxHP,
		TurnCount:     ns.TurnCount,
		Log:           ns.Log,
		IsOver:        ns.IsOver,
		Winner:        ns.Winner,
	}
	if state.IsOver {
		state.Winner = "defender"
	}
	return state
}

// BattleState represents the current state of a turn-based battle.
type BattleState struct {
	AttackerHP    int      `json:"attackerHP"`
	DefenderHP    int      `json:"defenderHP"`
	AttackerMaxHP int      `json:"attackerMaxHP"`
	DefenderMaxHP int      `json:"defenderMaxHP"`
	TurnCount     int      `json:"turnCount"`
	Log           []string `json:"log"`
	IsOver        bool     `json:"isOver"`
	Winner        string   `json:"winner"` // "attacker" | "defender" | ""
}

// TurnInput contains everything needed to simulate one turn.
type TurnInput struct {
	State         BattleState   `json:"state"`
	AttackerStats Stats         `json:"attackerStats"`
	DefenderStats Stats         `json:"defenderStats"`
	AttackerTypes []PokemonType `json:"attackerTypes"`
	DefenderTypes []PokemonType `json:"defenderTypes"`
	AttackerLevel int           `json:"attackerLevel"`
	DefenderLevel int           `json:"defenderLevel"`
	Move          Move          `json:"move"`
	AttackerName  string        `json:"attackerName"`
}

// TurnResult is the outcome of one executed turn.
type TurnResult struct {
	NewState BattleState  `json:"newState"`
	Damage   DamageResult `json:"damage"`
	LogEntry string       `json:"logEntry"`
	Missed   bool         `json:"missed"`
}

// InitBattle returns a fresh BattleState with full HP for both sides.
func InitBattle(attackerMaxHP, defenderMaxHP int) BattleState {
	return BattleState{
		AttackerHP:    attackerMaxHP,
		DefenderHP:    defenderMaxHP,
		AttackerMaxHP: attackerMaxHP,
		DefenderMaxHP: defenderMaxHP,
		TurnCount:     0,
		Log:           []string{},
		IsOver:        false,
		Winner:        "",
	}
}

// CheckAccuracy determines whether a move hits. Accuracy 0 means never-miss.
// randSource(n) returns [0, n). Returns true if the move hits.
func CheckAccuracy(accuracy int, randSource func(n int) int) bool {
	if accuracy <= 0 || accuracy >= 100 {
		return true
	}
	return randSource(100) < accuracy
}

// ExecuteTurn is a pure function: given the current state and the attacker's chosen move,
// it returns the new state and the damage result. No side effects.
// If randSource is nil, damage is deterministic (average) with no accuracy/crit checks.
func ExecuteTurn(input TurnInput, randSource func(n int) int) TurnResult {
	state := input.State

	// No-op if battle is already over.
	if state.IsOver {
		return TurnResult{NewState: state}
	}

	state.TurnCount++

	// Accuracy check
	if randSource != nil && input.Move.Power > 0 && input.Move.Category != "status" {
		if !CheckAccuracy(input.Move.Accuracy, randSource) {
			logEntry := fmt.Sprintf("[T%d] %s usó %s → ¡Falló!",
				state.TurnCount, input.AttackerName, input.Move.Name)
			newLog := make([]string, len(state.Log)+1)
			copy(newLog, state.Log)
			newLog[len(state.Log)] = logEntry
			state.Log = newLog
			return TurnResult{NewState: state, LogEntry: logEntry, Missed: true}
		}
	}

	var dmg DamageResult
	if randSource != nil {
		dmg = CalculateBattleDamage(DamageInput{
			AttackerStats: input.AttackerStats,
			DefenderStats: input.DefenderStats,
			Move:          input.Move,
			AttackerTypes: input.AttackerTypes,
			DefenderTypes: input.DefenderTypes,
			Level:         input.AttackerLevel,
			IsCritical:    false,
			WeatherBonus:  1.0,
		}, randSource)
	} else {
		dmg = CalculateDamage(DamageInput{
			AttackerStats: input.AttackerStats,
			DefenderStats: input.DefenderStats,
			Move:          input.Move,
			AttackerTypes: input.AttackerTypes,
			DefenderTypes: input.DefenderTypes,
			Level:         input.AttackerLevel,
			IsCritical:    false,
			WeatherBonus:  1.0,
		})
	}

	var logEntry string

	if input.Move.Power <= 0 || input.Move.Category == "status" {
		logEntry = fmt.Sprintf("[T%d] %s usó %s → sin efecto de daño | HP Defensor: %d/%d",
			state.TurnCount, input.AttackerName, input.Move.Name, state.DefenderHP, state.DefenderMaxHP)
	} else {
		applied := dmg.ActualDamage
		if applied < 1 && !dmg.HasNoEffect {
			applied = 1
		}

		newHP := state.DefenderHP - applied
		if newHP < 0 {
			newHP = 0
		}
		state.DefenderHP = newHP

		effStr := fmt.Sprintf("×%.1f", dmg.Multiplier)
		if dmg.HasNoEffect {
			effStr = "Sin efecto (inmune)"
		} else if dmg.IsSuperEffective {
			effStr = fmt.Sprintf("¡Super eficaz! ×%.1f", dmg.Multiplier)
		} else if dmg.IsNotVeryEffective {
			effStr = fmt.Sprintf("Poco eficaz ×%.1f", dmg.Multiplier)
		}

		critStr := ""
		if dmg.WasCritical {
			critStr = "¡Golpe crítico! "
		}

		stabStr := ""
		if dmg.HasSTAB {
			stabStr = "(STAB) "
		}

		logEntry = fmt.Sprintf("%s%s[T%d] %s usó %s → %d daño (%s) | HP Defensor: %d/%d",
			critStr, stabStr, state.TurnCount, input.AttackerName, input.Move.Name, applied, effStr, state.DefenderHP, state.DefenderMaxHP)

		if state.DefenderHP <= 0 {
			state.IsOver = true
			state.Winner = "attacker"
		}
	}

	newLog := make([]string, len(state.Log)+1)
	copy(newLog, state.Log)
	newLog[len(state.Log)] = logEntry
	state.Log = newLog

	return TurnResult{
		NewState: state,
		Damage:   dmg,
		LogEntry: logEntry,
	}
}

// BattleSummary holds the key stats from a single battle simulation.
type BattleSummary struct {
	Winner      string `json:"winner"`
	Turns       int    `json:"turns"`
	AttackerHP  int    `json:"attackerHP"`
	DefenderHP  int    `json:"defenderHP"`
}

// BattleReport aggregates statistics from N battle simulations.
type BattleReport struct {
	TotalSimulations int     `json:"totalSimulations"`
	AttackerWins     int     `json:"attackerWins"`
	DefenderWins     int     `json:"defenderWins"`
	Draws            int     `json:"draws"`
	AttackerWinPct   float64 `json:"attackerWinPct"`
	DefenderWinPct   float64 `json:"defenderWinPct"`
	DrawPct          float64 `json:"drawPct"`
	AvgTurns         float64 `json:"avgTurns"`
	MinTurns         int     `json:"minTurns"`
	MaxTurns         int     `json:"maxTurns"`
	MedianTurns      int     `json:"medianTurns"`
	AvgWinnerHP      float64 `json:"avgWinnerHP"`
}

// SimulateTeamBattle simulates a full team-vs-team battle.
// When one Pokemon faints, the next on that team enters. The winner's HP carries over.
func SimulateTeamBattle(input TeamBattleInput, randSource func(int) int) TeamBattleState {
	if len(input.Team1Members) == 0 || len(input.Team2Members) == 0 {
		return TeamBattleState{IsOver: true, Winner: "draw"}
	}

	idx1, idx2 := 0, 0
	carryHP1, carryHP2 := 0, 0 // 0 means use full HP
	totalTurns := 0
	var rounds []BattleState
	var log []string

	for idx1 < len(input.Team1Members) && idx2 < len(input.Team2Members) {
		m1 := input.Team1Members[idx1]
		m2 := input.Team2Members[idx2]

		bi := FullBattleInput{
			AttackerStats: m1.Stats,
			DefenderStats: m2.Stats,
			AttackerTypes: m1.Types,
			DefenderTypes: m2.Types,
			AttackerLevel: m1.Level,
			DefenderLevel: m2.Level,
			AttackerMoves: m1.Moves,
			DefenderMoves: m2.Moves,
			AttackerName:  m1.PokemonName,
			DefenderName:  m2.PokemonName,
		}

		// Apply carry-over HP
		if carryHP1 > 0 {
			bi.AttackerStats.HP = carryHP1
		}
		if carryHP2 > 0 {
			bi.DefenderStats.HP = carryHP2
		}

		result := SimulateFullBattle(bi, randSource)
		rounds = append(rounds, result)
		totalTurns += result.TurnCount

		roundNum := len(rounds)
		if result.Winner == "attacker" {
			log = append(log, fmt.Sprintf("[Ronda %d] %s venció a %s (HP restante: %d)",
				roundNum, m1.PokemonName, m2.PokemonName, result.AttackerHP))
			carryHP1 = result.AttackerHP
			carryHP2 = 0
			idx2++
		} else {
			log = append(log, fmt.Sprintf("[Ronda %d] %s venció a %s (HP restante: %d)",
				roundNum, m2.PokemonName, m1.PokemonName, result.DefenderHP))
			carryHP2 = result.DefenderHP
			carryHP1 = 0
			idx1++
		}
	}

	winner := "team1"
	if idx1 >= len(input.Team1Members) {
		winner = "team2"
	}

	return TeamBattleState{
		Team1Remaining: len(input.Team1Members) - idx1,
		Team2Remaining: len(input.Team2Members) - idx2,
		TotalTurns:     totalTurns,
		Rounds:         rounds,
		Log:            log,
		IsOver:         true,
		Winner:         winner,
	}
}

// SimulateMultipleTeamBattles runs N team battle simulations and returns aggregated statistics.
func SimulateMultipleTeamBattles(input TeamBattleInput, n int, randSource func(int) int) TeamBattleReport {
	if n <= 0 {
		return TeamBattleReport{}
	}

	report := TeamBattleReport{TotalSimulations: n}
	totalTurns := 0
	totalT1Rem := 0
	totalT2Rem := 0

	for i := 0; i < n; i++ {
		result := SimulateTeamBattle(input, randSource)
		totalTurns += result.TotalTurns
		totalT1Rem += result.Team1Remaining
		totalT2Rem += result.Team2Remaining
		switch result.Winner {
		case "team1":
			report.Team1Wins++
		case "team2":
			report.Team2Wins++
		default:
			report.Draws++
		}
	}

	report.Team1WinPct = float64(report.Team1Wins) / float64(n) * 100
	report.Team2WinPct = float64(report.Team2Wins) / float64(n) * 100
	report.DrawPct = float64(report.Draws) / float64(n) * 100
	report.AvgTotalTurns = float64(totalTurns) / float64(n)
	report.AvgTeam1Remaining = float64(totalT1Rem) / float64(n)
	report.AvgTeam2Remaining = float64(totalT2Rem) / float64(n)

	return report
}

// SimulateMultipleBattles runs N full battle simulations and returns aggregated statistics.
// randSource is injected for testability; use rand.Intn in production.
func SimulateMultipleBattles(input FullBattleInput, n int, randSource func(int) int) BattleReport {
	if n <= 0 {
		return BattleReport{}
	}

	summaries := make([]BattleSummary, n)
	for i := 0; i < n; i++ {
		result := SimulateFullBattle(input, randSource)
		summaries[i] = BattleSummary{
			Winner:     result.Winner,
			Turns:      result.TurnCount,
			AttackerHP: result.AttackerHP,
			DefenderHP: result.DefenderHP,
		}
	}

	report := BattleReport{TotalSimulations: n}
	totalTurns := 0
	totalWinnerHP := 0
	turns := make([]int, n)

	for i, s := range summaries {
		turns[i] = s.Turns
		totalTurns += s.Turns
		switch s.Winner {
		case "attacker":
			report.AttackerWins++
			totalWinnerHP += s.AttackerHP
		case "defender":
			report.DefenderWins++
			totalWinnerHP += s.DefenderHP
		default:
			report.Draws++
		}
	}

	report.AttackerWinPct = float64(report.AttackerWins) / float64(n) * 100
	report.DefenderWinPct = float64(report.DefenderWins) / float64(n) * 100
	report.DrawPct = float64(report.Draws) / float64(n) * 100
	report.AvgTurns = float64(totalTurns) / float64(n)

	sort.Ints(turns)
	report.MinTurns = turns[0]
	report.MaxTurns = turns[n-1]
	if n%2 == 0 {
		report.MedianTurns = (turns[n/2-1] + turns[n/2]) / 2
	} else {
		report.MedianTurns = turns[n/2]
	}

	winners := report.AttackerWins + report.DefenderWins
	if winners > 0 {
		report.AvgWinnerHP = float64(totalWinnerHP) / float64(winners)
	}

	return report
}
