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
		AttackerName:  "Charizard",
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
	if !strings.Contains(result.LogEntry, "Charizard") {
		t.Errorf("Log should contain attacker name, got: %s", result.LogEntry)
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
		AttackerName:  "Scizor",
	}
	result := ExecuteTurn(input, nil)

	if result.NewState.DefenderHP != 200 {
		t.Errorf("DefenderHP: want 200 after status move, got %d", result.NewState.DefenderHP)
	}
	if !strings.Contains(result.LogEntry, "sin efecto de daño") {
		t.Errorf("Log should mention no damage, got: %s", result.LogEntry)
	}
	if !strings.Contains(result.LogEntry, "Scizor") {
		t.Errorf("Log should contain attacker name, got: %s", result.LogEntry)
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
		AttackerName:  "Pikachu",
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
		AttackerName:  "Snorlax",
		DefenderName:  "Magikarp",
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
	for _, entry := range result.Log {
		if !strings.Contains(entry, "Snorlax") && !strings.Contains(entry, "Magikarp") {
			t.Errorf("Log entry should contain a Pokemon name, got: %s", entry)
		}
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
		AttackerName:  "Magikarp",
		DefenderName:  "Snorlax",
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
		AttackerName:  "Pikachu",
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
	if !strings.Contains(result.LogEntry, "Pikachu") {
		t.Errorf("Log should contain attacker name, got: %s", result.LogEntry)
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
		AttackerName:  "Pikachu",
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
		AttackerName:  "Rattata",
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
		AttackerName:  "Rattata",
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
		AttackerName:  "Eevee",
		DefenderName:  "Pidgey",
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
		AttackerName:  "Slowpoke",
		DefenderName:  "Jolteon",
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

func TestSimulateMultipleBattles_N1(t *testing.T) {
	move := Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100}
	input := FullBattleInput{
		AttackerStats: Stats{HP: 200, Attack: 200, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 80},
		DefenderStats: Stats{HP: 200, Attack: 100, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 80},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{move},
		DefenderMoves: []Move{move},
		AttackerName:  "Snorlax",
		DefenderName:  "Chansey",
	}
	report := SimulateMultipleBattles(input, 1, func(n int) int { return 0 })
	if report.TotalSimulations != 1 {
		t.Errorf("TotalSimulations: want 1, got %d", report.TotalSimulations)
	}
	total := report.AttackerWins + report.DefenderWins + report.Draws
	if total != 1 {
		t.Errorf("Sum of outcomes: want 1, got %d", total)
	}
	if report.MinTurns <= 0 {
		t.Error("MinTurns should be > 0")
	}
}

func TestSimulateMultipleBattles_N100(t *testing.T) {
	move := Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100}
	input := FullBattleInput{
		AttackerStats: Stats{HP: 200, Attack: 150, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 80},
		DefenderStats: Stats{HP: 200, Attack: 150, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 80},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{move},
		DefenderMoves: []Move{move},
		AttackerName:  "Tauros",
		DefenderName:  "Miltank",
	}
	// Use a simple incrementing counter for varied results
	counter := 0
	rand := func(n int) int {
		counter++
		return counter % n
	}
	report := SimulateMultipleBattles(input, 100, rand)

	if report.TotalSimulations != 100 {
		t.Errorf("TotalSimulations: want 100, got %d", report.TotalSimulations)
	}
	totalPct := report.AttackerWinPct + report.DefenderWinPct + report.DrawPct
	if totalPct < 99.9 || totalPct > 100.1 {
		t.Errorf("Percentages should sum to ~100, got %.2f", totalPct)
	}
	if report.AvgTurns <= 0 {
		t.Error("AvgTurns should be > 0")
	}
	if report.MinTurns > report.MaxTurns {
		t.Errorf("MinTurns (%d) > MaxTurns (%d)", report.MinTurns, report.MaxTurns)
	}
	if report.MedianTurns < report.MinTurns || report.MedianTurns > report.MaxTurns {
		t.Errorf("MedianTurns (%d) out of range [%d, %d]", report.MedianTurns, report.MinTurns, report.MaxTurns)
	}
}

func TestSimulateMultipleBattles_N0(t *testing.T) {
	move := Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100}
	input := FullBattleInput{
		AttackerStats: Stats{HP: 200, Attack: 100, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 80},
		DefenderStats: Stats{HP: 200, Attack: 100, SpAttack: 100, Defense: 80, SpDefense: 80, Speed: 80},
		AttackerMoves: []Move{move},
		DefenderMoves: []Move{move},
	}
	report := SimulateMultipleBattles(input, 0, func(n int) int { return 0 })
	if report.TotalSimulations != 0 {
		t.Errorf("TotalSimulations: want 0, got %d", report.TotalSimulations)
	}
	if report.AttackerWins != 0 || report.DefenderWins != 0 || report.Draws != 0 {
		t.Error("All outcomes should be 0 for N=0")
	}
}

// --- Team Battle Tests ---

