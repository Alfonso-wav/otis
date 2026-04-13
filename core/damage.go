package core

import "math"

// DamageInput contiene todos los parámetros necesarios para calcular daño
type DamageInput struct {
	AttackerStats Stats         `json:"attackerStats"`
	DefenderStats Stats         `json:"defenderStats"`
	Move          Move          `json:"move"`
	AttackerTypes []PokemonType `json:"attackerTypes"`
	DefenderTypes []PokemonType `json:"defenderTypes"`
	Level         int           `json:"level"`
	IsCritical    bool          `json:"isCritical"`
	CriticalStage int           `json:"criticalStage"`
	// WeatherBonus is a direct numeric multiplier override (legacy). Weather is preferred.
	WeatherBonus float64 `json:"weatherBonus"`
	// Weather active during the move. If set, overrides WeatherBonus via WeatherDamageMultiplier.
	Weather  Weather `json:"weather,omitempty"`
	IsBurned bool    `json:"isBurned"`
	// Stat stages of attacker/defender. Zero value = no stage modifiers.
	AttackerStages StatStages `json:"attackerStages,omitempty"`
	DefenderStages StatStages `json:"defenderStages,omitempty"`
	// Ability slugs. Empty string = no ability.
	AttackerAbility string `json:"attackerAbility,omitempty"`
	DefenderAbility string `json:"defenderAbility,omitempty"`
	// State snapshot so ability hooks can query HP ratios / flags.
	State BattleState `json:"-"`
}

// DamageResult contiene el resultado del cálculo de daño
type DamageResult struct {
	Min                int     `json:"min"`
	Max                int     `json:"max"`
	Average            int     `json:"average"`
	ActualDamage       int     `json:"actualDamage"`
	Multiplier         float64 `json:"multiplier"`
	IsSuperEffective   bool    `json:"isSuperEffective"`
	IsNotVeryEffective bool    `json:"isNotVeryEffective"`
	HasNoEffect        bool    `json:"hasNoEffect"`
	WasCritical        bool    `json:"wasCritical"`
	HasSTAB            bool    `json:"hasSTAB"`
	STABMultiplier     float64 `json:"stabMultiplier"`
	BurnApplied        bool    `json:"burnApplied"`
}

// typeChart mapea [atacante][defensor] → multiplicador (solo entradas no-1.0)
var typeChart = map[string]map[string]float64{
	"normal": {
		"rock": 0.5, "ghost": 0, "steel": 0.5,
	},
	"fire": {
		"fire": 0.5, "water": 0.5, "grass": 2, "ice": 2,
		"bug": 2, "rock": 0.5, "dragon": 0.5, "steel": 2,
	},
	"water": {
		"fire": 2, "water": 0.5, "grass": 0.5, "ground": 2,
		"rock": 2, "dragon": 0.5,
	},
	"electric": {
		"water": 2, "electric": 0.5, "grass": 0.5, "ground": 0,
		"flying": 2, "dragon": 0.5,
	},
	"grass": {
		"fire": 0.5, "water": 2, "grass": 0.5, "poison": 0.5,
		"ground": 2, "flying": 0.5, "bug": 0.5, "rock": 2,
		"dragon": 0.5, "steel": 0.5,
	},
	"ice": {
		"water": 0.5, "grass": 2, "ice": 0.5, "ground": 2,
		"flying": 2, "dragon": 2, "steel": 0.5,
	},
	"fighting": {
		"normal": 2, "ice": 2, "poison": 0.5, "flying": 0.5,
		"psychic": 0.5, "bug": 0.5, "rock": 2, "ghost": 0,
		"dark": 2, "steel": 2, "fairy": 0.5,
	},
	"poison": {
		"grass": 2, "poison": 0.5, "ground": 0.5, "rock": 0.5,
		"ghost": 0.5, "steel": 0, "fairy": 2,
	},
	"ground": {
		"fire": 2, "electric": 2, "grass": 0.5, "poison": 2,
		"flying": 0, "bug": 0.5, "rock": 2, "steel": 2,
	},
	"flying": {
		"electric": 0.5, "grass": 2, "fighting": 2, "bug": 2,
		"rock": 0.5, "steel": 0.5,
	},
	"psychic": {
		"fighting": 2, "poison": 2, "psychic": 0.5, "dark": 0,
		"steel": 0.5,
	},
	"bug": {
		"fire": 0.5, "grass": 2, "fighting": 0.5, "poison": 0.5,
		"flying": 0.5, "psychic": 2, "ghost": 0.5, "dark": 2,
		"steel": 0.5, "fairy": 0.5,
	},
	"rock": {
		"fire": 2, "ice": 2, "fighting": 0.5, "ground": 0.5,
		"flying": 2, "bug": 2, "steel": 0.5,
	},
	"ghost": {
		"normal": 0, "psychic": 2, "ghost": 2, "dark": 0.5,
	},
	"dragon": {
		"dragon": 2, "steel": 0.5, "fairy": 0,
	},
	"dark": {
		"fighting": 0.5, "psychic": 2, "ghost": 2, "dark": 0.5,
		"fairy": 0.5,
	},
	"steel": {
		"fire": 0.5, "water": 0.5, "electric": 0.5, "ice": 2,
		"rock": 2, "steel": 0.5, "fairy": 2,
	},
	"fairy": {
		"fire": 0.5, "fighting": 2, "poison": 0.5, "dragon": 2,
		"dark": 2, "steel": 0.5,
	},
}

