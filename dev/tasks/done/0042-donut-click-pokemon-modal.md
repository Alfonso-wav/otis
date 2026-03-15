# Modal de Pokémon al hacer click en segmento del donut de distribución de tipos

**ID**: 0042-donut-click-pokemon-modal
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

En la sección Explorar > Regiones, al expandir una tarjeta de región se muestra un donut chart (ECharts) con la distribución de tipos. Actualmente es solo visual, sin interacción.

Se quiere que al hacer click en un segmento del donut (ej: "ghost: 3 Pokémon"), se abra una ventana emergente (modal) mostrando los Pokémon concretos de ese tipo en esa región, con su sprite y nombre.

## Capas afectadas

- **Core**: Nueva función pura `FilterPokedexByType` que cruza las entradas de un Pokédex regional con los Pokémon de un tipo dado.
- **Shell**: Sin cambios (ya existen `FetchPokedex`, `FetchType`).
- **APP**: Nuevo binding `GetRegionPokemonByType(region, typeName string)` que orquesta Shell para obtener los Pokémon de un tipo en una región.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/pokemon.go` | modificar | Añadir `FilterPokedexByType(entries []string, typePokemon []string) []string` |
| `core/pokemon_test.go` | modificar | Test para `FilterPokedexByType` |
| `app/bindings.go` | modificar | Añadir binding `GetRegionPokemonByType(region, typeName string)` |
| `frontend/src/charts/type-distribution.ts` | modificar | Añadir evento click en segmentos del pie chart, callback que invoca al modal |
| `frontend/src/pages/explore/regions.ts` | modificar | Importar y usar el modal, pasar callback al chart |
| `frontend/src/components/pokemon-type-modal.ts` | crear | Componente modal que muestra lista de Pokémon con sprite y nombre |
| `frontend/src/styles/_explore.scss` | modificar | Estilos del modal de Pokémon por tipo |

## Plan de implementacion

1. **Core**: Crear `FilterPokedexByType` — función pura que recibe una lista de nombres del pokédex regional y una lista de nombres del tipo, y devuelve la intersección (los que están en ambas listas).
2. **APP**: Crear binding `GetRegionPokemonByType(region, typeName)`:
   - Obtener el pokédex regional correspondiente (mapear región → pokédex name, ej: kanto → kanto).
   - Obtener el detalle del tipo con `FetchType(typeName)` para tener la lista de Pokémon de ese tipo.
   - Cruzar ambas listas con `FilterPokedexByType`.
   - Retornar la lista de nombres resultante.
3. **Frontend — Modal**: Crear `pokemon-type-modal.ts` con funciones `openTypeModal(regionName, typeName, pokemonNames)` y `closeTypeModal()`. Renderizar un overlay con un grid de Pokémon (sprite + nombre). Reusar el patrón de modal existente en `builds.ts`.
4. **Frontend — Chart click**: En `renderTypeDistributionChart`, añadir `chart.on('click', callback)` que reciba el nombre del tipo clickeado y dispare la carga de datos + apertura del modal.
5. **Frontend — Integración**: En `regions.ts`, pasar el nombre de la región al chart para que el callback pueda invocar `GetRegionPokemonByType` y luego `openTypeModal`.
6. **Estilos**: Añadir estilos del modal en `_explore.scss` (overlay, grid de Pokémon, responsive).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/pokemon_test.go` | `FilterPokedexByType` con intersección parcial, vacía y completa |
| `app/bindings_test.go` | `GetRegionPokemonByType` con mock fetcher |

## Criterios de aceptacion

- [ ] Al hacer click en un segmento del donut se abre un modal
- [ ] El modal muestra los Pokémon del tipo seleccionado en esa región
- [ ] Cada Pokémon muestra su sprite (local con fallback CDN) y nombre
- [ ] El modal se cierra con click en X, click en overlay, o tecla Escape
- [ ] Funciona para todas las regiones que tienen datos de distribución
- [ ] Si no hay Pokémon (error o lista vacía), se muestra un mensaje apropiado
- [ ] Responsive: el modal se adapta a pantallas pequeñas

## Notas

- Los datos de distribución en `REGION_TYPE_DATA` son hardcoded e ilustrativos. Los Pokémon reales vendrán de cruzar el pokédex regional (PokeAPI) con el endpoint de tipo. Los conteos pueden no coincidir exactamente con los datos hardcoded.
- Mapeo región → pokédex: kanto→kanto, johto→original-johto, hoenn→hoenn, sinnoh→original-sinnoh, unova→original-unova, kalos→kalos-central. Verificar nombres exactos en la API.
- Ya existe un patrón de modal en `builds.ts` (save modal) que se puede reusar como referencia de estilo y comportamiento.
- Los sprites locales están servidos por el proyecto (tarea 0033). Usar la misma lógica de fallback a CDN.
