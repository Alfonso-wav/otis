# Task 0118 — Vista tabla: overlay de carga al entrar + controles encima

## Estado: done

## Goal

Dos mejoras en la vista de tabla de la Pokédex:

1. **Mr. Mime no aparece al cambiar de tarjetas a tabla**: durante el tiempo de carga al entrar en la vista de tabla desde la vista de tarjetas, no se muestra el overlay de carga (Mr. Mime). Hay que mostrar el overlay mientras se cargan y renderizan los datos de la tabla.

2. **Controles de paginación encima de la tabla**: actualmente el `#row-limit-control` (selector de límite de filas + paginador) aparece debajo de la tabla. El usuario quiere que esté encima.

## Contexto técnico

### 1. Mr. Mime ausente durante la transición tarjetas → tabla

- El listener del botón `#view-toggle-btn` está en `frontend/src/pages/pokedex.ts` (~línea 1462).
- Al cambiar de grid a table (líneas 1469-1482), se llama a `morphToTable(grid, callback)` donde el callback carga y renderiza la tabla.
- **Problema**: no se llama a `showSortingOverlay()` antes del callback ni se oculta con `hideSortingOverlay()` después. Durante la carga de datos (que puede ser lenta), la UI queda en blanco sin indicación visual.
- `showSortingOverlay` y `hideSortingOverlay` ya están importados (~línea 18). El overlay sí se usa en la ordenación de columnas (~línea 501) y en el filtrado (~línea 1277).
- **Solución**: llamar a `showSortingOverlay()` justo antes de `morphToTable(...)` y a `hideSortingOverlay()` al final del callback, dentro del bloque `if (oldMode === "grid" && viewMode === "table")`.

### 2. Mover `#row-limit-control` encima de la tabla

- HTML: `frontend/index.html`.
- Estructura actual (líneas 59-72):
  ```html
  <div id="pokemon-grid"></div>
  <div id="scroll-sentinel" ...></div>
  <div id="row-limit-control" class="row-limit-control hidden">...</div>
  ```
- **Solución**: mover el `<div id="row-limit-control">` **antes** de `<div id="pokemon-grid">`, quedando:
  ```html
  <div id="row-limit-control" class="row-limit-control hidden">...</div>
  <div id="pokemon-grid"></div>
  <div id="scroll-sentinel" ...></div>
  ```
- No es necesario cambiar ningún JS/TS ya que los elementos se referencian por ID y no por posición DOM.

## Acceptance criteria

### Mr. Mime al cambiar a tabla
- [ ] Al pulsar el botón de cambio de vista (tarjetas → tabla), aparece el overlay de Mr. Mime mientras se cargan los datos.
- [ ] El overlay desaparece cuando la tabla ya está renderizada.
- [ ] No se rompe la animación de `morphToTable` ni el stagger de filas.

### Controles encima de la tabla
- [ ] El selector de límite de filas y los botones de paginación aparecen visualmente encima de la tabla cuando el modo tabla está activo.
- [ ] Siguen estando ocultos en modo tarjetas (comportamiento `hidden` sin cambios).
- [ ] El orden de los controles y su funcionamiento no cambia.

## Archivos afectados

- `frontend/index.html` — mover `#row-limit-control` antes de `#pokemon-grid`
- `frontend/src/pages/pokedex.ts` — añadir `showSortingOverlay()` / `hideSortingOverlay()` en el listener de `#view-toggle-btn`
