# Fix: generación de equipos aleatorios siempre devuelve el mismo resultado

**ID**: 0053-fix-random-team-generation
**Estado**: todo
**Fecha**: 2026-03-15
**Depende de**: 0046-random-team-fill

---

## Descripcion

La función "Rellenar aleatorio" en la sección "Mis Equipos" genera siempre el mismo equipo. El problema está en `app/bindings.go` líneas 460-468: se itera la lista de 151 Pokémon pero se rompe el bucle en cuanto se tienen `slotsNeeded + len(team.Members)` Pokémon con datos detallados. Para un equipo vacío esto significa que solo se obtienen los primeros 6 Pokémon de la lista (Bulbasaur→Charizard), y `core.FillTeamRandom` no tiene un pool real del que elegir aleatoriamente.

## Causa raíz

```go
// app/bindings.go — líneas 459-468
var pokemon []core.Pokemon
for _, item := range list.Results {
    if len(pokemon) >= slotsNeeded+len(team.Members) {
        break // ← corta demasiado pronto, siempre toma los primeros N
    }
    p, ferr := a.fetcher.FetchPokemon(item.Name)
    if ferr == nil {
        pokemon = append(pokemon, p)
    }
}
```

El pool de candidatos debe ser mucho mayor que los slots necesarios para que el shuffle de Fisher-Yates en `core.FillTeamRandom` tenga variedad real.

## Capas afectadas

- **Core**: Sin cambios. La lógica de `FillTeamRandom` y su Fisher-Yates shuffle ya son correctos — el problema es que reciben un pool demasiado pequeño.
- **Shell**: Sin cambios.
- **APP**: Modificar `FillTeamRandom` en `app/bindings.go` para pasar todos los 151 Pokémon como pool de candidatos en lugar de cortar en `slotsNeeded`.
- **Frontend**: Sin cambios.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `app/bindings.go` | modificar | Eliminar el early-break del bucle de fetch de Pokémon para que se pasen todos los 151 como candidatos a `core.FillTeamRandom` |

## Plan de implementacion

1. **APP** — En `app/bindings.go`, función `FillTeamRandom`:
   - Eliminar la condición `if len(pokemon) >= slotsNeeded+len(team.Members) { break }`.
   - Iterar toda la lista de 151 Pokémon y obtener datos detallados de todos.
   - Pasar el pool completo a `core.FillTeamRandom`, que ya se encarga de hacer shuffle y seleccionar aleatoriamente.
   - **Nota sobre rendimiento**: si el fetch de 151 Pokémon es lento, considerar cachear los datos o usar una lista precargada. Verificar si ya existe cache en el fetcher.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Ejecutar "Rellenar aleatorio" varias veces consecutivas y verificar que genera equipos distintos |
| `core/teams_test.go` | Los tests existentes de `FillTeamRandom` ya validan la lógica de shuffle — no requieren cambios |

## Criterios de aceptacion

- [ ] Al pulsar "Rellenar aleatorio" varias veces (creando equipos nuevos), se obtienen Pokémon distintos cada vez
- [ ] El pool de candidatos incluye los 151 Pokémon disponibles, no solo los primeros N
- [ ] No se degrada significativamente el rendimiento (verificar que existe cache en el fetcher o implementar si es necesario)
- [ ] Los tests existentes siguen pasando
