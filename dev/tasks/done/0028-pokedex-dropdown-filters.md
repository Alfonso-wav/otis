# Filtros multi-selección en lista desplegable (Pokédex)

**ID**: 0028-pokedex-dropdown-filters
**Estado**: done
**Fecha**: 2026-03-14

---

## Descripcion

Convertir los filtros de Generación y Tipo de la vista principal de la Pokédex de chips horizontales con scroll a listas desplegables (dropdown). Cada dropdown muestra un botón que al hacer clic despliega un panel con las opciones seleccionables (chips/checkboxes). Las selecciones activas se reflejan en el texto del botón (ej: "Tipo (2)") y el panel se cierra al hacer clic fuera. Esto mejora la UX al ocupar menos espacio horizontal y escalar mejor con muchas opciones.

## Capas afectadas

- **Core**: sin cambios
- **Shell**: sin cambios
- **APP**: frontend — HTML, TypeScript y SCSS

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Reemplazar `.filter-multi-group` por estructura de dropdown con botón trigger y panel desplegable |
| `frontend/src/pages/pokedex.ts` | modificar | Adaptar `populateFilters()` y `resetFilterUI()` para crear chips dentro del panel dropdown; añadir lógica de abrir/cerrar dropdown |
| `frontend/src/styles/_components.scss` | modificar | Añadir estilos para `.filter-dropdown`, `.filter-dropdown__trigger`, `.filter-dropdown__panel`; actualizar/eliminar estilos de `.filter-multi-group` y `.filter-multi-select` que ya no se usen |

## Plan de implementacion

1. Modificar `frontend/index.html`: reemplazar cada `.filter-multi-group` por un contenedor `.filter-dropdown` que contenga un `<button class="filter-dropdown__trigger">` y un `<div class="filter-dropdown__panel">` donde se inyectarán los chips.
2. Añadir estilos en `_components.scss` para el dropdown: posición relativa en el contenedor, panel con `position: absolute`, `display: none` por defecto, y clase `.open` para mostrarlo. Los chips dentro del panel se muestran en `flex-wrap: wrap`.
3. Modificar `populateFilters()` en `pokedex.ts` para inyectar los chips dentro de `.filter-dropdown__panel` en lugar de `.filter-multi-select`. Actualizar el texto del trigger cuando cambian las selecciones (ej: "Generación (3)").
4. Añadir lógica de toggle en el botón trigger y cierre al hacer clic fuera del dropdown (event listener en `document`).
5. Adaptar `resetFilterUI()` para cerrar los dropdowns y resetear el texto de los triggers.
6. Limpiar estilos huérfanos de `.filter-multi-group` y `.filter-multi-select` si ya no se usan en ningún otro lugar.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Abrir/cerrar dropdown de Generación y Tipo, seleccionar múltiples opciones, verificar filtrado correcto, verificar cierre al clic externo, verificar reset |

## Criterios de aceptacion

- [x] Cada filtro (Generación, Tipo) se muestra como un botón dropdown en la barra de filtros
- [x] Al hacer clic en el botón se despliega un panel con los chips seleccionables
- [x] Se pueden seleccionar múltiples chips y la lógica de filtrado existente sigue funcionando
- [x] El texto del botón refleja cuántas opciones están seleccionadas (ej: "Tipo (2)")
- [x] El panel se cierra al hacer clic fuera del dropdown
- [x] El botón "Todos" (reset) cierra los dropdowns y limpia las selecciones
- [x] Funciona correctamente en mobile (responsive)
- [x] Solo un dropdown puede estar abierto a la vez

## Notas

- La lógica de filtrado (`applyFilters`, `hasFilter`, `loadFiltered`) no necesita cambios, solo cambia la forma de presentar las opciones.
- Reutilizar la clase `.filter-chip` existente dentro del panel dropdown para mantener consistencia visual.
- Considerar que el panel no se salga del viewport en pantallas pequeñas.
