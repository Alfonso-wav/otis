# 0016 — Simulación de batalla por turnos en la pestaña Builds

## Estado

done

## Descripción

Añadir una sección de **simulación de batalla por turnos** dentro del tab **Builds** (`builds.ts`). Una vez que el usuario ha configurado su atacante y defensor (con IVs, EVs, naturaleza y movimientos), podrá iniciar una batalla simulada donde se eligen movimientos por turno, se actualizan las HP y se declara un ganador.

## Contexto

- El tab Builds ya tiene configuración completa de atacante/defensor con IVs, EVs, naturaleza y 4 slots de movimiento.
- El backend ya expone `SimulateDamage(input)` → `DamageResult` (min/max/avg damage, multiplicador, flags de efectividad).
- El backend ya expone `CalculateStats(pokemon, ivs, evs, nature, level)` → stats finales.
- El core ya tiene `CalculateDamage()` en `core/damage.go` con toda la mecánica gen 5+.
- **Lo que falta**: lógica de batalla multi-turno (HP pool, selección de movimiento, turnos alternados, PP implícitos, victoria/derrota).

## Capas involucradas

- **Core**: añadir `BattleState`, `BattleResult` y función pura `SimulateTurn()` en `core/battle.go`
- **APP (bindings)**: exponer `StartBattle()` y `ExecuteTurn()` en `app/bindings.go`
- **Frontend (APP)**: añadir sección de batalla en `builds.ts` con HP bars, log de turnos y botones de movimiento

## Plan de implementación

### Paso 1 — Core: tipos y lógica de batalla (`core/battle.go`)

Crear `core/battle.go` con:

```go
type BattleState struct {
    AttackerHP    int
    DefenderHP    int
    AttackerMaxHP int
    DefenderMaxHP int
    TurnCount     int
    Log           []string
    IsOver        bool
    Winner        string // "attacker" | "defender" | ""
}

type TurnInput struct {
    State         BattleState
    AttackerStats Stats
    DefenderStats Stats
    AttackerTypes []string
    DefenderTypes []string
    AttackerLevel int
    DefenderLevel int
    Move          Move
}

type TurnResult struct {
    NewState BattleState
    Damage   DamageResult
    LogEntry string
}

func InitBattle(attackerMaxHP, defenderMaxHP int) BattleState
func ExecuteTurn(input TurnInput) TurnResult
```

- `ExecuteTurn` es una función pura: recibe el estado actual + input del turno, devuelve nuevo estado + resultado
- Aplica el daño al HP del defensor; si HP ≤ 0, marca `IsOver = true` y `Winner = "attacker"`
- Añade entrada al log: `"[T1] Pikachu usó Thunderbolt → 45 daño (x2.0) | HP Charmander: 120/165"`
- No hay AI para el defensor en esta iteración: el usuario controla ambos lados
- Añadir tests en `core/battle_test.go`

### Paso 2 — APP: bindings (`app/bindings.go`)

Añadir dos métodos al struct `App`:

```go
func (a *App) InitBattle(attackerMaxHP, defenderMaxHP int) core.BattleState
func (a *App) ExecuteTurn(input core.TurnInput) core.TurnResult
```

Ambos son thin wrappers sobre las funciones de core.

### Paso 3 — Frontend: sección de batalla en `builds.ts`

Añadir debajo de la tabla de daño actual una sección colapsable `#battle-section` que se activa cuando ambos Pokémon están configurados y hay al menos un movimiento seleccionado.

**UI de batalla:**
- Botón "⚔ Iniciar batalla" que llama a `InitBattle()` con los HP máximos calculados
- Dos barras de HP (una por Pokémon) con nombre, sprite pequeño y valor `HP actual / HP máx`
- Cuatro botones de movimiento (uno por slot configurado) — sólo activos si tienen PP > 0 (simplificado: PP = potencia estimada, no implementar PP reales todavía)
- Log de batalla scrollable (últimos 10 turnos) con entradas coloreadas por efectividad
- Indicador de turno actual
- Botón "Reiniciar" que llama a `InitBattle()` de nuevo

**Flujo:**
1. Clic en movimiento → llama `ExecuteTurn()` con el estado actual y el movimiento elegido
2. Aplica daño al defensor (HP bar animada)
3. El usuario elige el movimiento del defensor para el contra-ataque
4. Se actualiza log y estado
5. Si `IsOver`, mostrar banner de ganador y deshabilitar botones

**Estado local del frontend:**
```typescript
interface BattleState {
  state: core.BattleState | null;
  phase: "idle" | "attacker-turn" | "defender-turn" | "over";
}
```

### Paso 4 — Estilos SCSS

Añadir en el archivo de estilos existente (o en uno nuevo si hay módulo SCSS por tab):
- `.battle-section`: contenedor con borde separador del resto del tab
- `.hp-bar`: barra de progreso con color verde→amarillo→rojo según porcentaje
- `.battle-log`: área scrollable de 150px con fondo oscuro y texto monoespaciado
- `.move-btn`: botones de movimiento con color por tipo, deshabilitado si batalla terminada
- `.battle-winner-banner`: overlay o banner visible cuando `IsOver = true`

## Criterios de aceptación

- [ ] Con atacante y defensor configurados, aparece botón "Iniciar batalla"
- [ ] Las HP se inicializan correctamente usando `CalculateStats`
- [ ] Clic en movimiento del atacante aplica daño al HP del defensor con animación en la barra
- [ ] El log muestra: turno, movimiento usado, daño, multiplicador de tipo y HP restante
- [ ] Tras el turno del atacante, se activan los botones del defensor
- [ ] Cuando HP ≤ 0 la batalla termina y se muestra el ganador
- [ ] Botón "Reiniciar" restaura las HP sin necesidad de reconfigurar el build
- [ ] `ExecuteTurn` en Core es una función pura con tests unitarios

## Notas

- No implementar AI para el defensor ni gestión de PP reales: queda para iteración futura
- Los movimientos sin power (status moves) aplican 0 daño y aparecen en el log como "sin efecto de daño"
- El turno siempre es: atacante primero, luego defensor (sin cálculo de velocidad por ahora)
- No modificar la lógica de la tabla de daño existente
