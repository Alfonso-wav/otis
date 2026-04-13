# Vista detalle Pokemon: añadir poder y clase (fisico/especial) a movimientos

**ID**: 0159-detail-moves-power-category
**Estado**: todo
**Fecha**: 2026-04-13

---

## Descripcion

La tabla de movimientos en la vista individual de un Pokemon (`frontend/src/pages/pokedex.ts:1141 renderMoves`) actualmente muestra solo: nombre del movimiento, nivel y metodo de aprendizaje. Añadir dos columnas:

- **Poder** (`Power`) — numero, o `—` si es movimiento de estado (`status`).
- **Clase** (`Category`) — `fisico` / `especial` / `estado` (con icono o badge, y i18n).

Los datos ya existen en `core.Move` (`core/domain.go:14`) con campos `Power int` y `Category string` y los expone `FetchMove(name)` / `GetAllMoves()`. El objeto actual `PokemonMoveEntry` (`core/domain.go:49`) **no** los lleva. Enfoque: cargar `GetAllMoves()` una sola vez en el frontend y construir un lookup `slug -> { Power, Category }` para enriquecer la tabla al renderizar.

## Capas afectadas

- **Core**: ninguna (los tipos ya existen).
- **Shell**: ninguna (endpoint `GetAllMoves` ya existe).
- **APP**: solo frontend.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/utils/move-names.ts` | modificar | Ampliar el cache: ademas de los nombres localizados, guardar `power` y `category` por slug. Exponer `getMovePower(slug)` y `getMoveCategory(slug)`. |
| `frontend/src/pages/pokedex.ts` | modificar | En `renderMoves` añadir dos columnas (`power`, `category`) a la tabla. Añadir sort por power (numerico). Badge con icono/color por categoria. |
| `frontend/src/styles/_pokedex.scss` o `_detail.scss` | modificar | Estilos de badge fisico/especial/estado (colores: fisico rojo-naranja, especial azul, estado gris). |
| `frontend/src/locales/en.json` | modificar | Claves: `detail.movePowerCol`, `detail.moveCategoryCol`, `detail.moveCategory.physical`, `detail.moveCategory.special`, `detail.moveCategory.status`. |
| `frontend/src/locales/es.json` | modificar | Mismas claves: `"Poder"`, `"Clase"`, `"Fisico"`, `"Especial"`, `"Estado"`. |

## Plan de implementacion

1. En `move-names.ts`:
   - Revisar si ya usa `GetAllMoves`. Si solo carga nombres, extender la carga para guardar tambien `Power` y `Category`.
   - Exportar `getMovePower(slug: string): number | null` y `getMoveCategory(slug: string): "physical" | "special" | "status" | null`.
2. En `pokedex.ts:renderMoves`:
   - Tras `await loadMoveNames();` el cache ya tiene todo.
   - Añadir tipo `SortCol = "name" | "method" | "level" | "power"`.
   - En la fila renderizada incluir:
     ```html
     <td class="move-power">${power ?? "—"}</td>
     <td class="move-category"><span class="move-cat-badge move-cat-${category}">${t("detail.moveCategory." + category)}</span></td>
     ```
   - Manejar moves no encontrados en el cache → `—` y `status`.
   - Añadir `<th>` para `power` y `category` con `data-sort` en power.
3. Sort por power: numerico; `null`/`status` al final en `asc`, al principio en `desc` (consistencia con otras tablas).
4. CSS badges (3 variantes):
   ```scss
   .move-cat-badge { padding: 2px 8px; border-radius: 4px; font-size: 0.8rem; font-weight: 600; }
   .move-cat-physical { background: #c43b2f; color: #fff; }
   .move-cat-special  { background: #4a79d0; color: #fff; }
   .move-cat-status   { background: #8a8a8a; color: #fff; }
   ```
5. i18n: añadir las claves en `en.json` y `es.json` simultaneamente. Escuchar `locale-changed` ya esta implementado en la vista detalle (confirmar).
6. Verificar sort combinado: activar y desactivar sorts, verificar que las filas filtradas por metodo tambien se ordenan por power correctamente.
7. Mobile: la tabla ya tiene `overflow-x: auto`; verificar que con 2 columnas nuevas no se rompe en 360px. Si hace falta, ocultar `level` en mobile o similar (ver patrones de `0156-type-chart-android-responsive.md`).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual desktop | Abrir detalle de Charizard: Flamethrower muestra `Poder 90`, `Especial` |
| Manual desktop | Tackle muestra `Poder 40`, `Fisico` |
| Manual desktop | Toxic muestra `Poder —`, `Estado` |
| Manual desktop | Sort por power asc/desc funciona, resetea en el 3er click |
| Manual desktop | Cambio de idioma (locale-changed) re-renderiza header y badges |
| Manual mobile 360px | Tabla sigue legible, scroll horizontal activa si hace falta |
| Manual APK Android | Al abrir detalle por primera vez, `GetAllMoves` resuelve via REST y las columnas se pueblan |

## Criterios de aceptacion

- [ ] Tabla de moves del detalle incluye columnas `Poder` y `Clase`
- [ ] Movimientos de estado muestran `—` como poder y badge `Estado`
- [ ] Badges con color distintivo por categoria (fisico/especial/estado)
- [ ] Sort por poder funciona asc/desc/off
- [ ] Paridad i18n EN/ES (`en.json` y `es.json` con las mismas claves)
- [ ] Cambio de locale re-renderiza la tabla correctamente
- [ ] Mobile 360px sin overflow roto
- [ ] Sin regresion en el sort existente de name/method/level
- [ ] APK Android: funciona igual (REST handler `GetAllMoves` ya existe, verificar)

## Notas

- El cache de moves se carga una sola vez por sesion via `FetchAllMoves` (en backend) con concurrencia y cache en memoria — seguro llamarlo desde el detalle.
- El texto localizado de la categoria va en las claves nuevas, no reutilizar strings de otras tablas sin verificar consistencia.
- Si se detecta que `getLocalizedMoveName` y los datos de power viven en sitios distintos, consolidarlos en un unico modulo `move-cache.ts` o similar (refactor opcional, no bloqueante).
- El `0155-pokemon-detail-moves-i18n.md` ya se ocupo de la i18n de nombres; esta tarea extiende ese flujo.
