# Fix overflow en Builds mobile (equipos y búsqueda)

**ID**: 0064-fix-builds-mobile-overflow
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Dos problemas de overflow en la pestaña Builds cuando se visualiza en móvil en orientación vertical:

1. **Equipos guardados — botones cortados**: Al desplegar un equipo guardado, la fila de cada miembro (`.team-member-row`) usa `display: flex` sin `flex-wrap` ni ajuste de overflow. Los botones "Atk", "Def" y el botón de eliminar se cortan por el borde derecho de la pantalla porque la fila no cabe en el ancho disponible (~360px).

2. **Búsqueda individual de atacante/defensor — recuadro cortado**: El contenedor `.build-search-row` dentro de `.build-col` se corta por la derecha en pantallas estrechas. El input y el botón no se adaptan correctamente al ancho disponible en la vista vertical del móvil.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Solo frontend (SCSS).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_builds.scss` | modificar | Añadir reglas mobile para `.team-member-row` (wrap o layout vertical) y asegurar que `.build-search-row` / `.build-col` no desborden |

## Plan de implementacion

1. **Problema 1 — `.team-member-row` en mobile**:
   - En el bloque `@include mobile`, hacer que `.team-member-row` use `flex-wrap: wrap` o cambie a layout de dos líneas (nombre+sprite arriba, botones abajo).
   - Reducir `min-width` de `.team-member-name` y `.team-member-detail` en mobile.
   - Asegurar que los botones `.team-import-btn` y `.team-member-delete-btn` tengan touch target mínimo de 44px.

2. **Problema 2 — `.build-search-row` / `.build-col` en mobile**:
   - Verificar que `.build-col` tiene `overflow: hidden` o `max-width: 100%` y no excede el viewport.
   - Asegurar que `.build-search-input` con `flex: 1` y `.build-search-btn` caben sin overflow. Si no, hacer que el botón baje a segunda línea con `flex-wrap: wrap` o reducir padding.
   - Comprobar que no hay `min-width` rígido que cause desbordamiento.

3. Verificar en viewport 360x800 que ambos problemas quedan resueltos sin afectar la vista desktop.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual | En móvil vertical (360px), los botones Atk/Def/eliminar de los miembros del equipo son visibles y pulsables |
| Test manual | En móvil vertical (360px), el campo de búsqueda y botón de atacante/defensor no se cortan por la derecha |
| Test manual | La vista desktop no se ve afectada |

## Criterios de aceptacion

- [x] Los botones Atk, Def y eliminar son completamente visibles en móvil vertical al desplegar equipos
- [x] Los botones tienen touch target >= 44px en mobile
- [x] El recuadro de búsqueda individual (atacante/defensor) no se corta por la derecha en móvil vertical
- [x] No hay scroll horizontal no deseado en la pestaña Builds en mobile
- [x] La vista desktop no se ve afectada

## Notas

- Los estilos actuales de `.team-member-row` (línea ~1221 de `_builds.scss`) usan flex sin wrap y con `min-width` fijos en nombre (90px) y detalle (80px), lo que suma demasiado para 360px.
- El bloque `@include mobile` actual (línea ~1503) solo ajusta `.team-member-edit-moves`, no las filas de miembros del equipo ni la búsqueda.
- Relacionada con tarea 0061 (optimización mobile general).
