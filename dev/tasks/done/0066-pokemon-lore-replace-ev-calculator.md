# Reemplazar calculador de EVs por lore/descripciones del Pokemon

**ID**: 0066-pokemon-lore-replace-ev-calculator
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Eliminar el calculador de EVs de la vista de detalle individual del Pokemon y reemplazarlo con una seccion de informacion narrativa: descripciones de la Pokedex, lore, categoria, habitat, flavor text, etc. Los datos se obtendran de la PokéAPI (endpoint `pokemon-species`) que ya tiene un binding parcial (`GetPokemonSpecies`), y se puede ampliar para traer mas campos.

## Capas afectadas

- **Core**: Nuevo tipo de dominio `PokemonLore` con campos para flavor text, genera, categoria, habitat, color, forma, descripciones, etc.
- **Shell**: Ampliar `pokeapi.go` para parsear los campos narrativos del endpoint `/pokemon-species/{name}` (flavor_text_entries, genera, habitat, color, shape, form_descriptions).
- **APP**: Nuevo binding/handler que exponga el lore. Eliminar binding `CalculateEVs` y handler POST `/api/calculator/evs` (o dejar deprecado).
- **Frontend**: Eliminar `ev-calculator.ts` y su invocacion en `pokedex.ts`. Crear nueva seccion de lore en la vista de detalle.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Agregar struct `PokemonLore` con campos: FlavorTexts, Genera (categoria), Habitat, Color, Shape, IsLegendary, IsMythical, FormDescriptions |
| `shell/pokeapi.go` | modificar | Ampliar `FetchPokemonSpecies` para mapear campos narrativos al nuevo `PokemonLore` |
| `app/bindings.go` | modificar | Agregar metodo `GetPokemonLore(name)` que llame a Shell y retorne `PokemonLore`. Eliminar o deprecar `CalculateEVs` |
| `app/mobile/handlers.go` | modificar | Agregar endpoint `GET /api/pokemon/{name}/lore`. Eliminar handler `POST /api/calculator/evs` |
| `frontend/src/ev-calculator.ts` | eliminar | Eliminar completamente el archivo del calculador de EVs |
| `frontend/src/pages/pokedex.ts` | modificar | Eliminar import y llamada a `renderEVCalculatorForm` / `initEVCalculator`. Agregar nueva seccion de lore en `showDetail()` |
| `frontend/src/api.ts` | modificar | Agregar funcion `GetPokemonLore(name)`. Eliminar `CalculateEVs` y `GetNatures` si ya no se usan |
| `core/ev_calc.go` | eliminar | Eliminar logica pura del calculador de EVs (ya no se usa) |
| `core/ev_calc_test.go` | eliminar | Eliminar tests del calculador (si existe) |

## Plan de implementacion

1. **Core**: Definir struct `PokemonLore` en `core/domain.go` con los campos narrativos necesarios.
2. **Shell**: Ampliar `FetchPokemonSpecies` en `shell/pokeapi.go` para parsear `flavor_text_entries` (filtrar por idioma ingles/espanol), `genera`, `habitat`, `color`, `shape`, `form_descriptions`.
3. **APP backend**: Crear binding `GetPokemonLore(name)` en `app/bindings.go` y endpoint REST `GET /api/pokemon/{name}/lore` en `app/mobile/handlers.go`.
4. **Frontend API**: Agregar `GetPokemonLore(name)` en `frontend/src/api.ts`.
5. **Frontend UI**: En `pokedex.ts`, eliminar la llamada a `renderEVCalculatorForm(p)` e `initEVCalculator(p)` en `showDetail()`. Reemplazar con nueva seccion que muestre flavor text, categoria, habitat y demas lore.
6. **Limpieza**: Eliminar `frontend/src/ev-calculator.ts`, `core/ev_calc.go`, y los bindings/handlers de `CalculateEVs`/`GetNatures` si ya no tienen uso.
7. **Probar** end-to-end en desktop (Wails) y mobile (REST).

## Fuente de datos

- **PokéAPI** endpoint `https://pokeapi.co/api/v2/pokemon-species/{name}` (gratuito, sin auth):
  - `flavor_text_entries[]`: descripciones de la Pokedex por version y idioma
  - `genera[]`: categoria del Pokemon (e.g., "Seed Pokemon") por idioma
  - `habitat`: habitat natural
  - `color`: color principal
  - `shape`: forma corporal
  - `is_legendary`, `is_mythical`: flags
  - `form_descriptions[]`: descripciones de formas alternativas

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/domain_test.go` | Validar que `PokemonLore` se inicializa correctamente |
| Tests manuales | Verificar que la vista de detalle muestra lore en lugar del calculador de EVs |

## Criterios de aceptacion

- [ ] El calculador de EVs ya no aparece en la vista de detalle del Pokemon
- [ ] Los archivos `ev-calculator.ts` y `core/ev_calc.go` han sido eliminados
- [ ] La vista de detalle muestra flavor text (descripcion de Pokedex) del Pokemon
- [ ] Se muestra la categoria del Pokemon (e.g., "Seed Pokemon")
- [ ] Se muestra habitat, color y/o forma si estan disponibles
- [ ] Se indica si el Pokemon es legendario o mitico
- [ ] Funciona tanto en modo desktop (Wails) como mobile (REST API)
- [ ] Core no importa ningun paquete externo
- [ ] La seccion de lore tiene un diseno coherente con el resto de la app

## Notas

- PokéAPI `pokemon-species` ya se consume parcialmente (`GetPokemonSpecies` en bindings). Se amplia el parsing para extraer campos narrativos.
- Filtrar `flavor_text_entries` por idioma: preferir `es` (espanol) con fallback a `en` (ingles), o mostrar ambos.
- Los flavor texts vienen con `\n` y `\f` que hay que limpiar antes de mostrar.
- Considerar cachear el lore ya que cambia poco y la API tiene rate limits.
