# UI: selector de habilidad en team member y atacante/defensor 1v1

**ID**: 0170-ui-ability-selector
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

Añadir selector de habilidad en la interfaz:
1. Modal de edición de team member (creación/edición de equipo).
2. Selector de atacante 1v1 en `builds.ts`.
3. Selector de defensor 1v1 en `builds.ts`.

Muestra habilidades normales + ocultas (con badge "HA"). Permite opción "Ninguna" (string vacío). Al simular, pasa la ability al backend.

Depende de 0168 (catálogo) y 0169 (persistencia + flag `isHidden`).

## Capas afectadas

- **Core**: ninguno.
- **Shell**: ninguno (el binding ya acepta `ability` tras 0169).
- **APP/Frontend**: `frontend/src/pages/builds.ts` + subcomponentes de team member + API client.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | Añadir estado `attackerAbility`, `defenderAbility`. Render `<select>` con placeholder "-- Ninguna --" + opciones normales + ocultas marcadas `[HA]`. Pasar a `SimulateBattle` y afines. |
| `frontend/src/pages/teams/*.ts` | modificar | Modal de team member gana un `<select>` de ability con las mismas reglas. Persiste vía `SaveToTeam`/`UpdateTeamMember`. |
| `frontend/src/api.ts` | modificar | Extender tipos `TeamMember`, `SimulateBattleInput` con `ability` (string). Extender `Pokemon.abilities` al nuevo shape `{ name, isHidden }`. |
| `frontend/src/components/ability-badge.ts` | crear (opcional) | Pequeño componente visual para badge "HA". |
| `frontend/src/locales/en.json` | modificar | Keys: `builds.ability`, `builds.abilityNone`, `builds.abilityHidden`, logs del motor emitidos por 0168. |
| `frontend/src/locales/es.json` | modificar | Idem ES. |
| `frontend/src/styles/_pokemon.scss` (o similar) | modificar | Estilos del badge HA y layout del selector. |

## Plan de implementacion

1. Extender tipo `Pokemon` en `api.ts`: `abilities: { name: string; isHidden: boolean }[]`.
2. Extender tipo `TeamMember`: `ability: string`.
3. En modal de team member: al cargar opciones del Pokémon seleccionado, render `<select>`:
   ```html
   <select>
     <option value="" disabled selected>-- Habilidad --</option>
     <option value="">Ninguna</option>
     <option value="blaze">Mar Llamas</option>
     <option value="solar-power">Fuerza Solar [HA]</option>
   </select>
   ```
   (placeholder disabled según convención CLAUDE.md).
4. Mostrar nombre localizado (fetch desde `abilities.ts` cache si ya existe).
5. En `builds.ts` 1v1: tras seleccionar atacante/defensor, cargar sus abilities y render selector análogo. Estado module-scoped (para no perder en re-render).
6. Al simular: pasar `attackerAbility`, `defenderAbility` al binding.
7. i18n: escuchar `locale-changed` y re-render.
8. Probar build de producción.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual QA | Selector aparece, carga opciones, persiste al guardar equipo, se refleja en simulación (log muestra efecto de habilidad). |
| Manual QA 360px | Selector accesible en móvil, no rompe layout. |
| Manual APK | Paridad con desktop. |

## Criterios de aceptacion

- [ ] Selector de ability visible en modal team member.
- [ ] Selector de ability visible en atacante 1v1.
- [ ] Selector de ability visible en defensor 1v1.
- [ ] Habilidades ocultas marcadas con `[HA]` o badge.
- [ ] Opción "Ninguna" funcional.
- [ ] Valor persiste al guardar equipo y se recupera al editar.
- [ ] Simulación 1v1 + batalla equipo refleja efecto (visible en log).
- [ ] i18n EN/ES completo, `locale-changed` respetado.
- [ ] Placeholder disabled en `<select>` (convención CLAUDE.md).
- [ ] Touch target ≥44x44px en móvil.
- [ ] Probado en build de producción.
- [ ] Probado en viewport 360px.
- [ ] Paridad Wails/REST ya resuelta por 0169.

## Notas

- Nombres localizados provienen de `Ability.NameEs` / `Ability.Name` (ya en core).
- Si el Pokémon no tiene habilidades (teóricamente imposible tras gen 3), selector deshabilitado con "Ninguna".
- No añadir auto-guardado en el selector — guardar solo al pulsar "Guardar" del modal.
- Al cambiar el Pokémon del slot, resetear ability a "".
