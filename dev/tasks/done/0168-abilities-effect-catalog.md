# Abilities: catálogo de efectos y hooks en motor de batalla

**ID**: 0168-abilities-effect-catalog
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

Modelar habilidades Pokémon como efectos tipados con hooks en el ciclo de batalla. Catálogo inicial ~80-100 habilidades combate-relevantes. Integrar en `ExecuteTurn`, `CalculateDamage`, `SimulateFullBattle`.

Depende de 0166 (clima) y 0167 (stat stages).

## Capas afectadas

- **Core**: nuevo paquete lógico `core/abilities.go`, integración en `battle.go` y `damage.go`. Extensión de `BattleState` con campos `AttackerAbility *Ability`, `DefenderAbility *Ability`.
- **Shell**: ninguno. Los datos de habilidades ya existen (`core.Ability` domain + PokéAPI).
- **APP**: sin cambios — la habilidad llega como campo en `FullBattleInput`.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/abilities.go` | crear | `AbilityEffect` struct con hooks opcionales: `OnSwitchIn(state) state`, `OnDamageReceived(state, dmg, move) (state, newDmg)`, `OnDamageDealt(state, dmg, move) (state, newDmg)`, `PassiveAtkModifier(ctx) float64`, `PassiveImmunity(moveType) bool`. Catálogo `AbilitiesCatalog map[string]AbilityEffect`. |
| `core/abilities_test.go` | crear | Tests unitarios por habilidad (pura, input/output). |
| `core/domain.go` | modificar | Añadir `AttackerAbility string`, `DefenderAbility string` a `FullBattleInput`. Extender `BattleState` con stat stages / clima ya cubiertos por 0166/0167. |
| `core/battle.go` | modificar | Al inicio de `SimulateFullBattle` disparar `OnSwitchIn` de ambas habilidades (aplica clima, stages, etc.). En `executeAttackerTurn`/`executeDefenderTurn` pasar `Ability` a `CalculateDamage` y aplicar hooks `OnDamageReceived`/`OnDamageDealt`. |
| `core/damage.go` | modificar | `DamageInput` acepta `AttackerAbility`, `DefenderAbility`. Aplicar modificadores pasivos (STAB boost, type immunity por habilidad, weather interactions). |
| `core/battle_test.go` | modificar | Tests de integración: Llovizna setea Rain al entrar. Intimidación baja Atk rival. Absorbe Agua inmune a Water + cura 25%. Espesura x1.5 a Grass con HP<1/3. |

## Plan de implementacion

Diseño del `AbilityEffect`:

```go
type AbilityEffect struct {
  Name        string // kebab-case, coincide con PokéAPI (ej. "drizzle")
  DisplayEn   string
  DisplayEs   string
  OnSwitchIn  func(state BattleState, side Side) BattleState
  OnTakeHit   func(state BattleState, side Side, move Move, incoming int) (BattleState, int) // puede anular, curar, reducir
  ModifyAttack func(ctx DamageCtx) float64 // multiplicador, 1.0 = no-op
  ModifyDefense func(ctx DamageCtx) float64
  ImmuneToType func(moveType string) bool
  EndOfTurn   func(state BattleState, side Side) BattleState
}
```

Catálogo inicial (~80-100) agrupado por categoría:

### Clima (7)
- `drizzle` → `OnSwitchIn`: setea Rain.
- `drought` → Sun.
- `sand-stream` → Sandstorm.
- `snow-warning` → Hail.
- `cloud-nine` / `air-lock` → flag `IgnoreWeather` (desactiva modificadores de clima mientras esté en campo).

### Inmunidades por tipo (8)
- `levitate` → inmune Ground.
- `water-absorb` → inmune Water + curar 25% HP.
- `volt-absorb` → inmune Electric + curar 25%.
- `flash-fire` → inmune Fire + boost Fire x1.5 (flag activado).
- `lightning-rod` → redirige Electric a sí mismo + inmune.
- `storm-drain` → ídem Water.
- `motor-drive` → inmune Electric + Spe +1.
- `sap-sipper` → inmune Grass + Atk +1.

### Stat-stage on switch-in (4)
- `intimidate` → Atk rival -1 al entrar.
- `download` → compara Def/SpD rival; sube Atk o SpA +1.
- `competitive` → al recibir bajada de stat → SpA +2.
- `defiant` → ídem, Atk +2.

### STAB/tipo boost condicional (10)
- `blaze` → Fire x1.5 si HP<1/3.
- `torrent` → Water x1.5 si HP<1/3.
- `overgrow` → Grass x1.5 si HP<1/3.
- `swarm` → Bug x1.5 si HP<1/3.
- `flash-fire` → Fire x1.5 si activado.
- `solar-power` → SpA x1.5 en Sun (pierde 1/8 HP por turno).
- `chlorophyll` → Spe x2 en Sun.
- `swift-swim` → Spe x2 en Rain.
- `sand-rush` → Spe x2 en Sandstorm.
- `slush-rush` → Spe x2 en Hail.

### Daño por contacto / al atacar (6)
- `static` → 30% paralizar al rival por contacto (si implementamos status paralizar, sino flag Spe -1 persistente).
- `flame-body` → 30% quemar (si hay status burn → Atk x0.5; sino log).
- `poison-point` → 30% envenenar.
- `rough-skin` / `iron-barbs` → rival pierde 1/8 HP por contacto.
- `cursed-body` → 30% disable move (skip si no hay PP tracking).

### Reducción de daño (5)
- `thick-fat` → daño recibido Fire/Ice x0.5.
- `heatproof` → daño Fire x0.5.
- `filter` / `solid-rock` → super efectivo x0.75.
- `multiscale` → a HP completo → daño recibido x0.5.

### Misceláneas (10)
- `guts` → Atk x1.5 con status (sin sistema status → skip o placeholder).
- `huge-power` / `pure-power` → Atk x2.
- `sheer-force` → moves con efecto → x1.3, sin efecto secundario (placeholder: x1.3 a todos los que Power>0).
- `tough-claws` → contacto x1.3.
- `adaptability` → STAB x2 (en vez de x1.5).
- `technician` → moves Power<=60 → x1.5.
- `reckless` → moves de retroceso → x1.2.
- `iron-fist` → punch moves → x1.2.
- `mega-launcher` → pulse moves → x1.5.
- `strong-jaw` → bite moves → x1.5.

### Placeholder/no-op documentado (resto hasta ~90)
- Habilidades sin infra (status conditions, PP, items): `truant`, `speed-boost`, `moody`, `regenerator`, `natural-cure`, etc. → no en este catálogo inicial. Fuera de scope.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/abilities_test.go` | Cada habilidad del catálogo: input → output esperado. Ej. Llovizna: `OnSwitchIn` setea clima Rain. Espesura: multiplier 1.5 si HP<1/3, 1.0 si >=1/3. |
| `core/battle_test.go` | Integración: Llovizna vs move Fire → daño reducido. Intimidación baja Atk rival antes del primer ataque. |
| `core/damage_test.go` | `Torrent` activado → Water x1.5. `Thick Fat` recibiendo Fire → daño /2. |

