import gsap from "gsap";
import { ComparePokemons } from "../../wailsjs/go/app/App";
import type { core } from "../../wailsjs/go/models";

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

function spriteURL(id: number): string {
  return `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${id}.png`;
}

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
          <img class="compare-sprite" src="${spriteURL(a.ID)}" alt="${a.Name}" />
          <div class="compare-pokemon-name">${a.Name}</div>
          <div class="compare-pokemon-types">${typeBadges(a.Types)}</div>
        </div>
        <div class="compare-winner ${winnerClass}">
          ${Winner === "tie" ? "⚖️" : "🏆"} ${winnerLabel}
        </div>
        <div class="compare-pokemon-card">
          <img class="compare-sprite" src="${spriteURL(b.ID)}" alt="${b.Name}" />
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
}

let initialized = false;

export function initCompare(): void {
  container = document.getElementById("tab-compare") as HTMLElement;
  if (!container) return;

  const tabBtn = document.querySelector<HTMLButtonElement>('[data-tab="compare"]');
  if (!tabBtn) return;

  tabBtn.addEventListener("click", () => {
    if (initialized) return;
    initialized = true;
    buildInitialLayout();
    gsap.fromTo(
      container.querySelector(".section-header"),
      { opacity: 0, y: -8 },
      { opacity: 1, y: 0, duration: 0.3, ease: "power2.out" },
    );
  });
}
