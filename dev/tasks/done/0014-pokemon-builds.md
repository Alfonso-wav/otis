# 0014 — Pokémon Builds & Damage Simulator Tab

## Descripción
Nueva pestaña **"Builds"** donde el usuario puede:
1. Seleccionar un Pokémon.
2. Ver y equipar hasta 4 ataques **solo compatibles** con ese Pokémon.
3. Introducir sus stats (nivel, IVs, EVs, naturaleza) o usar los stats base.
4. Simular el daño que haría/recibiría con cada ataque contra un Pokémon defensor configurable.

## Estado
- [x] Done

## Contexto
- `core.Pokemon` no incluye movimientos actualmente. La PokéAPI devuelve `moves` en `/pokemon/{name}` — hay que extender el dominio.
- Ya existe `core.EV` calculator y `core.CalculateAllStats` — base para calcular stats finales.
- `core.Move` ya tiene `Power`, `Type`, `Category` (físico/especial/estado) — suficiente para la fórmula de daño.
- La fórmula de daño Gen 5+: `((2*L/5+2) * Poder * Atk/Def / 50 + 2) * modificadores`.
- Ya existe binding `CalculateStats` y `GetNatures` en `app/bindings.go`.

## Capas afectadas
- **Core** (`core/domain.go`, nuevo `core/damage.go`): extender `Pokemon` con moves, añadir fórmula de daño.
- **Shell** (`shell/pokeapi.go`): extender `apiPokemon` para incluir moves, actualizar `toDomainPokemon`.
- **App** (`app/bindings.go`): nuevo binding `GetPokemonMoves`, `SimulateDamage`.
- **Frontend**: nueva página `pages/builds.ts`, `_builds.scss`, tab en HTML.

## Plan de implementación

### Paso 1 — Core: extender dominio con movimientos de Pokémon

En `core/domain.go`, añadir:
```go
// PokemonMoveEntry representa un movimiento que puede aprender un Pokémon
type PokemonMoveEntry struct {
    Name   string
    Method string // level-up, machine, egg, tutor
    Level  int    // nivel al que se aprende (0 si no aplica)
}
```

Extender `core.Pokemon`:
```go
type Pokemon struct {
    // ... campos existentes ...
    Moves []PokemonMoveEntry
}
```

### Paso 2 — Shell: parsear moves de la API
En `shell/pokeapi.go`, extender `apiPokemon`:
```go
Moves []struct {
    Move struct { Name string `json:"name"` } `json:"move"`
    VersionGroupDetails []struct {
        LevelLearnedAt  int    `json:"level_learned_at"`
        MoveLearnMethod struct { Name string `json:"name"` } `json:"move_learn_method"`
    } `json:"version_group_details"`
} `json:"moves"`
```

Actualizar `toDomainPokemon` para mapear moves a `[]PokemonMoveEntry`.
Usar solo el primer `version_group_details` para simplicidad (o el más reciente).

### Paso 3 — Core: fórmula de daño (función pura)
Crear `core/damage.go`:

```go
type DamageInput struct {
    AttackerStats  Stats
    DefenderStats  Stats
    Move           Move
    AttackerTypes  []PokemonType
    DefenderTypes  []PokemonType
    Level          int
    IsCritical     bool
    WeatherBonus   float64 // 1.0 por defecto
}

type DamageResult struct {
    Min        int
    Max        int
    Average    int
    Multiplier float64 // efectividad de tipo
    IsSuperEffective bool
    IsNotVeryEffective bool
    HasNoEffect bool
}

func CalculateDamage(input DamageInput) DamageResult
func TypeEffectiveness(moveType string, defenderTypes []PokemonType) float64
```

La tabla de efectividad de tipos es una constante en Core (25 combinaciones de tipos × 18 tipos = mapa estático).
La fórmula: `((2*Level/5+2) * Power * AtkStat/DefStat / 50 + 2) * STAB * TypeEff * Critical * Random`.
- STAB: 1.5 si el tipo del move coincide con algún tipo del atacante.
- Critical: 1.5 si `IsCritical`.
- Random: range 0.85–1.00 → calcular Min (×0.85) y Max (×1.00).

### Paso 4 — App: bindings
En `app/bindings.go`:

```go
// GetPokemonMoves retorna los movimientos compatibles de un Pokémon (ya incluidos en GetPokemon tras el paso 2).
// SimulateDamage calcula el daño de un movimiento dado atacante y defensor configurados.
func (a *App) SimulateDamage(input core.DamageInput) (core.DamageResult, error)
```

`SimulateDamage` es llamada directamente con el input construido en el frontend — es pura delegación a Core.

### Paso 5 — Frontend: página builds
Crear `frontend/src/pages/builds.ts`:

**Sección 1 — Selección de Pokémon atacante:**
- Input de búsqueda + botón → llama `GetPokemon(name)`.
- Muestra sprite, tipos, stats base.

**Sección 2 — Equipar ataques (max 4):**
- Lista desplegable/searchable con los moves del Pokémon (del campo `Moves`).
- Al seleccionar un move, muestra su nombre, tipo (badge), categoría, potencia, precisión.
- Hasta 4 slots.
- Para cada move seleccionado: botón "Ver detalle" que llama `GetMove(name)`.

**Sección 3 — Configurar stats del atacante:**
- Inputs: Nivel (1-100), Naturaleza (select), IVs (6 inputs), EVs (6 inputs con validación ≤252 cada uno, ≤510 total).
- Botón "Calcular stats" → llama `CalculateStats` y muestra stats finales calculados.
- Opción "Usar stats base" para simplificar.

**Sección 4 — Defensor:**
- Input búsqueda del Pokémon defensor.
- Misma configuración de stats que el atacante.

**Sección 5 — Simulación:**
- Por cada move equipado: tabla con columnas Movimiento | Tipo | Categoría | Daño Mín | Daño Máx | Efectividad.
- Botón "¿Critico?" toggle por move.
- Resaltar moves super efectivos (verde) / no muy efectivos (rojo) / sin efecto (gris).

### Paso 6 — Frontend: estilos
Crear `frontend/src/styles/_builds.scss`:
- Layout de 2 columnas en pantallas grandes (atacante / defensor) + tabla de daño abajo.
- `.build-slot` para los 4 slots de moves con borde del color del tipo del move.
- `.damage-row` para cada fila de la tabla de simulación.
- Consistente con el diseño de la Pokédex.

### Paso 7 — HTML y router
- Añadir tab "Builds" en `index.html`.
- Registrar en `main.ts`.
- Importar `initBuilds()`.

## Dependencias
- Requiere completar **0012** (design consistency) antes para heredar los estilos base.
- Se puede desarrollar en paralelo con **0013** (comparator) ya que no comparten lógica.

## Criterios de éxito
- Al seleccionar un Pokémon, solo aparecen sus moves compatibles (sin moves de otros Pokémon).
- Los stats calculados coinciden con los de herramientas externas (Showdown, Serebii).
- La simulación de daño muestra rango min-max correcto incluyendo efectividad de tipo y STAB.
- `CalculateDamage` y `TypeEffectiveness` en Core son funciones puras con tests unitarios.
- El total de EVs no puede superar 510 (validación en frontend).
