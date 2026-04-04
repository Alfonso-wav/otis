# Tabla Pokédex: campo "Límite de filas" y paginación con flechas

**ID**: 0102-table-row-limit-label-and-pagination
**Estado**: done
**Fecha**: 2026-04-04

---

## Descripcion

Fix del cambio anterior (0101). Se necesitan dos ajustes en la vista de tabla de la Pokédex:

1. **Campo "Límite de filas"**: El control actual de ±50 debe mostrarse como un campo con label "Límite de filas" (i18n), con valor por defecto 50. El usuario puede ajustar cuántas filas se muestran por página.

2. **Paginación con flechas**: Independientemente del límite de filas configurado, se necesitan botones/flechas de navegación para ir a la página siguiente y anterior. Actualmente no existen — al cambiar el límite se resetea al offset 0 y no hay forma de avanzar a la siguiente página de resultados.

## Capas afectadas

- **Core**: No aplica
- **Shell**: No aplica
- **APP**: Frontend — `pokedex.ts`, `index.html`, `_pokemon.scss`, archivos de locales

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Añadir label "Límite de filas" al control existente `#row-limit-control`; añadir botones de paginación (anterior/siguiente) con indicador de página actual |
| `frontend/src/pages/pokedex.ts` | modificar | Implementar lógica de paginación: navegar entre páginas usando `offset`, actualizar `getCurrentPageItems()` para paginar correctamente, deshabilitar botones prev/next en los extremos |
| `frontend/src/styles/_pokemon.scss` | modificar | Estilos para los botones de paginación y el label del campo de límite |
| `frontend/src/locales/es.json` | modificar | Añadir claves i18n: `pokedex.rowLimit` ("Límite de filas"), `pokedex.page` ("Página"), etc. |
| `frontend/src/locales/en.json` | modificar | Añadir claves i18n: `pokedex.rowLimit` ("Row limit"), `pokedex.page` ("Page"), etc. |

## Plan de implementacion

1. En `index.html`, añadir un `<label>` con texto i18n "Límite de filas" al `#row-limit-control` existente
2. En `index.html`, añadir un bloque de paginación debajo o junto al control de límite: botón `←` (anterior), indicador `Página X / Y`, botón `→` (siguiente)
3. En `pokedex.ts`, implementar la navegación por páginas:
   - Calcular `totalPages = Math.ceil(totalItems / rowLimit)`
   - Botón siguiente: `offset += rowLimit`, re-render con `getCurrentPageItems()`
   - Botón anterior: `offset -= rowLimit`, re-render
   - Deshabilitar `←` en página 1, `→` en última página
   - Al cambiar `rowLimit`, resetear `offset = 0` (ya existe)
4. Actualizar `getCurrentPageItems()` para usar `offset` y `rowLimit` correctamente con todas las fuentes de datos (sortedFullList, filteredList, lista base)
5. Actualizar `updateRowLimitControl()` para también actualizar el indicador de página
6. Añadir claves i18n en `es.json` y `en.json`
7. Aplicar estilos coherentes con el diseño actual

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | El label "Límite de filas" aparece junto al control |
| Manual | Botones ← → navegan entre páginas correctamente |
| Manual | El indicador muestra "Página X / Y" actualizado |
| Manual | ← deshabilitado en página 1, → en última página |
| Manual | Al cambiar el límite se vuelve a página 1 |
| Manual | Funciona con filtros activos y con ordenación |

## Criterios de aceptacion

- [ ] El control muestra el label "Límite de filas" (traducido según idioma)
- [ ] Valor por defecto del límite es 50
- [ ] Existen botones de página anterior (←) y siguiente (→)
- [ ] Se muestra indicador de página actual / total
- [ ] ← se deshabilita en la primera página
- [ ] → se deshabilita en la última página
- [ ] Al cambiar el límite de filas, se resetea a la página 1
- [ ] La paginación funciona correctamente con filtros y ordenación activos
- [ ] Los textos están internacionalizados (es/en)
- [ ] No hay regresiones en la vista de grid ni en móvil

## Notas

- El `offset` ya existe en `pokedex.ts` pero solo se usa en grid view (infinite scroll). Hay que reutilizarlo para la paginación de tabla.
- `getCurrentPageItems()` (línea ~1096) ya hace slice con `offset` y `rowLimit` pero solo para `sortedFullList`. Hay que cubrir también `filteredList` y la lista base.
- Los botones ±50 actuales (líneas 63-67 de index.html) se mantienen pero con el label añadido.
