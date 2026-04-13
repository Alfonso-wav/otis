# Traducir movimientos en vista detalle Pokemon

**ID**: 0155-pokemon-detail-moves-i18n
**Estado**: todo
**Fecha**: 2026-04-13

---

## Descripcion

En la version Española, la tabla de movimientos de la vista individual de Pokemon (`renderMoves` en `frontend/src/pages/pokedex.ts`) muestra los nombres de movimiento en ingles (ej. "Thunder Punch" en vez de "Puño Trueno"). El render actual usa `capitalize(m.Name.replace(/-/g, " "))` sobre el slug ingles. El tipo `core.Move` ya expone `NameEs` (ver `frontend/src/pages/explore/moves.ts:86`), pero `core.PokemonMoveEntry` (`core/domain.go:49`) solo tiene `Name`.

Objetivo: cuando `getLocale() === "es"`, mostrar el nombre en Español del movimiento; si no hay traduccion, fallback a ingles.

## Capas afectadas

- **Core**: añadir campo `NameEs string` a `PokemonMoveEntry` en `core/domain.go`.
- **Shell**: poblar `NameEs` en `shell/pokeapi.go` (`toDomainPokemon`, ~linea 148-157). Opciones:
  - A) Extraer traduccion del endpoint `/move/{id}` (campo `names[].name` con `language.name == "es"`). Requiere N llamadas extra por Pokemon — caro.
  - B) (Preferida) Reusar el cache existente de moves (el que sirve `GetAllMoves`/explore/moves) para mapear `Name → NameEs` en el frontend sin tocar backend. Cargar la lista una sola vez, cachear Map<string,string>, y en `renderMoves` resolver el nombre localizado.
- **APP**: sin cambios si se elige opcion B. Si se elige A, regenerar bindings Wails + verificar paridad REST (`app/mobile/handlers.go`).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | `renderMoves` resuelve nombre localizado via cache de moves; escuchar `locale-changed` para re-render |
| `frontend/src/utils/move-names.ts` | crear | Helper `getLocalizedMoveName(slug)` con cache `Map<string,string>` cargado lazy desde `GetAllMoves` |
| `frontend/src/api.ts` | revisar | Confirmar que `GetAllMoves` devuelve `NameEs` tanto en Wails como REST |
| `core/domain.go` | modificar (solo si opcion A) | Añadir `NameEs` a `PokemonMoveEntry` |
| `shell/pokeapi.go` | modificar (solo si opcion A) | Poblar `NameEs` en `toDomainPokemon` |

## Plan de implementacion

1. Confirmar que `GetAllMoves` devuelve `NameEs` en mobile/REST (revisar `app/mobile/handlers.go`).
2. Crear `frontend/src/utils/move-names.ts` con:
   - `loadMoveNames(): Promise<void>` que llama `GetAllMoves` una vez y llena `Map<slug, NameEs>`.
   - `getLocalizedMoveName(slug): string` que devuelve NameEs si locale es "es" y existe, si no `capitalize(slug.replace(/-/g," "))`.
3. En `pokedex.ts` `renderMoves`: llamar `await loadMoveNames()` antes del primer render; sustituir `capitalize(m.Name.replace(/-/g, " "))` por `getLocalizedMoveName(m.Name)`.
4. Registrar listener `locale-changed` en la seccion moves para re-renderizar la tabla.
5. Verificar que filtro/orden por `name` ordena por el nombre mostrado (localeCompare sobre nombre localizado).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Cambiar idioma a ES en settings, abrir detalle Pokemon, verificar movimientos en Español |
| Manual | Cambiar a EN, verificar que vuelve a ingles sin recargar |
| Manual | Mismo flujo en APK Android (paridad REST) |

## Criterios de aceptacion

- [ ] En locale "es", los nombres de movimientos en detalle Pokemon salen en Español
- [ ] En locale "en", salen en ingles (comportamiento actual)
- [ ] `locale-changed` actualiza la tabla sin recargar pagina
- [ ] Funciona identico en build desktop (Wails) y APK Android (REST)
- [ ] Ordenacion por columna "name" usa nombre localizado
- [ ] Fallback a ingles si `NameEs` vacio

## Notas

- Reusar cache de moves: no inflar red con llamadas por Pokemon.
- Ver `frontend/src/pages/explore/moves.ts:86` como patron de busqueda con NameEs existente.
- La seccion ya escucha `locale-changed` a nivel de vista detalle — verificar que llega hasta el render de moves.
