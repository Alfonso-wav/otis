# Battle: sistema de clima (Rain/Sun/Sandstorm/Hail)

**ID**: 0166-battle-weather-system
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

El motor de batalla actual (`core/battle.go`) no modela clima. Este es prerequisito para habilidades tipo Llovizna/Sequía/Arena/Nevada y para los movimientos climáticos (Danza Lluvia, Día Soleado, Tormenta Arena, Granizo). Implementar sistema de clima canónico: estado + turnos restantes + modificadores en cálculo de daño + efectos pasivos por turno.

Es la base de la tarea 0168 (habilidades) — habilidades meteorológicas como Llovizna aplicarán clima pasivamente al entrar en combate.

## Capas afectadas

- **Core**: `Weather` enum, extensión de `BattleState`, modificador en `CalculateDamage`, efectos de fin de turno (daño Arena/Granizo), moves que setean clima.
- **Shell**: ninguno directo (datos de moves ya existen en PokéAPI).
- **APP**: logs de batalla emiten mensajes de clima. Frontend renderiza sin cambios (lee del log).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir tipo `Weather` (None/Rain/Sun/Sandstorm/Hail) y campos `Weather Weather`, `WeatherTurnsLeft int` en `BattleState`. |
| `core/battle.go` | modificar | `ExecuteTurn` detecta moves climáticos (`rain-dance`, `sunny-day`, `sandstorm`, `hail`) y setea clima 5 turnos. `SimulateFullBattle` decrementa turnos al final de cada ronda y aplica daño residual de Sandstorm/Hail a tipos no inmunes. |
| `core/damage.go` | modificar | `CalculateDamage` acepta clima; aplica x1.5 a Water en Rain / Fire en Sun, x0.5 a Fire en Rain / Water en Sun. Sandstorm sube SpD x1.5 a tipo Rock. |
| `core/battle_test.go` | modificar | Tests: moves climáticos setean clima, clima expira a los 5 turnos, daño residual correcto, modificador STAB+clima correcto. |
| `core/damage_test.go` | modificar | Tests adicionales con parámetro clima. |

## Plan de implementacion

1. Definir `type Weather string` con constantes `WeatherNone`, `WeatherRain`, `WeatherSun`, `WeatherSandstorm`, `WeatherHail`.
2. Extender `BattleState` con `Weather`, `WeatherTurnsLeft` (value types, inmutables turno a turno).
3. Refactorizar `CalculateDamage` para aceptar `Weather` (parámetro explícito o campo en `DamageInput`).
4. En `ExecuteTurn` (o helper nuevo `applyMoveEffect`), detectar moves climáticos por `Name` y producir nuevo `BattleState` con clima seteado y 5 turnos.
5. En `SimulateFullBattle` al cierre de cada turno completo: decrementar `WeatherTurnsLeft`; si llega a 0 → `WeatherNone` + log "El clima ha vuelto a la normalidad".
6. Daño residual: fin de turno si Sandstorm → dañar (1/16 HP) a todos menos tipos Rock/Ground/Steel. Si Hail → dañar (1/16 HP) a todos menos tipo Ice.
7. Logs i18n: emitir entradas "Empezó a llover", "La lluvia cesó", "La tormenta de arena hirió a {nombre}", etc.
8. Tests unitarios puros (sin mocks).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/battle_test.go` | Rain Dance → clima Rain 5 turnos; expira al turno 6. Sandstorm aplica daño residual a non-Rock. Hail respeta inmunidad Ice. |
| `core/damage_test.go` | Water move en Rain = x1.5; Fire en Rain = x0.5; Rock SpD x1.5 en Sandstorm. |

## Criterios de aceptacion

- [ ] `Weather` enum + campos en `BattleState`.
- [ ] Moves `rain-dance`, `sunny-day`, `sandstorm`, `hail` setean clima 5 turnos.
- [ ] Modificadores de daño aplicados en `CalculateDamage`.
- [ ] Daño residual Sandstorm/Hail al final de turno con inmunidades correctas.
- [ ] Clima expira automáticamente; log emitido.
- [ ] Logs i18n en EN/ES (`locales/en.json`, `locales/es.json`) para eventos de clima.
- [ ] Tests unitarios pasan.
- [ ] No rompe simulaciones existentes (incluida Transform de 0164).
- [ ] Funciona igual en batalla turno-a-turno y `SimulateFullBattle`.

## Notas

- Mantener funciones puras: el clima viaja dentro de `BattleState`, no estado global.
- No incluir climas gen 6+ (Niebla Misteriosa, Fuerte Lluvia, Sol Fuerte) en esta tarea — scope acotado.
- Sin sistema de objetos/rocas (Heat Rock etc.) — todos los climas duran 5 turnos fijos.
- Prerequisito de: 0168 abilities catalog (Llovizna, Sequía, Clima Despejado, etc.).
