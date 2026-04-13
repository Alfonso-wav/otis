package core

import "strings"

// Side identifies which combatant a hook applies to.
type Side string

const (
	SideAttacker Side = "attacker"
	SideDefender Side = "defender"
)

// DamageCtx is passed to ability damage-modifier hooks (ModifyAttack /
// ModifyDefense). Side is the *owner* of the ability being queried.
type DamageCtx struct {
	State             BattleState
	Side              Side
	Move              Move
	AttackerTypes     []PokemonType
	DefenderTypes     []PokemonType
	TypeEffectiveness float64
}

// AbilityEffect is the behavior of a Pokémon ability inside battle.
// All hooks are optional; zero-value means no-op.
// Hook functions must be pure: return new BattleState, never mutate input.
type AbilityEffect struct {
	Name      string
	DisplayEn string
	DisplayEs string

	// OnSwitchIn fires once when a Pokémon enters battle (start of battle
	// in a 1v1). May set weather, change stat stages, activate flags, etc.
	OnSwitchIn func(state BattleState, side Side) BattleState

	// OnTakeHit is called after raw damage is computed but before HP is
	// applied. Can redirect to healing (Water Absorb), activate a flag
	// (Flash Fire) or nullify with side effects.
	OnTakeHit func(state BattleState, side Side, move Move, incoming int) (BattleState, int)

	// ModifyAttack multiplies the attacker's offensive stat / damage.
	ModifyAttack func(ctx DamageCtx) float64

	// ModifyDefense multiplies the defender's defensive stat / damage.
	ModifyDefense func(ctx DamageCtx) float64

	// ImmuneToType returns true if the ability grants complete immunity
	// to a given move type.
	ImmuneToType func(moveType string) bool

	// SpeedMultiplier modifies Speed during turn-order resolution.
	SpeedMultiplier func(state BattleState, side Side) float64

	// EndOfTurn runs after the weather tick each turn.
	EndOfTurn func(state BattleState, side Side) BattleState
}

// GetAbility returns the catalog entry for an ability name. Accepts PokéAPI
// slugs and names case-insensitively. Returns (_, false) for unknown / empty.
func GetAbility(name string) (AbilityEffect, bool) {
	key := normalizeAbilityName(name)
	if key == "" {
		return AbilityEffect{}, false
	}
	ab, ok := abilitiesCatalog[key]
	return ab, ok
}

// AbilityNames returns all registered ability slugs. Stable across calls.
func AbilityNames() []string {
	out := make([]string, 0, len(abilitiesCatalog))
	for k := range abilitiesCatalog {
		out = append(out, k)
	}
	return out
}

func normalizeAbilityName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	return strings.ReplaceAll(s, " ", "-")
}

// --- Helpers used by hook implementations ---

func sideHP(state BattleState, side Side) (hp, maxHP int) {
	if side == SideAttacker {
		return state.AttackerHP, state.AttackerMaxHP
	}
	return state.DefenderHP, state.DefenderMaxHP
}

func setSideHP(state BattleState, side Side, hp int) BattleState {
	if hp < 0 {
		hp = 0
	}
	if side == SideAttacker {
		if hp > state.AttackerMaxHP {
			hp = state.AttackerMaxHP
		}
		state.AttackerHP = hp
	} else {
		if hp > state.DefenderMaxHP {
			hp = state.DefenderMaxHP
		}
		state.DefenderHP = hp
	}
	return state
}

func opposite(side Side) Side {
	if side == SideAttacker {
		return SideDefender
	}
	return SideAttacker
}

func applyStageToSide(state BattleState, side Side, stat StatKey, delta int) BattleState {
	if side == SideAttacker {
		state.AttackerStages = ApplyStage(state.AttackerStages, stat, delta)
	} else {
		state.DefenderStages = ApplyStage(state.DefenderStages, stat, delta)
	}
	return state
}

// hpRatio returns currentHP/maxHP for a given side. Safe when max<=0.
func hpRatio(state BattleState, side Side) float64 {
	hp, max := sideHP(state, side)
	if max <= 0 {
		return 0
	}
	return float64(hp) / float64(max)
}

