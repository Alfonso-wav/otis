# 0073 — Mostrar ubicaciones y encuentros del Pokémon en la vista de detalle

## Descripción

Añadir una nueva sección en la vista de detalle del Pokémon (al hacer clic en uno en la Pokédex) que muestre dónde se puede encontrar ese Pokémon en las diferentes regiones: ubicaciones, tasa de aparición, método de encuentro (hierba, pesca, surf…), y condiciones relevantes (día/noche, versión del juego, etc.).

## Capas afectadas

- **Core** — nuevos tipos de dominio para encuentros por Pokémon
- **Shell** — nuevo fetcher para el endpoint `pokemon/{name}/encounters` de PokeAPI
- **APP** — nuevo binding + wiring; nuevo endpoint REST para móvil
- **APP (frontend)** — nueva sección en la vista de detalle + función en api.ts

## Archivos a crear / modificar

### Core

1. **`core/domain.go`** — Añadir tipos nuevos:
   ```go
   type EncounterCondition struct {
       Name   string // e.g. "time-night", "season-winter"
       Values []string
   }

   type EncounterMethodDetail struct {
       Method      string // e.g. "walk", "surf", "old-rod"
       Chance      int    // probabilidad (max_chance)
       MinLevel    int
       MaxLevel    int
       Conditions  []EncounterCondition
   }

   type VersionEncounter struct {
       Version    string // e.g. "red", "platinum"
       MaxChance  int
       Details    []EncounterMethodDetail
   }

   type PokemonLocationEncounter struct {
       LocationArea string
       Region       string // extraído del nombre o con lookup adicional
       Versions     []VersionEncounter
   }
   ```

### Shell

2. **`shell/pokeapi_encounters.go`** (nuevo) — Implementar `FetchPokemonEncounters(name string) ([]core.PokemonLocationEncounter, error)` consumiendo `GET /pokemon/{name}/encounters` de PokeAPI. Este endpoint devuelve un array con:
   - `location_area.name`
   - `version_details[].version.name`
   - `version_details[].max_chance`
   - `version_details[].encounter_details[].method.name`
   - `version_details[].encounter_details[].chance`
   - `version_details[].encounter_details[].min_level`, `max_level`
   - `version_details[].encounter_details[].condition_values[].name`

### APP (backend)

3. **`app/bindings.go`** — Añadir método:
   ```go
   func (a *App) GetPokemonEncounters(name string) ([]core.PokemonLocationEncounter, error) {
       return a.fetcher.FetchPokemonEncounters(core.NormalizeName(name))
   }
   ```

4. **`app/mobile/handlers.go`** (o equivalente) — Añadir endpoint REST `GET /pokemon/{name}/encounters` para móvil.

### APP (frontend)

5. **`frontend/src/api.ts`** — Añadir función `GetPokemonEncounters(name)` con el patrón dual (Wails IPC / HTTP).

6. **`frontend/src/pages/pokedex.ts`** — Añadir función `loadEncounters(name)` similar a `loadLore()`:
   - Se invoca desde `renderDetail()` después de `loadLore()`.
   - Renderiza una sección `#pokemon-encounters` debajo de Lore.
   - Agrupa las ubicaciones por región (si es posible inferir la región del nombre del location area).
   - Muestra para cada ubicación: nombre del área, versiones del juego, método de encuentro, probabilidad, niveles, y condiciones (día/noche, etc.).
   - Si no hay encuentros (Pokémon legendario, starter, regalo, etc.), mostrar mensaje tipo "This Pokémon is not found in the wild."

### Estilos

7. **`frontend/src/styles/_pokedex.scss`** — Estilos para la sección de encuentros (tabla o lista colapsable por región).
8. **`frontend/src/styles/_dark.scss`** — Override dark mode para la nueva sección.

## Requisitos

1. La sección debe cargarse de forma asíncrona (como Lore) para no bloquear el renderizado inicial.
2. Mostrar un spinner/loading mientras se cargan los datos.
3. Agrupar encuentros por región o al menos por location area.
4. Mostrar información de versión del juego, método (caminar, surfear, pescar, etc.), probabilidad (%), rango de niveles, y condiciones (día/noche, estación).
5. Si el Pokémon no tiene encuentros salvajes, mostrar un mensaje informativo.
6. Responsive: debe verse bien en desktop y en móvil.
7. Compatible con dark mode.

## Endpoint de referencia (PokeAPI)

```
GET https://pokeapi.co/api/v2/pokemon/{id or name}/encounters
```

Devuelve un array de objetos con esta estructura:
```json
[
  {
    "location_area": { "name": "cerulean-city-area", "url": "..." },
    "version_details": [
      {
        "version": { "name": "red", "url": "..." },
        "max_chance": 100,
        "encounter_details": [
          {
            "method": { "name": "gift", "url": "..." },
            "chance": 100,
            "min_level": 10,
            "max_level": 10,
            "condition_values": []
          }
        ]
      }
    ]
  }
]
```

## Tests

- Verificar que la sección de encuentros aparece bajo Lore en la vista de detalle.
- Verificar que Pokémon con encuentros salvajes (e.g., Pidgey) muestran ubicaciones correctas.
- Verificar que Pokémon sin encuentros salvajes (e.g., starters, legendarios con gift) muestran mensaje adecuado.
- Verificar carga asíncrona (spinner visible brevemente).
- Verificar dark mode.
- Verificar layout responsive en móvil.
- Verificar que funciona tanto en desktop (Wails) como en Android (HTTP).
