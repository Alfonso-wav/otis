# Task 0133 — Limit moves table to 50 rows with scroll and total indicator

## Estado: done

## Goal
Limit the moves table in the Pokemon detail view (Pokedex tab) to show a maximum of 50 rows to improve performance and usability for Pokemon with many moves. When truncated, show an indicator of total moves count.

## Contexto tecnico

- The moves table is rendered by `renderTable()` inside `renderMoves()` in `frontend/src/pages/pokedex.ts` (~line 1029).
- Moves are filtered by method and sorted before rendering.
- The table wrapper `.moves-table-wrap` in `frontend/src/styles/_pokemon.scss` has scroll support via `overflow-y: auto`.
- i18n uses `t()` function with keys from `frontend/src/locales/en.json` and `es.json`.

## Acceptance criteria

- [x] The moves table shows a maximum of 50 rows after filtering and sorting
- [x] When there are more than 50 moves, an indicator shows "Showing 50 of {total} moves" (localized)
- [x] The indicator text supports i18n (Spanish and English)
- [x] The table wrapper max-height is adjusted to ~1750px to accommodate ~50 rows
- [x] Dark mode styling is applied to the indicator

## Archivos afectados

- `frontend/src/pages/pokedex.ts` — slice sorted moves to 50, render limit indicator
- `frontend/src/styles/_pokemon.scss` — adjust max-height, add `.moves-limit-info` styles
- `frontend/src/locales/en.json` — add `detail.movesLimitInfo` key
- `frontend/src/locales/es.json` — add `detail.movesLimitInfo` key