// pinchBoostMultiplier is the classic Blaze/Torrent/Overgrow/Swarm helper:
// if moveType matches matchType AND HP ratio < 1/3 → ×1.5.
func pinchBoostMultiplier(ctx DamageCtx, matchType string) float64 {
	if ctx.Move.Type != matchType {
		return 1.0
	}
	if hpRatio(ctx.State, ctx.Side) < 1.0/3.0 {
		return 1.5
	}
	return 1.0
}

// isContactMove is a heuristic: physical moves are contact, special/status are not.
// Real contact data is in PokéAPI move meta; approximation sufficient for catalog.
func isContactMove(m Move) bool {
	return m.Category == "physical" && m.Power > 0
}

// moveNameContainsAny reports whether a move name contains any of the given
// substrings (case-insensitive). Used for punch / bite / pulse detection.
func moveNameContainsAny(m Move, words ...string) bool {
	name := strings.ToLower(m.Name)
	for _, w := range words {
		if strings.Contains(name, w) {
			return true
		}
	}
	return false
}

// weatherOrNone returns the effective weather, respecting cloud-nine/air-lock
// flags on either side.
func weatherOrNone(state BattleState) Weather {
	if state.AttackerIgnoresWeather || state.DefenderIgnoresWeather {
		return WeatherNone
	}
	return state.Weather
}

// --- Catalog build ---

var abilitiesCatalog = buildAbilitiesCatalog()

