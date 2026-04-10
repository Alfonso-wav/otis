import { t } from "../../i18n";

const TYPES = [
  "normal", "fire", "water", "electric", "grass", "ice",
  "fighting", "poison", "ground", "flying", "psychic", "bug",
  "rock", "ghost", "dragon", "dark", "steel", "fairy",
] as const;

type PokemonType = typeof TYPES[number];

// chart[attacker][defender] = multiplier; only non-1 values listed
const OVERRIDES: Partial<Record<PokemonType, Partial<Record<PokemonType, number>>>> = {
  normal:   { rock: 0.5, ghost: 0, steel: 0.5 },
  fire:     { fire: 0.5, water: 0.5, grass: 2, ice: 2, bug: 2, rock: 0.5, dragon: 0.5, steel: 2 },
  water:    { fire: 2, water: 0.5, grass: 0.5, ground: 2, rock: 2, dragon: 0.5 },
  electric: { water: 2, electric: 0.5, grass: 0.5, ground: 0, flying: 2, dragon: 0.5 },
  grass:    { fire: 0.5, water: 2, grass: 0.5, poison: 0.5, ground: 2, flying: 0.5, bug: 0.5, rock: 2, dragon: 0.5, steel: 0.5 },
  ice:      { water: 0.5, grass: 2, ice: 0.5, ground: 2, flying: 2, dragon: 2, steel: 0.5 },
  fighting: { normal: 2, ice: 2, poison: 0.5, flying: 0.5, psychic: 0.5, bug: 0.5, rock: 2, ghost: 0, dark: 2, steel: 2, fairy: 0.5 },
  poison:   { grass: 2, poison: 0.5, ground: 0.5, rock: 0.5, ghost: 0.5, steel: 0, fairy: 2 },
  ground:   { fire: 2, electric: 2, grass: 0.5, poison: 2, flying: 0, bug: 0.5, rock: 2, steel: 2 },
  flying:   { electric: 0.5, grass: 2, fighting: 2, bug: 2, rock: 0.5, steel: 0.5 },
  psychic:  { fighting: 2, poison: 2, psychic: 0.5, dark: 0, steel: 0.5 },
  bug:      { fire: 0.5, grass: 2, fighting: 0.5, flying: 0.5, psychic: 2, ghost: 0.5, dark: 2, steel: 0.5, fairy: 0.5 },
  rock:     { fire: 2, ice: 2, fighting: 0.5, ground: 0.5, flying: 2, bug: 2, steel: 0.5 },
  ghost:    { normal: 0, psychic: 2, ghost: 2, dark: 0.5 },
  dragon:   { dragon: 2, steel: 0.5, fairy: 0 },
  dark:     { fighting: 0.5, psychic: 2, ghost: 2, dark: 0.5, fairy: 0.5 },
  steel:    { fire: 0.5, water: 0.5, electric: 0.5, ice: 2, rock: 2, steel: 0.5, fairy: 2 },
  fairy:    { fire: 0.5, fighting: 2, poison: 0.5, dragon: 2, dark: 2, steel: 0.5 },
};

function effectiveness(attacker: PokemonType, defender: PokemonType): number {
  return OVERRIDES[attacker]?.[defender] ?? 1;
}

function cellLabel(mult: number): string {
  if (mult === 2) return "2×";
  if (mult === 0.5) return "½";
  if (mult === 0) return "0";
  return "";
}

function cellClass(mult: number): string {
  if (mult === 2) return "tc-super";
  if (mult === 0.5) return "tc-resisted";
  if (mult === 0) return "tc-immune";
  return "tc-neutral";
}

function typeHeader(type: PokemonType): string {
  const name = t(`typeNames.${type}`);
  return `<img src="/assets/types/${type}.svg" class="tc-icon" alt="${name}">`;
}

const filteredTypes = new Set<PokemonType>();
let isDragging = false;
let dragAxis: "col" | "row" | null = null;

