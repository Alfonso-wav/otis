# Task 0130 — Regiones: datos de Pokemon vacios en Gen VII+

## Estado: pending

## Goal

A partir de la Gen VII (Alola, Galar, Hisui, Paldea), las localizaciones dentro de las regiones no muestran datos de Pokemon al hacer click, y la grafica de distribucion de tipos aparece vacia. Diagnosticar y corregir la causa raiz, o en su defecto mostrar un mensaje informativo al usuario.

---

## Contexto tecnico

### Archivos principales

- `app/bindings.go` — `GetRegionPokemonByType()` (lineas 71–107), contiene el `pokedexMap`
- `shell/pokeapi.go` — `FetchRegion()`, `FetchPokedex()`
- `shell/pokeapi_locations.go` — `FetchLocationEncounters()`, `FetchLocationArea()`
- `frontend/src/pages/explore/regions.ts` — renderizado de localizaciones y grafica
- `frontend/src/charts/type-distribution.ts` — grafica ECharts de distribucion de tipos
- `frontend/src/components/location-encounter-modal.ts` — modal de encuentros

### Causa raiz: dos problemas independientes

#### Problema 1: `pokedexMap` incompleto

`GetRegionPokemonByType()` en `app/bindings.go` tiene un mapa estatico que solo cubre Gen I–VI:

```go
pokedexMap := map[string]string{
    "kanto":  "kanto",
    "johto":  "original-johto",
    "hoenn":  "hoenn",
    "sinnoh": "original-sinnoh",
    "unova":  "original-unova",
    "kalos":  "kalos-central",
}
```

Para Gen VII+ el fallback es usar el nombre de la region directamente (`pokedexName = core.NormalizeName(region)`), pero los nombres de pokedex en PokeAPI no coinciden con los nombres de region:

| Region | Pokedex en PokeAPI |
|---|---|
| alola | `original-alola` o `updated-alola` |
| galar | `galar` |
| hisui | `hisui` |
| paldea | `paldea` |

**Fix**: Añadir las entradas faltantes al `pokedexMap`.

#### Problema 2: PokeAPI no tiene datos de encuentros para Gen VII+

El endpoint `/location-area/{name}` de PokeAPI tiene `pokemon_encounters` vacio para localizaciones de juegos Gen VII+ (Sun/Moon, Sword/Shield, Legends Arceus, Scarlet/Violet). Esto es una limitacion de la API externa, no un bug del codigo.

Cuando el usuario hace click en una localizacion de Alola/Galar/etc., el modal de encuentros muestra el mensaje "locationEmpty" correctamente, pero no hay forma de obtener esos datos de PokeAPI.

**Fix para UX**: Informar al usuario de forma clara que los datos de encuentros no estan disponibles para estas generaciones. Opciones:
- Mostrar un aviso en las tarjetas de region Gen VII+ indicando la limitacion.
- Deshabilitar el click en location tags para regiones sin datos de encuentros.
- Mostrar un tooltip o nota bajo la grilla de localizaciones.

---

## Cambios requeridos

### 1. Completar `pokedexMap` — `app/bindings.go`

```go
pokedexMap := map[string]string{
    "kanto":  "kanto",
    "johto":  "original-johto",
    "hoenn":  "hoenn",
    "sinnoh": "original-sinnoh",
    "unova":  "original-unova",
    "kalos":  "kalos-central",
    // Gen VII+
    "alola":  "original-alola",
    "galar":  "galar",
    "hisui":  "hisui",
    "paldea": "paldea",
}
```

**Nota**: verificar los nombres exactos de pokedex en PokeAPI (`/pokedex/{name}`). Algunos pueden ser `updated-alola`, `isle-of-armor`, `crown-tundra`, etc. Usar la version "principal" de cada region.

### 2. Grafica de distribucion de tipos — verificar

Con el `pokedexMap` corregido, `GetRegionPokemonByType()` deberia devolver datos correctos para Gen VII+, y la grafica de distribucion se rellenaria automaticamente.

### 3. Aviso de datos de encuentros limitados — `regions.ts`

Añadir una nota visual bajo las localizaciones para regiones Gen VII+:

```typescript
const GEN7_PLUS = new Set(["alola", "galar", "hisui", "paldea"]);

// Dentro de loadRegionDetail, tras la grilla de localizaciones:
const encounterNote = GEN7_PLUS.has(regionName)
  ? `<p class="region-encounter-note">${t("regions.encounterLimited")}</p>`
  : "";
```

### 4. Añadir clave i18n — `locales/es.json` y `locales/en.json`

```json
// es.json
"regions": {
  "encounterLimited": "Los datos de encuentros por localizacion no estan disponibles para esta generacion."
}

// en.json
"regions": {
  "encounterLimited": "Location encounter data is not available for this generation."
}
```

### 5. Opcional: deshabilitar click en locations sin datos

Para regiones Gen VII+, hacer las location tags no-clickeables o con estilo "disabled", ya que el modal siempre mostrara vacio.

---

## Archivos afectados

### Backend
- `app/bindings.go` — completar `pokedexMap` con Gen VII+

### Frontend
- `frontend/src/pages/explore/regions.ts` — aviso de datos limitados, opcional deshabilitar click
- `frontend/src/locales/es.json` — nueva clave `regions.encounterLimited`
- `frontend/src/locales/en.json` — nueva clave `regions.encounterLimited`

### Core
Ninguno.

### Shell
Ninguno.

---

## Acceptance criteria

- [ ] La grafica de distribucion de tipos se muestra correctamente para Alola, Galar, Hisui y Paldea.
- [ ] El `pokedexMap` incluye todas las regiones Gen VII–IX con los nombres correctos de PokeAPI.
- [ ] Las regiones Gen VII+ muestran un mensaje informativo sobre la limitacion de datos de encuentros.
- [ ] Las regiones Gen I–VI no se ven afectadas (sin regresiones).
- [ ] Los location tags de Gen VII+ estan deshabilitados o muestran un mensaje claro al hacer click.
- [ ] Las nuevas claves i18n existen en ambos idiomas (es y en).
- [ ] Sin errores en consola al expandir regiones Gen VII+.

---

## Dependencias

- **Tarea 0129** (regions-locations-i18n): si se implementa antes, los nombres localizados aplican tambien a Gen VII+.