func buildAbilitiesCatalog() map[string]AbilityEffect {
	c := map[string]AbilityEffect{}
	register := func(effects ...AbilityEffect) {
		for _, e := range effects {
			c[e.Name] = e
		}
	}

	// ---------- Weather abilities (OnSwitchIn sets permanent weather) ----------
	register(
		AbilityEffect{
			Name: "drizzle", DisplayEn: "Drizzle", DisplayEs: "Llovizna",
			OnSwitchIn: func(state BattleState, side Side) BattleState {
				state.Weather = WeatherRain
				state.WeatherTurnsLeft = WeatherDefaultTurns
				return state
			},
		},
		AbilityEffect{
			Name: "drought", DisplayEn: "Drought", DisplayEs: "Sequía",
			OnSwitchIn: func(state BattleState, side Side) BattleState {
				state.Weather = WeatherSun
				state.WeatherTurnsLeft = WeatherDefaultTurns
				return state
			},
		},
		AbilityEffect{
			Name: "sand-stream", DisplayEn: "Sand Stream", DisplayEs: "Chorro Arena",
			OnSwitchIn: func(state BattleState, side Side) BattleState {
				state.Weather = WeatherSandstorm
				state.WeatherTurnsLeft = WeatherDefaultTurns
				return state
			},
		},
		AbilityEffect{
			Name: "snow-warning", DisplayEn: "Snow Warning", DisplayEs: "Nevada",
			OnSwitchIn: func(state BattleState, side Side) BattleState {
				state.Weather = WeatherHail
				state.WeatherTurnsLeft = WeatherDefaultTurns
				return state
			},
		},
		AbilityEffect{
			Name: "cloud-nine", DisplayEn: "Cloud Nine", DisplayEs: "Antinubes",
			OnSwitchIn: func(state BattleState, side Side) BattleState {
				if side == SideAttacker {
					state.AttackerIgnoresWeather = true
				} else {
					state.DefenderIgnoresWeather = true
				}
				return state
			},
		},
		AbilityEffect{
			Name: "air-lock", DisplayEn: "Air Lock", DisplayEs: "Esclusa de Aire",
			OnSwitchIn: func(state BattleState, side Side) BattleState {
				if side == SideAttacker {
					state.AttackerIgnoresWeather = true
				} else {
					state.DefenderIgnoresWeather = true
				}
				return state
			},
		},
	)

	// ---------- Type immunities / absorbs ----------
	register(
		AbilityEffect{
			Name: "levitate", DisplayEn: "Levitate", DisplayEs: "Levitación",
			ImmuneToType: func(mt string) bool { return mt == "ground" },
		},
		AbilityEffect{
			Name: "water-absorb", DisplayEn: "Water Absorb", DisplayEs: "Absorbe Agua",
			OnTakeHit: func(state BattleState, side Side, move Move, incoming int) (BattleState, int) {
				if move.Type != "water" {
					return state, incoming
				}
				_, max := sideHP(state, side)
				state = setSideHP(state, side, func() int {
					hp, _ := sideHP(state, side)
					return hp + max/4
				}())
				return state, 0
			},
			ImmuneToType: func(mt string) bool { return mt == "water" },
		},
		AbilityEffect{
			Name: "volt-absorb", DisplayEn: "Volt Absorb", DisplayEs: "Absorbe Electricidad",
			OnTakeHit: func(state BattleState, side Side, move Move, incoming int) (BattleState, int) {
				if move.Type != "electric" {
					return state, incoming
				}
				_, max := sideHP(state, side)
				state = setSideHP(state, side, func() int {
					hp, _ := sideHP(state, side)
					return hp + max/4
				}())
				return state, 0
			},
			ImmuneToType: func(mt string) bool { return mt == "electric" },
		},
		AbilityEffect{
			Name: "flash-fire", DisplayEn: "Flash Fire", DisplayEs: "Absorbe Fuego",
			OnTakeHit: func(state BattleState, side Side, move Move, incoming int) (BattleState, int) {
				if move.Type != "fire" {
					return state, incoming
				}
				if side == SideAttacker {
					state.AttackerFlashFireActive = true
				} else {
					state.DefenderFlashFireActive = true
				}
				return state, 0
			},
			ImmuneToType: func(mt string) bool { return mt == "fire" },
			ModifyAttack: func(ctx DamageCtx) float64 {
				if ctx.Move.Type != "fire" {
					return 1.0
				}
				if ctx.Side == SideAttacker && ctx.State.AttackerFlashFireActive {
					return 1.5
				}
				if ctx.Side == SideDefender && ctx.State.DefenderFlashFireActive {
					return 1.5
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "lightning-rod", DisplayEn: "Lightning Rod", DisplayEs: "Pararrayos",
			ImmuneToType: func(mt string) bool { return mt == "electric" },
		},
		AbilityEffect{
			Name: "storm-drain", DisplayEn: "Storm Drain", DisplayEs: "Colector",
			ImmuneToType: func(mt string) bool { return mt == "water" },
		},
		AbilityEffect{
			Name: "motor-drive", DisplayEn: "Motor Drive", DisplayEs: "Electromotor",
			OnTakeHit: func(state BattleState, side Side, move Move, incoming int) (BattleState, int) {
				if move.Type != "electric" {
					return state, incoming
				}
				state = applyStageToSide(state, side, StatSpe, 1)
				return state, 0
			},
			ImmuneToType: func(mt string) bool { return mt == "electric" },
		},
		AbilityEffect{
			Name: "sap-sipper", DisplayEn: "Sap Sipper", DisplayEs: "Herbívoro",
			OnTakeHit: func(state BattleState, side Side, move Move, incoming int) (BattleState, int) {
				if move.Type != "grass" {
					return state, incoming
				}
				state = applyStageToSide(state, side, StatAtk, 1)
				return state, 0
			},
			ImmuneToType: func(mt string) bool { return mt == "grass" },
		},
	)

	// ---------- Stat-stage on switch-in ----------
	register(
		AbilityEffect{
			Name: "intimidate", DisplayEn: "Intimidate", DisplayEs: "Intimidación",
			OnSwitchIn: func(state BattleState, side Side) BattleState {
				return applyStageToSide(state, opposite(side), StatAtk, -1)
			},
		},
		AbilityEffect{
			Name: "download", DisplayEn: "Download", DisplayEs: "Descarga",
			// Without access to live opponent stats at switch-in, approximate
			// by boosting SpA if defender's SpD stage <= Def stage, else Atk.
			// Real Download compares base stats, not stages — this is best-effort.
			OnSwitchIn: func(state BattleState, side Side) BattleState {
				var oppStages StatStages
				if side == SideAttacker {
					oppStages = state.DefenderStages
				} else {
					oppStages = state.AttackerStages
				}
				if oppStages.SpD <= oppStages.Def {
					return applyStageToSide(state, side, StatSpA, 1)
				}
				return applyStageToSide(state, side, StatAtk, 1)
			},
		},
		AbilityEffect{
			Name: "competitive", DisplayEn: "Competitive", DisplayEs: "Tenacidad",
			// Hook on stat drops is not implemented; placeholder registered.
		},
		AbilityEffect{
			Name: "defiant", DisplayEn: "Defiant", DisplayEs: "Competitivo",
			// Placeholder — reaction-on-drop not hooked.
		},
	)

	// ---------- Conditional type boost (pinch / weather / flag) ----------
	register(
		AbilityEffect{
			Name: "blaze", DisplayEn: "Blaze", DisplayEs: "Mar Llamas",
			ModifyAttack: func(ctx DamageCtx) float64 { return pinchBoostMultiplier(ctx, "fire") },
		},
		AbilityEffect{
			Name: "torrent", DisplayEn: "Torrent", DisplayEs: "Torrente",
			ModifyAttack: func(ctx DamageCtx) float64 { return pinchBoostMultiplier(ctx, "water") },
		},
		AbilityEffect{
			Name: "overgrow", DisplayEn: "Overgrow", DisplayEs: "Espesura",
			ModifyAttack: func(ctx DamageCtx) float64 { return pinchBoostMultiplier(ctx, "grass") },
		},
		AbilityEffect{
			Name: "swarm", DisplayEn: "Swarm", DisplayEs: "Enjambre",
			ModifyAttack: func(ctx DamageCtx) float64 { return pinchBoostMultiplier(ctx, "bug") },
		},
		AbilityEffect{
			Name: "solar-power", DisplayEn: "Solar Power", DisplayEs: "Poder Solar",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if weatherOrNone(ctx.State) == WeatherSun && ctx.Move.Category == "special" {
					return 1.5
				}
				return 1.0
			},
			EndOfTurn: func(state BattleState, side Side) BattleState {
				if weatherOrNone(state) != WeatherSun {
					return state
				}
				_, max := sideHP(state, side)
				hp, _ := sideHP(state, side)
				return setSideHP(state, side, hp-max/8)
			},
		},
		AbilityEffect{
			Name: "chlorophyll", DisplayEn: "Chlorophyll", DisplayEs: "Clorofila",
			SpeedMultiplier: func(state BattleState, side Side) float64 {
				if weatherOrNone(state) == WeatherSun {
					return 2.0
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "swift-swim", DisplayEn: "Swift Swim", DisplayEs: "Nado Rápido",
			SpeedMultiplier: func(state BattleState, side Side) float64 {
				if weatherOrNone(state) == WeatherRain {
					return 2.0
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "sand-rush", DisplayEn: "Sand Rush", DisplayEs: "Ímpetu Arena",
			SpeedMultiplier: func(state BattleState, side Side) float64 {
				if weatherOrNone(state) == WeatherSandstorm {
					return 2.0
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "slush-rush", DisplayEn: "Slush Rush", DisplayEs: "Quitanieves",
			SpeedMultiplier: func(state BattleState, side Side) float64 {
				if weatherOrNone(state) == WeatherHail {
					return 2.0
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "sand-force", DisplayEn: "Sand Force", DisplayEs: "Poder Arena",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if weatherOrNone(ctx.State) != WeatherSandstorm {
					return 1.0
				}
				switch ctx.Move.Type {
				case "rock", "ground", "steel":
					return 1.3
				}
				return 1.0
			},
		},
	)

	// ---------- Contact / attack side-effects (placeholders without status) ----------
	register(
		AbilityEffect{Name: "static", DisplayEn: "Static", DisplayEs: "Electricidad Estática"},
		AbilityEffect{Name: "flame-body", DisplayEn: "Flame Body", DisplayEs: "Cuerpo Llama"},
		AbilityEffect{Name: "poison-point", DisplayEn: "Poison Point", DisplayEs: "Punto Tóxico"},
		AbilityEffect{
			Name: "rough-skin", DisplayEn: "Rough Skin", DisplayEs: "Piel Tosca",
			OnTakeHit: func(state BattleState, side Side, move Move, incoming int) (BattleState, int) {
				if !isContactMove(move) || incoming <= 0 {
					return state, incoming
				}
				// Attacker (opposite of this side) loses 1/8 max HP.
				att := opposite(side)
				_, max := sideHP(state, att)
				hp, _ := sideHP(state, att)
				state = setSideHP(state, att, hp-max/8)
				return state, incoming
			},
		},
		AbilityEffect{
			Name: "iron-barbs", DisplayEn: "Iron Barbs", DisplayEs: "Punta Acero",
			OnTakeHit: func(state BattleState, side Side, move Move, incoming int) (BattleState, int) {
				if !isContactMove(move) || incoming <= 0 {
					return state, incoming
				}
				att := opposite(side)
				_, max := sideHP(state, att)
				hp, _ := sideHP(state, att)
				state = setSideHP(state, att, hp-max/8)
				return state, incoming
			},
		},
		AbilityEffect{Name: "cursed-body", DisplayEn: "Cursed Body", DisplayEs: "Cuerpo Maldito"},
	)

	// ---------- Damage reduction ----------
	register(
		AbilityEffect{
			Name: "thick-fat", DisplayEn: "Thick Fat", DisplayEs: "Sebo",
			ModifyDefense: func(ctx DamageCtx) float64 {
				if ctx.Move.Type == "fire" || ctx.Move.Type == "ice" {
					return 0.5
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "heatproof", DisplayEn: "Heatproof", DisplayEs: "Ignífugo",
			ModifyDefense: func(ctx DamageCtx) float64 {
				if ctx.Move.Type == "fire" {
					return 0.5
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "filter", DisplayEn: "Filter", DisplayEs: "Filtro",
			ModifyDefense: func(ctx DamageCtx) float64 {
				if ctx.TypeEffectiveness > 1.0 {
					return 0.75
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "solid-rock", DisplayEn: "Solid Rock", DisplayEs: "Roca Sólida",
			ModifyDefense: func(ctx DamageCtx) float64 {
				if ctx.TypeEffectiveness > 1.0 {
					return 0.75
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "prism-armor", DisplayEn: "Prism Armor", DisplayEs: "Armadura Prisma",
			ModifyDefense: func(ctx DamageCtx) float64 {
				if ctx.TypeEffectiveness > 1.0 {
					return 0.75
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "multiscale", DisplayEn: "Multiscale", DisplayEs: "Multiescamas",
			ModifyDefense: func(ctx DamageCtx) float64 {
				if hpRatio(ctx.State, ctx.Side) >= 1.0 {
					return 0.5
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "shadow-shield", DisplayEn: "Shadow Shield", DisplayEs: "Escudo Fantasma",
			ModifyDefense: func(ctx DamageCtx) float64 {
				if hpRatio(ctx.State, ctx.Side) >= 1.0 {
					return 0.5
				}
				return 1.0
			},
		},
	)

	// ---------- Miscellaneous offensive ----------
	register(
		AbilityEffect{
			Name: "huge-power", DisplayEn: "Huge Power", DisplayEs: "Potencia",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if ctx.Move.Category == "physical" {
					return 2.0
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "pure-power", DisplayEn: "Pure Power", DisplayEs: "Energía Pura",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if ctx.Move.Category == "physical" {
					return 2.0
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "sheer-force", DisplayEn: "Sheer Force", DisplayEs: "Potencia Bruta",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if ctx.Move.Power > 0 {
					return 1.3
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "tough-claws", DisplayEn: "Tough Claws", DisplayEs: "Garra Dura",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if isContactMove(ctx.Move) {
					return 1.3
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "adaptability", DisplayEn: "Adaptability", DisplayEs: "Adaptable",
			// STAB from 1.5 to 2.0 = add x (2.0/1.5) when STAB applies.
			ModifyAttack: func(ctx DamageCtx) float64 {
				for _, t := range ctx.AttackerTypes {
					if t.Name == ctx.Move.Type {
						return 2.0 / 1.5
					}
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "technician", DisplayEn: "Technician", DisplayEs: "Experto",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if ctx.Move.Power > 0 && ctx.Move.Power <= 60 {
					return 1.5
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "reckless", DisplayEn: "Reckless", DisplayEs: "Audaz",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if moveNameContainsAny(ctx.Move, "take-down", "double-edge", "submission", "head-charge", "head-smash", "wood-hammer", "brave-bird", "flare-blitz", "volt-tackle", "wild-charge") {
					return 1.2
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "iron-fist", DisplayEn: "Iron Fist", DisplayEs: "Puño Férreo",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if moveNameContainsAny(ctx.Move, "punch") {
					return 1.2
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "mega-launcher", DisplayEn: "Mega Launcher", DisplayEs: "Megadisparador",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if moveNameContainsAny(ctx.Move, "pulse", "aura-sphere", "dark-pulse", "dragon-pulse", "water-pulse", "heal-pulse") {
					return 1.5
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "strong-jaw", DisplayEn: "Strong Jaw", DisplayEs: "Mandíbula Fuerte",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if moveNameContainsAny(ctx.Move, "bite", "crunch", "fang", "chomp", "jaw") {
					return 1.5
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "tinted-lens", DisplayEn: "Tinted Lens", DisplayEs: "Cristal Tintado",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if ctx.TypeEffectiveness > 0 && ctx.TypeEffectiveness < 1 {
					return 2.0
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "neuroforce", DisplayEn: "Neuroforce", DisplayEs: "Inerciafuerza",
			ModifyAttack: func(ctx DamageCtx) float64 {
				if ctx.TypeEffectiveness > 1.0 {
					return 1.25
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "analytic", DisplayEn: "Analytic", DisplayEs: "Cálculo Final",
			// Approximation: always +1.3 (we don't track who went last).
			ModifyAttack: func(ctx DamageCtx) float64 { return 1.3 },
		},
	)

	// ---------- Weather-passive healing + variants ----------
	register(
		AbilityEffect{
			Name: "rain-dish", DisplayEn: "Rain Dish", DisplayEs: "Cura Lluvia",
			EndOfTurn: func(state BattleState, side Side) BattleState {
				if weatherOrNone(state) != WeatherRain {
					return state
				}
				hp, max := sideHP(state, side)
				return setSideHP(state, side, hp+max/16)
			},
		},
		AbilityEffect{
			Name: "ice-body", DisplayEn: "Ice Body", DisplayEs: "Gélido",
			EndOfTurn: func(state BattleState, side Side) BattleState {
				if weatherOrNone(state) != WeatherHail {
					return state
				}
				hp, max := sideHP(state, side)
				return setSideHP(state, side, hp+max/16)
			},
		},
		AbilityEffect{
			Name: "dry-skin", DisplayEn: "Dry Skin", DisplayEs: "Piel Seca",
			OnTakeHit: func(state BattleState, side Side, move Move, incoming int) (BattleState, int) {
				if move.Type != "water" {
					return state, incoming
				}
				_, max := sideHP(state, side)
				hp, _ := sideHP(state, side)
				return setSideHP(state, side, hp+max/4), 0
			},
			ImmuneToType: func(mt string) bool { return mt == "water" },
			ModifyDefense: func(ctx DamageCtx) float64 {
				if ctx.Move.Type == "fire" {
					return 1.25
				}
				return 1.0
			},
			EndOfTurn: func(state BattleState, side Side) BattleState {
				switch weatherOrNone(state) {
				case WeatherRain:
					hp, max := sideHP(state, side)
					return setSideHP(state, side, hp+max/8)
				case WeatherSun:
					hp, max := sideHP(state, side)
					return setSideHP(state, side, hp-max/8)
				}
				return state
			},
		},
	)

	// ---------- Evasion weather ----------
	// (Accuracy is not yet piped through ability hooks — registered as metadata.)
	register(
		AbilityEffect{Name: "sand-veil", DisplayEn: "Sand Veil", DisplayEs: "Velo Arena"},
		AbilityEffect{Name: "snow-cloak", DisplayEn: "Snow Cloak", DisplayEs: "Manto Níveo"},
	)

	// ---------- Immunities / utility ----------
	register(
		AbilityEffect{
			Name: "wonder-guard", DisplayEn: "Wonder Guard", DisplayEs: "Superguarda",
			ModifyDefense: func(ctx DamageCtx) float64 {
				// Only super-effective hits deal damage. Below 1.0 → zero damage.
				if ctx.Move.Power <= 0 {
					return 1.0
				}
				if ctx.TypeEffectiveness <= 1.0 {
					return 0.0
				}
				return 1.0
			},
		},
		AbilityEffect{
			Name: "unaware", DisplayEn: "Unaware", DisplayEs: "Ignorante",
			// Ignores opponent stat stages. Handled specially in CalculateDamage by
			// checking the *defender's* ability for unaware before applying attacker
			// stages. ModifyDefense cannot zero-out stage effects easily; marker only.
		},
		AbilityEffect{
			Name: "speed-boost", DisplayEn: "Speed Boost", DisplayEs: "Impulso",
			EndOfTurn: func(state BattleState, side Side) BattleState {
				return applyStageToSide(state, side, StatSpe, 1)
			},
		},
		AbilityEffect{Name: "contrary", DisplayEn: "Contrary", DisplayEs: "Respondón"},
		AbilityEffect{Name: "simple", DisplayEn: "Simple", DisplayEs: "Simple"},
		AbilityEffect{Name: "clear-body", DisplayEn: "Clear Body", DisplayEs: "Cuerpo Puro"},
		AbilityEffect{Name: "white-smoke", DisplayEn: "White Smoke", DisplayEs: "Humo Blanco"},
		AbilityEffect{Name: "hyper-cutter", DisplayEn: "Hyper Cutter", DisplayEs: "Corte Fuerte"},
		AbilityEffect{Name: "big-pecks", DisplayEn: "Big Pecks", DisplayEs: "Sacapecho"},
		AbilityEffect{Name: "keen-eye", DisplayEn: "Keen Eye", DisplayEs: "Vista Lince"},
		AbilityEffect{Name: "compound-eyes", DisplayEn: "Compound Eyes", DisplayEs: "Ojo Compuesto"},
		AbilityEffect{Name: "no-guard", DisplayEn: "No Guard", DisplayEs: "Indefenso"},
		AbilityEffect{Name: "sniper", DisplayEn: "Sniper", DisplayEs: "Francotirador"},
		AbilityEffect{Name: "serene-grace", DisplayEn: "Serene Grace", DisplayEs: "Dicha"},
		AbilityEffect{Name: "super-luck", DisplayEn: "Super Luck", DisplayEs: "Afortunado"},
		AbilityEffect{Name: "moxie", DisplayEn: "Moxie", DisplayEs: "Autoestima"},
		AbilityEffect{Name: "aftermath", DisplayEn: "Aftermath", DisplayEs: "Detonación"},
		AbilityEffect{Name: "guts", DisplayEn: "Guts", DisplayEs: "Agallas"},
		AbilityEffect{Name: "marvel-scale", DisplayEn: "Marvel Scale", DisplayEs: "Escama Especial"},
		AbilityEffect{Name: "quick-feet", DisplayEn: "Quick Feet", DisplayEs: "Pies Rápidos"},
		AbilityEffect{Name: "toxic-boost", DisplayEn: "Toxic Boost", DisplayEs: "Ímpetu Tóxico"},
		AbilityEffect{Name: "flare-boost", DisplayEn: "Flare Boost", DisplayEs: "Ímpetu Ardiente"},
		AbilityEffect{Name: "regenerator", DisplayEn: "Regenerator", DisplayEs: "Regeneración"},
		AbilityEffect{Name: "natural-cure", DisplayEn: "Natural Cure", DisplayEs: "Cura Natural"},
		AbilityEffect{Name: "truant", DisplayEn: "Truant", DisplayEs: "Holgazán"},
	)

	return c
}
