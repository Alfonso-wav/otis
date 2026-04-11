# Task 0134 — Type chart: replace click-to-fade filter with red X button per type

## Estado: todo

## Goal
Change the filtering mechanism in the Explore > Type Chart table. Instead of clicking/dragging on type headers to fade them, add a small red "X" button next to each type icon. Clicking the X filters out that type. Remove drag-to-filter behavior. Keep the existing "Restore types" button to clear all filters.

## Contexto tecnico

- The type chart is rendered in `frontend/src/pages/explore/type-chart.ts`.
- Current filtering uses `mousedown`/`mouseenter`/`touchstart`/`touchmove` on header cells (`th.tc-col-header`, `th.tc-row-header`) to add types to a `filteredTypes` Set, then re-renders.
- Drag filtering tracks `isDragging` and `dragAxis` state variables.
- The restore button (`.tc-restore-btn`) clears `filteredTypes` and re-renders — must be preserved as-is.
- Styling lives in `frontend/src/styles/_explore.scss` (lines 474–667) with responsive breakpoints at 991px, 767px, and 479px.
- Type icons use `<img src="/assets/types/{type}.svg" class="tc-icon">` inside a colored circular background.
- No i18n changes needed (the X is a visual symbol, not text).

## Acceptance criteria

- [ ] Each type header (column and row) shows a small red "X" button positioned near the type icon
- [ ] Clicking the X button on a type filters it out (adds to `filteredTypes` and re-renders)
- [ ] Click/drag on the type header itself no longer triggers filtering (remove `mousedown`, `mouseenter`, `touchstart`, `touchmove` filter handlers)
- [ ] Remove `isDragging` and `dragAxis` state variables and related `mouseup`/`touchend` document listeners
- [ ] The "Restore types" button continues to work as before
- [ ] The X button is styled in red, small, and doesn't overlap or obstruct the type icon
- [ ] Responsive: the X button scales appropriately at all breakpoints (991px, 767px, 479px)
- [ ] Headers no longer show `cursor: pointer` or hover opacity change (since clicking them does nothing now)

## Archivos afectados

- `frontend/src/pages/explore/type-chart.ts` — replace drag/click handlers with X button click handler, update `typeHeader()` to include X button, simplify `attachFilterListeners()`, remove drag state
- `frontend/src/styles/_explore.scss` — add `.tc-remove-btn` styles, remove interactive cursor/hover from headers, responsive sizing for X button
