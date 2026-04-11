# Fix porcentajes de encuentro que superan 100% en localizaciones

**ID**: 0140-fix-location-percentages
**Estado**: todo
**Fecha**: 2026-04-11

---

## Descripcion

En Regiones > Localizaciones, al clickar una localizacion, el modal de encuentros muestra porcentajes (`MaxChance`) que pueden superar el 100%. Esto ocurre porque PokeAPI devuelve valores `max_chance` que no siempre estan acotados a 100, y la app los pasa directamente sin validacion ni normalizacion en ninguna capa.

**Flujo actual del bug:**
1. `shell/pokeapi_locations.go:66-70` — extrae `max_chance` de PokeAPI sin validar
2. `core/logic_encounters.go:9-11` — agrega tomando el maximo, sin acotar
3. `app/bindings.go:423` — pass-through
4. `frontend/src/components/location-encounter-modal.ts:101` — muestra `${enc.MaxChance}%` directamente

## Capas afectadas

- **Core**: `core/logic_encounters.go` — normalizar porcentajes
- **Shell**: ninguna (los datos de PokeAPI son lo que son)
- **APP**: ninguna (la normalizacion debe estar en Core)

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/logic_encounters.go` | modificar | Acotar MaxChance a rango 0-100 en AggregateEncounters |
| `core/logic_encounters_test.go` | modificar | Agregar test para valores > 100 |

## Plan de implementacion

1. En `core/logic_encounters.go`, dentro de `AggregateEncounters()`, despues de calcular el MaxChance maximo para cada Pokemon, acotar el valor a 100 si supera ese limite: `if maxChance > 100 { maxChance = 100 }`.
2. Agregar tests que verifiquen que valores > 100 se normalizan a 100.
3. Verificar visualmente en el frontend que los porcentajes ahora son <= 100%.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/logic_encounters_test.go` | Que MaxChance nunca supere 100 despues de AggregateEncounters |
| `core/logic_encounters_test.go` | Que valores normales (< 100) no se modifican |

## Criterios de aceptacion

- [ ] Ningun porcentaje en el modal de encuentros supera 100%
- [ ] Valores normales (< 100) no se alteran
- [ ] La logica de normalizacion esta en Core (funcion pura)
- [ ] Tests unitarios cubren el caso de valores > 100
- [ ] La correccion aplica tanto para Wails como para la API REST (APK)

## Notas

La correccion es simple y va en Core porque es logica pura de dominio. No se toca Shell ni frontend. El cap a 100 se aplica al resultado final de la agregacion, no a los datos crudos de PokeAPI.
