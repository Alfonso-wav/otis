package core

import (
	"math"
	"testing"
)

func TestGetAbility_NormalizesAndMatchesKnown(t *testing.T) {
	ab, ok := GetAbility("Drizzle")
	if !ok || ab.Name != "drizzle" {
		t.Errorf("Drizzle: want name=drizzle ok=true, got %+v ok=%v", ab, ok)
	}
	ab2, ok2 := GetAbility("sand stream")
	if !ok2 || ab2.Name != "sand-stream" {
		t.Errorf("sand stream: want sand-stream, got %+v ok=%v", ab2, ok2)
	}
	if _, ok3 := GetAbility(""); ok3 {
		t.Error("empty name must return ok=false")
	}
	if _, ok4 := GetAbility("not-a-real-ability"); ok4 {
		t.Error("unknown name must return ok=false")
	}
}

func TestAbilityCatalog_MinimumSize(t *testing.T) {
	// Task 0168 requires ≥80 abilities registered.
	if n := len(AbilityNames()); n < 80 {
		t.Errorf("ability catalog size = %d, want ≥80", n)
	}
}

// --- Weather on switch-in ---

func TestDrizzle_SetsRainOnSwitchIn(t *testing.T) {
	s := InitBattle(200, 200)
	ab, _ := GetAbility("drizzle")
	s = ab.OnSwitchIn(s, SideAttacker)
	if s.Weather != WeatherRain {
		t.Errorf("Drizzle weather: want Rain, got %q", s.Weather)
	}
	if s.WeatherTurnsLeft != WeatherDefaultTurns {
		t.Errorf("Drizzle turns: want %d, got %d", WeatherDefaultTurns, s.WeatherTurnsLeft)
	}
}

func TestDroughtSnowSandStream_SetWeather(t *testing.T) {
	cases := []struct {
		ab   string
		want Weather
	}{
		{"drought", WeatherSun},
		{"sand-stream", WeatherSandstorm},
		{"snow-warning", WeatherHail},
	}
	for _, c := range cases {
		s := InitBattle(100, 100)
		ab, _ := GetAbility(c.ab)
		s = ab.OnSwitchIn(s, SideAttacker)
		if s.Weather != c.want {
			t.Errorf("%s: weather = %q, want %q", c.ab, s.Weather, c.want)
		}
	}
}

// --- Intimidate ---

func TestIntimidate_DropsOpponentAtk(t *testing.T) {
	s := InitBattle(100, 100)
	ab, _ := GetAbility("intimidate")
	s = ab.OnSwitchIn(s, SideAttacker)
	if s.DefenderStages.Atk != -1 {
		t.Errorf("Intimidate (attacker side): defender Atk want -1, got %d", s.DefenderStages.Atk)
	}
	s2 := InitBattle(100, 100)
	s2 = ab.OnSwitchIn(s2, SideDefender)
	if s2.AttackerStages.Atk != -1 {
		t.Errorf("Intimidate (defender side): attacker Atk want -1, got %d", s2.AttackerStages.Atk)
	}
}

// --- Water Absorb: immunity + heal 25% ---

func TestWaterAbsorb_HealsAndNullifies(t *testing.T) {
	s := InitBattle(200, 200)
	s.DefenderHP = 100 // wounded
	ab, _ := GetAbility("water-absorb")
	if ab.ImmuneToType == nil || !ab.ImmuneToType("water") {
		t.Error("water-absorb must be immune to water")
	}
	s2, dmg := ab.OnTakeHit(s, SideDefender, Move{Type: "water"}, 50)
	if dmg != 0 {
		t.Errorf("water-absorb damage: want 0, got %d", dmg)
	}
	// Heal 200/4 = 50
	if s2.DefenderHP != 150 {
		t.Errorf("water-absorb heal: want HP 150, got %d", s2.DefenderHP)
	}
}

func TestWaterAbsorb_DoesNotAffectOtherTypes(t *testing.T) {
	s := InitBattle(200, 200)
	s.DefenderHP = 100
	ab, _ := GetAbility("water-absorb")
	s2, dmg := ab.OnTakeHit(s, SideDefender, Move{Type: "fire"}, 50)
	if dmg != 50 || s2.DefenderHP != 100 {
		t.Errorf("fire move through water-absorb should be unaffected, got dmg=%d hp=%d", dmg, s2.DefenderHP)
	}
}

// --- Pinch boosts (Blaze/Torrent/Overgrow/Swarm) ---

