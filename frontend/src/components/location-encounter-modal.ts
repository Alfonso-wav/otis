import gsap from "gsap";
import { GetLocationEncounters } from "../api";
import { t } from "../i18n";
import { createInlineDiglett, removeInlineDiglett } from "./sorting-overlay";

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
  if (e.key === "Escape") closeModal();
}

function closeModal(): void {
  if (overlayEl) {
    gsap.to(overlayEl, {
      opacity: 0,
      duration: 0.15,
      onComplete() {
        overlayEl?.remove();
        overlayEl = null;
      },
    });
    document.removeEventListener("keydown", handleKeydown);
  }
}

function formatLocationName(name: string): string {
  return name.replace(/-/g, " ");
}

export async function openLocationEncounterModal(locationName: string): Promise<void> {
  closeModal();

  const overlay = document.createElement("div");
  overlay.className = "type-modal-overlay";
  overlay.innerHTML = `
    <div class="type-modal">
      <div class="type-modal-header">
        <span>${formatLocationName(locationName)}</span>
        <button class="type-modal-close">&times;</button>
      </div>
      <div class="type-modal-body"></div>
    </div>`;

  document.body.appendChild(overlay);
  overlayEl = overlay;

  const bodyEl = overlay.querySelector<HTMLElement>(".type-modal-body")!;
  const diglettEl = createInlineDiglett(bodyEl, t("modals.locationLoading"));

  overlay.querySelector(".type-modal-close")!.addEventListener("click", closeModal);
  overlay.addEventListener("click", (e) => {
    if (e.target === overlay) closeModal();
  });
  document.addEventListener("keydown", handleKeydown);

  gsap.fromTo(overlay, { opacity: 0 }, { opacity: 1, duration: 0.15 });
  gsap.fromTo(
    overlay.querySelector(".type-modal")!,
    { scale: 0.9, y: 20 },
    { scale: 1, y: 0, duration: 0.2, ease: "power2.out" },
  );

  try {
    const encounters = await GetLocationEncounters(locationName);
    removeInlineDiglett(diglettEl);

    if (!encounters || encounters.length === 0) {
      bodyEl.innerHTML = `<p class="type-modal-empty">${t("modals.locationEmpty")}</p>`;
      return;
    }

    // Sort by MaxChance descending
    encounters.sort((a, b) => b.MaxChance - a.MaxChance);

    bodyEl.innerHTML = `
      <div class="type-modal-grid">
        ${encounters
          .map(
            (enc) => `
          <div class="type-modal-pokemon encounter-pokemon-item" data-name="${enc.PokemonName}">
            <img
              src="${spriteUrl(enc.PokemonName)}"
              alt="${enc.PokemonName}"
              loading="lazy"
              data-fallback="0"
              onerror="${spriteOnerror(enc.PokemonName)}"
            />
            <span class="type-modal-pokemon-name">${enc.PokemonName.replace(/-/g, " ")}</span>
            <span class="encounter-chance">${enc.MaxChance}%</span>
          </div>`,
          )
          .join("")}
      </div>`;

    const items = bodyEl.querySelectorAll<HTMLDivElement>(".encounter-pokemon-item");
    gsap.fromTo(
      items,
      { opacity: 0, y: 8 },
      { opacity: 1, y: 0, duration: 0.2, stagger: 0.02, ease: "power2.out" },
    );

    items.forEach((item) => {
      item.addEventListener("click", () => {
        const name = item.dataset.name;
        if (!name) return;
        closeModal();
        document
          .querySelector<HTMLButtonElement>('[data-tab="pokedex"]')
          ?.click();
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
  } catch {
    removeInlineDiglett(diglettEl);
    bodyEl.innerHTML = `<p class="type-modal-empty">${t("modals.locationError")}</p>`;
  }
}