## Criterios de aceptacion

- [ ] `AbilityEffect` struct con hooks definidos.
- [ ] Catálogo ≥80 habilidades implementadas (las listadas arriba).
- [ ] `FullBattleInput` acepta `AttackerAbility`/`DefenderAbility` por nombre kebab-case.
- [ ] `SimulateFullBattle` dispara `OnSwitchIn` al comienzo.
- [ ] `CalculateDamage` respeta `ModifyAttack`, `ModifyDefense`, `ImmuneToType`.
- [ ] `OnTakeHit` permite curar (Water Absorb) o anular daño.
- [ ] Habilidad desconocida / vacía = no-op (no rompe).
- [ ] Tests unitarios por habilidad pasan.
- [ ] Tests de integración pasan.
- [ ] Logs i18n EN/ES para eventos de habilidad ("{Pokemon} tiene Intimidación", "El ataque de {rival} bajó", "Absorbió el agua", etc.).
- [ ] No rompe batallas existentes cuando ambos lados ausentes de habilidad.

## Notas

- Las habilidades que requieren sistemas no implementados (status conditions quemado/paralizado/envenenado, PP, items) quedan FUERA de este catálogo. Documentar como skipped.
- Prioridad en el catálogo: las más comunes y jugables, no exhaustividad.
- Mantener `AbilitiesCatalog` como `map[string]AbilityEffect` inmutable (inicializado una vez).
- Todas las funciones hook deben ser puras (value types, no punteros a estado mutable).
- Pokémon sin habilidad = string vacío → no-op.
