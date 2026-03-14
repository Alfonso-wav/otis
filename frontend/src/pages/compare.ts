import gsap from "gsap";
import { ComparePokemons, ListPokemon } from "../../wailsjs/go/app/App";
import type { core } from "../../wailsjs/go/models";
import { createAutocomplete } from "../autocomplete";

const MAX_STAT = 255;

const STAT_LABELS: Record<string, string> = {
  hp: "HP",
  attack: "Ataque",
  defense: "Defensa",
  "special-attack": "Sp. Atk",
  "special-defense": "Sp. Def",
  speed: "Velocidad",
};

let container: HTMLElement;
let pokemonNames: string[] = [];

function spriteURL(name: string): string {
  const safeName = name.toLowerCase().replace(/[^a-z0-9-]/g, "");
  return `https://img.pokemondb.net/sprites/home/normal/${safeName}.png`;
}

const ICON_TROPHY = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-label="Winner"><path stroke-linecap="round" stroke-linejoin="round" d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M7.73 9.728a6.726 6.726 0 0 0 2.748 1.35m8.272-6.842V4.5c0 2.108-.966 3.99-2.48 5.228m2.48-5.492a46.32 46.32 0 0 1 2.916.52 6.003 6.003 0 0 1-5.395 4.972m0 0a6.726 6.726 0 0 1-2.749 1.35m0 0a6.772 6.772 0 0 1-3.044 0" /></svg>`;

const ICON_SCALE = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-label="Tie"><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v17.25m0 0c-1.472 0-2.882.265-4.185.75M12 20.25c1.472 0 2.882.265 4.185.75M18.75 4.97A48.416 48.416 0 0 0 12 4.5c-2.291 0-4.545.16-6.75.47m13.5 0c1.01.143 2.01.317 3 .52m-3-.52 2.62 10.726c.122.499-.106 1.028-.589 1.202a5.988 5.988 0 0 1-2.031.352 5.988 5.988 0 0 1-2.031-.352c-.483-.174-.711-.703-.59-1.202L18.75 4.971Zm-16.5.52c.99-.203 1.99-.377 3-.52m0 0 2.62 10.726c.122.499-.106 1.028-.589 1.202a5.989 5.989 0 0 1-2.031.352 5.989 5.989 0 0 1-2.031-.352c-.483-.174-.711-.703-.59-1.202L5.25 4.971Z" /></svg>`;

function typeBadges(types: Array<{ Name: string }>): string {
  return types
    .map((t) => `<span class="type-badge type-${t.Name}">${t.Name}</span>`)
    .join(" ");
}

function statBar(value: number, winner: "a" | "b" | "tie", side: "a" | "b"): string {
  const pct = Math.round((value / MAX_STAT) * 100);
  let colorClass = "compare-bar--tie";
  if (winner === side) colorClass = "compare-bar--win";
  else if (winner !== "tie") colorClass = "compare-bar--lose";
  return `<div class="compare-bar-track">
    <div class="compare-bar ${colorClass}" style="width:${pct}%"></div>
  </div>`;
}

