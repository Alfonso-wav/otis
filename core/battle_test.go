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
	result := ExecuteTurn(input)

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
	result := ExecuteTurn(input)

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
	result := ExecuteTurn(input)

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

func TestExecuteTurn_NoopWhenOver(t *testing.T) {
	state := InitBattle(200, 0)
	state.IsOver = true
	state.Winner = "attacker"
	input := TurnInput{
		State: state,
		Move:  Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical"},
	}
	result := ExecuteTurn(input)

	if result.NewState.TurnCount != 0 {
		t.Error("Should not increment turn count when battle is already over")
	}
	if len(result.NewState.Log) != 0 {
		t.Error("Should not add log entries when battle is already over")
	}
}
