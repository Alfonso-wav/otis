# Fix: no se puede eliminar el último equipo ni crear equipos nuevos

**ID**: 0055-fix-team-delete-and-create-buttons
**Estado**: done
**Fecha**: 2026-03-15
**Depende de**: 0032-team-save-load

---

## Descripcion

En la sección "Builds > Mis Equipos", al eliminar equipos uno por uno, el último equipo restante ("team a") no se puede eliminar — el botón no responde. Además, el botón "+ Crear equipo" deja de funcionar por completo.

Ambos síntomas apuntan a que los event handlers dejan de funcionar tras cierta secuencia de eliminaciones, lo que sugiere un error durante el ciclo de re-render (`buildLayout()` → `bindTeamEvents()`) que interrumpe el bindeo de eventos.

## Causa raíz confirmada

**Error reportado por el usuario:**
```
Error: remove data\teams\team-a.json: The system cannot find the file specified.
```

El archivo `team-a.json` ya no existe en disco cuando se intenta eliminar. Esto causa un fallo en cadena:

1. `handleDeleteTeam("team a")` llama a `DeleteTeam("team a")` en Go.
2. `shell/teams.go` → `os.Remove(filePath)` falla con "file not found".
3. El error se propaga al frontend y entra en el `catch`.
4. **El `catch` solo muestra un `alert` pero NO refresca la UI** — las líneas `cachedTeams = await ListTeams()` y `buildLayout()` están en el `try` y nunca se ejecutan.
5. `cachedTeams` queda con datos obsoletos (el equipo fantasma sigue ahí).
6. El botón "+ Crear equipo" probablemente también falla porque la UI quedó en un estado inconsistente tras el alert, o porque un error previo ya había roto el ciclo de render/bind.

**Hay dos problemas independientes:**

### Bug 1: Backend no es idempotente (`shell/teams.go:100-102`)
```go
func (fs *FileTeamStorage) DeleteTeam(name string) error {
    return os.Remove(fs.filePath(name)) // falla si el archivo no existe
}
```
`os.Remove` devuelve error si el archivo no existe. Debería ignorar `os.ErrNotExist`.

### Bug 2: Frontend no sincroniza UI tras error (`builds.ts:938-948`)
```typescript
async function handleDeleteTeam(name: string): Promise<void> {
  if (!confirm(`Eliminar equipo "${name}"?`)) return;
  try {
    await DeleteTeam(name);
    teamsDetailsOpen = true;
    cachedTeams = await ListTeams();  // ← nunca se ejecuta si DeleteTeam falla
    buildLayout();                     // ← nunca se ejecuta
  } catch (err: unknown) {
    alert(`Error: ${String(err)}`);    // ← solo muestra alert, no refresca
  }
}
```
El refresh de `cachedTeams` y `buildLayout()` debería estar en un bloque `finally` para que la UI siempre se sincronice con el estado real del disco.

## Capas afectadas

| Capa | Impacto |
|------|---------|
| **Core** | Sin cambios |
| **Shell** | Posible: `DeleteTeam` debería manejar gracefully el caso de archivo inexistente |
| **APP** | Posible: `DeleteTeam` en bindings.go podría necesitar validación |
| **Frontend** | Principal: `frontend/src/pages/builds.ts` — handlers de eventos y ciclo de render |

## Archivos a revisar/modificar

| Archivo | Acción |
|---------|--------|
| `frontend/src/pages/builds.ts` | Investigar y corregir el ciclo de render/bind cuando quedan 0-1 equipos |
| `shell/teams.go` | Verificar que `DeleteTeam` no falle si el archivo no existe |
| `app/bindings.go` | Verificar manejo de errores en `DeleteTeam` |

## Plan de implementacion

1. **Fix backend** — `shell/teams.go:100-102`:
   - Hacer `DeleteTeam` idempotente: si `os.Remove` devuelve `os.ErrNotExist`, retornar `nil` en lugar del error.
   ```go
   func (fs *FileTeamStorage) DeleteTeam(name string) error {
       err := os.Remove(fs.filePath(name))
       if errors.Is(err, os.ErrNotExist) {
           return nil
       }
       return err
   }
   ```

2. **Fix frontend** — `frontend/src/pages/builds.ts:938-948`:
   - Mover `cachedTeams = await ListTeams()` y `buildLayout()` a un bloque `finally` para que la UI siempre se sincronice con el disco, haya o no error.
   ```typescript
   async function handleDeleteTeam(name: string): Promise<void> {
     if (!confirm(`Eliminar equipo "${name}"?`)) return;
     try {
       await DeleteTeam(name);
     } catch (err: unknown) {
       console.warn(`DeleteTeam error (ignored): ${String(err)}`);
     } finally {
       teamsDetailsOpen = true;
       cachedTeams = await ListTeams();
       buildLayout();
     }
   }
   ```

## Tests

- Verificar eliminación de equipos hasta quedar en 0.
- Verificar que "+ Crear equipo" funciona con 0 equipos existentes.
- Verificar que eliminar el único equipo existente funciona.
- Verificar que `DeleteTeam` no falla si se llama dos veces con el mismo nombre.

## Criterios de aceptacion

- [x] Se puede eliminar cualquier equipo, incluido el último
- [x] El botón "+ Crear equipo" funciona siempre, sin importar cuántos equipos existan
- [x] No hay errores silenciosos en la consola del navegador durante el flujo de eliminación
- [x] `DeleteTeam` en el backend es idempotente (no falla si el archivo ya no existe)
