# Vista de tabla con todas las habilidades en Explorar > Habilidades

**ID**: 0035-abilities-table-view-all
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Reemplazar la vista de tarjetas en la pestaña Explorar > Habilidades por una vista de tabla ordenable (idéntica al patrón de Movimientos), y cargar **todas** las habilidades disponibles en PokeAPI (~300) en lugar de las 45 hardcodeadas en `KNOWN_ABILITIES`.

**Problemas actuales**:
1. Solo se muestran 45 habilidades hardcodeadas de las ~300 disponibles.
2. La vista de tarjetas no es consistente con la vista de movimientos (que ya usa tabla).
3. Las habilidades se cargan una a una con `GetAbility()` (lazy fetch individual).

**Solución**:
1. Crear `FetchAllAbilities()` en Shell con bulk fetch concurrente + cache en memoria (mismo patrón que `FetchAllMoves()`).
2. Exponer binding `GetAllAbilities()` en APP.
3. Reescribir `abilities.ts` con tabla ordenable en lugar de grid de tarjetas.

## Capas afectadas
- **Core**: Agregar `FetchAllAbilities` al interface `PokemonFetcher` en `ports.go`.
- **Shell**: Implementar `FetchAllAbilities()` en `pokeapi_meta.go` con concurrencia controlada y cache.
- **APP**: Exponer nuevo binding `GetAllAbilities()`.
- **Frontend**: Reescribir `abilities.ts` con tabla ordenable (columnas: Nombre, Descripción, Pokémon).

## Dependencias externas nuevas
Ninguna.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/ports.go` | modificar | Agregar `FetchAllAbilities() ([]Ability, error)` al interface |
| `shell/pokeapi_meta.go` | modificar | Implementar `FetchAllAbilities()` con bulk fetch + cache |
| `app/bindings.go` | modificar | Exponer `GetAllAbilities() ([]core.Ability, error)` |
| `frontend/src/pages/explore/abilities.ts` | reescribir | Tabla ordenable con búsqueda |
| `frontend/src/styles/_explore.scss` | modificar | Reusar/extender estilos de tabla de movimientos |

## Plan de implementacion

### Fase 1 — Backend: fetch masivo de habilidades
1. En `core/ports.go`, agregar `FetchAllAbilities() ([]Ability, error)` al interface `PokemonFetcher`.
2. En `shell/pokeapi_meta.go`, implementar `FetchAllAbilities()`:
   - Llamar a `FetchAbilityList(0, 2000)` para obtener todos los nombres.
   - Iterar con concurrencia controlada (15 goroutines, mismo patrón que moves).
   - Cachear resultado en memoria (`allAbilitiesCache`).
   - Retornar `[]core.Ability`.
3. En `app/bindings.go`, agregar `GetAllAbilities()` que llame a `fetcher.FetchAllAbilities()`.

### Fase 2 — Frontend: tabla ordenable
4. Reescribir `initAbilities()` en `abilities.ts`:
   - Al montar, llamar `GetAllAbilities()` con indicador de carga.
   - Renderizar tabla con columnas: **Nombre**, **Descripción**, **Pokémon** (cantidad).
   - Headers clickeables para ordenar asc/desc.
   - Filtro de búsqueda por nombre.
   - Contador de resultados.
5. Actualizar estilos en `_explore.scss`:
   - Reusar clases `poke-table` y patrón de `moves-table`.
   - Responsive: scroll horizontal en móvil.

### Fase 3 — Pulido
6. Asegurar rendimiento con ~300 habilidades.
7. Animación fade-in al cargar la tabla.

## Tests
| Archivo | Que se testea |
|---------|---------------|
| Manual | `GetAllAbilities()` retorna ~300 habilidades |
| Manual | La tabla renderiza correctamente |
| Manual | Ordenación por cada columna funciona (asc/desc) |
| Manual | Filtro de búsqueda funciona |
| Manual | Responsive: scroll horizontal en móvil |

## Criterios de aceptacion
- [x] Se muestran todas las habilidades de PokeAPI (~300) en la pestaña
- [x] Vista de tabla con columnas: Nombre, Descripción, Pokémon (cantidad)
- [x] Click en header ordena la columna asc/desc
- [x] Filtro de búsqueda por nombre funciona
- [x] Contador de resultados visible
- [x] Carga completa en tiempo razonable (cache en backend)
- [x] Responsive en móvil con scroll horizontal
- [x] Consistente visualmente con la vista de Movimientos

## Notas
- Seguir exactamente el mismo patrón de `FetchAllMoves()` en `shell/pokeapi_moves.go` para la implementación del bulk fetch.
- Las habilidades son ~300 vs ~920 movimientos, así que la carga será más rápida.
- La columna Descripción puede ser larga; truncar con CSS (`text-overflow: ellipsis`) o limitar a una línea.
