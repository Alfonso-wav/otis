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
import { showView, staggerCards } from "../animations/transitions";
import { renderEVCalculatorForm, initEVCalculator } from "../ev-calculator";

const LIMIT = 20;

interface FilterState {
  generation: string | null;
  type: string | null;
  legendary: boolean;
  mythical: boolean;
}

let filter: FilterState = { generation: null, type: null, legendary: false, mythical: false };
let filteredList: PokemonListItem[] = [];
let offset = 0;
let totalCount = 0;

// Cache para evitar llamadas repetidas a GetPokemonSpecies
const legendaryCache = new Map<string, { isLegendary: boolean; isMythical: boolean }>();

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
let filterGenSelect: HTMLSelectElement;
let filterTypeSelect: HTMLSelectElement;
let filterLegendaryBtn: HTMLButtonElement;
let filterMythicalBtn: HTMLButtonElement;
let filterResetBtn: HTMLButtonElement;

function hasFilter(): boolean {
  return (
    filter.generation !== null ||
    filter.type !== null ||
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
    renderGrid(data.Results);
    updatePagination();
  } catch (err: unknown) {
    grid.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}

// -- Lista con filtro --------------------------------------------------------

async function loadFiltered(): Promise<void> {
  grid.innerHTML = '<p class="loading">Aplicando filtros...</p>';
  try {
    let base: PokemonListItem[] = [];

    if (filter.generation && filter.type) {
      // Ambos filtros: intersección
      const [gen, typeData] = await Promise.all([
        GetGeneration(filter.generation),
        GetType(filter.type),
      ]);
      const genNames = new Set(gen.PokemonSpecies.map((p) => p.Name));
      base = typeData.Pokemon.filter((p) => genNames.has(p.Name)).map((p) => ({
        Name: p.Name,
        URL: p.URL,
      }));
    } else if (filter.generation) {
      const gen = await GetGeneration(filter.generation);
      base = gen.PokemonSpecies.map((p) => ({ Name: p.Name, URL: p.URL }));
    } else if (filter.type) {
      const typeData = await GetType(filter.type);
      base = typeData.Pokemon.map((p) => ({ Name: p.Name, URL: p.URL }));
    } else if (filter.legendary || filter.mythical) {
      grid.innerHTML =
        '<p class="loading">Selecciona una generación o tipo para filtrar por legendario / mítico.</p>';
      return;
    }

    if ((filter.legendary || filter.mythical) && base.length > 0) {
      base = await filterByLegendary(base);
    }

    filteredList = base;
    totalCount = base.length;
    offset = 0;

    const page = filteredList.slice(0, LIMIT);
    renderGrid(page);
    updatePagination();
  } catch (err: unknown) {
    grid.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
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

function updatePagination(): void {
  const page = Math.floor(offset / LIMIT) + 1;
  const pages = Math.ceil(totalCount / LIMIT) || 1;
  pageInfo.textContent = `Pág. ${page} / ${pages}`;
  prevBtn.disabled = offset === 0;
  nextBtn.disabled = hasFilter()
    ? offset + LIMIT >= filteredList.length
    : offset + LIMIT >= totalCount;
}

// -- Paginación --------------------------------------------------------------

function prevPage(): void {
  if (offset <= 0) return;
  offset -= LIMIT;
  if (hasFilter()) {
    renderGrid(filteredList.slice(offset, offset + LIMIT));
    updatePagination();
  } else {
    loadList();
  }
}

function nextPage(): void {
  const max = hasFilter() ? filteredList.length : totalCount;
  if (offset + LIMIT >= max) return;
  offset += LIMIT;
  if (hasFilter()) {
    renderGrid(filteredList.slice(offset, offset + LIMIT));
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
    filter = { generation: null, type: null, legendary: false, mythical: false };
    filteredList = [];
    resetFilterUI();
    loadList();
    return;
  }
  await showDetail(query);
}

// -- Filtros -----------------------------------------------------------------

async function populateFilters(): Promise<void> {
  try {
    const [gens, types] = await Promise.all([ListGenerations(), ListTypes()]);

    gens.forEach((g) => {
      const opt = document.createElement("option");
      opt.value = g.Name;
      opt.textContent = formatGenName(g.Name);
      filterGenSelect.appendChild(opt);
    });

    types.Results.forEach((t) => {
      if (t.Name === "shadow" || t.Name === "unknown") return;
      const opt = document.createElement("option");
      opt.value = t.Name;
      opt.textContent = t.Name.charAt(0).toUpperCase() + t.Name.slice(1);
      filterTypeSelect.appendChild(opt);
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
  filter.generation = filterGenSelect.value || null;
  filter.type = filterTypeSelect.value || null;

  if (hasFilter()) {
    loadFiltered();
  } else {
    loadList();
  }
}

function resetFilterUI(): void {
  filterGenSelect.value = "";
  filterTypeSelect.value = "";
  filter.legendary = false;
  filter.mythical = false;
  filterLegendaryBtn.classList.remove("active");
  filterMythicalBtn.classList.remove("active");
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
  filterGenSelect = document.getElementById("filter-gen") as HTMLSelectElement;
  filterTypeSelect = document.getElementById("filter-type") as HTMLSelectElement;
  filterLegendaryBtn = document.getElementById("filter-legendary") as HTMLButtonElement;
  filterMythicalBtn = document.getElementById("filter-mythical") as HTMLButtonElement;
  filterResetBtn = document.getElementById("filter-reset") as HTMLButtonElement;

  prevBtn.addEventListener("click", prevPage);
  nextBtn.addEventListener("click", nextPage);

  backBtn.addEventListener("click", async () => {
    const { disposeChart } = await import("../charts/stats-chart");
    disposeChart();
    await showView(listView, detailView);
  });

  searchBtn.addEventListener("click", search);
  searchInput.addEventListener("keydown", (e: KeyboardEvent) => {
    if (e.key === "Enter") search();
  });

  filterGenSelect.addEventListener("change", applyFilters);
  filterTypeSelect.addEventListener("change", applyFilters);

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
    filter = { generation: null, type: null, legendary: false, mythical: false };
    filteredList = [];
    resetFilterUI();
    loadList();
  });

  populateFilters();
  loadList();
}
