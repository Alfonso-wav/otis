import { CalculateEVs, GetNatures } from "./api";
import { core } from "../wailsjs/go/models";
import type { Pokemon } from "./types";

let cachedNatures: core.Nature[] = [];

const STAT_NAMES: Record<string, string> = {
  hp: "HP",
  attack: "Ataque",
  defense: "Defensa",
  spAttack: "At. Esp.",
  spDefense: "Def. Esp.",
  speed: "Velocidad",
};

export async function loadNatures(): Promise<core.Nature[]> {
  if (cachedNatures.length === 0) {
    cachedNatures = await GetNatures();
    cachedNatures.sort((a, b) => a.name.localeCompare(b.name));
  }
  return cachedNatures;
}

export function renderEVCalculatorForm(pokemon: Pokemon): string {
  return `
    <div class="ev-calculator">
      <h3>Calculadora de EVs</h3>
      <p class="ev-calc-info">Ingresa los datos de tu ${pokemon.Name} para calcular sus EVs</p>

      <form id="ev-calc-form" class="ev-form">
        <div class="ev-form-row">
          <label>
            Nivel
            <input type="number" id="ev-level" min="1" max="100" value="50" required />
          </label>
          <label>
            Naturaleza
            <select id="ev-nature" required>
              <option value="">Cargando...</option>
            </select>
          </label>
        </div>

        <div class="ev-stats-input">
          <h4>Stats actuales (del juego)</h4>
          <div class="ev-stats-grid">
            <label>HP <input type="number" id="stat-hp" min="1" max="999" required /></label>
            <label>Ataque <input type="number" id="stat-attack" min="1" max="999" required /></label>
            <label>Defensa <input type="number" id="stat-defense" min="1" max="999" required /></label>
            <label>At. Esp. <input type="number" id="stat-spAttack" min="1" max="999" required /></label>
            <label>Def. Esp. <input type="number" id="stat-spDefense" min="1" max="999" required /></label>
            <label>Velocidad <input type="number" id="stat-speed" min="1" max="999" required /></label>
          </div>
        </div>

        <details class="ev-ivs-section">
          <summary>IVs conocidos (opcional, asume 31 si no se especifica)</summary>
          <div class="ev-stats-grid">
            <label>HP <input type="number" id="iv-hp" min="0" max="31" placeholder="31" /></label>
            <label>Ataque <input type="number" id="iv-attack" min="0" max="31" placeholder="31" /></label>
            <label>Defensa <input type="number" id="iv-defense" min="0" max="31" placeholder="31" /></label>
            <label>At. Esp. <input type="number" id="iv-spAttack" min="0" max="31" placeholder="31" /></label>
            <label>Def. Esp. <input type="number" id="iv-spDefense" min="0" max="31" placeholder="31" /></label>
            <label>Velocidad <input type="number" id="iv-speed" min="0" max="31" placeholder="31" /></label>
          </div>
        </details>

        <button type="submit" class="btn btn-primary">Calcular EVs</button>
      </form>

      <div id="ev-results" class="ev-results hidden"></div>
    </div>
  `;
}

export async function initEVCalculator(pokemon: Pokemon): Promise<void> {
  const natures = await loadNatures();
  const natureSelect = document.getElementById("ev-nature") as HTMLSelectElement;

  if (natureSelect) {
    natureSelect.innerHTML = natures
      .map((n) => {
        const effect =
          n.increasedStat && n.decreasedStat
            ? ` (+${STAT_NAMES[n.increasedStat] || n.increasedStat}, -${STAT_NAMES[n.decreasedStat] || n.decreasedStat})`
            : " (Neutral)";
        return `<option value="${n.name}">${n.name}${effect}</option>`;
      })
      .join("");
  }

  const form = document.getElementById("ev-calc-form") as HTMLFormElement;
  form?.addEventListener("submit", async (e) => {
    e.preventDefault();
    await calculateAndShowEVs(pokemon.Name);
  });
}

