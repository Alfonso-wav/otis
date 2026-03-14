package core

import (
	"strings"
	"testing"
)

func TestInitBattle(t *testing.T) {
	s := InitBattle(200, 150)
	if s.AttackerHP != 200 {
		t.Errorf("AttackerHP: want 200, got %d", s.AttackerHP)
	}
	if s.DefenderHP != 150 {
		t.Errorf("DefenderHP: want 150, got %d", s.DefenderHP)
	}
	if s.AttackerMaxHP != 200 {
		t.Errorf("AttackerMaxHP: want 200, got %d", s.AttackerMaxHP)
	}
	if s.DefenderMaxHP != 150 {
		t.Errorf("DefenderMaxHP: want 150, got %d", s.DefenderMaxHP)
	}
	if s.TurnCount != 0 {
		t.Errorf("TurnCount: want 0, got %d", s.TurnCount)
	}
	if s.IsOver {
		t.Error("IsOver should be false on init")
	}
	if s.Winner != "" {
		t.Errorf("Winner: want '', got %q", s.Winner)
	}
	if len(s.Log) != 0 {
		t.Errorf("Log: want empty, got %d entries", len(s.Log))
	}
}

func TestExecuteTurn_DealsDamage(t *testing.T) {
	state := InitBattle(200, 200)
	input := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 100, Defense: 80, SpAttack: 100, SpDefense: 80, Speed: 80, HP: 200},
		DefenderStats: Stats{Attack: 80, Defense: 80, SpAttack: 80, SpDefense: 80, Speed: 80, HP: 200},
		AttackerTypes: []PokemonType{{Name: "fire"}},
		DefenderTypes: []PokemonType{{Name: "grass"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		Move:          Move{Name: "Flamethrower", Type: "fire", Power: 90, Category: "special", Accuracy: 100},
	}
	result := ExecuteTurn(input, nil)

	if result.NewState.TurnCount != 1 {
		t.Errorf("TurnCount: want 1, got %d", result.NewState.TurnCount)
	}
	if result.NewState.DefenderHP >= 200 {
		t.Error("DefenderHP should decrease after a damaging move")
	}
	if len(result.NewState.Log) != 1 {
		t.Errorf("Log: want 1 entry, got %d", len(result.NewState.Log))
	}
	if result.Damage.IsSuperEffective == false {
		t.Error("Fire vs Grass should be super effective")
	}
}

func TestExecuteTurn_StatusMoveNoDamage(t *testing.T) {
	state := InitBattle(200, 200)
	input := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 100, HP: 200},
		DefenderStats: Stats{Defense: 80, HP: 200},
		AttackerLevel: 50,
		Move:          Move{Name: "Swords Dance", Type: "normal", Power: 0, Category: "status"},
	}
	result := ExecuteTurn(input, nil)

	if result.NewState.DefenderHP != 200 {
		t.Errorf("DefenderHP: want 200 after status move, got %d", result.NewState.DefenderHP)
	}
	if !strings.Contains(result.LogEntry, "sin efecto de daño") {
		t.Errorf("Log should mention no damage, got: %s", result.LogEntry)
	}
	if result.NewState.IsOver {
		t.Error("Battle should not end from a status move")
	}
}

func TestExecuteTurn_WinCondition(t *testing.T) {
	state := InitBattle(200, 1)
	input := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 200, HP: 200},
		DefenderStats: Stats{Defense: 1, HP: 1},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		Move:          Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100},
	}
	result := ExecuteTurn(input, nil)

	if !result.NewState.IsOver {
		t.Error("Battle should be over when defender HP reaches 0")
	}
	if result.NewState.Winner != "attacker" {
		t.Errorf("Winner: want 'attacker', got %q", result.NewState.Winner)
	}
	if result.NewState.DefenderHP != 0 {
		t.Errorf("DefenderHP: want 0, got %d", result.NewState.DefenderHP)
	}
}

func TestSimulateFullBattle_AttackerWins(t *testing.T) {
	strongMove := Move{Name: "Hyper Beam", Type: "normal", Power: 150, Category: "special", Accuracy: 100}
	weakMove := Move{Name: "Splash", Type: "water", Power: 0, Category: "status"}

	input := FullBattleInput{
		AttackerStats: Stats{HP: 200, Attack: 200, SpAttack: 200, Defense: 80, SpDefense: 80, Speed: 80},
		DefenderStats: Stats{HP: 10, Attack: 10, SpAttack: 10, Defense: 10, SpDefense: 10, Speed: 10},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 100,
		DefenderLevel: 50,
		AttackerMoves: []Move{strongMove},
		DefenderMoves: []Move{weakMove},
	}

	alwaysZero := func(n int) int { return 0 }
	result := SimulateFullBattle(input, alwaysZero)

	if !result.IsOver {
		t.Error("Battle should be over")
	}
	if result.Winner != "attacker" {
		t.Errorf("Winner: want 'attacker', got %q", result.Winner)
	}
	if result.TurnCount == 0 {
		t.Error("TurnCount should be > 0")
	}
	if len(result.Log) == 0 {
		t.Error("Log should not be empty")
	}
}

