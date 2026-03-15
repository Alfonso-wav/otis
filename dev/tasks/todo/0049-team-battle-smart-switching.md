# Orden aleatorio inicial y cambio inteligente por debilidades de tipo

**ID**: 0049-team-battle-smart-switching
**Estado**: todo
**Fecha**: 2026-03-15

---

## Descripcion

Actualmente la batalla por equipos usa orden secuencial fijo (idx1=0, idx2=0 e incrementa al perdedor). Se necesitan dos cambios:

1. **Orden de comienzo aleatorio**: al iniciar la batalla, aleatorizar el orden de los miembros de cada equipo.
2. **Cambio inteligente**: cuando un Pokémon es derrotado, en vez de enviar al siguiente en la lista, seleccionar al miembro del equipo que tenga mejor ventaja de tipo contra el oponente actual (usando efectividad de tipos).

## Capas afectadas

- **Core**: Nueva función pura para seleccionar el mejor miembro según ventaja de tipo. Modificar `SimulateTeamBattle` para barajar el orden inicial y usar selección inteligente al cambiar.
- **Shell**: No afectada.
- **APP**: No afectada.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/battle.go` | modificar | Añadir shuffling inicial de miembros y lógica de selección inteligente por tipo |
| `core/types.go` | revisar | Verificar que `TypeEffectiveness` está accesible para calcular ventaja |
| `core/battle_test.go` | modificar | Tests para verificar que el cambio inteligente elige miembros con ventaja de tipo |

## Plan de implementacion

1. Crear función pura `ChooseBestMember(available []TeamBattleMember, opponentTypes []PokemonType) int` que retorna el índice del miembro con mayor ventaja de tipo contra el oponente.
2. La ventaja se calcula sumando los multiplicadores de efectividad de los tipos de ataque del miembro contra los tipos del defensor.
3. Modificar `SimulateTeamBattle`: al inicio, barajar ambos equipos con `randSource`. Cuando un miembro cae, usar `ChooseBestMember` en vez de incrementar índice secuencialmente.
4. Cambiar la estructura interna para manejar un slice de miembros disponibles en vez de un simple índice.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/battle_test.go` | Que `ChooseBestMember` selecciona al miembro con ventaja de tipo. Que el orden inicial es aleatorio (con seed fija). |

## Criterios de aceptacion

- [ ] El orden de inicio de los miembros es aleatorio en cada batalla
- [ ] Al caer un Pokémon, se envía al que tiene mejor ventaja de tipo contra el oponente actual
- [ ] Si no hay ventaja clara, se elige cualquiera de los disponibles
- [ ] Los tests verifican la selección inteligente con tipos conocidos
- [ ] La función `ChooseBestMember` es pura (sin efectos secundarios)

## Notas

La tabla de efectividad de tipos ya existe en `core/types.go` (`TypeEffectiveness`). La ventaja se puede calcular como el producto de multiplicadores de los tipos del atacante contra los tipos del defensor, considerando los movimientos disponibles del miembro candidato.
