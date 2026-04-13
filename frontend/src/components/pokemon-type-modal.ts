import { t, typeName } from "../i18n";
import { getActiveId } from "../router";
import { setDetailOrigin } from "../pages/pokedex";

let overlayEl: HTMLDivElement | null = null;

function spriteUrl(name: string): string {
  return `https://img.pokemondb.net/sprites/black-white/normal/${name}.png`;
}

function spriteOnerror(name: string): string {
  const fb1 = `https://img.pokemondb.net/sprites/x-y/normal/${name}.png`;
  const fb2 = `https://img.pokemondb.net/sprites/home/normal/1x/${name}.png`;
  const fb3 = `/assets/sprites/home-normal/${name}.png`;
  return `var f=parseInt(this.dataset.fallback||'0');if(f===0){this.dataset.fallback='1';this.src='${fb1}'}else if(f===1){this.dataset.fallback='2';this.src='${fb2}'}else if(f===2){this.dataset.fallback='3';this.src='${fb3}'}else{this.onerror=null;this.style.visibility='hidden'}`;
}

function handleKeydown(e: KeyboardEvent): void {
  if (e.key === "Escape") closeTypeModal();
}

export function openTypeModal(
  regionName: string,
  typeKey: string,
  pokemonNames: string[],
): void {
  closeTypeModal();

  const overlay = document.createElement("div");
  overlay.className = "type-modal-overlay";

  const isEmpty = pokemonNames.length === 0;
  const gridHtml = isEmpty
    ? `<p class="type-modal-empty">${t("modals.typeRegionEmpty", { type: typeName(typeKey), region: regionName })}</p>`
    : `<div class="type-modal-grid">${pokemonNames
        .map(
          (name) => `
      <div class="type-modal-pokemon" data-name="${name}">
        <img
          src="${spriteUrl(name)}"
          alt="${name}"
          loading="lazy"
          data-fallback="0"
          onerror="${spriteOnerror(name)}"
        />
        <span class="type-modal-pokemon-name">${name.replace(/-/g, " ")}</span>
      </div>`,
        )
        .join("")}</div>`;

  overlay.innerHTML = `
    <div class="type-modal">
      <div class="type-modal-header">
        <span>${t("modals.typeInRegion", { type: typeName(typeKey), region: regionName.charAt(0).toUpperCase() + regionName.slice(1) })}</span>
        <button class="type-modal-close" id="type-modal-close">&times;</button>
      </div>
      <div class="type-modal-body">
        ${gridHtml}
      </div>
    </div>`;

  document.body.appendChild(overlay);
  overlayEl = overlay;

  overlay.querySelector("#type-modal-close")!.addEventListener("click", closeTypeModal);
  overlay.addEventListener("click", (e) => {
    if (e.target === overlay) closeTypeModal();
  });
  document.addEventListener("keydown", handleKeydown);

  overlay
    .querySelectorAll<HTMLDivElement>(".type-modal-pokemon")
    .forEach((item) => {
      item.addEventListener("click", () => {
        const name = item.dataset.name;
        if (!name) return;
        const origin = getActiveId();
        closeTypeModal();
        document
          .querySelector<HTMLButtonElement>('[data-tab="pokedex"]')
          ?.click();
        setDetailOrigin(origin);
        const input = document.getElementById(
          "search-input",
        ) as HTMLInputElement;
        const btn = document.getElementById(
          "search-btn",
        ) as HTMLButtonElement;
        if (input && btn) {
          input.value = name;
          btn.click();
        }
      });
    });
}

export function closeTypeModal(): void {
  if (overlayEl) {
    overlayEl.remove();
    overlayEl = null;
    document.removeEventListener("keydown", handleKeydown);
  }
}
