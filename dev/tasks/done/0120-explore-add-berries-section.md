# Task 0120 — Explorar: añadir sección de Bayas (Berries) con vista tarjetas y tabla

## Estado: done

## Goal

Añadir una nueva subpestaña **"Bayas"** dentro de la pestaña **Explorar**, con datos obtenidos de PokeAPI (`/api/v2/berry/`). La sección debe ofrecer:
- **Vista de tarjetas** (por defecto), igual que la Pokédex.
- **Vista de tabla**, convertible desde tarjetas con el mismo botón de toggle.

Cada baya mostrará: nombre, imagen del item (desde PokeAPI), firmeza, tipo de regalo natural, poder, y sabores principales.

## Contexto técnico

### Backend — Core

**`core/domain.go`** — añadir nuevos tipos:
```go
type Berry struct {
    ID               int
    Name             string
    GrowthTime       int
    MaxHarvest       int
    NaturalGiftPower int
    Size             int
    Smoothness       int
    SoilDryness      int
    Firmness         string
    Flavors          []BerryFlavor
    NaturalGiftType  string
    ItemName         string   // nombre del item correspondiente (para imagen)
    ItemSprite       string   // URL de la imagen del item desde PokeAPI
}

type BerryFlavor struct {
    Flavor  string
    Potency int
}

type BerryListItem struct {
    Name string
    URL  string
}

type BerryListResponse struct {
    Count   int
    Results []BerryListItem
}
```

**`core/ports.go`** — añadir al interface `PokemonFetcher`:
```go
// --- Bayas ---
FetchBerryList() (BerryListResponse, error)
FetchBerry(name string) (Berry, error)
```

### Backend — Shell

**`shell/pokeapi_berries.go`** — nuevo fichero:

Endpoint PokeAPI: `GET {baseURL}/berry/{name}`

Estructura raw de la respuesta de PokeAPI:
```go
type apiBerry struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    GrowthTime  int    `json:"growth_time"`
    MaxHarvest  int    `json:"max_harvest"`
    NaturalGiftPower int `json:"natural_gift_power"`
    Size        int    `json:"size"`
    Smoothness  int    `json:"smoothness"`
    SoilDryness int    `json:"soil_dryness"`
    Firmness    struct { Name string `json:"name"` } `json:"firmness"`
    Flavors     []struct {
        Flavor  struct { Name string `json:"name"` } `json:"flavor"`
        Potency int `json:"potency"`
    } `json:"flavors"`
    NaturalGiftType struct { Name string `json:"name"` } `json:"natural_gift_type"`
    Item struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"item"`
}
```

Para la imagen del item, hacer una segunda llamada a `/item/{name}` de PokeAPI y extraer `sprites.default`. Alternativamente, la URL del sprite del item sigue el patrón conocido de PokeAPI CDN.

Implementar:
```go
func (c *PokeAPIClient) FetchBerryList() (core.BerryListResponse, error)
func (c *PokeAPIClient) FetchBerry(name string) (core.Berry, error)
```

### Backend — App (bindings)

**`app/bindings.go`** — añadir:
```go
// ListBerries retorna la lista de todas las bayas.
func (a *App) ListBerries() (core.BerryListResponse, error) {
    return a.fetcher.FetchBerryList()
}

// GetBerry retorna el detalle de una baya por nombre.
func (a *App) GetBerry(name string) (core.Berry, error) {
    return a.fetcher.FetchBerry(core.NormalizeName(name))
}
```

### Backend — Mobile server

**`app/mobile/handlers.go`** — añadir a `RegisterRoutes`:
```
ListBerries  → GET  /api/berries
GetBerry     → GET  /api/berries/{name}
```

Añadir los handlers correspondientes siguiendo el patrón existente.

### Frontend — API

**`frontend/src/api.ts`** — añadir:
```ts
export function ListBerries(): Promise<core.BerryListResponse> {
  if (isWails()) return wails("ListBerries");
  return get("/api/berries");
}

