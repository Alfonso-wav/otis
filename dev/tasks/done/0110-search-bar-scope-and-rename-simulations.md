# Task 0110 — Search bar scoped to Pokédex + rename Builds → Simulations

## Goal
Two small UI changes:
1. The search bar (`#search-bar`) in the header must only be visible when the **Pokédex** tab is active. Switching to Explorar or Simulations hides it.
2. Rename the **Builds** tab to **Simulations** everywhere (HTML, i18n es + en).

## Context
- Search bar is a `<div id="search-bar">` inside `<header class="pokedex-header">` in `frontend/index.html` (line 25-28). It is always visible regardless of active tab.
- Tab navigation is handled by `frontend/src/router.ts` — `navigate(id)` is called whenever a tab changes.
- i18n keys live in `frontend/src/locales/es.json` and `frontend/src/locales/en.json`. Key: `tabs.builds` (currently `"Builds"` in both).
- The HTML button has `data-tab="builds"` and `data-i18n="tabs.builds"`.

## Scope
- **Frontend only** — no backend changes.
- Layers: APP (index.html) + Shell-equivalent (router.ts, i18n locales).

## Steps

### 1. Scope search bar to Pokédex tab
In `frontend/src/router.ts`, after updating active tab styles in `navigate()`, toggle the visibility of `#search-bar` based on whether the destination tab is `"pokedex"`.

```ts
const searchBar = document.getElementById("search-bar");
if (searchBar) {
  searchBar.classList.toggle("hidden", id !== "pokedex");
}
```

Also ensure it's visible on initial load (already active tab is `"pokedex"` via `initRouter("pokedex")`).

### 2. Rename Builds → Simulations in HTML
In `frontend/index.html` line 33, the button text is set via `data-i18n="tabs.builds"`. The button `id` and `data-tab` attribute stay as `"builds"` (internal routing key), only the displayed label changes.

### 3. Update i18n keys
- `frontend/src/locales/es.json`: change `tabs.builds` from `"Builds"` to `"Simulaciones"`
- `frontend/src/locales/en.json`: change `tabs.builds` from `"Builds"` to `"Simulations"`
- Also check `builds.title` key (currently `"Builds & Simulador"` in es.json) — update to `"Simulaciones"` or similar.

## Files to modify
- `frontend/index.html`
- `frontend/src/router.ts`
- `frontend/src/locales/es.json`
- `frontend/src/locales/en.json`

## Acceptance criteria
- [ ] Navegando a Explorar o Simulations, el buscador desaparece del header
- [ ] Al volver a Pokédex, el buscador reaparece
- [ ] La pestaña muestra "Simulaciones" (es) / "Simulations" (en)
- [ ] El routing interno sigue funcionando con `data-tab="builds"`
