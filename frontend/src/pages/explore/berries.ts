import gsap from "gsap";
import { ListBerries, GetBerry } from "../../api";
import type { core } from "../../../wailsjs/go/models";
import { SortCache } from "../../utils/sort-cache";
import { showSortingOverlay, hideSortingOverlay, createInlineHeracross } from "../../components/sorting-overlay";
import { t, typeName } from "../../i18n";

type ViewMode = "cards" | "table";
type SortColumn = "id" | "name" | "naturalGiftPower" | "size" | null;
type SortDirection = "asc" | "desc" | null;

interface BerryState {
  allBerries: core.Berry[];
  listNames: string[];
  searchQuery: string;
  viewMode: ViewMode;
  sortColumn: SortColumn;
  sortDirection: SortDirection;
  rowLimit: number;
  page: number;
  loadedCount: number;
  loading: boolean;
  tableLoaded: boolean;
}

const TYPE_COLORS: Record<string, string> = {
  normal: "#a0aec0", fire: "#f6ad55", water: "#63b3ed", grass: "#68d391",
  electric: "#f6e05e", psychic: "#f687b3", ice: "#76e4f7", fighting: "#c05621",
  poison: "#9f7aea", ground: "#d69e2e", flying: "#90cdf4", bug: "#a8e063",
  rock: "#b7791f", ghost: "#553c9a", dragon: "#7f9cf5", dark: "#4a5568",
  steel: "#718096", fairy: "#fbb6ce",
};

const state: BerryState = {
  allBerries: [],
  listNames: [],
  searchQuery: "",
  viewMode: "cards",
  sortColumn: null,
  sortDirection: null,
  rowLimit: 20,
  page: 1,
  loadedCount: 0,
  loading: false,
  tableLoaded: false,
};

const berriesSortCache = new SortCache<core.Berry>([
  { key: "id",               compare: (a, b) => a.ID - b.ID },
  { key: "name",             compare: (a, b) => a.Name.localeCompare(b.Name) },
  { key: "naturalGiftPower", compare: (a, b) => a.NaturalGiftPower - b.NaturalGiftPower },
  { key: "size",             compare: (a, b) => a.Size - b.Size },
]);

function filteredBerries(): core.Berry[] {
  let berries = state.allBerries;
  if (state.searchQuery !== "") {
    const q = state.searchQuery.toLowerCase();
    berries = berries.filter((b) => b.Name.toLowerCase().includes(q));
  }
  if (state.sortColumn && state.sortDirection) {
    berriesSortCache.setData(berries);
    berries = berriesSortCache.get(state.sortColumn, state.sortDirection);
  }
  return berries;
}

function typeColor(type: string): string {
  return TYPE_COLORS[type] ?? "#a0aec0";
}

function typeIcon(type: string): string {
  if (!type) return "";
  return `<img src="/assets/types/${type}.svg" class="type-icon-sm" alt="${typeName(type)}" title="${typeName(type)}">`;
}

function firmnessBadge(firmness: string): string {
  return `<span class="berry-firmness-badge">${firmness.replace(/-/g, " ")}</span>`;
}

function flavorsText(berry: core.Berry): string {
  if (!berry.Flavors || berry.Flavors.length === 0) return "—";
  return berry.Flavors
    .slice()
    .sort((a, b) => b.Potency - a.Potency)
    .slice(0, 2)
    .map((f) => `${f.Flavor} (${f.Potency})`)
    .join(", ");
}

// ─── Cards view ─────────────────────────────────────────────────────────────

const CARDS_PAGE_SIZE = 50;

