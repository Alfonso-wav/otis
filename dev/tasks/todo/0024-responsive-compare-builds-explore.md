# 0024 — Responsive Compare, Builds y Explore pages

## Descripcion

Hacer que las paginas de Compare (comparador de Pokemon), Builds (simulador de batalla y daño) y Explore (regiones, movimientos, habilidades) se adapten completamente a todos los tamaños de ventana.

## Capas involucradas

- **Frontend (APP)**: `frontend/src/styles/_compare.scss`, `frontend/src/styles/_builds.scss`, `frontend/src/styles/_explore.scss`.
- **Core**: No se requieren cambios.
- **Shell**: No se requieren cambios.

## Contexto actual

- Compare: layout side-by-side que colapsa a 1 columna en 480px, pero no hay breakpoint intermedio para tablets. La tabla de stats y los charts no se ajustan bien.
- Builds: layout 2 columnas (attacker/defender) colapsa a 768px. Los move slots pasan de 4 a 2 columnas a 900px. La seccion de batalla (HP bars, log, botones) no tiene ajustes responsive robustos.
- Explore: region cards y grids de moves/abilities usan auto-fill pero no tienen ajustes especificos para mobile.

## Plan de implementacion

### Paso 1 — Compare page responsive

- Inputs side-by-side: añadir breakpoint tablet (768px) para stack parcial.
- Stats comparison table: scroll horizontal o layout alternativo en mobile.
- Radar chart (ECharts): reducir tamaño del canvas en mobile, ajustar opciones del chart.
- Pokemon cards en compare: reducir tamaño de imagen y stats en mobile.
- "VS" divider: ocultar o reducir en mobile.

- [ ] Completado

### Paso 2 — Builds page responsive

- Layout attacker/defender: ya colapsa a 768px, verificar que funcione bien.
- Move slots grid: ajustar a 2 columnas en tablet, 1 en mobile.
- Damage table: scroll horizontal en mobile, reducir font-size.
- EV calculator: stack inputs en mobile, reducir el grid de stats de 3 a 2 o 1 columna.

- [ ] Completado

### Paso 3 — Battle section responsive

- HP bars: asegurar que flex-wrap funcione correctamente, labels legibles.
- Move buttons: grid responsive, tamaño tactil adecuado.
- Battle log: altura reducida en mobile, scroll funcional.
- Control buttons (auto-simulate, reset): stack vertical o wrap en mobile.

- [ ] Completado

### Paso 4 — Explore pages responsive

- Region cards: reducir padding y font-size en mobile, location grid a 1-2 columnas.
- Moves grid: ajustar minmax para mobile, move cards mas compactas.
- Abilities grid: similar a moves, cards compactas en mobile.
- Mini-tabs de explore: scroll horizontal si no caben.

- [ ] Completado

## Criterios de aceptacion

- [ ] Compare page usable en mobile: stats legibles, chart visible, inputs accesibles.
- [ ] Builds page: todos los controles accesibles en mobile, damage table legible.
- [ ] Battle section: botones de tamaño tactil, log scrollable, HP bars visibles.
- [ ] Explore: todas las sub-paginas (regions, moves, abilities) usables en mobile.
- [ ] Transiciones entre tamaños de ventana son fluidas (sin saltos bruscos).
- [ ] No se introducen cambios en Core ni Shell.

## Dependencias externas nuevas

Ninguna.

## Notas

- Depende de 0022 para los mixins de breakpoints.
- Para los charts de ECharts, puede ser necesario llamar `chart.resize()` en el evento `window.resize` si no se hace ya.
- Considerar reducir la cantidad de datos visibles en mobile (ej. menos columnas en tablas) sin perder funcionalidad.
