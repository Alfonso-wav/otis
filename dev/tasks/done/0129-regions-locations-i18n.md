# Task 0129 — Regiones: localizaciones traducidas al español

## Estado: pending

## Goal

Los nombres de las localizaciones dentro de cada region (Explorar > Regiones) se muestran siempre en ingles (slugs como "pallet-town" → "pallet town"), incluso cuando la app esta en español. Deben mostrarse en el idioma activo del usuario, usando los nombres localizados que proporciona PokeAPI.

---

## Contexto tecnico

### Archivos principales

- `core/domain.go` — struct `Location` (solo tiene `Name` y `Region`)
- `shell/pokeapi.go` — `FetchRegion()` (lineas 284–309), parsea `apiRegion` que solo guarda `name` (slug)
- `shell/pokeapi_locations.go` — `FetchLocation()`, parsea `apiLocationDetail` (solo guarda `name`)
- `frontend/src/pages/explore/regions.ts` — `loadRegionDetail()`, muestra `l.Name.replace(/-/g, " ")`
- `frontend/src/api.ts` — capa de transporte

### Problema raiz

PokeAPI devuelve los nombres de localizaciones como slugs en ingles (ej: `"pallet-town"`). El endpoint `/location/{name}` tiene un campo `names` con traducciones por idioma:

```json
{
  "name": "pallet-town",
  "names": [
    { "name": "Pueblo Paleta", "language": { "name": "es" } },
    { "name": "Pallet Town", "language": { "name": "en" } }
  ]
}
```

Actualmente el backend no extrae este campo. El struct `Location` solo tiene `Name` (slug) y `Region`.

### Solucion propuesta

**Opcion A (recomendada): enriquecer Location con nombres localizados**

1. Ampliar `Location` en `core/domain.go` con un campo `Names map[string]string` (idioma → nombre localizado).
2. En `FetchRegion()`, para cada location, hacer un fetch de `/location/{slug}` y extraer el array `names`.
3. En el frontend, usar `l.Names[locale]` con fallback a `l.Name.replace(/-/g, " ")`.

**Concern de rendimiento**: una region puede tener 50+ localizaciones. Hacer 50 fetches secuenciales seria lento. Mitigar con:
- **Concurrencia limitada** en Go (ej: pool de 5 goroutines con `errgroup`).
- **Cache en memoria** de localizaciones ya resueltas (el `PokeAPIClient` ya tiene cache).

**Opcion B (mas simple): fetch lazy en frontend**

No cambiar el backend. En el frontend, cuando se expande una region, hacer un fetch por cada location para obtener su nombre localizado. Mas simple pero mas lento visualmente.

**Opcion C: mapa estatico de traducciones**

Mantener un JSON con las traducciones de todas las localizaciones. No requiere fetches extra pero es fragil y requiere mantenimiento manual.

Se recomienda **Opcion A** por ser coherente con la arquitectura (Shell enriquece datos, Core los consume puros).

---

## Cambios requeridos

### 1. Ampliar struct Location — `core/domain.go`

```go
type Location struct {
    Name   string
    Region string
    Names  map[string]string  // language code → localized name
}
```

### 2. Ampliar apiRegion o fetch individual — `shell/pokeapi.go`

Opcion: crear un struct `apiLocationNames` para parsear `/location/{name}`:

```go
type apiLocationNames struct {
    Name  string `json:"name"`
    Names []struct {
        Name     string `json:"name"`
        Language struct {
            Name string `json:"name"`
        } `json:"language"`
    } `json:"names"`
}
```

En `FetchRegion()`, despues de obtener la lista de locations, hacer fetch concurrente de cada `/location/{slug}` para extraer los `names`.

### 3. Exponer locale al frontend — `frontend/src/pages/explore/regions.ts`

Cambiar el renderizado de location tags:

```typescript
// Antes:
`<span class="region-location-tag" data-location="${l.Name}">${l.Name.replace(/-/g, " ")}</span>`

// Despues:
const displayName = l.Names?.[getLocale()] ?? l.Names?.["en"] ?? l.Name.replace(/-/g, " ");
`<span class="region-location-tag" data-location="${l.Name}">${displayName}</span>`
```

### 4. Re-renderizar al cambiar idioma

El listener `locale-changed` ya reinicializa las regiones (`initialized = false; initRegions(lastContainer)`), lo cual re-fetcha y re-renderiza todo. Si los nombres vienen del backend (cached), el cambio de idioma reflejaria el nuevo locale automaticamente.

---

## Archivos afectados

### Core
- `core/domain.go` — ampliar `Location` con `Names`

### Shell
- `shell/pokeapi.go` — enriquecer `FetchRegion()` con nombres localizados (fetch concurrente)

### Frontend
- `frontend/src/pages/explore/regions.ts` — usar nombre localizado en lugar del slug

### Backend
Ninguno adicional (los bindings y handlers ya serializan el struct completo).

---

## Acceptance criteria

- [ ] Las localizaciones se muestran en español cuando la app esta en español.
- [ ] Las localizaciones se muestran en ingles cuando la app esta en ingles.
- [ ] Si no hay traduccion disponible para un idioma, se usa el fallback al nombre ingles o al slug.
- [ ] El rendimiento al expandir una region no se degrada notablemente (< 3s con cache frio).
- [ ] El cache evita re-fetches en expansiones posteriores de la misma region.
- [ ] Al cambiar el idioma, las localizaciones ya cargadas se actualizan al nuevo idioma.
- [ ] Sin regresiones en el modal de encuentros (sigue usando el slug como identificador).
