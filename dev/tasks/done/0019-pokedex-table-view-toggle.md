# 0019 — Toggle vista tabla / tarjetas en la Pokédex

## Descripción

Añadir un botón en la vista principal de la Pokédex que permita alternar entre la vista de tarjetas actual y una nueva vista de tabla con nombre, tipo(s) y todas las stats base (HP, Ataque, Defensa, Sp. Atk, Sp. Def, Velocidad). La transición entre vistas debe ser animada y visualmente llamativa usando GSAP.

## Capas involucradas

- **Frontend (APP)**: `frontend/src/pages/pokedex.ts`, `frontend/index.html`, `frontend/src/styles/_pokemon.scss`, `frontend/src/animations/transitions.ts`
- **Core / Shell**: No se requieren cambios en backend. Se reutiliza la función `GetPokemon` existente para cargar datos al cambiar a vista tabla.

## Contexto actual

- `renderGrid()` genera tarjetas con sprite, número y nombre desde `PokemonListItem[]` (solo `Name` y `URL`).
- La tabla necesita tipos y stats → hay que llamar a `GetPokemon(name)` para los 20 Pokémon visibles al activar la vista tabla.
- Ya existe GSAP y la función `staggerCards` en `animations/transitions.ts`.
- El grid se define en `_pokemon.scss` bajo `#pokemon-grid`.

## Plan de implementación

### Paso 1 — Botón de toggle en el HTML

En `frontend/index.html`, dentro de `#filter-bar`, añadir un botón con id `view-toggle-btn` y un icono de tabla/grid (puede ser texto o SVG simple). Colocarlo al final de la barra de filtros.

```html
<button id="view-toggle-btn" class="filter-pill" title="Cambiar vista">⊞ Tabla</button>
```

### Paso 2 — Estado de vista en `pokedex.ts`

- Añadir variable de estado `let viewMode: 'grid' | 'table' = 'grid'`.
- Añadir ref DOM `let viewToggleBtn: HTMLButtonElement`.
- Cache para datos ya cargados: `const pokemonDataCache = new Map<string, Pokemon>()` (reutilizar la info descargada).

### Paso 3 — Función `renderTable(items: PokemonListItem[])`

- Mostrar un spinner mientras se cargan los datos (puede ser la misma línea de "Cargando…").
- Llamar en paralelo a `GetPokemon(name)` para cada item de la página actual, usando el cache para no repetir llamadas.
- Generar un `<table>` con columnas: `#`, Sprite (mini 40px), Nombre, Tipo(s), HP, Atk, Def, SpA, SpD, Vel, Total.
- Cada fila debe ser clickable para abrir el detalle (igual que las tarjetas).
- Llamar a la animación de entrada de filas al terminar.

### Paso 4 — Animación de transición (GSAP)

En `animations/transitions.ts`, añadir dos funciones exportadas:

**`morphToTable(container: HTMLElement)`**
- `gsap.to(container.children, { opacity: 0, scale: 0.85, y: -10, stagger: 0.015, duration: 0.18, ease: 'power2.in' })` para desvanecer las tarjetas.
- Al completar, reemplazar el innerHTML por la tabla y animar las filas de entrada con `fromTo` (slide desde la izquierda: `x: -20 → 0`, con stagger).

**`morphToGrid(container: HTMLElement)`**
- Fade + slide de salida de filas de tabla, luego entrada de tarjetas con `staggerCards`.

### Paso 5 — Estilos en `_pokemon.scss`

Añadir estilos para:
- `.poke-table` — tabla full-width, `border-collapse: collapse`, `font-family: $font-body`.
- `.poke-table th` — header fijo, fondo sutil, `font-size: 0.78rem`, sticky top.
- `.poke-table td` — padding compacto, border-bottom suave.
- `.poke-table tr:hover` — highlight row con `background` suave y cursor pointer.
- `.poke-table__sprite` — 40×40px, `image-rendering: pixelated`.
- `.poke-table__name` — capitalize, `font-family: $font-display`.
- `.stat-cell` — `text-align: right`, `font-weight: 700`.
- `.stat-total` — color destacado (e.g. `$pokedex-red`).
- Clases de tipo reutilizando `.type-badge` ya existente.

### Paso 6 — Lógica de toggle en `initPokedex()`

```ts
viewToggleBtn.addEventListener('click', async () => {
  viewMode = viewMode === 'grid' ? 'table' : 'grid';
  viewToggleBtn.textContent = viewMode === 'grid' ? '⊞ Tabla' : '⊟ Tarjetas';
  await rerender(); // llama renderGrid o renderTable según viewMode
});
```

Refactorizar `renderGrid` y `renderTable` para que `loadList`, `loadFiltered`, `prevPage`, `nextPage` respeten el `viewMode` actual.

## Criterios de aceptación

- [ ] Botón visible en la barra de filtros que alterna el texto/icono entre "Tabla" y "Tarjetas".
- [ ] Vista tabla muestra: sprite mini, nombre, tipo(s) como badges, HP, Atk, Def, SpA, SpD, Vel, Total.
- [ ] Clicar una fila de la tabla abre el detalle igual que las tarjetas.
- [ ] La transición tarjetas→tabla usa GSAP con desvanecer + slide de entrada de filas.
- [ ] La transición tabla→tarjetas usa GSAP con desvanecer + `staggerCards`.
- [ ] La paginación funciona correctamente en ambas vistas.
- [ ] Los filtros (generación, tipo, legendario, mítico) funcionan en ambas vistas.
- [ ] Los datos de Pokémon se cachean para no repetir llamadas al backend.
- [ ] No se introducen cambios en Core ni Shell.

## Dependencias externas nuevas

Ninguna. Se reutilizan GSAP (ya instalado) y las llamadas `GetPokemon` existentes.

## Notas

- El spinner de carga al cambiar a tabla puede ser el mismo `<p class="loading">Cargando datos…</p>` que ya existe.
- El `total` de stats es la suma de los 6 valores base.
- Los type-badges ya tienen estilos CSS en `_types.scss` — reutilizarlos directamente.
