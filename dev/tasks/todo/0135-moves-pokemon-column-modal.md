# Task 0135 — Moves table: add "PKMN" column with Pokemon learners modal

## Estado: done

## Goal
Add a new "PKMN" column to the Explore > Moves table. Each cell shows a red button labeled "PKMN" that opens a modal listing all Pokemon that can learn that move. Clicking a Pokemon in the modal navigates to its individual Pokedex detail view.

## Contexto tecnico

### Backend (Go)
- `core.Move` struct (`core/domain.go:14-25`) does NOT have a field for Pokemon learners.
- PokeAPI's `/move/{name}` endpoint returns `learned_by_pokemon` (array of `{name, url}`), but `apiMove` struct (`shell/pokeapi.go:369-385`) does not parse it.
- `FetchMove()` (`shell/pokeapi.go:387-447`) maps the raw response to `core.Move` — needs to include the new field.
- `FetchAllMoves()` (`shell/pokeapi_moves.go:39-100`) fetches all ~920 moves using `FetchMove()` per move — the new field will propagate automatically.
- **Warning**: Adding `LearnedBy` to all ~920 moves increases memory and payload size significantly. Consider returning only count in `GetAllMoves` and lazy-loading the full list via `GetMove(name)` when the modal opens.

### Frontend
- Moves table: `frontend/src/pages/explore/moves.ts` — 7 columns currently (name, type, category, power, accuracy, pp, priority).
- Column config: `movesTableColumns()` (line 43-53), table HTML (lines 170-186), row render (lines 116-127).
- Existing modal pattern to reuse: `frontend/src/components/ability-pokemon-modal.ts` — simple modal with Pokemon grid, sprite images, click-to-navigate.
- Navigation to Pokemon detail from modals uses pattern: close modal → click Pokédex tab → set search input → click search button.
- CSS: modal styles already exist in `_explore.scss` (`.type-modal-overlay`, `.type-modal`, `.type-modal-grid`, etc.).
- i18n: keys in `frontend/src/locales/en.json` and `es.json` under `moves.columns.*`.
- The button style (red rectangular background) needs a new CSS class (e.g., `.pkmn-btn`).

### UI language
- Column labels are in English abbreviated form (e.g., "Cat.", "Acc.", "Prio.") — the button label "PKMN" is consistent with this style.
- Modal title can follow existing pattern: "Move Name (N)" showing count.

## Acceptance criteria

- [ ] `core.Move` struct has a new `LearnedBy []string` field with Pokemon names
- [ ] `apiMove` struct parses `learned_by_pokemon` from PokeAPI
- [ ] `FetchMove()` populates `LearnedBy` in the returned `core.Move`
- [ ] Moves table has a new "PKMN" column (after priority or as last column)
- [ ] Each cell in the PKMN column shows a red rectangular button labeled "PKMN"
- [ ] Clicking the button opens a modal with a grid of Pokemon that can learn the move
- [ ] Each Pokemon in the modal shows sprite + name
- [ ] Clicking a Pokemon in the modal navigates to its Pokedex detail view
- [ ] Modal closes on Escape, click outside, or close button
- [ ] Column toggle supports the new PKMN column
- [ ] i18n keys added for EN and ES (column label, modal empty state)
- [ ] Dark mode styling works correctly for button and modal

## Archivos afectados

- `core/domain.go` — add `LearnedBy []string` to `Move` struct
- `shell/pokeapi.go` — add `LearnedByPokemon` to `apiMove`, map in `FetchMove()`
- `frontend/src/pages/explore/moves.ts` — add PKMN column config, header, row cell, button click handler, modal integration
- `frontend/src/components/move-pokemon-modal.ts` — new file, modal for move learners (based on ability-pokemon-modal pattern)
- `frontend/src/locales/en.json` — add `moves.columns.pokemon` and `modals.moveEmpty` keys
- `frontend/src/locales/es.json` — add corresponding Spanish translations
- `frontend/src/styles/_explore.scss` — add `.pkmn-btn` styles (red background, white text)
- `frontend/src/styles/_dark.scss` — dark mode adjustments if needed
