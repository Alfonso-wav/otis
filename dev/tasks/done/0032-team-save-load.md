# Guardar y cargar equipos de Pokemon

**ID**: 0032-team-save-load
**Estado**: done
**Fecha**: 2026-03-14

---

## Descripcion

Permitir al usuario guardar Pokemon buildeados (con sus moves, IVs, EVs, naturaleza y nivel) en equipos con nombre personalizado, y poder importar esos equipos al configurar simulaciones de batalla.

### Funcionalidad concreta

1. **Guardar equipo**: En la seccion Builds, un boton para guardar el Pokemon actual (attacker o defender) en un equipo. El usuario elige el nombre del equipo (nuevo o existente) y el Pokemon se agrega con toda su configuracion (moves, stats, nivel, naturaleza, IVs, EVs).
2. **Gestionar equipos**: Ver equipos guardados, ver los Pokemon dentro de cada equipo, eliminar Pokemon individuales o equipos completos.
3. **Importar desde equipo**: Al configurar una simulacion (attacker/defender), opcion de importar un Pokemon desde un equipo guardado, cargando automaticamente todos sus datos (moves, IVs, EVs, naturaleza, nivel).

## Capas afectadas
- **Core**: Definir tipos de dominio `Team`, `TeamMember`. Funciones puras de validacion.
- **Shell**: Persistencia en archivos JSON locales (directorio `data/teams/`).
- **APP**: Bindings Wails para CRUD de equipos y miembros.
- **Frontend**: UI de gestion de equipos en Builds + selector de importacion.

## Dependencias externas nuevas
Ninguna. Se usa `encoding/json` y `os` de la stdlib de Go.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Agregar tipos `Team`, `TeamMember` |
| `core/teams.go` | crear | Funciones puras: validacion de equipo (nombre no vacio, max 6 miembros, etc.) |
| `core/ports.go` | modificar | Agregar interfaz `TeamStorage` con metodos CRUD |
| `shell/teams.go` | crear | Implementacion de `TeamStorage` con persistencia JSON en disco |
| `app/bindings.go` | modificar | Agregar bindings: `SaveTeamMember`, `ListTeams`, `GetTeam`, `DeleteTeam`, `DeleteTeamMember` |
| `app/main.go` (root) | modificar | Inyectar `TeamStorage` en el wiring |
| `frontend/src/pages/builds.ts` | modificar | UI para guardar en equipo, listar equipos, importar desde equipo |

## Plan de implementacion

### Fase 1 — Tipos de dominio y validacion (Core)
1. Agregar en `core/domain.go`:
   ```go
   type TeamMember struct {
       PokemonName string       `json:"pokemonName"`
       Moves       []string     `json:"moves"`       // nombres de los moves equipados
       Level       int          `json:"level"`
       Nature      string       `json:"nature"`
       IVs         Stats        `json:"ivs"`
       EVs         Stats        `json:"evs"`
   }

   type Team struct {
       Name    string       `json:"name"`
       Members []TeamMember `json:"members"`
   }
   ```
2. Crear `core/teams.go` con funciones puras:
   - `ValidateTeam(team Team) error` — nombre no vacio, max 6 miembros
   - `ValidateTeamMember(member TeamMember) error` — EVs <= 510 total, IVs 0-31, nivel 1-100
   - `AddMemberToTeam(team Team, member TeamMember) (Team, error)` — retorna nuevo Team con el miembro agregado
   - `RemoveMemberFromTeam(team Team, index int) (Team, error)` — retorna nuevo Team sin el miembro

### Fase 2 — Persistencia en disco (Shell)
3. Agregar interfaz `TeamStorage` en `core/ports.go`:
   ```go
   type TeamStorage interface {
       SaveTeam(team Team) error
       ListTeams() ([]Team, error)
       GetTeam(name string) (Team, error)
       DeleteTeam(name string) error
   }
   ```
4. Crear `shell/teams.go`:
   - Struct `FileTeamStorage` con directorio base (`data/teams/`)
   - Cada equipo se guarda como `{slug-del-nombre}.json`
   - Implementar los 4 metodos de la interfaz
   - Crear el directorio automaticamente si no existe

### Fase 3 — Bindings Wails (APP)
5. Modificar `app/bindings.go` para agregar:
   - `SaveToTeam(teamName string, member core.TeamMember) error` — carga equipo existente o crea nuevo, agrega miembro, guarda
   - `ListTeams() ([]core.Team, error)`
   - `GetTeam(name string) (core.Team, error)`
   - `DeleteTeam(name string) error`
   - `DeleteTeamMember(teamName string, memberIndex int) error`
6. Modificar `main.go` para inyectar `FileTeamStorage`

### Fase 4 — Frontend: Guardar en equipo
7. En `builds.ts`, agregar boton "Guardar en equipo" junto a cada Pokemon configurado (attacker/defender)
8. Al pulsar, mostrar modal/dropdown con:
   - Input de texto para nombre del equipo (con autocomplete de equipos existentes)
   - Boton confirmar
9. Al confirmar, llamar a `SaveToTeam()` con el nombre y los datos del Pokemon actual

### Fase 5 — Frontend: Importar desde equipo
10. En `builds.ts`, agregar boton "Importar de equipo" junto al buscador de Pokemon (attacker/defender)
11. Al pulsar, mostrar panel con:
    - Lista de equipos guardados
    - Al seleccionar un equipo, mostrar sus miembros con nombre y moves
    - Al seleccionar un miembro, cargar todos sus datos en el slot correspondiente (attacker/defender)
12. Al importar: cargar Pokemon via `GetPokemon()`, setear nivel/naturaleza/IVs/EVs, cargar moves via `GetMove()`

### Fase 6 — Frontend: Gestion de equipos
13. Agregar seccion "Mis equipos" en la pagina de Builds (colapsable)
14. Mostrar lista de equipos con sus miembros
15. Botones para eliminar equipo completo o miembro individual

## Tests
| Archivo | Que se testea |
|---------|---------------|
| `core/teams_test.go` | Validacion de equipos y miembros, agregar/eliminar miembros (funciones puras) |
| `shell/teams_test.go` | CRUD de archivos JSON: guardar, leer, listar, eliminar equipos en directorio temporal |

## Criterios de aceptacion
- [ ] Se pueden guardar Pokemon buildeados en equipos con nombre personalizado
- [ ] Los equipos persisten entre sesiones (archivos JSON en disco)
- [ ] Se pueden ver los equipos guardados y sus miembros
- [ ] Se puede importar un Pokemon desde un equipo guardado al configurar attacker/defender
- [ ] Al importar se cargan todos los datos: Pokemon, moves, nivel, naturaleza, IVs, EVs
- [ ] Se pueden eliminar equipos y miembros individuales
- [ ] Max 6 miembros por equipo
- [ ] Validacion de datos al guardar (EVs, IVs, nivel en rango)
- [ ] Tests unitarios de Core y Shell pasan

## Notas
- Los equipos se guardan como archivos JSON individuales para simplicidad. No se necesita base de datos.
- El slug del nombre se genera normalizando: lowercase, espacios a guiones, sin caracteres especiales.
- Si un equipo ya tiene 6 miembros, mostrar error al intentar agregar otro.
- El mismo Pokemon puede aparecer varias veces en un equipo (con diferentes builds).