func TestBlaze_BoostsFireOnlyWhenLowHP(t *testing.T) {
	ab, _ := GetAbility("blaze")
	state := BattleState{AttackerHP: 100, AttackerMaxHP: 300, DefenderMaxHP: 100}
	ctx := DamageCtx{State: state, Side: SideAttacker, Move: Move{Type: "fire", Power: 40}}
	// 100/300 = 0.33 → not strictly < 1/3. Drop further.
	state.AttackerHP = 80
	ctx.State = state
	got := ab.ModifyAttack(ctx)
	if math.Abs(got-1.5) > 1e-9 {
		t.Errorf("Blaze with low HP and fire: want 1.5, got %v", got)
	}
	state.AttackerHP = 200
	ctx.State = state
	got = ab.ModifyAttack(ctx)
	if got != 1.0 {
		t.Errorf("Blaze with high HP: want 1.0, got %v", got)
	}
	// Wrong move type → no boost even in pinch.
	state.AttackerHP = 80
	ctx.State = state
	ctx.Move = Move{Type: "water", Power: 40}
	got = ab.ModifyAttack(ctx)
	if got != 1.0 {
		t.Errorf("Blaze with non-fire move: want 1.0, got %v", got)
	}
}

func TestTorrentOvergrowSwarm_PinchBoost(t *testing.T) {
	cases := []struct {
		ab   string
		mType string
	}{
		{"torrent", "water"},
		{"overgrow", "grass"},
		{"swarm", "bug"},
	}
	for _, c := range cases {
		ab, _ := GetAbility(c.ab)
		state := BattleState{AttackerHP: 80, AttackerMaxHP: 300}
		ctx := DamageCtx{State: state, Side: SideAttacker, Move: Move{Type: c.mType, Power: 40}}
		if got := ab.ModifyAttack(ctx); got != 1.5 {
			t.Errorf("%s pinch: want 1.5, got %v", c.ab, got)
		}
	}
}

// --- Damage reduction ---

func TestThickFat_HalvesFireAndIce(t *testing.T) {
	ab, _ := GetAbility("thick-fat")
	ctx := DamageCtx{Move: Move{Type: "fire"}}
	if ab.ModifyDefense(ctx) != 0.5 {
		t.Error("thick-fat fire: want 0.5")
	}
	ctx.Move.Type = "ice"
	if ab.ModifyDefense(ctx) != 0.5 {
		t.Error("thick-fat ice: want 0.5")
	}
	ctx.Move.Type = "water"
	if ab.ModifyDefense(ctx) != 1.0 {
		t.Error("thick-fat water: want 1.0")
	}
}

func TestFilterSolidRock_ReduceSuperEffective(t *testing.T) {
	for _, name := range []string{"filter", "solid-rock", "prism-armor"} {
		ab, _ := GetAbility(name)
		ctx := DamageCtx{TypeEffectiveness: 2.0}
		if ab.ModifyDefense(ctx) != 0.75 {
			t.Errorf("%s SE: want 0.75", name)
		}
		ctx.TypeEffectiveness = 0.5
		if ab.ModifyDefense(ctx) != 1.0 {
			t.Errorf("%s NVE: want 1.0", name)
		}
	}
}

func TestMultiscale_ReducesAtFullHP(t *testing.T) {
	ab, _ := GetAbility("multiscale")
	state := BattleState{DefenderHP: 200, DefenderMaxHP: 200}
	if ab.ModifyDefense(DamageCtx{State: state, Side: SideDefender}) != 0.5 {
		t.Error("Multiscale full HP: want 0.5")
	}
	state.DefenderHP = 100
	if ab.ModifyDefense(DamageCtx{State: state, Side: SideDefender}) != 1.0 {
		t.Error("Multiscale half HP: want 1.0")
	}
}

// --- Offensive multipliers ---

func TestHugePower_DoublesPhysical(t *testing.T) {
	ab, _ := GetAbility("huge-power")
	ctx := DamageCtx{Move: Move{Category: "physical", Power: 40}}
	if ab.ModifyAttack(ctx) != 2.0 {
		t.Error("Huge Power physical: want 2.0")
	}
	ctx.Move.Category = "special"
	if ab.ModifyAttack(ctx) != 1.0 {
		t.Error("Huge Power special: want 1.0")
	}
}

func TestTechnician_BoostsLowPowerMoves(t *testing.T) {
	ab, _ := GetAbility("technician")
	if ab.ModifyAttack(DamageCtx{Move: Move{Power: 60}}) != 1.5 {
		t.Error("Technician Power 60: want 1.5")
	}
	if ab.ModifyAttack(DamageCtx{Move: Move{Power: 61}}) != 1.0 {
		t.Error("Technician Power 61: want 1.0")
	}
}

func TestAdaptability_BoostsSTAB(t *testing.T) {
	ab, _ := GetAbility("adaptability")
	ctx := DamageCtx{
		Move:          Move{Type: "water"},
		AttackerTypes: []PokemonType{{Name: "water"}},
	}
	got := ab.ModifyAttack(ctx)
	if math.Abs(got-2.0/1.5) > 1e-9 {
		t.Errorf("Adaptability STAB: want 2.0/1.5, got %v", got)
	}
	ctx.AttackerTypes = []PokemonType{{Name: "fire"}}
	if ab.ModifyAttack(ctx) != 1.0 {
		t.Error("Adaptability non-STAB: want 1.0")
	}
}

