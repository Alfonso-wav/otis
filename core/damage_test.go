package core

import (
	"math"
	"testing"
)

func TestTypeEffectiveness_SuperEffective(t *testing.T) {
	// Agua vs Fuego = ×2
	mult := TypeEffectiveness("water", []PokemonType{{Name: "fire"}})
	if mult != 2.0 {
		t.Errorf("water vs fire = %.2f, want 2.0", mult)
	}
}

func TestTypeEffectiveness_NotVeryEffective(t *testing.T) {
	// Fuego vs Agua = ×0.5
	mult := TypeEffectiveness("fire", []PokemonType{{Name: "water"}})
	if mult != 0.5 {
		t.Errorf("fire vs water = %.2f, want 0.5", mult)
	}
}

func TestTypeEffectiveness_NoEffect(t *testing.T) {
	// Normal vs Ghost = ×0
	mult := TypeEffectiveness("normal", []PokemonType{{Name: "ghost"}})
	if mult != 0.0 {
		t.Errorf("normal vs ghost = %.2f, want 0.0", mult)
	}
}

func TestTypeEffectiveness_Neutral(t *testing.T) {
	// Normal vs Normal = ×1
	mult := TypeEffectiveness("normal", []PokemonType{{Name: "normal"}})
	if mult != 1.0 {
		t.Errorf("normal vs normal = %.2f, want 1.0", mult)
	}
}

func TestTypeEffectiveness_DualType(t *testing.T) {
	// Eléctrico vs Agua+Volador = ×2 × ×2 = ×4
	mult := TypeEffectiveness("electric", []PokemonType{{Name: "water"}, {Name: "flying"}})
	if mult != 4.0 {
		t.Errorf("electric vs water/flying = %.2f, want 4.0", mult)
	}
}

func TestTypeEffectiveness_DualTypeImmune(t *testing.T) {
	// Eléctrico vs Agua+Tierra = ×2 × ×0 = ×0
	mult := TypeEffectiveness("electric", []PokemonType{{Name: "water"}, {Name: "ground"}})
	if mult != 0.0 {
		t.Errorf("electric vs water/ground = %.2f, want 0.0", mult)
	}
}

func TestCalculateDamage_StatusMove(t *testing.T) {
	input := DamageInput{
		Move:  Move{Name: "toxic", Category: "status", Power: 0},
		Level: 50,
	}
	result := CalculateDamage(input)
	if result.Min != 0 || result.Max != 0 {
		t.Errorf("status move should do 0 damage, got min=%d max=%d", result.Min, result.Max)
	}
}

func TestCalculateDamage_PhysicalSTAB(t *testing.T) {
	// Pikachu (Eléctrico) usando Thunderbolt (Especial) vs Slowpoke (Agua/Psíquico)
	// Simplificado con atk=100, def=100, pwr=90, nivel=50
	input := DamageInput{
		AttackerStats: Stats{SpAttack: 100},
		DefenderStats: Stats{SpDefense: 100},
		Move:          Move{Name: "thunderbolt", Type: "electric", Category: "special", Power: 90},
		AttackerTypes: []PokemonType{{Name: "electric"}},
		DefenderTypes: []PokemonType{{Name: "water"}},
		Level:         50,
		IsCritical:    false,
		WeatherBonus:  1.0,
	}
	result := CalculateDamage(input)

	// base = floor(floor(2*50/5+2) * 90 * 100/100 / 50) + 2
	//      = floor(floor(22) * 90 / 50) + 2
	//      = floor(39.6) + 2 = 39 + 2 = 41
	// STAB = 1.5 (electric attacker, electric move)
	// TypeEff = 2.0 (electric vs water)
	// modifier = 1.5 * 2.0 = 3.0
	// min = floor(41 * 3.0 * 0.85) = floor(104.55) = 104
	// max = floor(41 * 3.0 * 1.00) = floor(123) = 123
	baseDmg := math.Floor(float64(2*50)/5+2)*90*100/100/50
	baseDmg = math.Floor(baseDmg) + 2
	wantMin := int(math.Floor(baseDmg * 1.5 * 2.0 * 0.85))
	wantMax := int(math.Floor(baseDmg * 1.5 * 2.0 * 1.00))

	if result.Min != wantMin {
		t.Errorf("min damage = %d, want %d", result.Min, wantMin)
	}
	if result.Max != wantMax {
		t.Errorf("max damage = %d, want %d", result.Max, wantMax)
	}
	if !result.IsSuperEffective {
		t.Error("expected IsSuperEffective = true")
	}
	if result.Multiplier != 2.0 {
		t.Errorf("multiplier = %.2f, want 2.0", result.Multiplier)
	}
}

func TestCalculateDamage_Critical(t *testing.T) {
	base := DamageInput{
		AttackerStats: Stats{Attack: 100},
		DefenderStats: Stats{Defense: 100},
		Move:          Move{Name: "tackle", Type: "normal", Category: "physical", Power: 40},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		Level:         50,
		IsCritical:    false,
		WeatherBonus:  1.0,
	}
	noCrit := CalculateDamage(base)
	base.IsCritical = true
	crit := CalculateDamage(base)

	if crit.Max <= noCrit.Max {
		t.Errorf("critical max (%d) should be greater than normal max (%d)", crit.Max, noCrit.Max)
	}
}

