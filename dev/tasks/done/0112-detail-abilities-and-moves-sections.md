# Task 0112 — Secciones de Habilidades y Movimientos en la vista de detalle

## Goal
Añadir dos nuevas secciones en la vista de detalle de un Pokémon (debajo de Encounters):
1. **Habilidades**: lista de habilidades del Pokémon con su descripción.
2. **Movimientos**: lista de movimientos que puede aprender, con el método de aprendizaje (level-up, machine/MT, egg, tutor) y el nivel (si aplica). La lista es filtrable por método de aprendizaje.

## Context

### Backend — Habilidades
- `core.Pokemon` (en `core/domain.go`) **no incluye** habilidades actualmente. La respuesta raw de PokeAPI `/pokemon/{name}` SÍ devuelve un campo `abilities` con nombre y si es oculta.
- `shell/pokeapi.go`: `apiPokemon` struct (línea 27) no tiene el campo `Abilities`. Hay que añadirlo.
- `core.Ability` tiene Name, NameEs, Description, DescriptionEs y Pokemon[].
- Para obtener la descripción de cada habilidad se puede llamar `FetchAbility(name)` (ya existe). Un Pokémon tiene 1-3 habilidades, por lo que es asumible.
- Alternativa más simple (sin necesidad de llamadas adicionales): añadir `Abilities []string` (solo nombres) a `core.Pokemon` y mostrarlos sin descripción, o fetchear la descripción en el frontend con `GetAbility(name)`.

### Backend — Movimientos
- `core.Pokemon` ya tiene `Moves []PokemonMoveEntry` con `{Name, Method, Level}`.
- No se necesitan cambios en el backend para movimientos.

### Frontend
- `frontend/src/types.ts`: `Pokemon` interface no tiene `Abilities`.
- `frontend/src/pages/pokedex.ts`: `renderDetail(p)` genera el HTML de detalle. Añadir secciones después del bloque `#pokemon-encounters`.
- `frontend/src/api.ts`: existe `GetAbility(name)` para fetchear descripción por nombre.

## Approach elegido

### Habilidades
1. Añadir `Abilities []string` (nombres) a `core.Pokemon` en `core/domain.go`.
2. Añadir campo `abilities` a `apiPokemon` en `shell/pokeapi.go` y popularlo en `toDomainPokemon`.
3. En el frontend (`types.ts`), añadir `Abilities: string[]` a `Pokemon`.
4. En `renderDetail`, renderizar la sección de habilidades: para cada nombre, llamar `GetAbility(name)` en el frontend para obtener descripción en el idioma activo.

### Movimientos
- Usar `p.Moves` ya disponible.
- Renderizar tabla/lista con columnas: Nombre, Método, Nivel.
- Añadir filtros de método: "Todos" | "level-up" | "machine" | "egg" | "tutor".
- Los métodos relevantes de PokeAPI: `level-up`, `machine`, `egg`, `tutor`.
- El nivel se muestra solo cuando `method === "level-up"` y `level > 0`.

## Steps

### Backend

#### 1. `core/domain.go`
Añadir campo a `Pokemon`:
```go
type Pokemon struct {
    // ... campos existentes ...
    Abilities []string
}
```

#### 2. `shell/pokeapi.go`
En `apiPokemon` añadir:
```go
Abilities []struct {
    Ability struct {
        Name string `json:"name"`
    } `json:"ability"`
    IsHidden bool `json:"is_hidden"`
} `json:"abilities"`
```

En `toDomainPokemon`, popularlo:
```go
abilities := make([]string, len(raw.Abilities))
for i, a := range raw.Abilities {
    abilities[i] = a.Ability.Name
}
// ... en el return:
Abilities: abilities,
```

### Frontend

#### 3. `frontend/src/types.ts`
Añadir `Abilities: string[]` a `Pokemon` interface.

#### 4. `frontend/src/pages/pokedex.ts` — sección habilidades
En `renderDetail`, añadir al HTML:
```html
<div id="pokemon-abilities" class="pokemon-abilities">
  <p class="loading">Cargando habilidades...</p>
</div>
```
Llamar `loadAbilities(p.Name, p.Abilities)` al renderizar.

`loadAbilities` carga cada habilidad con `GetAbility(name)` (en paralelo con `Promise.all`) y renderiza:
```html
<h3>Habilidades</h3>
<div class="abilities-list">
  <div class="ability-card">
    <span class="ability-name">Torrent</span>
    <p class="ability-desc">Potencia los movimientos de agua...</p>
    <!-- badge si es oculta -->
  </div>
</div>
```
Nota: como `Abilities` solo tiene nombres (sin saber cuál es oculta), se puede omitir el badge de oculta en esta iteración.

#### 5. `frontend/src/pages/pokedex.ts` — sección movimientos
Añadir al HTML de detalle:
```html
<div id="pokemon-moves" class="pokemon-moves">
  <!-- renderizado por renderMoves -->
</div>
```

Función `renderMoves(moves: PokemonMoveEntry[])`:
- Genera botones de filtro por método: "Todos", "level-up", "machine", "egg", "tutor".
- Renderiza tabla de movimientos filtrados.
- Al cambiar filtro, re-renderiza la tabla.

Formato de la tabla:
| Nombre | Método | Nivel |
|--------|--------|-------|
| tackle | level-up | 1 |
| surf | machine | — |

#### 6. i18n
Añadir claves a `es.json` y `en.json`:
- `detail.abilities`: "Habilidades" / "Abilities"
- `detail.moves`: "Movimientos" / "Moves"
- `detail.moveMethod.levelUp`: "Nivel" / "Level-up"
- `detail.moveMethod.machine`: "MT/MO" / "TM/HM"
- `detail.moveMethod.egg`: "Huevo" / "Egg"
- `detail.moveMethod.tutor`: "Tutor" / "Tutor"
- `detail.moveMethod.all`: "Todos" / "All"

## Files to modify
- `core/domain.go`
- `shell/pokeapi.go`
- `frontend/src/types.ts`
- `frontend/src/pages/pokedex.ts`
- `frontend/src/locales/es.json`
- `frontend/src/locales/en.json`
- Posiblemente SCSS para estilos de `.abilities-list`, `.pokemon-moves`, tabla de movimientos

## Acceptance criteria
- [ ] La vista de detalle muestra una sección "Habilidades" con nombre y descripción de cada habilidad
- [ ] La vista de detalle muestra una sección "Movimientos" con la lista completa de movimientos
- [ ] Cada movimiento muestra cómo se aprende (método + nivel si aplica)
- [ ] La lista de movimientos es filtrable por método de aprendizaje
- [ ] Los datos de habilidades están en el campo `Abilities` de `core.Pokemon` (backend)
- [ ] Todo funciona tanto en Wails (desktop) como en HTTP (mobile)
