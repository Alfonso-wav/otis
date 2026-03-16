# Filtro de Game para la tabla de Encounters

**ID**: 0086-encounters-game-filter
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Actualmente la tabla de Encounters muestra todas las filas de todos los juegos sin posibilidad de filtrar. Cuando un Pokémon aparece en muchos juegos, la tabla se vuelve muy larga y difícil de navegar.

**Solución**: Añadir un filtro de tipo dropdown multi-select (siguiendo el patrón existente de filtros en el Pokedex — generación/tipo) que permita seleccionar uno o más juegos (versions) para mostrar solo los encounters de esos juegos. El filtro se posiciona encima de la tabla de encounters.

## Capas afectadas

- **APP (frontend)**: UI del filtro y lógica de filtrado client-side en `pokedex.ts`.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Añadir dropdown multi-select de games sobre la tabla de encounters, extraer lista única de games de los datos cargados, filtrar `encounterRows` antes de renderizar |
| `frontend/src/styles/_pokemon.scss` | modificar | Estilos para el dropdown de filtro de game en la sección encounters |
| `frontend/src/styles/_dark.scss` | modificar | Estilos dark mode para el filtro de game |

## Plan de implementacion

### Parte 1: Extraer games disponibles

1. **En `pokedex.ts`**, después de obtener los datos de encounters y construir `encounterRows`, extraer la lista única de games (valores del campo `game`) y ordenarla alfabéticamente. Almacenar en una variable de estado (e.g. `availableGames: string[]`).

### Parte 2: UI del filtro

2. **Renderizar un dropdown multi-select** encima de la tabla de encounters, siguiendo el mismo patrón visual que los filtros de generación/tipo del Pokedex:
   - Un botón trigger con texto "Game" que muestra el conteo de seleccionados: "Game (3)".
   - Al hacer click, despliega una lista de chips/botones con cada game disponible.
   - Cada chip es toggleable (seleccionar/deseleccionar).
   - Incluir un botón de reset para limpiar la selección.
   - Cerrar el dropdown al hacer click fuera.

3. **Estado del filtro**: mantener un `Set<string>` o `string[]` con los games seleccionados (e.g. `selectedGames`).

### Parte 3: Lógica de filtrado

4. **Filtrar `encounterRows`** antes de renderizar la tabla:
   - Si `selectedGames` está vacío, mostrar todas las filas (sin filtro activo).
   - Si tiene valores, mostrar solo las filas cuyo campo `game` esté en `selectedGames`.

5. **Re-renderizar la tabla** cada vez que cambie la selección de games. Asegurarse de:
   - Mantener el sorting actual (si hay un `SortCache` para encounters, invalidarlo o re-aplicar con los datos filtrados).
   - Re-aplicar `reapplyColumnVisibility("encounters")` tras re-renderizar.

### Parte 4: Estilos

6. **Estilos en `_pokemon.scss`**: posicionar el filtro encima de la tabla, reutilizar los estilos existentes de dropdowns de filtro del Pokedex (`.filter-dropdown`, `.filter-chip`, etc.) o crear clases específicas si el contexto es diferente.

7. **Dark mode en `_dark.scss`**: asegurar que el filtro se vea correcto en modo oscuro.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | El dropdown de Game aparece encima de la tabla de encounters |
| (visual) | La lista de games se extrae correctamente de los datos cargados (sin duplicados, ordenada) |
| (visual) | Seleccionar un game filtra la tabla mostrando solo encounters de ese game |
| (visual) | Seleccionar múltiples games muestra encounters de todos los seleccionados |
| (visual) | Deseleccionar todos los games vuelve a mostrar todas las filas |
| (visual) | El botón reset limpia la selección |
| (visual) | El sorting sigue funcionando correctamente con el filtro activo |
| (visual) | La visibilidad de columnas se mantiene tras filtrar |
| (visual) | Dark mode se ve correctamente |
| (visual) | En móvil el dropdown es usable y no se desborda |

## Criterios de aceptacion

- [ ] Existe un dropdown multi-select de "Game" encima de la tabla de encounters
- [ ] La lista de games se genera dinámicamente de los datos del Pokémon actual
- [ ] Seleccionar games filtra correctamente las filas de la tabla
- [ ] El conteo de games seleccionados se muestra en el trigger del dropdown
- [ ] El botón reset limpia el filtro
- [ ] El sorting funciona correctamente con datos filtrados
- [ ] La visibilidad de columnas se mantiene tras filtrar
- [ ] Estilos correctos en light y dark mode
- [ ] Responsive: usable en móvil

## Notas

- Los valores de `game` vienen del campo `Version` de `VersionEncounter` en el backend (e.g. "red", "blue", "yellow", "gold", etc.).
- Al cambiar de Pokémon en el modal, el filtro debe resetearse ya que cada Pokémon tiene encounters en juegos distintos.
- Reutilizar al máximo los patrones de UI y estilos existentes de los filtros del Pokedex para mantener consistencia visual.
