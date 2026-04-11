import gsap from "gsap";
import { GetAllMoves } from "../../api";
import type { core } from "../../../wailsjs/go/models";
import { initColumnToggle, reapplyColumnVisibility, type ColumnConfig } from "../../components/column-toggle";
import { SortCache } from "../../utils/sort-cache";
import { showSortingOverlay, hideSortingOverlay, showLoadingOverlay } from "../../components/sorting-overlay";
import { t, typeName, getLocale } from "../../i18n";
import { openMovePokemonModal } from "../../components/move-pokemon-modal";

type Category = "all" | "physical" | "special" | "status";
type SortColumn = "name" | "type" | "category" | "power" | "accuracy" | "pp" | "priority" | null;
type SortDirection = "asc" | "desc" | null;

interface MoveState {
  allMoves: core.Move[];
  selectedCategory: Category;
  selectedType: string;
  searchQuery: string;
  sortColumn: SortColumn;
  sortDirection: SortDirection;
  loading: boolean;
}

const state: MoveState = {
  allMoves: [],
  selectedCategory: "all",
  selectedType: "all",
  searchQuery: "",
  sortColumn: null,
  sortDirection: null,
  loading: false,
};

const movesSortCache = new SortCache<core.Move>([
  { key: "name", compare: (a, b) => a.Name.localeCompare(b.Name) },
  { key: "type", compare: (a, b) => a.Type.localeCompare(b.Type) },
  { key: "category", compare: (a, b) => a.Category.localeCompare(b.Category) },
  { key: "power", compare: (a, b) => a.Power - b.Power },
  { key: "accuracy", compare: (a, b) => a.Accuracy - b.Accuracy },
  { key: "pp", compare: (a, b) => a.PP - b.PP },
  { key: "priority", compare: (a, b) => a.Priority - b.Priority },
]);

function movesTableColumns(): ColumnConfig[] {
  return [
    { key: "name", label: t("moves.columns.name"), fixed: true },
    { key: "type", label: t("moves.columns.type") },
    { key: "category", label: t("moves.columns.category") },
    { key: "power", label: t("moves.columns.power") },
    { key: "accuracy", label: t("moves.columns.accuracy") },
    { key: "pp", label: t("moves.columns.pp") },
    { key: "priority", label: t("moves.columns.priority") },
    { key: "pokemon", label: t("moves.columns.pokemon") },
  ];
}

const TYPE_COLORS: Record<string, string> = {
  normal: "#a0aec0", fire: "#f6ad55", water: "#63b3ed", grass: "#68d391",
  electric: "#f6e05e", psychic: "#f687b3", ice: "#76e4f7", fighting: "#c05621",
  poison: "#9f7aea", ground: "#d69e2e", flying: "#90cdf4", bug: "#a8e063",
  rock: "#b7791f", ghost: "#553c9a", dragon: "#7f9cf5", dark: "#4a5568",
  steel: "#718096", fairy: "#fbb6ce",
};

function typeColor(t: string): string {
  return TYPE_COLORS[t] ?? "#718096";
}

function categoryIcon(cat: string): string {
  const base = "https://img.pokemondb.net/images/icons";
  const map: Record<string, string> = {
    physical: `<img src="${base}/move-physical.png" class="move-cat-icon" alt="Physical" title="Physical">`,
    special:  `<img src="${base}/move-special.png"  class="move-cat-icon" alt="Special"  title="Special">`,
    status:   `<img src="${base}/move-status.png"   class="move-cat-icon" alt="Status"   title="Status">`,
  };
  return map[cat] ?? `<span class="move-cat-unknown">?</span>`;
}

function filteredMoves(): core.Move[] {
  let moves = state.allMoves;

  if (state.searchQuery !== "") {
    const q = state.searchQuery;
    moves = moves.filter((m) =>
      m.Name.includes(q) ||
      (m.NameEs && m.NameEs.toLowerCase().includes(q.replace(/-/g, " ")))
    );
  }
  if (state.selectedCategory !== "all") {
    moves = moves.filter((m) => m.Category === state.selectedCategory);
  }
  if (state.selectedType !== "all") {
    moves = moves.filter((m) => m.Type === state.selectedType);
  }

  if (state.sortColumn && state.sortDirection) {
    movesSortCache.setData(moves);
    moves = movesSortCache.get(state.sortColumn, state.sortDirection);
  }

  return moves;
}

function renderTable(container: HTMLElement): void {
  const tbody = container.querySelector<HTMLElement>("#moves-tbody");
  if (!tbody) return;

  const moves = filteredMoves();
  const countEl = container.querySelector<HTMLElement>("#moves-count");
  if (countEl) countEl.textContent = t("moves.count", { count: moves.length });

  if (moves.length === 0) {
    tbody.innerHTML = `<tr><td colspan="8" class="loading">${t("moves.noResults")}</td></tr>`;
    return;
  }

  const isEs = getLocale() === "es";
  tbody.innerHTML = moves.map((m) => {
    const displayName = isEs && m.NameEs ? m.NameEs : m.Name.replace(/-/g, " ");
    return `<tr>
    <td class="move-name-cell" data-col="name">${displayName}</td>
    <td data-col="type"><span class="type-badge type-badge--icon-only" style="background:${typeColor(m.Type)}" title="${typeName(m.Type)}"><img src="/assets/types/${m.Type}.svg" alt="${typeName(m.Type)}" class="type-icon"></span></td>
    <td class="move-cat-cell" data-col="category">${categoryIcon(m.Category)}</td>
    <td class="num-cell" data-col="power">${m.Power || "—"}</td>
    <td class="num-cell" data-col="accuracy">${m.Accuracy ? m.Accuracy + "%" : "—"}</td>
    <td class="num-cell" data-col="pp">${m.PP}</td>
    <td class="num-cell" data-col="priority">${m.Priority}</td>
    <td data-col="pokemon"><button class="pkmn-btn" data-move="${m.Name}">PKMN</button></td>
  </tr>`;
  }).join("");

  tbody.querySelectorAll<HTMLButtonElement>(".pkmn-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const moveName = btn.dataset.move;
      if (moveName) openMovePokemonModal(moveName);
    });
  });

  reapplyColumnVisibility("moves");
  gsap.fromTo(tbody, { opacity: 0 }, { opacity: 1, duration: 0.25, ease: "power2.out" });
}

