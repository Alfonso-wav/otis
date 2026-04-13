import {
  ListPokemon,
  GetPokemon,
  ListGenerations,
  GetGeneration,
  ListTypes,
  GetType,
  GetPokemonSpecies,
  GetPokemonEncounters,
  GetAllSpeciesClassifications,
  GetAbility,
} from "../api";
import type { Pokemon, PokemonListItem, PokemonMoveEntry } from "../types";
import { t, typeName, statName, getLocale } from "../i18n";
import { showView, staggerCards, morphToTable, morphToGrid } from "../animations/transitions";
import { initColumnToggle, reapplyColumnVisibility, type ColumnConfig } from "../components/column-toggle";
import { SortCache } from "../utils/sort-cache";
import { showSortingOverlay, updateSortingOverlayText, hideSortingOverlay, createInlineDiglett } from "../components/sorting-overlay";
import { createAutocomplete } from "../autocomplete";
import { loadMoveNames, getLocalizedMoveName } from "../utils/move-names";

const LIMIT_STEP = 50;
let rowLimit = 50;

interface FilterState {
  generations: string[];
  types: string[];
  legendary: boolean;
  mythical: boolean;
}

let filter: FilterState = { generations: [], types: [], legendary: false, mythical: false };
let filteredList: PokemonListItem[] = [];
let isFiltering = false;
let offset = 0;
let totalCount = 0;
let viewMode: "grid" | "table" = "grid";
let infiniteLoading = false;
let scrollObserver: IntersectionObserver | null = null;
let scrollSentinel: HTMLDivElement;

// -- Sorting state -----------------------------------------------------------

type SortDirection = 'asc' | 'desc' | null;
type SortColumn = 'id' | 'name' | 'hp' | 'atk' | 'def' | 'spa' | 'spd' | 'vel' | 'total' | null;

let currentSortColumn: SortColumn = null;
let currentSortDirection: SortDirection = null;
let sortedFullList: Pokemon[] | null = null;
let sortingLoading = false;

const pokemonSortCache = new SortCache<Pokemon>([
  { key: "id", compare: (a, b) => a.ID - b.ID },
  { key: "name", compare: (a, b) => a.Name.localeCompare(b.Name) },
  { key: "hp", compare: (a, b) => a.Stats[0].BaseStat - b.Stats[0].BaseStat },
  { key: "atk", compare: (a, b) => a.Stats[1].BaseStat - b.Stats[1].BaseStat },
  { key: "def", compare: (a, b) => a.Stats[2].BaseStat - b.Stats[2].BaseStat },
  { key: "spa", compare: (a, b) => a.Stats[3].BaseStat - b.Stats[3].BaseStat },
  { key: "spd", compare: (a, b) => a.Stats[4].BaseStat - b.Stats[4].BaseStat },
  { key: "vel", compare: (a, b) => a.Stats[5].BaseStat - b.Stats[5].BaseStat },
  { key: "total", compare: (a, b) => {
    const totalA = a.Stats.reduce((s, st) => s + st.BaseStat, 0);
    const totalB = b.Stats.reduce((s, st) => s + st.BaseStat, 0);
    return totalA - totalB;
  }},
]);

function pokedexTableColumns(): ColumnConfig[] {
  return [
    { key: "id", label: t("pokedex.columns.id") },
    { key: "sprite", label: t("pokedex.columns.sprite") },
    { key: "name", label: t("pokedex.columns.name"), fixed: true },
    { key: "types", label: t("pokedex.columns.types") },
    { key: "hp", label: t("pokedex.columns.hp") },
    { key: "atk", label: t("pokedex.columns.atk") },
    { key: "def", label: t("pokedex.columns.def") },
    { key: "spa", label: t("pokedex.columns.spa") },
    { key: "spd", label: t("pokedex.columns.spd") },
    { key: "vel", label: t("pokedex.columns.vel") },
    { key: "total", label: t("pokedex.columns.total") },
  ];
}

function encounterTableColumns(): ColumnConfig[] {
  return [
    { key: "location", label: t("encounters.columns.location"), fixed: true },
    { key: "game", label: t("encounters.columns.game") },
    { key: "method", label: t("encounters.columns.method") },
    { key: "chance", label: t("encounters.columns.chance") },
    { key: "levels", label: t("encounters.columns.levels") },
    { key: "conditions", label: t("encounters.columns.conditions") },
  ];
}

function sortPokemonData(data: Pokemon[], column: SortColumn, direction: SortDirection): Pokemon[] {
  if (!column || !direction) return data;
  const sorted = [...data];
  const mult = direction === 'asc' ? 1 : -1;
  sorted.sort((a, b) => {
    switch (column) {
      case 'id': return mult * (a.ID - b.ID);
      case 'name': return mult * a.Name.localeCompare(b.Name);
      case 'hp': return mult * (a.Stats[0].BaseStat - b.Stats[0].BaseStat);
      case 'atk': return mult * (a.Stats[1].BaseStat - b.Stats[1].BaseStat);
      case 'def': return mult * (a.Stats[2].BaseStat - b.Stats[2].BaseStat);
      case 'spa': return mult * (a.Stats[3].BaseStat - b.Stats[3].BaseStat);
      case 'spd': return mult * (a.Stats[4].BaseStat - b.Stats[4].BaseStat);
      case 'vel': return mult * (a.Stats[5].BaseStat - b.Stats[5].BaseStat);
      case 'total': {
        const totalA = a.Stats.reduce((s, st) => s + st.BaseStat, 0);
        const totalB = b.Stats.reduce((s, st) => s + st.BaseStat, 0);
        return mult * (totalA - totalB);
      }
      default: return 0;
    }
  });
  return sorted;
}

async function ensureAllPokemonLoaded(): Promise<Pokemon[]> {
  let itemsToLoad: PokemonListItem[] = [];

  if (hasFilter()) {
    itemsToLoad = filteredList;
  } else {
    // Get full list from backend
    const fullList = await ListPokemon(0, totalCount);
    itemsToLoad = fullList.Results;
  }

  // Find which ones are not cached
  const uncached = itemsToLoad.filter((item) => !pokemonDataCache.has(item.Name));

  // Load in batches of 50
  const BATCH = 50;
  for (let i = 0; i < uncached.length; i += BATCH) {
    const batch = uncached.slice(i, i + BATCH);
    await Promise.all(
      batch.map(async (item) => {
        const p = await GetPokemon(item.Name);
        pokemonDataCache.set(item.Name, p);
      }),
    );
    // Update loading progress
    const loaded = Math.min(i + BATCH, uncached.length);
    const total = uncached.length;
    const pct = Math.round((loaded / total) * 100);
    updateSortingOverlayText(t("pokedex.sortingPct", { pct }));
  }

  return itemsToLoad.map((item) => pokemonDataCache.get(item.Name)!);
}

// Cache para clasificaciones species (legendary/mythical) — se carga una sola vez desde el backend
let classificationsCache: Record<string, { isLegendary: boolean; isMythical: boolean }> | null = null;

// Cache para datos completos de Pokémon (tabla)
const pokemonDataCache = new Map<string, Pokemon>();

// Nombres de todos los Pokémon (para el autocomplete del comparador)
let allPokemonNames: string[] = [];

// DOM refs
let grid: HTMLDivElement;
let listView: HTMLElement;
let detailView: HTMLElement;
let detailEl: HTMLDivElement;
let limitDecreaseBtn: HTMLButtonElement;
let limitIncreaseBtn: HTMLButtonElement;
let limitDisplay: HTMLSpanElement;
let rowLimitControl: HTMLDivElement;
let pagePrevBtn: HTMLButtonElement;
let pageNextBtn: HTMLButtonElement;
let pageIndicator: HTMLSpanElement;
let searchInput: HTMLInputElement;
let searchBtn: HTMLButtonElement;
let backBtn: HTMLButtonElement;
let filterGenContainer: HTMLDivElement;
let filterTypeContainer: HTMLDivElement;
let filterLegendaryBtn: HTMLButtonElement;
let filterMythicalBtn: HTMLButtonElement;
let filterResetBtn: HTMLButtonElement;
let viewToggleBtn: HTMLButtonElement;

function hasFilter(): boolean {
  return (
    filter.generations.length > 0 ||
    filter.types.length > 0 ||
    filter.legendary ||
    filter.mythical
  );
}

// -- Lista sin filtro --------------------------------------------------------

