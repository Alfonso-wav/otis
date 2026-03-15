# Bug: no se pueden eliminar todos los equipos tras crearlos en Builds

**ID**: 0056-fix-team-delete-after-create
**Estado**: done
**Fecha**: 2026-03-15
**Depende de**: 0055-fix-team-delete-and-create-buttons

---

## Descripcion

En la sección "Builds > Mis Equipos", después de crear equipos nuevos, al intentar eliminarlos siempre queda uno que no se puede borrar. El botón de eliminar no responde para el último equipo restante. El usuario tiene que salir y volver a la sección (vaciarla) para que funcione de nuevo.

Esto ocurre **después** del fix de la tarea 0055, que corrigió el backend no-idempotente y movió el refresh de la UI a un bloque `finally`. El problema actual es diferente: los equipos SÍ existen en disco (fueron recién creados), pero la UI deja de responder para el último equipo.

## Hipótesis

1. **Race condition en eliminación secuencial**: al eliminar equipos rápidamente, `handleDeleteTeam` es `async` y no tiene protección contra invocaciones concurrentes. Si el usuario hace clic en "Eliminar" de un equipo mientras el `buildLayout()` de la eliminación anterior aún no terminó, los event listeners del DOM podrían quedar desincronizados.

2. **Re-binding incompleto**: `buildLayout()` reemplaza `container.innerHTML` y luego llama a `bindTeamEvents()`. Si `cachedTeams` tiene exactamente 1 equipo al momento del render, podría haber un caso borde en `renderTeamsSection()` o `bindTeamEvents()` que no bindea correctamente el botón de eliminar cuando solo queda un equipo.

3. **Estado `teamsDetailsOpen` o `<details>` toggle**: el `<details>` element tiene un listener de `toggle` que actualiza `teamsDetailsOpen`. Si al reconstruir el DOM el `<details>` se cierra momentáneamente, el estado podría quedar inconsistente y el botón no ser accesible.

## Capas afectadas

| Capa | Impacto |
|------|---------|
| **Core** | Sin cambios |
| **Shell** | Sin cambios (ya es idempotente desde 0055) |
| **APP** | Sin cambios |
| **Frontend** | Principal: `frontend/src/pages/builds.ts` — handlers de eventos, ciclo de render/bind |

## Archivos a revisar/modificar

| Archivo | Acción |
|---------|--------|
| `frontend/src/pages/builds.ts` | Investigar `handleDeleteTeam` (L938-948), `bindTeamEvents` (L1192-1371), `renderTeamsSection` (L1110-1190), `buildLayout` (L1375-1464) |

## Plan de implementacion

1. **Reproducir el bug**: crear 2-3 equipos, intentar eliminarlos todos en secuencia.
2. **Añadir logs temporales** en `handleDeleteTeam`, `bindTeamEvents` y `buildLayout` para rastrear el flujo de ejecución.
3. **Verificar race conditions**: añadir un flag `isDeleting` que prevenga invocaciones concurrentes de `handleDeleteTeam`:
   ```typescript
   let isDeleting = false;
   async function handleDeleteTeam(name: string): Promise<void> {
     if (isDeleting) return;
     if (!confirm(`Eliminar equipo "${name}"?`)) return;
     isDeleting = true;
     try {
       await DeleteTeam(name);
     } catch (err: unknown) {
       console.warn(`DeleteTeam error (ignored): ${String(err)}`);
     } finally {
       teamsDetailsOpen = true;
       cachedTeams = await ListTeams();
       buildLayout();
       isDeleting = false;
     }
   }
   ```
4. **Verificar binding con 1 equipo**: inspeccionar el DOM tras eliminar hasta dejar 1 equipo y confirmar que `.team-delete-btn` tiene event listener.
5. **Verificar estado del `<details>`**: confirmar que `teamsDetailsOpen` se mantiene como `true` durante todo el ciclo de eliminación.

## Tests

- Crear 3 equipos → eliminarlos todos uno a uno → verificar que se eliminan todos (0 equipos).
- Crear 1 equipo → eliminarlo → verificar que funciona.
- Eliminar equipos rápidamente (clic rápido) → verificar que no queda ninguno.

## Criterios de aceptacion

- [x] Se pueden eliminar TODOS los equipos creados, incluido el último, sin necesidad de recargar la sección.
- [x] No hay race conditions al eliminar equipos en secuencia rápida.
- [x] El botón de eliminar responde siempre, independientemente del número de equipos restantes.
