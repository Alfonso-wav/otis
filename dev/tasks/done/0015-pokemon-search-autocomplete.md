# 0015 — Autocompletado en buscadores de Pokémon

## Estado

done

## Descripción

Añadir funcionalidad de autocompletado en los campos de búsqueda de Pokémon del tab **Comparador** (`compare.ts`) y del tab **Builds** (`builds.ts`). Al escribir en los inputs, debe aparecer un dropdown con sugerencias filtradas de nombres de Pokémon.

## Contexto

- Los inputs afectados son:
  - Comparador: `#compare-input-a` y `#compare-input-b`
  - Builds: `#atk-input` y `#def-input`
- El backend ya expone `ListPokemon(offset, limit)` → `PokemonListResponse { Results: [{Name, URL}] }`
- Con `ListPokemon(0, 2000)` se obtienen todos los nombres de Pokémon de una sola llamada
- La lista se puede cargar una vez al inicializar el tab y reutilizarla para todos los inputs

## Capas involucradas

- **Frontend (APP)**: `compare.ts`, `builds.ts` — añadir lógica de autocompletado
- **Backend (APP)**: `bindings.go` — no requiere cambios; usar `ListPokemon` existente
- **Core/Shell**: sin cambios

## Plan de implementación

### Paso 1 — Crear módulo compartido `autocomplete.ts`
- Crear `frontend/src/autocomplete.ts` con una función pura `createAutocomplete(input, names, onSelect)`
- La función adjunta listeners al input y renderiza un `<ul class="autocomplete-dropdown">` como sibling
- Filtra nombres que contengan el texto escrito (case-insensitive), limitar a 8 sugerencias
- Navegar con teclado (↑↓ para seleccionar, Enter para confirmar, Escape para cerrar)
- Al hacer clic o Enter sobre una sugerencia, llama `onSelect(name)` y cierra el dropdown

### Paso 2 — Estilos en SCSS
- Añadir estilos para `.autocomplete-dropdown` en el archivo de estilos global o en un archivo dedicado
- Dropdown absoluto bajo el input, con fondo oscuro, bordes redondeados y hover highlight
- Consistente con el diseño existente (variables CSS del proyecto)

### Paso 3 — Cargar lista de nombres al inicializar tabs
- En `initCompare()`: llamar `ListPokemon(0, 2000)` al inicializar, guardar nombres en array local
- En `initBuilds()`: aprovechar el tick de inicialización del tab (ya hay `initialized` guard) para cargar nombres junto a `GetNatures()`

### Paso 4 — Integrar `createAutocomplete` en ambos tabs
- Comparador: llamar `createAutocomplete` en `bindInputs()` para los dos inputs
- Builds: llamar `createAutocomplete` en `bindEvents()` para `#atk-input` y `#def-input`
- El `onSelect` debe actualizar el valor del input y disparar la búsqueda automáticamente

## Criterios de aceptación

- [ ] Al escribir ≥1 carácter en cualquier buscador de Pokémon aparece un dropdown con sugerencias
- [ ] El dropdown filtra por prefijo o subcadena (case-insensitive)
- [ ] Se puede navegar con ↑↓ y confirmar con Enter
- [ ] Al hacer clic en una sugerencia el Pokémon se carga automáticamente
- [ ] Escape cierra el dropdown
- [ ] El dropdown desaparece al hacer clic fuera
- [ ] Máximo 8 sugerencias visibles
- [ ] Los estilos son coherentes con el resto del diseño

## Notas

- No se necesitan cambios en Core ni en Shell
- La lista de nombres se carga una vez por sesión de tab (no persiste entre reinicios de la app)
- `ListPokemon(0, 2000)` cubre la Pokédex nacional completa (~1302 entradas)
