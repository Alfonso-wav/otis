# Task 0164 — Ditto: movimiento "Transformacion" adopta stats y moves del rival

## Estado: done

## Goal
En el simulador de batallas, el movimiento **Transformation** de Ditto no tiene efecto actualmente. Implementar la regla canonica: tras usarlo, Ditto **adopta los stats y moves del rival** (salvo HP). Este es el gameplay central de Ditto.

## Contexto tecnico

### Backend (Core)
- Motor batalla: `core/battle.go` (560 lineas).
- `ExecuteTurn()` (lineas 197-243): solo distingue `move.Power > 0` vs `Category == "status"`. **No hay sistema de efectos** por movimiento — todos los moves con `Power == 0` se loggean como "sin efecto de daño".
- `SimulateFullBattle()` (lineas 38-80): loop de turnos, ambos lados eligen move random, resuelve orden por prioridad/velocidad.
- `FullBattleInput` (lineas 9-20): contiene stats, types, level, moves, name por lado. Stats son **value types** (`Stats` struct) — se pueden clonar sin problema.
- `Move` struct (`core/domain.go:14-26`): `Name, Type, Power, Accuracy, Priority, Category, Description`. No hay `Effect` field.

### Regla canonica de Transformation
Al ejecutarla sobre rival X:
- Ditto copia: stats (Atk, Def, SpA, SpD, Speed — **NO HP**), types, moves, (nivel del propio Ditto se mantiene en este simulador simplificado).
- No falla (accuracy=0 en juego real, always hits).
- PP en moves copiados quedaria a 5, pero este simulador no trackea PP — omitir.

### Diseño propuesto
Opciones:
1. **Special-case en ExecuteTurn**: `if move.Name == "transform"` → devolver `TurnResult` con flag `StatsOverride`/`MovesOverride`. Requiere que `SimulateFullBattle` aplique el override al input en turnos siguientes.
2. **Campo `State.AttackerStatsOverride`** (opcional) en `BattleState` leido por `CalculateDamage`.
3. **Mutar `FullBattleInput` entre turnos** (feo, rompe pureza pero pragmatico).

Preferido: **opcion 1 + 2** — `BattleState` gana campos opcionales `AttackerStatsOverride *Stats`, `AttackerMovesOverride []Move`, `AttackerTypesOverride []PokemonType`. `SimulateFullBattle` los respeta al elegir move y al calcular damage. `ExecuteTurn` al ejecutar Transform setea esos campos copiando los del defensor. Si lo usa el defensor, campos simetricos `DefenderStatsOverride`, etc.

### Nombre del move
Buscar name canonico en PokeAPI: `"transform"` (lowercase, kebab-case). Ditto lo aprende — confirmar via `GetMove("transform")` o via datos almacenados.

### Frontend
- Log de batalla lo mostrara automaticamente si `ExecuteTurn` retorna `LogEntry` adecuado (ej. "Ditto usó Transformation → se transformó en Pikachu").
- No requiere cambios de UI criticos.

### Tests
- Core es testeable trivialmente (funciones puras). Agregar test en `core/battle_test.go` (si no existe, crearlo):
  - Ditto(stats bajos) vs Pikachu(stats altos): tras Transform, siguiente turno usa stats de Pikachu.
  - Verificar que moves disponibles post-Transform son los de Pikachu.
  - Verificar que HP de Ditto **no** cambia.

## Acceptance criteria

- [ ] `Move.Name == "transform"` detectado en `ExecuteTurn` (o nueva funcion `applyMoveEffect`).
- [ ] Tras Transform, stats del usuario (atacante en ese turno) se sustituyen por los del rival en turnos posteriores (excepto HP).
- [ ] Moves del usuario pasan a ser los del rival en turnos posteriores.
- [ ] Types del usuario pasan a los del rival (afecta STAB y efectividad).
- [ ] HP del usuario **no** se modifica.
- [ ] Log incluye entrada clara: "Ditto usó Transformation → se transformó en {defender}".
- [ ] Simetrico: funciona si lo usa attacker o defender.
- [ ] Test unitario en `core/battle_test.go` verificando comportamiento.
- [ ] No rompe simulaciones de Pokemon que no tienen Transform.
- [ ] Funciona tanto en batalla manual (turno a turno) como en `SimulateFullBattle`.

## Archivos afectados

- `core/domain.go` — posible extension de `BattleState` con campos override opcionales.
- `core/battle.go` — logica de Transform en `ExecuteTurn`, respetar overrides en `SimulateFullBattle` y `CalculateDamage` lookup.
- `core/battle_test.go` — tests nuevos.
- `frontend/src/pages/builds.ts` — verificar que log se renderiza (probablemente no requiere cambio).
- `frontend/src/locales/{en,es}.json` — si se traduce el texto del log, agregar keys.

## Notas

- Mantener funciones puras: overrides viajan dentro de `BattleState`, no estado global.
- Si el simulador batch (`SimulateFullBattle`) no trackea overrides, los resultados estadisticos de Ditto seran incorrectos — asegurar que loop de turnos lee overrides al elegir move y calcular damage.
- Otros moves de efecto (Swords Dance, Recover, etc.) fuera de scope — solo Transform en esta tarea. Pero el sistema de overrides deberia ser extensible.
