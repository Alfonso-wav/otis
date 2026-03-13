# 0010 — Exhaustive PokéAPI Coverage

## Objetivo
Ampliar las conexiones a PokéAPI para cubrir **todos los endpoints públicos relevantes** de la API. Actualmente el proyecto cubre ~9 endpoints; la API expone más de 50 recursos. Esta tarea los implementa por grupos de dominio siguiendo la arquitectura Core / Shell / APP.

## Estado actual
Ya implementado en `shell/pokeapi.go` + `core/ports.go`:
- `GET /pokemon/{name}` → `FetchPokemon`
- `GET /pokemon` → `FetchPokemonList`
- `GET /type` → `FetchTypeList`
- `GET /type/{name}` → `FetchType`
- `GET /region` → `FetchRegions`
- `GET /region/{name}` → `FetchRegion`
- `GET /move/{name}` → `FetchMove`
- `GET /ability/{name}` → `FetchAbility`
- `GET /evolution-chain/{id}` → `FetchEvolutionChain`

## Grupos de trabajo

### Grupo A — Pokémon extendido (alta prioridad)
| Endpoint | Shell method | Descripción |
|---|---|---|
| `GET /pokemon-species/{name}` | `FetchPokemonSpecies` | Flavor text, gender rate, capture rate, egg groups, hatch counter, baby/legendary/mythical, genus |
| `GET /pokemon-form/{name}` | `FetchPokemonForm` | Formas alternativas: mega, regional, gigantamax |
| `GET /pokemon-color` | `FetchPokemonColors` | Lista de colores Pokédex |
| `GET /pokemon-shape` | `FetchPokemonShapes` | Lista de formas Pokédex |
| `GET /pokemon-habitat` | `FetchPokemonHabitats` | Lista de hábitats |

### Grupo B — Naturalezas y cría (alta prioridad)
| Endpoint | Shell method | Descripción |
|---|---|---|
| `GET /nature` | `FetchNatureList` | Lista de naturalezas (reemplaza el hardcode en `core/ev_calc.go`) |
| `GET /nature/{name}` | `FetchNature` | Detalle: stat aumentado, stat disminuido, preferencias de bayas |
| `GET /egg-group/{name}` | `FetchEggGroup` | Pokémon en el grupo de huevo |
| `GET /gender/{name}` | `FetchGender` | Pokémon con ese género |
| `GET /growth-rate/{name}` | `FetchGrowthRate` | Curva de experiencia, niveles, Pokémon |
| `GET /characteristic/{id}` | `FetchCharacteristic` | Característica de IV alto |

### Grupo C — Movimientos completos (alta prioridad)
| Endpoint | Shell method | Descripción |
|---|---|---|
| `GET /move` | `FetchMoveList` | Lista paginada de todos los movimientos |
| `GET /move-damage-class/{name}` | `FetchMoveDamageClass` | physical / special / status |
| `GET /move-ailment/{name}` | `FetchMoveAilment` | Estado alterado que causa |
| `GET /move-target/{name}` | `FetchMoveTarget` | A quién afecta el movimiento |
| `GET /machine/{id}` | `FetchMachine` | TM/HM: qué movimiento enseña |

### Grupo D — Habilidades completas (media prioridad)
| Endpoint | Shell method | Descripción |
|---|---|---|
| `GET /ability` | `FetchAbilityList` | Lista paginada de todas las habilidades |

### Grupo E — Items (media prioridad)
| Endpoint | Shell method | Descripción |
|---|---|---|
| `GET /item` | `FetchItemList` | Lista paginada de items |
| `GET /item/{name}` | `FetchItem` | Detalle: efecto, coste, categoría, máquina |
| `GET /item-category/{name}` | `FetchItemCategory` | Items en la categoría |
| `GET /item-pocket/{name}` | `FetchItemPocket` | Bolsillo de la mochila |

