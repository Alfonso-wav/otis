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
  if (e.key === "Escape") closeAbilityPokemonModal();
}

export function openAbilityPokemonModal(
  abilityName: string,
  pokemonNames: string[],
): void {
  closeAbilityPokemonModal();

  const overlay = document.createElement("div");
  overlay.className = "type-modal-overlay";

  const displayName = abilityName.replace(/-/g, " ");
  const capitalized = displayName.charAt(0).toUpperCase() + displayName.slice(1);

  const isEmpty = pokemonNames.length === 0;
  const gridHtml = isEmpty
    ? `<p class="type-modal-empty">No se encontraron Pokémon con la habilidad <strong>${capitalized}</strong>.</p>`
    : `<div class="type-modal-grid">${pokemonNames
        .map(
          (name) => `
      <div class="type-modal-pokemon">
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

  overlay.innerHTML = `
    <div class="type-modal">
      <div class="type-modal-header">
        <span>${capitalized} (${pokemonNames.length})</span>
        <button class="type-modal-close" id="ability-modal-close">&times;</button>
      </div>
      <div class="type-modal-body">
        ${gridHtml}
      </div>
    </div>`;

  document.body.appendChild(overlay);
  overlayEl = overlay;

  overlay.querySelector("#ability-modal-close")!.addEventListener("click", closeAbilityPokemonModal);
  overlay.addEventListener("click", (e) => {
    if (e.target === overlay) closeAbilityPokemonModal();
  });
  document.addEventListener("keydown", handleKeydown);
}

export function closeAbilityPokemonModal(): void {
  if (overlayEl) {
    overlayEl.remove();
    overlayEl = null;
    document.removeEventListener("keydown", handleKeydown);
  }
}