function renderCards(container: HTMLElement): void {
  const grid = container.querySelector<HTMLElement>(".berries-grid");
  if (!grid) return;

  const visible = filteredBerries();
  grid.innerHTML = visible.slice(0, state.loadedCount).map((berry) => {
    const sprite = berry.ItemSprite
      ? `<img src="${berry.ItemSprite}" class="berry-card__sprite" alt="${berry.Name}" loading="lazy">`
      : `<div class="berry-card__sprite-placeholder"></div>`;

    const typeCell = berry.NaturalGiftType
      ? `${typeIcon(berry.NaturalGiftType)} <span class="type-badge" style="background:${typeColor(berry.NaturalGiftType)}">${typeName(berry.NaturalGiftType)}</span>`
      : "—";

    const effectLine = berry.Effect
      ? `<div class="berry-card__effect" title="${berry.Effect}">${berry.Effect}</div>`
      : "";

    return `
      <div class="berry-card">
        <div class="berry-card__img">${sprite}</div>
        <div class="berry-card__body">
          <div class="berry-card__name">${capitalize(berry.Name.replace(/-/g, " "))}</div>
          <div class="berry-card__type">${typeCell}</div>
          <div class="berry-card__meta">${firmnessBadge(berry.Firmness)}</div>
          ${effectLine}
        </div>
      </div>`;
  }).join("");
}

// ─── Table view ──────────────────────────────────────────────────────────────

function renderTableBody(container: HTMLElement): void {
  const tbody = container.querySelector<HTMLElement>("#berries-tbody");
  if (!tbody) return;

  const berries = filteredBerries();
  const total = berries.length;
  const pages = Math.max(1, Math.ceil(total / state.rowLimit));
  if (state.page > pages) state.page = pages;

  const start = (state.page - 1) * state.rowLimit;
  const slice = berries.slice(start, start + state.rowLimit);

  const countEl = container.querySelector<HTMLElement>("#berries-count");
  if (countEl) countEl.textContent = t("berries.count", { count: total });

  const pagEl = container.querySelector<HTMLElement>("#berries-page-info");
  if (pagEl) pagEl.textContent = t("pokedex.page", { page: state.page, pages });

  const prevBtn = container.querySelector<HTMLButtonElement>("#berries-prev");
  const nextBtn = container.querySelector<HTMLButtonElement>("#berries-next");
  if (prevBtn) prevBtn.disabled = state.page <= 1;
  if (nextBtn) nextBtn.disabled = state.page >= pages;

  if (slice.length === 0) {
    tbody.innerHTML = `<tr><td colspan="11" class="loading">${t("berries.noResults")}</td></tr>`;
    return;
  }

  tbody.innerHTML = slice.map((berry) => {
    const sprite = berry.ItemSprite
      ? `<img src="${berry.ItemSprite}" class="berry-tbl-sprite" alt="${berry.Name}" loading="lazy">`
      : "—";

    const typeCell = berry.NaturalGiftType
      ? `${typeIcon(berry.NaturalGiftType)} <span class="type-badge" style="background:${typeColor(berry.NaturalGiftType)};font-size:0.65rem">${typeName(berry.NaturalGiftType)}</span>`
      : "—";

    const effectText = berry.Effect ? berry.Effect.slice(0, 60) + (berry.Effect.length > 60 ? "…" : "") : "—";

    return `<tr>
      <td class="num-cell">${berry.ID}</td>
      <td class="berry-name-cell">${capitalize(berry.Name.replace(/-/g, " "))}</td>
      <td>${sprite}</td>
      <td>${firmnessBadge(berry.Firmness)}</td>
      <td>${typeCell}</td>
      <td class="num-cell">${berry.NaturalGiftPower}</td>
      <td class="num-cell">${berry.Size}</td>
      <td class="num-cell">${berry.Smoothness}</td>
      <td class="num-cell">${berry.MaxHarvest}</td>
      <td class="berry-flavors-cell">${flavorsText(berry)}</td>
      <td class="berry-effect-cell" title="${berry.Effect}">${effectText}</td>
    </tr>`;
  }).join("");
}

// ─── Sorting ─────────────────────────────────────────────────────────────────

function sortIndicator(col: string): string {
  if (state.sortColumn !== col) return `<span class="sort-indicator">↕</span>`;
  return state.sortDirection === "asc" ? `<span class="sort-indicator active">↑</span>` : `<span class="sort-indicator active">↓</span>`;
}

