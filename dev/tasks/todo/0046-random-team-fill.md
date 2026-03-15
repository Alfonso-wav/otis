# Rellenar equipo aleatoriamente

**ID**: 0046-random-team-fill
**Estado**: todo
**Fecha**: 2026-03-15
**Depende de**: 0045-create-team-from-panel

---

## Descripcion

Añadir un botón "Rellenar aleatoriamente" en cada equipo del panel "Mis Equipos" que complete las plazas vacías (hasta 6 miembros) con Pokémon aleatorios, cada uno con naturaleza, nivel, EVs, IVs y movimientos generados aleatoriamente de forma razonable.

## Capas afectadas

- **Core**: Añadir función pura `GenerateRandomTeamMember` que genere un TeamMember aleatorio a partir de datos de Pokémon disponibles. Añadir `FillTeamRandom` que rellene plazas vacías.
- **Shell**: Sin cambios directos (se usa el fetcher de PokeAPI existente para obtener lista de Pokémon y sus movimientos).
- **APP**: Añadir binding `FillTeamRandom(teamName string) (core.Team, error)` que obtenga datos necesarios, genere miembros random y guarde.
- **Frontend**: Añadir botón "Rellenar aleatorio" por equipo.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/teams.go` | modificar | Añadir `GenerateRandomTeamMember(pokemon core.Pokemon, allMoves []core.Move, rng func(int) int) core.TeamMember` — genera un miembro con naturaleza random, nivel 50, IVs 31, EVs distribuidos aleatoriamente, y 4 movimientos random del pool del Pokémon |
| `core/teams.go` | modificar | Añadir `FillTeamRandom(team core.Team, availablePokemon []core.Pokemon, allMoves []core.Move, rng func(int) int) core.Team` — rellena hasta 6 con miembros random sin repetir Pokémon |
| `app/bindings.go` | modificar | Añadir `FillTeamRandom(teamName string) (core.Team, error)` que cargue el equipo, obtenga lista de Pokémon, llame a core y guarde |
| `frontend/src/pages/builds.ts` | modificar | Añadir botón "Rellenar aleatorio" en cada tarjeta de equipo. Al click, llamar binding y re-renderizar |

## Plan de implementacion

1. **Core** — Implementar `GenerateRandomTeamMember`:
   - Recibir un `Pokemon` completo (con su lista de movimientos aprendibles).
   - Naturaleza: elegir aleatoria de las 25 naturalezas.
   - Nivel: 50 (estándar competitivo).
   - IVs: todos 31 (simplificación razonable).
   - EVs: distribuir 510 puntos aleatoriamente entre 6 stats (máx 252 por stat).
   - Movimientos: elegir 4 movimientos aleatorios del pool del Pokémon (solo movimientos con power > 0 si hay suficientes, sino cualquiera).

2. **Core** — Implementar `FillTeamRandom`:
   - Calcular plazas libres (`6 - len(team.Members)`).
   - Elegir N Pokémon aleatorios del pool disponible (sin repetir entre sí ni con los existentes).
   - Generar un `TeamMember` por cada uno.
   - Retornar nuevo Team con los miembros añadidos.

3. **APP** — Binding que orqueste: cargar equipo → obtener Pokémon disponibles → llamar core → guardar.

4. **Frontend** — Botón por equipo. Solo visible si el equipo tiene < 6 miembros.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/teams_test.go` | `GenerateRandomTeamMember` genera miembro válido (EVs ≤ 510, 4 movimientos, nivel correcto) |
| `core/teams_test.go` | `FillTeamRandom` rellena hasta 6, no repite Pokémon existentes |
| Manual | Verificar UI: click "Rellenar aleatorio" → equipo se completa con Pokémon variados |

## Criterios de aceptacion

- [ ] Botón "Rellenar aleatorio" visible en equipos con < 6 miembros
- [ ] Al rellenar, se completan las plazas vacías hasta 6
- [ ] No se repiten Pokémon dentro del mismo equipo
- [ ] Cada miembro generado tiene 4 movimientos, naturaleza válida, EVs ≤ 510
- [ ] El equipo relleno se persiste correctamente
- [ ] Si el equipo ya tiene 6 miembros, el botón no aparece o está deshabilitado
