# Simulación de N batallas completas con informe de resultados

**ID**: 0043-batch-battle-simulation
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Añadir en el simulador de batallas la opción de lanzar N simulaciones completas del combate (donde N lo introduce el usuario desde el frontend). El resultado es un informe detallado que muestra estadísticas agregadas de todas las ejecuciones: porcentaje de victorias de cada lado, media de turnos, daño medio infligido, etc.

## Capas afectadas

- **Core**: Nuevos tipos (`BattleReport`, `BattleSummary`) y función pura `SimulateMultipleBattles` que ejecuta N simulaciones y agrega los resultados en un informe estadístico.
- **Shell**: Sin cambios.
- **APP**: Nuevo binding `SimulateMultipleBattles(input FullBattleInput, n int)` que conecta frontend con Core.
- **Frontend**: Input numérico para N, botón para lanzar batch, y sección de renderizado del informe con estadísticas detalladas.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/battle.go` | modificar | Añadir tipos `BattleReport` y `BattleSummary`, y función `SimulateMultipleBattles(input FullBattleInput, n int, randSource func(int) int) BattleReport` |
| `core/battle_test.go` | modificar | Tests para `SimulateMultipleBattles` con distintos valores de N |
| `app/bindings.go` | modificar | Añadir binding `SimulateMultipleBattles(input FullBattleInput, n int) core.BattleReport` |
| `frontend/src/pages/builds.ts` | modificar | Añadir input numérico para N, botón de lanzar batch, función `simulateMultipleBattles()`, y renderizado del informe |
| `frontend/src/styles/_builds.scss` | modificar | Estilos para el informe de resultados batch (tarjetas de stats, barras de porcentaje, tabla resumen) |

## Plan de implementacion

1. **Core — Tipos**: Definir `BattleSummary` (resultado individual resumido: ganador, turnos, HP restante de cada lado) y `BattleReport` (agregado: total simulaciones, victorias atacante/defensor/empates, porcentajes, media de turnos, media de HP restante del ganador, distribución de turnos).
2. **Core — Función pura**: Implementar `SimulateMultipleBattles(input FullBattleInput, n int, randSource func(int) int) BattleReport`. Internamente ejecuta `SimulateFullBattle` N veces, recoge cada `BattleState` final, extrae estadísticas y las agrega en el `BattleReport`. Sin logs individuales (solo stats agregadas) para evitar consumo excesivo de memoria.
3. **Core — Tests**: Tests unitarios con randSource determinista: N=1 (caso base), N=100 (verificar que los porcentajes suman 100, turnos medios > 0), N=0 (edge case, report vacío).
4. **APP — Binding**: Exponer `SimulateMultipleBattles(input FullBattleInput, n int) core.BattleReport` en `bindings.go`. Usar `math/rand.Intn` como randSource. Validar que n >= 1 y n <= 10000 (límite razonable).
5. **Frontend — UI**: En la sección de batalla (fase idle), añadir un input numérico (min=1, max=10000, default=100) y un botón "Simular N batallas". Solo visible si ambos lados tienen al menos 1 movimiento.
6. **Frontend — Lógica**: Implementar `simulateMultipleBattles()` que lee el valor de N del input, construye el `FullBattleInput` (reutilizar lógica de `simulateFullBattle()`), invoca el binding y renderiza el informe.
7. **Frontend — Informe**: Renderizar una sección de resultados con: tarjetas de victoria (atacante %, defensor %, empates %), stat cards (media turnos, mediana turnos, min/max turnos), y una barra visual de distribución de victorias. Estilo coherente con el theme Game Boy existente.
8. **Estilos**: Añadir estilos para `.battle-report`, `.report-stat-card`, `.report-win-bar` en `_builds.scss`.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/battle_test.go` | `SimulateMultipleBattles` — N=1, N=100, N=0, verificar sumas de porcentajes, rangos de turnos |

## Criterios de aceptacion

- [ ] El usuario puede introducir un número N (1-10000) en un input numérico
- [ ] Al pulsar el botón se ejecutan N simulaciones completas en el backend
- [ ] Se muestra un informe con: victorias atacante (%), victorias defensor (%), empates (%)
- [ ] El informe incluye estadísticas de turnos: media, mínimo, máximo
- [ ] Se muestra una barra visual de distribución de victorias (atacante vs defensor vs empate)
- [ ] La UI permanece responsive durante la simulación (no se bloquea en valores altos de N)
- [ ] El informe mantiene el estilo visual coherente con el simulador existente (theme Game Boy)
- [ ] Edge cases: N=1 funciona, valores fuera de rango se rechazan con feedback al usuario

## Notas

- No guardar los logs individuales de cada batalla en el report — solo estadísticas agregadas. N=10000 con 200 turnos por batalla generaría demasiada memoria.
- Reutilizar `SimulateFullBattle` internamente para cada iteración. No duplicar lógica de combate.
- El `randSource` debe ser diferente en cada iteración (no usar la misma seed) para resultados variados.
- Considerar mostrar un indicador de progreso si N es grande, aunque Go es suficientemente rápido para que N=10000 tarde < 1s.
