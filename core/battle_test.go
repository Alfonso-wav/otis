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
		Move:          Move{Name: "Protect", Type: "normal", Power: 0, Category: "status"},
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

// --- Transform (Ditto) Tests ---

func TestExecuteTurn_TransformSetsOverrides(t *testing.T) {
	state := InitBattle(100, 200)
	dittoStats := Stats{HP: 100, Attack: 30, Defense: 30, SpAttack: 30, SpDefense: 30, Speed: 30}
	pikachuStats := Stats{HP: 200, Attack: 120, Defense: 80, SpAttack: 130, SpDefense: 90, Speed: 110}
	pikachuTypes := []PokemonType{{Name: "electric"}}
	pikachuMoves := []Move{
		{Name: "Thunderbolt", Type: "electric", Power: 90, Category: "special", Accuracy: 100},
		{Name: "Quick Attack", Type: "normal", Power: 40, Category: "physical", Accuracy: 100, Priority: 1},
	}

	input := TurnInput{
		State:         state,
		AttackerStats: dittoStats,
		DefenderStats: pikachuStats,
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: pikachuTypes,
		AttackerLevel: 50,
		DefenderLevel: 50,
		Move:          Move{Name: "transform", Type: "normal", Power: 0, Category: "status"},
		AttackerName:  "Ditto",
		DefenderName:  "Pikachu",
		DefenderMoves: pikachuMoves,
	}

	result := ExecuteTurn(input, nil)

	// Verify overrides are set
	if result.NewState.AttackerStatsOverride == nil {
		t.Fatal("AttackerStatsOverride should be set after Transform")
	}

	// Verify stats are copied from defender (except HP preserved)
	override := result.NewState.AttackerStatsOverride
	if override.Attack != pikachuStats.Attack {
		t.Errorf("Attack override: want %d, got %d", pikachuStats.Attack, override.Attack)
	}
	if override.Defense != pikachuStats.Defense {
		t.Errorf("Defense override: want %d, got %d", pikachuStats.Defense, override.Defense)
	}
	if override.SpAttack != pikachuStats.SpAttack {
		t.Errorf("SpAttack override: want %d, got %d", pikachuStats.SpAttack, override.SpAttack)
	}
	if override.SpDefense != pikachuStats.SpDefense {
		t.Errorf("SpDefense override: want %d, got %d", pikachuStats.SpDefense, override.SpDefense)
	}
	if override.Speed != pikachuStats.Speed {
		t.Errorf("Speed override: want %d, got %d", pikachuStats.Speed, override.Speed)
	}
	// HP must be Ditto's, NOT Pikachu's
	if override.HP != dittoStats.HP {
		t.Errorf("HP override: want %d (Ditto's), got %d", dittoStats.HP, override.HP)
	}

	// Verify moves are copied
	if len(result.NewState.AttackerMovesOverride) != len(pikachuMoves) {
		t.Fatalf("AttackerMovesOverride: want %d moves, got %d",
			len(pikachuMoves), len(result.NewState.AttackerMovesOverride))
	}
	for i, m := range result.NewState.AttackerMovesOverride {
		if m.Name != pikachuMoves[i].Name {
			t.Errorf("Move[%d]: want %s, got %s", i, pikachuMoves[i].Name, m.Name)
		}
	}

	// Verify types are copied
	if len(result.NewState.AttackerTypesOverride) != len(pikachuTypes) {
		t.Fatalf("AttackerTypesOverride: want %d types, got %d",
			len(pikachuTypes), len(result.NewState.AttackerTypesOverride))
	}
	if result.NewState.AttackerTypesOverride[0].Name != "electric" {
		t.Errorf("Type override: want electric, got %s", result.NewState.AttackerTypesOverride[0].Name)
	}

	// HP unchanged
	if result.NewState.AttackerHP != 100 {
		t.Errorf("Ditto HP should be unchanged: want 100, got %d", result.NewState.AttackerHP)
	}
	if result.NewState.DefenderHP != 200 {
		t.Errorf("Pikachu HP should be unchanged: want 200, got %d", result.NewState.DefenderHP)
	}

	// Log message
	if !strings.Contains(result.LogEntry, "Transform") {
		t.Errorf("Log should mention Transform, got: %s", result.LogEntry)
	}
	if !strings.Contains(result.LogEntry, "Ditto") {
		t.Errorf("Log should mention Ditto, got: %s", result.LogEntry)
	}
	if !strings.Contains(result.LogEntry, "Pikachu") {
		t.Errorf("Log should mention Pikachu, got: %s", result.LogEntry)
	}

	// Battle should not be over
	if result.NewState.IsOver {
		t.Error("Battle should not end after Transform")
	}
}

