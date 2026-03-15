# Simulador de combates completos entre equipos

**ID**: 0047-team-battle-simulator
**Estado**: done
**Fecha**: 2026-03-15
**Depende de**: 0045-create-team-from-panel, 0046-random-team-fill

---

## Descripcion

Añadir al simulador de batallas un modo de combate de equipos completos (6v6 o NvN). Dos equipos seleccionados del panel "Mis Equipos" se enfrentan en una serie de combates 1v1 secuenciales: cuando un Pokémon cae, el siguiente del equipo entra. El equipo que agote todos los Pokémon del rival gana. Se soportará simulación automática completa y batch (N simulaciones con estadísticas).

## Capas afectadas

- **Core**: Añadir structs `TeamBattleInput`, `TeamBattleState`, `TeamBattleReport`. Implementar `SimulateTeamBattle` y `SimulateMultipleTeamBattles`.
- **Shell**: Sin cambios.
- **APP**: Añadir bindings `SimulateTeamBattle` y `SimulateMultipleTeamBattles` que resuelvan stats de cada miembro y llamen a core.
- **Frontend**: Añadir sección de selección de equipos y UI de batalla de equipos con progreso visual.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir tipos `TeamBattleInput`, `TeamBattleMember`, `TeamBattleState`, `TeamBattleReport` |
| `core/battle.go` | modificar | Añadir `SimulateTeamBattle(input TeamBattleInput, rng func(int) int) TeamBattleState` — simula combate completo de equipos |
| `core/battle.go` | modificar | Añadir `SimulateMultipleTeamBattles(input TeamBattleInput, n int, rng func(int) int) TeamBattleReport` — N simulaciones con estadísticas |
| `app/bindings.go` | modificar | Añadir `SimulateTeamBattle(team1Name, team2Name string) (core.TeamBattleState, error)` y `SimulateMultipleTeamBattles(team1Name, team2Name string, n int) (core.TeamBattleReport, error)` |
| `frontend/src/pages/builds.ts` | modificar | Añadir sección "Batalla de equipos" con selectores de equipo, botón simular, visualización de progreso (qué Pokémon quedan vivos por equipo), log de combate, y reporte batch |
| `frontend/src/types.ts` | modificar | Añadir interfaces TypeScript para los nuevos tipos de batalla de equipos |

## Plan de implementacion

### 1. Core — Tipos

```go
type TeamBattleMember struct {
    Pokemon   Pokemon     `json:"pokemon"`
    Member    TeamMember  `json:"member"`
    Stats     Stats       `json:"stats"`
    Types     []PokemonType `json:"types"`
}

type TeamBattleInput struct {
    Team1Name    string              `json:"team1Name"`
    Team1Members []TeamBattleMember  `json:"team1Members"`
    Team2Name    string              `json:"team2Name"`
    Team2Members []TeamBattleMember  `json:"team2Members"`
}

type TeamBattleState struct {
    Team1Remaining int        `json:"team1Remaining"`
    Team2Remaining int        `json:"team2Remaining"`
    TotalTurns     int        `json:"totalTurns"`
    Rounds         []BattleState `json:"rounds"` // cada combate 1v1
    Log            []string   `json:"log"`
    IsOver         bool       `json:"isOver"`
    Winner         string     `json:"winner"` // "team1" | "team2"
}

type TeamBattleReport struct {
    TotalSimulations int     `json:"totalSimulations"`
    Team1Wins        int     `json:"team1Wins"`
    Team2Wins        int     `json:"team2Wins"`
    Draws            int     `json:"draws"`
    Team1WinPct      float64 `json:"team1WinPct"`
    Team2WinPct      float64 `json:"team2WinPct"`
    DrawPct          float64 `json:"drawPct"`
    AvgTotalTurns    float64 `json:"avgTotalTurns"`
    AvgTeam1Remaining float64 `json:"avgTeam1Remaining"`
    AvgTeam2Remaining float64 `json:"avgTeam2Remaining"`
}
```

### 2. Core — Lógica de batalla de equipos

`SimulateTeamBattle`:
- Mantener índice activo por equipo (empieza en 0).
- Resolver combate 1v1 entre los Pokémon activos usando `SimulateFullBattle` existente.
- Cuando uno cae, avanzar índice del equipo que perdió. El ganador mantiene su HP restante.
- Repetir hasta que un equipo se quede sin Pokémon.
- El HP restante del ganador del 1v1 anterior se arrastra al siguiente combate (carry-over HP).

### 3. APP — Bindings

- Resolver cada `TeamMember` a `TeamBattleMember` calculando stats reales con `CalculateStats`.
- Obtener tipos y movimientos completos de cada Pokémon.
- Pasar a core y retornar resultado.

### 4. Frontend — UI

- **Selector de equipos**: dos dropdowns con los equipos disponibles (mínimo 1 miembro cada uno).
- **Vista previa**: mostrar sprites de ambos equipos lado a lado.
- **Botón "Simular batalla de equipos"**: ejecuta una simulación completa.
- **Progreso visual**: mostrar qué Pokémon siguen vivos (sprites en color) y cuáles cayeron (sprites en gris/tachados).
- **Log**: resumen de cada ronda (quién venció a quién, HP restante).
- **Batch**: input numérico + botón para N simulaciones con reporte estadístico.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/battle_test.go` | `SimulateTeamBattle` con equipos de 1, 3 y 6 miembros — verificar que siempre hay ganador |
| `core/battle_test.go` | `SimulateTeamBattle` con carry-over HP — el ganador de un 1v1 conserva HP |
| `core/battle_test.go` | `SimulateMultipleTeamBattles` — verificar estadísticas consistentes |
| Manual | Flujo completo UI: seleccionar equipos → simular → ver progreso → batch |

## Criterios de aceptacion

- [x] Se pueden seleccionar dos equipos para enfrentarse
- [x] La simulación resuelve combates 1v1 secuenciales hasta que un equipo se queda sin Pokémon
- [x] El HP del ganador de cada 1v1 se conserva para el siguiente combate
- [x] Se muestra visualmente qué Pokémon quedan vivos y cuáles cayeron
- [x] El log muestra el desarrollo de cada ronda
- [x] La simulación batch funciona y muestra estadísticas (win%, turnos promedio, Pokémon restantes promedio)
- [x] Funciona con equipos de cualquier tamaño (1-6 miembros)
- [x] No hay regresiones en el simulador de batallas 1v1 existente

## Notas

- Reutilizar `SimulateFullBattle` existente para cada combate 1v1 — no reimplementar lógica de batalla.
- El carry-over de HP es clave para que la simulación sea realista (el que gana el primer combate entra herido al siguiente).
- Considerar que los equipos pueden tener distinto número de miembros (3v6 es válido).