// TypeEffectiveness calcula el multiplicador de efectividad de tipo
func TypeEffectiveness(moveType string, defenderTypes []PokemonType) float64 {
	mult := 1.0
	defChart := typeChart[moveType]
	for _, dt := range defenderTypes {
		if defChart != nil {
			if v, ok := defChart[dt.Name]; ok {
				mult *= v
				continue
			}
		}
		mult *= 1.0
	}
	return mult
}

// WeatherDamageMultiplier returns the damage multiplier on a move given active weather.
// Rain: Water ×1.5, Fire ×0.5. Sun: Fire ×1.5, Water ×0.5.
// Hail/Sandstorm: no damage multiplier (Sandstorm SpD boost is applied separately).
func WeatherDamageMultiplier(weather Weather, moveType string) float64 {
	switch weather {
	case WeatherRain:
		switch moveType {
		case "water":
			return 1.5
		case "fire":
			return 0.5
		}
	case WeatherSun:
		switch moveType {
		case "fire":
			return 1.5
		case "water":
			return 0.5
		}
	}
	return 1.0
}

// hasType reports whether types contains a type with the given name.
func hasType(types []PokemonType, name string) bool {
	for _, t := range types {
		if t.Name == name {
			return true
		}
	}
	return false
}

// CalculateDamage calcula el rango de daño usando la fórmula Gen 5+
func CalculateDamage(input DamageInput) DamageResult {
	if input.Move.Power <= 0 || input.Move.Category == "status" {
		return DamageResult{}
	}

	// Ability lookups (empty name → zero value / no-op).
	atkAbility, _ := GetAbility(input.AttackerAbility)
	defAbility, _ := GetAbility(input.DefenderAbility)

	// Defender's ability may grant complete type immunity (Levitate, Water Absorb…).
	if defAbility.ImmuneToType != nil && defAbility.ImmuneToType(input.Move.Type) {
		return DamageResult{HasNoEffect: true}
	}

	// Weather multiplier respects cloud-nine / air-lock flags on either side.
	effectiveWeather := input.Weather
	if input.State.AttackerIgnoresWeather || input.State.DefenderIgnoresWeather {
		effectiveWeather = WeatherNone
	}
	weather := input.WeatherBonus
	if effectiveWeather != WeatherNone {
		weather = WeatherDamageMultiplier(effectiveWeather, input.Move.Type)
	}
	if weather == 0 {
		weather = 1.0
	}

	// Unaware on the opposite side zeroes out the relevant stat stages.
	atkStages := input.AttackerStages
	defStages := input.DefenderStages
	if defAbility.Name == "unaware" {
		atkStages.Atk, atkStages.SpA = 0, 0
	}
	if atkAbility.Name == "unaware" {
		defStages.Def, defStages.SpD = 0, 0
	}

	var atkStat, defStat int
	// Apply stat-stage multipliers for Atk/Def/SpA/SpD. Crit hits ignore negative
	// attacker-offensive stages and positive defender-defensive stages (Gen 2+).
	if input.Move.Category == "special" {
		atkStage := atkStages.SpA
		defStage := defStages.SpD
		if input.IsCritical {
			if atkStage < 0 {
				atkStage = 0
			}
			if defStage > 0 {
				defStage = 0
			}
		}
		atkStat = int(math.Floor(float64(input.AttackerStats.SpAttack) * StageMultiplier(atkStage)))
		defStat = int(math.Floor(float64(input.DefenderStats.SpDefense) * StageMultiplier(defStage)))
		// Sandstorm boosts Rock-type SpDef ×1.5.
		if effectiveWeather == WeatherSandstorm && hasType(input.DefenderTypes, "rock") {
			defStat = int(math.Floor(float64(defStat) * 1.5))
		}
	} else {
		atkStage := atkStages.Atk
		defStage := defStages.Def
		if input.IsCritical {
			if atkStage < 0 {
				atkStage = 0
			}
			if defStage > 0 {
				defStage = 0
			}
		}
		atkStat = int(math.Floor(float64(input.AttackerStats.Attack) * StageMultiplier(atkStage)))
		defStat = int(math.Floor(float64(input.DefenderStats.Defense) * StageMultiplier(defStage)))
	}

	if atkStat <= 0 {
		atkStat = 1
	}
	if defStat <= 0 {
		defStat = 1
	}

	// Fórmula base: floor((floor(2*L/5+2) * Power * Atk/Def / 50) + 2)
	baseDmg := math.Floor(float64(2*input.Level)/5+2) *
		float64(input.Move.Power) *
		float64(atkStat) / float64(defStat) / 50
	baseDmg = math.Floor(baseDmg) + 2

	// STAB
	stab := 1.0
	hasSTAB := false
	for _, at := range input.AttackerTypes {
		if at.Name == input.Move.Type {
			stab = 1.5
			hasSTAB = true
			break
		}
	}

	// Efectividad de tipo
	typeEff := TypeEffectiveness(input.Move.Type, input.DefenderTypes)

	// Crítico
	crit := 1.0
	if input.IsCritical {
		crit = 1.5
	}

	// Quemadura: ×0.5 al daño físico si el atacante está quemado
	burn := 1.0
	burnApplied := false
	if input.IsBurned && input.Move.Category == "physical" {
		burn = 0.5
		burnApplied = true
	}

	// Ability damage multipliers. Attacker's ModifyAttack and defender's
	// ModifyDefense are applied as multiplicative modifiers. Adaptability's
	// STAB bump is baked into its ModifyAttack (2.0/1.5 when STAB applies).
	abilityMod := 1.0
	ctx := DamageCtx{
		State:             input.State,
		Move:              input.Move,
		AttackerTypes:     input.AttackerTypes,
		DefenderTypes:     input.DefenderTypes,
		TypeEffectiveness: typeEff,
	}
	if atkAbility.ModifyAttack != nil {
		ctx.Side = SideAttacker
		abilityMod *= atkAbility.ModifyAttack(ctx)
	}
	if defAbility.ModifyDefense != nil {
		ctx.Side = SideDefender
		abilityMod *= defAbility.ModifyDefense(ctx)
	}

	modifier := stab * typeEff * crit * weather * burn * abilityMod

	minDmg := int(math.Floor(baseDmg * modifier * 0.85))
	maxDmg := int(math.Floor(baseDmg * modifier * 1.00))

	if minDmg < 1 && maxDmg > 0 {
		minDmg = 1
	}

	return DamageResult{
		Min:                minDmg,
		Max:                maxDmg,
		Average:            (minDmg + maxDmg) / 2,
		ActualDamage:       (minDmg + maxDmg) / 2,
		Multiplier:         typeEff,
		IsSuperEffective:   typeEff > 1,
		IsNotVeryEffective: typeEff > 0 && typeEff < 1,
		HasNoEffect:        typeEff == 0 || abilityMod == 0,
		HasSTAB:            hasSTAB,
		STABMultiplier:     stab,
		BurnApplied:        burnApplied,
	}
}

// critThresholds maps critical stage to the denominator for probability.
// Stage 0 → 1/24, Stage 1 → 1/8, Stage 2 → 1/2, Stage ≥3 → 1/1 (always).
var critThresholds = [4]int{24, 8, 2, 1}

// CalculateBattleDamage computes damage with random roll (0.85–1.00) and
// probabilistic critical hits. randSource(n) must return [0, n).
func CalculateBattleDamage(input DamageInput, randSource func(n int) int) DamageResult {
	// Determine probabilistic critical hit
	stage := input.CriticalStage
	if stage < 0 {
		stage = 0
	}
	if stage > 3 {
		stage = 3
	}
	threshold := critThresholds[stage]
	wasCrit := input.IsCritical
	if !wasCrit && threshold > 0 {
		wasCrit = randSource(threshold) == 0
	}

	modInput := input
	modInput.IsCritical = wasCrit

	result := CalculateDamage(modInput)

	// Apply random roll: integer between 85 and 100 inclusive (16 values)
	roll := 85 + randSource(16)
	actual := int(math.Floor(float64(result.Max) * float64(roll) / 100.0))
	if actual < 1 && result.Max > 0 {
		actual = 1
	}

	result.ActualDamage = actual
	result.WasCritical = wasCrit
	return result
}
