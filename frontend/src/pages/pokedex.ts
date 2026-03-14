import {
  ListPokemon,
  GetPokemon,
  ListGenerations,
  GetGeneration,
  ListTypes,
  GetType,
  GetPokemonSpecies,
} from "../../wailsjs/go/app/App";
import type { Pokemon, PokemonListItem } from "../types";
import { showView, staggerCards, morphToTable, morphToGrid } from "../animations/transitions";
import { renderEVCalculatorForm, initEVCalculator } from "../ev-calculator";

const LIMIT = 20;

interface FilterState {
  generations: string[];
  types: string[];
  legendary: boolean;
  mythical: boolean;
}

let filter: FilterState = { generations: [], types: [], legendary: false, mythical: false };
let filteredList: PokemonListItem[] = [];
let offset = 0;
let totalCount = 0;
let viewMode: "grid" | "table" = "grid";

// -- Sorting state -----------------------------------------------------------

type SortDirection = 'asc' | 'desc' | null;
type SortColumn = 'id' | 'name' | 'hp' | 'atk' | 'def' | 'spa' | 'spd' | 'vel' | 'total' | null;

let currentSortColumn: SortColumn = null;
let currentSortDirection: SortDirection = null;
let sortedFullList: Pokemon[] | null = null;
let sortingLoading = false;

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
    grid.innerHTML = `<p class="loading">Cargando Pokémon para ordenar... ${pct}% (${loaded}/${total})</p>`;
  }

  return itemsToLoad.map((item) => pokemonDataCache.get(item.Name)!);
}

// Cache para evitar llamadas repetidas a GetPokemonSpecies
const legendaryCache = new Map<string, { isLegendary: boolean; isMythical: boolean }>();

// Cache para datos completos de Pokémon (tabla)
const pokemonDataCache = new Map<string, Pokemon>();

// DOM refs
let grid: HTMLDivElement;
let listView: HTMLElement;
let detailView: HTMLElement;
let detailEl: HTMLDivElement;
let prevBtn: HTMLButtonElement;
let nextBtn: HTMLButtonElement;
let pageInfo: HTMLSpanElement;
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
  grid.innerHTML = '<p class="loading">Cargando...</p>';
  try {
    const data = await ListPokemon(offset, LIMIT);
    totalCount = data.Count;
    await renderCurrentView(data.Results);
    updatePagination();
  } catch (err: unknown) {
    grid.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}

// -- Lista con filtro --------------------------------------------------------

async function loadFiltered(): Promise<void> {
  resetSorting();
  grid.innerHTML = '<p class="loading">Aplicando filtros...</p>';
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

    const page = filteredList.slice(0, LIMIT);
    await renderCurrentView(page);
    updatePagination();
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
    grid.innerHTML = `<p class="loading">Cargando lista completa... ${pct}% (${currentOffset}/${total})</p>`;
    const batch = await ListPokemon(currentOffset, BATCH_SIZE);
    all.push(...batch.Results);
    currentOffset += BATCH_SIZE;
  }

  return all;
}

async function filterByLegendary(list: PokemonListItem[]): Promise<PokemonListItem[]> {
  const BATCH = 10;
  const results: PokemonListItem[] = [];

  for (let i = 0; i < list.length; i += BATCH) {
    const batch = list.slice(i, i + BATCH);
    await Promise.allSettled(
      batch.map(async (item) => {
        if (!legendaryCache.has(item.Name)) {
          try {
            const species = await GetPokemonSpecies(item.Name);
            legendaryCache.set(item.Name, {
              isLegendary: species.IsLegendary,
              isMythical: species.IsMythical,
            });
          } catch {
            legendaryCache.set(item.Name, { isLegendary: false, isMythical: false });
          }
        }
      })
    );

    for (const item of batch) {
      const data = legendaryCache.get(item.Name);
      if (!data) continue;
      if (filter.legendary && filter.mythical) {
        if (data.isLegendary || data.isMythical) results.push(item);
      } else if (filter.legendary) {
        if (data.isLegendary) results.push(item);
      } else if (filter.mythical) {
        if (data.isMythical) results.push(item);
      }
    }
  }

  return results;
}

// -- Render ------------------------------------------------------------------

