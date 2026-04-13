import { GetMove, GetPokemon } from "../api";
import { t, typeName } from "../i18n";
import { Pokemon } from "../types";
import { createInlineDiglett, removeInlineDiglett } from "./sorting-overlay";
import { getActiveId } from "../router";
import { setDetailOrigin } from "../pages/pokedex";

const MAX_POKEMON = 50;

let overlayEl: HTMLDivElement | null = null;
let isLoading = false;

function handleKeydown(e: KeyboardEvent): void {
  if (e.key === "Escape") closeMovePokemonModal();
}

function navigateToPokemon(name: string): void {
  const origin = getActiveId();
  closeMovePokemonModal();
  (document.querySelector('[data-tab="pokedex"]') as HTMLElement)?.click();
  setDetailOrigin(origin);
  const input = document.getElementById("search-input") as HTMLInputElement;
  const btn = document.getElementById("search-btn") as HTMLButtonElement;
  if (input && btn) {
    input.value = name;
    btn.click();
  }
}

function buildTableHTML(pokemonData: Pokemon[], totalCount: number, shownCount: number): string {
  const statColKeys = ["hp", "atk", "def", "spa", "spd", "vel"];

  const rows = pokemonData
    .map((p) => {
      const sprite = p.Sprites.FrontDefault
        ? `<img class="poke-table__sprite" src="${p.Sprites.FrontDefault}" alt="${p.Name}" loading="lazy" />`
        : "";
      const types = (p.Types || [])
        .map(
          (tp) =>
            `<span class="type-badge type-badge--icon-only type-${tp.Name}" title="${typeName(tp.Name)}"><img src="/assets/types/${tp.Name}.svg" alt="${typeName(tp.Name)}" class="type-icon"></span>`,
        )
        .join(" ");
      const stats = (p.Stats || []).map((s) => s.BaseStat);
      const total = stats.reduce((a, b) => a + b, 0);
      const statCells = (p.Stats || [])
        .map((s, i) => `<td class="stat-cell" data-col="${statColKeys[i]}">${s.BaseStat}</td>`)
        .join("");

      return `<tr class="poke-table__row" data-name="${p.Name}">
        <td data-col="sprite">${sprite}</td>
        <td class="poke-table__name" data-col="name">${p.Name}</td>
        <td data-col="types">${types}</td>
        ${statCells}
        <td class="stat-cell stat-total" data-col="total">${total}</td>
      </tr>`;
    })
    .join("");

  const limitNote =
    totalCount > shownCount
      ? `<p class="move-modal-limit-note">${t("modals.moveShowingOf", { shown: shownCount, total: totalCount })}</p>`
      : "";

  return `${limitNote}<div class="move-modal-table-wrap"><table class="poke-table" data-table-id="move-pokemon-stats">
    <thead><tr>
      <th data-col="sprite"></th>
      <th data-col="name">${t("pokedex.columns.name")}</th>
      <th data-col="types">${t("pokedex.columns.types")}</th>
      <th class="stat-cell" data-col="hp">${t("pokedex.columns.hp")}</th>
      <th class="stat-cell" data-col="atk">${t("pokedex.columns.atk")}</th>
      <th class="stat-cell" data-col="def">${t("pokedex.columns.def")}</th>
      <th class="stat-cell" data-col="spa">${t("pokedex.columns.spa")}</th>
      <th class="stat-cell" data-col="spd">${t("pokedex.columns.spd")}</th>
      <th class="stat-cell" data-col="vel">${t("pokedex.columns.vel")}</th>
      <th class="stat-cell" data-col="total">${t("pokedex.columns.total")}</th>
    </tr></thead>
    <tbody>${rows}</tbody>
  </table></div>`;
}

export async function openMovePokemonModal(moveName: string): Promise<void> {
  if (isLoading) return;
  closeMovePokemonModal();

  isLoading = true;
  try {
    const overlay = document.createElement("div");
    overlay.className = "type-modal-overlay";

    const displayName = moveName.replace(/-/g, " ");
    const capitalized =
      displayName.charAt(0).toUpperCase() + displayName.slice(1);

    overlay.innerHTML = `
      <div class="type-modal move-modal--wide">
        <div class="type-modal-header">
          <span>${capitalized}</span>
          <button class="type-modal-close" id="move-modal-close">&times;</button>
        </div>
        <div class="type-modal-body"></div>
      </div>`;

    document.body.appendChild(overlay);
    overlayEl = overlay;

    const bodyEl = overlay.querySelector<HTMLElement>(".type-modal-body")!;
    const diglettEl = createInlineDiglett(bodyEl, t("modals.moveLoading"));

    overlay
      .querySelector("#move-modal-close")!
      .addEventListener("click", closeMovePokemonModal);
    overlay.addEventListener("click", (e) => {
      if (e.target === overlay) closeMovePokemonModal();
    });
    document.addEventListener("keydown", handleKeydown);

    const move = await GetMove(moveName);
    const allNames = move.LearnedBy ?? [];
    const totalCount = allNames.length;
    const namesToFetch = allNames.slice(0, MAX_POKEMON);

    if (totalCount === 0) {
      removeInlineDiglett(diglettEl);
      const header = overlay.querySelector<HTMLElement>(
        ".type-modal-header span",
      );
      if (header) header.textContent = `${capitalized} (0)`;
      bodyEl.innerHTML = `<p class="type-modal-empty">${t("modals.moveEmpty", { name: capitalized })}</p>`;
      return;
    }

    const pokemonData: Pokemon[] = await Promise.all(
      namesToFetch.map((name) => GetPokemon(name)),
    );
    removeInlineDiglett(diglettEl);

    const header = overlay.querySelector<HTMLElement>(
      ".type-modal-header span",
    );
    if (header) header.textContent = `${capitalized} (${totalCount})`;

    bodyEl.innerHTML = buildTableHTML(
      pokemonData,
      totalCount,
      namesToFetch.length,
    );

    bodyEl
      .querySelectorAll<HTMLTableRowElement>(".poke-table__row")
      .forEach((row) => {
        row.addEventListener("click", () => {
          const name = row.dataset.name;
          if (name) navigateToPokemon(name);
        });
      });
  } catch {
    if (overlayEl) {
      const bodyEl =
        overlayEl.querySelector<HTMLElement>(".type-modal-body");
      if (bodyEl) {
        const displayName = moveName.replace(/-/g, " ");
        const capitalized =
          displayName.charAt(0).toUpperCase() + displayName.slice(1);
        bodyEl.innerHTML = `<p class="type-modal-empty">${t("modals.moveEmpty", { name: capitalized })}</p>`;
      }
    }
  } finally {
    isLoading = false;
  }
}

export function closeMovePokemonModal(): void {
  if (overlayEl) {
    overlayEl.remove();
    overlayEl = null;
    document.removeEventListener("keydown", handleKeydown);
  }
}