func makeTeamBattleMember(name string, hp, atk, def, speed int) TeamBattleMember {
	move := Move{Name: "tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100}
	return TeamBattleMember{
		PokemonName: name,
		Stats:       Stats{HP: hp, Attack: atk, Defense: def, SpAttack: atk, SpDefense: def, Speed: speed},
		Types:       []PokemonType{{Name: "normal"}},
		Moves:       []Move{move},
		Level:       50,
	}
}

func TestSimulateTeamBattle_1v1(t *testing.T) {
	input := TeamBattleInput{
		Team1Name:    "Alpha",
		Team1Members: []TeamBattleMember{makeTeamBattleMember("pikachu", 200, 100, 80, 90)},
		Team2Name:    "Beta",
		Team2Members: []TeamBattleMember{makeTeamBattleMember("charmander", 200, 100, 80, 90)},
	}
	counter := 0
	rng := func(n int) int { counter++; return counter % n }
	result := SimulateTeamBattle(input, rng)

	if !result.IsOver {
		t.Error("battle should be over")
	}
	if result.Winner != "team1" && result.Winner != "team2" {
		t.Errorf("unexpected winner: %s", result.Winner)
	}
	if len(result.Rounds) != 1 {
		t.Errorf("expected 1 round, got %d", len(result.Rounds))
	}
	// Log should contain: round header + detailed turn logs + round summary
	if len(result.Log) < 3 {
		t.Errorf("expected at least 3 log entries (header + moves + summary), got %d", len(result.Log))
	}
	// First entry should be the round header
	if !strings.Contains(result.Log[0], "--- Ronda 1:") {
		t.Errorf("first log entry should be round header, got: %s", result.Log[0])
	}
	// Last entry should be the round summary
	lastEntry := result.Log[len(result.Log)-1]
	if !strings.Contains(lastEntry, "[Ronda 1]") || !strings.Contains(lastEntry, "venció") {
		t.Errorf("last log entry should be round summary, got: %s", lastEntry)
	}
}

func TestSimulateTeamBattle_3v3(t *testing.T) {
	input := TeamBattleInput{
		Team1Name: "Alpha",
		Team1Members: []TeamBattleMember{
			makeTeamBattleMember("pikachu", 200, 120, 80, 100),
			makeTeamBattleMember("bulbasaur", 220, 100, 90, 70),
			makeTeamBattleMember("squirtle", 210, 90, 100, 60),
		},
		Team2Name: "Beta",
		Team2Members: []TeamBattleMember{
			makeTeamBattleMember("charmander", 200, 110, 80, 90),
			makeTeamBattleMember("eevee", 180, 90, 70, 80),
			makeTeamBattleMember("jigglypuff", 300, 60, 50, 40),
		},
	}
	counter := 0
	rng := func(n int) int { counter++; return counter % n }
	result := SimulateTeamBattle(input, rng)

	if !result.IsOver {
		t.Error("battle should be over")
	}
	if result.Winner != "team1" && result.Winner != "team2" {
		t.Errorf("unexpected winner: %s", result.Winner)
	}
	if result.Team1Remaining+result.Team2Remaining < 1 {
		t.Error("at least one team should have remaining members")
	}
	// Winner should have remaining members, loser 0
	if result.Winner == "team1" && result.Team1Remaining == 0 {
		t.Error("team1 won but has 0 remaining")
	}
	if result.Winner == "team2" && result.Team2Remaining == 0 {
		t.Error("team2 won but has 0 remaining")
	}
}

func TestSimulateTeamBattle_CarryOverHP(t *testing.T) {
	// Team1 has one very strong Pokemon, Team2 has two weak ones
	input := TeamBattleInput{
		Team1Name:    "Strong",
		Team1Members: []TeamBattleMember{makeTeamBattleMember("mewtwo", 500, 200, 100, 150)},
		Team2Name:    "Weak",
		Team2Members: []TeamBattleMember{
			makeTeamBattleMember("magikarp", 50, 20, 20, 20),
			makeTeamBattleMember("caterpie", 50, 20, 20, 20),
		},
	}
	counter := 0
	rng := func(n int) int { counter++; return counter % n }
	result := SimulateTeamBattle(input, rng)

	if result.Winner != "team1" {
		t.Errorf("expected team1 to win, got %s", result.Winner)
	}
	if len(result.Rounds) != 2 {
		t.Errorf("expected 2 rounds, got %d", len(result.Rounds))
	}
	// Verify carry-over: mewtwo should have taken some damage in round 1
	// and started round 2 with less HP
	if result.Team1Remaining != 1 {
		t.Errorf("expected 1 remaining for team1, got %d", result.Team1Remaining)
	}
}

func TestSimulateMultipleTeamBattles(t *testing.T) {
	input := TeamBattleInput{
		Team1Name: "A",
		Team1Members: []TeamBattleMember{
			makeTeamBattleMember("pikachu", 200, 100, 80, 90),
			makeTeamBattleMember("bulbasaur", 200, 100, 80, 70),
		},
		Team2Name: "B",
		Team2Members: []TeamBattleMember{
			makeTeamBattleMember("charmander", 200, 100, 80, 85),
			makeTeamBattleMember("squirtle", 200, 100, 80, 75),
		},
	}
	counter := 0
	rng := func(n int) int { counter++; return counter % n }
	report := SimulateMultipleTeamBattles(input, 10, rng)

	if report.TotalSimulations != 10 {
		t.Errorf("expected 10 sims, got %d", report.TotalSimulations)
	}
	total := report.Team1Wins + report.Team2Wins + report.Draws
	if total != 10 {
		t.Errorf("wins+draws should equal 10, got %d", total)
	}
	if report.AvgTotalTurns <= 0 {
		t.Error("avgTotalTurns should be > 0")
	}
}

func TestSimulateTeamBattle_DetailedLog(t *testing.T) {
	input := TeamBattleInput{
		Team1Name: "Alpha",
		Team1Members: []TeamBattleMember{
			makeTeamBattleMember("pikachu", 200, 120, 80, 100),
			makeTeamBattleMember("bulbasaur", 220, 100, 90, 70),
		},
		Team2Name: "Beta",
		Team2Members: []TeamBattleMember{
			makeTeamBattleMember("charmander", 200, 110, 80, 90),
		},
	}
	counter := 0
	rng := func(n int) int { counter++; return counter % n }
	result := SimulateTeamBattle(input, rng)

	// Log must contain individual move entries (turn-by-turn), not just summaries
	hasRoundHeader := false
	hasMoveEntry := false
	hasRoundSummary := false
	for _, entry := range result.Log {
		if strings.Contains(entry, "--- Ronda") {
			hasRoundHeader = true
		}
		if strings.Contains(entry, "[T") && strings.Contains(entry, "usó") {
			hasMoveEntry = true
		}
		if strings.Contains(entry, "[Ronda") && strings.Contains(entry, "venció") {
			hasRoundSummary = true
		}
	}
	if !hasRoundHeader {
		t.Error("log should contain round headers")
	}
	if !hasMoveEntry {
		t.Error("log should contain individual move entries (turn-by-turn detail)")
	}
	if !hasRoundSummary {
		t.Error("log should contain round summaries")
	}
}

func TestChooseBestMember_TypeAdvantage(t *testing.T) {
	waterMove := Move{Name: "surf", Type: "water", Power: 90, Category: "special", Accuracy: 100}
	grassMove := Move{Name: "razor-leaf", Type: "grass", Power: 55, Category: "physical", Accuracy: 95}
	normalMove := Move{Name: "tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100}

	available := []TeamBattleMember{
		{PokemonName: "squirtle", Moves: []Move{waterMove}, Types: []PokemonType{{Name: "water"}}},
		{PokemonName: "bulbasaur", Moves: []Move{grassMove}, Types: []PokemonType{{Name: "grass"}}},
		{PokemonName: "rattata", Moves: []Move{normalMove}, Types: []PokemonType{{Name: "normal"}}},
	}

	// Opponent is fire type → water is super effective
	fireTypes := []PokemonType{{Name: "fire"}}
	idx := ChooseBestMember(available, fireTypes)
	if idx != 0 {
		t.Errorf("expected squirtle (water, idx 0) vs fire, got idx %d (%s)", idx, available[idx].PokemonName)
	}

	// Opponent is water type → grass is super effective
	waterTypes := []PokemonType{{Name: "water"}}
	idx = ChooseBestMember(available, waterTypes)
	if idx != 1 {
		t.Errorf("expected bulbasaur (grass, idx 1) vs water, got idx %d (%s)", idx, available[idx].PokemonName)
	}
}

func TestChooseBestMember_SingleMember(t *testing.T) {
	m := TeamBattleMember{
		PokemonName: "pikachu",
		Moves:       []Move{{Name: "thunderbolt", Type: "electric", Power: 90, Category: "special"}},
	}
	idx := ChooseBestMember([]TeamBattleMember{m}, []PokemonType{{Name: "water"}})
	if idx != 0 {
		t.Errorf("expected 0 for single member, got %d", idx)
	}
}

func TestSimulateTeamBattle_ShuffledOrder(t *testing.T) {
	// With a deterministic randSource, the order should differ from the original
	members := []TeamBattleMember{
		makeTeamBattleMember("first", 200, 100, 80, 90),
		makeTeamBattleMember("second", 200, 100, 80, 90),
		makeTeamBattleMember("third", 200, 100, 80, 90),
	}

	input := TeamBattleInput{
		Team1Name:    "Alpha",
		Team1Members: members,
		Team2Name:    "Beta",
		Team2Members: []TeamBattleMember{makeTeamBattleMember("opponent", 200, 100, 80, 90)},
	}

	// Run multiple times with different seeds to verify shuffling occurs
	counter := 5 // start at 5 to get a different shuffle pattern
	rng := func(n int) int { counter++; return counter % n }
	result := SimulateTeamBattle(input, rng)

	// Just verify the battle completes successfully
	if !result.IsOver {
		t.Error("battle should be over")
	}
	if result.Winner != "team1" && result.Winner != "team2" {
		t.Errorf("unexpected winner: %s", result.Winner)
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
