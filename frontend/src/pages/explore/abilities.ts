import gsap from "gsap";
import { GetAbility } from "../../../wailsjs/go/app/App";
import type { core } from "../../../wailsjs/go/models";

const KNOWN_ABILITIES = [
  "overgrow", "blaze", "torrent", "shield-dust", "shed-skin", "compound-eyes",
  "sturdy", "keen-eye", "run-away", "intimidate", "static", "volt-absorb",
  "water-absorb", "flash-fire", "swift-swim", "chlorophyll", "hustle",
  "serene-grace", "synchronize", "levitate", "thick-fat", "natural-cure",
  "lightning-rod", "wonder-guard", "speed-boost", "battle-armor", "rough-skin",
  "drought", "drizzle", "sand-stream", "snow-warning", "pressure",
  "multiscale", "regenerator", "magic-guard", "adaptability", "moxie",
  "justified", "technician", "sand-rush", "unburden", "gale-wings",
  "protean", "mega-launcher", "strong-jaw", "refrigerate", "pixilate",
];

interface AbilityState {
  loaded: Map<string, core.Ability>;
  searchQuery: string;
}

const state: AbilityState = {
  loaded: new Map(),
  searchQuery: "",
};

function filteredAbilities(): string[] {
  const q = state.searchQuery.toLowerCase();
  return KNOWN_ABILITIES.filter((name) =>
    q === "" || name.includes(q),
  );
}

function renderAbilityCard(name: string): string {
  const ability = state.loaded.get(name);
  const displayName = name.replace(/-/g, " ");
  if (!ability) {
    return `<div class="ability-card ability-card--loading" data-ability="${name}">
      <span class="ability-card__name">${displayName}</span>
    </div>`;
  }
  const pokemonPreview = (ability.Pokemon ?? []).slice(0, 5);
  return `<div class="ability-card" data-ability="${name}">
    <div class="ability-card__header">
      <span class="ability-card__name">${displayName}</span>
      <span class="ability-card__count">${(ability.Pokemon ?? []).length} Pokémon</span>
    </div>
    ${ability.Description ? `<p class="ability-card__desc">${ability.Description.replace(/\n/g, " ")}</p>` : ""}
    ${
      pokemonPreview.length > 0
        ? `<div class="ability-pokemon-list">
            ${pokemonPreview.map((p) => `<span class="ability-pokemon-tag">${p}</span>`).join("")}
            ${(ability.Pokemon ?? []).length > 5 ? `<span class="ability-pokemon-more">+${(ability.Pokemon ?? []).length - 5} más</span>` : ""}
          </div>`
        : ""
    }
  </div>`;
}

function renderGrid(container: HTMLElement): void {
  const grid = container.querySelector<HTMLElement>("#abilities-grid");
  if (!grid) return;

  const names = filteredAbilities();
  if (names.length === 0) {
    grid.innerHTML = '<p class="explore-loading">No se encontraron habilidades.</p>';
    return;
  }

  grid.innerHTML = names.map(renderAbilityCard).join("");

  const cards = grid.querySelectorAll<HTMLElement>(
    ".ability-card:not(.ability-card--loading)",
  );
  gsap.fromTo(
    cards,
    { opacity: 0, y: 10 },
    { opacity: 1, y: 0, duration: 0.2, stagger: 0.02, ease: "power2.out" },
  );

  const unloaded = names.filter((n) => !state.loaded.has(n));
  loadAbilitiesLazy(unloaded, container);
}

async function loadAbilitiesLazy(
  names: string[],
  container: HTMLElement,
): Promise<void> {
  for (const name of names) {
    try {
      const ability = await GetAbility(name);
      state.loaded.set(name, ability);
      const card = container.querySelector<HTMLElement>(
        `[data-ability="${name}"]`,
      );
      if (card) {
        const tmp = document.createElement("div");
        tmp.innerHTML = renderAbilityCard(name);
        const newCard = tmp.firstElementChild as HTMLElement;
        gsap.fromTo(newCard, { opacity: 0 }, { opacity: 1, duration: 0.3 });
        card.replaceWith(newCard);
      }
    } catch {
      // silently skip
    }
  }
}

export function initAbilities(container: HTMLElement): void {
  container.innerHTML = `
    <div class="abilities-controls">
      <input
        type="text"
        id="abilities-search"
        class="explore-input"
        placeholder="Buscar habilidad..."
      />
    </div>
    <div id="abilities-grid" class="abilities-grid"></div>`;

  const searchInput =
    container.querySelector<HTMLInputElement>("#abilities-search")!;
  searchInput.addEventListener("input", () => {
    state.searchQuery = searchInput.value.trim().toLowerCase().replace(/\s+/g, "-");
    renderGrid(container);
  });

  renderGrid(container);
}
