import { t } from "../../i18n";
import { renderTypeHeatmap, disposeTypeHeatmap } from "../../charts/type-heatmap";
import { renderTypeRadar, disposeDefenseRadar } from "../../charts/type-defense-radar";

export const ALL_TYPES = [
  "normal", "fire", "water", "electric", "grass", "ice",
  "fighting", "poison", "ground", "flying", "psychic", "bug",
  "rock", "ghost", "dragon", "dark", "steel", "fairy",
] as const;

export type PokemonType = typeof ALL_TYPES[number];

// chart[attacker][defender] = multiplier; only non-1 values listed
export const OVERRIDES: Partial<Record<PokemonType, Partial<Record<PokemonType, number>>>> = {
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

export const TYPE_COLORS: Record<string, string> = {
  normal: "#a0aec0", fire: "#f6ad55", water: "#63b3ed", grass: "#68d391",
  electric: "#f6e05e", psychic: "#f687b3", ice: "#76e4f7", fighting: "#c05621",
  poison: "#9f7aea", ground: "#d69e2e", flying: "#90cdf4", bug: "#a8e063",
  rock: "#b7791f", ghost: "#553c9a", dragon: "#7f9cf5", dark: "#4a5568",
  steel: "#718096", fairy: "#fbb6ce",
};

export function effectiveness(attacker: string, defender: string): number {
  return (OVERRIDES as Record<string, Record<string, number> | undefined>)[attacker]?.[defender] ?? 1;
}

let selectedRadarType: string | null = null;
let showDefense = true;
let showOffense = true;

function renderChart(panel: HTMLElement): void {
  const title = t("typeChart.title");
  const radarHint = selectedRadarType
    ? `<h3 class="tc-radar-title">${t("typeChart.radarTitle")}: ${t(`typeNames.${selectedRadarType}`)}</h3>`
    : `<p class="tc-radar-hint">${t("typeChart.selectTypeForRadar")}</p>`;

  const togglesHtml = selectedRadarType
    ? `
      <div class="tc-radar-toggles" role="group" aria-label="${t("typeChart.radarToggles")}">
        <label class="tc-radar-toggle tc-radar-toggle--defense">
          <input type="checkbox" id="tc-toggle-defense" ${showDefense ? "checked" : ""} />
          <span class="tc-radar-toggle-dot tc-radar-toggle-dot--defense"></span>
          <span class="tc-radar-toggle-label">${t("typeChart.toggleDefensive")}</span>
        </label>
        <label class="tc-radar-toggle tc-radar-toggle--offense">
          <input type="checkbox" id="tc-toggle-offense" ${showOffense ? "checked" : ""} />
          <span class="tc-radar-toggle-dot tc-radar-toggle-dot--offense"></span>
          <span class="tc-radar-toggle-label">${t("typeChart.toggleOffensive")}</span>
        </label>
      </div>
    `
    : "";

  const anyVisible = showDefense || showOffense;

  panel.innerHTML = `
    <div class="type-chart-wrap">
      <h2 class="type-chart-title">${title}</h2>
      <div class="tc-heatmap-container" id="tc-heatmap"></div>
      <div class="tc-radar-section">
        ${radarHint}
        ${togglesHtml}
        <div class="tc-radar-container${selectedRadarType && anyVisible ? "" : " hidden"}" id="tc-radar"></div>
      </div>
    </div>
  `;

  const heatmapEl = panel.querySelector<HTMLElement>("#tc-heatmap");
  if (heatmapEl) {
    renderTypeHeatmap(heatmapEl, (type) => {
      selectedRadarType = type;
      renderChart(panel);
    });
  }

  if (selectedRadarType && anyVisible) {
    const radarEl = panel.querySelector<HTMLElement>("#tc-radar");
    if (radarEl) {
      renderTypeRadar(radarEl, selectedRadarType, { showDefense, showOffense });
    }
  } else {
    disposeDefenseRadar();
  }

  const defenseToggle = panel.querySelector<HTMLInputElement>("#tc-toggle-defense");
  const offenseToggle = panel.querySelector<HTMLInputElement>("#tc-toggle-offense");

  if (defenseToggle) {
    defenseToggle.addEventListener("change", () => {
      showDefense = defenseToggle.checked;
      renderChart(panel);
    });
  }
  if (offenseToggle) {
    offenseToggle.addEventListener("change", () => {
      showOffense = offenseToggle.checked;
      renderChart(panel);
    });
  }
}

export function initTypeChart(panel: HTMLElement): void {
  selectedRadarType = null;
  showDefense = true;
  showOffense = true;
  renderChart(panel);

  document.addEventListener("locale-changed", () => {
    disposeTypeHeatmap();
    disposeDefenseRadar();
    renderChart(panel);
  });
}