async function loadList(): Promise<void> {
  createInlineDiglett(grid, t("common.loading"));
  try {
    const data = await ListPokemon(offset, rowLimit);
    totalCount = data.Count;
    if (viewMode === "grid") {
      renderGrid(data.Results, false);
      updateScrollSentinel();
    } else {
      await renderTable(data.Results);
      updateRowLimitControl();
    }
  } catch (err: unknown) {
    grid.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}

// -- Lista con filtro --------------------------------------------------------

async function loadFiltered(): Promise<void> {
  resetSorting();
  createInlineDiglett(grid, t("pokedex.applyingFilters"));
  try {
    let base: PokemonListItem[] = [];
    const hasGens = filter.generations.length > 0;
    const hasTypes = filter.types.length > 0;

    if (hasGens && hasTypes) {
      // Both: union of gens INTERSECT union of types
      const [genResults, typeResults] = await Promise.all([
        Promise.all(filter.generations.map((g) => GetGeneration(g))),
        Promise.all(filter.types.map((t) => GetType(t))),
      ]);
      const genNames = new Set<string>();
      for (const gen of genResults) {
        for (const p of gen.PokemonSpecies) genNames.add(p.Name);
      }
      const typeMap = new Map<string, PokemonListItem>();
      for (const td of typeResults) {
        for (const p of td.Pokemon) {
          if (!typeMap.has(p.Name)) typeMap.set(p.Name, { Name: p.Name, URL: p.URL });
        }
      }
      base = [...typeMap.values()].filter((p) => genNames.has(p.Name));
    } else if (hasGens) {
      const genResults = await Promise.all(filter.generations.map((g) => GetGeneration(g)));
      const seen = new Set<string>();
      for (const gen of genResults) {
        for (const p of gen.PokemonSpecies) {
          if (!seen.has(p.Name)) {
            seen.add(p.Name);
            base.push({ Name: p.Name, URL: p.URL });
          }
        }
      }
    } else if (hasTypes) {
      const typeResults = await Promise.all(filter.types.map((t) => GetType(t)));
      const seen = new Set<string>();
      for (const td of typeResults) {
        for (const p of td.Pokemon) {
          if (!seen.has(p.Name)) {
            seen.add(p.Name);
            base.push({ Name: p.Name, URL: p.URL });
          }
        }
      }
    } else if (filter.legendary || filter.mythical) {
      // No gen/type selected: load all pokemon progressively
      base = await loadAllPokemonList();
    }

    if ((filter.legendary || filter.mythical) && base.length > 0) {
      base = await filterByLegendary(base);
    }

    filteredList = base;
    totalCount = base.length;
    offset = 0;

    const page = filteredList.slice(0, rowLimit);
    if (viewMode === "grid") {
      renderGrid(page, false);
      updateScrollSentinel();
    } else {
      await renderTable(page);
      updateRowLimitControl();
    }
  } catch (err: unknown) {
    grid.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}

async function loadAllPokemonList(): Promise<PokemonListItem[]> {
  const BATCH_SIZE = 100;
  const all: PokemonListItem[] = [];
  let currentOffset = 0;
  // First call to get count
  const first = await ListPokemon(0, BATCH_SIZE);
  const total = first.Count;
  all.push(...first.Results);
  currentOffset += BATCH_SIZE;

  while (currentOffset < total) {
    const pct = Math.round((currentOffset / total) * 100);
    const existingText = grid.querySelector<HTMLParagraphElement>(".diglett-inline__text");
    if (existingText) {
      existingText.textContent = t("pokedex.loadingFullList", { pct, current: currentOffset, total });
    } else {
      createInlineDiglett(grid, t("pokedex.loadingFullList", { pct, current: currentOffset, total }));
    }
    const batch = await ListPokemon(currentOffset, BATCH_SIZE);
    all.push(...batch.Results);
    currentOffset += BATCH_SIZE;
  }

  return all;
}

async function filterByLegendary(list: PokemonListItem[]): Promise<PokemonListItem[]> {
  if (!classificationsCache) {
    createInlineDiglett(grid, t("pokedex.loadingClassifications"));
    classificationsCache = await GetAllSpeciesClassifications();
  }

  const results: PokemonListItem[] = [];
  for (const item of list) {
    const data = classificationsCache[item.Name];
    if (!data) continue;
    if (filter.legendary && filter.mythical) {
      if (data.isLegendary || data.isMythical) results.push(item);
    } else if (filter.legendary) {
      if (data.isLegendary) results.push(item);
    } else if (filter.mythical) {
      if (data.isMythical) results.push(item);
    }
  }

  return results;
}

// -- Render ------------------------------------------------------------------

function renderGrid(items: PokemonListItem[], append = false): void {
  lastRenderedItems = items;
  if (!items || items.length === 0) {
    if (!append) grid.innerHTML = `<p class="loading">${t("pokedex.noResults")}</p>`;
    return;
  }

  const html = items
    .map((item) => {
      const id = idFromURL(item.URL);
      const sprite = spriteURL(item.Name);
      const numId = parseInt(id, 10);
      const paddedId = isNaN(numId) ? "" : `#${String(numId).padStart(3, "0")}`;
      return `<div class="poke-card" data-name="${item.Name}">
        <span class="poke-card__number">${paddedId}</span>
        <img class="poke-card__sprite" src="${sprite}" data-fallback="0" onerror="${spriteOnerror(item.Name)}" alt="${item.Name}" loading="lazy" />
        <div class="poke-card__name">${item.Name}</div>
      </div>`;
    })
    .join("");

  if (append) {
    const countBefore = grid.querySelectorAll(".poke-card").length;
    grid.insertAdjacentHTML("beforeend", html);
    // Attach click listeners and animate only new cards
    const allCards = grid.querySelectorAll<HTMLDivElement>(".poke-card");
    allCards.forEach((card, i) => {
      if (i >= countBefore) {
        card.addEventListener("click", () => {
          const name = card.dataset.name;
          if (name) showDetail(name);
        });
      }
    });
    staggerNewCards(grid, countBefore);
  } else {
    grid.innerHTML = html;
    grid.querySelectorAll<HTMLDivElement>(".poke-card").forEach((card) => {
      card.addEventListener("click", () => {
        const name = card.dataset.name;
        if (name) showDetail(name);
      });
    });
    staggerCards(grid);
  }
}

async function renderCurrentView(items: PokemonListItem[]): Promise<void> {
  lastRenderedItems = items;
  if (viewMode === "table") {
    await renderTable(items);
  } else {
    renderGrid(items, false);
  }
}

async function renderTable(items: PokemonListItem[]): Promise<void> {
  if (!items || items.length === 0) {
    grid.innerHTML = `<p class="loading">${t("pokedex.noResults")}</p>`;
    return;
  }

  createInlineDiglett(grid, t("pokedex.loadingData"));

  const pokemonData: Pokemon[] = await Promise.all(
    items.map(async (item) => {
      if (pokemonDataCache.has(item.Name)) {
        return pokemonDataCache.get(item.Name)!;
      }
      const p = await GetPokemon(item.Name);
      pokemonDataCache.set(item.Name, p);
      return p;
    }),
  );

  const sortedData = (currentSortColumn && currentSortDirection)
    ? sortPokemonData(pokemonData, currentSortColumn, currentSortDirection)
    : pokemonData;

  const statColKeys = ["hp", "atk", "def", "spa", "spd", "vel"];
  const rows = sortedData
    .map((p) => {
      const sprite = p.Sprites.FrontDefault
        ? `<img class="poke-table__sprite" src="${p.Sprites.FrontDefault}" alt="${p.Name}" loading="lazy" />`
        : "";
      const types = (p.Types || [])
        .map((t) => `<span class="type-badge type-badge--icon-only type-${t.Name}" title="${typeName(t.Name)}"><img src="/assets/types/${t.Name}.svg" alt="${typeName(t.Name)}" class="type-icon"></span>`)
        .join(" ");
      const stats = (p.Stats || []).map((s) => s.BaseStat);
      const total = stats.reduce((a, b) => a + b, 0);
      const statCells = (p.Stats || [])
        .map((s, i) => `<td class="stat-cell" data-col="${statColKeys[i]}">${s.BaseStat}</td>`)
        .join("");

      return `<tr class="poke-table__row" data-name="${p.Name}">
        <td class="poke-table__id" data-col="id">#${String(p.ID).padStart(3, "0")}</td>
        <td data-col="sprite">${sprite}</td>
        <td class="poke-table__name" data-col="name">${p.Name}</td>
        <td data-col="types">${types}</td>
        ${statCells}
        <td class="stat-cell stat-total" data-col="total">${total}</td>
      </tr>`;
    })
    .join("");

  const sortableStats: [string, SortColumn][] = [
    [t("pokedex.columns.hp"), "hp"], [t("pokedex.columns.atk"), "atk"], [t("pokedex.columns.def"), "def"],
    [t("pokedex.columns.spa"), "spa"], [t("pokedex.columns.spd"), "spd"], [t("pokedex.columns.vel"), "vel"],
  ];
  const statHeaders = sortableStats
    .map(([label, col]) => {
      const ind = currentSortColumn === col ? (currentSortDirection === 'asc' ? 'asc' : currentSortDirection === 'desc' ? 'desc' : '') : '';
      const activeClass = currentSortColumn === col ? ' active' : '';
      return `<th class="stat-cell sortable${activeClass}" data-sort="${col}" data-col="${col}">${label} <span class="sort-indicator ${ind}"></span></th>`;
    })
    .join("");

  const idInd = currentSortColumn === 'id' ? (currentSortDirection === 'asc' ? 'asc' : currentSortDirection === 'desc' ? 'desc' : '') : '';
  const nameInd = currentSortColumn === 'name' ? (currentSortDirection === 'asc' ? 'asc' : currentSortDirection === 'desc' ? 'desc' : '') : '';
  const totalInd = currentSortColumn === 'total' ? (currentSortDirection === 'asc' ? 'asc' : currentSortDirection === 'desc' ? 'desc' : '') : '';

  grid.innerHTML = `<table class="poke-table" data-table-id="pokedex-stats">
    <thead><tr>
      <th class="sortable${currentSortColumn === 'id' ? ' active' : ''}" data-sort="id" data-col="id">${t("pokedex.columns.id")} <span class="sort-indicator ${idInd}"></span></th>
      <th data-col="sprite"></th>
      <th class="sortable${currentSortColumn === 'name' ? ' active' : ''}" data-sort="name" data-col="name">${t("pokedex.columns.name")} <span class="sort-indicator ${nameInd}"></span></th>
      <th data-col="types">${t("pokedex.columns.types")}</th>
      ${statHeaders}
      <th class="stat-cell sortable${currentSortColumn === 'total' ? ' active' : ''}" data-sort="total" data-col="total">${t("pokedex.columns.total")} <span class="sort-indicator ${totalInd}"></span></th>
    </tr></thead>
    <tbody>${rows}</tbody>
  </table>`;

  initColumnToggle("pokedex-stats", pokedexTableColumns());

  grid.querySelectorAll<HTMLTableRowElement>(".poke-table__row").forEach((row) => {
    row.addEventListener("click", () => {
      const name = row.dataset.name;
      if (name) showDetail(name);
    });
  });

  grid.querySelectorAll<HTMLTableCellElement>("th.sortable").forEach((th) => {
    th.addEventListener("click", async (e) => {
      e.stopPropagation();
      if (sortingLoading) return;
      const col = th.dataset.sort as SortColumn;
      if (currentSortColumn === col) {
        if (currentSortDirection === 'asc') {
          currentSortDirection = 'desc';
        } else if (currentSortDirection === 'desc') {
          currentSortDirection = null;
          currentSortColumn = null;
          sortedFullList = null;
          offset = 0;
          if (hasFilter()) {
            await renderCurrentView(filteredList.slice(offset, offset + rowLimit));
          } else {
            await loadList();
          }
          updateRowLimitControl();
          return;
        }
      } else {
        currentSortColumn = col;
        currentSortDirection = 'asc';
      }

      sortingLoading = true;
      showSortingOverlay();
      try {
        const allPokemon = await ensureAllPokemonLoaded();
        updateSortingOverlayText(t("pokedex.sorting"));
        await new Promise((r) => requestAnimationFrame(r));
        pokemonSortCache.setData(allPokemon);
        sortedFullList = pokemonSortCache.get(currentSortColumn!, currentSortDirection as 'asc' | 'desc');
        offset = 0;
        hideSortingOverlay();
        await renderTable(getCurrentPageItems());
        updateRowLimitControl();
      } finally {
        hideSortingOverlay();
        sortingLoading = false;
      }
    });
  });

  staggerTableRows(grid);
}

function staggerTableRows(container: HTMLElement): void {
  const rows = container.querySelectorAll(".poke-table__row");
  if (rows.length === 0) return;

  import("gsap").then(({ default: gsap }) => {
    gsap.fromTo(
      rows,
      { opacity: 0, x: -20 },
      { opacity: 1, x: 0, duration: 0.2, stagger: 0.02, ease: "power2.out" },
    );
  });
}

function updateRowLimitControl(): void {
  limitDisplay.textContent = String(rowLimit);
  limitDecreaseBtn.disabled = rowLimit <= LIMIT_STEP;

  const total = getTotalItemCount();
  const totalPages = Math.max(1, Math.ceil(total / rowLimit));
  const currentPage = Math.floor(offset / rowLimit) + 1;

  pageIndicator.textContent = t("pokedex.page", { page: String(currentPage), pages: String(totalPages) });
  pagePrevBtn.disabled = currentPage <= 1;
  pageNextBtn.disabled = currentPage >= totalPages;
}

function getTotalItemCount(): number {
  if (sortedFullList) return sortedFullList.length;
  if (hasFilter() && filteredList.length > 0) return filteredList.length;
  return totalCount;
}

// -- Paginación --------------------------------------------------------------

function resetSorting(): void {
  currentSortColumn = null;
  currentSortDirection = null;
  sortedFullList = null;
  pokemonSortCache.invalidate();
}

async function changeRowLimit(delta: number): Promise<void> {
  const newLimit = rowLimit + delta;
  if (newLimit < LIMIT_STEP) return;
  rowLimit = newLimit;
  offset = 0;
  updateRowLimitControl();
  if (sortedFullList) {
    // Re-sort with new limit
    const pageItems = sortedFullList.slice(0, rowLimit).map((p) => ({ Name: p.Name, URL: "" }));
    await renderCurrentView(pageItems);
  } else if (hasFilter()) {
    await renderCurrentView(filteredList.slice(0, rowLimit));
  } else {
    await loadList();
  }
}

async function goToPrevPage(): Promise<void> {
  if (offset <= 0) return;
  offset -= rowLimit;
  if (offset < 0) offset = 0;
  updateRowLimitControl();
  await renderCurrentView(getCurrentPageItems());
}

async function goToNextPage(): Promise<void> {
  const total = getTotalItemCount();
  if (offset + rowLimit >= total) return;
  offset += rowLimit;
  updateRowLimitControl();

  const items = getCurrentPageItems();
  if (items.length > 0) {
    await renderCurrentView(items);
  } else {
    // Need to fetch from API (no sorted/filtered list available)
    const data = await ListPokemon(offset, rowLimit);
    await renderCurrentView(data.Results);
  }
}

// -- Detalle -----------------------------------------------------------------

async function showDetail(name: string): Promise<void> {
  createInlineDiglett(detailEl, t("common.loading"));
  await showView(detailView, listView);
  try {
    const p = await GetPokemon(name);
    renderDetail(p);
  } catch (err: unknown) {
    detailEl.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}

async function renderDetail(p: Pokemon): Promise<void> {
  currentDetailPokemon = p;
  const types = (p.Types || [])
    .map((t) => `<span class="type-badge type-${t.Name}"><img src="/assets/types/${t.Name}.svg" alt="" class="type-icon">${typeName(t.Name)}</span>`)
    .join("");

  const sprites = `
    <div class="sprites">
      ${p.Sprites.FrontDefault ? `<div><img src="${p.Sprites.FrontDefault}" alt="default"/><span>${t("detail.normal")}</span></div>` : ""}
      ${p.Sprites.FrontShiny ? `<div><img src="${p.Sprites.FrontShiny}" alt="shiny"/><span>${t("detail.shiny")}</span></div>` : ""}
    </div>`;

  detailEl.innerHTML = `
    <h2>#${p.ID} ${p.Name}</h2>
    ${sprites}
    <div class="types">${types}</div>
    <p class="meta">${t("pokedex.height")}: ${p.Height / 10} m &nbsp;&middot;&nbsp; ${t("pokedex.weight")}: ${p.Weight / 10} kg</p>
    <div id="stats-chart" style="width:100%;height:420px;"></div>
    <div id="compare-controls" class="compare-controls">
      <button id="compare-btn" class="btn btn-secondary">${t("detail.compare")}</button>
    </div>
    <div id="chart-legend" class="chart-legend hidden"></div>
    <div id="stat-totals" class="stat-totals hidden"></div>
    <div id="pokemon-lore" class="pokemon-lore"><p class="loading">${t("detail.loadingLore")}</p></div>
    <div id="pokemon-encounters" class="pokemon-encounters"><p class="loading">${t("encounters.loading")}</p></div>
    <div id="pokemon-abilities" class="pokemon-abilities"><p class="loading">${t("common.loading")}</p></div>
    <div id="pokemon-moves" class="pokemon-moves"></div>`;

  const chartContainer = document.getElementById("stats-chart") as HTMLDivElement;
  const { renderStatsChart } = await import("../charts/stats-chart");

  const COMPARE_COLORS = [
    "#e53e3e", // base — red
    "#3182ce", // blue
    "#38a169", // green
    "#d69e2e", // gold
    "#805ad5", // purple
    "#dd6b20", // orange
    "#319795", // teal
    "#e53e8e", // pink
  ];

  const STAT_COLORS: Record<string, string> = {
    hp: "#ff5959",
    attack: "#f5ac78",
    defense: "#fae078",
    "special-attack": "#9db7f5",
    "special-defense": "#a7db8d",
    speed: "#fa92b2",
  };

  const primarySeries = { label: p.Name, stats: p.Stats || [], color: COMPARE_COLORS[0] };
  const comparedSeries = [primarySeries];
  renderStatsChart(chartContainer, comparedSeries);

  const compareBtn = document.getElementById("compare-btn") as HTMLButtonElement;
  const compareControls = document.getElementById("compare-controls") as HTMLDivElement;
  const chartLegend = document.getElementById("chart-legend") as HTMLDivElement;
  const statTotalsEl = document.getElementById("stat-totals") as HTMLDivElement;

  function rebuildLegend(): void {
    if (comparedSeries.length === 1) {
      chartLegend.classList.add("hidden");
      chartLegend.innerHTML = "";
      return;
    }
    chartLegend.innerHTML = comparedSeries
      .map((s, i) => {
        const removeBtn = i === 0
          ? ""
          : `<button class="compare-remove-btn" data-index="${i}" aria-label="Remove">×</button>`;
        return `<span class="chart-legend-entry" style="color:${s.color}">■ ${s.label}${removeBtn}</span>`;
      })
      .join("");
    chartLegend.classList.remove("hidden");

    chartLegend.querySelectorAll<HTMLButtonElement>(".compare-remove-btn").forEach((btn) => {
      btn.addEventListener("click", () => {
        const idx = parseInt(btn.dataset.index ?? "0", 10);
        comparedSeries.splice(idx, 1);
        renderStatsChart(chartContainer, comparedSeries);
        rebuildLegend();
        rebuildStatTotals();
      });
    });
  }

  function bst(stats: { BaseStat: number }[]): number {
    return stats.reduce((sum, s) => sum + s.BaseStat, 0);
  }

  function rebuildStatTotals(): void {
    if (comparedSeries.length < 2) {
      statTotalsEl.classList.add("hidden");
      statTotalsEl.innerHTML = "";
      return;
    }

    const baseBst = bst(comparedSeries[0].stats);

    const bars = comparedSeries.map((s, idx) => {
      const total = bst(s.stats);
      const diff = total - baseBst;

      const segments = s.stats
        .map((stat) => {
          const pct = (stat.BaseStat / total) * 100;
          const color = STAT_COLORS[stat.Name] || "#999";
          return `<div class="bst-segment" data-stat="${stat.Name}" style="width:${pct}%;background:${color}" title="${statName(stat.Name)}: ${stat.BaseStat}"></div>`;
        })
        .join("");

      const diffRow = idx === 0
        ? ""
        : `<div class="bst-diff-row">
            <span class="bst-diff-spacer"></span>
            <div class="bst-diff-bar">${s.stats.map((stat, i) => {
              const d = stat.BaseStat - comparedSeries[0].stats[i].BaseStat;
              const pct = (stat.BaseStat / total) * 100;
              const cls = d > 0 ? "bst-diff--pos" : d < 0 ? "bst-diff--neg" : "bst-diff--neutral";
              return `<div class="bst-diff-segment ${cls}" style="width:${pct}%">${d !== 0 ? (d > 0 ? "+" : "") + d : ""}</div>`;
            }).join("")}</div>
            <span class="bst-diff-total ${diff > 0 ? "bst-diff--pos" : diff < 0 ? "bst-diff--neg" : ""}">${diff > 0 ? "+" : ""}${diff}</span>
          </div>`;

      return `<div class="bst-group">
        <div class="bst-row">
          <span class="bst-label">${s.label}</span>
          <div class="bst-bar">${segments}</div>
          <span class="bst-value">${total}</span>
        </div>
        ${diffRow}
      </div>`;
    });

    // Stat color legend
    const statLegend = comparedSeries[0].stats
      .map((s) => `<span class="bst-stat-legend-entry"><span class="bst-stat-swatch" style="background:${STAT_COLORS[s.Name] || "#999"}"></span>${statName(s.Name)}</span>`)
      .join("");

    // KPIs — only when exactly 2 Pokémon
    let kpisHtml = "";
    if (comparedSeries.length === 2) {
      const baseStats = comparedSeries[0].stats;
      const otherStats = comparedSeries[1].stats;

      // Highest stat for each
      const baseMax = baseStats.reduce((best, s) => (s.BaseStat > best.BaseStat ? s : best), baseStats[0]);
      const otherMax = otherStats.reduce((best, s) => (s.BaseStat > best.BaseStat ? s : best), otherStats[0]);

      // Biggest advantage / disadvantage (base vs other)
      let bestAdv = { name: "", diff: 0 };
      let worstAdv = { name: "", diff: 0 };
      for (let i = 0; i < baseStats.length; i++) {
        const d = baseStats[i].BaseStat - otherStats[i].BaseStat;
        if (d > bestAdv.diff) bestAdv = { name: baseStats[i].Name, diff: d };
        if (d < worstAdv.diff) worstAdv = { name: baseStats[i].Name, diff: d };
      }

      kpisHtml = `<div class="bst-kpis">
        <div class="bst-kpi"><span class="bst-kpi-label">${t("detail.highestStat")}</span> <strong>${comparedSeries[0].label}</strong> — ${statName(baseMax.Name)} (${baseMax.BaseStat})</div>
        <div class="bst-kpi"><span class="bst-kpi-label">${t("detail.highestStat")}</span> <strong>${comparedSeries[1].label}</strong> — ${statName(otherMax.Name)} (${otherMax.BaseStat})</div>
        ${bestAdv.diff > 0 ? `<div class="bst-kpi"><span class="bst-kpi-label">${t("detail.biggestAdvantage")}</span> ${statName(bestAdv.name)} (+${bestAdv.diff})</div>` : ""}
        ${worstAdv.diff < 0 ? `<div class="bst-kpi"><span class="bst-kpi-label">${t("detail.biggestDisadvantage")}</span> ${statName(worstAdv.name)} (${worstAdv.diff})</div>` : ""}
      </div>`;
    }

    statTotalsEl.innerHTML = `
      <h4 class="bst-title">${t("detail.statTotals")}</h4>
      <div class="bst-stat-legend">${statLegend}</div>
      ${bars.join("")}
      ${kpisHtml}
    `;
    statTotalsEl.classList.remove("hidden");

    // Click on a segment to show stat comparison tooltip
    statTotalsEl.querySelectorAll<HTMLDivElement>(".bst-segment").forEach((seg) => {
      seg.style.cursor = "pointer";
      seg.addEventListener("click", (e) => {
        const sName = seg.dataset.stat!;
        const color = STAT_COLORS[sName] || "#999";

        // Remove any existing tooltip
        statTotalsEl.querySelector(".bst-tooltip")?.remove();

        const rows = comparedSeries.map((s) => {
          const val = s.stats.find((st) => st.Name === sName)?.BaseStat ?? 0;
          return `<div class="bst-tooltip-row"><span class="bst-tooltip-name">${s.label}</span><span class="bst-tooltip-val">${val}</span></div>`;
        }).join("");

        const tooltip = document.createElement("div");
        tooltip.className = "bst-tooltip";
        tooltip.innerHTML = `<div class="bst-tooltip-header" style="border-color:${color}">${statName(sName)}</div>${rows}`;

        // Position near the clicked segment
        const rect = seg.getBoundingClientRect();
        const parentRect = statTotalsEl.getBoundingClientRect();
        tooltip.style.left = `${rect.left - parentRect.left + rect.width / 2}px`;
        tooltip.style.top = `${rect.top - parentRect.top - 4}px`;

        statTotalsEl.appendChild(tooltip);

        // Close on click outside
        const close = (ev: MouseEvent) => {
          if (!tooltip.contains(ev.target as Node)) {
            tooltip.remove();
            document.removeEventListener("click", close);
          }
        };
        setTimeout(() => document.addEventListener("click", close), 0);

        e.stopPropagation();
      });
    });
  }

  compareBtn.addEventListener("click", () => {
    const existing = compareControls.querySelector(".compare-input-wrap");
    if (existing) {
      existing.remove();
      return;
    }

    const inputWrap = document.createElement("div");
    inputWrap.className = "compare-input-wrap";

    const input = document.createElement("input");
    input.type = "text";
    input.placeholder = t("detail.comparePlaceholder");
    input.className = "compare-input";
    inputWrap.appendChild(input);
    compareControls.appendChild(inputWrap);
    input.focus();

    createAutocomplete(input, allPokemonNames, async (name) => {
      inputWrap.remove();
      try {
        const p2 = await GetPokemon(name);
        const colorIndex = 1 + ((comparedSeries.length - 1) % (COMPARE_COLORS.length - 1));
        const color = COMPARE_COLORS[colorIndex];
        comparedSeries.push({ label: p2.Name, stats: p2.Stats || [], color });
        renderStatsChart(chartContainer, comparedSeries);
        rebuildLegend();
        rebuildStatTotals();
      } catch {
        // silently ignore fetch errors
      }
    });
  });

  loadLore(p.Name);
  loadEncounters(p.Name);
  loadAbilities(p.Abilities || []);
  renderMoves(p.Moves || []);
}

async function loadLore(name: string): Promise<void> {
  const loreEl = document.getElementById("pokemon-lore") as HTMLDivElement;
  if (!loreEl) return;

  try {
    const species = await GetPokemonSpecies(name, getLocale());

    const flavorText = cleanFlavorText(species.FlavorText);
    const badges: string[] = [];
    if (species.IsLegendary) badges.push(`<span class="lore-badge lore-badge--legendary">${t("detail.legendary")}</span>`);
    if (species.IsMythical) badges.push(`<span class="lore-badge lore-badge--mythical">${t("detail.mythical")}</span>`);

    const infoRows: string[] = [];
    if (species.Genus) infoRows.push(`<div class="lore-info-row"><span class="lore-label">${t("detail.category")}</span><span class="lore-value">${species.Genus}</span></div>`);
    if (species.Habitat) infoRows.push(`<div class="lore-info-row"><span class="lore-label">${t("detail.habitat")}</span><span class="lore-value">${capitalize(species.Habitat)}</span></div>`);
    if (species.Color) infoRows.push(`<div class="lore-info-row"><span class="lore-label">${t("detail.color")}</span><span class="lore-value">${capitalize(species.Color)}</span></div>`);
    if (species.Shape) infoRows.push(`<div class="lore-info-row"><span class="lore-label">${t("detail.shape")}</span><span class="lore-value">${capitalize(species.Shape)}</span></div>`);

    loreEl.innerHTML = `
      <h3>${t("detail.lore")}</h3>
      ${badges.length > 0 ? `<div class="lore-badges">${badges.join("")}</div>` : ""}
      ${flavorText ? `<p class="lore-flavor-text">${flavorText}</p>` : ""}
      ${infoRows.length > 0 ? `<div class="lore-info">${infoRows.join("")}</div>` : ""}
    `;
  } catch {
    loreEl.innerHTML = `<p class="lore-error">${t("detail.loreError")}</p>`;
  }
}

type EncounterSortColumn = "location" | "game" | "method" | "chance" | "levels" | "conditions" | null;
type EncounterSortDirection = "asc" | "desc" | null;

interface EncounterRow {
  location: string;
  game: string;
  method: string;
  chance: number;
  minLevel: number;
  maxLevel: number;
  conditions: string;
}

let encounterSortColumn: EncounterSortColumn = null;
let currentDetailPokemon: Pokemon | null = null;
let encounterSortDirection: EncounterSortDirection = null;
let encounterRows: EncounterRow[] = [];
let selectedGames: string[] = [];

function sortEncounterRows(rows: EncounterRow[], column: EncounterSortColumn, direction: EncounterSortDirection): EncounterRow[] {
  if (!column || !direction) return rows;
  const mult = direction === "asc" ? 1 : -1;
  return [...rows].sort((a, b) => {
    switch (column) {
      case "location": return mult * a.location.localeCompare(b.location);
      case "game": return mult * a.game.localeCompare(b.game);
      case "method": return mult * a.method.localeCompare(b.method);
      case "chance": return mult * (a.chance - b.chance);
      case "levels": return mult * (a.minLevel - b.minLevel || a.maxLevel - b.maxLevel);
      case "conditions": return mult * a.conditions.localeCompare(b.conditions);
      default: return 0;
    }
  });
}

function renderEncounterTbody(el: HTMLElement): void {
  const filtered = selectedGames.length > 0
    ? encounterRows.filter((r) => selectedGames.includes(r.game))
    : encounterRows;
  const sorted = sortEncounterRows(filtered, encounterSortColumn, encounterSortDirection);
  const tbody = el.querySelector(".encounters-table tbody");
  if (!tbody) return;
  tbody.innerHTML = sorted.map((r) => `<tr>
            <td data-col="location">${r.location}</td>
            <td data-col="game">${r.game}</td>
            <td data-col="method">${r.method}</td>
            <td data-col="chance">${r.chance}%</td>
            <td data-col="levels">${r.minLevel === r.maxLevel ? `Lv. ${r.minLevel}` : `Lv. ${r.minLevel}-${r.maxLevel}`}</td>
            <td data-col="conditions">${r.conditions}</td>
          </tr>`).join("");

  reapplyColumnVisibility("encounters");
}

function updateEncounterSortIndicators(el: HTMLElement): void {
  el.querySelectorAll<HTMLElement>("th.sortable").forEach((th) => {
    const col = th.dataset.col as EncounterSortColumn;
    const indicator = th.querySelector<HTMLElement>(".sort-indicator");
    if (!indicator) return;
    th.classList.toggle("active", col === encounterSortColumn);
    indicator.className = "sort-indicator";
    if (col === encounterSortColumn && encounterSortDirection) {
      indicator.classList.add(encounterSortDirection);
    }
  });
}

async function loadEncounters(name: string): Promise<void> {
  const el = document.getElementById("pokemon-encounters") as HTMLDivElement;
  if (!el) return;

  encounterSortColumn = null;
  encounterSortDirection = null;
  encounterRows = [];
  selectedGames = [];

  try {
    const encounters = await GetPokemonEncounters(name);

    if (!encounters || encounters.length === 0) {
      el.innerHTML = `
        <h3>${t("encounters.title")}</h3>
        <p class="encounters-empty">${t("encounters.empty")}</p>`;
      return;
    }

    for (const enc of encounters) {
      const location = formatLocationName(enc.LocationArea);
      for (const v of enc.Versions || []) {
        const game = capitalize(v.Version);
        for (const d of v.Details || []) {
          encounterRows.push({
            location,
            game,
            method: formatEncounterMethod(d.Method),
            chance: d.Chance,
            minLevel: d.MinLevel,
            maxLevel: d.MaxLevel,
            conditions: (d.Conditions || []).map((c: any) => c.Name).filter(Boolean).join(", ") || "\u2014",
          });
        }
      }
    }

    const availableGames = [...new Set(encounterRows.map((r) => r.game))].sort();

    el.innerHTML = `
      <h3>${t("encounters.title")}</h3>
      <div class="encounters-filter-bar">
        <div class="filter-dropdown" id="dropdown-encounter-game">
          <button class="filter-dropdown__trigger" type="button">${t("encounters.columns.game")}</button>
          <div class="filter-dropdown__panel" id="filter-encounter-game"></div>
        </div>
      </div>
      <div class="encounters-table-wrap">
        <table class="poke-table encounters-table" data-table-id="encounters">
          <thead>
            <tr>
              <th class="sortable" data-col="location">${t("encounters.columns.location")} <span class="sort-indicator"></span></th>
              <th class="sortable" data-col="game">${t("encounters.columns.game")} <span class="sort-indicator"></span></th>
              <th class="sortable" data-col="method">${t("encounters.columns.method")} <span class="sort-indicator"></span></th>
              <th class="sortable" data-col="chance">${t("encounters.columns.chance")} <span class="sort-indicator"></span></th>
              <th class="sortable" data-col="levels">${t("encounters.columns.levels")} <span class="sort-indicator"></span></th>
              <th class="sortable" data-col="conditions">${t("encounters.columns.conditions")} <span class="sort-indicator"></span></th>
            </tr>
          </thead>
          <tbody></tbody>
        </table>
      </div>`;

    // Setup game filter dropdown
    const gameDropdown = document.getElementById("dropdown-encounter-game")!;
    const gamePanel = document.getElementById("filter-encounter-game")!;

    gameDropdown.querySelector(".filter-dropdown__trigger")!.addEventListener("click", () => {
      toggleDropdown("dropdown-encounter-game");
    });

    availableGames.forEach((game) => {
      const chip = document.createElement("button");
      chip.className = "filter-chip";
      chip.dataset.value = game;
      chip.textContent = game;
      chip.addEventListener("click", () => {
        chip.classList.toggle("active");
        if (chip.classList.contains("active")) {
          selectedGames.push(game);
        } else {
          selectedGames = selectedGames.filter((g) => g !== game);
        }
        updateDropdownTrigger("dropdown-encounter-game", t("encounters.columns.game"), selectedGames.length);
        renderEncounterTbody(el);
      });
      gamePanel.appendChild(chip);
    });

    // Reset button
    const resetBtn = document.createElement("button");
    resetBtn.className = "filter-chip filter-chip--reset";
    resetBtn.textContent = t("encounters.reset");
    resetBtn.addEventListener("click", () => {
      selectedGames = [];
      gamePanel.querySelectorAll(".filter-chip.active").forEach((c) => c.classList.remove("active"));
      updateDropdownTrigger("dropdown-encounter-game", t("encounters.columns.game"), 0);
      renderEncounterTbody(el);
    });
    gamePanel.insertBefore(resetBtn, gamePanel.firstChild);

    renderEncounterTbody(el);
    initColumnToggle("encounters", encounterTableColumns());

    el.querySelectorAll<HTMLElement>("th.sortable").forEach((th) => {
      th.addEventListener("click", () => {
        const col = th.dataset.col as EncounterSortColumn;
        if (encounterSortColumn === col) {
          if (encounterSortDirection === "asc") {
            encounterSortDirection = "desc";
          } else if (encounterSortDirection === "desc") {
            encounterSortColumn = null;
            encounterSortDirection = null;
          }
        } else {
          encounterSortColumn = col;
          encounterSortDirection = "asc";
        }
        updateEncounterSortIndicators(el);
        renderEncounterTbody(el);
      });
    });
  } catch {
    el.innerHTML = `
      <h3>${t("encounters.title")}</h3>
      <p class="encounters-error">${t("encounters.error")}</p>`;
  }
}

// -- Abilities section -------------------------------------------------------

async function loadAbilities(abilityNames: string[]): Promise<void> {
  const el = document.getElementById("pokemon-abilities") as HTMLDivElement;
  if (!el) return;

  if (!abilityNames || abilityNames.length === 0) {
    el.innerHTML = "";
    return;
  }

  try {
    const abilities = await Promise.all(abilityNames.map((name) => GetAbility(name)));
    const lang = getLocale() === "es" ? "es" : "en";

    const cards = abilities.map((ab) => {
      const name = lang === "es" && ab.NameEs ? ab.NameEs : ab.Name;
      const desc = lang === "es" && ab.DescriptionEs ? ab.DescriptionEs : ab.Description;
      return `<div class="ability-card">
        <span class="ability-name">${capitalize(name)}</span>
        ${desc ? `<p class="ability-desc">${desc}</p>` : ""}
      </div>`;
    }).join("");

    el.innerHTML = `
      <h3>${t("detail.abilities")}</h3>
      <div class="abilities-list">${cards}</div>
    `;
  } catch {
    el.innerHTML = `<h3>${t("detail.abilities")}</h3><p class="lore-error">${t("detail.abilitiesError")}</p>`;
  }
}

// -- Moves section -----------------------------------------------------------

type MoveMethodFilter = "all" | "level-up" | "machine" | "egg" | "tutor";

async function renderMoves(moves: PokemonMoveEntry[]): Promise<void> {
  const el = document.getElementById("pokemon-moves") as HTMLDivElement;
  if (!el) return;

  if (!moves || moves.length === 0) {
    el.innerHTML = "";
    return;
  }

  await loadMoveNames();

  let activeFilters: Set<MoveMethodFilter> = new Set();
  type SortCol = "name" | "method" | "level" | null;
  let sortCol: SortCol = null;
  let sortDir: "asc" | "desc" = "asc";

  const methodFilters: MoveMethodFilter[] = ["all", "level-up", "machine", "egg", "tutor"];

  function getFilterLabel(method: MoveMethodFilter): string {
    switch (method) {
      case "all": return t("detail.moveMethod.all");
      case "level-up": return t("detail.moveMethod.levelUp");
      case "machine": return t("detail.moveMethod.machine");
      case "egg": return t("detail.moveMethod.egg");
      case "tutor": return t("detail.moveMethod.tutor");
    }
  }

  function getMethodLabel(method: string): string {
    switch (method) {
      case "level-up": return t("detail.moveMethod.levelUp");
      case "machine": return t("detail.moveMethod.machine");
      case "egg": return t("detail.moveMethod.egg");
      case "tutor": return t("detail.moveMethod.tutor");
      default: return capitalize(method.replace(/-/g, " "));
    }
  }

  function renderTable(): string {
    const filtered = activeFilters.size === 0
      ? moves
      : moves.filter((m) => activeFilters.has(m.Method as MoveMethodFilter));

    let sorted = [...filtered];
    if (sortCol !== null) {
      sorted.sort((a, b) => {
        let cmp = 0;
        if (sortCol === "name") {
          cmp = getLocalizedMoveName(a.Name).localeCompare(getLocalizedMoveName(b.Name));
        } else if (sortCol === "method") {
          cmp = a.Method.localeCompare(b.Method);
        } else if (sortCol === "level") {
          const la = a.Method === "level-up" && a.Level > 0 ? a.Level : Infinity;
          const lb = b.Method === "level-up" && b.Level > 0 ? b.Level : Infinity;
          cmp = la - lb;
        }
        return sortDir === "asc" ? cmp : -cmp;
      });
    }

    if (sorted.length === 0) {
      return `<p class="moves-empty">${t("pokedex.noResults")}</p>`;
    }

    const rows = sorted.map((m) => {
      const levelCell = m.Method === "level-up" && m.Level > 0 ? String(m.Level) : "\u2014";
      return `<tr>
        <td class="move-name">${getLocalizedMoveName(m.Name)}</td>
        <td class="move-level">${levelCell}</td>
        <td class="move-method">${getMethodLabel(m.Method)}</td>
      </tr>`;
    }).join("");

    function thClass(col: SortCol): string {
      if (sortCol !== col) return 'class="sort-none"';
      return `class="${sortDir === "asc" ? "sort-asc" : "sort-desc"}"`;
    }

    return `<div class="moves-table-wrap">
      <table class="poke-table moves-table">
        <thead><tr>
          <th ${thClass("name")} data-sort="name">${t("detail.moveName")}</th>
          <th ${thClass("level")} data-sort="level">${t("detail.moveLevelCol")}</th>
          <th ${thClass("method")} data-sort="method">${t("detail.moveMethodCol")}</th>
        </tr></thead>
        <tbody>${rows}</tbody>
      </table>
    </div>`;
  }

  function rebuildTable(): void {
    const tableWrap = el.querySelector(".moves-body");
    if (tableWrap) {
      tableWrap.innerHTML = renderTable();
      tableWrap.querySelectorAll<HTMLTableCellElement>("thead th[data-sort]").forEach((th) => {
        th.addEventListener("click", () => {
          const col = th.dataset.sort as SortCol;
          if (sortCol === col) {
            if (sortDir === "asc") {
              sortDir = "desc";
            } else {
              sortCol = null;
            }
          } else {
            sortCol = col;
            sortDir = "asc";
          }
          rebuildTable();
        });
      });
    }
  }

  function rebuildFilterBtns(): void {
    const bar = el.querySelector(".moves-filter-bar");
    if (!bar) return;
    bar.innerHTML = methodFilters.map((mf) => {
      const isAll = mf === "all";
      const isActive = isAll ? activeFilters.size === 0 : activeFilters.has(mf);
      return `<button class="filter-chip${isActive ? " active" : ""}" data-method="${mf}">${getFilterLabel(mf)}</button>`;
    }).join("");

    bar.querySelectorAll<HTMLButtonElement>(".filter-chip").forEach((btn) => {
      btn.addEventListener("click", () => {
        const method = btn.dataset.method as MoveMethodFilter;
        if (method === "all") {
          activeFilters.clear();
        } else if (activeFilters.has(method)) {
          activeFilters.delete(method);
        } else {
          activeFilters.add(method);
        }
        sortCol = null;
        rebuildFilterBtns();
        rebuildTable();
      });
    });
  }

  el.innerHTML = `
    <h3>${t("detail.moves")}</h3>
    <div class="moves-filter-bar"></div>
    <div class="moves-body"></div>
  `;

  rebuildFilterBtns();
  rebuildTable();
}

function formatLocationName(name: string): string {
  return name.replace(/-/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

function formatEncounterMethod(method: string): string {
  const key = `encounters.methods.${method}`;
  const translated = t(key);
  if (translated !== key) return translated;
  return method.replace(/-/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

function cleanFlavorText(text: string): string {
  return text.replace(/[\n\f\r]/g, " ").replace(/\s+/g, " ").trim();
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

// -- Búsqueda ----------------------------------------------------------------

async function search(): Promise<void> {
  const query = searchInput.value.trim().toLowerCase();
  if (!query) {
    resetInfiniteScroll();
    filter = { generations: [], types: [], legendary: false, mythical: false };
    filteredList = [];
    resetSorting();
    resetFilterUI();
    loadList();
    return;
  }
  await showDetail(query);
}

// -- Filtros -----------------------------------------------------------------

function updateDropdownTrigger(dropdownId: string, label: string, count: number): void {
  const dropdown = document.getElementById(dropdownId);
  if (!dropdown) return;
  const trigger = dropdown.querySelector(".filter-dropdown__trigger") as HTMLButtonElement;
  trigger.textContent = count > 0 ? `${label} (${count})` : label;
  trigger.classList.toggle("filter-dropdown__trigger--active", count > 0);
}

function toggleDropdown(dropdownId: string): void {
  const dropdown = document.getElementById(dropdownId)!;
  const isOpen = dropdown.classList.contains("open");

  // Close all dropdowns first
  document.querySelectorAll(".filter-dropdown.open").forEach((d) => d.classList.remove("open"));

  if (!isOpen) {
    dropdown.classList.add("open");
  }
}

function closeAllDropdowns(): void {
  document.querySelectorAll(".filter-dropdown.open").forEach((d) => d.classList.remove("open"));
}

async function populateFilters(): Promise<void> {
  try {
    const [gens, types] = await Promise.all([ListGenerations(), ListTypes()]);

    // Setup trigger click handlers
    const genDropdown = document.getElementById("dropdown-gen")!;
    const typeDropdown = document.getElementById("dropdown-type")!;

    genDropdown.querySelector(".filter-dropdown__trigger")!.addEventListener("click", () => {
      toggleDropdown("dropdown-gen");
    });
    typeDropdown.querySelector(".filter-dropdown__trigger")!.addEventListener("click", () => {
      toggleDropdown("dropdown-type");
    });

    // Close dropdowns on outside click
    document.addEventListener("click", (e) => {
      const target = e.target as HTMLElement;
      if (!target.closest(".filter-dropdown")) {
        closeAllDropdowns();
      }
    });

    gens.forEach((g) => {
      const chip = document.createElement("button");
      chip.className = "filter-chip";
      chip.dataset.value = g.Name;
      chip.textContent = formatGenName(g.Name);
      chip.addEventListener("click", () => {
        chip.classList.toggle("active");
        if (chip.classList.contains("active")) {
          filter.generations.push(g.Name);
        } else {
          filter.generations = filter.generations.filter((v) => v !== g.Name);
        }
        updateDropdownTrigger("dropdown-gen", t("filters.generation"), filter.generations.length);
        applyFilters();
      });
      filterGenContainer.appendChild(chip);
    });

    types.Results.forEach((tp) => {
      if (tp.Name === "shadow" || tp.Name === "unknown") return;
      const chip = document.createElement("button");
      chip.className = "filter-chip";
      chip.dataset.value = tp.Name;
      chip.textContent = typeName(tp.Name);
      chip.addEventListener("click", () => {
        chip.classList.toggle("active");
        if (chip.classList.contains("active")) {
          filter.types.push(tp.Name);
        } else {
          filter.types = filter.types.filter((v) => v !== tp.Name);
        }
        updateDropdownTrigger("dropdown-type", t("filters.type"), filter.types.length);
        applyFilters();
      });
      filterTypeContainer.appendChild(chip);
    });
  } catch (err) {
    console.error(t("pokedex.errorLoadingFilters"), err);
  }
}

function formatGenName(name: string): string {
  const match = name.match(/generation-(.+)/);
  if (!match) return name;
  return "Gen " + match[1].toUpperCase();
}

async function applyFilters(): Promise<void> {
  if (isFiltering) return;
  isFiltering = true;
  resetInfiniteScroll();
  filteredList = [];

  showSortingOverlay(t("pokedex.applyingFilters"));
  try {
    if (hasFilter()) {
      await loadFiltered();
    } else {
      resetSorting();
      await loadList();
    }
  } finally {
    hideSortingOverlay();
    isFiltering = false;
  }
}

function resetFilterUI(): void {
  filterGenContainer.querySelectorAll(".filter-chip.active").forEach((chip) => {
    chip.classList.remove("active");
  });
  filterTypeContainer.querySelectorAll(".filter-chip.active").forEach((chip) => {
    chip.classList.remove("active");
  });
  filter.generations = [];
  filter.types = [];
  filter.legendary = false;
  filter.mythical = false;
  filterLegendaryBtn.classList.remove("active");
  filterMythicalBtn.classList.remove("active");
  closeAllDropdowns();
  updateDropdownTrigger("dropdown-gen", t("filters.generation"), 0);
  updateDropdownTrigger("dropdown-type", t("filters.type"), 0);
}

// -- Helpers -----------------------------------------------------------------

function idFromURL(url: string): string {
  const parts = url.replace(/\/$/, "").split("/");
  return parts[parts.length - 1];
}

function spriteURL(name: string): string {
  const safeName = name.toLowerCase().replace(/[^a-z0-9-]/g, "");
  return `https://img.pokemondb.net/sprites/black-white/normal/${safeName}.png`;
}

function spriteOnerror(name: string): string {
  const safeName = name.toLowerCase().replace(/[^a-z0-9-]/g, "");
  const fb1 = `https://img.pokemondb.net/sprites/x-y/normal/${safeName}.png`;
  const fb2 = `https://img.pokemondb.net/sprites/home/normal/1x/${safeName}.png`;
  const fb3 = `/assets/sprites/home-normal/${safeName}.png`;
  return `var f=parseInt(this.dataset.fallback||'0');if(f===0){this.dataset.fallback='1';this.src='${fb1}'}else if(f===1){this.dataset.fallback='2';this.src='${fb2}'}else if(f===2){this.dataset.fallback='3';this.src='${fb3}'}else{this.onerror=null;this.style.visibility='hidden'}`;
}

// -- Infinite scroll ---------------------------------------------------------

function staggerNewCards(container: HTMLElement, startIndex: number): void {
  const cards = container.querySelectorAll(".poke-card");
  const newCards = Array.from(cards).slice(startIndex);
  if (newCards.length === 0) return;

  import("gsap").then(({ default: gsap }) => {
    gsap.fromTo(
      newCards,
      { opacity: 0, y: 15 },
      { opacity: 1, y: 0, duration: 0.25, stagger: 0.03, ease: "power2.out" },
    );
  });
}

async function loadNextBatch(): Promise<void> {
  if (infiniteLoading || viewMode !== "grid") return;
  const total = sortedFullList ? sortedFullList.length : totalCount;
  if (offset + rowLimit >= total) return;

  infiniteLoading = true;
  offset += rowLimit;

  try {
    let items: PokemonListItem[];
    if (sortedFullList) {
      items = getCurrentPageItems();
    } else if (hasFilter()) {
      items = filteredList.slice(offset, offset + rowLimit);
    } else {
      const data = await ListPokemon(offset, rowLimit);
      items = data.Results;
    }
    renderGrid(items, true);
    updateScrollSentinel();
  } catch (err: unknown) {
    console.error("Error loading next batch:", err);
    offset -= rowLimit;
  } finally {
    infiniteLoading = false;
  }
}

function setupScrollObserver(): void {
  if (scrollObserver) scrollObserver.disconnect();
  scrollObserver = new IntersectionObserver(
    (entries) => {
      if (entries[0].isIntersecting) {
        loadNextBatch();
      }
    },
    { rootMargin: "200px" },
  );
  scrollObserver.observe(scrollSentinel);
}

function disconnectScrollObserver(): void {
  if (scrollObserver) {
    scrollObserver.disconnect();
    scrollObserver = null;
  }
}

function updateScrollSentinel(): void {
  const total = sortedFullList ? sortedFullList.length : totalCount;
  const allLoaded = offset + rowLimit >= total;
  scrollSentinel.classList.toggle("hidden", allLoaded);
}

function updateRowLimitVisibility(): void {
  if (viewMode === "grid") {
    rowLimitControl.classList.add("hidden");
    scrollSentinel.classList.remove("hidden");
    setupScrollObserver();
  } else {
    rowLimitControl.classList.remove("hidden");
    scrollSentinel.classList.add("hidden");
    disconnectScrollObserver();
  }
}

function resetInfiniteScroll(): void {
  offset = 0;
  grid.innerHTML = "";
  updateScrollSentinel();
}

// -- Current page items helper -----------------------------------------------

let lastRenderedItems: PokemonListItem[] = [];

function getCurrentPageItems(): PokemonListItem[] {
  if (sortedFullList) {
    const pageItems = sortedFullList.slice(offset, offset + rowLimit);
    return pageItems.map((p) => ({ Name: p.Name, URL: "" }));
  }
  if (hasFilter() && filteredList.length > 0) {
    return filteredList.slice(offset, offset + rowLimit);
  }
  return lastRenderedItems;
}

// -- Init --------------------------------------------------------------------

export function initPokedex(): void {
  grid = document.getElementById("pokemon-grid") as HTMLDivElement;
  listView = document.getElementById("list-view") as HTMLElement;
  detailView = document.getElementById("detail-view") as HTMLElement;
  detailEl = document.getElementById("pokemon-detail") as HTMLDivElement;
  limitDecreaseBtn = document.getElementById("limit-decrease") as HTMLButtonElement;
  limitIncreaseBtn = document.getElementById("limit-increase") as HTMLButtonElement;
  limitDisplay = document.getElementById("limit-display") as HTMLSpanElement;
  rowLimitControl = document.getElementById("row-limit-control") as HTMLDivElement;
  pagePrevBtn = document.getElementById("page-prev") as HTMLButtonElement;
  pageNextBtn = document.getElementById("page-next") as HTMLButtonElement;
  pageIndicator = document.getElementById("page-indicator") as HTMLSpanElement;
  searchInput = document.getElementById("search-input") as HTMLInputElement;
  searchBtn = document.getElementById("search-btn") as HTMLButtonElement;
  backBtn = document.getElementById("back-btn") as HTMLButtonElement;
  filterGenContainer = document.getElementById("filter-gen") as HTMLDivElement;
  filterTypeContainer = document.getElementById("filter-type") as HTMLDivElement;
  filterLegendaryBtn = document.getElementById("filter-legendary") as HTMLButtonElement;
  filterMythicalBtn = document.getElementById("filter-mythical") as HTMLButtonElement;
  filterResetBtn = document.getElementById("filter-reset") as HTMLButtonElement;
  viewToggleBtn = document.getElementById("view-toggle-btn") as HTMLButtonElement;
  scrollSentinel = document.getElementById("scroll-sentinel") as HTMLDivElement;

  limitDecreaseBtn.addEventListener("click", () => changeRowLimit(-LIMIT_STEP));
  limitIncreaseBtn.addEventListener("click", () => changeRowLimit(LIMIT_STEP));
  pagePrevBtn.addEventListener("click", () => goToPrevPage());
  pageNextBtn.addEventListener("click", () => goToNextPage());

  viewToggleBtn.addEventListener("click", async () => {
    const oldMode = viewMode;
    viewMode = viewMode === "grid" ? "table" : "grid";
    resetSorting();
    viewToggleBtn.textContent = viewMode === "grid" ? t("filters.tableView") : t("filters.cardView");
    updateRowLimitVisibility();

    if (oldMode === "grid" && viewMode === "table") {
      // Switch to table: reset offset for table pagination, render first page
      offset = 0;
      lastRenderedItems = [];
      showSortingOverlay();
      await morphToTable(grid, async () => {
        if (sortedFullList || (hasFilter() && filteredList.length > 0)) {
          await renderTable(getCurrentPageItems());
        } else {
          // Fetch fresh data from offset 0 to avoid stale lastRenderedItems from grid scroll
          const data = await ListPokemon(offset, rowLimit);
          totalCount = data.Count;
          await renderTable(data.Results);
        }
        hideSortingOverlay();
      });
      updateRowLimitControl();
    } else {
      // Switch to grid: reset and load with infinite scroll
      offset = 0;
      await morphToGrid(grid, () => {
        // Will be replaced by loadList below
      });
      if (hasFilter()) {
        renderGrid(filteredList.slice(0, rowLimit), false);
      } else {
        await loadList();
        return;
      }
      updateScrollSentinel();
    }
  });

  backBtn.addEventListener("click", async () => {
    hideSortingOverlay();
    const { disposeChart } = await import("../charts/stats-chart");
    disposeChart();
    await showView(listView, detailView);
  });

  searchBtn.addEventListener("click", search);
  searchInput.addEventListener("keydown", (e: KeyboardEvent) => {
    if (e.key === "Enter") search();
  });

  filterLegendaryBtn.addEventListener("click", async () => {
    if (isFiltering) return;
    isFiltering = true;
    filter.legendary = !filter.legendary;
    filterLegendaryBtn.classList.toggle("active", filter.legendary);
    resetInfiniteScroll();
    filteredList = [];
    showSortingOverlay(t("pokedex.applyingFilters"));
    try {
      if (hasFilter()) await loadFiltered();
      else await loadList();
    } finally {
      hideSortingOverlay();
      isFiltering = false;
    }
  });

  filterMythicalBtn.addEventListener("click", async () => {
    if (isFiltering) return;
    isFiltering = true;
    filter.mythical = !filter.mythical;
    filterMythicalBtn.classList.toggle("active", filter.mythical);
    resetInfiniteScroll();
    filteredList = [];
    showSortingOverlay(t("pokedex.applyingFilters"));
    try {
      if (hasFilter()) await loadFiltered();
      else await loadList();
    } finally {
      hideSortingOverlay();
      isFiltering = false;
    }
  });

  filterResetBtn.addEventListener("click", () => {
    if (isFiltering) return;
    resetInfiniteScroll();
    filter = { generations: [], types: [], legendary: false, mythical: false };
    filteredList = [];
    resetSorting();
    resetFilterUI();
    loadList();
  });

  populateFilters();
  updateRowLimitVisibility();
  loadList();

  // Load all pokemon names for autocomplete (non-blocking)
  ListPokemon(0, 2000).then((data) => {
    const names = data.Results.map((r) => r.Name);
    allPokemonNames = names;
    createAutocomplete(searchInput, names, (name) => {
      searchInput.value = name;
      showDetail(name);
    });
  });

  document.addEventListener("locale-changed", () => {
    viewToggleBtn.textContent = viewMode === "grid" ? t("filters.tableView") : t("filters.cardView");
    updateDropdownTrigger("dropdown-gen", t("filters.generation"), filter.generations.length);
    updateDropdownTrigger("dropdown-type", t("filters.type"), filter.types.length);
    if (viewMode === "table") {
      renderCurrentView(getCurrentPageItems());
    }
    updateRowLimitControl();
    if (currentDetailPokemon !== null && detailView.style.display !== "none") {
      renderDetail(currentDetailPokemon);
    }
  });
}