function updateSortIndicators(container: HTMLElement): void {
  container.querySelectorAll<HTMLElement>("th.sortable").forEach((th) => {
    const col = th.dataset.col as SortColumn;
    const indicator = th.querySelector<HTMLElement>(".sort-indicator");
    if (!indicator) return;

    th.classList.toggle("active", col === state.sortColumn);
    indicator.className = "sort-indicator";
    if (col === state.sortColumn && state.sortDirection) {
      indicator.classList.add(state.sortDirection);
    }
  });
}

function buildTypeFilterOptions(): string {
  const types = [...new Set(state.allMoves.map((m) => m.Type))].sort();
  return `<option value="all">${t("moves.allTypes")}</option>` +
    types.map((tp) => `<option value="${tp}">${typeName(tp)}</option>`).join("");
}

export async function initMoves(container: HTMLElement): Promise<void> {
  container.innerHTML = `
    <div class="section-header"><h2>${t("moves.title")}</h2></div>
    <div class="moves-controls">
      <input type="text" id="moves-search" class="explore-input" placeholder="${t("moves.searchPlaceholder")}" />
      <div class="moves-filters">
        <button class="filter-btn active" data-cat="all">${t("moves.categories.all")}</button>
        <button class="filter-btn" data-cat="physical"><img src="https://img.pokemondb.net/images/icons/move-physical.png" class="move-cat-icon" alt="Physical"> ${t("moves.categories.physical")}</button>
        <button class="filter-btn" data-cat="special"><img src="https://img.pokemondb.net/images/icons/move-special.png" class="move-cat-icon" alt="Special"> ${t("moves.categories.special")}</button>
        <button class="filter-btn" data-cat="status"><img src="https://img.pokemondb.net/images/icons/move-status.png" class="move-cat-icon" alt="Status"> ${t("moves.categories.status")}</button>
        <select id="moves-type-filter" class="explore-input moves-type-select">
          <option value="all">${t("moves.allTypes")}</option>
        </select>
      </div>
    </div>
    <span id="moves-count" class="moves-count"></span>
    <div class="moves-table-wrap">
      <table class="poke-table moves-table" data-table-id="moves">
        <thead>
          <tr>
            <th class="sortable" data-col="name">${t("moves.columns.name")} <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="type">${t("moves.columns.type")} <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="category">${t("moves.columns.category")} <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="power">${t("moves.columns.power")} <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="accuracy">${t("moves.columns.accuracy")} <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="pp">${t("moves.columns.pp")} <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="priority">${t("moves.columns.priority")} <span class="sort-indicator"></span></th>
            <th data-col="pokemon">${t("moves.columns.pokemon")}</th>
          </tr>
        </thead>
        <tbody id="moves-tbody">
          <tr><td colspan="8" class="loading">${t("moves.loading")}</td></tr>
        </tbody>
      </table>
    </div>`;

  // Search
  const searchInput = container.querySelector<HTMLInputElement>("#moves-search")!;
  searchInput.addEventListener("input", () => {
    state.searchQuery = searchInput.value.trim().toLowerCase().replace(/\s+/g, "-");
    renderTable(container);
  });

  // Category filter buttons
  container.querySelectorAll<HTMLButtonElement>(".filter-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      container.querySelectorAll(".filter-btn").forEach((b) => b.classList.remove("active"));
      btn.classList.add("active");
      state.selectedCategory = (btn.dataset.cat ?? "all") as Category;
      renderTable(container);
    });
  });

  // Type filter select
  const typeSelect = container.querySelector<HTMLSelectElement>("#moves-type-filter")!;
  typeSelect.addEventListener("change", () => {
    state.selectedType = typeSelect.value;
    renderTable(container);
  });

  // Sortable headers
  container.querySelectorAll<HTMLElement>("th.sortable").forEach((th) => {
    th.addEventListener("click", async () => {
      const col = th.dataset.col as SortColumn;
      if (state.sortColumn === col) {
        if (state.sortDirection === "asc") {
          state.sortDirection = "desc";
        } else if (state.sortDirection === "desc") {
          state.sortColumn = null;
          state.sortDirection = null;
        }
      } else {
        state.sortColumn = col;
        state.sortDirection = "asc";
      }
      updateSortIndicators(container);
      showSortingOverlay();
      await new Promise((r) => requestAnimationFrame(r));
      renderTable(container);
      hideSortingOverlay();
    });
  });

  // Load all moves
  if (state.allMoves.length === 0) {
    state.loading = true;
    showLoadingOverlay(t("moves.loading"));
    try {
      state.allMoves = await GetAllMoves();
    } catch (err) {
      const tbody = container.querySelector<HTMLElement>("#moves-tbody");
      if (tbody) tbody.innerHTML = `<tr><td colspan="8" class="loading">${t("moves.error")}</td></tr>`;
      return;
    } finally {
      state.loading = false;
      hideSortingOverlay();
    }
  }

  // Populate type filter
  typeSelect.innerHTML = buildTypeFilterOptions();

  renderTable(container);
  initColumnToggle("moves", movesTableColumns());

  document.addEventListener("locale-changed", () => {
    if (state.allMoves.length === 0) return;
    initMoves(container);
  });
}
