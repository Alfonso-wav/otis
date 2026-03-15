# Mostrar todos los movimientos de la batalla por equipos en la terminal

**ID**: 0048-team-battle-detailed-log
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Actualmente el log de la batalla por equipos solo muestra el resultado de cada ronda ("[Ronda N] X venció a Y (HP restante: Z)"). Se necesita que la terminal muestre **todos los movimientos realizados** turno a turno dentro de cada ronda, de forma similar al log de batalla 1v1.

Cada `BattleState` en `TeamBattleState.Rounds` ya contiene un `Log []string` con los detalles turno a turno. Solo hay que:
1. Incluir esos logs detallados en el `TeamBattleState.Log` final (o exponer un log combinado).
2. Renderizar ese log detallado en el frontend.

## Capas afectadas

- **Core**: Modificar `SimulateTeamBattle` para incluir los logs de cada ronda individual en el log general del `TeamBattleState`.
- **Shell**: No afectada.
- **APP**: No afectada (los datos ya fluyen al frontend).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/battle.go` | modificar | En `SimulateTeamBattle`, concatenar `result.Log` de cada ronda al log general, con separadores de ronda |
| `core/battle_test.go` | modificar | Verificar que el log detallado incluye entradas de movimientos individuales |
| `frontend/src/pages/builds.ts` | modificar | Renderizar el log detallado con formato visual claro (separadores de ronda, indentación) |
| `frontend/src/styles/_builds.scss` | modificar | Estilos para sub-entradas de log y separadores de ronda |

## Plan de implementacion

1. En `SimulateTeamBattle`, después de cada `SimulateFullBattle`, insertar un encabezado de ronda y luego todas las entradas de `result.Log` en el log general.
2. Mantener la entrada resumen existente al final de cada ronda.
3. En el frontend, renderizar cada entrada del log como ya se hace, pero añadir clases CSS para distinguir encabezados de ronda vs movimientos individuales vs resumen.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/battle_test.go` | Que el log de TeamBattleState contiene entradas de movimientos individuales, no solo resúmenes de ronda |

## Criterios de aceptacion

- [x] El log de batalla por equipos muestra cada movimiento realizado en cada turno de cada ronda
- [x] Los movimientos se agrupan visualmente por rondas con separadores claros
- [x] El resumen de cada ronda (quién venció a quién) sigue visible
- [x] El log es scrollable si es muy largo

## Notas

Los datos ya están disponibles en `TeamBattleState.Rounds[i].Log`. El cambio principal es concatenarlos al log global y mejorar el renderizado.
