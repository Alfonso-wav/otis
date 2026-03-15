# Crear equipo nuevo desde el panel "Mis Equipos"

**ID**: 0045-create-team-from-panel
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Actualmente los equipos se crean implícitamente al guardar un Pokémon en un equipo nuevo desde el modal de guardado. Se quiere añadir un botón explícito "Crear equipo" en el panel de "Mis Equipos" (pestaña Build) que permita crear un equipo vacío con un nombre, para luego ir añadiéndole miembros.

## Capas afectadas

- **Core**: Sin cambios (ya existe `ValidateTeam` que acepta equipos con 0 miembros válidos).
- **Shell**: Sin cambios (`SaveTeam` ya persiste cualquier equipo válido).
- **APP**: Añadir binding `CreateTeam(name string) error` que cree un equipo vacío.
- **Frontend**: Añadir botón y mini-formulario en el panel de equipos.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `app/bindings.go` | modificar | Añadir método `CreateTeam(name string) error` que valide nombre y cree equipo vacío via `teamStorage.SaveTeam` |
| `frontend/src/pages/builds.ts` | modificar | Añadir botón "Crear equipo" en `renderTeamsSection()` con input de nombre inline. Al confirmar, llamar a `CreateTeam` y re-renderizar la lista |

## Plan de implementacion

1. **Backend** — En `app/bindings.go`, añadir:
   ```go
   func (a *App) CreateTeam(name string) error {
       team := core.Team{Name: name, Members: []core.TeamMember{}}
       if err := core.ValidateTeam(team); err != nil {
           return err
       }
       return a.teamStorage.SaveTeam(team)
   }
   ```
2. **Frontend** — En `renderTeamsSection()`, añadir un botón "Crear equipo" al inicio del panel. Al hacer click, mostrar un input inline con nombre + botón confirmar. Al confirmar, llamar al binding `CreateTeam`, recargar equipos y re-renderizar.
3. Validar que no se pueda crear un equipo con nombre vacío o duplicado.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `app/bindings_test.go` | Crear equipo vacío con nombre válido; rechazar nombre vacío |
| Manual | Verificar flujo UI: click "Crear equipo" → input nombre → confirmar → equipo aparece en lista vacío |

## Criterios de aceptacion

- [x] Botón "Crear equipo" visible en el panel de Mis Equipos
- [x] Al crear equipo, aparece inmediatamente en la lista (vacío, 0/6 miembros)
- [x] No se pueden crear equipos con nombre vacío
- [x] No se pueden crear equipos con nombre duplicado (mismo slug)
- [x] El equipo creado persiste tras recargar la app