### Grupo F — Ubicaciones extendidas (media prioridad)
| Endpoint | Shell method | Descripción |
|---|---|---|
| `GET /location/{name}` | `FetchLocation` | Ubicación con sus áreas |
| `GET /location-area/{name}` | `FetchLocationArea` | Pokémon que aparecen en el área |
| `GET /pal-park-area/{name}` | `FetchPalParkArea` | Áreas del Pal Park |

### Grupo G — Stats y generaciones (media prioridad)
| Endpoint | Shell method | Descripción |
|---|---|---|
| `GET /stat/{name}` | `FetchStat` | Detalles de un stat (qué Pokémon lo tienen alto) |
| `GET /generation/{id}` | `FetchGeneration` | Juegos, tipos y Pokémon introducidos |
| `GET /pokedex/{name}` | `FetchPokedex` | Pokédex regional con entradas |
| `GET /version-group/{name}` | `FetchVersionGroup` | Juegos agrupados |

### Grupo H — Bayas (baja prioridad)
| Endpoint | Shell method | Descripción |
|---|---|---|
| `GET /berry` | `FetchBerryList` | Lista de bayas |
| `GET /berry/{name}` | `FetchBerry` | Detalle: firmeza, sabores, efecto natural |
| `GET /berry-flavor/{name}` | `FetchBerryFlavor` | Sabor con sus bayas y naturalezas relacionadas |

## Plan de implementación

### Paso 1 — Tipos de dominio en `core/domain.go`
Añadir structs para cada nuevo recurso:
- `PokemonSpecies`, `PokemonForm`, `NatureDetail`, `EggGroup`, `GrowthRate`
- `MoveList`, `MoveDamageClass`, `MoveAilment`, `MoveTarget`, `Machine`
- `AbilityListResponse`
- `Item`, `ItemListResponse`, `ItemCategory`, `ItemPocket`
- `LocationDetail`, `LocationArea`, `PalParkArea`
- `StatDetail`, `Generation`, `Pokedex`, `VersionGroup`
- `Berry`, `BerryListResponse`, `BerryFlavor`

### Paso 2 — Ampliar `core/ports.go`
Añadir todos los métodos nuevos a la interfaz `PokemonFetcher`.

### Paso 3 — Implementar en `shell/pokeapi.go`
Por cada endpoint:
1. Definir struct `api*` raw con los campos JSON relevantes.
2. Implementar método `Fetch*` en `PokeAPIClient`.
3. Escribir función `toDomain*` pura que convierte raw → domain.

### Paso 4 — Exponer en `app/bindings.go`
Añadir métodos públicos en `App` para cada nuevo Fetch.

### Paso 5 — Generar bindings de Wails
Ejecutar `wails generate bindings` o rebuild para que el frontend pueda llamar a los nuevos métodos.

### Paso 6 — Tests
Añadir tests de las funciones `toDomain*` en `core/` (son puras, sin mocks).

## Restricciones
- Los structs `api*` solo viven en `shell/`. Los tipos de dominio en `core/`.
- Ninguna función de `core/` hace HTTP.
- El cliente HTTP existente en `PokeAPIClient` se reutiliza para todos los endpoints.
- Respetar el timeout de 10s existente; para endpoints que devuelven listas grandes (moves, abilities) usar `limit=10000` para obtener todo en una sola llamada.
- Los campos opcionales de la API (`*int`, `*string`) se convierten a zero values en el dominio salvo que el negocio necesite distinguir null de 0.

## Archivos afectados
- `core/domain.go` — nuevos tipos
- `core/ports.go` — interfaz ampliada
- `shell/pokeapi.go` — implementaciones
- `app/bindings.go` — métodos expuestos
- `core/domain_test.go` (nuevo) — tests de conversión

## Notas
- Esta tarea es grande. Se recomienda implementar por grupos (A → H) y hacer commit al terminar cada grupo.
- La cobertura de Grupos A–D es suficiente para las vistas actuales del frontend; E–H son opcionales para futuras vistas.
