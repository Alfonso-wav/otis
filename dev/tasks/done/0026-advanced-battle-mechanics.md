# 0026 — Implementar mecánicas avanzadas de batalla desde el documento de referencia

## Estado

done

## Descripción

Mejorar el simulador de batallas (tanto simulación de daño como batalla por turnos) implementando mecánicas adicionales documentadas en `docs/pokemon_battle_mechanics.md` que actualmente no están cubiertas. El objetivo es hacer que la simulación sea más fiel a las mecánicas reales de Pokémon Gen V+.

## Contexto

- `core/damage.go` implementa la fórmula base Gen V+ con STAB, efectividad de tipo, crítico y weather. Pero le faltan varias mecánicas documentadas.
- `core/battle.go` implementa `ExecuteTurn` y `SimulateFullBattle` pero sin variación aleatoria en daño, sin cálculo de velocidad/prioridad, sin precisión/evasión y sin condiciones de estado.
- El documento `docs/pokemon_battle_mechanics.md` (844 líneas) cubre exhaustivamente todas las mecánicas de Gen I-IX.

### Mecánicas que faltan por implementar (priorizadas)

**Prioridad alta** (impacto directo en fidelidad de simulación):
1. **Random roll** (§8): factor aleatorio 0.85–1.00 en cada ataque. Actualmente el daño es determinístico (usa average).
2. **Golpe crítico probabilístico** (§7): probabilidad de 1/24 por defecto. Actualmente `IsCritical` siempre es `false`.
3. **Orden de turno por velocidad** (§17): el más rápido ataca primero. Actualmente siempre ataca primero el atacante.
4. **Precisión/evasión** (§4.3, Apéndice B): los movimientos pueden fallar según su precisión base. Actualmente todos los movimientos aciertan siempre.

**Prioridad media** (mejoran realismo pero son más complejos):
5. **Quemadura** (§9.1): reducción ×0.5 al daño físico del Pokémon quemado.
6. **Etapas de estadísticas** (§4): modificadores -6 a +6 en stats de combate para movimientos que las modifican.
7. **Movimientos de prioridad** (§17.1): Aqua Jet, Quick Attack, etc. van antes independientemente de velocidad.

**Prioridad baja** (para futuras iteraciones):
8. Clima y terreno dinámico
9. Habilidades que modifican daño
10. Objetos equipados
11. Pantallas defensivas (Reflect, Light Screen)
12. Condiciones de estado completas (parálisis, sueño, congelación)

## Capas involucradas

- **Core**: `core/damage.go` (añadir random roll, critical hit probabilístico, precisión), `core/battle.go` (velocidad, orden de turno)
- **Core (nuevos)**: posiblemente `core/accuracy.go`, `core/speed.go`
- **APP (bindings)**: `app/bindings.go` si cambian firmas de funciones
- **Frontend**: `frontend/src/pages/builds.ts` para mostrar info adicional (críticos, fallos, orden)

## Plan de implementación

### Paso 1 — Core: Random roll en cálculo de daño (`core/damage.go`)

Modificar `CalculateDamage` para aceptar un `randSource func(n int) int` inyectable (igual que `SimulateFullBattle`):
- Generar un entero aleatorio entre 85 y 100 (16 valores posibles)
- Aplicar como multiplicador `roll/100.0` al daño base
- `DamageResult` sigue devolviendo min/max/average para la tabla de daño (sin random), pero añadir un campo `ActualDamage` que se usa en batallas
- Si `randSource` es nil, comportamiento actual (determinístico con min/max/average)

### Paso 2 — Core: Golpe crítico probabilístico (`core/damage.go`)

- Añadir campo `CriticalStage int` a `DamageInput` (default 0)
- Calcular probabilidad: stage 0 → 1/24, stage 1 → 1/8, stage 2 → 1/2, stage ≥3 → 1/1
- Usar `randSource` para determinar si es crítico
- Aplicar multiplicador ×1.5 si es crítico
- Retornar flag `WasCritical bool` en `DamageResult`

### Paso 3 — Core: Precisión y fallo de movimientos (`core/battle.go`)

- Antes de calcular daño en `ExecuteTurn`, verificar si el movimiento acierta
- Fórmula: `P = accuracy_base × (stages_accuracy / stages_evasion)`
- Si el movimiento tiene accuracy "—" (never miss), siempre acierta
- Si falla, el log muestra "¡El ataque falló!" y no se aplica daño
- Añadir campo `Accuracy int` al tipo `Move` si no existe (0 o 100 = always hits)

### Paso 4 — Core: Orden de turno por velocidad (`core/battle.go`)

- Modificar `SimulateFullBattle` y la lógica de batalla para que el Pokémon con mayor Speed ataque primero cada turno
- Si hay empate de velocidad, resolver al azar
- Aplicar bracket de prioridad del movimiento si existe (campo `Priority int` en `Move`)
- Esto requiere que ambos movimientos se seleccionen antes de ejecutar el turno

### Paso 5 — Frontend: mostrar info adicional en batalla

- En el log, indicar si un ataque fue crítico: `"¡Golpe crítico! [T3] usó..."`
- En el log, indicar si un ataque falló: `"[T3] usó Trueno → ¡Falló!"`
- Mostrar quién ataca primero cada turno basándose en velocidad
- Añadir indicador de golpe crítico en la UI (flash o color diferente)

### Paso 6 — Tests

- Actualizar `core/damage_test.go` con tests para random roll y critical hit probabilístico
- Actualizar `core/battle_test.go` con tests para accuracy check y speed ordering
- Usar `randSource` inyectado para tests determinísticos

## Criterios de aceptación

- [ ] El daño en batalla varía entre min y max (random roll 0.85–1.00)
- [ ] Los golpes críticos ocurren con probabilidad ~4.17% y hacen ×1.5 daño
- [ ] Los movimientos con precisión < 100% pueden fallar y se muestra en el log
- [ ] El Pokémon más rápido ataca primero cada turno
- [ ] Los movimientos con prioridad (Quick Attack, etc.) respetan su bracket
- [ ] La tabla de daño existente (min/max/avg) sigue funcionando sin cambios para la preview
- [ ] Todos los tests existentes siguen pasando + nuevos tests para las mecánicas añadidas
- [ ] La simulación completa (`SimulateFullBattle`) usa todas las mecánicas nuevas

## Notas

- Implementar mecánicas de prioridad alta primero (pasos 1-4). Las de prioridad media/baja quedan para futuras tareas.
- Mantener backward compatibility: la tabla de daño de preview sigue mostrando min/max/avg determinístico.
- `randSource` se inyecta para testabilidad, en producción se usa `rand.Intn`.
- No implementar habilidades, objetos ni condiciones de estado en esta tarea — son demasiado complejas y quedan para iteraciones futuras.
- Referencia principal: `docs/pokemon_battle_mechanics.md` secciones 1, 4, 7, 8, 17.