export function GetBerry(name: string): Promise<core.Berry> {
  if (isWails()) return wails("GetBerry", name);
  return get(`/api/berries/${encodeURIComponent(name)}`);
}
```

Los tipos `core.BerryListResponse` y `core.Berry` estarán disponibles en `wailsjs/go/models` tras regenerar los bindings de Wails (`wails generate module` o `wails dev`).

### Frontend — Vista de Bayas

**`frontend/src/pages/explore/berries.ts`** — nuevo fichero.

**Vista tarjetas**:
- Grid CSS con tarjetas, siguiendo el estilo `.poke-card` de la Pokédex.
- Cada tarjeta muestra:
  - Imagen del item (sprite de PokeAPI o fallback a placeholder)
  - Nombre de la baya (capitalizado)
  - Tipo de regalo natural con icono de tipo (SVG de `/assets/types/{name}.svg`)
  - Firmeza (badge pequeño)
- Carga paginada: 50 bayas por bloque con scroll infinito (Intersection Observer), igual que la Pokédex.
- Al hacer click en una tarjeta, expandir o mostrar detalle (modal simple o inline).

**Vista tabla**:
- Columnas: #, Nombre, Imagen, Firmeza, Tipo Natural, Poder Natural, Talla, Suavidad, Cosecha Máx., Sabores (top 1-2 con potencia).
- Ordenación por columnas (al menos: ID, Nombre, Poder, Talla).
- Paginación con selector de límite de filas (misma lógica que Pokédex).

**Toggle tarjetas/tabla**:
- Botón de toggle igual al de la Pokédex (`.filter-pill`) en la barra de controles.
- Mantener el estado de vista al cambiar entre subpestañas y volver.

**Carga**:
- Llamar a `ListBerries()` para obtener la lista (~64 bayas).
- Cargar detalles de cada baya con `GetBerry(name)` de forma lazy (en scroll o al cambiar a tabla se cargan todas con overlay similar al de la Pokédex).
- Mostrar indicador de carga mientras se obtienen los datos.

### Frontend — Integración en Explorar

**`frontend/src/pages/explore.ts`**:
- Añadir `"berries"` a `ExploreTab`, `TAB_KEYS`, `tabInited`.
- `import { initBerries } from "./explore/berries"`.
- `case "berries": initBerries(panel)` en `initPanel()`.
- Icono sugerido: un icono de hoja/fruta de Heroicons (e.g., `ICON_LEAF` o similar).
- Label i18n: `explore.tabs.berries`.

**`frontend/index.html`** — no requiere cambios si los paneles se generan dinámicamente (igual que las otras subpestañas de Explorar).

### Localización

Añadir en `es.json` y `en.json`:
```json
// es.json — dentro de "explore.tabs"
"berries": "Bayas"

// es.json — nueva sección raíz
"berries": {
  "title": "Bayas",
  "loading": "Cargando bayas...",
  "searchPlaceholder": "Buscar baya...",
  "noResults": "No se encontraron bayas.",
  "error": "Error al cargar bayas.",
  "count": "{count} bayas",
  "tableView": "⊞ Tabla",
  "cardView": "⊟ Tarjetas",
  "columns": {
    "id": "#",
    "name": "Nombre",
    "sprite": "Imagen",
    "firmness": "Firmeza",
    "naturalGiftType": "Tipo Natural",
    "naturalGiftPower": "Poder",
    "size": "Talla",
    "smoothness": "Suavidad",
    "maxHarvest": "Cosecha Máx.",
    "flavors": "Sabores"
  }
}
```

## Acceptance criteria

- [ ] Aparece la subpestaña "Bayas" en la navegación de Explorar con su icono.
- [ ] Al entrar en la subpestaña se muestra la vista de tarjetas con las ~64 bayas.
- [ ] Cada tarjeta muestra imagen, nombre, tipo natural y firmeza.
- [ ] El scroll infinito carga más bayas al llegar al final de la página.
- [ ] El botón de toggle cambia entre vista tarjetas y vista tabla.
- [ ] La vista tabla muestra todas las columnas especificadas y permite ordenar.
- [ ] La paginación de la tabla funciona (límite de filas + botones anterior/siguiente).
- [ ] Las etiquetas se actualizan al cambiar el idioma (locale-changed).
- [ ] No hay errores de compilación Go (ports satisfechos, tipos correctos).
- [ ] No hay errores en consola del frontend.

## Dependencias

- **Depende de 0119**: La tarea 0119 modifica `explore.ts` (elimina "types" y añade "typeChart"). Esta tarea añade "berries" a ese mismo fichero. Implementar después de que 0119 esté mergeado en `main`.

## Archivos afectados

### Backend
- `core/domain.go` — añadir `Berry`, `BerryFlavor`, `BerryListItem`, `BerryListResponse`
- `core/ports.go` — añadir `FetchBerryList()` y `FetchBerry()` a `PokemonFetcher`
- `shell/pokeapi_berries.go` — **crear nuevo** con implementación PokeAPI
- `app/bindings.go` — añadir `ListBerries()` y `GetBerry()`
- `app/mobile/handlers.go` — añadir rutas `/api/berries` y `/api/berries/{name}`

### Frontend
- `frontend/src/api.ts` — añadir `ListBerries()` y `GetBerry()`
- `frontend/src/pages/explore/berries.ts` — **crear nuevo** con vista tarjetas + tabla
- `frontend/src/pages/explore.ts` — añadir subpestaña "berries"
- `frontend/src/locales/es.json` — añadir claves berries
- `frontend/src/locales/en.json` — añadir claves berries