func TestSimulateFullBattle_DefenderWins(t *testing.T) {
	strongMove := Move{Name: "Hyper Beam", Type: "normal", Power: 150, Category: "special", Accuracy: 100}
	weakMove := Move{Name: "Splash", Type: "water", Power: 0, Category: "status"}

	input := FullBattleInput{
		AttackerStats: Stats{HP: 10, Attack: 10, SpAttack: 10, Defense: 10, SpDefense: 10, Speed: 10},
		DefenderStats: Stats{HP: 200, Attack: 200, SpAttack: 200, Defense: 80, SpDefense: 80, Speed: 80},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		DefenderLevel: 100,
		AttackerMoves: []Move{weakMove},
		DefenderMoves: []Move{strongMove},
	}

	alwaysZero := func(n int) int { return 0 }
	result := SimulateFullBattle(input, alwaysZero)

	if !result.IsOver {
		t.Error("Battle should be over")
	}
	if result.Winner != "defender" {
		t.Errorf("Winner: want 'defender', got %q", result.Winner)
	}
}

func TestSimulateFullBattle_EmptyMoves(t *testing.T) {
	input := FullBattleInput{
		AttackerStats: Stats{HP: 100},
		DefenderStats: Stats{HP: 100},
		AttackerMoves: []Move{},
		DefenderMoves: []Move{{Name: "Tackle", Type: "normal", Power: 40, Category: "physical"}},
	}
	result := SimulateFullBattle(input, func(n int) int { return 0 })
	if result.IsOver {
		t.Error("Should return empty state when attacker has no moves")
	}
}

func TestSimulateFullBattle_MaxTurns(t *testing.T) {
	// Both use status moves → no damage → resolves by HP (equal → draw)
	statusMove := Move{Name: "Splash", Type: "water", Power: 0, Category: "status"}
	input := FullBattleInput{
		AttackerStats: Stats{HP: 100, Attack: 80, SpAttack: 80, Defense: 80, SpDefense: 80, Speed: 80},
		DefenderStats: Stats{HP: 100, Attack: 80, SpAttack: 80, Defense: 80, SpDefense: 80, Speed: 80},
		AttackerMoves: []Move{statusMove},
		DefenderMoves: []Move{statusMove},
	}
	result := SimulateFullBattle(input, func(n int) int { return 0 })
	if !result.IsOver {
		t.Error("Battle should be over after max turns")
	}
	if result.Winner != "draw" {
		t.Errorf("Winner: want 'draw', got %q", result.Winner)
	}
}

func TestExecuteTurn_AccuracyMiss(t *testing.T) {
	state := InitBattle(200, 200)
	input := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 100, HP: 200},
		DefenderStats: Stats{Defense: 80, HP: 200},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		Move:          Move{Name: "Thunder", Type: "electric", Power: 110, Category: "special", Accuracy: 70},
	}
	// randSource returns 99 for accuracy check → 99 >= 70 → miss
	callCount := 0
	alwaysMiss := func(n int) int {
		callCount++
		if callCount == 1 {
			return 99 // accuracy check: 99 >= 70 → miss
		}
		return 0
	}
	result := ExecuteTurn(input, alwaysMiss)

	if !result.Missed {
		t.Error("Expected move to miss with accuracy 70 and roll 99")
	}
	if result.NewState.DefenderHP != 200 {
		t.Errorf("DefenderHP should be unchanged after miss, got %d", result.NewState.DefenderHP)
	}
	if !strings.Contains(result.LogEntry, "¡Falló!") {
		t.Errorf("Log should contain miss message, got: %s", result.LogEntry)
	}
}