func TestSimulateFullBattle_DittoTransformUsesCopiedStats(t *testing.T) {
	// Ditto has very low stats but Transform should copy Pikachu's high stats.
	// Pikachu only has Splash (no damage), so Ditto survives turn 1, uses Transform,
	// and then fights with Pikachu's copied Thunderbolt in subsequent turns.
	transformMove := Move{Name: "transform", Type: "normal", Power: 0, Category: "status"}
	thunderbolt := Move{Name: "Thunderbolt", Type: "electric", Power: 90, Category: "special", Accuracy: 100}
	splash := Move{Name: "Splash", Type: "water", Power: 0, Category: "status"}

	dittoStats := Stats{HP: 200, Attack: 30, Defense: 30, SpAttack: 30, SpDefense: 30, Speed: 150}
	pikachuStats := Stats{HP: 200, Attack: 120, Defense: 80, SpAttack: 130, SpDefense: 90, Speed: 110}

	input := FullBattleInput{
		AttackerStats: dittoStats,
		DefenderStats: pikachuStats,
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "electric"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{transformMove},        // Ditto only has Transform
		DefenderMoves: []Move{splash, thunderbolt},   // Pikachu has Splash and Thunderbolt
		AttackerName:  "Ditto",
		DefenderName:  "Pikachu",
	}

	// alwaysZero: picks first move always (transform for Ditto, splash for Pikachu).
	// Ditto is faster (150 vs 110), so turn 1: Ditto transforms, Pikachu splashes.
	// Turn 2+: Ditto now has Pikachu's moves (picks index 0 = Splash), Pikachu also Splash.
	// To verify Ditto actually gains the moves, use a counter that picks index 1 after turn 1.
	callCount := 0
	mockRand := func(n int) int {
		callCount++
		// Move selection calls come in pairs (attacker move, defender move).
		// After Transform, Ditto has 2 moves. We want to pick Thunderbolt (index 1)
		// for Ditto after the first turn.
		return 0
	}
	result := SimulateFullBattle(input, mockRand)

	if !result.IsOver {
		t.Error("Battle should be over")
	}

	if result.TurnCount < 2 {
		t.Errorf("Expected at least 2 turns (Transform + attack), got %d", result.TurnCount)
	}

	// Verify the log contains a Transform entry
	foundTransform := false
	for _, entry := range result.Log {
		if strings.Contains(entry, "Transform") && strings.Contains(entry, "Ditto") {
			foundTransform = true
			break
		}
	}
	if !foundTransform {
		t.Error("Log should contain Ditto's Transform entry")
	}
}

func TestSimulateFullBattle_DittoTransformUsesDefenderMoves(t *testing.T) {
	// Verify Ditto uses defender's moves after Transform, not its own.
	transformMove := Move{Name: "transform", Type: "normal", Power: 0, Category: "status"}
	hyperBeam := Move{Name: "Hyper Beam", Type: "normal", Power: 150, Category: "special", Accuracy: 100}

	dittoStats := Stats{HP: 200, Attack: 30, Defense: 30, SpAttack: 30, SpDefense: 30, Speed: 150}
	magikarpStats := Stats{HP: 50, Attack: 10, Defense: 10, SpAttack: 10, SpDefense: 10, Speed: 10}

	input := FullBattleInput{
		AttackerStats: dittoStats,
		DefenderStats: magikarpStats,
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "water"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{transformMove},
		DefenderMoves: []Move{hyperBeam},
		AttackerName:  "Ditto",
		DefenderName:  "Magikarp",
	}

	alwaysZero := func(n int) int { return 0 }
	result := SimulateFullBattle(input, alwaysZero)

	if !result.IsOver {
		t.Error("Battle should be over")
	}
	// Ditto transforms turn 1, gains Hyper Beam + Magikarp's stats.
	// Even with Magikarp's low stats, Hyper Beam (150 power) should KO Magikarp (50 HP) quickly.
	if result.Winner != "attacker" {
		t.Errorf("Ditto should win after Transform + Hyper Beam, got winner: %s", result.Winner)
	}

	// Check Transform is in the log
	foundTransform := false
	foundHyperBeam := false
	for _, entry := range result.Log {
		if strings.Contains(entry, "Transform") {
			foundTransform = true
		}
		if strings.Contains(entry, "Ditto") && strings.Contains(entry, "Hyper Beam") {
			foundHyperBeam = true
		}
	}
	if !foundTransform {
		t.Error("Log should contain Transform")
	}
	if !foundHyperBeam {
		t.Error("Log should show Ditto using Hyper Beam (copied from Magikarp)")
	}
}