func TestTintedLens_DoublesNotVeryEffective(t *testing.T) {
	ab, _ := GetAbility("tinted-lens")
	if ab.ModifyAttack(DamageCtx{TypeEffectiveness: 0.5}) != 2.0 {
		t.Error("Tinted Lens NVE: want 2.0")
	}
	if ab.ModifyAttack(DamageCtx{TypeEffectiveness: 2.0}) != 1.0 {
		t.Error("Tinted Lens SE: want 1.0")
	}
}

// --- Speed-weather abilities ---

func TestSwiftSwim_DoublesInRain(t *testing.T) {
	ab, _ := GetAbility("swift-swim")
	s := BattleState{Weather: WeatherRain, WeatherTurnsLeft: 5}
	if ab.SpeedMultiplier(s, SideAttacker) != 2.0 {
		t.Error("Swift Swim in Rain: want 2.0")
	}
	s.Weather = WeatherNone
	if ab.SpeedMultiplier(s, SideAttacker) != 1.0 {
		t.Error("Swift Swim no weather: want 1.0")
	}
}

// --- Cloud Nine / Air Lock ---

func TestCloudNine_DisablesWeatherDamageModifier(t *testing.T) {
	// With Rain active + Cloud Nine on one side, WeatherDamageMultiplier
	// must NOT apply in CalculateDamage.
	atkInput := DamageInput{
		AttackerStats: Stats{SpAttack: 100},
		DefenderStats: Stats{SpDefense: 100},
		Move:          Move{Name: "surf", Type: "water", Power: 90, Category: "special"},
		AttackerTypes: []PokemonType{{Name: "normal"}},
		DefenderTypes: []PokemonType{{Name: "normal"}},
		Level:         50,
		Weather:       WeatherRain,
		State:         BattleState{AttackerIgnoresWeather: true},
	}
	withRain := atkInput
	withRain.State = BattleState{}
	noRainEffect := CalculateDamage(atkInput)
	rainy := CalculateDamage(withRain)
	if rainy.Max <= noRainEffect.Max {
		t.Errorf("Cloud Nine should suppress Rain boost: cloud9=%d rain=%d", noRainEffect.Max, rainy.Max)
	}
}

// --- Flash Fire ---

func TestFlashFire_ActivatesOnFireHitAndBoostsFire(t *testing.T) {
	ab, _ := GetAbility("flash-fire")
	s := BattleState{DefenderHP: 100, DefenderMaxHP: 100}
	s2, dmg := ab.OnTakeHit(s, SideDefender, Move{Type: "fire"}, 50)
	if dmg != 0 {
		t.Errorf("Flash Fire fire damage: want 0, got %d", dmg)
	}
	if !s2.DefenderFlashFireActive {
		t.Error("Flash Fire flag should activate on defender")
	}
	// Boost only when activated + fire move.
	ctx := DamageCtx{
		State: s2, Side: SideDefender, Move: Move{Type: "fire"},
	}
	if got := ab.ModifyAttack(ctx); got != 1.5 {
		t.Errorf("Flash Fire boost: want 1.5, got %v", got)
	}
}

// --- Wonder Guard ---

func TestWonderGuard_OnlySuperEffectiveDamages(t *testing.T) {
	ab, _ := GetAbility("wonder-guard")
	// Neutral hit
	if ab.ModifyDefense(DamageCtx{Move: Move{Power: 40}, TypeEffectiveness: 1.0}) != 0.0 {
		t.Error("Wonder Guard neutral: want 0.0")
	}
	// Super effective hit
	if ab.ModifyDefense(DamageCtx{Move: Move{Power: 40}, TypeEffectiveness: 2.0}) != 1.0 {
		t.Error("Wonder Guard SE: want 1.0")
	}
}

// --- Speed Boost ---

func TestSpeedBoost_EndOfTurnPlus1Spe(t *testing.T) {
	ab, _ := GetAbility("speed-boost")
	s := BattleState{}
	s = ab.EndOfTurn(s, SideAttacker)
	if s.AttackerStages.Spe != 1 {
		t.Errorf("Speed Boost: want +1 Spe, got %d", s.AttackerStages.Spe)
	}
}

// --- Sap Sipper / Motor Drive ---

func TestSapSipper_ImmuneGrassPlusAtkBoost(t *testing.T) {
	ab, _ := GetAbility("sap-sipper")
	s := InitBattle(100, 100)
	s2, dmg := ab.OnTakeHit(s, SideDefender, Move{Type: "grass"}, 40)
	if dmg != 0 || s2.DefenderStages.Atk != 1 {
		t.Errorf("Sap Sipper: want dmg=0 atk+1, got dmg=%d atk=%d", dmg, s2.DefenderStages.Atk)
	}
}

func TestMotorDrive_ImmuneElectricPlusSpeBoost(t *testing.T) {
	ab, _ := GetAbility("motor-drive")
	s := InitBattle(100, 100)
	s2, dmg := ab.OnTakeHit(s, SideDefender, Move{Type: "electric"}, 40)
	if dmg != 0 || s2.DefenderStages.Spe != 1 {
		t.Errorf("Motor Drive: want dmg=0 spe+1, got dmg=%d spe=%d", dmg, s2.DefenderStages.Spe)
	}
}
