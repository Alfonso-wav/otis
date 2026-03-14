# 0021 — Ordenamiento global en la tabla Pokedex (todas las paginas)

## Descripcion

Actualmente, al ordenar una columna en la tabla de la Pokedex, solo se ordenan los ~20 Pokemon de la pagina visible. El usuario espera que al ordenar por una stat (ej. Ataque), se ordenen **todos** los Pokemon disponibles y se pagine sobre el resultado ordenado. Esto requiere cargar todos los datos antes de ordenar y paginar sobre la lista ordenada completa.

## Capas involucradas

- **Frontend (APP)**: `frontend/src/pages/pokedex.ts` — logica de sorting, paginacion y cache.
- **Core**: No se requieren cambios.
- **Shell**: No se requieren cambios. Los datos ya se obtienen de PokeAPI y se cachean en `pokemonDataCache`.

## Contexto actual

- `ListPokemon(offset, limit)` trae solo 20 Pokemon por pagina desde el backend.
- `pokemonDataCache` guarda los datos completos de cada Pokemon ya consultado.
- `sortPokemonData()` ordena un array de `Pokemon[]`, pero solo recibe los ~20 de la pagina actual.
- Al cambiar de pagina se llama `resetSorting()`, perdiendo el estado de orden.
- Cuando hay filtros activos, `filteredList` contiene todos los IDs que coinciden, pero solo se cargan los datos de la pagina actual.

## Plan de implementacion

### Paso 1 — Cargar todos los datos de Pokemon al activar sorting

Cuando el usuario haga clic en una cabecera para ordenar y no se tengan todos los Pokemon en cache, cargar los datos completos de todos los Pokemon (o al menos los del conjunto actual: sin filtros = todos; con filtros = los de `filteredList`). Mostrar un indicador de carga mientras se obtienen.

- Crear una funcion `ensureAllPokemonLoaded()` que:
  - Sin filtros: obtenga la lista completa de Pokemon (usar el `totalCount` ya conocido) y llame a `GetPokemon()` para cada uno que no este en cache.
  - Con filtros: haga lo mismo pero solo para los IDs en `filteredList`.
- Usar `Promise.all` con batches para no saturar el backend.

- [x] Completado

### Paso 2 — Mantener una lista ordenada global

Crear una variable `sortedFullList: Pokemon[] | null` que contenga todos los Pokemon ordenados cuando hay un sort activo:

```ts
let sortedFullList: Pokemon[] | null = null;
```

- Cuando se activa sorting: cargar todos los datos (Paso 1), aplicar `sortPokemonData()` sobre el array completo, guardar en `sortedFullList`.
- Cuando se desactiva sorting (`direction = null`): poner `sortedFullList = null` y volver al flujo de paginacion normal.

- [x] Completado

### Paso 3 — Paginar sobre la lista ordenada

Modificar la logica de paginacion para que, cuando `sortedFullList` no sea null, pagine sobre esa lista en lugar de hacer llamadas al backend:

- `getCurrentPageItems()` debe devolver un slice de `sortedFullList` cuando haya sort activo.
- `prevPage()` y `nextPage()` **no** deben llamar a `resetSorting()` cuando hay sort activo; deben paginar sobre `sortedFullList`.
- Actualizar `totalPages` basandose en `sortedFullList.length` cuando hay sort activo.

- [x] Completado

### Paso 4 — Mantener el sort al cambiar de pagina

Eliminar la llamada a `resetSorting()` en `prevPage()` y `nextPage()`. En su lugar, si hay sort activo, simplemente re-renderizar la pagina correspondiente desde `sortedFullList`.

- [x] Completado

### Paso 5 — Resetear sorting al cambiar filtros o vista

Cuando el usuario cambia filtros o alterna entre grid/table:
- Poner `sortedFullList = null`.
- Llamar `resetSorting()` como antes.
- Esto es correcto porque el conjunto base de datos cambia.

- [x] Completado

### Paso 6 — Indicador de carga durante la carga masiva

Mostrar un spinner o mensaje "Cargando todos los Pokemon para ordenar..." mientras `ensureAllPokemonLoaded()` trabaja. Deshabitar los clicks en cabeceras hasta que termine.

- [x] Completado

## Criterios de aceptacion

- [x]Al ordenar por cualquier columna, se ordenan TODOS los Pokemon (no solo la pagina actual).
- [x]La paginacion navega sobre la lista ordenada completa.
- [x]El indicador de sort se mantiene al cambiar de pagina.
- [x]Se muestra indicador de carga mientras se obtienen todos los datos.
- [x]Al quitar el sort (tercer clic), se vuelve al comportamiento normal de paginacion.
- [x]Al cambiar filtros o vista, se resetea el sorting.
- [x]Con filtros activos, el sort opera sobre el subconjunto filtrado completo.
- [x]No se introducen cambios en Core ni Shell.
- [x]Performance aceptable: la carga inicial de todos los datos puede tardar, pero una vez cacheados, los re-sorts son instantaneos.

## Dependencias externas nuevas

Ninguna.

## Notas

- Los datos ya cacheados en `pokemonDataCache` se reutilizan; solo se cargan los que faltan.
- Una vez que todos los datos estan en cache, cambiar de columna o direccion de sort es instantaneo (no requiere nuevas llamadas al backend).
- Considerar un batch size razonable (ej. 50 concurrentes) para no sobrecargar el backend con `GetPokemon()` simultaneos.
