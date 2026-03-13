# Vista de Pokémon agrupados por tipo

**ID**: 0007-pokemon-type-grouping-view
**Estado**: done
**Fecha**: 2026-03-13
**Depende de**: 0006-tab-navigation-system
**Agente**: .claude/agents/frontend

---

## Descripcion

Implementar la pestaña "Por Tipo" que muestra los Pokémon agrupados según su tipo (fuego, agua, planta, etc.). Requiere un nuevo endpoint en el backend para obtener la lista de tipos y los Pokémon de cada tipo desde PokeAPI. La vista mostrará secciones colapsables o cards por tipo, cada una con los Pokémon que pertenecen a ese tipo, usando los colores de tipo ya definidos en el SCSS.

## Capas afectadas

- **Core**: nuevos tipos de dominio para `PokemonTypeDetail` y `TypeListResponse`
- **Shell**: nuevo método en `PokeAPIClient` para fetch de tipos (`/type` y `/type/{id}`)
- **APP**: nuevos bindings para exponer los endpoints de tipos al frontend
- **Frontend**: implementar la vista completa en `pages/types.ts`

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir structs: `PokemonTypeDetail`, `TypePokemonEntry`, `TypeListResponse` |
| `core/ports.go` | modificar | Añadir métodos a `PokemonFetcher`: `FetchTypeList()`, `FetchType(name)` |
| `core/pokemon.go` | modificar | Funciones puras para agrupar/ordenar Pokémon por tipo si se necesitan |
| `shell/pokeapi.go` | modificar | Implementar `FetchTypeList()` y `FetchType(name)` contra PokeAPI |
| `app/bindings.go` | modificar | Nuevos métodos Wails: `ListTypes()`, `GetType(name)` |
| `frontend/src/pages/types.ts` | modificar | Implementar vista completa: grid de tipos, expansión con Pokémon de cada tipo |
| `frontend/src/styles/_types.scss` | crear | Estilos específicos para la vista de tipos (cards de tipo, grid de Pokémon dentro) |
| `frontend/src/styles/main.scss` | modificar | Importar `_types.scss` |

## Plan de implementacion

1. Añadir tipos de dominio en `core/domain.go`: `PokemonTypeDetail` (nombre, pokémon que lo tienen), `TypeListResponse`
2. Extender interfaz `PokemonFetcher` en `core/ports.go` con los nuevos métodos
3. Implementar en `shell/pokeapi.go` las llamadas a `GET /type` (lista) y `GET /type/{id}` (detalle con pokémon)
4. Añadir bindings en `app/bindings.go`: `ListTypes()` y `GetType(name string)`
5. Regenerar bindings Wails (`wails generate module`)
6. Implementar la vista en `pages/types.ts`:
   - Al entrar al tab, cargar la lista de tipos
   - Mostrar cards/badges por tipo con sus colores
   - Al hacer click en un tipo, expandir y mostrar los Pokémon de ese tipo (lazy load)
   - Cada Pokémon muestra sprite + nombre (reutilizar diseño de cards del grid)
7. Añadir animaciones GSAP: stagger en la aparición de tipos, expansión animada
8. Estilos con personalidad "Aura": cards de tipo con bordes redondeados, sombras sutiles, colores de tipo como fondo

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/pokemon_test.go` | Funciones puras de agrupación/filtrado por tipo |
| Manual | Vista carga todos los tipos disponibles |
| Manual | Click en tipo muestra los Pokémon correctos |
| Manual | Sprites y nombres se muestran correctamente |
| `npm run build` | Build frontend compila sin errores |
| `go build ./...` | Build backend compila sin errores |

## Criterios de aceptacion

- [ ] La pestaña "Por Tipo" muestra todos los tipos de Pokémon disponibles
- [ ] Cada tipo tiene su color correspondiente (usando los colores ya definidos en SCSS)
- [ ] Click en un tipo muestra los Pokémon de ese tipo con sprite + nombre
- [ ] Carga lazy: los Pokémon de un tipo se cargan al expandir, no todos de golpe
- [ ] Animaciones fluidas al expandir/colapsar tipos y al aparecer Pokémon
- [ ] Se puede navegar de esta vista al detalle de un Pokémon individual
- [ ] Diseño responsive y coherente con la estética "Aura"

## Notas

- PokeAPI endpoints: `GET /type` devuelve ~20 tipos, `GET /type/{id}` devuelve los Pokémon de ese tipo.
- Algunos tipos tienen cientos de Pokémon — considerar paginación o límite dentro de cada tipo.
- Reutilizar la card de Pokémon de la vista principal para mantener consistencia visual.
- Los colores de tipo ya están definidos en `_variables.scss` — reutilizarlos.