function renderTableHeaders(container: HTMLElement): void {
  const thead = container.querySelector<HTMLElement>("#berries-thead");
  if (!thead) return;
  const cols: [string, string, boolean][] = [
    ["id",               t("berries.columns.id"),               true],
    ["name",             t("berries.columns.name"),             true],
    ["-",                t("berries.columns.sprite"),           false],
    ["-",                t("berries.columns.firmness"),         false],
    ["-",                t("berries.columns.naturalGiftType"),  false],
    ["naturalGiftPower", t("berries.columns.naturalGiftPower"), true],
    ["size",             t("berries.columns.size"),             true],
    ["-",                t("berries.columns.smoothness"),       false],
    ["-",                t("berries.columns.maxHarvest"),       false],
    ["-",                t("berries.columns.flavors"),          false],
    ["-",                t("berries.columns.effect"),           false],
  ];
  thead.innerHTML = `<tr>${cols.map(([key, label, sortable]) =>
    sortable && key !== "-"
      ? `<th class="sortable" data-col="${key}">${label} ${sortIndicator(key)}</th>`
      : `<th>${label}</th>`
  ).join("")}</tr>`;

  thead.querySelectorAll<HTMLElement>("[data-col]").forEach((th) => {
    th.addEventListener("click", () => {
      const col = th.dataset.col as SortColumn;
      if (state.sortColumn === col) {
        state.sortDirection = state.sortDirection === "asc" ? "desc" : "asc";
      } else {
        state.sortColumn = col;
        state.sortDirection = "asc";
      }
      showSortingOverlay();
      setTimeout(() => {
        renderTableHeaders(container);
        renderTableBody(container);
        hideSortingOverlay();
      }, 0);
    });
  });
}

// ─── Layout ──────────────────────────────────────────────────────────────────

function capitalize(s: string): string {
  return s ? s.charAt(0).toUpperCase() + s.slice(1) : s;
}

function buildTableView(container: HTMLElement): void {
  const content = container.querySelector<HTMLElement>(".berries-content");
  if (!content) return;

  content.innerHTML = `
    <div class="berries-table-controls">
      <span id="berries-count" class="berries-count"></span>
      <div class="berries-pagination">
        <label>${t("pokedex.rowLimit")}:
          <select id="berries-row-limit" class="filter-select" style="margin-left:0.4rem">
            ${[10, 20, 50, 100].map((n) => `<option value="${n}"${n === state.rowLimit ? " selected" : ""}>${n}</option>`).join("")}
          </select>
        </label>
        <button id="berries-prev" class="filter-pill">${t("common.previous")}</button>
        <span id="berries-page-info"></span>
        <button id="berries-next" class="filter-pill">${t("common.next")}</button>
      </div>
    </div>
    <div class="berries-table-wrap">
      <table class="table table-hover berries-table">
        <thead id="berries-thead"></thead>
        <tbody id="berries-tbody"></tbody>
      </table>
    </div>`;

  content.querySelector<HTMLSelectElement>("#berries-row-limit")?.addEventListener("change", (e) => {
    state.rowLimit = parseInt((e.target as HTMLSelectElement).value, 10);
    state.page = 1;
    renderTableBody(container);
  });
  content.querySelector<HTMLButtonElement>("#berries-prev")?.addEventListener("click", () => {
    if (state.page > 1) { state.page--; renderTableBody(container); }
  });
  content.querySelector<HTMLButtonElement>("#berries-next")?.addEventListener("click", () => {
    const pages = Math.ceil(filteredBerries().length / state.rowLimit);
    if (state.page < pages) { state.page++; renderTableBody(container); }
  });

  renderTableHeaders(container);
  renderTableBody(container);
}

function buildCardsView(container: HTMLElement): void {
  const content = container.querySelector<HTMLElement>(".berries-content");
  if (!content) return;

  content.innerHTML = `<div class="berries-grid"></div><div class="berries-sentinel"></div>`;
  renderCards(container);

  const sentinel = content.querySelector<HTMLElement>(".berries-sentinel");
  if (!sentinel) return;

  const observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting) {
      const visible = filteredBerries();
      if (state.loadedCount < visible.length) {
        state.loadedCount = Math.min(state.loadedCount + CARDS_PAGE_SIZE, visible.length);
        renderCards(container);
      }
    }
  }, { threshold: 0.1 });
  observer.observe(sentinel);
}

function rebuildView(container: HTMLElement): void {
  if (state.viewMode === "cards") {
    buildCardsView(container);
  } else {
    buildTableView(container);
  }
}

