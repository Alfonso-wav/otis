# Guardar y editar movimientos de Pokémon en equipos

**ID**: 0051-save-moves-with-team-members
**Estado**: todo
**Fecha**: 2026-03-15

---

## Descripcion

Actualmente los Pokémon no se guardan en los equipos con sus movimientos. El campo `Moves []string` en `TeamMember` existe en la estructura pero no se está rellenando al guardar un miembro en el equipo. Esto causa que en batalla se use el fallback a "struggle".

Se necesita:
1. **Guardar movimientos**: Al añadir un Pokémon a un equipo, guardar también los movimientos seleccionados.
2. **Editar movimientos**: Al seleccionar un Pokémon ya guardado en un equipo, poder editar sus movimientos.

## Capas afectadas

- **Core**: No afectada (la estructura `TeamMember.Moves` ya existe).
- **Shell**: No afectada (la persistencia ya serializa `Moves` en JSON).
- **APP**: Frontend — modificar la lógica de guardado para incluir los movimientos, y añadir UI de edición de movimientos para miembros existentes.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | 1) Al guardar un miembro, incluir los movimientos seleccionados en el objeto `TeamMember`. 2) Al seleccionar un miembro existente del equipo, mostrar sus movimientos como editables. 3) Al editar movimientos, persistir el cambio vía binding de guardado. |
| `app/bindings.go` | revisar | Verificar que `SaveTeam`/`UpdateTeam` persiste correctamente los movimientos (probablemente ya funciona si el frontend envía los datos) |
| `frontend/src/styles/_builds.scss` | modificar | Estilos para la sección de edición de movimientos en miembros del equipo |

## Plan de implementacion

1. Identificar dónde en `builds.ts` se construye el objeto `TeamMember` al guardar — verificar que `Moves` se incluye con los valores de los selectores de movimiento.
2. Al renderizar un miembro del equipo ya guardado, mostrar selectores de movimientos prellenados con los movimientos guardados.
3. Añadir lógica para que al cambiar un movimiento de un miembro guardado se actualice el equipo y se guarde.
4. Verificar el flujo completo: guardar con movimientos → cargar equipo → ver movimientos → editar → guardar de nuevo.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| N/A | Test manual del flujo: guardar Pokémon con movimientos, recargar, verificar que los movimientos persisten y son editables |

## Criterios de aceptacion

- [ ] Al añadir un Pokémon a un equipo, sus movimientos seleccionados se guardan en el JSON
- [ ] Al cargar un equipo, los movimientos de cada miembro se muestran correctamente
- [ ] Los movimientos son editables desde la vista del equipo (selectores de movimiento funcionales)
- [ ] Al cambiar un movimiento y guardar, el cambio persiste
- [ ] En batalla, los Pokémon usan sus movimientos guardados (no "struggle")

## Notas

La estructura de datos ya soporta movimientos (`TeamMember.Moves []string`). El `resolveTeamBattleMembers` en `app/bindings.go` ya resuelve los nombres de movimientos a objetos `Move` completos, y ya tiene fallback a "struggle" si no hay movimientos. El problema es exclusivamente que el frontend no envía los movimientos al guardar.
