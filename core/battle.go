package core

import "fmt"

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
