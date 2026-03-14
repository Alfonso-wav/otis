# 0027 — Filtros multi-seleccion y toggles independientes en Pokedex

## Estado

todo

## Descripcion

Mejorar los filtros de la vista principal de la Pokedex:
1. **Tipo y Generacion**: cambiar de `<select>` simple a seleccion multiple (poder elegir varios tipos y/o varias generaciones a la vez).
2. **Mitico y Legendario**: convertirlos en toggles on/off que afecten a toda la tabla sin requerir que haya una generacion o tipo seleccionado. Actualmente muestran un mensaje pidiendo seleccionar generacion/tipo primero.

## Contexto

- `frontend/src/pages/pokedex.ts` — toda la logica de filtros, estado `FilterState`, funcion `loadFiltered()`, `filterByLegendary()`, event listeners.
- `frontend/index.html` — markup del `#filter-bar` con los `<select>` y botones pill.
- `frontend/src/styles/_components.scss` — estilos de filtros.
- El filtro legendario/mitico actualmente requiere gen o tipo porque necesita una lista base sobre la que iterar y llamar a `GetPokemonSpecies()` uno por uno. Sin base, habria que recorrer los ~1025 pokemon.
- Existe un cache `legendaryCache` que evita re-fetches.

## Plan de implementacion

### Paso 1 — Reemplazar `<select>` por componente multi-select para Tipo y Generacion
- [ ] En `index.html`, reemplazar los `<select>` de generacion y tipo por contenedores de pills/chips clickables que permitan seleccion multiple.
- [ ] En `FilterState`, cambiar `generation: string | null` a `generations: string[]` y `type: string | null` a `types: string[]`.
- [ ] Actualizar `populateFilters()` para generar pills en lugar de `<option>`.
- [ ] Actualizar `applyFilters()` y `hasFilter()` para trabajar con arrays.
- [ ] Estilos: pills seleccionados con clase `.active`, scroll horizontal si hay muchos.

### Paso 2 — Actualizar logica de `loadFiltered()` para multi-seleccion
- [ ] Cuando hay multiples generaciones: hacer fetch de todas y unir las listas (union).
- [ ] Cuando hay multiples tipos: hacer fetch de todos y unir las listas (union).
- [ ] Cuando hay generaciones + tipos: interseccion (pokemon que este en alguna de las generaciones seleccionadas Y en alguno de los tipos seleccionados).
- [ ] Mantener deduplicacion por nombre.

### Paso 3 — Legendario/Mitico independientes (sin requerir gen/tipo)
- [ ] Eliminar la restriccion en `loadFiltered()` que muestra mensaje "Selecciona una generacion o tipo".
- [ ] Cuando solo hay filtro legendario/mitico sin gen/tipo: cargar la lista completa de pokemon via `ListPokemon` con offset incremental y aplicar `filterByLegendary()` en batches, mostrando resultados progresivamente.
- [ ] Aprovechar el `legendaryCache` existente para no re-fetchear species ya conocidas.

### Paso 4 — Actualizar reset y UI
- [ ] Actualizar `resetFilterUI()` para limpiar pills activos en lugar de resetear selects.
- [ ] Actualizar el boton de reset para limpiar todos los filtros multi-select.
- [ ] Asegurar que el sorting se resetea correctamente al cambiar filtros.

### Paso 5 — Testing visual y edge cases
- [ ] Verificar combinaciones: multi-gen + multi-tipo + legendario + mitico.
- [ ] Verificar que la paginacion funciona con las nuevas listas filtradas.
- [ ] Verificar que grid y table view funcionan con los nuevos filtros.
- [ ] Verificar que el boton reset limpia todo correctamente.

## Capas afectadas

- **APP/Frontend**: `pokedex.ts`, `index.html`, `_components.scss`
- **Core**: sin cambios
- **Shell**: sin cambios

## Notas

- La carga de legendarios/miticos sin filtro base sera mas lenta (~1025 species a verificar). Mostrar feedback de progreso al usuario.
- Considerar un limite de carga progresiva (cargar en batches de 50-100 y mostrar resultados parciales).
