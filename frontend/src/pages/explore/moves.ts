import gsap from "gsap";
import { GetAllMoves } from "../../../wailsjs/go/app/App";
import type { core } from "../../../wailsjs/go/models";

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
    moves = moves.filter((m) => m.Name.includes(state.searchQuery));
  }
  if (state.selectedCategory !== "all") {
    moves = moves.filter((m) => m.Category === state.selectedCategory);
  }
  if (state.selectedType !== "all") {
    moves = moves.filter((m) => m.Type === state.selectedType);
  }

  if (state.sortColumn && state.sortDirection) {
    const mult = state.sortDirection === "asc" ? 1 : -1;
    moves = [...moves].sort((a, b) => {
      switch (state.sortColumn) {
        case "name": return mult * a.Name.localeCompare(b.Name);
        case "type": return mult * a.Type.localeCompare(b.Type);
        case "category": return mult * a.Category.localeCompare(b.Category);
        case "power": return mult * (a.Power - b.Power);
        case "accuracy": return mult * (a.Accuracy - b.Accuracy);
        case "pp": return mult * (a.PP - b.PP);
        case "priority": return mult * (a.Priority - b.Priority);
        default: return 0;
      }
    });
  }

  return moves;
}

function renderTable(container: HTMLElement): void {
  const tbody = container.querySelector<HTMLElement>("#moves-tbody");
  if (!tbody) return;

  const moves = filteredMoves();
  const countEl = container.querySelector<HTMLElement>("#moves-count");
  if (countEl) countEl.textContent = `${moves.length} movimientos`;

  if (moves.length === 0) {
    tbody.innerHTML = '<tr><td colspan="7" class="loading">No se encontraron movimientos.</td></tr>';
    return;
  }

  tbody.innerHTML = moves.map((m) => `<tr>
    <td class="move-name-cell">${m.Name.replace(/-/g, " ")}</td>
    <td><span class="type-badge" style="background:${typeColor(m.Type)}">${m.Type}</span></td>
    <td class="move-cat-cell">${categoryIcon(m.Category)}</td>
    <td class="num-cell">${m.Power || "—"}</td>
    <td class="num-cell">${m.Accuracy ? m.Accuracy + "%" : "—"}</td>
    <td class="num-cell">${m.PP}</td>
    <td class="num-cell">${m.Priority}</td>
  </tr>`).join("");

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
  return `<option value="all">Todos los tipos</option>` +
    types.map((t) => `<option value="${t}">${t}</option>`).join("");
}

export async function initMoves(container: HTMLElement): Promise<void> {
  container.innerHTML = `
    <div class="section-header"><h2>Movimientos</h2></div>
    <div class="moves-controls">
      <input type="text" id="moves-search" class="explore-input" placeholder="Buscar movimiento..." />
      <div class="moves-filters">
        <button class="filter-btn active" data-cat="all">Todos</button>
        <button class="filter-btn" data-cat="physical"><img src="https://img.pokemondb.net/images/icons/move-physical.png" class="move-cat-icon" alt="Physical"> Fisico</button>
        <button class="filter-btn" data-cat="special"><img src="https://img.pokemondb.net/images/icons/move-special.png" class="move-cat-icon" alt="Special"> Especial</button>
        <button class="filter-btn" data-cat="status"><img src="https://img.pokemondb.net/images/icons/move-status.png" class="move-cat-icon" alt="Status"> Estado</button>
        <select id="moves-type-filter" class="explore-input moves-type-select">
          <option value="all">Todos los tipos</option>
        </select>
      </div>
    </div>
    <span id="moves-count" class="moves-count"></span>
    <div class="moves-table-wrap">
      <table class="poke-table moves-table">
        <thead>
          <tr>
            <th class="sortable" data-col="name">Nombre <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="type">Tipo <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="category">Cat. <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="power">Poder <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="accuracy">Prec. <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="pp">PP <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="priority">Prio. <span class="sort-indicator"></span></th>
          </tr>
        </thead>
        <tbody id="moves-tbody">
          <tr><td colspan="7" class="loading">Cargando movimientos...</td></tr>
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
    th.addEventListener("click", () => {
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
      renderTable(container);
    });
  });

  // Load all moves
  if (state.allMoves.length === 0) {
    state.loading = true;
    try {
      state.allMoves = await GetAllMoves();
    } catch (err) {
      const tbody = container.querySelector<HTMLElement>("#moves-tbody");
      if (tbody) tbody.innerHTML = '<tr><td colspan="7" class="loading">Error al cargar movimientos.</td></tr>';
      return;
    } finally {
      state.loading = false;
    }
  }

  // Populate type filter
  typeSelect.innerHTML = buildTypeFilterOptions();

  renderTable(container);
}
