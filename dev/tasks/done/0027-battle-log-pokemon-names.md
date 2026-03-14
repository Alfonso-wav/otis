# 0027 — Mostrar nombre del Pokémon en cada línea del log de batalla

## Estado

todo

## Descripción

En la terminal del simulador de batallas, cada línea del log muestra `[T1] usó Move → ...` indicando el turno, pero no se identifica qué Pokémon está realizando el movimiento. Se necesita incluir el nombre del Pokémon atacante en cada entrada del log para que el usuario sepa quién actúa en cada turno.

Formato deseado: `[T1] Charizard usó Flamethrower → ...`

## Contexto

- `core/battle.go` genera las entradas del log en `ExecuteTurn` (líneas 202, 240, 268) usando `fmt.Sprintf`.
- Actualmente ni `TurnInput` ni `FullBattleInput` tienen campos para los nombres de los Pokémon.
- El frontend (`frontend/src/pages/builds.ts`) pasa los datos de batalla al backend pero no incluye nombres.
- `app/bindings.go` expone `ExecuteTurn` y `SimulateFullBattle` al frontend vía Wails.

## Capas involucradas

- **Core**: `core/battle.go` — añadir campo `AttackerName` a `TurnInput` y `FullBattleInput`, incluirlo en los `fmt.Sprintf` del log.
- **APP**: `app/bindings.go` — propagar el nuevo campo si cambian las firmas.
- **Frontend**: `frontend/src/pages/builds.ts` — pasar los nombres de los Pokémon seleccionados al invocar la batalla.

## Plan de implementación

### Paso 1 — Core: Añadir nombres a los structs de entrada (`core/battle.go`)

- Añadir campo `AttackerName string` a `FullBattleInput` y `DefenderName string`.
- Añadir campo `AttackerName string` a `TurnInput`.
- [ ] Completado

### Paso 2 — Core: Incluir nombre en las entradas del log (`core/battle.go`)

Modificar los 3 puntos donde se genera `logEntry` en `ExecuteTurn`:

1. **Fallo por precisión** (línea ~202):
   - De: `[T%d] usó %s → ¡Falló!`
   - A: `[T%d] %s usó %s → ¡Falló!`

2. **Movimiento de estado** (línea ~240):
   - De: `[T%d] usó %s → sin efecto de daño | HP Defensor: %d/%d`
   - A: `[T%d] %s usó %s → sin efecto de daño | HP Defensor: %d/%d`

3. **Movimiento de daño** (línea ~268):
   - De: `%s[T%d] usó %s → %d daño (%s) | HP Defensor: %d/%d`
   - A: `%s[T%d] %s usó %s → %d daño (%s) | HP Defensor: %d/%d`

- [ ] Completado

### Paso 3 — Core: Propagar nombres en SimulateFullBattle (`core/battle.go`)

- En `executeAttackerTurn`, pasar `AttackerName` del `FullBattleInput` al `TurnInput`.
- En `executeDefenderTurn`, pasar `DefenderName` del `FullBattleInput` como `AttackerName` del `TurnInput` (porque el defensor actúa como atacante en su turno).

- [ ] Completado

### Paso 4 — APP: Actualizar bindings si es necesario (`app/bindings.go`)

- Verificar que los bindings de Wails exponen correctamente los nuevos campos JSON.
- No debería requerir cambios si los structs se serializan automáticamente.

- [ ] Completado

### Paso 5 — Frontend: Enviar nombres de Pokémon (`frontend/src/pages/builds.ts`)

- Al construir `FullBattleInput` para `SimulateFullBattle`, incluir `attackerName` y `defenderName` con los nombres de los Pokémon seleccionados.
- Al construir `TurnInput` para `ExecuteTurn` (si se usa modo turno a turno), incluir `attackerName`.

- [ ] Completado

### Paso 6 — Tests: Actualizar tests existentes (`core/battle_test.go`)

- Actualizar los tests para incluir `AttackerName` en los inputs.
- Verificar que las entradas del log contienen el nombre del Pokémon.

- [ ] Completado

## Criterios de aceptación

- [ ] Cada línea del log de batalla incluye el nombre del Pokémon que realiza la acción
- [ ] El formato es `[T#] NombrePokemon usó Move → ...`
- [ ] Funciona tanto en simulación completa como turno a turno
- [ ] Los tests existentes siguen pasando
- [ ] El frontend envía los nombres correctos de los Pokémon seleccionados

## Notas

- Cambio sencillo y localizado. No afecta la lógica de combate, solo la presentación del log.
- Los nombres ya están disponibles en el frontend (el usuario los selecciona), solo falta propagarlos al backend.