function renderGrid(items: PokemonListItem[]): void {
  lastRenderedItems = items;
  if (!items || items.length === 0) {
    grid.innerHTML = '<p class="loading">No se encontraron Pokémon.</p>';
    return;
  }

  grid.innerHTML = items
    .map((item) => {
      const id = idFromURL(item.URL);
      const sprite = spriteURL(id);
      const numId = parseInt(id, 10);
      const paddedId = isNaN(numId) ? "" : `#${String(numId).padStart(3, "0")}`;
      return `<div class="poke-card" data-name="${item.Name}">
        <span class="poke-card__number">${paddedId}</span>
        <img class="poke-card__sprite" src="${sprite}" alt="${item.Name}" loading="lazy" />
        <div class="poke-card__name">${item.Name}</div>
      </div>`;
    })
    .join("");

  grid.querySelectorAll<HTMLDivElement>(".poke-card").forEach((card) => {
    card.addEventListener("click", () => {
      const name = card.dataset.name;
      if (name) showDetail(name);
    });
  });

  staggerCards(grid);
}

async function renderCurrentView(items: PokemonListItem[]): Promise<void> {
  lastRenderedItems = items;
  if (viewMode === "table") {
    await renderTable(items);
  } else {
    renderGrid(items);
  }
}

async function renderTable(items: PokemonListItem[]): Promise<void> {
  if (!items || items.length === 0) {
    grid.innerHTML = '<p class="loading">No se encontraron Pokémon.</p>';
    return;
  }

  grid.innerHTML = '<p class="loading">Cargando datos...</p>';

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

  const sortedData = sortPokemonData(pokemonData, currentSortColumn, currentSortDirection);

  const rows = sortedData
    .map((p) => {
      const sprite = p.Sprites.FrontDefault
        ? `<img class="poke-table__sprite" src="${p.Sprites.FrontDefault}" alt="${p.Name}" loading="lazy" />`
        : "";
      const types = (p.Types || [])
        .map((t) => `<span class="type-badge type-${t.Name}">${t.Name}</span>`)
        .join(" ");
      const stats = (p.Stats || []).map((s) => s.BaseStat);
      const total = stats.reduce((a, b) => a + b, 0);
      const statCells = (p.Stats || [])
        .map((s) => `<td class="stat-cell">${s.BaseStat}</td>`)
        .join("");

      return `<tr class="poke-table__row" data-name="${p.Name}">
        <td class="poke-table__id">#${String(p.ID).padStart(3, "0")}</td>
        <td>${sprite}</td>
        <td class="poke-table__name">${p.Name}</td>
        <td>${types}</td>
        ${statCells}
        <td class="stat-cell stat-total">${total}</td>
      </tr>`;
    })
    .join("");

  const sortableStats: [string, SortColumn][] = [
    ["HP", "hp"], ["Atk", "atk"], ["Def", "def"],
    ["SpA", "spa"], ["SpD", "spd"], ["Vel", "vel"],
  ];
  const statHeaders = sortableStats
    .map(([label, col]) => {
      const ind = currentSortColumn === col ? (currentSortDirection === 'asc' ? 'asc' : currentSortDirection === 'desc' ? 'desc' : '') : '';
      const activeClass = currentSortColumn === col ? ' active' : '';
      return `<th class="stat-cell sortable${activeClass}" data-sort="${col}">${label} <span class="sort-indicator ${ind}"></span></th>`;
    })
    .join("");

  const idInd = currentSortColumn === 'id' ? (currentSortDirection === 'asc' ? 'asc' : currentSortDirection === 'desc' ? 'desc' : '') : '';
  const nameInd = currentSortColumn === 'name' ? (currentSortDirection === 'asc' ? 'asc' : currentSortDirection === 'desc' ? 'desc' : '') : '';
  const totalInd = currentSortColumn === 'total' ? (currentSortDirection === 'asc' ? 'asc' : currentSortDirection === 'desc' ? 'desc' : '') : '';

  grid.innerHTML = `<table class="poke-table">
    <thead><tr>
      <th class="sortable${currentSortColumn === 'id' ? ' active' : ''}" data-sort="id"># <span class="sort-indicator ${idInd}"></span></th>
      <th></th>
      <th class="sortable${currentSortColumn === 'name' ? ' active' : ''}" data-sort="name">Nombre <span class="sort-indicator ${nameInd}"></span></th>
      <th>Tipo</th>
      ${statHeaders}
      <th class="stat-cell sortable${currentSortColumn === 'total' ? ' active' : ''}" data-sort="total">Total <span class="sort-indicator ${totalInd}"></span></th>
    </tr></thead>
    <tbody>${rows}</tbody>
  </table>`;

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
            await renderCurrentView(filteredList.slice(offset, offset + LIMIT));
          } else {
            await loadList();
          }
          updatePagination();
          return;
        }
      } else {
        currentSortColumn = col;
        currentSortDirection = 'asc';
      }

      sortingLoading = true;
      try {
        const allPokemon = await ensureAllPokemonLoaded();
        sortedFullList = sortPokemonData(allPokemon, currentSortColumn, currentSortDirection);
        offset = 0;
        await renderTable(getCurrentPageItems());
        updatePagination();
      } finally {
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

function updatePagination(): void {
  const total = sortedFullList ? sortedFullList.length : totalCount;
  const page = Math.floor(offset / LIMIT) + 1;
  const pages = Math.ceil(total / LIMIT) || 1;
  pageInfo.textContent = `Pág. ${page} / ${pages}`;
  prevBtn.disabled = offset === 0;
  nextBtn.disabled = offset + LIMIT >= total;
}

// -- Paginación --------------------------------------------------------------

function resetSorting(): void {
  currentSortColumn = null;
  currentSortDirection = null;
  sortedFullList = null;
}

async function prevPage(): Promise<void> {
  if (offset <= 0) return;
  offset -= LIMIT;
  if (sortedFullList) {
    await renderCurrentView(getCurrentPageItems());
    updatePagination();
  } else if (hasFilter()) {
    await renderCurrentView(filteredList.slice(offset, offset + LIMIT));
    updatePagination();
  } else {
    loadList();
  }
}

async function nextPage(): Promise<void> {
  const max = sortedFullList ? sortedFullList.length : hasFilter() ? filteredList.length : totalCount;
  if (offset + LIMIT >= max) return;
  offset += LIMIT;
  if (sortedFullList) {
    await renderCurrentView(getCurrentPageItems());
    updatePagination();
  } else if (hasFilter()) {
    await renderCurrentView(filteredList.slice(offset, offset + LIMIT));
    updatePagination();
  } else {
    loadList();
  }
}

// -- Detalle -----------------------------------------------------------------

async function showDetail(name: string): Promise<void> {
  await showView(detailView, listView);
  detailEl.innerHTML = '<p class="loading">Cargando...</p>';
  try {
    const p = await GetPokemon(name);
    renderDetail(p);
  } catch (err: unknown) {
    detailEl.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}

async function renderDetail(p: Pokemon): Promise<void> {
  const types = (p.Types || [])
    .map((t) => `<span class="type-badge type-${t.Name}">${t.Name}</span>`)
    .join("");

  const sprites = `
    <div class="sprites">
      ${p.Sprites.FrontDefault ? `<div><img src="${p.Sprites.FrontDefault}" alt="default"/><span>Normal</span></div>` : ""}
      ${p.Sprites.FrontShiny ? `<div><img src="${p.Sprites.FrontShiny}" alt="shiny"/><span>Shiny</span></div>` : ""}
    </div>`;

  detailEl.innerHTML = `
    <h2>#${p.ID} ${p.Name}</h2>
    ${sprites}
    <div class="types">${types}</div>
    <p class="meta">Altura: ${p.Height / 10} m &nbsp;&middot;&nbsp; Peso: ${p.Weight / 10} kg</p>
    <div id="stats-chart" style="width:100%;height:300px;"></div>
    ${renderEVCalculatorForm(p)}`;

  const chartContainer = document.getElementById("stats-chart") as HTMLDivElement;
  const { renderStatsChart } = await import("../charts/stats-chart");
  renderStatsChart(chartContainer, p.Stats || []);
  await initEVCalculator(p);
}

// -- Búsqueda ----------------------------------------------------------------

async function search(): Promise<void> {
  const query = searchInput.value.trim().toLowerCase();
  if (!query) {
    offset = 0;
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
        updateDropdownTrigger("dropdown-gen", "Generación", filter.generations.length);
        applyFilters();
      });
      filterGenContainer.appendChild(chip);
    });

    types.Results.forEach((t) => {
      if (t.Name === "shadow" || t.Name === "unknown") return;
      const chip = document.createElement("button");
      chip.className = "filter-chip";
      chip.dataset.value = t.Name;
      chip.textContent = t.Name.charAt(0).toUpperCase() + t.Name.slice(1);
      chip.addEventListener("click", () => {
        chip.classList.toggle("active");
        if (chip.classList.contains("active")) {
          filter.types.push(t.Name);
        } else {
          filter.types = filter.types.filter((v) => v !== t.Name);
        }
        updateDropdownTrigger("dropdown-type", "Tipo", filter.types.length);
        applyFilters();
      });
      filterTypeContainer.appendChild(chip);
    });
  } catch (err) {
    console.error("Error cargando filtros:", err);
  }
}

