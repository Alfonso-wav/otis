# 0023 — Responsive Pokedex y Types pages

## Descripcion

Hacer que las paginas de Pokedex (grid + tabla) y Types se adapten completamente a todos los tamaños de ventana: cards mas pequeñas en mobile, tabla con scroll horizontal, grids que colapsen correctamente, y type cards que se reorganicen.

## Capas involucradas

- **Frontend (APP)**: `frontend/src/styles/_pokemon.scss`, `frontend/src/styles/_types.scss`, posibles ajustes en `frontend/src/pages/pokedex.ts`.
- **Core**: No se requieren cambios.
- **Shell**: No se requieren cambios.

## Contexto actual

- Pokemon grid usa `minmax(148px, 1fr)` — funciona razonablemente pero las cards pueden ser muy pequeñas en mobile.
- La tabla `.poke-table` tiene sticky header pero no tiene scroll horizontal en mobile.
- Las type cards usan `minmax(220px, 1fr)` sin ajustes responsive.
- La paginacion no se adapta bien a mobile.

## Plan de implementacion

### Paso 1 — Pokemon grid responsive

- Ajustar `minmax` para mobile: cards mas grandes en mobile (minmax ~120px) para que quepan 2-3 por fila.
- En desktop wide: permitir mas columnas.
- Asegurar que las imagenes de los Pokemon escalen proporcionalmente.
- Card hover effects: desactivar en touch devices.

- [ ] Completado

### Paso 2 — Pokemon tabla responsive

- Envolver la tabla en un contenedor con `overflow-x: auto` para scroll horizontal en mobile.
- Reducir font-size y padding de celdas en mobile.
- Mantener la primera columna (nombre/imagen) sticky en mobile si es posible.
- Ajustar los indicadores de sorting para que sean visibles en mobile.

- [ ] Completado

### Paso 3 — Paginacion responsive

- Reducir tamaño de botones de paginacion en mobile.
- Simplificar el texto de paginacion ("1/50" en vez de "Pagina 1 de 50").
- Centrar controles de paginacion.

- [ ] Completado

### Paso 4 — Types page responsive

- Type cards: ajustar minmax para mobile (minmax ~180px).
- Contenido expandible dentro de type cards: stack vertical en mobile.
- Asegurar que las badges de tipo sean legibles en todos los tamaños.

- [ ] Completado

## Criterios de aceptacion

- [ ] Pokemon grid muestra 2-3 cards por fila en mobile, 4+ en desktop.
- [ ] Tabla de Pokemon tiene scroll horizontal en mobile sin romper el layout.
- [ ] Paginacion usable en mobile con botones de tamaño tactil.
- [ ] Type cards se reorganizan correctamente en todos los breakpoints.
- [ ] No se pierde funcionalidad (sorting, filtros, toggle grid/tabla) en ningún tamaño.
- [ ] No se introducen cambios en Core ni Shell.

## Dependencias externas nuevas

Ninguna.

## Notas

- Depende de 0022 para los mixins de breakpoints.
- Considerar usar `@media (hover: hover)` para diferenciar touch vs mouse.