async function calculateAndShowEVs(pokemonName: string): Promise<void> {
  const resultsEl = document.getElementById("ev-results") as HTMLDivElement;
  resultsEl.classList.remove("hidden");
  resultsEl.innerHTML = '<p class="loading">Calculando...</p>';

  try {
    const level = parseInt((document.getElementById("ev-level") as HTMLInputElement).value);
    const natureName = (document.getElementById("ev-nature") as HTMLSelectElement).value;

    const currentStats = new core.Stats({
      hp: parseInt((document.getElementById("stat-hp") as HTMLInputElement).value),
      attack: parseInt((document.getElementById("stat-attack") as HTMLInputElement).value),
      defense: parseInt((document.getElementById("stat-defense") as HTMLInputElement).value),
      spAttack: parseInt((document.getElementById("stat-spAttack") as HTMLInputElement).value),
      spDefense: parseInt((document.getElementById("stat-spDefense") as HTMLInputElement).value),
      speed: parseInt((document.getElementById("stat-speed") as HTMLInputElement).value),
    });

    // IVs opcionales
    let knownIVs: core.Stats | undefined;
    const ivHp = (document.getElementById("iv-hp") as HTMLInputElement).value;
    if (ivHp !== "") {
      knownIVs = new core.Stats({
        hp: parseInt(ivHp) || 31,
        attack: parseInt((document.getElementById("iv-attack") as HTMLInputElement).value) || 31,
        defense: parseInt((document.getElementById("iv-defense") as HTMLInputElement).value) || 31,
        spAttack: parseInt((document.getElementById("iv-spAttack") as HTMLInputElement).value) || 31,
        spDefense: parseInt((document.getElementById("iv-spDefense") as HTMLInputElement).value) || 31,
        speed: parseInt((document.getElementById("iv-speed") as HTMLInputElement).value) || 31,
      });
    }

    const input = core.EVCalculatorInput.createFrom({
      pokemonName,
      level,
      natureName,
      currentStats,
      knownIVs,
    });

    const result = await CalculateEVs(input);

    renderEVResults(result);
  } catch (err) {
    resultsEl.innerHTML = `<p class="error">Error: ${String(err)}</p>`;
  }
}

function renderEVResults(result: core.EVCalculatorResult): void {
  const resultsEl = document.getElementById("ev-results") as HTMLDivElement;

  const statKeys: (keyof core.Stats)[] = ["hp", "attack", "defense", "spAttack", "spDefense", "speed"];

  const evBars = statKeys
    .map((key) => {
      const ev = result.estimatedEVs[key];
      const range = result.evRanges[key];
      const percent = Math.min(100, (ev / 252) * 100);
      const rangeText = range && range.min !== range.max ? ` (${range.min}-${range.max})` : "";

      return `
      <div class="ev-bar-row">
        <span class="ev-stat-name">${STAT_NAMES[key]}</span>
        <div class="ev-bar-container">
          <div class="ev-bar" style="width: ${percent}%"></div>
        </div>
        <span class="ev-value">${ev}${rangeText}</span>
      </div>
    `;
    })
    .join("");

  const maxStatsList = statKeys
    .map((key) => `<span>${STAT_NAMES[key]}: ${result.maxPossibleStats[key]}</span>`)
    .join("");

  resultsEl.innerHTML = `
    <h4>EVs Estimados</h4>
    <div class="ev-bars">${evBars}</div>

    <div class="ev-summary">
      <div class="ev-summary-item">
        <strong>EVs usados:</strong> ${result.totalEVsUsed} / 510
      </div>
      <div class="ev-summary-item">
        <strong>EVs restantes:</strong> ${result.evsRemaining}
      </div>
    </div>

    <div class="ev-progress-total">
      <div class="ev-bar-container">
        <div class="ev-bar ev-bar-total" style="width: ${(result.totalEVsUsed / 510) * 100}%"></div>
      </div>
    </div>

    <details class="ev-max-stats">
      <summary>Stats maximos posibles (252 EVs)</summary>
      <div class="ev-max-stats-list">${maxStatsList}</div>
    </details>
  `;
}