function formatGenName(name: string): string {
  const match = name.match(/generation-(.+)/);
  if (!match) return name;
  return "Gen " + match[1].toUpperCase();
}

function applyFilters(): void {
  offset = 0;
  filteredList = [];

  if (hasFilter()) {
    loadFiltered();
  } else {
    resetSorting();
    loadList();
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
  updateDropdownTrigger("dropdown-gen", "Generación", 0);
  updateDropdownTrigger("dropdown-type", "Tipo", 0);
}

// -- Helpers -----------------------------------------------------------------

function idFromURL(url: string): string {
  const parts = url.replace(/\/$/, "").split("/");
  return parts[parts.length - 1];
}

function spriteURL(id: string): string {
  const numId = parseInt(id, 10);
  if (isNaN(numId)) return "";
  return `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${numId}.png`;
}

// -- Current page items helper -----------------------------------------------

let lastRenderedItems: PokemonListItem[] = [];

function getCurrentPageItems(): PokemonListItem[] {
  if (sortedFullList) {
    const pageItems = sortedFullList.slice(offset, offset + LIMIT);
    return pageItems.map((p) => ({ Name: p.Name, URL: "" }));
  }
  return lastRenderedItems;
}

// -- Init --------------------------------------------------------------------

export function initPokedex(): void {
  grid = document.getElementById("pokemon-grid") as HTMLDivElement;
  listView = document.getElementById("list-view") as HTMLElement;
  detailView = document.getElementById("detail-view") as HTMLElement;
  detailEl = document.getElementById("pokemon-detail") as HTMLDivElement;
  prevBtn = document.getElementById("prev-btn") as HTMLButtonElement;
  nextBtn = document.getElementById("next-btn") as HTMLButtonElement;
  pageInfo = document.getElementById("page-info") as HTMLSpanElement;
  searchInput = document.getElementById("search-input") as HTMLInputElement;
  searchBtn = document.getElementById("search-btn") as HTMLButtonElement;
  backBtn = document.getElementById("back-btn") as HTMLButtonElement;
  filterGenContainer = document.getElementById("filter-gen") as HTMLDivElement;
  filterTypeContainer = document.getElementById("filter-type") as HTMLDivElement;
  filterLegendaryBtn = document.getElementById("filter-legendary") as HTMLButtonElement;
  filterMythicalBtn = document.getElementById("filter-mythical") as HTMLButtonElement;
  filterResetBtn = document.getElementById("filter-reset") as HTMLButtonElement;
  viewToggleBtn = document.getElementById("view-toggle-btn") as HTMLButtonElement;

  prevBtn.addEventListener("click", prevPage);
  nextBtn.addEventListener("click", nextPage);

  viewToggleBtn.addEventListener("click", async () => {
    const oldMode = viewMode;
    viewMode = viewMode === "grid" ? "table" : "grid";
    resetSorting();
    viewToggleBtn.textContent = viewMode === "grid" ? "⊞ Tabla" : "⊟ Tarjetas";

    const items = getCurrentPageItems();
    if (oldMode === "grid" && viewMode === "table") {
      await morphToTable(grid, async () => {
        await renderTable(items);
      });
    } else {
      await morphToGrid(grid, () => {
        renderGrid(items);
      });
    }
  });

  backBtn.addEventListener("click", async () => {
    const { disposeChart } = await import("../charts/stats-chart");
    disposeChart();
    await showView(listView, detailView);
  });

  searchBtn.addEventListener("click", search);
  searchInput.addEventListener("keydown", (e: KeyboardEvent) => {
    if (e.key === "Enter") search();
  });

  filterLegendaryBtn.addEventListener("click", () => {
    filter.legendary = !filter.legendary;
    filterLegendaryBtn.classList.toggle("active", filter.legendary);
    offset = 0;
    filteredList = [];
    if (hasFilter()) loadFiltered();
    else loadList();
  });

  filterMythicalBtn.addEventListener("click", () => {
    filter.mythical = !filter.mythical;
    filterMythicalBtn.classList.toggle("active", filter.mythical);
    offset = 0;
    filteredList = [];
    if (hasFilter()) loadFiltered();
    else loadList();
  });

  filterResetBtn.addEventListener("click", () => {
    offset = 0;
    filter = { generations: [], types: [], legendary: false, mythical: false };
    filteredList = [];
    resetSorting();
    resetFilterUI();
    loadList();
  });

  populateFilters();
  loadList();
}
