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
  const abbr = name.slice(0, 3);
  return `<img src="/assets/types/${type}.svg" class="tc-icon" alt="${name}"><span class="tc-abbr">${abbr}</span>`;
}

export function initTypeChart(panel: HTMLElement): void {
  const title = t("typeChart.title");
  const attackingLabel = t("typeChart.attackingLabel");
  const defendingLabel = t("typeChart.defendingLabel");

  const headerCells = TYPES.map(
    (def) => `<th class="tc-col-header" title="${t(`typeNames.${def}`)}">${typeHeader(def)}</th>`,
  ).join("");

  const rows = TYPES.map((atk) => {
    const cells = TYPES.map((def) => {
      const mult = effectiveness(atk, def);
      return `<td class="tc-cell ${cellClass(mult)}">${cellLabel(mult)}</td>`;
    }).join("");
    return `<tr><th class="tc-row-header">${typeHeader(atk)}</th>${cells}</tr>`;
  }).join("");

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
          <tbody>
            ${rows}
          </tbody>
        </table>
      </div>
    </div>
  `;
}