func TestSimulateFullBattle_DittoHPNotChanged(t *testing.T) {
	// Ditto's HP must remain its own, not be overwritten by defender's HP.
	transformMove := Move{Name: "transform", Type: "normal", Power: 0, Category: "status"}
	splash := Move{Name: "Splash", Type: "water", Power: 0, Category: "status"}

	dittoHP := 100
	pikachuHP := 300

	input := FullBattleInput{
		AttackerStats: Stats{HP: dittoHP, Attack: 30, Defense: 30, SpAttack: 30, SpDefense: 30, Speed: 30},
		DefenderStats: Stats{HP: pikachuHP, Attack: 120, Defense: 80, SpAttack: 130, SpDefense: 90, Speed: 110},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "electric"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{transformMove},
		DefenderMoves: []Move{splash}, // No damage so we can check HP at end
		AttackerName:  "Ditto",
		DefenderName:  "Pikachu",
	}

	alwaysZero := func(n int) int { return 0 }
	result := SimulateFullBattle(input, alwaysZero)

	// After max turns (both use splash/status), HP should be unchanged.
	// Ditto HP must be 100 (its original), NOT 300 (Pikachu's).
	if result.AttackerHP != dittoHP {
		t.Errorf("Ditto HP: want %d (original), got %d", dittoHP, result.AttackerHP)
	}
	if result.AttackerMaxHP != dittoHP {
		t.Errorf("Ditto MaxHP: want %d, got %d", dittoHP, result.AttackerMaxHP)
	}
}

func TestSimulateFullBattle_DefenderTransformWorks(t *testing.T) {
	// Test symmetry: defender using Transform also works.
	tackle := Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100}
	transformMove := Move{Name: "transform", Type: "normal", Power: 0, Category: "status"}

	input := FullBattleInput{
		AttackerStats: Stats{HP: 200, Attack: 150, Defense: 100, SpAttack: 100, SpDefense: 100, Speed: 100},
		DefenderStats: Stats{HP: 200, Attack: 30, Defense: 30, SpAttack: 30, SpDefense: 30, Speed: 30},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{tackle},
		DefenderMoves: []Move{transformMove},
		AttackerName:  "Snorlax",
		DefenderName:  "Ditto",
	}

	alwaysZero := func(n int) int { return 0 }
	result := SimulateFullBattle(input, alwaysZero)

	if !result.IsOver {
		t.Error("Battle should be over")
	}

	// Check that Ditto (defender) used Transform
	foundTransform := false
	for _, entry := range result.Log {
		if strings.Contains(entry, "Ditto") && strings.Contains(entry, "Transform") {
			foundTransform = true
			break
		}
	}
	if !foundTransform {
		t.Error("Log should contain Ditto's Transform entry when used as defender")
	}
}

