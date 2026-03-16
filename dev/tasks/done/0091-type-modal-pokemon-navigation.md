# Navegar a detalle de Pokémon desde el modal de tipo en regiones

**ID**: 0091-type-modal-pokemon-navigation
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

En Explorar > Regiones, al hacer clic en un segmento del donut de distribución de tipos se abre un modal (`pokemon-type-modal.ts`) con la lista de Pokémon de ese tipo en esa región. Actualmente esos Pokémon no son clickables. Se necesita que al hacer clic en un Pokémon del modal, se navegue a su vista de detalle en la Pokédex (la misma vista que se muestra al buscar/clickar un Pokémon desde la pestaña Pokédex).

## Capas afectadas

- **APP (frontend)**: `pokemon-type-modal.ts` — añadir click handler en cada `.type-modal-pokemon` que navegue al detalle del Pokémon.
- **Frontend styles**: `_explore.scss` — añadir `cursor: pointer` y hover feedback en las cards de Pokémon del modal.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/components/pokemon-type-modal.ts` | modificar | Añadir click handler en cada card de Pokémon para navegar a su vista de detalle en la Pokédex |
| `frontend/src/styles/_explore.scss` | modificar | Añadir `cursor: pointer` en `.type-modal-pokemon` si no lo tiene ya |

## Plan de implementacion

1. En `pokemon-type-modal.ts`, después de insertar el overlay en el DOM, seleccionar todos los `.type-modal-pokemon` y añadir un event listener de click que:
   - Obtenga el nombre del Pokémon (del `alt` del `<img>` o añadir un `data-name` al div).
   - Cierre el modal (`closeTypeModal()`).
   - Navegue a la pestaña Pokédex haciendo click programático en `[data-tab="pokedex"]`.
   - Escriba el nombre en el input de búsqueda (`#search-input`) y haga click en el botón de búsqueda (`#search-btn`).
   - Este patrón ya está implementado en `frontend/src/pages/explore/types.ts` líneas 132-144.
2. Añadir atributo `data-name="${name}"` al div `.type-modal-pokemon` para facilitar la obtención del nombre.
3. En `_explore.scss`, asegurar que `.type-modal-pokemon` tiene `cursor: pointer` y un hover visual claro.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Al hacer clic en un Pokémon del modal de tipo, se cierra el modal y se navega a la vista de detalle en la Pokédex |
| (visual) | El cursor cambia a pointer al pasar sobre un Pokémon en el modal |
| (visual) | Funciona correctamente en dark mode |
| (visual) | Funciona correctamente en mobile (responsive) |

## Criterios de aceptacion

- [x] Al hacer clic en un Pokémon del modal de distribución de tipos, se cierra el modal
- [x] Se navega a la pestaña Pokédex y se muestra el detalle de ese Pokémon
- [x] Los Pokémon del modal muestran cursor pointer y feedback visual al hover
- [x] Compatible con dark mode
- [x] Compatible con mobile (responsive)

## Notas

- El patrón de navegación cross-tab ya existe en `types.ts` (líneas 132-144): click en tab pokedex → escribir nombre en search input → click en search btn.
- No se necesitan cambios en backend.
- Cambio pequeño y localizado, solo frontend.
