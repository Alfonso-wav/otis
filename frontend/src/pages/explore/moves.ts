import gsap from "gsap";
import { GetMove } from "../../../wailsjs/go/app/App";
import type { core } from "../../../wailsjs/go/models";

// Lista de movimientos representativos para búsqueda/demo
// (PokeAPI tiene +900 movimientos; usamos lazy fetch individual)
const KNOWN_MOVES = [
  "pound", "tackle", "scratch", "growl", "tail-whip", "ember", "water-gun",
  "vine-whip", "thunderbolt", "flamethrower", "surf", "fly", "earthquake",
  "psychic", "ice-beam", "blizzard", "thunder", "solar-beam", "dragon-claw",
  "shadow-ball", "brick-break", "rock-slide", "iron-tail", "aerial-ace",
  "hyper-beam", "giga-impact", "close-combat", "flare-blitz", "waterfall",
  "leaf-blade", "night-slash", "swords-dance", "nasty-plot", "calm-mind",
  "dragon-dance", "will-o-wisp", "toxic", "stealth-rock", "spore",
  "stone-edge", "crunch", "dark-pulse", "flash-cannon", "focus-blast",
  "energy-ball", "aura-sphere", "draco-meteor", "outrage", "moonblast",
];

type Category = "all" | "physical" | "special" | "status";
type FilterType = "all" | string;

interface MoveState {
  loaded: Map<string, core.Move>;
  selectedCategory: Category;
  selectedType: FilterType;
  searchQuery: string;
}

const state: MoveState = {
  loaded: new Map(),
  selectedCategory: "all",
  selectedType: "all",
  searchQuery: "",
};

const TYPE_COLORS: Record<string, string> = {
  normal: "#a0aec0", fire: "#f6ad55", water: "#63b3ed", grass: "#68d391",
  electric: "#f6e05e", psychic: "#f687b3", ice: "#76e4f7", fighting: "#c05621",
  poison: "#9f7aea", ground: "#d69e2e", flying: "#90cdf4", bug: "#a8e063",
  rock: "#b7791f", ghost: "#553c9a", dragon: "#7f9cf5", dark: "#4a5568",
  steel: "#718096", fairy: "#fbb6ce",
};

function typeColor(type: string): string {
  return TYPE_COLORS[type] ?? "#718096";
}

function categoryIcon(cat: string): string {
  const icons: Record<string, string> = {
    physical: "⚔️",
    special: "✨",
    status: "💤",
  };
  return icons[cat] ?? "❓";
}

function filteredMoves(): string[] {
  return KNOWN_MOVES.filter((name) => {
    const move = state.loaded.get(name);
    if (!move) return state.searchQuery === "" || name.includes(state.searchQuery);

    const matchSearch =
      state.searchQuery === "" ||
      move.Name.includes(state.searchQuery.toLowerCase());
    const matchCat =
      state.selectedCategory === "all" ||
      move.Category === state.selectedCategory;
    const matchType =
      state.selectedType === "all" || move.Type === state.selectedType;

    return matchSearch && matchCat && matchType;
  });
}

function renderMoveCard(name: string): string {
  const move = state.loaded.get(name);
  if (!move) {
    return `<div class="move-card move-card--loading" data-move="${name}">
      <span class="move-card__name">${name.replace(/-/g, " ")}</span>
    </div>`;
  }
  return `<div class="move-card" data-move="${name}">
    <div class="move-card__header">
      <span class="move-card__name">${move.Name.replace(/-/g, " ")}</span>
      <span class="move-card__type" style="background:${typeColor(move.Type)}">${move.Type}</span>
      <span class="move-card__cat" title="${move.Category}">${categoryIcon(move.Category)}</span>
    </div>
    <div class="move-card__stats">
      <span class="move-stat"><b>PWR</b> ${move.Power || "—"}</span>
      <span class="move-stat"><b>ACC</b> ${move.Accuracy || "—"}%</span>
      <span class="move-stat"><b>PP</b> ${move.PP}</span>
    </div>
    ${move.Description ? `<p class="move-card__desc">${move.Description.replace(/\n/g, " ")}</p>` : ""}
  </div>`;
}

function renderGrid(container: HTMLElement): void {
  const grid = container.querySelector<HTMLElement>("#moves-grid");
  if (!grid) return;

  const names = filteredMoves();
  if (names.length === 0) {
    grid.innerHTML = '<p class="loading">No se encontraron movimientos.</p>';
    return;
  }

  grid.innerHTML = names.map(renderMoveCard).join("");

  const cards = grid.querySelectorAll<HTMLElement>(".move-card:not(.move-card--loading)");
  gsap.fromTo(
    cards,
    { opacity: 0, y: 10 },
    { opacity: 1, y: 0, duration: 0.2, stagger: 0.02, ease: "power2.out" },
  );

  // Lazy-load unloaded moves
  const unloaded = names.filter((n) => !state.loaded.has(n));
  loadMovesLazy(unloaded, container);
}

async function loadMovesLazy(
  names: string[],
  container: HTMLElement,
): Promise<void> {
  for (const name of names) {
    try {
      const move = await GetMove(name);
      state.loaded.set(name, move);
      const card = container.querySelector<HTMLElement>(
        `[data-move="${name}"]`,
      );
      if (card) {
        const tmp = document.createElement("div");
        tmp.innerHTML = renderMoveCard(name);
        const newCard = tmp.firstElementChild as HTMLElement;
        gsap.fromTo(newCard, { opacity: 0 }, { opacity: 1, duration: 0.3 });
        card.replaceWith(newCard);
      }
    } catch {
      // silently skip unavailable moves
    }
  }
}

export function initMoves(container: HTMLElement): void {
  container.innerHTML = `
    <div class="section-header"><h2>Movimientos</h2></div>
    <div class="moves-controls">
      <input
        type="text"
        id="moves-search"
        class="explore-input"
        placeholder="Buscar movimiento..."
      />
      <div class="moves-filters">
        <button class="filter-btn active" data-cat="all">Todos</button>
        <button class="filter-btn" data-cat="physical">⚔️ Físico</button>
        <button class="filter-btn" data-cat="special">✨ Especial</button>
        <button class="filter-btn" data-cat="status">💤 Estado</button>
      </div>
    </div>
    <div id="moves-grid" class="moves-grid"></div>`;

  const searchInput = container.querySelector<HTMLInputElement>("#moves-search")!;
  searchInput.addEventListener("input", () => {
    state.searchQuery = searchInput.value.trim().toLowerCase().replace(/\s+/g, "-");
    renderGrid(container);
  });

  container.querySelectorAll<HTMLButtonElement>(".filter-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      container
        .querySelectorAll(".filter-btn")
        .forEach((b) => b.classList.remove("active"));
      btn.classList.add("active");
      state.selectedCategory = (btn.dataset.cat ?? "all") as Category;
      renderGrid(container);
    });
  });

  renderGrid(container);
}
