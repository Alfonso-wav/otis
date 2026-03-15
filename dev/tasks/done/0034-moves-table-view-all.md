# Vista de tabla con todos los movimientos en Explorar > Movimientos

**ID**: 0034-moves-table-view-all
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Reemplazar la vista de tarjetas en la pestaña Explorar > Movimientos por una vista de tabla con ordenación columnar, y cargar **todos** los movimientos disponibles en PokeAPI (~920) en lugar de los 47 hardcodeados actualmente en `KNOWN_MOVES`.

**Problemas actuales**:
1. Solo se muestran 47 movimientos hardcodeados de los ~920 disponibles.
2. La vista de tarjetas no es eficiente para explorar una lista tan larga.
3. Los movimientos se cargan uno a uno con `GetMove()` (lazy fetch individual), lo cual sería inviable para 920 movimientos.

**Solución**:
1. Crear un endpoint backend que devuelva todos los movimientos de una sola vez (bulk fetch + cache).
2. Reemplazar la grid de tarjetas por una tabla HTML con columnas ordenables (como la tabla del Pokédex).

## Capas afectadas
- **Core**: Agregar tipo para lista completa de movimientos con datos resumidos.
- **Shell**: Implementar fetch masivo de movimientos usando `FetchMoveList` + `FetchMove` con cache en memoria.
- **APP**: Exponer nuevo binding `GetAllMoves()` al frontend.
- **Frontend**: Reescribir `moves.ts` con tabla ordenable en lugar de tarjetas.

## Dependencias externas nuevas
Ninguna.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Agregar tipo `MoveEntry` (versión resumida para tabla) si es necesario, o reusar `Move` |
| `shell/pokeapi.go` | modificar | Agregar método `FetchAllMoves()` que haga bulk fetch con cache |
| `app/bindings.go` | modificar | Exponer `GetAllMoves() ([]core.Move, error)` |
| `frontend/src/pages/explore/moves.ts` | reescribir | Tabla ordenable con todas las columnas |
| `frontend/src/styles/_explore.scss` | modificar | Estilos de tabla (reusar patrón de `_pokemon.scss`) |

## Plan de implementacion

### Fase 1 — Backend: fetch masivo de movimientos
1. En `shell/pokeapi.go`, agregar método `FetchAllMoves()`:
   - Llamar a `FetchMoveList(0, 2000)` para obtener la lista completa de nombres.
   - Iterar y llamar `FetchMove(name)` para cada uno, con concurrencia controlada (ej. 10 goroutines).
   - Cachear el resultado en memoria para no repetir las ~920 llamadas HTTP en cada consulta.
   - Retornar `[]core.Move`.
2. En `app/bindings.go`, agregar `GetAllMoves()` que llame a `FetchAllMoves()`.

### Fase 2 — Frontend: tabla ordenable
3. Reescribir `initMoves()` en `moves.ts`:
   - Al montar, llamar `GetAllMoves()` y mostrar un indicador de carga.
   - Renderizar tabla HTML con columnas: **Nombre**, **Tipo**, **Categoría**, **Poder**, **Precisión**, **PP**, **Prioridad**.
   - Cada header es clickeable para ordenar asc/desc (seguir patrón de `pokedex.ts`).
   - Mantener los filtros existentes (búsqueda, categoría) y el filtro por tipo.
4. Actualizar estilos en `_explore.scss`:
   - Reusar clases de tabla de `_pokemon.scss` (`poke-table`) o crear variante `moves-table`.
   - Badges de tipo con color.
   - Iconos de categoría (physical/special/status).
   - Responsive: scroll horizontal en móvil con columnas sticky.

### Fase 3 — Pulido
5. Asegurar que la tabla soporta los 920+ movimientos sin problemas de rendimiento.
6. Mantener animación sutil al cargar (fade-in de la tabla, no stagger por fila).

## Tests
| Archivo | Que se testea |
|---------|---------------|
| Manual | `GetAllMoves()` retorna 900+ movimientos |
| Manual | La tabla renderiza correctamente con todos los movimientos |
| Manual | Ordenación por cada columna funciona (asc/desc) |
| Manual | Filtros de búsqueda, categoría y tipo funcionan |
| Manual | Responsive: scroll horizontal en móvil |

## Criterios de aceptacion
- [x] Se muestran todos los movimientos de PokeAPI (~920) en la pestaña
- [x] Vista de tabla con columnas: Nombre, Tipo, Categoría, Poder, Precisión, PP, Prioridad
- [x] Click en header ordena la columna asc/desc
- [x] Filtro de búsqueda por nombre funciona
- [x] Filtro por categoría (physical/special/status) funciona
- [x] Badge de tipo con color en cada fila
- [x] Icono de categoría en cada fila
- [x] Carga completa en tiempo razonable (cache en backend)
- [x] Responsive en móvil con scroll horizontal

## Notas
- La primera carga será lenta (~920 requests a PokeAPI). El cache en memoria evita que se repita. Considerar persistir cache en disco o usar batch si PokeAPI lo soporta.
- PokeAPI permite `?limit=2000` en el endpoint de lista, pero solo devuelve nombres. Los detalles requieren fetch individual por movimiento.
- Alternativa: construir un JSON estático con todos los movimientos y servirlo como asset (evita 920 requests). Evaluar con el usuario.
