# Team: persistencia de habilidad + exposición con flag isHidden

**ID**: 0169-team-ability-persistence
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

Añadir campo `ability` al esquema de team JSON. Exponer a frontend la lista de habilidades de cada Pokémon marcando las ocultas (`isHidden: true`). Permitir opción "ninguna" (`ability: ""`).

Depende de 0168 (catálogo) para que la habilidad guardada sea efectiva en simulación.

## Capas afectadas

- **Core**: `TeamMember` struct gana campo `Ability string`. `Ability` domain gana flag `IsHidden` por Pokémon (o bien `PokemonAbilityEntry { Name, IsHidden }` reemplazando `Pokemon.Abilities []string`).
- **Shell**: `pokeapi.go` ya parsea `abilities[].is_hidden` — asegurar propagación al dominio. Repositorio teams serializa/deserializa el nuevo campo.
- **APP**: bindings Wails `SaveToTeam`, `CreateTeam`, `UpdateTeamMember`, `SimulateTeamBattle` aceptan/transmiten ability. Handlers REST equivalentes en `app/mobile/handlers.go`.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | `Pokemon.Abilities []string` → `[]PokemonAbilityEntry { Name string; IsHidden bool }`. Nuevo campo `Ability string` en `TeamMember`. |
| `shell/pokeapi.go` | modificar | Mapear `ability.name` + `is_hidden` al nuevo tipo. Respetar slot order. |
| `shell/teams_repo.go` (o equivalente) | modificar | JSON ahora incluye `"ability": "drizzle"`. Retrocompatibilidad: miembros sin campo → `""`. |
| `app/bindings.go` | modificar | `SaveToTeam`, `CreateTeam`, `UpdateTeamMember` aceptan ability. Getter `ListTeams` devuelve ability. |
| `app/mobile/handlers.go` | modificar | Handlers REST equivalentes (regla CLAUDE.md: paridad Wails/REST). |
| `core/domain_test.go` / `shell/pokeapi_test.go` / repo test | modificar | Tests de parse con `is_hidden`. Test de serialización JSON. |

## Plan de implementacion

1. Definir `PokemonAbilityEntry { Name string; IsHidden bool }`. Migrar todo uso de `Pokemon.Abilities []string` al nuevo tipo. Actualizar grep call sites.
2. En shell `pokeapi.go`, al parsear `abilities[]`, guardar ambos campos.
3. Extender `TeamMember` con `Ability string` (valor canónico kebab-case o "" = ninguna).
4. Extender JSON schema de team: `"ability": "drizzle"`.
5. Backfill: archivos de team existentes sin `ability` → loader pone `""`.
6. Validación en backend al guardar: si `ability != ""`, debe existir en `Pokemon.Abilities` (o permitir libre y advertir?). **Decisión: validar estricto** — rechazar con error si no coincide.
7. Extender `SimulateTeamBattle` y `SimulateMultipleTeamBattles` para pasar la ability al `FullBattleInput` (campos `AttackerAbility`/`DefenderAbility` ya definidos en 0168).
8. Paridad Wails/REST: replicar cualquier binding que maneje team members.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `shell/pokeapi_test.go` | Parse de `is_hidden` correcto por slot. |
| `shell/teams_repo_test.go` | Round-trip JSON con ability presente/ausente. |
| `core/domain_test.go` | Retrocompat: team member sin ability carga como `""`. |

## Criterios de aceptacion

- [ ] `PokemonAbilityEntry` con `IsHidden` reemplaza `[]string`.
- [ ] Team JSON acepta `"ability"` por miembro.
- [ ] Loading retrocompatible: equipos sin `ability` cargan sin error.
- [ ] Validación al guardar: ability debe pertenecer al Pokémon o ser "".
- [ ] `SimulateTeamBattle` y `SimulateMultipleTeamBattles` propagan ability al motor.
- [ ] Paridad Wails/REST en `app/bindings.go` ↔ `app/mobile/handlers.go`.
- [ ] Tests de parsing y serialización pasan.
- [ ] `data/teams/asas.json` y `data/teams/dfdf.json` siguen cargando (sin ability = no-op).

## Notas

- No romper APK móvil — revisar `handlers.go` siempre.
- Depende de 0168: sin catálogo no hay efecto, pero el campo puede existir sin él.
- Frontend (UI selector) en 0170 — no tocar aquí.
