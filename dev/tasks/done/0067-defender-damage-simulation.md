# Simulacion de dano del defensor en la pestana Builds

**ID**: 0067-defender-damage-simulation
**Estado**: todo
**Fecha**: 2026-03-16

---

## Descripcion

En la pestana Builds, la seccion "Simulacion de dano" (`renderDamageSection`) solo muestra el dano que los movimientos del atacante hacen al defensor. Se necesita mostrar tambien una segunda tabla con el dano que los movimientos del defensor hacen al atacante, para que el usuario tenga una vision completa del matchup.

## Capas afectadas

- **Core**: Sin cambios. `CalculateDamage()` ya soporta cualquier combinacion atacante/defensor.
- **Shell**: Sin cambios.
- **APP**: Solo frontend (TypeScript). El backend (`SimulateDamage` en `app/bindings.go`) ya acepta parametros intercambiados.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | Extender `renderDamageSection()` para incluir una segunda tabla con el dano del defensor. Extender o crear funcion complementaria a `loadDamageTable()` que itere sobre `state.defenderSlots`, llame a `SimulateDamage()` con stats/tipos intercambiados (defensor como atacante y viceversa), y renderice los resultados en la nueva tabla. |

## Plan de implementacion

1. **Modificar `renderDamageSection()` (~linea 310)**:
   - Anadir un segundo contenedor HTML para la tabla de dano del defensor debajo de la tabla existente del atacante.
   - Anadir titulo "Dano del Atacante" a la tabla existente y "Dano del Defensor" a la nueva tabla.
   - La seccion del defensor solo se muestra si el defensor tiene al menos un movimiento asignado en `state.defenderSlots`.

2. **Crear funcion `loadDefenderDamageTable()` o extender `loadDamageTable()`**:
   - Filtrar `state.defenderSlots` para obtener los slots con movimiento asignado.
   - Para cada movimiento del defensor, llamar a `SimulateDamage()` con los parametros intercambiados:
     - `attackerStats` = stats del defensor (con sus EVs/IVs y nivel)
     - `defenderStats` = stats del atacante
     - `attackerTypes` = tipos del defensor
     - `defenderTypes` = tipos del atacante
     - `level` = `state.defenderLevel`
   - Renderizar los resultados en la tabla del defensor con el mismo formato (Move, Type, Category, Min, Max, Effectiveness).

3. **Actualizar llamadas**:
   - Asegurar que cuando se cambian los movimientos del defensor o se recarga la simulacion, la tabla del defensor se actualice tambien.

4. **Estilo**:
   - Usar el mismo estilo de tabla que la existente para mantener coherencia visual.
   - Separar visualmente ambas tablas con un encabezado claro para cada una.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual | Al seleccionar atacante, defensor y movimientos de ambos, se muestran dos tablas de dano |
| Test manual | La tabla del atacante muestra el dano de sus movimientos contra el defensor (comportamiento existente) |
| Test manual | La tabla del defensor muestra el dano de sus movimientos contra el atacante (nueva funcionalidad) |
| Test manual | Si el defensor no tiene movimientos asignados, solo se muestra la tabla del atacante |
| Test manual | Los valores de efectividad, STAB y dano min/max son correctos en ambas tablas |
| Test manual | Al cambiar movimientos del defensor, la tabla del defensor se actualiza |

## Criterios de aceptacion

- [ ] Se muestra una tabla "Dano del Atacante" con el dano de los movimientos del atacante al defensor
- [ ] Se muestra una tabla "Dano del Defensor" con el dano de los movimientos del defensor al atacante
- [ ] La tabla del defensor solo aparece si el defensor tiene al menos un movimiento asignado
- [ ] Los calculos usan los stats, tipos y nivel correctos (intercambiados)
- [ ] Ambas tablas muestran movimiento, tipo, categoria, dano min/max, efectividad y STAB
- [ ] La UI es coherente visualmente entre ambas tablas

## Notas

- No se requieren cambios en el backend: `SimulateDamage()` ya acepta cualquier combinacion de atacante/defensor.
- Reutilizar al maximo la logica existente de `loadDamageTable()` para evitar duplicacion.
- Considerar extraer la logica de renderizado de tabla a una funcion comun que reciba los slots y los parametros de calculo.
