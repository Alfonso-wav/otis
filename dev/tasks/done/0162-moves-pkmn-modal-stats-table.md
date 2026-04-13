# Task 0162 — Moves PKMN modal: Pokemon stats table instead of icons grid

## Estado: done

## Goal
En Explore > Movimientos, la columna **PKMN** abre un modal con grid de sprites + nombre. Cambiarlo por una **tabla de Pokemon con stats** (mismo patron visual que la tabla de Pokedex), manteniendo navegacion al detalle al click.

## Contexto tecnico

### Frontend
- Modal actual: `frontend/src/components/move-pokemon-modal.ts` (111 lineas). Recibe `moveName`, hace `GetMove(name)`, usa `move.LearnedBy` (array de nombres).
- Actualmente renderiza `.type-modal-grid` con sprites + nombres (lineas 76-89).
- Navegacion al detalle: `navigateToPokemon(name)` (lineas 22-31) — cierra modal, click en tab Pokedex, rellena search, click buscar.
- Patron de tabla a reutilizar: `frontend/src/pages/pokedex.ts` (lineas 394-479). Columnas: ID, sprite, nombre, types, HP, Atk, Def, SpA, SpD, Spe, Total. Stats via `GetPokemon()` batch (lineas 402-411). Orden stats: `[HP, Atk, Def, SpA, SpD, Spe]`.
- i18n existentes reutilizables: `pokedex.columns.{id,sprite,name,types,hp,atk,def,spa,spd,vel,total}`, `modals.moveLoading`, `modals.moveEmpty`.
- `GetMove` devuelve solo nombres — para stats hay que hacer fetch en paralelo de cada Pokemon (`GetPokemon(name)` para cada entrada en `LearnedBy`).

### Consideraciones
- `LearnedBy` puede tener 100+ entradas (ej. Tackle, Rest). Fetch secuencial seria lento. Opciones:
  - A) Paralelo con `Promise.all` (riesgo rate limit / memoria).
  - B) Limitar a primeras N (ej. 50) con nota "mostrando X de Y". Consistente con limite de 50 rows en otras tablas del proyecto (ver 0128, 0133).
  - **Preferido**: paralelo con spinner Diglett; si lista > 100, paginar o limitar a 50 iniciales con boton "ver mas".
- Modal debe crecer a ancho mayor para tabla (actualmente `.type-modal` tiene max-width chico). Nueva clase `.move-modal--wide` o similar.
- Tabla scrollable vertical (`max-height` + `overflow-y:auto`) — evitar modal mas alto que viewport.
- Ordenacion de columnas: opcional. MVP sin sort; si trivial de agregar reusando logica de pokedex, incluir.

### UI language
- Spanish / i18n activo. Column labels via `t("pokedex.columns.*")` ya existentes.

## Acceptance criteria

- [ ] Click en boton PKMN de tabla moves abre modal con **tabla de stats**, no grid de iconos.
- [ ] Tabla tiene columnas: sprite, nombre, types, HP, Atk, Def, SpA, SpD, Spe, Total.
- [ ] Click en fila navega al detalle del Pokemon (cierra modal → tab Pokedex → search → click).
- [ ] Loading spinner (Diglett) mientras se fetch stats de cada Pokemon.
- [ ] Si `LearnedBy` > 50, limitar a 50 iniciales o paginar.
- [ ] Modal scroll vertical si tabla excede viewport.
- [ ] Dark mode coherente con resto de tablas.
- [ ] Mobile (360px): tabla con `overflow-x:auto`, no rompe layout.
- [ ] i18n: reutiliza keys existentes. Si hay texto nuevo, agregar a `en.json` Y `es.json`.

## Archivos afectados

- `frontend/src/components/move-pokemon-modal.ts` — reemplazar grid por render de tabla, fetch paralelo de stats.
- `frontend/src/styles/_explore.scss` o `_modals.scss` — ancho mayor para modal, estilos tabla reusables.
- `frontend/src/locales/{en,es}.json` — keys nuevas si aplica.

## Notas

- No tocar backend. `LearnedBy` ya devuelve nombres; stats via `GetPokemon()` ya existente.
- Si fetch paralelo de 200+ Pokemon causa lag, limitar a 50 y mostrar contador total.