func TestCalculateDamage_NoEffect(t *testing.T) {
	input := DamageInput{
		AttackerStats: Stats{Attack: 100},
		DefenderStats: Stats{Defense: 100},
		Move:          Move{Name: "tackle", Type: "normal", Category: "physical", Power: 40},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "ghost"}},
		Level:         50,
		WeatherBonus:  1.0,
	}
	result := CalculateDamage(input)
	if !result.HasNoEffect {
		t.Error("expected HasNoEffect = true for normal vs ghost")
	}
	if result.Min != 0 || result.Max != 0 {
		t.Errorf("no effect move should do 0, got min=%d max=%d", result.Min, result.Max)
	}
}

func TestCalculateDamage_NotVeryEffective(t *testing.T) {
	input := DamageInput{
		AttackerStats: Stats{SpAttack: 100},
		DefenderStats: Stats{SpDefense: 100},
		Move:          Move{Name: "flamethrower", Type: "fire", Category: "special", Power: 90},
		AttackerTypes: []PokemonType{{Name: "fire"}},
		DefenderTypes: []PokemonType{{Name: "water"}},
		Level:         50,
		WeatherBonus:  1.0,
	}
	result := CalculateDamage(input)
	if !result.IsNotVeryEffective {
		t.Error("expected IsNotVeryEffective = true for fire vs water")
	}
	if result.Multiplier != 0.5 {
		t.Errorf("multiplier = %.2f, want 0.5", result.Multiplier)
	}
}

func TestCalculateBattleDamage_RandomRoll(t *testing.T) {
	input := DamageInput{
		AttackerStats: Stats{Attack: 100},
		DefenderStats: Stats{Defense: 100},
		Move:          Move{Name: "tackle", Type: "normal", Category: "physical", Power: 40},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		Level:         50,
		WeatherBonus:  1.0,
	}

	detResult := CalculateDamage(input)

	// Roll=0 → factor 85/100 (min), crit roll=5 (no crit at stage 0)
	callCount := 0
	minRoll := func(n int) int {
		callCount++
		switch callCount {
		case 1:
			return 5 // crit: 5 != 0 → no crit
		case 2:
			return 0 // roll: 85+0=85 → min roll
		}
		return 0
	}
	minResult := CalculateBattleDamage(input, minRoll)

	// Roll=15 → factor 100/100 (max), crit roll=5 (no crit)
	callCount = 0
	maxRoll := func(n int) int {
		callCount++
		switch callCount {
		case 1:
			return 5 // crit: no crit
		case 2:
			return 15 // roll: 85+15=100 → max roll
		}
		return 0
	}
	maxResult := CalculateBattleDamage(input, maxRoll)

	if minResult.ActualDamage > maxResult.ActualDamage {
		t.Errorf("min roll (%d) should be <= max roll (%d)", minResult.ActualDamage, maxResult.ActualDamage)
	}
	if maxResult.ActualDamage != detResult.Max {
		t.Errorf("max roll actual (%d) should equal deterministic max (%d)", maxResult.ActualDamage, detResult.Max)
	}
}

func TestCalculateBattleDamage_ProbabilisticCrit(t *testing.T) {
	input := DamageInput{
		AttackerStats: Stats{Attack: 100},
		DefenderStats: Stats{Defense: 100},
		Move:          Move{Name: "tackle", Type: "normal", Category: "physical", Power: 40},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		Level:         50,
		WeatherBonus:  1.0,
		CriticalStage: 0,
	}

	// Crit roll = 0 → crit (1/24 threshold, 0 == 0 → crit)
	critRand := func(n int) int { return 0 }
	critResult := CalculateBattleDamage(input, critRand)
	if !critResult.WasCritical {
		t.Error("Expected critical hit when crit roll is 0 at stage 0")
	}

	// Crit roll = 5 → no crit (5 != 0)
	noCritRand := func(n int) int { return 5 }
	noCritResult := CalculateBattleDamage(input, noCritRand)
	if noCritResult.WasCritical {
		t.Error("Expected no critical hit when crit roll is 5 at stage 0")
	}

	if critResult.ActualDamage <= noCritResult.ActualDamage {
		t.Errorf("crit damage (%d) should be > non-crit damage (%d)", critResult.ActualDamage, noCritResult.ActualDamage)
	}
}

func TestCalculateBattleDamage_CritStage3AlwaysCrits(t *testing.T) {
	input := DamageInput{
		AttackerStats: Stats{Attack: 100},
		DefenderStats: Stats{Defense: 100},
		Move:          Move{Name: "tackle", Type: "normal", Category: "physical", Power: 40},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		Level:         50,
		WeatherBonus:  1.0,
		CriticalStage: 3,
	}

	// At stage 3, threshold is 1, so randSource(1) always returns 0 → always crit
	anyRand := func(n int) int { return 0 }
	result := CalculateBattleDamage(input, anyRand)
	if !result.WasCritical {
		t.Error("Expected always critical at stage 3")
	}
}

func TestCheckAccuracy(t *testing.T) {
	alwaysHit := func(n int) int { return 0 }
	alwaysMiss := func(n int) int { return 99 }

	// Accuracy 100 → always hits
	if !CheckAccuracy(100, alwaysMiss) {
		t.Error("Accuracy 100 should always hit")
	}

	// Accuracy 0 → always hits (never-miss)
	if !CheckAccuracy(0, alwaysMiss) {
		t.Error("Accuracy 0 should always hit (never-miss)")
	}

	// Accuracy 70, roll 0 → hit
	if !CheckAccuracy(70, alwaysHit) {
		t.Error("Accuracy 70, roll 0 should hit")
	}

	// Accuracy 70, roll 99 → miss
	if CheckAccuracy(70, alwaysMiss) {
		t.Error("Accuracy 70, roll 99 should miss")
	}
}
