# Rellenar movimientos aleatorios solo con los aprendidos por nivel

**ID**: 0172-random-moves-level-up-only
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

Al rellenar movimientos aleatoriamente en el simulador (tanto en 1v1 atacante/defensor como en la generación aleatoria de miembros de equipo), los movimientos elegidos deben provenir **exclusivamente** del pool `Method == "level-up"` del Pokémon.

Los movimientos aprendidos por MT, huevo o tutor siguen disponibles para selección manual en los dropdowns de movimiento (no se filtran en la vista general), pero el botón "Rellenar aleatoriamente" y la generación de miembro aleatorio no los incluyen.

Si un Pokémon tiene menos de 4 movimientos por nivel, se rellena con los que haya (sin repetir) y los slots restantes quedan vacíos.

## Capas afectadas

- **Core**: `core/teams.go` — `GenerateRandomTeamMember` filtra `pokemon.Moves` por `Method == "level-up"` antes de elegir.
- **Shell**: ninguna.
- **APP/Frontend**: `frontend/src/pages/builds.ts` — `randomFillSlots` filtra `pokemon.Moves` por `Method === "level-up"` antes de elegir.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/teams.go` | modificar | En `GenerateRandomTeamMember`, reemplazar el pool de moves por `m.Method == "level-up"` solamente. |
| `core/teams_test.go` | modificar | Añadir test que verifica que los movimientos generados pertenecen todos al pool level-up, y que si hay menos de 4 level-up se rellena lo posible. |
| `frontend/src/pages/builds.ts` | modificar | En `randomFillSlots` (≈ línea 2027), filtrar `pokemon.Moves` por `m.Method === "level-up"` antes de shuffle/slice. |

## Plan de implementacion

1. **Core** (`core/teams.go:132-138`):
   ```go
   var movePool []string
   for _, m := range pokemon.Moves {
       if m.Method == "level-up" {
           movePool = append(movePool, m.Name)
       }
   }
   moves := pickRandomUnique(movePool, 4, rng)
   ```
   Eliminar la variable muerta `otherMoves` y la asignación trivial.

2. **Frontend** (`frontend/src/pages/builds.ts:2031`):
   ```ts
   const available = (pokemon.Moves ?? []).filter(m => m.Method === "level-up");
   ```
   Si `available.length === 0`, salir igual que hoy (no hay nada que rellenar).

3. **Tests core**: cubrir dos casos:
   - Pokémon con mezcla (level-up + machine + egg) → resultado subset de level-up.
   - Pokémon con <4 movimientos level-up → devuelve solo los disponibles, sin duplicados.

4. **QA manual**:
   - Equipos → "Rellenar aleatoriamente": los 4 movimientos del nuevo miembro son todos de nivel (verificar en modal de edición contra el detalle del Pokémon).
   - 1v1 atacante/defensor → "Rellenar aleatoriamente" + revisar que ninguno sea MT/huevo/tutor.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/teams_test.go` | `GenerateRandomTeamMember` solo devuelve moves con `Method=="level-up"`; comportamiento cuando hay <4 level-up. |
| Manual QA | 1v1 "Rellenar aleatoriamente" nunca introduce moves MT/huevo/tutor. |
| Manual QA | "Rellenar aleatoriamente" de equipo completo produce miembros con moves solo de nivel. |

## Criterios de aceptacion

- [ ] `GenerateRandomTeamMember` filtra por `Method=="level-up"`.
- [ ] `randomFillSlots` (builds 1v1) filtra por `Method==="level-up"`.
- [ ] Ningún movimiento MT/huevo/tutor aparece al rellenar aleatoriamente.
- [ ] Los dropdowns manuales de movimiento siguen mostrando toda la lista (sin cambios).
- [ ] Test unitario de core cubre el filtro y el caso de pool reducido.
- [ ] Probado en build de producción.
- [ ] Paridad Wails/REST intacta (no se toca el binding).

## Notas

- El valor `"level-up"` coincide con el usado por el scraper/PokéAPI en `PokemonMoveEntry.Method`; verificar en `core/domain.go:78` y en los datos actuales antes de commit.
- No tocar el flujo de "añadir movimiento manual": ahí el usuario elige de toda la lista para poder incluir MTs deliberadamente.
- No añadir fallback a otros Methods si no hay level-up: la especificación del usuario es que los slots queden vacíos y se rellenen a mano.
