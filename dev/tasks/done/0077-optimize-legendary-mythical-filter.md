# Optimizar filtro Legendario/Mítico sin generación seleccionada

**ID**: 0077-optimize-legendary-mythical-filter
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

El filtro de Legendario y Mítico en la Pokédex es extremadamente lento cuando no hay ninguna generación seleccionada. Actualmente se hacen ~1036 peticiones HTTP: ~11 para listar los ~1025 Pokémon y luego ~1025 llamadas individuales a `GetPokemonSpecies()` en lotes de 10 secuenciales. La solución es crear un endpoint backend que devuelva directamente la lista de Pokémon legendarios/míticos, precalculada y cacheada, eliminando las ~1025 llamadas individuales.

## Capas afectadas

- **Core**: Nuevo tipo `SpeciesClassification` y función pura para filtrar por legendary/mythical desde un mapa precalculado.
- **Shell**: Nuevo método que recorre todas las species y construye un mapa `name → {isLegendary, isMythical}`, con caché en memoria.
- **APP**: Nuevo binding `GetAllSpeciesClassifications()` que expone el mapa al frontend.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir struct `SpeciesClassification` |
| `core/pokemon.go` | modificar | Función pura para filtrar lista por clasificación |
| `core/ports.go` | modificar | Añadir método al interface `PokemonFetcher` |
| `shell/pokeapi_species.go` | modificar | Implementar fetch masivo de clasificaciones con caché |
| `app/bindings.go` | modificar | Exponer nuevo binding `GetAllSpeciesClassifications()` |
| `frontend/src/pages/pokedex.ts` | modificar | Usar nuevo endpoint en vez de llamadas individuales a `GetPokemonSpecies` |

## Plan de implementacion

1. **Core**: Añadir `SpeciesClassification` struct con campos `IsLegendary bool` e `IsMythical bool`. Añadir función `FilterByClassification(items []PokemonListItem, classifications map[string]SpeciesClassification, legendary, mythical bool) []PokemonListItem`.
2. **Core ports**: Añadir `GetAllSpeciesClassifications() (map[string]SpeciesClassification, error)` al interface.
3. **Shell**: Implementar el método que itere sobre todas las species (usando el listado base de Pokemon), obtenga `is_legendary`/`is_mythical` de cada una, y cachee el resultado completo en memoria con `sync.Once` o similar. Usar concurrencia (goroutines con semáforo) para paralelizar las peticiones al PokeAPI.
4. **APP bindings**: Exponer `GetAllSpeciesClassifications()` como binding de Wails.
5. **Frontend**: En `loadFiltered()`, cuando `filter.legendary || filter.mythical` y no hay gens/types, llamar a `GetAllSpeciesClassifications()` una sola vez (1 llamada IPC), filtrar el listado completo en memoria del frontend, y eliminar el bucle de `filterByLegendary()` con sus ~1025 llamadas individuales.
6. **Frontend cache**: Almacenar el mapa de clasificaciones en una variable de módulo para reutilización inmediata en filtrados posteriores.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/pokemon_test.go` | `FilterByClassification` filtra correctamente por legendary, mythical, y ambos |

## Criterios de aceptacion

- [ ] Filtrar por Legendario sin generación carga en < 3 segundos (vs ~30+ actual)
- [ ] Filtrar por Mítico sin generación carga en < 3 segundos
- [ ] El mapa de clasificaciones se cachea: la segunda aplicación del filtro es instantánea
- [ ] Filtros combinados (gen + legendary) siguen funcionando correctamente
- [ ] No hay regresiones en los filtros existentes (gen, tipo, legendary+gen, etc.)

## Notas

- El cuello de botella actual está en `filterByLegendary()` (pokedex.ts:236-272) que hace ~1025 llamadas a `GetPokemonSpecies()` en lotes secuenciales de 10.
- PokeAPI no tiene un endpoint `/legendary` directo, por lo que el backend debe construir el mapa iterando species.
- La caché backend con `sync.Once` garantiza que solo se construye una vez por sesión de la app.
- Considerar usar `errgroup` con límite de concurrencia (e.g., 20-50 goroutines) para el fetch inicial en el backend.
