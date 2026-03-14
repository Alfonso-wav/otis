# 0029 — Hacer visible el STAB en la UI y añadir modificador de quemadura

## Estado

done

## Descripción

El STAB (Same Type Attack Bonus ×1.5) está implementado en el cálculo de daño (`core/damage.go:153-160`) pero **no es visible para el usuario**: `DamageResult.Multiplier` solo contiene la efectividad de tipo, y el frontend no muestra ningún indicador de STAB. Esto hace que el usuario perciba que el STAB no se está aplicando.

Además, revisando `docs/pokemon_battle_mechanics.md`, el **modificador de quemadura** (§9.1: ×0.5 al daño físico si el atacante está quemado) es la siguiente mecánica pendiente de prioridad media que debería implementarse.

## Contexto

- `core/damage.go` calcula STAB correctamente (×1.5 si el tipo del movimiento coincide con algún tipo del atacante), pero no lo expone en `DamageResult`.
- `DamageResult.Multiplier` = solo `typeEff`. El frontend (`frontend/src/pages/builds.ts`) usa este campo para mostrar "¡Super eficaz! ×2.0" etc., pero nunca menciona STAB.
- El log de batalla (`core/battle.go`) tampoco indica cuándo hay STAB.
- La quemadura (Burn) no está implementada: no hay campo `IsBurned` ni en `DamageInput` ni lógica de ×0.5 para movimientos físicos.
- Referencia: `docs/pokemon_battle_mechanics.md` §5 (STAB), §9.1 (Quemadura).

## Capas involucradas

- **Core**: `core/damage.go` (exponer STAB en resultado, añadir burn modifier)
- **Core (tests)**: `core/damage_test.go` (tests para STAB visible y burn)
- **APP (bindings)**: `app/bindings.go` si cambian las firmas
- **Frontend**: `frontend/src/pages/builds.ts` (mostrar indicador STAB en tabla de daño y log), `frontend/wailsjs/go/models.ts` (regenerar tipos)
- **Frontend (estilos)**: `frontend/src/styles/_builds.scss` (estilo para badge STAB)

## Plan de implementación

### Paso 1 — Core: Exponer STAB en DamageResult (`core/damage.go`)

- Añadir campo `HasSTAB bool` y `STABMultiplier float64` a `DamageResult`.
- En `CalculateDamage`, setear `HasSTAB = true` y `STABMultiplier = 1.5` cuando aplique STAB.
- Mantener el campo `Multiplier` existente como está (solo typeEff) para no romper lógica del frontend.

### Paso 2 — Core: Añadir modificador de quemadura (`core/damage.go`)

- Añadir campo `IsBurned bool` a `DamageInput`.
- Si `IsBurned == true` y `Move.Category == "physical"`, aplicar multiplicador ×0.5 al daño.
- Añadir campo `BurnApplied bool` a `DamageResult` para que el frontend pueda mostrarlo.

### Paso 3 — Core: Tests (`core/damage_test.go`)

- Test: verificar que `HasSTAB == true` y `STABMultiplier == 1.5` cuando atacante y movimiento comparten tipo.
- Test: verificar que `HasSTAB == false` cuando no comparten tipo.
- Test: verificar que quemadura reduce daño físico a la mitad.
- Test: verificar que quemadura NO reduce daño especial.

### Paso 4 — Frontend: Mostrar STAB en tabla de daño (`frontend/src/pages/builds.ts`)

- En `loadDamageTable()`, añadir indicador visual (badge/texto "STAB") cuando `HasSTAB == true` en la celda de efectividad o como columna adicional.
- En el log de batalla, indicar "(STAB)" cuando el movimiento tiene bonus de mismo tipo.

### Paso 5 — Frontend: Mostrar estado de quemadura (si aplica)

- Si en el futuro se añade UI para condiciones de estado, el campo `BurnApplied` estará listo para mostrarse.
- Por ahora, el campo queda disponible para el log de batalla.

### Paso 6 — Regenerar bindings de Wails

- Ejecutar `wails generate module` para actualizar `frontend/wailsjs/go/models.ts` con los nuevos campos.

## Criterios de aceptación

- [x] `DamageResult` incluye `HasSTAB` y `STABMultiplier` con valores correctos
- [x] La tabla de daño en la UI muestra un indicador visual cuando un movimiento tiene STAB
- [x] El log de batalla indica "(STAB)" cuando aplica
- [x] La quemadura reduce ×0.5 el daño de movimientos físicos cuando `IsBurned == true`
- [x] La quemadura NO afecta movimientos especiales
- [x] `DamageResult` incluye `BurnApplied` para indicar cuándo se aplicó la reducción
- [x] Tests unitarios cubren todos los casos de STAB visibilidad y quemadura
- [x] Los tests existentes siguen pasando
- [x] Los bindings de Wails están actualizados

## Notas

- El STAB ya funciona correctamente en el cálculo — esto es un fix de **visibilidad/UX**, no de lógica.
- La quemadura es la primera condición de estado con impacto en daño. Prepara el terreno para futuras condiciones (parálisis, etc.).
- No implementar Guts, Facade ni otras excepciones a quemadura en esta tarea — quedan para cuando se implementen habilidades.
- Referencia: `docs/pokemon_battle_mechanics.md` §5, §9.1, §9.3.
