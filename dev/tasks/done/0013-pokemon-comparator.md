# 0013 — Pokémon Comparator Tab

## Descripción
Nueva pestaña **"Comparar"** que permite seleccionar 2 Pokémon y ver sus stats base comparados lado a lado, con indicadores visuales de diferencia (qué stat es mayor en cada uno y por cuánto).

## Estado
- [x] Done

## Contexto
- Ya existe `core.Pokemon` con `Stats []Stat` (HP, Attack, Defense, Sp.Atk, Sp.Def, Speed).
- Ya existe `app.GetPokemon(name)` como binding Wails accesible desde frontend.
- El chart de stats usa ECharts (`charts/stats-chart.ts`) — se puede reutilizar o extender para overlay comparativo.
- El router (`router.ts`) tiene un sistema `registerPage` simple para añadir tabs.

## Capas afectadas
- **Core** (`core/`): nueva función pura `CompareStats`.
- **App** (`app/bindings.go`): nuevo binding `ComparePokemons`.
- **Frontend**: nueva página `pages/compare.ts`, nuevo SCSS `_compare.scss`, nuevo tab en HTML.

## Plan de implementación

### Paso 1 — Core: función pura de comparación
En `core/pokemon.go` (o nuevo `core/compare.go`), añadir:
```go
type StatComparison struct {
    Name    string
    StatA   int
    StatB   int
    Diff    int    // StatA - StatB (positivo = A gana, negativo = B gana)
    Winner  string // "a", "b", "tie"
}

type PokemonComparison struct {
    PokemonA core.Pokemon
    PokemonB core.Pokemon
    Stats    []StatComparison
    TotalA   int
    TotalB   int
    Winner   string // "a", "b", "tie"
}

func ComparePokemons(a, b Pokemon) PokemonComparison
```

### Paso 2 — App: binding
En `app/bindings.go`, añadir:
```go
func (a *App) ComparePokemons(nameA, nameB string) (core.PokemonComparison, error)
```
Fetcha ambos Pokémon y llama a `core.ComparePokemons`.

### Paso 3 — Wails: regenerar bindings
Ejecutar `wails generate module` (o el proyecto lo hace en build) para que `ComparePokemons` esté disponible en `wailsjs/go/app/App.ts`.

### Paso 4 — Frontend: página compare
Crear `frontend/src/pages/compare.ts`:
- Dos inputs de búsqueda (uno para cada Pokémon) con autocomplete o campo de texto simple.
- Botón "Comparar".
- Al comparar, llamar `ComparePokemons(nameA, nameB)`.
- Renderizar:
  - Sprites de ambos Pokémon en columnas.
  - Tabla de stats con barra de progreso horizontal para cada stat.
  - Color verde en el stat mayor, rojo en el menor, gris en empate.
  - Fila de "Total" al final con el ganador global.
- Alternativamente, un chart de radar ECharts con overlay de los dos Pokémon.

### Paso 5 — Frontend: estilos
Crear `frontend/src/styles/_compare.scss`:
- Layout de dos columnas con `.compare-panel` para cada Pokémon.
- `.stat-row` con barras de progreso relativas (max stat = 255).
- Indicadores de diferencia (`.stat-diff`).

### Paso 6 — HTML y router
- Añadir tab "Comparar" en `index.html`.
- Registrar la página en `main.ts` con `registerPage`.
- Importar y llamar `initCompare()`.

## Criterios de éxito
- Se pueden comparar 2 Pokémon cualesquiera.
- Se muestra claramente qué stat gana cada uno y por cuánto.
- El diseño es coherente con el resto de la app (ver tarea 0012).
- `ComparePokemons` en Core es una función pura y testeable.
