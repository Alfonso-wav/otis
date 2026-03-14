# 0018 — Movimientos del defensor y simulación automática de batalla

## Estado

todo

## Descripción

Mejorar el simulador de batallas con tres features:

1. **Movimientos del defensor**: Añadir 4 slots de movimiento para el defensor (igual que el atacante), con un botón "Rellenar aleatoriamente" que elige 4 movimientos compatibles con ese Pokémon al azar.
2. **El defensor ataca**: En la simulación turno a turno, cuando es el turno del defensor, se muestran y usan sus propios movimientos (no los del atacante como ahora).
3. **Botón "Simular batalla completa"**: Simula toda la batalla de forma aleatoria (ambos lados eligen movimientos al azar cada turno) y muestra el log completo y el ganador.

## Contexto

- El estado actual (`builds.ts`) solo tiene `slots` para el atacante. El defensor no tiene movimientos configurables.
- En la batalla, cuando es el turno del defensor, se muestran los mismos botones del atacante — esto es incorrecto.
- `core/battle.go` ya tiene `ExecuteTurn` (función pura) e `InitBattle`.
- El defensor ya tiene `state.defender.Moves` disponible (lista de movimientos aprendibles del Pokémon).
- `GetMove(name)` ya existe en el backend para cargar el detalle de cada movimiento.

## Capas involucradas

- **Core**: Añadir función pura `SimulateFullBattle` en `core/battle.go`
- **APP (bindings)**: Exponer `SimulateFullBattle` en `app/bindings.go`
- **Frontend**: Refactor de estado y UI en `frontend/src/pages/builds.ts`

## Plan de implementación

### Paso 1 — Core: `SimulateFullBattle` (`core/battle.go`)

Añadir tipo de input y función pura:

```go
type FullBattleInput struct {
    AttackerStats  Stats         `json:"attackerStats"`
    DefenderStats  Stats         `json:"defenderStats"`
    AttackerTypes  []PokemonType `json:"attackerTypes"`
    DefenderTypes  []PokemonType `json:"defenderTypes"`
    AttackerLevel  int           `json:"attackerLevel"`
    DefenderLevel  int           `json:"defenderLevel"`
    AttackerMoves  []Move        `json:"attackerMoves"`
    DefenderMoves  []Move        `json:"defenderMoves"`
}

func SimulateFullBattle(input FullBattleInput, randSource func(n int) int) BattleState
```

- `randSource` es una función para inyectar aleatoriedad (facilita tests).
- Cada turno: atacante elige movimiento random de `AttackerMoves`, defensor de `DefenderMoves`.
- El turno del atacante usa `ExecuteTurn` con los stats/tipos del atacante → aplica daño al HP del defensor.
- El turno del defensor usa `ExecuteTurn` con perspectiva invertida → aplica daño al HP del atacante.
- Se alterna atacante/defensor hasta que `IsOver = true` (HP ≤ 0) o max 200 turnos (anti-loop infinito).
- El campo `Winner` es `"attacker"` o `"defender"` según quién sobreviva.
- Añadir tests en `core/battle_test.go`.

### Paso 2 — APP: binding (`app/bindings.go`)

```go
func (a *App) SimulateFullBattle(input core.FullBattleInput) core.BattleState {
    return core.SimulateFullBattle(input, func(n int) int {
        return rand.Intn(n)
    })
}
```

### Paso 3 — Frontend: estado defensor con movimientos

En `builds.ts`, añadir al estado:

```typescript
defenderSlots: [BuildSlot, BuildSlot, BuildSlot, BuildSlot];
```

Extraer la lógica de `renderMoveSlots` para que acepte un parámetro `prefix: "atk" | "def"` y use los slots y Pokémon correspondientes.

Añadir botón "Rellenar aleatoriamente" en la sección de movimientos del defensor:
- Filtra `state.defender.Moves` para quedarse solo con los que son aprendibles (todos los que aparezcan).
- Elige hasta 4 al azar.
- Llama `GetMove(name)` para cada uno y rellena `state.defenderSlots`.

### Paso 4 — Frontend: corrección del turno del defensor

En `renderBattleSection` y `handleMoveClick`:

- Cuando `phase === "defender-turn"`, mostrar los botones de `state.defenderSlots` (no `state.slots`).
- La lógica de `handleMoveClick` ya hace la inversión de perspectiva de HP correctamente; solo hay que asegurarse de que el movimiento seleccionado provenga del slot correcto (atacante o defensor según el turno).

### Paso 5 — Frontend: botón "Simular batalla completa"

En `renderBattleSection`, cuando la batalla está activa y hay movimientos en ambos lados, mostrar el botón "Simular batalla entera" que:

1. Recoge los moves de atacante (`state.slots`) y defensor (`state.defenderSlots`) que no sean null.
2. Si alguno de los dos no tiene movimientos, mostrar alerta y abortar.
3. Llama `SimulateFullBattle(input)` y recibe el `BattleState` final.
4. Actualiza `battleUI.battleState` con el resultado y `battleUI.phase = "over"`.
5. Muestra el log completo y el banner de ganador.

También exponer el botón en la fase `"idle"` si ambos lados tienen movimientos configurados, para simular directamente sin hacer turnos manuales.

### Paso 6 — SCSS

Añadir estilos para:
- `.defender-moves-section`: sección de movimientos del defensor (mismo estilo que la del atacante).
- `.battle-random-fill-btn`: botón de rellenar aleatoriamente (estilo secundario/outline).
- `.battle-auto-btn`: botón de simular batalla completa (estilo destacado/accent).

## Criterios de aceptación

- [ ] El defensor tiene su propia sección de 4 slots de movimiento.
- [ ] El botón "Rellenar aleatoriamente" del defensor carga 4 movimientos compatibles random.
- [ ] En la simulación turno a turno, el turno del defensor muestra y usa sus propios movimientos.
- [ ] El botón "Simular batalla completa" ejecuta toda la batalla de forma aleatoria y muestra el resultado.
- [ ] `SimulateFullBattle` en Core es función pura con tests.
- [ ] Si un lado no tiene movimientos, "Simular batalla completa" no está disponible o muestra aviso.
- [ ] El log de la batalla completa muestra todos los turnos (o los últimos N si son muchos).

## Notas

- El "rellenar aleatoriamente" para el defensor debe filtrar movimientos compatibles con ese Pokémon (`state.defender.Moves`), no una lista global de todos los movimientos.
- `SimulateFullBattle` en Core no tiene efectos secundarios: recibe `randSource` inyectado para que sea testeable.
- No cambiar la lógica existente de la tabla de daño ni los slots del atacante.
- El atacante también debería tener opción de rellenar aleatoriamente para consistencia — incluir en esta misma tarea.