// ─── Load data ───────────────────────────────────────────────────────────────

async function loadAllBerryDetails(container: HTMLElement): Promise<void> {
  if (state.tableLoaded) return;
  state.tableLoaded = true;

  const names = state.listNames;
  const batchSize = 10;
  const results: core.Berry[] = [];

  for (let i = 0; i < names.length; i += batchSize) {
    const batch = names.slice(i, i + batchSize);
    const fetched = await Promise.allSettled(batch.map((n) => GetBerry(n)));
    fetched.forEach((r) => {
      if (r.status === "fulfilled") results.push(r.value);
    });
  }

  results.sort((a, b) => a.ID - b.ID);
  state.allBerries = results;
  berriesSortCache.setData(results);
  state.loadedCount = Math.min(CARDS_PAGE_SIZE, results.length);

  if (state.viewMode === "cards") {
    buildCardsView(container);
  } else {
    buildTableView(container);
  }
}

// ─── Entry point ─────────────────────────────────────────────────────────────

export async function initBerries(panel: HTMLElement): Promise<void> {
  panel.innerHTML = `
    <div class="berries-wrap">
      <div class="berries-controls">
        <input type="search" class="explore-input" id="berries-search"
          placeholder="${t("berries.searchPlaceholder")}">
        <div class="berries-view-toggle">
          <button class="filter-pill${state.viewMode === "cards" ? " active" : ""}" data-view="cards">${t("berries.cardView")}</button>
          <button class="filter-pill${state.viewMode === "table" ? " active" : ""}" data-view="table">${t("berries.tableView")}</button>
        </div>
      </div>
      <div class="berries-content">
        <div id="berries-loading"></div>
      </div>
    </div>`;

  // Show Heracross loading animation
  const loadingEl = panel.querySelector<HTMLElement>("#berries-loading");
  if (loadingEl) createInlineHeracross(loadingEl, t("berries.loading"));

  // Toggle view
  panel.querySelectorAll<HTMLButtonElement>("[data-view]").forEach((btn) => {
    btn.addEventListener("click", () => {
      const v = btn.dataset.view as ViewMode;
      if (v === state.viewMode) return;
      state.viewMode = v;
      panel.querySelectorAll<HTMLButtonElement>("[data-view]").forEach((b) =>
        b.classList.toggle("active", b.dataset.view === v)
      );
      if (state.allBerries.length > 0) {
        rebuildView(panel);
      } else if (!state.loading) {
        loadAllBerryDetails(panel);
      }
    });
  });

  // Search
  panel.querySelector<HTMLInputElement>("#berries-search")?.addEventListener("input", (e) => {
    state.searchQuery = (e.target as HTMLInputElement).value.toLowerCase();
    state.page = 1;
    state.loadedCount = Math.min(CARDS_PAGE_SIZE, filteredBerries().length);
    if (state.viewMode === "cards") {
      renderCards(panel);
    } else {
      renderTableHeaders(panel);
      renderTableBody(panel);
    }
  });

  // Locale updates
  panel.addEventListener("locale-changed", () => {
    panel.querySelector<HTMLButtonElement>("[data-view='cards']")!.textContent = t("berries.cardView");
    panel.querySelector<HTMLButtonElement>("[data-view='table']")!.textContent = t("berries.tableView");
    const searchInput = panel.querySelector<HTMLInputElement>("#berries-search");
    if (searchInput) searchInput.placeholder = t("berries.searchPlaceholder");
    if (state.viewMode === "table" && state.allBerries.length > 0) {
      renderTableHeaders(panel);
      renderTableBody(panel);
    }
  });

  // Load berry list, then details
  state.loading = true;
  try {
    const list = await ListBerries();
    state.listNames = list.Results.map((r: { Name: string }) => r.Name);
  } catch {
    panel.querySelector<HTMLElement>(".berries-content")!.innerHTML =
      `<p class="loading">${t("berries.error")}</p>`;
    state.loading = false;
    return;
  }
  state.loading = false;
  await loadAllBerryDetails(panel);

  gsap.fromTo(
    panel.querySelector(".berries-content"),
    { opacity: 0, y: 8 },
    { opacity: 1, y: 0, duration: 0.25 }
  );
}