function renderResult(result: core.PokemonComparison): void {
  const { PokemonA: a, PokemonB: b, Stats: stats, TotalA, TotalB, Winner } = result;

  const winnerLabel =
    Winner === "a" ? a.Name : Winner === "b" ? b.Name : "¡Empate!";
  const winnerClass =
    Winner === "tie" ? "compare-winner--tie" : "compare-winner--win";

  const rows = stats
    .map((s) => {
      const label = STAT_LABELS[s.Name] ?? s.Name;
      const diffText =
        s.Diff === 0 ? "=" : s.Diff > 0 ? `+${s.Diff}` : `${s.Diff}`;
      return `<tr class="compare-stat-row">
        <td class="compare-stat-val compare-stat-val--a ${s.Winner === "a" ? "compare-stat-val--win" : s.Winner === "b" ? "compare-stat-val--lose" : ""}">
          ${s.StatA}
          ${statBar(s.StatA, s.Winner as "a" | "b" | "tie", "a")}
        </td>
        <td class="compare-stat-name">${label}</td>
        <td class="compare-stat-val compare-stat-val--b ${s.Winner === "b" ? "compare-stat-val--win" : s.Winner === "a" ? "compare-stat-val--lose" : ""}">
          ${statBar(s.StatB, s.Winner as "a" | "b" | "tie", "b")}
          ${s.StatB}
        </td>
        <td class="compare-diff ${s.Winner === "a" ? "compare-diff--a" : s.Winner === "b" ? "compare-diff--b" : ""}">${diffText}</td>
      </tr>`;
    })
    .join("");

  container.innerHTML = `
    <div class="section-header"><h2>Comparar Pokémon</h2></div>
    <div class="compare-inputs">
      <input id="compare-input-a" class="compare-input" type="text" placeholder="Pokémon A (nombre o #)..." value="${a.Name}" />
      <span class="compare-vs">VS</span>
      <input id="compare-input-b" class="compare-input" type="text" placeholder="Pokémon B (nombre o #)..." value="${b.Name}" />
      <button id="compare-btn" class="compare-btn">Comparar</button>
    </div>
    <div id="compare-result" class="compare-result">
      <div class="compare-header">
        <div class="compare-pokemon-card">
          <img class="compare-sprite" src="${spriteURL(a.Name)}" onerror="this.onerror=null;this.src='https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${a.ID}.png'" alt="${a.Name}" />
          <div class="compare-pokemon-name">${a.Name}</div>
          <div class="compare-pokemon-types">${typeBadges(a.Types)}</div>
        </div>
        <div class="compare-winner ${winnerClass}">
          ${Winner === "tie" ? ICON_SCALE : ICON_TROPHY} ${winnerLabel}
        </div>
        <div class="compare-pokemon-card">
          <img class="compare-sprite" src="${spriteURL(b.Name)}" onerror="this.onerror=null;this.src='https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${b.ID}.png'" alt="${b.Name}" />
          <div class="compare-pokemon-name">${b.Name}</div>
          <div class="compare-pokemon-types">${typeBadges(b.Types)}</div>
        </div>
      </div>
      <table class="compare-table">
        <tbody>
          ${rows}
          <tr class="compare-total-row">
            <td class="compare-stat-val ${TotalA > TotalB ? "compare-stat-val--win" : TotalA < TotalB ? "compare-stat-val--lose" : ""}">${TotalA}</td>
            <td class="compare-stat-name compare-total-label">Total BST</td>
            <td class="compare-stat-val ${TotalB > TotalA ? "compare-stat-val--win" : TotalB < TotalA ? "compare-stat-val--lose" : ""}">${TotalB}</td>
            <td class="compare-diff"></td>
          </tr>
        </tbody>
      </table>
    </div>`;

  bindInputs();

  const result2 = container.querySelector<HTMLElement>("#compare-result")!;
  gsap.fromTo(
    result2,
    { opacity: 0, y: 16 },
    { opacity: 1, y: 0, duration: 0.35, ease: "power2.out" },
  );
}

function buildInitialLayout(): void {
  container.innerHTML = `
    <div class="section-header"><h2>Comparar Pokémon</h2></div>
    <div class="compare-inputs">
      <input id="compare-input-a" class="compare-input" type="text" placeholder="Pokémon A (nombre o #)..." />
      <span class="compare-vs">VS</span>
      <input id="compare-input-b" class="compare-input" type="text" placeholder="Pokémon B (nombre o #)..." />
      <button id="compare-btn" class="compare-btn">Comparar</button>
    </div>
    <div id="compare-result" class="compare-result compare-result--empty">
      <p class="loading">Introduce dos Pokémon y pulsa Comparar</p>
    </div>`;

  bindInputs();
}

async function runComparison(nameA: string, nameB: string): Promise<void> {
  const resultEl = container.querySelector<HTMLElement>("#compare-result")!;
  resultEl.innerHTML = '<p class="loading">Comparando...</p>';
  resultEl.classList.remove("compare-result--empty");

  try {
    const data = await ComparePokemons(nameA.trim().toLowerCase(), nameB.trim().toLowerCase());
    renderResult(data);
  } catch (err: unknown) {
    resultEl.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}

function bindInputs(): void {
  const btn = container.querySelector<HTMLButtonElement>("#compare-btn")!;
  const inputA = container.querySelector<HTMLInputElement>("#compare-input-a")!;
  const inputB = container.querySelector<HTMLInputElement>("#compare-input-b")!;

  const onCompare = () => {
    if (!inputA.value.trim() || !inputB.value.trim()) return;
    runComparison(inputA.value, inputB.value);
  };

  btn.addEventListener("click", onCompare);
  [inputA, inputB].forEach((inp) => {
    inp.addEventListener("keydown", (e) => {
      if (e.key === "Enter") onCompare();
    });
  });

  if (pokemonNames.length > 0) {
    createAutocomplete(inputA, pokemonNames, (name) => {
      inputA.value = name;
    });
    createAutocomplete(inputB, pokemonNames, (name) => {
      inputB.value = name;
    });
  }
}

let initialized = false;

export function initCompare(): void {
  container = document.getElementById("tab-compare") as HTMLElement;
  if (!container) return;

  const tabBtn = document.querySelector<HTMLButtonElement>('[data-tab="compare"]');
  if (!tabBtn) return;

  tabBtn.addEventListener("click", async () => {
    if (initialized) return;
    initialized = true;

    try {
      const resp = await ListPokemon(0, 2000);
      pokemonNames = resp.Results.map((r) => r.Name);
    } catch {
      pokemonNames = [];
    }

    buildInitialLayout();
    gsap.fromTo(
      container.querySelector(".section-header"),
      { opacity: 0, y: -8 },
      { opacity: 1, y: 0, duration: 0.3, ease: "power2.out" },
    );
  });
}
