# Task 0122 — Explorar: Bayas — fondo oscuro en tarjetas y columna de efecto

## Estado: done

## Goal

Dos correcciones/mejoras sobre la subpestaña **"Bayas"** dentro de Explorar:

1. **Fix dark mode en tarjetas**: el fondo de `.berry-card` es blanco en modo oscuro; debe adaptarse al tema oscuro igual que `.poke-card`.

2. **Columna de efecto**: añadir una columna/dato que muestre el **efecto** de la baya (p.ej. "Restaura 10 PS", "Cura el estado de quemadura"). Este dato se extrae del campo `short_effect` del endpoint de items de PokeAPI y requiere cambios en backend y frontend.

## Contexto técnico

### Fix 1: Dark mode en tarjetas de baya

**`frontend/src/styles/_dark.scss`** — dentro del bloque `[data-bs-theme="dark"]`, añadir sección "Berries":

```scss
// ─── Berries ─────────────────────────────────────────────────────────────
.berry-card {
  background: #2d3748;
  border-color: #4a5568;
}

.berry-card__img {
  background: #1a202c;
}

.berry-card__name {
  color: #e2e8f0;
}

.berry-firmness-badge {
  background: #4a5568;
  color: #e2e8f0;
}
```

No se requieren cambios en backend ni en TypeScript para este fix.

---

### Fix 2: Columna de efecto de la baya

#### Backend — Core

**`core/domain.go`** — añadir campo `Effect` a `Berry`:
```go
type Berry struct {
    // ... campos existentes ...
    Effect string  // efecto corto del item (p.ej. "Restores 10 HP.")
}
```

#### Backend — Shell

**`shell/pokeapi_berries.go`** — extender `apiItem` para leer `effect_entries`:
```go
type apiItem struct {
    Sprites struct {
        Default string `json:"default"`
    } `json:"sprites"`
    EffectEntries []struct {
        ShortEffect string `json:"short_effect"`
        Language    struct {
            Name string `json:"name"`
        } `json:"language"`
    } `json:"effect_entries"`
}
```

En `FetchBerry()`, ya se hace la llamada al endpoint `/item/{name}`. Tras obtener `item`, extraer el `ShortEffect` en inglés:
```go
effect := ""
for _, e := range item.EffectEntries {
    if e.Language.Name == "en" {
        effect = e.ShortEffect
        break
    }
}
```

Añadir `Effect: effect` al return de `core.Berry{...}`.

No se requieren cambios en `core/ports.go`, `app/bindings.go` ni en `app/mobile/handlers.go` porque los endpoints existentes (`GetBerry`) ya devuelven el struct `Berry` completo — el nuevo campo se añade automáticamente a la respuesta JSON.

#### Frontend — Bindings Wails

Tras añadir `Effect string` al struct `core.Berry`, regenerar los bindings de Wails (`wails generate module` o `wails dev`) para que el tipo TypeScript `core.Berry` en `wailsjs/go/models.ts` incluya el nuevo campo.

#### Frontend — Vista de Bayas

**`frontend/src/pages/explore/berries.ts`**:

- **Vista tabla**: añadir columna "Efecto" entre "Sabores" y el final (o donde mejor encaje). No sortable. Renderizar el texto del campo `Effect` truncado a ~60 caracteres si es muy largo (usar `title` en el `<td>` con el texto completo).

- **Vista tarjetas**: mostrar el efecto como una línea de texto pequeña debajo del nombre o de la firmeza, solo si `Effect` no está vacío. Clase sugerida: `berry-card__effect`. Limitar a 1 línea con `text-overflow: ellipsis`.

#### Frontend — Estilos

**`frontend/src/styles/_explore.scss`** — añadir dentro de `.berry-card`:
```scss
.berry-card__effect {
  font-size: 0.62rem;
  color: #718096;
  text-align: center;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding: 0 0.25rem;
}
```

**`frontend/src/styles/_dark.scss`** — dentro de la sección Berries (ya añadida en Fix 1):
```scss
.berry-card__effect {
  color: #a0aec0;
}
```

#### Localización

**`frontend/src/locales/es.json`** y **`en.json`** — añadir en `"berries.columns"`:
```json
"effect": "Efecto"      // es.json
"effect": "Effect"      // en.json
```

## Archivos afectados

### Backend
- `core/domain.go` — añadir campo `Effect string` a `Berry`
- `shell/pokeapi_berries.go` — extender `apiItem` con `EffectEntries`, extraer y mapear `Effect`

### Frontend
- `frontend/src/pages/explore/berries.ts` — añadir columna Efecto en tabla, texto en tarjeta
- `frontend/src/styles/_explore.scss` — añadir `.berry-card__effect`
- `frontend/src/styles/_dark.scss` — añadir overrides dark mode para `.berry-card` y `.berry-card__effect`
- `frontend/src/locales/es.json` — añadir `berries.columns.effect`
- `frontend/src/locales/en.json` — añadir `berries.columns.effect`

## Acceptance criteria

- [ ] En modo oscuro, las tarjetas de baya tienen fondo oscuro (#2d3748), no blanco.
- [ ] En modo oscuro, el área de imagen, el nombre y la firmeza de la tarjeta se adaptan al tema.
- [ ] La vista tabla de bayas muestra una columna "Efecto" con el efecto corto de cada baya.
- [ ] Las tarjetas de baya muestran el efecto corto (si existe) debajo del nombre/firmeza.
- [ ] El texto de efecto no desborda su contenedor (truncado con ellipsis si procede).
- [ ] Las etiquetas se actualizan al cambiar el idioma.
- [ ] No hay errores de compilación Go.
- [ ] No hay errores en consola del frontend.

## Dependencias

Ninguna. Tarea independiente sobre código ya existente en `main`.
