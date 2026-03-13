# 0011 — Pokedex: Filtros + Rediseño Visual

## Descripción

Mejorar la vista principal de la Pokédex añadiendo filtros por generación, tipo y estado legendario, y rediseñar las tarjetas y tipografía con más personalidad visual (aura "Out of the box REAL").

---

## Alcance

### Capas involucradas
- **Frontend** (TypeScript, SCSS, HTML): filtros, rediseño tarjetas y fuentes
- **APP** (`app/bindings.go`): nuevos métodos IPC si son necesarios
- **Shell** (ya tiene los endpoints base): sin cambios previstos

---

## Contexto técnico

### API disponible (ya expuesta al frontend via Wails IPC)
| Método | Uso para filtro |
|--------|----------------|
| `ListGenerations()` | Lista de generaciones (gen-i … gen-ix) |
| `GetGeneration(name)` | Devuelve `Generation{PokemonSpecies []string}` → base del filtro por gen |
| `ListTypes()` | Lista de todos los tipos |
| `GetType(name)` | Devuelve `PokemonTypeDetail{Pokemon []TypePokemonEntry}` → base del filtro por tipo |
| `GetPokemonSpecies(name)` | Devuelve `PokemonSpecies{IsLegendary, IsMythical}` → para el filtro legendario |

### Nota sobre el filtro de legendarios
PokéAPI no tiene un endpoint de "todos los legendarios". La estrategia: cuando se activa el filtro, se itera sobre la lista filtrada actual (gen o tipo, o todos si no hay filtro) y se consulta species en paralelo (batching) para determinar legendario. Para el modo "todos" se necesita un `GetPokemonByFilter` nuevo en el backend (ver paso 3).

---

## Plan de implementación

### Paso 1 — Nuevo binding `ListPokemonByGeneration` (Go)
**Archivo:** `app/bindings.go`

Añadir:
```go
// ListPokemonByGeneration retorna la lista de species de una generación.
func (a *App) ListPokemonByGeneration(name string) (core.Generation, error) {
    return a.fetcher.FetchGeneration(core.NormalizeName(name))
}
```
*(Ya existe `GetGeneration`, podría reutilizarse — revisar si el binding actual ya está expuesto al frontend con ese nombre.)*

Verificar que `GetGeneration` ya genera binding Wails correcto. Si es así, no hace falta cambio en Go.

### Paso 2 — Barra de filtros (HTML)
**Archivo:** `frontend/index.html`

Añadir dentro de `#list-view`, antes del grid:
```html
<div id="filter-bar">
  <select id="filter-gen">...</select>
  <select id="filter-type">...</select>
  <button id="filter-legendary" class="filter-toggle">⭐ Legendario</button>
  <button id="filter-mythical" class="filter-toggle">✨ Mítico</button>
  <button id="filter-reset">Todos</button>
</div>
```

### Paso 3 — Lógica de filtros (TypeScript)
**Archivo:** `frontend/src/pages/pokedex.ts`

Estrategia de datos:
- **Sin filtro**: comportamiento actual (paginación de API)
- **Filtro por generación**: `GetGeneration(name)` → `PokemonSpecies[]` → paginar client-side
- **Filtro por tipo**: `GetType(name)` → `TypePokemonEntry[]` → paginar client-side
- **Filtro legendario/mítico**: requiere species data → nuevo binding en backend (Paso 3b)

**Paso 3b — Binding `ListLegendaryPokemon`** (si la carga perezosa es insuficiente)
Alternativa pragmática: al seleccionar legendario, iterar las primeras N generaciones y batching de `GetPokemonSpecies` con `Promise.allSettled`. Limitar a 50 simultáneas.

Estado del filtro como objeto:
```typescript
interface FilterState {
  generation: string | null;
  type: string | null;
  legendary: boolean;
  mythical: boolean;
}
```

Cuando hay filtro activo: deshabilitar paginación de API y usar paginación client-side sobre la lista filtrada.

### Paso 4 — Rediseño visual de tarjetas (SCSS)
**Archivos:** `frontend/src/styles/_pokemon.scss`, `frontend/src/styles/_variables.scss`

**Fuente:** Añadir Google Font con personalidad. Candidatas:
- **`Exo 2`**: tech/sci-fi, legible, atrevida pero funcional → primera opción
- **`Bebas Neue`**: condensada, impacto, para nombres y headings
- Mezclar: Bebas Neue para títulos, Exo 2 para texto normal

**Tarjetas rediseñadas:**
- Sprite más grande (128px → 140px con sombra de color del tipo primario)
- Background de la tarjeta con gradiente sutil del color del tipo primario (10% opacidad)
- Número de Pokédex visible (#001) en esquina top-left, fuente pequeña y grisácea
- Badge de tipo(s) dentro de la tarjeta (pill, color del tipo)
- Nombre en Bebas Neue o similar, más grande
- Hover: elevación + glow del color del tipo
- Border-radius asimétrico (18px 12px 18px 12px) para "imperfección intencionada" (aura)

**Barra de filtros:**
- Pills horizontales con scroll si overflow
- Toggle activo con color sólido del tipo/gen seleccionado
- Transición suave al cambiar filtros (fade out → nueva lista → stagger entrada)

### Paso 5 — Animación de transición al filtrar (TypeScript + GSAP)
**Archivo:** `frontend/src/animations/transitions.ts`

Añadir `fadeAndReplace(grid, renderFn)`: fade out del grid, ejecutar render, stagger entrada de nuevas cards.

---

## Archivos a modificar

| Archivo | Cambio |
|---------|--------|
| `frontend/index.html` | Añadir `#filter-bar` con selects y toggles |
| `frontend/src/pages/pokedex.ts` | Lógica de filtros, estado, paginación client-side |
| `frontend/src/styles/_pokemon.scss` | Rediseño de `.card`, `.filter-bar`, layout |
| `frontend/src/styles/_variables.scss` | Nueva variable de fuente, quizás ajuste de colores |
| `frontend/src/animations/transitions.ts` | `fadeAndReplace` para cambio de filtro |
| `app/bindings.go` | Verificar/exponer binding de generación (probablemente ya ok) |

---

## Criterios de aceptación

- [ ] Dropdown de generación funcional (gen-i … gen-ix)
- [ ] Dropdown de tipo funcional (todos los tipos disponibles)
- [ ] Toggle legendario funciona y muestra solo legendarios
- [ ] Toggle mítico funciona y muestra solo míticos
- [ ] Botón "Todos" limpia filtros y vuelve a la lista paginada normal
- [ ] Los filtros son combinables (gen + tipo, gen + legendario, etc.)
- [ ] Tarjetas muestran: sprite grande, número, nombre bold, badge(s) de tipo
- [ ] Background de tarjeta con tinte del tipo primario
- [ ] Hover con glow del tipo
- [ ] Fuente con personalidad (Bebas Neue / Exo 2) aplicada
- [ ] Transición suave al cambiar filtros
- [ ] La paginación client-side funciona cuando hay filtro activo

---

## Dependencias externas nuevas

- Google Fonts (Bebas Neue + Exo 2) — añadir al `<head>` del HTML
- No hay nuevas dependencias npm

---

## Notas de diseño

Aura "Out of the box REAL":
- Los números de pokédex en esquina: pequeños, sin protagonismo pero ahí
- Las tarjetas no son simétricas entre sí (border-radius variable opcionalmente)
- El tinte del tipo es sutil, no invasivo
- El grid no tiene columnas rígidas: `auto-fill minmax(160px, 1fr)` con gap variado
- Los filtros activos se "marcan" con el color del tipo/gen, no un gris genérico
