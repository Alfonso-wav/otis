# Fix: botón de simulación masiva desaparece tras iniciar una batalla

**ID**: 0044-fix-batch-button-disappears
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

El botón "Simular N batallas" (simulación masiva/batch) desaparece después de iniciar una batalla turno a turno y no vuelve a aparecer, ni siquiera al pulsar "Reiniciar". La causa raíz es que el botón "Reiniciar" llama a `startBattle()`, que establece `battleUI.phase = "attacker-turn"`, pero la sección batch solo se renderiza cuando `phase === "idle"`. No existe ningún camino en la UI para volver al estado idle una vez iniciada una batalla (salvo cambiar de Pokémon, lo que borra todos los movimientos).

## Causa raíz

En `frontend/src/pages/builds.ts`:
- **Línea 717**: `#battle-reset-btn` ejecuta `startBattle()` → `phase = "attacker-turn"`.
- **Líneas 397-424**: La fila batch (`battle-batch-row`) solo se renderiza si `phase === "idle"`.
- **Líneas 454-475**: La vista de batalla en progreso NO incluye la fila batch ni un botón para volver a idle.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Sin cambios.
- **Frontend**: Cambio en la lógica del botón "Reiniciar" para volver al estado idle en lugar de iniciar una nueva batalla.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | Cambiar el handler del botón "Reiniciar" para que vuelva al estado idle (`battleUI = { battleState: null, phase: "idle" }`) y re-renderice, en lugar de llamar a `startBattle()` |

## Plan de implementacion

1. En `bindBattleEvents()` (línea 717), cambiar el listener de `#battle-reset-btn` para que en lugar de llamar a `startBattle()`, resetee `battleUI` a estado idle:
   ```typescript
   container.querySelector<HTMLButtonElement>("#battle-reset-btn")?.addEventListener("click", () => {
     battleUI = { battleState: null, phase: "idle" };
     renderBattleInPlace();
   });
   ```
2. Limpiar también `batchReport` si se desea un reset completo (opcional, según preferencia UX).
3. Verificar que tras pulsar "Reiniciar", la sección idle se muestra con el botón batch y el botón de batalla turno a turno disponibles.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que tras iniciar una batalla y pulsar "Reiniciar", el botón "Simular N batallas" reaparece |
| Manual | Verificar que el flujo completo funciona: idle → batalla → reiniciar → idle (con batch visible) → batch simulation |

## Criterios de aceptacion

- [x] Al pulsar "Reiniciar" durante o tras una batalla, se vuelve al estado idle
- [x] En estado idle, el botón "Simular N batallas" es visible (si ambos lados tienen movimientos)
- [x] El botón "Iniciar batalla turno a turno" también es visible tras reiniciar
- [x] La arena de batalla muestra los sprites correctos en estado idle tras reiniciar
- [x] No hay regresiones en el flujo de batalla turno a turno ni en simulación completa

## Notas

- El cambio es mínimo: solo afecta al event listener del botón reset en `bindBattleEvents()`.
- Considerar si `batchReport` debe limpiarse al reiniciar o mantenerse visible como referencia del último batch ejecutado.
