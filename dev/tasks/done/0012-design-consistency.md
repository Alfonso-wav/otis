# 0012 — Design Consistency: Apply Pokédex Style to All Tabs

## Descripción
Aplicar el sistema visual de la Pokédex (filtro-bar, poke-cards, header rojo, tipografías, badges de tipo, animaciones stagger) al resto de las pestañas: **Types** y **Explore** (Regions, Moves, Abilities). El objetivo es que toda la app tenga una identidad visual coherente.

## Estado
- [x] Done

## Contexto
- La Pokédex (`_pokemon.scss`, `pages/pokedex.ts`) tiene un diseño completo: header con fondo rojo, filtros pill, tarjetas con hover, badges de tipo con colores, chart de stats, animaciones GSAP.
- **Types** (`pages/types.ts`, `_types.scss`): lista de tipos como acordeones con Pokémon dentro.
- **Explore** (`pages/explore.ts`, `_explore.scss`): subsecciones Regions, Moves, Abilities con listas simples.
- Variables de diseño centralizadas en `_variables.scss`.

## Capas afectadas
- **Frontend únicamente**: SCSS y TypeScript en `frontend/src/`.
- No requiere cambios en Go.

## Plan de implementación

### Paso 1 — Auditar estilos existentes
- Leer `_variables.scss`, `_types.scss`, `_explore.scss` completos.
- Identificar qué componentes de `_pokemon.scss` son reutilizables (clases utilitarias vs. específicas).

### Paso 2 — Extraer componentes compartidos
- Crear `_components.scss` con clases comunes:
  - `.section-header` (header con fondo rojo y título)
  - `.filter-bar` (ya existe en `_pokemon.scss`, moverla a `_components`)
  - `.item-card` (tarjeta genérica similar a `.poke-card`)
  - `.type-badge` (ya existe, asegurarse que se usa en todas las vistas)
  - `.loading`, `.error-text` (estados de carga)
- Importar `_components.scss` en `main.scss` y eliminar duplicados.

### Paso 3 — Rediseñar Types view
- Reemplazar el acordeón simple por un layout de tarjetas de tipo (grid de type-badges grandes clicables).
- Al hacer click en un tipo, mostrar un panel con sus Pokémon en formato grid con sprites (igual que la Pokédex).
- Añadir header `.section-header` con título "Tipos".
- Aplicar animaciones `staggerCards` al revelar el grid de Pokémon.

### Paso 4 — Rediseñar Explore view
- **Regions**: tarjetas de región con nombre y número de locations. Al hacer click, lista de locations.
- **Moves**: tabla/lista con move-cards que muestren tipo (badge), categoría (físico/especial/estado), potencia y precisión. Igual que en el detalle de move actual pero en grid.
- **Abilities**: cards con nombre y descripción corta. Al click, detalle completo.
- Añadir `.section-header` en cada subsección.

### Paso 5 — HTML templates
- Actualizar el HTML en `index.html` si es necesario para los nuevos contenedores.
- Asegurar que cada sección tiene su estructura DOM esperada por los TS pages.

### Paso 6 — Review visual
- Verificar coherencia de colores, espaciados y tipografías entre las 3 pestañas.
- Asegurar que `staggerCards` y `showView` se usan consistentemente.

## Criterios de éxito
- Las 3 pestañas usan el mismo header rojo, tipografías, tarjetas y badges.
- No hay duplicación de CSS entre `_pokemon.scss`, `_types.scss` y `_explore.scss`.
- Las animaciones de transición son consistentes.
- No hay regresiones en la Pokédex.
