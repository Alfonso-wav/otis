# Autocompletado en el buscador principal de la Pokédex

**ID**: 0088-pokedex-search-autocomplete
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Integrar la función de autocompletado existente (`createAutocomplete`) en el buscador principal de la Pokédex (`#search-input`), ubicado en la parte superior derecha del header a la altura del título "Pokédex". Actualmente el buscador solo permite búsqueda exacta al pulsar Enter/Buscar; con esta tarea, al escribir aparecerá un dropdown con sugerencias de nombres de Pokémon.

## Capas afectadas

- **APP (frontend)**: `pokedex.ts` — integrar `createAutocomplete` con el input de búsqueda existente.
- **Backend/Core/Shell**: sin cambios. Se reutiliza `ListPokemon` ya disponible en la API.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Importar `createAutocomplete`, cargar lista de nombres con `ListPokemon(0, 2000)` al inicializar y conectar el autocomplete al `#search-input` |
| `frontend/src/styles/_components.scss` | modificar (si necesario) | Ajustar estilos del dropdown autocomplete para que se vea bien en el contexto del header (z-index, posicionamiento) |
| `frontend/index.html` | modificar (si necesario) | Ajustar el contenedor `#search-bar` si el position relative no está aplicado correctamente para el dropdown |

## Plan de implementacion

### Parte 1: Cargar lista de nombres de Pokémon

1. En `initPokedex()` de `pokedex.ts`, llamar a `ListPokemon(0, 2000)` para obtener todos los nombres de Pokémon.
2. Almacenar los nombres en una variable de módulo (e.g., `pokemonNames: string[]`).
3. Esta carga puede ser asíncrona y no bloquear la inicialización del listado.

### Parte 2: Integrar `createAutocomplete` en el search input

4. Importar `createAutocomplete` desde `../autocomplete`.
5. Una vez cargados los nombres, llamar `createAutocomplete(searchInput, pokemonNames, onSelect)`.
6. El callback `onSelect` debe:
   - Actualizar el valor del input con el nombre seleccionado.
   - Llamar a `showDetail(name)` para cargar y mostrar el Pokémon directamente.

### Parte 3: Coordinar con el comportamiento existente de búsqueda

7. Asegurar que si el usuario selecciona una sugerencia con Enter (dentro del autocomplete), NO se dispare también el handler de `search()` que escucha `keydown Enter` en el input. El módulo `autocomplete.ts` ya hace `e.stopPropagation()` en este caso (línea 89).
8. El botón "Buscar" y Enter sin sugerencia activa siguen funcionando como antes (búsqueda exacta).

### Parte 4: Ajustes de estilo

9. Verificar que el dropdown autocomplete se posiciona correctamente debajo del input del header. El contenedor `#search-bar` debe tener `position: relative` para que el dropdown absoluto se anide correctamente.
10. Verificar z-index suficiente para que el dropdown aparezca sobre el grid/tabla de Pokémon.
11. Verificar que los estilos funcionan tanto en light como dark mode.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Al escribir ≥1 carácter aparece un dropdown con sugerencias |
| (visual) | Las sugerencias filtran por subcadena (case-insensitive) |
| (visual) | Navegar con ↑↓ y confirmar con Enter carga el Pokémon |
| (visual) | Clic en una sugerencia carga el Pokémon |
| (visual) | Escape cierra el dropdown |
| (visual) | Clic fuera cierra el dropdown |
| (visual) | Máximo 8 sugerencias visibles |
| (visual) | El dropdown se ve bien en dark mode |
| (visual) | El dropdown no se solapa con los filtros |
| (visual) | Enter sin sugerencia activa sigue haciendo búsqueda exacta como antes |

## Criterios de aceptacion

- [ ] Al escribir en el buscador principal de la Pokédex aparece un dropdown con sugerencias de nombres de Pokémon
- [ ] El filtrado es por subcadena, case-insensitive, máximo 8 sugerencias
- [ ] Se puede navegar con ↑↓ y confirmar con Enter o clic
- [ ] Al seleccionar una sugerencia se carga el detalle del Pokémon automáticamente
- [ ] Escape y clic fuera cierran el dropdown
- [ ] Compatible con dark mode
- [ ] No rompe el comportamiento existente de búsqueda con Enter/botón Buscar

## Notas

- El módulo `createAutocomplete` ya existe en `frontend/src/autocomplete.ts` y se usa en builds.ts. No requiere modificaciones.
- La tarea 0015 implementó autocomplete para Compare y Builds; esta tarea extiende esa misma funcionalidad al buscador principal.
- `ListPokemon(0, 2000)` cubre toda la Pokédex nacional (~1302 entradas).
