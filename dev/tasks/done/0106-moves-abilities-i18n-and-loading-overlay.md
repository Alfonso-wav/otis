# Traducir nombres de movimientos y habilidades + overlay de carga con Machamp

**ID**: 0106-moves-abilities-i18n-and-loading-overlay
**Estado**: done
**Fecha**: 2026-04-04

---

## Descripcion

Los nombres de movimientos y habilidades (y las descripciones de habilidades) se muestran siempre en inglés, incluso cuando el idioma seleccionado es español. Se necesita obtener las traducciones desde la PokeAPI (campo `names` y `flavor_text_entries` en español) y mostrarlas en el frontend según el locale activo. Además, al abrir las pestañas de movimientos o habilidades debe aparecer un overlay de carga con el Pokémon **Machamp** (igual que Mr. Mime aparece en el overlay de ordenación).

## Capas afectadas

- **Core**: Añadir campo `NameEs` y `DescriptionEs` a los structs `Move` y `Ability` en `domain.go`.
- **Shell**: Modificar `FetchMove()` y `FetchAbility()` en `pokeapi.go` para extraer los nombres y descripciones en español desde la respuesta de PokeAPI (campo `names` con `language.name == "es"` y `flavor_text_entries` con `language.name == "es"`).
- **APP**: Sin cambios significativos (los bindings ya exponen los structs).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir campos `NameEs` y `DescriptionEs` a `Move` y `Ability` |
| `shell/pokeapi.go` | modificar | Extraer nombre español de `names[]` y descripción española de `flavor_text_entries[]` al parsear movimientos y habilidades |
| `frontend/src/pages/explore/moves.ts` | modificar | Usar nombre traducido según locale; mostrar overlay Machamp durante carga inicial |
| `frontend/src/pages/explore/abilities.ts` | modificar | Usar nombre/descripción traducidos según locale; mostrar overlay Machamp durante carga inicial |
| `frontend/src/components/sorting-overlay.ts` | modificar | Añadir soporte para sprite configurable (Machamp) o crear función de overlay de carga separada |
| `frontend/src/locales/es.json` | modificar | Añadir claves de texto para overlay de carga de movimientos/habilidades si es necesario |
| `frontend/src/locales/en.json` | modificar | Añadir claves equivalentes en inglés |

## Plan de implementacion

1. **Core** — Añadir `NameEs string` y `DescriptionEs string` a los structs `Move` y `Ability` en `core/domain.go`.
2. **Shell (moves)** — En `shell/pokeapi.go` → `FetchMove()`: parsear el array `names` de la respuesta PokeAPI para extraer el nombre con `language.name == "es"`. Parsear `flavor_text_entries` con `language.name == "es"` para la descripción en español.
3. **Shell (abilities)** — En `shell/pokeapi.go` → `FetchAbility()`: mismo patrón — extraer nombre y descripción en español del array `names` y `flavor_text_entries`.
4. **Frontend (overlay)** — Modificar `sorting-overlay.ts` para aceptar un sprite personalizable, o crear una función `showLoadingOverlay(sprite, text)` que use Machamp (`pokemon/other/official-artwork/68.png`).
5. **Frontend (moves)** — En `moves.ts`: mostrar overlay Machamp al inicio de `initMoves()` mientras se cargan los datos. En `renderTable()`, usar `m.NameEs` o `m.Name` según el locale activo.
6. **Frontend (abilities)** — En `abilities.ts`: mostrar overlay Machamp al inicio de `initAbilities()` mientras se cargan los datos. En `renderTable()`, usar `a.NameEs`/`a.DescriptionEs` o los campos ingleses según el locale.
7. **Filtro de búsqueda** — Asegurar que la búsqueda funcione tanto con el nombre original como con el traducido.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/domain_test.go` | Verificar que los structs Move y Ability incluyen los campos de traducción |
| Test manual | Cambiar idioma a ES y verificar que nombres/descripciones aparecen en español |
| Test manual | Verificar overlay de Machamp al abrir pestañas de movimientos y habilidades |

## Criterios de aceptacion

- [x] Los nombres de movimientos aparecen en español cuando el locale es "es"
- [x] Los nombres de habilidades aparecen en español cuando el locale es "es"
- [x] Las descripciones de habilidades aparecen en español cuando el locale es "es"
- [x] Al abrir la pestaña de movimientos se muestra un overlay con Machamp mientras cargan los datos
- [x] Al abrir la pestaña de habilidades se muestra un overlay con Machamp mientras cargan los datos
- [x] La búsqueda de movimientos funciona con nombres en español
- [x] La búsqueda de habilidades funciona con nombres en español
- [x] En inglés todo sigue funcionando igual que antes
- [x] El overlay de ordenación (Mr. Mime) no se ve afectado

## Notas

- PokeAPI expone traducciones en el campo `names[]` con `{ name: "...", language: { name: "es" } }` tanto para movimientos como habilidades.
- Machamp es el Pokémon #68, sprite oficial: `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/68.png`
- El overlay de carga de Machamp es independiente del overlay de ordenación de Mr. Mime — se puede parametrizar `showSortingOverlay` o crear una función nueva.
- Las descripciones de movimientos (`flavor_text_entries` en español) también podrían traducirse, pero el usuario no las pidió explícitamente — solo nombres de movimientos y nombres+descripciones de habilidades.
