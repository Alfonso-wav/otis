import gsap from "gsap";
import { ListTypes, GetType } from "../../wailsjs/go/app/App";
import type { PokemonListItem } from "../types";

const POKEMON_DISPLAY_LIMIT = 50;

const GAME_TYPES = new Set([
  "normal", "fire", "water", "electric", "grass", "ice",
  "fighting", "poison", "ground", "flying", "psychic", "bug",
  "rock", "ghost", "dragon", "dark", "steel", "fairy",
]);

let container: HTMLElement;
let initialized = false;

function idFromURL(url: string): string {
  const parts = url.replace(/\/$/, "").split("/");
  return parts[parts.length - 1];
}

function spriteURL(id: string): string {
  return `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${id}.png`;
}

function typeHeaderClass(typeName: string): string {
  return GAME_TYPES.has(typeName) ? `type-header-${typeName}` : "type-header-default";
}

function renderTypeCards(types: PokemonListItem[]): void {
  const filtered = types.filter((t) => GAME_TYPES.has(t.Name));
  container.innerHTML = `<div class="types-grid" id="types-grid"></div>`;
  const grid = document.getElementById("types-grid") as HTMLDivElement;

  grid.innerHTML = filtered
    .map(
      (t) => `
    <div class="type-card" data-type="${t.Name}">
      <div class="type-card__header ${typeHeaderClass(t.Name)}">
        <span class="type-card__name">${t.Name}</span>
        <span class="type-card__chevron">▼</span>
      </div>
      <div class="type-card__pokemon hidden"></div>
    </div>`,
    )
    .join("");

  const cards = grid.querySelectorAll<HTMLDivElement>(".type-card");
  gsap.fromTo(
    cards,
    { opacity: 0, y: 20 },
    { opacity: 1, y: 0, duration: 0.3, stagger: 0.05, ease: "power2.out" },
  );

  cards.forEach((card) => {
    card.querySelector(".type-card__header")!.addEventListener("click", () =>
      toggleType(card),
    );
  });
}

async function toggleType(card: HTMLDivElement): Promise<void> {
  const typeName = card.dataset.type!;
  const pokemonPanel = card.querySelector<HTMLDivElement>(".type-card__pokemon")!;
  const isExpanded = card.classList.contains("expanded");

  if (isExpanded) {
    card.classList.remove("expanded");
    gsap.to(pokemonPanel, {
      height: 0,
      opacity: 0,
      duration: 0.25,
      ease: "power2.in",
      onComplete() {
        pokemonPanel.classList.add("hidden");
        pokemonPanel.style.removeProperty("height");
        pokemonPanel.style.removeProperty("opacity");
      },
    });
    return;
  }

  card.classList.add("expanded");

  // Already loaded — just show
  if (pokemonPanel.innerHTML.trim() !== "") {
    pokemonPanel.classList.remove("hidden");
    gsap.fromTo(pokemonPanel, { opacity: 0 }, { opacity: 1, duration: 0.2 });
    return;
  }

  pokemonPanel.innerHTML =
    '<p class="types-loading" style="padding:1rem">Cargando...</p>';
  pokemonPanel.classList.remove("hidden");

  try {
    const detail = await GetType(typeName);
    const allEntries = detail.Pokemon || [];
    const entries = allEntries.slice(0, POKEMON_DISPLAY_LIMIT);
    const remaining = allEntries.length - entries.length;

    const pokemonHTML = entries
      .map((p) => {
        const id = idFromURL(p.URL);
        return `<div class="type-pokemon-item" data-name="${p.Name}">
          <img src="${spriteURL(id)}" alt="${p.Name}" loading="lazy" />
          <span>${p.Name}</span>
        </div>`;
      })
      .join("");

    const moreHTML =
      remaining > 0
        ? `<div class="type-pokemon-more">+${remaining} más</div>`
        : "";

    pokemonPanel.innerHTML = `<div class="type-pokemon-grid">${pokemonHTML}</div>${moreHTML}`;

    const items = pokemonPanel.querySelectorAll(".type-pokemon-item");
    gsap.fromTo(
      items,
      { opacity: 0, scale: 0.85 },
      { opacity: 1, scale: 1, duration: 0.2, stagger: 0.02, ease: "back.out(1.2)" },
    );

    pokemonPanel
      .querySelectorAll<HTMLDivElement>(".type-pokemon-item")
      .forEach((item) => {
        item.addEventListener("click", () => {
          const name = item.dataset.name;
          if (!name) return;
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
  } catch (err: unknown) {
    pokemonPanel.innerHTML = `<p class="types-loading" style="color:#e53e3e;padding:1rem">${String(err)}</p>`;
  }
}

export function initTypes(): void {
  container = document.getElementById("tab-types") as HTMLElement;

  const tabBtn = document.querySelector<HTMLButtonElement>(
    '[data-tab="types"]',
  );
  if (!tabBtn) return;

  tabBtn.addEventListener("click", async () => {
    if (initialized) return;
    initialized = true;
    container.innerHTML = '<p class="types-loading">Cargando tipos...</p>';
    try {
      const data = await ListTypes();
      renderTypeCards(data.Results || []);
    } catch (err: unknown) {
      container.innerHTML = `<p class="types-loading" style="color:#e53e3e">${String(err)}</p>`;
    }
  });
}
