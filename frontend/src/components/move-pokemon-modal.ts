import { GetMove } from "../api";
import { t } from "../i18n";
import { createInlineDiglett, removeInlineDiglett } from "./sorting-overlay";

const SPRITE_BASE = "/assets/sprites";
const CDN_FALLBACK = "https://img.pokemondb.net/sprites/home/normal";

let overlayEl: HTMLDivElement | null = null;

function spriteUrl(name: string): string {
  return `${SPRITE_BASE}/${name}.png`;
}

function cdnUrl(name: string): string {
  return `${CDN_FALLBACK}/${name}.png`;
}

function handleKeydown(e: KeyboardEvent): void {
  if (e.key === "Escape") closeMovePokemonModal();
}

function navigateToPokemon(name: string): void {
  closeMovePokemonModal();
  (document.querySelector('[data-tab="pokedex"]') as HTMLElement)?.click();
  const input = document.getElementById("search-input") as HTMLInputElement;
  const btn = document.getElementById("search-btn") as HTMLButtonElement;
  if (input && btn) {
    input.value = name;
    btn.click();
  }
}

export async function openMovePokemonModal(moveName: string): Promise<void> {
  closeMovePokemonModal();

  const overlay = document.createElement("div");
  overlay.className = "type-modal-overlay";

  const displayName = moveName.replace(/-/g, " ");
  const capitalized = displayName.charAt(0).toUpperCase() + displayName.slice(1);

  // Show loading state
  overlay.innerHTML = `
    <div class="type-modal">
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

  overlay.querySelector("#move-modal-close")!.addEventListener("click", closeMovePokemonModal);
  overlay.addEventListener("click", (e) => {
    if (e.target === overlay) closeMovePokemonModal();
  });
  document.addEventListener("keydown", handleKeydown);

  // Lazy fetch LearnedBy
  try {
    const move = await GetMove(moveName);
    const pokemonNames = move.LearnedBy ?? [];
    removeInlineDiglett(diglettEl);

    const header = overlay.querySelector<HTMLElement>(".type-modal-header span");
    if (header) header.textContent = `${capitalized} (${pokemonNames.length})`;

    if (pokemonNames.length === 0) {
      bodyEl.innerHTML = `<p class="type-modal-empty">${t("modals.moveEmpty", { name: capitalized })}</p>`;
    } else {
      bodyEl.innerHTML = `<div class="type-modal-grid">${pokemonNames
        .map(
          (name) => `
        <div class="type-modal-pokemon" data-name="${name}">
          <img
            src="${spriteUrl(name)}"
            alt="${name}"
            loading="lazy"
            onerror="if(!this.dataset.fallback){this.dataset.fallback='1';this.src='${cdnUrl(name)}'}"
          />
          <span class="type-modal-pokemon-name">${name.replace(/-/g, " ")}</span>
        </div>`,
        )
        .join("")}</div>`;

      bodyEl.querySelectorAll<HTMLElement>(".type-modal-pokemon").forEach((el) => {
        el.addEventListener("click", () => {
          const pName = el.dataset.name;
          if (pName) navigateToPokemon(pName);
        });
      });
    }
  } catch {
    removeInlineDiglett(diglettEl);
    bodyEl.innerHTML = `<p class="type-modal-empty">${t("modals.moveEmpty", { name: capitalized })}</p>`;
  }
}

export function closeMovePokemonModal(): void {
  if (overlayEl) {
    overlayEl.remove();
    overlayEl = null;
    document.removeEventListener("keydown", handleKeydown);
  }
}