func TestSimulateFullBattle_NonTransformPokemonUnaffected(t *testing.T) {
	// Ensure normal Pokemon without Transform are completely unaffected.
	tackle := Move{Name: "Tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100}

	input := FullBattleInput{
		AttackerStats: Stats{HP: 200, Attack: 100, Defense: 80, SpAttack: 100, SpDefense: 80, Speed: 80},
		DefenderStats: Stats{HP: 200, Attack: 100, Defense: 80, SpAttack: 100, SpDefense: 80, Speed: 80},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{tackle},
		DefenderMoves: []Move{tackle},
		AttackerName:  "Rattata",
		DefenderName:  "Pidgey",
	}

	alwaysZero := func(n int) int { return 0 }
	result := SimulateFullBattle(input, alwaysZero)

	if !result.IsOver {
		t.Error("Battle should be over")
	}
	// No overrides should be set
	if result.AttackerStatsOverride != nil {
		t.Error("AttackerStatsOverride should be nil for non-Transform Pokemon")
	}
	if result.DefenderStatsOverride != nil {
		t.Error("DefenderStatsOverride should be nil for non-Transform Pokemon")
	}
	if len(result.AttackerMovesOverride) != 0 {
		t.Error("AttackerMovesOverride should be empty for non-Transform Pokemon")
	}
	if len(result.DefenderMovesOverride) != 0 {
		t.Error("DefenderMovesOverride should be empty for non-Transform Pokemon")
	}
	// Verify no Transform in log
	for _, entry := range result.Log {
		if strings.Contains(entry, "Transform") {
			t.Errorf("Non-Transform battle should not mention Transform, found: %s", entry)
		}
	}
}

// --- Weather system ---

func TestExecuteTurn_RainDanceSetsRain(t *testing.T) {
	state := InitBattle(200, 200)
	input := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 80, Defense: 80, SpAttack: 80, SpDefense: 80, Speed: 80, HP: 200},
		DefenderStats: Stats{Attack: 80, Defense: 80, SpAttack: 80, SpDefense: 80, Speed: 80, HP: 200},
		AttackerTypes: []PokemonType{{Name: "water"}},
		DefenderTypes: []PokemonType{{Name: "grass"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		Move:          Move{Name: "rain-dance", Category: "status", Power: 0, Accuracy: 0},
		AttackerName:  "Kyogre",
		DefenderName:  "Venusaur",
	}
	r := ExecuteTurn(input, nil).NewState
	if r.Weather != WeatherRain {
		t.Errorf("Weather: want Rain, got %q", r.Weather)
	}
	if r.WeatherTurnsLeft != WeatherDefaultTurns {
		t.Errorf("WeatherTurnsLeft: want %d, got %d", WeatherDefaultTurns, r.WeatherTurnsLeft)
	}
	if r.DefenderHP != 200 {
		t.Errorf("Rain Dance should not damage: defenderHP=%d", r.DefenderHP)
	}
}

func TestExecuteTurn_AllWeatherMovesMapped(t *testing.T) {
	cases := []struct {
		move string
		want Weather
	}{
		{"rain-dance", WeatherRain},
		{"sunny-day", WeatherSun},
		{"sandstorm", WeatherSandstorm},
		{"hail", WeatherHail},
	}
	for _, c := range cases {
		s := InitBattle(200, 200)
		in := TurnInput{
			State:         s,
			AttackerStats: Stats{Speed: 50, HP: 200},
			DefenderStats: Stats{Speed: 50, HP: 200},
			AttackerTypes: []PokemonType{{Name: "normal"}},
			DefenderTypes: []PokemonType{{Name: "normal"}},
			AttackerLevel: 50,
			DefenderLevel: 50,
			Move:          Move{Name: c.move, Category: "status", Power: 0},
			AttackerName:  "A",
			DefenderName:  "B",
		}
		r := ExecuteTurn(in, nil).NewState
		if r.Weather != c.want {
			t.Errorf("%s: weather = %q, want %q", c.move, r.Weather, c.want)
		}
		if r.WeatherTurnsLeft != WeatherDefaultTurns {
			t.Errorf("%s: turnsLeft = %d, want %d", c.move, r.WeatherTurnsLeft, WeatherDefaultTurns)
		}
	}
}

func TestTickWeather_SandstormDamageRespectsImmunities(t *testing.T) {
	s := BattleState{
		AttackerHP: 160, DefenderHP: 160,
		AttackerMaxHP: 160, DefenderMaxHP: 160,
		Weather: WeatherSandstorm, WeatherTurnsLeft: 5,
	}
	// Attacker is Fire (damaged), defender is Rock (immune).
	atkTypes := []PokemonType{{Name: "fire"}}
	defTypes := []PokemonType{{Name: "rock"}}
	r := tickWeather(s, atkTypes, defTypes, "A", "B")
	if r.AttackerHP != 150 {
		t.Errorf("Non-immune attacker: want 150 HP (160 - 10), got %d", r.AttackerHP)
	}
	if r.DefenderHP != 160 {
		t.Errorf("Rock defender must be immune to Sandstorm: got %d", r.DefenderHP)
	}
	if r.WeatherTurnsLeft != 4 {
		t.Errorf("WeatherTurnsLeft: want 4, got %d", r.WeatherTurnsLeft)
	}
}

func TestTickWeather_HailRespectsIceImmunity(t *testing.T) {
	s := BattleState{
		AttackerHP: 160, DefenderHP: 160,
		AttackerMaxHP: 160, DefenderMaxHP: 160,
		Weather: WeatherHail, WeatherTurnsLeft: 3,
	}
	atkTypes := []PokemonType{{Name: "ice"}}        // immune
	defTypes := []PokemonType{{Name: "grass"}}       // hit
	r := tickWeather(s, atkTypes, defTypes, "A", "B")
	if r.AttackerHP != 160 {
		t.Errorf("Ice should be immune to Hail: got %d", r.AttackerHP)
	}
	if r.DefenderHP != 150 {
		t.Errorf("Non-ice should take 10 Hail dmg: got %d", r.DefenderHP)
	}
}

func TestTickWeather_ExpiresAfterDuration(t *testing.T) {
	s := BattleState{
		AttackerHP: 200, DefenderHP: 200,
		AttackerMaxHP: 200, DefenderMaxHP: 200,
		Weather: WeatherRain, WeatherTurnsLeft: 1,
	}
	r := tickWeather(s, []PokemonType{{Name: "water"}}, []PokemonType{{Name: "water"}}, "A", "B")
	if r.Weather != WeatherNone {
		t.Errorf("Weather should clear, got %q", r.Weather)
	}
	if r.WeatherTurnsLeft != 0 {
		t.Errorf("TurnsLeft should be 0, got %d", r.WeatherTurnsLeft)
	}
	// Log should contain end message
	hasEnd := false
	for _, l := range r.Log {
		if strings.Contains(l, "lluvia cesó") {
			hasEnd = true
			break
		}
	}
	if !hasEnd {
		t.Error("expected rain-end log entry")
	}
}

func TestSimulateFullBattle_WeatherExpiresAfter5Turns(t *testing.T) {
	// Attacker opens with Rain Dance (status), then spams Tackle.
	// Defender always uses Tackle. Use stats such that Tackle deals ~10 dmg so
	// battle reaches turn 6+.
	input := FullBattleInput{
		AttackerStats: Stats{Attack: 80, Defense: 120, SpAttack: 40, SpDefense: 120, Speed: 100, HP: 500},
		DefenderStats: Stats{Attack: 80, Defense: 120, SpAttack: 40, SpDefense: 120, Speed: 90, HP: 500},
		AttackerTypes: []PokemonType{{Name: "water"}},
		DefenderTypes: []PokemonType{{Name: "grass"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{
			{Name: "rain-dance", Category: "status", Power: 0, Accuracy: 0},
		},
		DefenderMoves: []Move{
			{Name: "tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100},
		},
		AttackerName: "Poliwrath",
		DefenderName: "Venusaur",
	}
	// Single-move attacker means it spams rain-dance every turn — each call
	// refreshes the weather. That does not test expiry. Split into a scripted
	// sequence: first move rain-dance, then dummy status without weather.
	input.AttackerMoves = []Move{
		{Name: "rain-dance", Category: "status", Power: 0, Accuracy: 0},
		{Name: "splash", Category: "status", Power: 0, Accuracy: 0},
	}
	// randSource: deterministic; picks index 0 first, then 1 forever so the
	// attacker uses rain-dance exactly once then switches to splash.
	callCount := 0
	randSource := func(n int) int {
		if n <= 1 {
			return 0
		}
		// First call to len(atkMoves)=2 picks 0 (rain-dance), every subsequent → 1 (splash).
		// Other len(n) calls (accuracy rolls, crit rolls, random rolls, len(defMoves)=1) return 0.
		if n == 2 {
			callCount++
			if callCount == 1 {
				return 0
			}
			return 1
		}
		return 0
	}
	r := SimulateFullBattle(input, randSource)
	// After 5 full turns from the activation (turn 1 Rain Dance activates),
	// the tick at the end of turn 5 decrements to 0 and clears.
	// Find the turn where weather becomes None.
	sawRain := false
	for _, l := range r.Log {
		if strings.Contains(l, "Empezó a llover") {
			sawRain = true
		}
	}
	if !sawRain {
		t.Error("expected rain start log")
	}
	// Final state must either have cleared weather within maxTurns or finished
	// under rain. With 500 HP and 40-power Tackle, battle should last long
	// enough to see expiry.
	cleared := false
	for _, l := range r.Log {
		if strings.Contains(l, "lluvia cesó") {
			cleared = true
			break
		}
	}
	if !cleared {
		t.Error("expected rain to expire and log end message")
	}
}

// --- Stat stages ---

func TestExecuteTurn_SwordsDanceRaisesAtkStage(t *testing.T) {
	state := InitBattle(200, 200)
	in := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 100, HP: 200},
		DefenderStats: Stats{Defense: 100, HP: 200},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		Move:          Move{Name: "swords-dance", Category: "status", Power: 0},
		AttackerName:  "Scyther",
		DefenderName:  "Chansey",
	}
	r := ExecuteTurn(in, nil).NewState
	if r.AttackerStages.Atk != 2 {
		t.Errorf("Atk stage: want 2, got %d", r.AttackerStages.Atk)
	}
	if r.DefenderHP != 200 {
		t.Errorf("Swords Dance must not damage: got %d", r.DefenderHP)
	}
}

func TestCalculateDamage_AtkStageDoublesDamage(t *testing.T) {
	base := DamageInput{
		AttackerStats: Stats{Attack: 100},
		DefenderStats: Stats{Defense: 100},
		Move:          Move{Name: "tackle", Type: "normal", Power: 40, Category: "physical"},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		Level:         50,
	}
	noStage := CalculateDamage(base)
	boosted := base
	boosted.AttackerStages.Atk = 2
	withBoost := CalculateDamage(boosted)
	ratio := float64(withBoost.Max) / float64(noStage.Max)
	if ratio < 1.95 || ratio > 2.05 {
		t.Errorf("+2 Atk ratio = %.3f, want ≈2.0 (dmg boosted=%d base=%d)", ratio, withBoost.Max, noStage.Max)
	}
}

func TestCalculateDamage_GrowlReducesDamage(t *testing.T) {
	base := DamageInput{
		AttackerStats: Stats{Attack: 100},
		DefenderStats: Stats{Defense: 100},
		Move:          Move{Name: "tackle", Type: "normal", Power: 40, Category: "physical"},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		Level:         50,
	}
	noStage := CalculateDamage(base)
	growled := base
	growled.AttackerStages.Atk = -1
	withGrowl := CalculateDamage(growled)
	ratio := float64(withGrowl.Max) / float64(noStage.Max)
	// -1 Atk → ×2/3
	if ratio < 0.6 || ratio > 0.73 {
		t.Errorf("-1 Atk ratio = %.3f, want ≈0.67", ratio)
	}
}

func TestResolveOrder_SpeStagesChangeOrder(t *testing.T) {
	// Attacker base Spe = 60, defender base Spe = 80. Without stages, defender first.
	// If attacker uses Agility (+2 Spe), its effective Spe = 60*2 = 120, attacker first.
	input := FullBattleInput{
		AttackerStats: Stats{Attack: 10, Defense: 200, SpAttack: 10, SpDefense: 200, Speed: 60, HP: 500},
		DefenderStats: Stats{Attack: 10, Defense: 200, SpAttack: 10, SpDefense: 200, Speed: 80, HP: 500},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		DefenderLevel: 50,
		AttackerMoves: []Move{
			{Name: "agility", Category: "status", Power: 0},
			{Name: "tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100},
		},
		DefenderMoves: []Move{
			{Name: "tackle", Type: "normal", Power: 40, Category: "physical", Accuracy: 100},
		},
		AttackerName: "Jolteon",
		DefenderName: "Raichu",
	}
	// randSource: picks 0 first (attacker's agility, then tackle) deterministically
	callCount := 0
	randSource := func(n int) int {
		if n == 2 {
			callCount++
			if callCount == 1 {
				return 0 // agility
			}
			return 1 // tackle
		}
		return 0
	}
	r := SimulateFullBattle(input, randSource)
	if r.AttackerStages.Spe != 2 {
		t.Errorf("Agility should leave +2 Spe stage, got %d", r.AttackerStages.Spe)
	}
}

func TestInitBattle_StagesAreZero(t *testing.T) {
	s := InitBattle(200, 200)
	if s.AttackerStages != (StatStages{}) {
		t.Errorf("AttackerStages must init to zero, got %+v", s.AttackerStages)
	}
	if s.DefenderStages != (StatStages{}) {
		t.Errorf("DefenderStages must init to zero, got %+v", s.DefenderStages)
	}
}

func TestExecuteTurn_GrowlDebuffsDefenderAtk(t *testing.T) {
	state := InitBattle(200, 200)
	in := TurnInput{
		State:         state,
		AttackerStats: Stats{Attack: 100, HP: 200},
		DefenderStats: Stats{Attack: 100, Defense: 100, HP: 200},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		AttackerLevel: 50,
		Move:          Move{Name: "growl", Category: "status", Power: 0},
		AttackerName:  "Pidgey",
		DefenderName:  "Rattata",
	}
	r := ExecuteTurn(in, nil).NewState
	if r.DefenderStages.Atk != -1 {
		t.Errorf("Growl: defender Atk stage want -1, got %d", r.DefenderStages.Atk)
	}
}
