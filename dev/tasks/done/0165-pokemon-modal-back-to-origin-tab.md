# Task 0165 — Volver a la pestaña de origen desde el detalle de Pokemon

## Estado: done

## Goal
Cuando se abre el detalle de un Pokemon desde una pestaña distinta a Pokedex (Explore > Types, Explore > Moves, modales de localizacion/encuentro, etc.), al pulsar la flecha de atras el usuario debe regresar a la pestaña desde la que abrio el detalle, no siempre a Pokedex.

## Contexto tecnico

### Comportamiento actual
- El detalle de Pokemon (`detail-view`) vive dentro del contenedor de la pestaña **Pokedex** (`frontend/src/pages/pokedex.ts:612 showDetail`, `:1598 listView`, `:1599 detailView`).
- Los 4 puntos de entrada desde otras pestañas hacen `click()` sobre `[data-tab="pokedex"]` antes de buscar:
  - `frontend/src/components/move-pokemon-modal.ts:17`
  - `frontend/src/components/pokemon-type-modal.ts:77`
  - `frontend/src/components/location-encounter-modal.ts:119`
  - `frontend/src/pages/explore/types.ts:136`
- Resultado: al abrir el detalle, el router (`frontend/src/router.ts:24 navigate`) cambia a `pokedex` y marca el tab activo. Al pulsar back (`pokedex.ts:1664`), solo se oculta el `detailView` → el usuario queda en Pokedex aunque venga de Explore u otra.

### Diseño propuesto
1. Exponer en `router.ts` un getter `getActiveId()` que devuelva `activeId`.
2. En `pokedex.ts`, añadir variable de modulo `detailOriginTab: string | null = null` y una funcion `setDetailOrigin(tab: string | null)` exportada.
3. Cada llamador externo (los 4 ficheros listados) debe:
   - Capturar `getActiveId()` ANTES de cambiar a pokedex.
   - Llamar a `setDetailOrigin(origin)` si `origin && origin !== "pokedex"`.
4. Modificar el handler del `back-btn` (`pokedex.ts:1664`):
   - Si `detailOriginTab` existe y no es `"pokedex"`: ocultar `detailView` (equivalente a `showView(listView, detailView)` actual) **y** `navigate(detailOriginTab)`.
   - Resetear `detailOriginTab = null` en todos los casos.
5. Cuando el detalle se abre desde dentro de Pokedex (click en grid/tabla/autocomplete, `pokedex.ts:368,378,476,1334,1730`), mantener `detailOriginTab = null` para conservar comportamiento actual.

### Edge cases
- Si el usuario navega manualmente a otra pestaña mientras el detalle esta abierto (cambio de tab sin back), resetear `detailOriginTab`. Puede hacerse en `router.ts navigate()` emitiendo un evento `tab-changed` que pokedex escuche, o añadiendo un hook de limpieza en `navigate()` via import de pokedex. Preferible: evento custom `document.dispatchEvent(new Event("tab-changed"))` desde el router para desacoplar.
- Abrir detalle desde builds (si existe algun path) — verificar durante la implementacion que no haya otros callers ademas de los 4 listados.

### Paridad Wails / REST
No aplica — cambio 100% frontend.

### i18n
No hay textos nuevos visibles.

## Acceptance criteria

- [ ] Abrir detalle desde Explore > Moves → click en fila → back → usuario queda en Explore > Moves con el submenu correcto.
- [ ] Abrir detalle desde Explore > Types (modal) → click en pokemon → back → usuario queda en Explore > Types.
- [ ] Abrir detalle desde modal de localizacion/encuentro → back → usuario queda en la pestaña de origen (Explore > Regions o donde estuviera).
- [ ] Abrir detalle desde Pokedex (grid, tabla, autocomplete, busqueda) → back → usuario queda en la lista de Pokedex (comportamiento actual intacto).
- [ ] Si el usuario cambia de pestaña manualmente con el detalle abierto y luego vuelve, back no teletransporta a una pestaña obsoleta.
- [ ] Probado en desktop (Wails) y en build de produccion del frontend.
- [ ] Probado en viewport 360px.

## Archivos afectados

- `frontend/src/router.ts` — exportar `getActiveId()`; opcionalmente emitir evento `tab-changed` en `navigate()`.
- `frontend/src/pages/pokedex.ts` — `detailOriginTab`, `setDetailOrigin()`, ajuste del handler `back-btn`, reset en `tab-changed`.
- `frontend/src/components/move-pokemon-modal.ts` — capturar origen antes de click.
- `frontend/src/components/pokemon-type-modal.ts` — idem.
- `frontend/src/components/location-encounter-modal.ts` — idem.
- `frontend/src/pages/explore/types.ts` — idem.

## Notas

- Mantener la logica dentro del frontend — no tocar backend Go.
- No cambiar la estructura DOM del `detail-view` ni su ubicacion dentro del contenedor de Pokedex. Basta con cambiar de tab al salir.
- Guardia frente a clicks rapidos en back-btn: el handler actual es async pero no tiene guard; si la tarea fuerza cambios mayores, añadir `isBackNavigating` flag + try/finally (ver CLAUDE.md).
