# Task 0113 — Corregir i18n sección Movimientos y recarga al cambio de idioma

## Estado: done

## Goal
La sección de Movimientos en la vista de detalle tiene dos problemas:
1. **Spanglish en `es.json`**: claves del bloque `detail` que están en inglés aunque el idioma activo sea español.
2. **No reacciona al cambio de idioma**: cuando el usuario cambia de idioma en Settings, el detalle ya renderizado (incluyendo botones de filtro de movimientos, cabeceras de columna y etiquetas de método) no se actualiza. Es necesario volver a renderizar el detalle activo con `renderDetail(p)` cuando se dispara `locale-changed`.

## Context

### Spanglish en `es.json` — bloque `detail`

Archivo: `frontend/src/locales/es.json` (líneas ~61-90)

Claves en inglés que deben traducirse:
| Clave | Valor actual (incorrecto) | Valor correcto |
|---|---|---|
| `detail.lore` | `"Lore"` | `"Historia"` |
| `detail.legendary` | `"Legendary"` | `"Legendario"` |
| `detail.mythical` | `"Mythical"` | `"Mítico"` |
| `detail.category` | `"Category"` | `"Categoría"` |
| `detail.habitat` | `"Habitat"` | `"Hábitat"` |
| `detail.shape` | `"Shape"` | `"Forma"` |
| `detail.loreError` | `"Could not load lore data."` | `"No se pudo cargar la información."` |
| `detail.moveMethod.levelUp` | `"Nivel"` | `"Subida de nivel"` (choca con la columna "Nivel") |

El bloque `encounters` en `es.json` también tiene valores en inglés — revisar y traducir en el mismo paso.

### Re-render al cambio de idioma

Archivo: `frontend/src/pages/pokedex.ts` (línea ~1482)

El handler `locale-changed` no re-renderiza el detalle:
```typescript
document.addEventListener("locale-changed", () => {
  // solo actualiza la lista, no el detalle
});
```

Cuando el usuario está en la vista de detalle y cambia el idioma, los filtros de movimientos, cabeceras y etiquetas de método quedan en el idioma anterior hasta que se recarga el detalle manualmente.

**Fix**: mantener una referencia al Pokémon activo en la vista de detalle y llamar `renderDetail(currentPokemon)` dentro del handler `locale-changed` cuando la vista activa sea el detalle.

## Approach

1. **`es.json`**: corregir las claves listadas en la tabla de arriba. No cambiar claves de `en.json` salvo que haya algún error de contenido.
2. **`pokedex.ts`**: 
   - Guardar el `Pokemon` activo en una variable de módulo (e.g., `let currentDetailPokemon: Pokemon | null = null`) que se asigna en `renderDetail`.
   - En el handler de `locale-changed`, si `currentDetailPokemon !== null` y la vista activa es el detalle, llamar `renderDetail(currentDetailPokemon)`.

## Files to modify
- `frontend/src/locales/es.json`
- `frontend/src/pages/pokedex.ts`

## Acceptance criteria
- [x] Todos los textos visibles en la sección de detalle aparecen en español cuando el idioma activo es español
- [x] `detail.moveMethod.levelUp` ya no choca visualmente con la columna "Nivel"
- [x] Al cambiar de idioma con el detalle de un Pokémon abierto, los filtros de movimientos, cabeceras y etiquetas se actualizan sin necesidad de volver a abrir el detalle
- [x] `en.json` no tiene cambios de contenido
