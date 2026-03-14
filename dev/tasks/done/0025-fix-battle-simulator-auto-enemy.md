# 0025 — Arreglar simulador manual de batallas + respuesta automática del enemigo

## Estado

done

## Descripción

El simulador manual de batallas está roto por un bug en el frontend: la variable `isAttackerTurn` se usa antes de ser declarada en `handleMoveClick()` (`builds.ts:422` usa la variable, pero se declara en la línea 426). Esto causa un `ReferenceError` en tiempo de ejecución que impide seleccionar movimientos.

Además, cambiar el flujo de batalla manual para que cuando el usuario haga clic en un movimiento del atacante, el defensor **responda automáticamente** con un movimiento aleatorio de sus slots. Actualmente el sistema alterna turnos manualmente (el usuario elige movimiento del atacante, luego del defensor). El nuevo flujo es: usuario elige movimiento → se ejecuta turno del atacante → automáticamente se ejecuta turno del defensor con un movimiento aleatorio → se actualizan HP y log de ambos turnos.

## Contexto

- **Bug crítico**: en `frontend/src/pages/builds.ts`, función `handleMoveClick()`, línea 422 usa `isAttackerTurn` antes de su declaración en línea 426.
- El sistema de turnos actual alterna `attacker-turn` / `defender-turn` en `battleUI.phase`, requiriendo que el usuario elija manualmente el movimiento de ambos lados.
- `core/battle.go` ya tiene `ExecuteTurn` como función pura que recibe un `TurnInput` y devuelve `TurnResult`.
- El defensor ya tiene `state.defenderSlots` con hasta 4 movimientos configurados.
- `SimulateFullBattle` en el backend ya implementa la lógica de elegir movimiento aleatorio para cada lado — se puede reusar el patrón para el turno del defensor.

## Capas involucradas

- **Frontend (APP)**: arreglar bug + refactorizar `handleMoveClick()` en `frontend/src/pages/builds.ts`

## Plan de implementación

### Paso 1 — Arreglar el bug de declaración de variable

En `handleMoveClick()` de `builds.ts`, mover la declaración `const isAttackerTurn = battleUI.phase === "attacker-turn"` antes de su primer uso (antes de la línea que accede `activeSlots`).

### Paso 2 — Implementar respuesta automática del defensor

Refactorizar `handleMoveClick()` para que:

1. Siempre ejecute el turno del atacante (el usuario elige su movimiento).
2. Si la batalla no ha terminado tras el turno del atacante, seleccione automáticamente un movimiento aleatorio de `state.defenderSlots` (filtrando solo los que tengan `.move` definido).
3. Ejecute el turno del defensor con ese movimiento aleatorio (usando la misma lógica de swap de perspectiva HP que ya existe).
4. Actualice `battleUI.battleState` con el resultado de ambos turnos.
5. Eliminar la fase `defender-turn` — ahora siempre es `attacker-turn` o `over`.

### Paso 3 — Actualizar renderizado de UI

- Eliminar la lógica que muestra los botones del defensor como clicables para elegir movimiento.
- Los botones del atacante siempre están activos durante la batalla (si no ha terminado).
- El log muestra ambos turnos (atacante + respuesta del defensor) en secuencia.
- Añadir una breve indicación visual/textual del movimiento que eligió el defensor automáticamente.

### Paso 4 — Manejar edge cases

- Si el defensor no tiene movimientos configurados en sus slots, mostrar mensaje de error o no iniciar la batalla.
- Si el atacante derrota al defensor en su turno, no ejecutar el turno del defensor.
- Si el defensor derrota al atacante en su respuesta, marcar la batalla como terminada.

## Criterios de aceptación

- [ ] Hacer clic en un movimiento del atacante ya no lanza error (bug de `isAttackerTurn` corregido)
- [ ] Tras elegir movimiento del atacante, el defensor responde automáticamente con un movimiento aleatorio
- [ ] El log muestra ambos turnos en secuencia (atacante + defensor)
- [ ] Si el atacante noquea al defensor, no se ejecuta el turno de respuesta
- [ ] Si el defensor noquea al atacante, se muestra correctamente el ganador
- [ ] No se puede iniciar batalla si el defensor no tiene movimientos configurados
- [ ] Los tests existentes de `core/battle.go` siguen pasando (no se toca el Core)

## Notas

- No se modifica `core/battle.go` ni `app/bindings.go` — los cambios son solo en frontend.
- El patrón de swap de perspectiva HP para el turno del defensor ya existe en `handleMoveClick()` — se mantiene.
- La simulación completa (`SimulateFullBattle`) sigue funcionando igual — no se toca.
