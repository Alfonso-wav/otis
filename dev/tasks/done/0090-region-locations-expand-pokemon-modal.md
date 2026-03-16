# Expandir localizaciones de región y modal de Pokémon por localización

**ID**: 0090-region-locations-expand-pokemon-modal
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

En Explorar > Regiones, cada tarjeta de región muestra las primeras 20 localizaciones como tags y un indicador estático "+X más". Se necesitan dos funcionalidades nuevas:

1. **Expandir todas las localizaciones**: Al pulsar en el tag "+X más", se despliegan todas las localizaciones restantes con animación (actualmente es un indicador estático sin interacción).
2. **Modal de Pokémon por localización**: Al hacer clic en cualquier tag de localización, se abre una ventana/modal que muestra los Pokémon que se pueden encontrar en esa localización (usando los datos de encounter de la PokeAPI).

## Capas afectadas

- **Core**: Añadir función pura para agregar encounters de múltiples áreas de una localización en una sola lista.
- **Shell**: Ya existe `FetchLocation` y `FetchLocationArea` en `pokeapi_locations.go`. Puede necesitar una función que agregue encounters de todas las áreas de una localización.
- **APP (backend)**: Exponer un nuevo endpoint o binding que devuelva los Pokémon de una localización concreta agregando sus áreas.
- **APP (frontend)**: `regions.ts` — hacer clickable el tag "+X más" para expandir, hacer clickable cada tag de localización para abrir modal con encounters.
- **Frontend styles**: `_explore.scss` — cursor pointer en tags, estilos del modal de encounters, animación de expansión.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir tipo `LocationEncounterSummary` si es necesario para agregar encounters |
| `core/logic.go` | modificar | Función pura para agregar/deduplicar encounters de múltiples áreas |
| `shell/pokeapi_locations.go` | modificar | Añadir método `FetchLocationEncounters(name)` que obtenga todas las áreas de una localización y agregue sus encounters |
| `app/bindings.go` | modificar | Exponer `GetLocationEncounters(name)` como binding para el frontend |
| `app/mobile/handlers.go` | modificar | Añadir endpoint REST `GET /api/locations/{name}/encounters` |
| `frontend/src/api.ts` | modificar | Añadir función `GetLocationEncounters(name)` |
| `frontend/src/pages/explore/regions.ts` | modificar | Implementar expand de tags y modal de encounters por localización |
| `frontend/src/styles/_explore.scss` | modificar | Estilos para tags clickables, modal de encounters, animación de expansión |

## Plan de implementacion

### Parte 1: Backend — Endpoint de encounters por localización

1. En `shell/pokeapi_locations.go`, crear `FetchLocationEncounters(name string)` que:
   - Llame a `FetchLocation(name)` para obtener las áreas de la localización.
   - Para cada área, llame a `FetchLocationArea(areaName)` para obtener los encounters.
   - Agregue todos los encounters en una sola lista (deduplicando por nombre de Pokémon, quedándose con la mayor `MaxChance`).
2. En `core/logic.go`, añadir función pura `AggregateEncounters(areas []LocationArea) []PokemonEncounter` que haga la agregación/deduplicación.
3. En `app/bindings.go`, exponer `GetLocationEncounters(name)`.
4. En `app/mobile/handlers.go`, añadir handler para `GET /api/locations/{name}/encounters`.
5. En `frontend/src/api.ts`, añadir `GetLocationEncounters(name)`.

### Parte 2: Frontend — Expandir "+X más"

6. En `regions.ts`, cambiar el tag `.region-location-more` para que sea clickable (`cursor: pointer`).
7. Al hacer clic en "+X más":
   - Renderizar las localizaciones restantes (las que estaban ocultas tras el límite de 20) con animación stagger (GSAP).
   - Ocultar o reemplazar el tag "+X más" por un tag "Mostrar menos" que colapse de vuelta a 20.
8. Almacenar el estado expandido/colapsado por región.

### Parte 3: Frontend — Modal de Pokémon por localización

9. Hacer cada `.region-location-tag` clickable (`cursor: pointer`).
10. Al hacer clic en un tag de localización:
    - Mostrar un indicador de carga (spinner) mientras se obtienen los datos.
    - Llamar a `GetLocationEncounters(locationName)` para obtener los Pokémon.
    - Renderizar un modal/popup con la lista de Pokémon encontrados (nombre, sprite si disponible, probabilidad de encuentro).
11. Reutilizar el estilo de modal existente en el proyecto (similar al modal de Pokémon por tipo en el chart de distribución).
12. Animación de entrada/salida del modal con GSAP.

### Parte 4: Estilos y polish

13. En `_explore.scss`:
    - `cursor: pointer` para `.region-location-tag` y `.region-location-more`.
    - Hover effect en los tags (ligero cambio de color/sombra).
    - Estilos del modal de encounters (reutilizar patrones existentes).
    - Transición suave para la expansión/colapso de tags.
14. Asegurar compatibilidad con dark mode.
15. Asegurar responsiveness en mobile.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/logic_test.go` | `AggregateEncounters` deduplica correctamente y mantiene la mayor `MaxChance` |
| (visual) | Al hacer clic en "+X más" se despliegan todas las localizaciones con animación |
| (visual) | Aparece "Mostrar menos" tras expandir y al pulsarlo se colapsa a 20 |
| (visual) | Al hacer clic en un tag de localización aparece modal con Pokémon |
| (visual) | El modal muestra correctamente nombre y probabilidad de cada Pokémon |
| (visual) | Loading spinner aparece mientras se cargan los encounters |
| (visual) | Si una localización no tiene encounters, el modal muestra un mensaje adecuado |
| (visual) | Funciona correctamente en dark mode |
| (visual) | Funciona correctamente en mobile (responsive) |

## Criterios de aceptacion

- [x] Al pulsar "+X más" se despliegan todas las localizaciones ocultas con animación
- [x] Aparece opción de colapsar tras expandir
- [x] Al hacer clic en cualquier tag de localización se abre un modal con los Pokémon de esa localización
- [x] El modal muestra nombre de Pokémon y probabilidad de encuentro
- [x] Se muestra loading mientras se obtienen datos de encounters
- [x] Si no hay encounters, se muestra mensaje informativo
- [x] Compatible con dark mode (N/A — app no tiene dark mode)
- [x] Compatible con mobile (responsive) — reutiliza estilos .type-modal responsive existentes
- [x] El endpoint/binding de encounters agrega correctamente las áreas de una localización

## Notas

- El backend ya tiene `FetchLocation` (devuelve áreas) y `FetchLocationArea` (devuelve encounters por área). Solo falta componer ambas llamadas.
- Actualmente se muestran solo 20 de N localizaciones y el tag "+X más" es estático (sin evento click).
- El modal de Pokémon por tipo en el chart de distribución de regiones puede servir como referencia de diseño.
- Las llamadas a la PokeAPI para encounters pueden ser lentas (una petición por área), considerar mostrar loading.
