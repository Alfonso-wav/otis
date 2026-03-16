# Caché de ordenamientos en tablas

**ID**: 0084-cache-table-sort-orders
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Al ordenar las tablas por columnas, el rendimiento es lento porque:

1. **Pokedex**: la primera vez necesita cargar ~1000+ Pokémon individualmente (batches de 50 via `GetPokemon`). Cada cambio de columna/dirección re-ordena todo el array desde cero (`sortPokemonData` crea copia + `Array.sort`).
2. **Moves/Abilities/Encounters**: cada cambio de columna también re-ordena todo el dataset.
3. No se cachean los resultados de ordenamiento — se descartan y recalculan cada vez.

**Solución**: implementar un `SortCache` genérico que pre-compute y almacene los índices ordenados (asc + desc) por cada columna. Al cambiar de columna o dirección, solo se hace un lookup en el caché en vez de re-ordenar. El caché se invalida cuando cambian los datos (filtros, nueva carga).

## Capas afectadas

- **APP (frontend)**: Lógica de sorting en pokedex.ts, moves.ts, abilities.ts. Nuevo módulo utilitario de caché.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/utils/sort-cache.ts` | crear | Módulo genérico `SortCache<T>` que pre-computa y cachea arrays ordenados por columna+dirección |
| `frontend/src/pages/pokedex.ts` | modificar | Reemplazar `sortPokemonData` por uso de `SortCache`, pre-computar ordenes al cargar datos, invalidar al filtrar |
| `frontend/src/pages/explore/moves.ts` | modificar | Usar `SortCache` para cachear ordenes de moves |
| `frontend/src/pages/explore/abilities.ts` | modificar | Usar `SortCache` para cachear ordenes de abilities |

## Plan de implementacion

### Parte 1: SortCache genérico

1. **Crear `frontend/src/utils/sort-cache.ts`** con una clase/función `SortCache<T>`:
   ```typescript
   interface SortDef<T> {
     key: string;
     compare: (a: T, b: T) => number; // comparador para asc
   }

   class SortCache<T> {
     private cache = new Map<string, T[]>(); // "column:asc" -> sorted array
     private sourceData: T[] = [];

     constructor(private columns: SortDef<T>[]) {}

     setData(data: T[]): void {
       this.sourceData = data;
       this.cache.clear();
       // Pre-computar todos los ordenes de golpe
       for (const col of this.columns) {
         const asc = [...data].sort(col.compare);
         this.cache.set(`${col.key}:asc`, asc);
         this.cache.set(`${col.key}:desc`, [...asc].reverse());
       }
     }

     get(column: string, direction: 'asc' | 'desc'): T[] {
       return this.cache.get(`${column}:${direction}`) ?? this.sourceData;
     }

     invalidate(): void {
       this.cache.clear();
     }
   }
   ```

2. **Estrategia de pre-cómputo**: al llamar `setData`, se ordenan todas las columnas de una vez. Para ~1000 Pokémon con ~8 columnas esto es ~8 sorts de 1000 items = trivial (~2-5ms total). El trade-off es un poco más de memoria (8 arrays de 1000 refs) pero el acceso posterior es O(1).

3. **Alternativa lazy**: si el pre-cómputo upfront es mucho (datasets muy grandes), computar solo bajo demanda y cachear. Para los tamaños actuales (~1000 Pokémon, ~900 moves) el pre-cómputo es preferible.

### Parte 2: Integrar en Pokedex

4. **En `pokedex.ts`**: crear instancia de `SortCache<Pokemon>` con los comparadores de cada columna (id, name, hp, atk, def, spa, spd, vel, total).

5. **Al completar `ensureAllPokemonLoaded`**: llamar `sortCache.setData(allPokemon)` para pre-computar todos los ordenes. Esto se hace una sola vez hasta que cambien los filtros.

6. **En el handler de click de columna**: reemplazar `sortPokemonData(allPokemon, col, dir)` por `sortCache.get(col, dir)`. Esto es instantáneo.

7. **Al cambiar filtros** (`resetSorting`): llamar `sortCache.invalidate()` para que se re-compute con el nuevo dataset filtrado.

### Parte 3: Integrar en Moves y Abilities

8. **En `moves.ts`**: crear `SortCache<core.Move>` con comparadores para name, type, category, power, accuracy, pp, priority. Llamar `setData` al cargar los moves y al filtrar.

9. **En `abilities.ts`**: crear `SortCache<core.Ability>` con comparadores para name, pokemonCount. Llamar `setData` al cargar.

### Parte 4: Optimización adicional (opcional)

10. **Considerar pre-cargar los datos de Pokémon en background** al entrar a la pestaña Pokedex (antes de que el usuario pulse ordenar). Así cuando pulse, el `SortCache` ya tiene todo listo.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `frontend/src/utils/sort-cache.test.ts` | SortCache devuelve datos ordenados correctamente por columna y dirección |
| `frontend/src/utils/sort-cache.test.ts` | SortCache invalida caché correctamente al llamar setData con nuevos datos |
| `frontend/src/utils/sort-cache.test.ts` | SortCache.get devuelve sourceData si se pide columna inexistente |
| (visual) | Ordenar pokedex por cualquier columna es instantáneo tras la carga inicial |
| (visual) | Cambiar entre asc/desc en la misma columna es instantáneo |
| (visual) | Cambiar de columna de ordenamiento es instantáneo |
| (visual) | Los filtros invalidan el caché y el re-ordenamiento funciona correctamente |
| (visual) | Moves table: ordenar por cualquier columna es rápido |
| (visual) | Abilities table: ordenar por cualquier columna es rápido |

## Criterios de aceptacion

- [x] Existe un módulo `SortCache` genérico y reutilizable en `frontend/src/utils/sort-cache.ts`
- [x] El Pokedex usa `SortCache` — ordenar por columna tras la carga inicial es instantáneo (<50ms)
- [x] Cambiar dirección (asc↔desc) es instantáneo en todas las tablas
- [x] Cambiar de columna de ordenamiento es instantáneo en todas las tablas
- [x] Los filtros invalidan correctamente el caché de ordenamiento
- [x] Moves table usa `SortCache`
- [x] Abilities table usa `SortCache`
- [x] No hay regresiones en el comportamiento actual de sorting
- [x] La memoria adicional es razonable (~8 arrays de referencias, no duplica objetos)

## Notas

- El cuello de botella real del primer sort en Pokedex es la carga de datos (`ensureAllPokemonLoaded`), no el sort en sí. El caché resuelve los sorts subsiguientes. Para la carga inicial, una mejora futura sería un endpoint backend que devuelva todos los Pokémon con stats de golpe (como ya se hace con moves/abilities).
- Los arrays cacheados almacenan referencias a los mismos objetos, no copias profundas, así que el overhead de memoria es mínimo.
- `[...asc].reverse()` para desc es más eficiente que re-ordenar con comparador invertido.
