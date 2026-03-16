# Ordenacion por columnas en tabla de encounters

**ID**: 0080-encounters-table-sorting
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

En la ficha individual de un Pokemon (detalle del Pokedex), la tabla de encounters no permite ordenar por columnas. Añadir ordenacion clickable en las cabeceras de la tabla, siguiendo el mismo patron que ya existe en las tablas de Moves, Abilities y la tabla del Pokedex.

## Capas afectadas

- **APP (frontend)**: Logica de ordenacion en la tabla de encounters y estilos de indicadores de sort.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | En la funcion `loadEncounters()`, hacer las cabeceras clickables con clase `.sortable`, añadir estado de ordenacion (`sortColumn`, `sortDirection`), implementar logica de sort sobre las filas aplanadas antes de renderizar, y añadir indicadores visuales `.sort-indicator` |
| `frontend/src/styles/_pokemon.scss` | modificar | Añadir estilos para cabeceras sortable y sort-indicator en `.encounters-table` (cursor pointer, iconos de flecha, hover) si no los hereda de `.poke-table` |
| `frontend/src/styles/_dark.scss` | modificar | Ajustar estilos de sort-indicator en dark mode para `.encounters-table` si es necesario |

## Plan de implementacion

1. En `pokedex.ts`, dentro de `loadEncounters()`, almacenar las filas aplanadas (array de objetos con location, game, method, chance, minLevel, maxLevel, conditions) en una variable accesible para re-renderizado.
2. Hacer las 6 cabeceras de la tabla (`Location`, `Game`, `Method`, `Chance`, `Levels`, `Conditions`) clickables con clase `sortable` y un `<span class="sort-indicator"></span>` dentro de cada `<th>`.
3. Añadir event listeners en cada `<th>` que ciclen el estado: `null` → `asc` → `desc` → `null` (mismo patron que `moves.ts` y `abilities.ts`).
4. Implementar funcion de sort que compare segun el tipo de columna:
   - **Location, Game, Method, Conditions**: comparacion alfabetica (`localeCompare`).
   - **Chance**: comparacion numerica (valor entero del porcentaje).
   - **Levels**: comparacion numerica por `minLevel`, desempate por `maxLevel`.
5. Tras ordenar (o resetear), re-renderizar el `<tbody>` con las filas en el nuevo orden sin reconstruir el `<thead>`.
6. Actualizar los indicadores visuales (clases `asc`/`desc`) en la cabecera activa y limpiar las demas.
7. Reutilizar los estilos de `.sortable` y `.sort-indicator` ya existentes en el proyecto (verificar si `.poke-table th` ya los tiene; si no, añadirlos scoped a `.encounters-table`).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Todas las cabeceras muestran cursor pointer y son clickables |
| (visual) | Al hacer click en una cabecera, las filas se reordenan ascendentemente |
| (visual) | Un segundo click en la misma cabecera invierte a descendente |
| (visual) | Un tercer click resetea al orden original |
| (visual) | El indicador de flecha aparece solo en la columna activa |
| (visual) | Ordenar por Chance ordena numericamente (5% < 10% < 25%) |
| (visual) | Ordenar por Levels ordena por nivel minimo y desempata por maximo |
| (visual) | Funciona correctamente en dark mode |
| (visual) | Funciona correctamente en mobile (responsive) |

## Criterios de aceptacion

- [x] Las 6 cabeceras de la tabla de encounters son clickables
- [x] Click en cabecera ordena asc → desc → reset (ciclo de 3 estados)
- [x] Columnas de texto (Location, Game, Method, Conditions) se ordenan alfabeticamente
- [x] Columna Chance se ordena numericamente
- [x] Columna Levels se ordena numericamente por nivel minimo
- [x] Indicador visual (flecha) muestra la direccion de ordenacion activa
- [x] Solo una columna puede estar activa a la vez
- [x] El sort no pierde ni duplica filas
- [x] Compatible con dark mode
- [x] Responsive: funciona en mobile

## Notas

- Seguir exactamente el patron de ordenacion de `explore/moves.ts` y `explore/abilities.ts` para consistencia.
- Las filas aplanadas ya se generan en `loadEncounters()` (location × version × method detail); usar ese array como fuente de verdad.
- No requiere cambios en backend: toda la ordenacion es client-side.
