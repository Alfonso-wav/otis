package core

import "strings"

// StatStageBounds defines the canonical -6..+6 range for stat stages.
const (
	StageMin = -6
	StageMax = 6
)

// Stat identifiers used by the stage system.
type StatKey string

const (
	StatAtk StatKey = "atk"
	StatDef StatKey = "def"
	StatSpA StatKey = "spa"
	StatSpD StatKey = "spd"
	StatSpe StatKey = "spe"
	StatAcc StatKey = "acc"
	StatEva StatKey = "eva"
)

// clampStage bounds a stage value into [StageMin, StageMax].
func clampStage(n int) int {
	if n < StageMin {
		return StageMin
	}
	if n > StageMax {
		return StageMax
	}
	return n
}

// ApplyStage returns a new StatStages with the given stat modified by delta,
// clamped to -6..+6. Pure function: the input is not mutated.
func ApplyStage(stages StatStages, stat StatKey, delta int) StatStages {
	out := stages
	switch stat {
	case StatAtk:
		out.Atk = clampStage(stages.Atk + delta)
	case StatDef:
		out.Def = clampStage(stages.Def + delta)
	case StatSpA:
		out.SpA = clampStage(stages.SpA + delta)
	case StatSpD:
		out.SpD = clampStage(stages.SpD + delta)
	case StatSpe:
		out.Spe = clampStage(stages.Spe + delta)
	case StatAcc:
		out.Acc = clampStage(stages.Acc + delta)
	case StatEva:
		out.Eva = clampStage(stages.Eva + delta)
	}
	return out
}

// StageMultiplier returns the canonical multiplier for Atk/Def/SpA/SpD/Spe.
// +n → (2+n)/2, -n → 2/(2+|n|). Clamped outside -6..+6.
func StageMultiplier(stage int) float64 {
	n := clampStage(stage)
	num, den := 2, 2
	if n >= 0 {
		num = 2 + n
	} else {
		den = 2 + (-n)
	}
	return float64(num) / float64(den)
}

// AccuracyStageMultiplier returns the canonical multiplier for Accuracy / Evasion.
// Base 3 instead of 2: +n → (3+n)/3, -n → 3/(3+|n|).
func AccuracyStageMultiplier(stage int) float64 {
	n := clampStage(stage)
	num, den := 3, 3
	if n >= 0 {
		num = 3 + n
	} else {
		den = 3 + (-n)
	}
	return float64(num) / float64(den)
}

// GetStage returns the current stage value for a given stat.
func GetStage(stages StatStages, stat StatKey) int {
	switch stat {
	case StatAtk:
		return stages.Atk
	case StatDef:
		return stages.Def
	case StatSpA:
		return stages.SpA
	case StatSpD:
		return stages.SpD
	case StatSpe:
		return stages.Spe
	case StatAcc:
		return stages.Acc
	case StatEva:
		return stages.Eva
	}
	return 0
}

// StatStageChange describes a single stat stage delta triggered by a move.
// Target = "self" applies to the attacker, "opponent" applies to the defender.
type StatStageChange struct {
	Target string // "self" | "opponent"
	Stat   StatKey
	Delta  int
}

// statStageMoveEffects maps PokéAPI move slugs to the stage changes they apply.
// Scope: canonical single-effect boosts / debuffs (see task 0167).
var statStageMoveEffects = map[string][]StatStageChange{
	// --- Self-boost ---
	"swords-dance":  {{Target: "self", Stat: StatAtk, Delta: 2}},
	"agility":       {{Target: "self", Stat: StatSpe, Delta: 2}},
	"nasty-plot":    {{Target: "self", Stat: StatSpA, Delta: 2}},
	"calm-mind":     {{Target: "self", Stat: StatSpA, Delta: 1}, {Target: "self", Stat: StatSpD, Delta: 1}},
	"bulk-up":       {{Target: "self", Stat: StatAtk, Delta: 1}, {Target: "self", Stat: StatDef, Delta: 1}},
	"iron-defense":  {{Target: "self", Stat: StatDef, Delta: 2}},
	"amnesia":       {{Target: "self", Stat: StatSpD, Delta: 2}},
	"barrier":       {{Target: "self", Stat: StatDef, Delta: 2}},
	// --- Opponent debuff ---
	"growl":        {{Target: "opponent", Stat: StatAtk, Delta: -1}},
	"leer":         {{Target: "opponent", Stat: StatDef, Delta: -1}},
	"tail-whip":    {{Target: "opponent", Stat: StatDef, Delta: -1}},
	"string-shot":  {{Target: "opponent", Stat: StatSpe, Delta: -1}},
	"smokescreen":  {{Target: "opponent", Stat: StatAcc, Delta: -1}},
	"sand-attack":  {{Target: "opponent", Stat: StatAcc, Delta: -1}},
}

// StatStageEffects returns the list of stage changes (if any) that the given
// move produces. Case-insensitive on move name. Accepts PokéAPI slugs.
func StatStageEffects(moveName string) ([]StatStageChange, bool) {
	key := strings.ToLower(strings.TrimSpace(moveName))
	key = strings.ReplaceAll(key, " ", "-")
	effects, ok := statStageMoveEffects[key]
	return effects, ok
}

// StageChangeLogFragment returns a short Spanish fragment describing the outcome
// of applying a delta to a stat that already had the given previous stage. The
// verb reflects magnitude ("subió" / "subió mucho" / "bajó" / "bajó mucho") and
// detects hard caps ("no puede subir más").
func StageChangeLogFragment(prev, delta int) string {
	if delta > 0 {
		if prev >= StageMax {
			return "no puede subir más"
		}
		if delta >= 2 {
			return "subió mucho"
		}
		return "subió"
	}
	if delta < 0 {
		if prev <= StageMin {
			return "no puede bajar más"
		}
		if delta <= -2 {
			return "bajó mucho"
		}
		return "bajó"
	}
	return ""
}

// StatLabelEs returns the Spanish label of a stat for log messages.
func StatLabelEs(stat StatKey) string {
	switch stat {
	case StatAtk:
		return "Ataque"
	case StatDef:
		return "Defensa"
	case StatSpA:
		return "Ataque Especial"
	case StatSpD:
		return "Defensa Especial"
	case StatSpe:
		return "Velocidad"
	case StatAcc:
		return "Precisión"
	case StatEva:
		return "Evasión"
	}
	return string(stat)
}
