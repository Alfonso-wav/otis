import gsap from "gsap";
import { GetAllAbilities } from "../../api";
import type { core } from "../../../wailsjs/go/models";
import { openAbilityPokemonModal } from "../../components/ability-pokemon-modal";
import { initColumnToggle, reapplyColumnVisibility, type ColumnConfig } from "../../components/column-toggle";
import { SortCache } from "../../utils/sort-cache";
import { showSortingOverlay, hideSortingOverlay } from "../../components/sorting-overlay";
import { t } from "../../i18n";

type SortColumn = "name" | "description" | "pokemon" | null;
type SortDirection = "asc" | "desc" | null;

interface AbilityState {
  allAbilities: core.Ability[];
  searchQuery: string;
  sortColumn: SortColumn;
  sortDirection: SortDirection;
  loading: boolean;
}

function abilitiesTableColumns(): ColumnConfig[] {
  return [
    { key: "name", label: t("abilities.columns.name"), fixed: true },
    { key: "description", label: t("abilities.columns.description") },
    { key: "pokemon", label: t("abilities.columns.pokemon") },
  ];
}

const state: AbilityState = {
  allAbilities: [],
  searchQuery: "",
  sortColumn: null,
  sortDirection: null,
  loading: false,
};

const abilitiesSortCache = new SortCache<core.Ability>([
  { key: "name", compare: (a, b) => a.Name.localeCompare(b.Name) },
  { key: "description", compare: (a, b) => a.Description.localeCompare(b.Description) },
  { key: "pokemon", compare: (a, b) => (a.Pokemon ?? []).length - (b.Pokemon ?? []).length },
]);

function filteredAbilities(): core.Ability[] {
  let abilities = state.allAbilities;

  if (state.searchQuery !== "") {
    abilities = abilities.filter((a) => a.Name.includes(state.searchQuery));
  }

  if (state.sortColumn && state.sortDirection) {
    abilitiesSortCache.setData(abilities);
    abilities = abilitiesSortCache.get(state.sortColumn, state.sortDirection);
  }

  return abilities;
}

function renderTable(container: HTMLElement): void {
  const tbody = container.querySelector<HTMLElement>("#abilities-tbody");
  if (!tbody) return;

  const abilities = filteredAbilities();
  const countEl = container.querySelector<HTMLElement>("#abilities-count");
  if (countEl) countEl.textContent = t("abilities.count", { count: abilities.length });

  if (abilities.length === 0) {
    tbody.innerHTML = `<tr><td colspan="3" class="loading">${t("abilities.noResults")}</td></tr>`;
    return;
  }

  tbody.innerHTML = abilities.map((a) => {
    const desc = a.Description ? a.Description.replace(/\n/g, " ") : "—";
    const pokemonList = a.Pokemon ?? [];
    const count = pokemonList.length;
    const countHtml = count > 0
      ? `<span class="ability-pokemon-count" data-ability="${a.Name}">${count}</span>`
      : `${count}`;
    return `<tr>
    <td class="ability-name-cell" data-col="name">${a.Name.replace(/-/g, " ")}</td>
    <td class="ability-desc-cell" data-col="description">${desc}</td>
    <td class="num-cell" data-col="pokemon">${countHtml}</td>
  </tr>`;
  }).join("");

  reapplyColumnVisibility("abilities");
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

export async function initAbilities(container: HTMLElement): Promise<void> {
  container.innerHTML = `
    <div class="section-header"><h2>${t("abilities.title")}</h2></div>
    <div class="abilities-controls">
      <input type="text" id="abilities-search" class="explore-input" placeholder="${t("abilities.searchPlaceholder")}" />
    </div>
    <span id="abilities-count" class="abilities-count"></span>
    <div class="abilities-table-wrap">
      <table class="poke-table abilities-table" data-table-id="abilities">
        <thead>
          <tr>
            <th class="sortable" data-col="name">${t("abilities.columns.name")} <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="description">${t("abilities.columns.description")} <span class="sort-indicator"></span></th>
            <th class="sortable" data-col="pokemon">${t("abilities.columns.pokemon")} <span class="sort-indicator"></span></th>
          </tr>
        </thead>
        <tbody id="abilities-tbody">
          <tr><td colspan="3" class="loading">${t("abilities.loading")}</td></tr>
        </tbody>
      </table>
    </div>`;

  // Search
  const searchInput = container.querySelector<HTMLInputElement>("#abilities-search")!;
  searchInput.addEventListener("input", () => {
    state.searchQuery = searchInput.value.trim().toLowerCase().replace(/\s+/g, "-");
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

  // Delegated click on pokemon count
  const tableWrap = container.querySelector<HTMLElement>(".abilities-table-wrap")!;
  tableWrap.addEventListener("click", (e) => {
    const target = (e.target as HTMLElement).closest<HTMLElement>(".ability-pokemon-count");
    if (!target) return;
    const abilityName = target.dataset.ability;
    if (!abilityName) return;
    const ability = state.allAbilities.find((a) => a.Name === abilityName);
    if (!ability || (ability.Pokemon ?? []).length === 0) return;
    openAbilityPokemonModal(abilityName, ability.Pokemon);
  });

  // Load all abilities
  if (state.allAbilities.length === 0) {
    state.loading = true;
    try {
      state.allAbilities = await GetAllAbilities();
    } catch {
      const tbody = container.querySelector<HTMLElement>("#abilities-tbody");
      if (tbody) tbody.innerHTML = `<tr><td colspan="3" class="loading">${t("abilities.error")}</td></tr>`;
      return;
    } finally {
      state.loading = false;
    }
  }

  renderTable(container);
  initColumnToggle("abilities", abilitiesTableColumns());
}
