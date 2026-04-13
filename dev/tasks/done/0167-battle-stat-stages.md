# Battle: sistema de stat stages (-6..+6)

**ID**: 0167-battle-stat-stages
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

El motor actual no modela stat stages. Habilidades como Intimidación y movimientos como Danza Espadas, Gruñido, Pantalla de Humo dependen de ello. Implementar sistema canónico de stages: -6..+6 por stat (Atk, Def, SpA, SpD, Spe, Accuracy, Evasion) con multiplicadores canónicos y aplicación en cálculo de daño / orden.

Prerequisito para 0168 (habilidades).

## Capas afectadas

- **Core**: `StatStages` struct en `BattleState`, helpers `ApplyStage` / `StageMultiplier`, integración en `CalculateDamage` (Atk/Def/SpA/SpD) y `resolveOrder` (Spe).
- **Shell**: ninguno.
- **APP**: logs de batalla reflejan cambios de stage.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir `StatStages` (Atk, Def, SpA, SpD, Spe, Acc, Eva int, rango -6..+6). Campos `AttackerStages StatStages`, `DefenderStages StatStages` en `BattleState`. |
| `core/stages.go` | crear | Funciones puras: `ApplyStage(stages, stat, delta) StatStages` (clamp -6..+6), `StageMultiplier(stage int) float64` (tabla canónica 2/8 .. 8/2 para Atk/Def/SpA/SpD/Spe; 3/3..9/3 para Acc/Eva). |
| `core/battle.go` | modificar | `ExecuteTurn` detecta moves de stat (`swords-dance`, `growl`, `leer`, `harden`, `agility`, etc.) y produce `BattleState` con stages actualizados. `resolveOrder` usa Spe * stage multiplier. |
| `core/damage.go` | modificar | `CalculateDamage` multiplica Atk/Def/SpA/SpD efectivos por `StageMultiplier` antes del cálculo. |
| `core/stages_test.go` | crear | Tests unitarios: clamp, tabla de multiplicadores, composición. |
| `core/battle_test.go` | modificar | Tests: Swords Dance sube Atk +2; siguiente ataque hace más daño. Growl baja Atk rival -1. |

## Plan de implementacion

1. Crear `StatStages` struct y funciones puras en `core/stages.go` (sin dependencias).
2. Tabla canónica: `multiplier(n) = max(2, 2+n)/max(2, 2-n)` para Atk/Def/SpA/SpD/Spe; para Acc/Eva la tabla es `max(3, 3+n)/max(3, 3-n)`.
3. Extender `BattleState` con ambos stages (valor inicial todo 0).
4. Catálogo de moves de stat stage (por `Name` y efecto conocido). Grupo mínimo inicial:
   - Self-boost: `swords-dance` (+2 Atk), `agility` (+2 Spe), `nasty-plot` (+2 SpA), `calm-mind` (+1 SpA/+1 SpD), `bulk-up` (+1 Atk/+1 Def), `iron-defense` (+2 Def), `amnesia` (+2 SpD), `barrier` (+2 Def).
   - Opponent debuff: `growl` (-1 Atk), `leer` (-1 Def), `tail-whip` (-1 Def), `string-shot` (-1 Spe), `smokescreen` (-1 Acc), `sand-attack` (-1 Acc).
5. Logs i18n: "Ataque de X subió", "Defensa de Y bajó mucho".
6. Tests unitarios puros.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/stages_test.go` | `StageMultiplier(0)=1.0`, `(+2)=2.0`, `(-2)=0.5`, `(+6)=4.0`, `(-6)=0.25`. Clamp -6..+6. |
| `core/battle_test.go` | Swords Dance → Atk efectivo x2 en siguiente ataque. Growl al rival → Atk rival x2/3. |

## Criterios de aceptacion

- [ ] `StatStages` struct + funciones puras.
- [ ] `BattleState` tracker ambos lados.
- [ ] `CalculateDamage` respeta stages Atk/Def/SpA/SpD.
- [ ] `resolveOrder` respeta stages Spe.
- [ ] Catálogo mínimo de moves de stat stage aplicado.
- [ ] Logs i18n EN/ES ("subió", "subió mucho", "bajó", "bajó mucho", "no puede subir más").
- [ ] Tests unitarios pasan.
- [ ] Stages no se resetean dentro de la batalla (persisten entre turnos).
- [ ] Stages se resetean entre batallas distintas (nuevo `InitBattle`).

## Notas

- Accuracy/Evasion tabla distinta (base 3/3, max 9/3).
- No incluir inversión (Topsy-Turvy) ni reset (Haze) en esta tarea — scope acotado, extensible.
- Prerequisito de: 0168 (Intimidación baja Atk rival -1 al entrar).