func TestExecuteTurn_AccuracyHit(t *testing.T) {
	state := InitBattle(200, 200)
	input := TurnInput{
		State:         state,
		AttackerStats: Stats{SpAttack: 100, HP: 200},
		DefenderStats: Stats{SpDefense: 80, HP: 200},
		AttackerTypes: []PokemonType{{Name: "electric"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		Move:          Move{Name: "Thunder", Type: "electric", Power: 110, Category: "special", Accuracy: 70},
	}
	// randSource returns 0 for all checks → accuracy hit, no crit, lowest roll
	alwaysZero := func(n int) int { return 0 }
	result := ExecuteTurn(input, alwaysZero)

	if result.Missed {
		t.Error("Expected move to hit with accuracy 70 and roll 0")
	}
	if result.NewState.DefenderHP >= 200 {
		t.Error("DefenderHP should decrease after hit")
	}
}

func TestExecuteTurn_CriticalHit(t *testing.T) {
	state := InitBattle(200, 200)
	input := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 100, HP: 200},
		DefenderStats: Stats{Defense: 100, HP: 200},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		Move:          Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100},
	}
	// randSource: accuracy=0(hit), crit=0(crit at stage 0 means 1/24, roll 0 → crit), roll=0(min)
	alwaysZero := func(n int) int { return 0 }
	result := ExecuteTurn(input, alwaysZero)

	if !result.Damage.WasCritical {
		t.Error("Expected critical hit when randSource returns 0 (1/24 chance)")
	}
	if !strings.Contains(result.LogEntry, "¡Golpe crítico!") {
		t.Errorf("Log should contain critical hit message, got: %s", result.LogEntry)
	}
}

func TestExecuteTurn_NoCritWhenRollHigh(t *testing.T) {
	state := InitBattle(200, 200)
	input := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 100, HP: 200},
		DefenderStats: Stats{Defense: 100, HP: 200},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		Move:          Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100},
	}
	// Accuracy 100 bypasses check. CalculateBattleDamage calls:
	// 1st: randSource(24) for crit → 5 != 0 → no crit
	// 2nd: randSource(16) for roll → 8
	callCount := 0
	mockRand := func(n int) int {
		callCount++
		switch callCount {
		case 1:
			return 5 // crit: 5 != 0 → no crit (stage 0, threshold 24)
		case 2:
			return 8 // roll: 85+8=93
		}
		return 0
	}
	result := ExecuteTurn(input, mockRand)

	if result.Damage.WasCritical {
		t.Error("Expected no critical hit when crit roll is 5 (threshold 24)")
	}
}

func TestSimulateFullBattle_SpeedOrdering(t *testing.T) {
	// Attacker is slower but has a priority move → should still go first
	priorityMove := Move{Name: "Quick Attack", Type: "normal", Power: 40, Category: "physical", Accuracy: 100, Priority: 1}
	normalMove := Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100, Priority: 0}

	input := FullBattleInput{
		AttackerStats: Stats{HP: 200, Attack: 200, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 10},
		DefenderStats: Stats{HP: 1, Attack: 100, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 100},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{priorityMove},
		DefenderMoves: []Move{normalMove},
	}

	alwaysZero := func(n int) int { return 0 }
	result := SimulateFullBattle(input, alwaysZero)

	if !result.IsOver {
		t.Error("Battle should be over")
	}
	// Attacker has priority move so should KO the 1 HP defender before defender can attack
	if result.Winner != "attacker" {
		t.Errorf("Winner: want 'attacker' (priority move), got %q", result.Winner)
	}
}

func TestSimulateFullBattle_FasterGoesFirst(t *testing.T) {
	// Defender is faster and has enough power to KO attacker (1 HP)
	move := Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100}

	input := FullBattleInput{
		AttackerStats: Stats{HP: 1, Attack: 200, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 10},
		DefenderStats: Stats{HP: 200, Attack: 200, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 100},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{move},
		DefenderMoves: []Move{move},
	}

	alwaysZero := func(n int) int { return 0 }
	result := SimulateFullBattle(input, alwaysZero)

	if !result.IsOver {
		t.Error("Battle should be over")
	}
	// Defender is faster → attacks first → KOs the 1 HP attacker
	if result.Winner != "defender" {
		t.Errorf("Winner: want 'defender' (faster), got %q", result.Winner)
	}
}

func TestExecuteTurn_NoopWhenOver(t *testing.T) {
	state := InitBattle(200, 0)
	state.IsOver = true
	state.Winner = "attacker"
	input := TurnInput{
		State: state,
		Move:  Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical"},
	}
	result := ExecuteTurn(input, nil)

	if result.NewState.TurnCount != 0 {
		t.Error("Should not increment turn count when battle is already over")
	}
	if len(result.NewState.Log) != 0 {
		t.Error("Should not add log entries when battle is already over")
	}
}