function renderChart(panel: HTMLElement): void {
  const title = t("typeChart.title");
  const attackingLabel = t("typeChart.attackingLabel");
  const defendingLabel = t("typeChart.defendingLabel");

  const visibleTypes = TYPES.filter((tp) => !filteredTypes.has(tp));

  const headerCells = visibleTypes.map(
    (def) =>
      `<th class="tc-col-header" data-type="${def}" title="${t(`typeNames.${def}`)}">${typeHeader(def)}</th>`,
  ).join("");

  const rows = visibleTypes.map((atk) => {
    const cells = visibleTypes.map((def) => {
      const mult = effectiveness(atk, def);
      return `<td class="tc-cell ${cellClass(mult)}">${cellLabel(mult)}</td>`;
    }).join("");
    return `<tr><th class="tc-row-header" data-type="${atk}">${typeHeader(atk)}</th>${cells}</tr>`;
  }).join("");

  const filterBar = filteredTypes.size > 0
    ? `<div class="tc-filter-bar"><button class="tc-restore-btn">${t("typeChart.restoreTypes")}</button></div>`
    : "";

  panel.innerHTML = `
    <div class="type-chart-wrap">
      <h2 class="type-chart-title">${title}</h2>
      <div class="type-chart-labels">
        <span class="tc-attacking-label">${attackingLabel}</span>
        <span class="tc-defending-label">${defendingLabel}</span>
      </div>
      <div class="type-chart-scroll">
        <table class="type-chart-table" aria-label="${title}">
          <thead>
            <tr>
              <th class="tc-corner"></th>
              ${headerCells}
            </tr>
          </thead>
          <tbody>${rows}</tbody>
        </table>
      </div>
      ${filterBar}
    </div>
  `;

  attachFilterListeners(panel);
}

function attachFilterListeners(panel: HTMLElement): void {
  const headers = panel.querySelectorAll<HTMLElement>(
    ".tc-col-header[data-type], .tc-row-header[data-type]",
  );

  headers.forEach((th) => {
    th.addEventListener("mousedown", (e) => {
      e.preventDefault();
      isDragging = true;
      dragAxis = th.classList.contains("tc-col-header") ? "col" : "row";
      const type = th.dataset.type as PokemonType;
      if (type) {
        filteredTypes.add(type);
        renderChart(panel);
      }
    });

    th.addEventListener("mouseenter", () => {
      if (!isDragging) return;
      const thAxis = th.classList.contains("tc-col-header") ? "col" : "row";
      if (thAxis !== dragAxis) return;
      const type = th.dataset.type as PokemonType;
      if (type && !filteredTypes.has(type)) {
        filteredTypes.add(type);
        renderChart(panel);
      }
    });

    th.addEventListener("touchstart", (e) => {
      e.preventDefault();
      isDragging = true;
      dragAxis = th.classList.contains("tc-col-header") ? "col" : "row";
      const type = th.dataset.type as PokemonType;
      if (type) {
        filteredTypes.add(type);
        renderChart(panel);
      }
    }, { passive: false });
  });

  // Touch drag — detect type under finger via elementFromPoint
  const scroll = panel.querySelector<HTMLElement>(".type-chart-scroll");
  if (scroll) {
    scroll.addEventListener("touchmove", (e) => {
      if (!isDragging) return;
      e.preventDefault();
      const touch = e.touches[0];
      const el = document.elementFromPoint(touch.clientX, touch.clientY);
      if (!el) return;
      const target = el.closest<HTMLElement>(".tc-col-header[data-type], .tc-row-header[data-type]");
      if (!target) return;
      const thAxis = target.classList.contains("tc-col-header") ? "col" : "row";
      if (thAxis !== dragAxis) return;
      const type = target.dataset.type as PokemonType;
      if (type && !filteredTypes.has(type)) {
        filteredTypes.add(type);
        renderChart(panel);
      }
    }, { passive: false });
  }

  const restoreBtn = panel.querySelector<HTMLButtonElement>(".tc-restore-btn");
  if (restoreBtn) {
    restoreBtn.addEventListener("click", () => {
      filteredTypes.clear();
      renderChart(panel);
    });
  }
}

export function initTypeChart(panel: HTMLElement): void {
  filteredTypes.clear();
  renderChart(panel);

  document.addEventListener("mouseup", () => {
    isDragging = false;
    dragAxis = null;
  });

  document.addEventListener("touchend", () => {
    isDragging = false;
    dragAxis = null;
  });

  document.addEventListener("locale-changed", () => {
    renderChart(panel);
  });
}
