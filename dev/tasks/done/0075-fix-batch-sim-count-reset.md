# 0075 — Fix: el número de simulaciones N se resetea a 100 tras cada ejecución

**Estado: done**

## Descripcion

En el simulador de batallas (individual y por equipos), el input del número de simulaciones masivas se resetea a `100` cada vez que se ejecuta una simulación. Esto ocurre porque `renderBattleInPlace()` y `buildLayout()` re-renderizan el HTML completo, recreando el `<input>` con `value="100"` hardcodeado. El valor introducido por el usuario se pierde.

**Comportamiento esperado**: tras ejecutar N simulaciones, el input debe conservar el último valor introducido por el usuario.

## Capas afectadas

- **APP (frontend)** — persistir el valor del input de simulaciones en el estado del módulo antes del re-render.

No se requieren cambios en Core ni Shell.

## Archivos a modificar

### Frontend

1. **`frontend/src/pages/builds.ts`**

#### Simulación individual (1v1)

- **Línea ~503**: el template HTML tiene `value="100"` hardcodeado.
- **Fix**: añadir una variable de estado del módulo (ej. `let lastBatchN = 100;`) y usarla en el template: `value="${lastBatchN}"`.
- **Línea ~696**: en `simulateBatchBattles()`, tras leer el valor del input, guardarlo en la variable de estado: `lastBatchN = n;`.

#### Simulación por equipos

- **Línea ~1156**: el template HTML tiene `value="100"` hardcodeado.
- **Fix**: añadir una variable de estado del módulo (ej. `let lastTeamBatchN = 100;`) y usarla en el template: `value="${lastTeamBatchN}"`.
- **Línea ~1567**: en el handler del click de `#tb-batch-btn`, tras leer el valor del input, guardarlo: `lastTeamBatchN = n;`.

## Requisitos

1. Ambos inputs (simulación individual y por equipos) deben conservar el último valor introducido tras cada ejecución.
2. El valor por defecto inicial sigue siendo `100`.
3. No se requiere persistencia entre sesiones (solo mantener el valor mientras el usuario está en la pestaña Builds).

## Tests

- Introducir un valor distinto de 100 (ej. 500), ejecutar la simulación, verificar que el input mantiene 500.
- Repetir para la simulación por equipos.
- Verificar que al entrar por primera vez a Builds, el valor por defecto es 100.
