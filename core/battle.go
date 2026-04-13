package core

import (
	"fmt"
	"sort"
	"strings"
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
	// Ability slugs (kebab-case). Empty string = no ability / unknown.
	AttackerAbility string `json:"attackerAbility"`
	DefenderAbility string `json:"defenderAbility"`
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
	// Record ability slugs on state and fire OnSwitchIn hooks for both sides.
	state.AttackerAbility = normalizeAbilityName(input.AttackerAbility)
	state.DefenderAbility = normalizeAbilityName(input.DefenderAbility)
	if ab, ok := GetAbility(state.AttackerAbility); ok && ab.OnSwitchIn != nil {
		state = ab.OnSwitchIn(state, SideAttacker)
	}
	if ab, ok := GetAbility(state.DefenderAbility); ok && ab.OnSwitchIn != nil {
		state = ab.OnSwitchIn(state, SideDefender)
	}
	const maxTurns = 200

	for !state.IsOver && state.TurnCount < maxTurns {
		// Use overridden moves/stats if Transform has been used.
		atkMoves := input.AttackerMoves
		if len(state.AttackerMovesOverride) > 0 {
			atkMoves = state.AttackerMovesOverride
		}
		defMoves := input.DefenderMoves
		if len(state.DefenderMovesOverride) > 0 {
			defMoves = state.DefenderMovesOverride
		}

		atkMove := atkMoves[randSource(len(atkMoves))]
		defMove := defMoves[randSource(len(defMoves))]

		atkSpeed := input.AttackerStats.Speed
		if state.AttackerStatsOverride != nil {
			atkSpeed = state.AttackerStatsOverride.Speed
		}
		defSpeed := input.DefenderStats.Speed
		if state.DefenderStatsOverride != nil {
			defSpeed = state.DefenderStatsOverride.Speed
		}
		// Apply Spe stage multipliers so Agility / String Shot affect turn order.
		atkSpeed = int(float64(atkSpeed) * StageMultiplier(state.AttackerStages.Spe))
		defSpeed = int(float64(defSpeed) * StageMultiplier(state.DefenderStages.Spe))
		// Ability speed modifiers (Chlorophyll / Swift Swim / etc.).
		if ab, ok := GetAbility(state.AttackerAbility); ok && ab.SpeedMultiplier != nil {
			atkSpeed = int(float64(atkSpeed) * ab.SpeedMultiplier(state, SideAttacker))
		}
		if ab, ok := GetAbility(state.DefenderAbility); ok && ab.SpeedMultiplier != nil {
			defSpeed = int(float64(defSpeed) * ab.SpeedMultiplier(state, SideDefender))
		}

		attackerFirst := resolveOrder(
			atkSpeed, defSpeed,
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

		// End-of-round weather tick: residual damage + duration decrement.
		if !state.IsOver {
			atkTypes := input.AttackerTypes
			if len(state.AttackerTypesOverride) > 0 {
				atkTypes = state.AttackerTypesOverride
			}
			defTypes := input.DefenderTypes
			if len(state.DefenderTypesOverride) > 0 {
				defTypes = state.DefenderTypesOverride
			}
			state = tickWeather(state, atkTypes, defTypes, input.AttackerName, input.DefenderName)
		}

		// End-of-turn ability hooks (Solar Power self-damage, Speed Boost, Rain Dish…).
		if !state.IsOver {
			if ab, ok := GetAbility(state.AttackerAbility); ok && ab.EndOfTurn != nil {
				state = ab.EndOfTurn(state, SideAttacker)
			}
			if state.AttackerHP <= 0 {
				state.IsOver = true
				state.Winner = "defender"
			}
		}
		if !state.IsOver {
			if ab, ok := GetAbility(state.DefenderAbility); ok && ab.EndOfTurn != nil {
				state = ab.EndOfTurn(state, SideDefender)
			}
			if state.DefenderHP <= 0 {
				state.IsOver = true
				state.Winner = "attacker"
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
	// Resolve effective stats/types using Transform overrides.
	atkStats := input.AttackerStats
	if state.AttackerStatsOverride != nil {
		atkStats = *state.AttackerStatsOverride
	}
	atkTypes := input.AttackerTypes
	if len(state.AttackerTypesOverride) > 0 {
		atkTypes = state.AttackerTypesOverride
	}
	defStats := input.DefenderStats
	if state.DefenderStatsOverride != nil {
		defStats = *state.DefenderStatsOverride
	}
	defTypes := input.DefenderTypes
	if len(state.DefenderTypesOverride) > 0 {
		defTypes = state.DefenderTypesOverride
	}
	defMoves := input.DefenderMoves
	if len(state.DefenderMovesOverride) > 0 {
		defMoves = state.DefenderMovesOverride
	}

	result := ExecuteTurn(TurnInput{
		State:         state,
		AttackerStats: atkStats,
		DefenderStats: defStats,
		AttackerTypes: atkTypes,
		DefenderTypes: defTypes,
		AttackerLevel: input.AttackerLevel,
		DefenderLevel: input.DefenderLevel,
		Move:          move,
		AttackerName:  input.AttackerName,
		DefenderName:  input.DefenderName,
		DefenderMoves: defMoves,
	}, randSource)
	return result.NewState
}

func executeDefenderTurn(state BattleState, input FullBattleInput, move Move, randSource func(n int) int) BattleState {
	// Swap overrides along with HP: defender's overrides become attacker's in the swapped view.
	swapped := BattleState{
		AttackerHP:            state.DefenderHP,
		DefenderHP:            state.AttackerHP,
		AttackerMaxHP:         state.DefenderMaxHP,
		DefenderMaxHP:         state.AttackerMaxHP,
		TurnCount:             state.TurnCount,
		Log:                   state.Log,
		IsOver:                state.IsOver,
		Winner:                state.Winner,
		AttackerStatsOverride: state.DefenderStatsOverride,
		AttackerMovesOverride: state.DefenderMovesOverride,
		AttackerTypesOverride: state.DefenderTypesOverride,
		DefenderStatsOverride: state.AttackerStatsOverride,
		DefenderMovesOverride: state.AttackerMovesOverride,
		DefenderTypesOverride: state.AttackerTypesOverride,
		Weather:               state.Weather,
		WeatherTurnsLeft:      state.WeatherTurnsLeft,
		AttackerStages:        state.DefenderStages,
		DefenderStages:        state.AttackerStages,
		AttackerAbility:         state.DefenderAbility,
		DefenderAbility:         state.AttackerAbility,
		AttackerFlashFireActive: state.DefenderFlashFireActive,
		DefenderFlashFireActive: state.AttackerFlashFireActive,
		AttackerIgnoresWeather:  state.DefenderIgnoresWeather,
		DefenderIgnoresWeather:  state.AttackerIgnoresWeather,
	}

	// Resolve effective stats/types using Transform overrides (from swapped perspective).
	atkStats := input.DefenderStats
	if swapped.AttackerStatsOverride != nil {
		atkStats = *swapped.AttackerStatsOverride
	}
	atkTypes := input.DefenderTypes
	if len(swapped.AttackerTypesOverride) > 0 {
		atkTypes = swapped.AttackerTypesOverride
	}
	defStats := input.AttackerStats
	if swapped.DefenderStatsOverride != nil {
		defStats = *swapped.DefenderStatsOverride
	}
	defTypes := input.AttackerTypes
	if len(swapped.DefenderTypesOverride) > 0 {
		defTypes = swapped.DefenderTypesOverride
	}
	defMoves := input.AttackerMoves
	if len(swapped.DefenderMovesOverride) > 0 {
		defMoves = swapped.DefenderMovesOverride
	}

	result := ExecuteTurn(TurnInput{
		State:         swapped,
		AttackerStats: atkStats,
		DefenderStats: defStats,
		AttackerTypes: atkTypes,
		DefenderTypes: defTypes,
		AttackerLevel: input.DefenderLevel,
		DefenderLevel: input.AttackerLevel,
		Move:          move,
		AttackerName:  input.DefenderName,
		DefenderName:  input.AttackerName,
		DefenderMoves: defMoves,
	}, randSource)
	ns := result.NewState
	// Swap back: un-rotate HP and overrides to the original perspective.
	state = BattleState{
		AttackerHP:            ns.DefenderHP,
		DefenderHP:            ns.AttackerHP,
		AttackerMaxHP:         ns.DefenderMaxHP,
		DefenderMaxHP:         ns.AttackerMaxHP,
		TurnCount:             ns.TurnCount,
		Log:                   ns.Log,
		IsOver:                ns.IsOver,
		Winner:                ns.Winner,
		AttackerStatsOverride: ns.DefenderStatsOverride,
		AttackerMovesOverride: ns.DefenderMovesOverride,
		AttackerTypesOverride: ns.DefenderTypesOverride,
		DefenderStatsOverride: ns.AttackerStatsOverride,
		DefenderMovesOverride: ns.AttackerMovesOverride,
		DefenderTypesOverride: ns.AttackerTypesOverride,
		Weather:               ns.Weather,
		WeatherTurnsLeft:      ns.WeatherTurnsLeft,
		AttackerStages:        ns.DefenderStages,
		DefenderStages:        ns.AttackerStages,
		AttackerAbility:         ns.DefenderAbility,
		DefenderAbility:         ns.AttackerAbility,
		AttackerFlashFireActive: ns.DefenderFlashFireActive,
		DefenderFlashFireActive: ns.AttackerFlashFireActive,
		AttackerIgnoresWeather:  ns.DefenderIgnoresWeather,
		DefenderIgnoresWeather:  ns.AttackerIgnoresWeather,
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

	// Transform overrides: when a Pokemon uses Transform, its stats/moves/types
	// are replaced by the opponent's for subsequent turns. HP is NOT overridden.
	// Nil means no override (use original values).
	AttackerStatsOverride *Stats        `json:"attackerStatsOverride,omitempty"`
	AttackerMovesOverride []Move        `json:"attackerMovesOverride,omitempty"`
	AttackerTypesOverride []PokemonType `json:"attackerTypesOverride,omitempty"`
	DefenderStatsOverride *Stats        `json:"defenderStatsOverride,omitempty"`
	DefenderMovesOverride []Move        `json:"defenderMovesOverride,omitempty"`
	DefenderTypesOverride []PokemonType `json:"defenderTypesOverride,omitempty"`

	// Weather field: active weather and turns remaining. WeatherNone + 0 when clear.
	Weather          Weather `json:"weather,omitempty"`
	WeatherTurnsLeft int     `json:"weatherTurnsLeft,omitempty"`

	// Stat stages per side. Zero value = no boosts / no debuffs.
	// Persist across turns in the same battle, reset by InitBattle.
	AttackerStages StatStages `json:"attackerStages,omitempty"`
	DefenderStages StatStages `json:"defenderStages,omitempty"`

	// Ability-related persistent flags (set by OnSwitchIn / OnTakeHit hooks).
	AttackerAbility         string `json:"attackerAbility,omitempty"`
	DefenderAbility         string `json:"defenderAbility,omitempty"`
	AttackerFlashFireActive bool   `json:"attackerFlashFireActive,omitempty"`
	DefenderFlashFireActive bool   `json:"defenderFlashFireActive,omitempty"`
	AttackerIgnoresWeather  bool   `json:"attackerIgnoresWeather,omitempty"`
	DefenderIgnoresWeather  bool   `json:"defenderIgnoresWeather,omitempty"`
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
	DefenderName  string        `json:"defenderName"`
	DefenderMoves []Move        `json:"defenderMoves"`
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

	// Transform: copy defender's stats (except HP), moves, and types.
	// Always hits (accuracy irrelevant). Returns immediately with overrides set.
	if strings.EqualFold(input.Move.Name, "transform") {
		// Copy defender stats but preserve attacker's current HP.
		copiedStats := input.DefenderStats
		copiedStats.HP = input.AttackerStats.HP

		// Copy defender moves into a new slice.
		copiedMoves := make([]Move, len(input.DefenderMoves))
		copy(copiedMoves, input.DefenderMoves)

		// Copy defender types into a new slice.
		copiedTypes := make([]PokemonType, len(input.DefenderTypes))
		copy(copiedTypes, input.DefenderTypes)

		state.AttackerStatsOverride = &copiedStats
		state.AttackerMovesOverride = copiedMoves
		state.AttackerTypesOverride = copiedTypes

		logEntry := fmt.Sprintf("[T%d] %s usó Transform → se transformó en %s",
			state.TurnCount, input.AttackerName, input.DefenderName)
		newLog := make([]string, len(state.Log)+1)
		copy(newLog, state.Log)
		newLog[len(state.Log)] = logEntry
		state.Log = newLog

		return TurnResult{NewState: state, LogEntry: logEntry}
	}

	// Weather-setting moves: Rain Dance / Sunny Day / Sandstorm / Hail.
	// Always succeed, set weather for WeatherDefaultTurns turns. No damage.
	if w, ok := weatherFromMoveName(input.Move.Name); ok {
		state.Weather = w
		state.WeatherTurnsLeft = WeatherDefaultTurns
		logEntry := fmt.Sprintf("[T%d] %s usó %s → %s",
			state.TurnCount, input.AttackerName, input.Move.Name, weatherStartMessage(w))
		newLog := make([]string, len(state.Log)+1)
		copy(newLog, state.Log)
		newLog[len(state.Log)] = logEntry
		state.Log = newLog
		return TurnResult{NewState: state, LogEntry: logEntry}
	}

	// Stat-stage moves: Swords Dance / Growl / etc. No damage, just stage changes.
	if effects, ok := StatStageEffects(input.Move.Name); ok {
		logs := make([]string, len(state.Log), len(state.Log)+1+len(effects))
		copy(logs, state.Log)
		logs = append(logs, fmt.Sprintf("[T%d] %s usó %s",
			state.TurnCount, input.AttackerName, input.Move.Name))
		for _, e := range effects {
			var target string
			var prev int
			if e.Target == "self" {
				prev = GetStage(state.AttackerStages, e.Stat)
				state.AttackerStages = ApplyStage(state.AttackerStages, e.Stat, e.Delta)
				target = input.AttackerName
			} else {
				prev = GetStage(state.DefenderStages, e.Stat)
				state.DefenderStages = ApplyStage(state.DefenderStages, e.Stat, e.Delta)
				target = input.DefenderName
			}
			fragment := StageChangeLogFragment(prev, e.Delta)
			logs = append(logs, fmt.Sprintf("%s de %s %s",
				StatLabelEs(e.Stat), target, fragment))
		}
		state.Log = logs
		return TurnResult{NewState: state, LogEntry: logs[len(logs)-1]}
	}

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

	dmgInput := DamageInput{
		AttackerStats:   input.AttackerStats,
		DefenderStats:   input.DefenderStats,
		Move:            input.Move,
		AttackerTypes:   input.AttackerTypes,
		DefenderTypes:   input.DefenderTypes,
		Level:           input.AttackerLevel,
		IsCritical:      false,
		WeatherBonus:    1.0,
		Weather:         state.Weather,
		AttackerStages:  state.AttackerStages,
		DefenderStages:  state.DefenderStages,
		AttackerAbility: state.AttackerAbility,
		DefenderAbility: state.DefenderAbility,
		State:           state,
	}
	var dmg DamageResult
	if randSource != nil {
		dmg = CalculateBattleDamage(dmgInput, randSource)
	} else {
		dmg = CalculateDamage(dmgInput)
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

		// Defender's OnTakeHit can heal (Water Absorb), activate flags
		// (Flash Fire) or zero out incoming damage. Attacker's OnTakeHit
		// also runs for recoil-style abilities on the attacker side (Rough
		// Skin style), but here we only invoke the defender's hook since
		// attacker-side reactions belong to the *defender's* ability.
		if ab, ok := GetAbility(state.DefenderAbility); ok && ab.OnTakeHit != nil {
			state, applied = ab.OnTakeHit(state, SideDefender, input.Move, applied)
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

// ChooseBestMember selects the member from available with the best type advantage
// against opponentTypes. Returns the index into the available slice.
// Advantage is calculated by finding the max effectiveness of any move the member
// has against the opponent's types. If no clear advantage, returns 0.
func ChooseBestMember(available []TeamBattleMember, opponentTypes []PokemonType) int {
	if len(available) <= 1 {
		return 0
	}
	bestIdx := 0
	bestScore := -1.0
	for i, m := range available {
		score := 0.0
		for _, mv := range m.Moves {
			if mv.Power <= 0 {
				continue
			}
			eff := TypeEffectiveness(mv.Type, opponentTypes)
			if eff > score {
				score = eff
			}
		}
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}
	return bestIdx
}

// shuffleMembers returns a new shuffled copy of the members slice using randSource.
func shuffleMembers(members []TeamBattleMember, randSource func(int) int) []TeamBattleMember {
	shuffled := make([]TeamBattleMember, len(members))
	copy(shuffled, members)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := randSource(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

// SimulateTeamBattle simulates a full team-vs-team battle.
// Members are shuffled at the start. When one Pokemon faints, the replacement
// is chosen by type advantage (smart switching). The winner's HP carries over.
func SimulateTeamBattle(input TeamBattleInput, randSource func(int) int) TeamBattleState {
	if len(input.Team1Members) == 0 || len(input.Team2Members) == 0 {
		return TeamBattleState{IsOver: true, Winner: "draw"}
	}

	// Shuffle initial order
	avail1 := shuffleMembers(input.Team1Members, randSource)
	avail2 := shuffleMembers(input.Team2Members, randSource)

	carryHP1, carryHP2 := 0, 0 // 0 means use full HP
	totalTurns := 0
	var rounds []BattleState
	var log []string

	for len(avail1) > 0 && len(avail2) > 0 {
		m1 := avail1[0]
		m2 := avail2[0]

		bi := FullBattleInput{
			AttackerStats:   m1.Stats,
			DefenderStats:   m2.Stats,
			AttackerTypes:   m1.Types,
			DefenderTypes:   m2.Types,
			AttackerLevel:   m1.Level,
			DefenderLevel:   m2.Level,
			AttackerMoves:   m1.Moves,
			DefenderMoves:   m2.Moves,
			AttackerName:    m1.PokemonName,
			DefenderName:    m2.PokemonName,
			AttackerAbility: m1.Ability,
			DefenderAbility: m2.Ability,
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

		// Add round header and detailed turn-by-turn log
		log = append(log, fmt.Sprintf("--- Ronda %d: %s vs %s ---",
			roundNum, m1.PokemonName, m2.PokemonName))
		log = append(log, result.Log...)

		if result.Winner == "attacker" {
			log = append(log, fmt.Sprintf("[Ronda %d] %s venció a %s (HP restante: %d)",
				roundNum, m1.PokemonName, m2.PokemonName, result.AttackerHP))
			carryHP1 = result.AttackerHP
			carryHP2 = 0
			// Remove defeated member from team2
			avail2 = avail2[1:]
			// Smart switch: choose best counter for surviving opponent
			if len(avail2) > 0 {
				bestIdx := ChooseBestMember(avail2, m1.Types)
				if bestIdx != 0 {
					avail2[0], avail2[bestIdx] = avail2[bestIdx], avail2[0]
				}
			}
		} else {
			log = append(log, fmt.Sprintf("[Ronda %d] %s venció a %s (HP restante: %d)",
				roundNum, m2.PokemonName, m1.PokemonName, result.DefenderHP))
			carryHP2 = result.DefenderHP
			carryHP1 = 0
			// Remove defeated member from team1
			avail1 = avail1[1:]
			// Smart switch: choose best counter for surviving opponent
			if len(avail1) > 0 {
				bestIdx := ChooseBestMember(avail1, m2.Types)
				if bestIdx != 0 {
					avail1[0], avail1[bestIdx] = avail1[bestIdx], avail1[0]
				}
			}
		}
	}

	winner := "team1"
	if len(avail1) == 0 {
		winner = "team2"
	}

	return TeamBattleState{
		Team1Remaining: len(avail1),
		Team2Remaining: len(avail2),
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

// weatherFromMoveName maps a weather-setting move to the Weather it applies.
// Accepts PokéAPI slugs (e.g. "rain-dance") case-insensitively.
func weatherFromMoveName(name string) (Weather, bool) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "rain-dance", "rain dance":
		return WeatherRain, true
	case "sunny-day", "sunny day":
		return WeatherSun, true
	case "sandstorm":
		return WeatherSandstorm, true
	case "hail":
		return WeatherHail, true
	}
	return WeatherNone, false
}

// weatherStartMessage returns the Spanish log line when a weather becomes active.
func weatherStartMessage(w Weather) string {
	switch w {
	case WeatherRain:
		return "¡Empezó a llover!"
	case WeatherSun:
		return "¡La luz solar se hizo intensa!"
	case WeatherSandstorm:
		return "¡Se desató una tormenta de arena!"
	case WeatherHail:
		return "¡Empezó a granizar!"
	}
	return ""
}

// weatherEndMessage returns the Spanish log line when a weather expires.
func weatherEndMessage(w Weather) string {
	switch w {
	case WeatherRain:
		return "La lluvia cesó."
	case WeatherSun:
		return "La luz solar volvió a la normalidad."
	case WeatherSandstorm:
		return "La tormenta de arena amainó."
	case WeatherHail:
		return "El granizo cesó."
	}
	return "El clima ha vuelto a la normalidad."
}

// sandstormImmune reports whether a Pokémon is immune to Sandstorm residual damage.
// Rock, Ground and Steel types take no Sandstorm chip damage.
func sandstormImmune(types []PokemonType) bool {
	return hasType(types, "rock") || hasType(types, "ground") || hasType(types, "steel")
}

// hailImmune reports whether a Pokémon is immune to Hail residual damage. Only Ice.
func hailImmune(types []PokemonType) bool {
	return hasType(types, "ice")
}

// residualWeatherDamage returns the residual HP damage a Pokémon suffers under
// the given weather. Sandstorm/Hail: 1/16 of max HP (min 1) unless immune.
func residualWeatherDamage(w Weather, maxHP int, types []PokemonType) int {
	switch w {
	case WeatherSandstorm:
		if sandstormImmune(types) {
			return 0
		}
	case WeatherHail:
		if hailImmune(types) {
			return 0
		}
	default:
		return 0
	}
	dmg := maxHP / 16
	if dmg < 1 {
		dmg = 1
	}
	return dmg
}

// tickWeather is a pure end-of-round helper. Given the current battle state and
// both sides' effective types/names, it applies residual weather damage to both
// sides, decrements WeatherTurnsLeft, and clears weather when it expires.
// Returns the new state. Marks the battle over if residual damage faints someone.
func tickWeather(state BattleState, atkTypes, defTypes []PokemonType, atkName, defName string) BattleState {
	if state.IsOver || state.Weather == WeatherNone {
		return state
	}

	logs := state.Log
	// Residual damage to attacker
	if d := residualWeatherDamage(state.Weather, state.AttackerMaxHP, atkTypes); d > 0 {
		newHP := state.AttackerHP - d
		if newHP < 0 {
			newHP = 0
		}
		state.AttackerHP = newHP
		logs = append(logs, fmt.Sprintf("[T%d] %s sufrió %d de daño por el clima | HP: %d/%d",
			state.TurnCount, atkName, d, state.AttackerHP, state.AttackerMaxHP))
		if state.AttackerHP == 0 {
			state.IsOver = true
			state.Winner = "defender"
		}
	}
	// Residual damage to defender
	if !state.IsOver {
		if d := residualWeatherDamage(state.Weather, state.DefenderMaxHP, defTypes); d > 0 {
			newHP := state.DefenderHP - d
			if newHP < 0 {
				newHP = 0
			}
			state.DefenderHP = newHP
			logs = append(logs, fmt.Sprintf("[T%d] %s sufrió %d de daño por el clima | HP: %d/%d",
				state.TurnCount, defName, d, state.DefenderHP, state.DefenderMaxHP))
			if state.DefenderHP == 0 {
				state.IsOver = true
				state.Winner = "attacker"
			}
		}
	}

	// Decrement duration
	if state.WeatherTurnsLeft > 0 {
		state.WeatherTurnsLeft--
	}
	if state.WeatherTurnsLeft <= 0 {
		logs = append(logs, fmt.Sprintf("[T%d] %s", state.TurnCount, weatherEndMessage(state.Weather)))
		state.Weather = WeatherNone
		state.WeatherTurnsLeft = 0
	}

	state.Log = logs
	return state
}
