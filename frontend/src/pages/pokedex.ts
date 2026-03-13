import { ListPokemon, GetPokemon } from "../../wailsjs/go/app/App";
import type { Pokemon, PokemonListItem } from "../types";
import { showView, staggerCards } from "../animations/transitions";
import { renderEVCalculatorForm, initEVCalculator } from "../ev-calculator";

const LIMIT = 20;
let offset = 0;
let totalCount = 0;

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

// -- Lista -------------------------------------------------------------------

async function loadList(): Promise<void> {
  grid.innerHTML = '<p class="loading">Cargando...</p>';
  try {
    const data = await ListPokemon(offset, LIMIT);
    totalCount = data.Count;
    renderGrid(data.Results);
    updatePagination();
  } catch (err: unknown) {
    grid.innerHTML = `<p class="loading" style="color:#e53e3e">${String(err)}</p>`;
  }
}

function renderGrid(items: PokemonListItem[]): void {
  if (!items || items.length === 0) {
    grid.innerHTML = '<p class="loading">No se encontraron Pokemon.</p>';
    return;
  }
  grid.innerHTML = items
    .map((item) => {
      const id = idFromURL(item.URL);
      const sprite = spriteURL(id);
      return `<div class="card" data-name="${item.Name}">
      <img src="${sprite}" alt="${item.Name}" loading="lazy" />
      <div class="poke-name">${item.Name}</div>
    </div>`;
    })
    .join("");

  grid.querySelectorAll<HTMLDivElement>(".card").forEach((card) => {
    card.addEventListener("click", () => {
      const name = card.dataset.name;
      if (name) showDetail(name);
    });
  });

  staggerCards(grid);
}

function updatePagination(): void {
  const page = Math.floor(offset / LIMIT) + 1;
  const pages = Math.ceil(totalCount / LIMIT);
  pageInfo.textContent = `Pagina ${page} / ${pages}`;
  prevBtn.disabled = offset === 0;
  nextBtn.disabled = offset + LIMIT >= totalCount;
}

// -- Detalle -----------------------------------------------------------------

async function showDetail(name: string): Promise<void> {
  await showView(detailView, listView);
  detailEl.innerHTML = '<p class="loading">Cargando...</p>';

  try {
    const p = await GetPokemon(name);
    renderDetail(p);
  } catch (err: unknown) {
    detailEl.innerHTML = `<p class="loading" style="color:#e53e3e">${String(err)}</p>`;
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

// -- Busqueda ----------------------------------------------------------------

async function search(): Promise<void> {
  const query = searchInput.value.trim().toLowerCase();
  if (!query) {
    loadList();
    return;
  }
  await showDetail(query);
}

// -- Helpers -----------------------------------------------------------------

function idFromURL(url: string): string {
  const parts = url.replace(/\/$/, "").split("/");
  return parts[parts.length - 1];
}

function spriteURL(id: string): string {
  return `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${id}.png`;
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

  prevBtn.addEventListener("click", () => {
    offset -= LIMIT;
    loadList();
  });
  nextBtn.addEventListener("click", () => {
    offset += LIMIT;
    loadList();
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

  loadList();
}
