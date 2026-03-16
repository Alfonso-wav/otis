# Navegar a detalle de Pokémon desde el modal de encuentros de localización

**ID**: 0092-region-encounter-modal-pokemon-navigation
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

En Explorar > Regiones, al hacer clic en una localización se abre un modal (`location-encounter-modal.ts`) que muestra los Pokémon que se pueden encontrar en esa zona con su porcentaje de encuentro. Actualmente esos Pokémon no son clickables. Se necesita que al hacer clic en un Pokémon del modal, se navegue a su vista de detalle en la Pokédex (la misma vista que se muestra al buscar/clickar un Pokémon desde la pestaña Pokédex).

## Capas afectadas

- **APP (frontend)**: `location-encounter-modal.ts` — añadir click handler en cada `.type-modal-pokemon` que navegue al detalle del Pokémon.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/components/location-encounter-modal.ts` | modificar | Añadir click handler en cada card de Pokémon para navegar a su vista de detalle en la Pokédex |

## Plan de implementacion

1. En `location-encounter-modal.ts`, después de insertar el overlay en el DOM, seleccionar todos los `.type-modal-pokemon` y añadir un event listener de click que:
   - Obtenga el nombre del Pokémon (del atributo `data-name` del div).
   - Cierre el modal (eliminar overlay del DOM).
   - Navegue a la pestaña Pokédex haciendo click programático en `[data-tab="pokedex"]`.
   - Escriba el nombre en el input de búsqueda (`#search-input`) y haga click en el botón de búsqueda (`#search-btn`).
   - Este patrón ya está implementado en la tarea 0091 en `pokemon-type-modal.ts`.
2. Asegurar que cada div `.type-modal-pokemon` tiene un atributo `data-name` con el nombre del Pokémon para facilitar la obtención del nombre.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Al hacer clic en un Pokémon del modal de encuentros, se cierra el modal y se navega a la vista de detalle en la Pokédex |
| (visual) | El cursor cambia a pointer al pasar sobre un Pokémon en el modal |
| (visual) | Funciona correctamente en dark mode |
| (visual) | Funciona correctamente en mobile (responsive) |

## Criterios de aceptacion

- [x] Al hacer clic en un Pokémon del modal de encuentros de localización, se cierra el modal
- [x] Se navega a la pestaña Pokédex y se muestra el detalle de ese Pokémon
- [x] Los Pokémon del modal muestran cursor pointer y feedback visual al hover
- [x] Compatible con dark mode
- [x] Compatible con mobile (responsive)

## Notas

- El patrón de navegación cross-tab ya existe en `pokemon-type-modal.ts` (tarea 0091): cerrar modal → click en tab pokedex → escribir nombre en search input → click en search btn.
- No se necesitan cambios en backend ni en estilos (`.type-modal-pokemon` ya tiene `cursor: pointer` desde la tarea 0091).
- Cambio pequeño y localizado, solo frontend.
