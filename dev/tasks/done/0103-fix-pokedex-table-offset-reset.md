# Fix offset no se resetea al volver a vista tabla en Pokedex

**ID**: 0103-fix-pokedex-table-offset-reset
**Estado**: done
**Fecha**: 2026-04-04

---

## Descripcion

Bug: al cambiar de vista tabla a tarjetas y volver a tabla, la tabla empieza por el pokemon N (donde N es el row limit) en lugar de empezar desde el #1. Esto ocurre porque `getCurrentPageItems()` devuelve `lastRenderedItems` (datos del ultimo batch del infinite scroll) en vez de obtener datos frescos desde offset 0.

**Flujo del bug:**
1. Usuario en grid view hace scroll infinito (offset sube a 50, 100, etc.)
2. `lastRenderedItems` se actualiza con el ultimo batch cargado (ej: pokemon #51-100)
3. Usuario cambia a tabla: `offset` se resetea a 0 correctamente
4. `getCurrentPageItems()` es llamado, pero como no hay filtro ni sort activo, retorna `lastRenderedItems` (stale data del scroll)
5. La tabla muestra pokemon desde #50 en vez de #1

**Causa raiz:** En `getCurrentPageItems()` (linea ~1145 de `pokedex.ts`), el fallback `return lastRenderedItems` no respeta el offset recien reseteado.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend - logica de vista Pokedex

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Corregir `getCurrentPageItems()` o el handler de toggle de vista para que al cambiar a tabla se obtengan los items correctos desde offset 0, no el cache stale de `lastRenderedItems` |

## Plan de implementacion

1. En el handler de cambio de vista (grid -> table, ~linea 1188), tras resetear `offset = 0`, asegurar que `lastRenderedItems` tambien se resetee o que se obtengan datos frescos en lugar de usar el fallback stale.
2. Opcion A: Resetear `lastRenderedItems` antes de llamar a `getCurrentPageItems()` forzando un fetch fresco.
3. Opcion B: Modificar `getCurrentPageItems()` para que en el caso sin filtro ni sort, haga slice desde offset sobre la lista completa de pokemon cargados en vez de devolver `lastRenderedItems`.
4. Verificar que el fix no rompe: primera carga en tabla, cambio de filtros, paginacion prev/next, sorting.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Cambiar grid -> table -> grid -> table y verificar que siempre empieza desde pokemon #1 |
| Manual | Probar con distintos row limits (50, 100) que el comportamiento es consistente |
| Manual | Verificar que filtros, sorting y paginacion siguen funcionando tras el fix |

## Criterios de aceptacion

- [x] Al cambiar de grid a tabla, la tabla siempre muestra desde el pokemon #1
- [x] Repetir el ciclo grid->tabla->grid->tabla multiples veces sin perder pokemon
- [x] La paginacion prev/next sigue funcionando correctamente
- [x] Filtros y sorting no se ven afectados por el fix

## Notas

- La funcion `getCurrentPageItems()` esta en `frontend/src/pages/pokedex.ts` ~linea 1137-1146
- El handler de toggle de vista esta ~linea 1178-1207
- `lastRenderedItems` se asigna en `renderGrid()` (linea ~328) y `renderCurrentView()` (linea ~375)
