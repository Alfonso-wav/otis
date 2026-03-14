import gsap from "gsap";
import {
  GetPokemon,
  GetMove,
  GetNatures,
  ListPokemon,
  CalculateStats,
  SimulateDamage,
} from "../../wailsjs/go/app/App";
import type { core } from "../../wailsjs/go/models";
import { createAutocomplete } from "../autocomplete";

// ─── State ────────────────────────────────────────────────────────────────────

interface BuildSlot {
  moveName: string | null;
  move: core.Move | null;
  isCritical: boolean;
}

interface BuildState {
  attacker: core.Pokemon | null;
  attackerStats: core.Stats | null;
  attackerLevel: number;
  attackerNature: string;
  attackerIVs: core.Stats;
  attackerEVs: core.Stats;

  defender: core.Pokemon | null;
  defenderStats: core.Stats | null;
  defenderLevel: number;
  defenderNature: string;
  defenderIVs: core.Stats;
  defenderEVs: core.Stats;

  slots: [BuildSlot, BuildSlot, BuildSlot, BuildSlot];
}

const defaultStats = (): core.Stats => ({
  hp: 0, attack: 0, defense: 0, spAttack: 0, spDefense: 0, speed: 0,
});
const defaultIVs = (): core.Stats => ({
  hp: 31, attack: 31, defense: 31, spAttack: 31, spDefense: 31, speed: 31,
});
const emptySlot = (): BuildSlot => ({ moveName: null, move: null, isCritical: false });

let state: BuildState = {
  attacker: null,
  attackerStats: null,
  attackerLevel: 50,
  attackerNature: "Hardy",
  attackerIVs: defaultIVs(),
  attackerEVs: defaultStats(),

  defender: null,
  defenderStats: null,
  defenderLevel: 50,
  defenderNature: "Hardy",
  defenderIVs: defaultIVs(),
  defenderEVs: defaultStats(),

  slots: [emptySlot(), emptySlot(), emptySlot(), emptySlot()],
};

let natures: core.Nature[] = [];
let pokemonNames: string[] = [];
let container: HTMLElement;

// ─── Helpers ──────────────────────────────────────────────────────────────────

function spriteURL(id: number): string {
  return `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${id}.png`;
}

function typeBadge(type: string): string {
  return `<span class="type-badge type-${type}">${type}</span>`;
}

function typeBadges(types: Array<{ Name: string }>): string {
  return types.map((t) => typeBadge(t.Name)).join(" ");
}

function categoryIcon(cat: string): string {
  const icons: Record<string, string> = {
    physical: "⚔️",
    special: "✨",
    status: "🛡️",
  };
  return icons[cat] ?? cat;
}

function effectLabel(result: core.DamageResult): string {
  if (result.hasNoEffect) return "Sin efecto";
  if (result.isSuperEffective) return `¡Super eficaz! ×${result.multiplier}`;
  if (result.isNotVeryEffective) return `Poco eficaz ×${result.multiplier}`;
  return `×${result.multiplier}`;
}

function effectClass(result: core.DamageResult): string {
  if (result.hasNoEffect) return "damage-row--no-effect";
  if (result.isSuperEffective) return "damage-row--super";
  if (result.isNotVeryEffective) return "damage-row--weak";
  return "";
}

function statsFromPokemon(pokemon: core.Pokemon): core.Stats {
  const s = defaultStats();
  for (const stat of pokemon.Stats) {
    switch (stat.Name) {
      case "hp": s.hp = stat.BaseStat; break;
      case "attack": s.attack = stat.BaseStat; break;
      case "defense": s.defense = stat.BaseStat; break;
      case "special-attack": s.spAttack = stat.BaseStat; break;
      case "special-defense": s.spDefense = stat.BaseStat; break;
      case "speed": s.speed = stat.BaseStat; break;
    }
  }
  return s;
}

function totalEVs(evs: core.Stats): number {
  return evs.hp + evs.attack + evs.defense + evs.spAttack + evs.spDefense + evs.speed;
}

// ─── Render ───────────────────────────────────────────────────────────────────

function renderPokemonCard(pokemon: core.Pokemon, stats: core.Stats | null): string {
  const bst = statsFromPokemon(pokemon);
  const displayStats = stats ?? bst;
  return `
    <div class="build-poke-card">
      <img class="build-sprite" src="${spriteURL(pokemon.ID)}" alt="${pokemon.Name}" />
      <div class="build-poke-name">${pokemon.Name}</div>
      <div class="build-poke-types">${typeBadges(pokemon.Types)}</div>
      <div class="build-stats-grid">
        <span class="build-stat-label">HP</span><span class="build-stat-val">${displayStats.hp}</span>
        <span class="build-stat-label">Atk</span><span class="build-stat-val">${displayStats.attack}</span>
        <span class="build-stat-label">Def</span><span class="build-stat-val">${displayStats.defense}</span>
        <span class="build-stat-label">SpA</span><span class="build-stat-val">${displayStats.spAttack}</span>
        <span class="build-stat-label">SpD</span><span class="build-stat-val">${displayStats.spDefense}</span>
        <span class="build-stat-label">Spe</span><span class="build-stat-val">${displayStats.speed}</span>
      </div>
    </div>`;
}

function renderMoveSlots(moves: core.PokemonMoveEntry[]): string {
  const moveOptions = [...moves]
    .sort((a, b) => a.Name.localeCompare(b.Name))
    .map(
      (m) =>
        `<option value="${m.Name}">${m.Name} (${m.Method}${m.Level > 0 ? " lv." + m.Level : ""})</option>`,
    )
    .join("");

  return state.slots
    .map(
      (slot, i) => `
    <div class="build-slot ${slot.move ? "type-" + slot.move.Type : ""}" data-slot="${i}">
      <div class="build-slot-header">
        <span class="build-slot-num">Slot ${i + 1}</span>
        ${slot.move ? `<button class="build-slot-clear" data-slot="${i}">✕</button>` : ""}
      </div>
      ${
        slot.move
          ? `<div class="build-move-info">
          <div class="build-move-name">${slot.move.Name}</div>
          <div class="build-move-meta">
            ${typeBadge(slot.move.Type)}
            <span class="build-move-cat">${categoryIcon(slot.move.Category)} ${slot.move.Category}</span>
            <span class="build-move-power">Pwr: ${slot.move.Power || "—"}</span>
            <span class="build-move-acc">Acc: ${slot.move.Accuracy || "—"}</span>
          </div>
          <label class="build-crit-toggle">
            <input type="checkbox" class="build-crit-cb" data-slot="${i}" ${slot.isCritical ? "checked" : ""} />
            Crítico
          </label>
        </div>`
          : `<select class="build-move-select" data-slot="${i}">
          <option value="">— Elegir movimiento —</option>
          ${moveOptions}
        </select>`
      }
    </div>`,
    )
    .join("");
}

function renderStatsConfig(prefix: "atk" | "def", level: number, nature: string): string {
  const natOptions = natures
    .map(
      (n) =>
        `<option value="${n.name}" ${n.name === nature ? "selected" : ""}>${n.name}${n.increasedStat ? ` (+${n.increasedStat})` : ""}</option>`,
    )
    .join("");

  const statKeys: Array<{ key: keyof core.Stats; label: string }> = [
    { key: "hp", label: "HP" },
    { key: "attack", label: "Atk" },
    { key: "defense", label: "Def" },
    { key: "spAttack", label: "SpA" },
    { key: "spDefense", label: "SpD" },
    { key: "speed", label: "Spe" },
  ];

  const ivState = prefix === "atk" ? state.attackerIVs : state.defenderIVs;
  const evState = prefix === "atk" ? state.attackerEVs : state.defenderEVs;

  const statRows = statKeys
    .map(
      ({ key, label }) => `
    <tr>
      <td class="sc-label">${label}</td>
      <td><input class="sc-input" type="number" min="0" max="31" data-prefix="${prefix}" data-stat="${key}" data-field="iv" value="${ivState[key]}" /></td>
      <td><input class="sc-input" type="number" min="0" max="252" data-prefix="${prefix}" data-stat="${key}" data-field="ev" value="${evState[key]}" /></td>
    </tr>`,
    )
    .join("");

  return `
    <div class="build-stats-config">
      <div class="build-stats-config-row">
        <label class="sc-field-label">Nivel</label>
        <input class="sc-input sc-level" type="number" min="1" max="100" data-prefix="${prefix}" data-field="level" value="${level}" />
        <label class="sc-field-label">Naturaleza</label>
        <select class="sc-select sc-nature" data-prefix="${prefix}" data-field="nature">${natOptions}</select>
      </div>
      <table class="sc-table">
        <thead><tr><th>Stat</th><th>IV</th><th>EV</th></tr></thead>
        <tbody>${statRows}</tbody>
      </table>
      <button class="build-calc-btn" data-prefix="${prefix}">Calcular stats</button>
    </div>`;
}

function renderDamageSection(): string {
  const filledSlots = state.slots.filter((s) => s.move !== null);
  if (!state.attacker || !state.defender || filledSlots.length === 0) return "";

  return `
    <div class="build-damage-section">
      <h3 class="build-section-title">Simulación de daño</h3>
      <div id="damage-table-content"><p class="loading">Calculando...</p></div>
    </div>`;
}

async function loadDamageTable(): Promise<void> {
  const el = container.querySelector<HTMLElement>("#damage-table-content");
  if (!el || !state.attacker || !state.defender) return;

  const filledSlots = state.slots.filter((s) => s.move !== null);
  if (filledSlots.length === 0) return;

  const attackerStats = state.attackerStats ?? statsFromPokemon(state.attacker);
  const defenderStats = state.defenderStats ?? statsFromPokemon(state.defender);

  const results = await Promise.all(
    filledSlots.map(async (slot) => {
      const result = await SimulateDamage({
        attackerStats,
        defenderStats,
        move: slot.move!,
        attackerTypes: state.attacker!.Types,
        defenderTypes: state.defender!.Types,
        level: state.attackerLevel,
        isCritical: slot.isCritical,
        weatherBonus: 1.0,
      } as core.DamageInput);
      return { slot, result };
    }),
  );

  const rows = results
    .map(
      ({ slot, result }) => `
    <tr class="damage-row ${effectClass(result)}">
      <td class="dmg-move">${slot.move!.Name}</td>
      <td>${typeBadge(slot.move!.Type)}</td>
      <td class="dmg-cat">${categoryIcon(slot.move!.Category)}</td>
      <td class="dmg-val">${result.min}</td>
      <td class="dmg-val">${result.max}</td>
      <td class="dmg-eff">${effectLabel(result)}</td>
    </tr>`,
    )
    .join("");

  el.innerHTML = `
    <table class="damage-table">
      <thead>
        <tr>
          <th>Movimiento</th><th>Tipo</th><th>Cat.</th>
          <th>Mín</th><th>Máx</th><th>Efectividad</th>
        </tr>
      </thead>
      <tbody>${rows}</tbody>
    </table>`;
}

// ─── Layout ───────────────────────────────────────────────────────────────────

function buildLayout(): void {
  const atkCard = state.attacker
    ? renderPokemonCard(state.attacker, state.attackerStats)
    : `<div class="build-poke-card build-poke-card--empty"><p>Selecciona un Pokémon atacante</p></div>`;

  const defCard = state.defender
    ? renderPokemonCard(state.defender, state.defenderStats)
    : `<div class="build-poke-card build-poke-card--empty"><p>Selecciona un Pokémon defensor</p></div>`;

  const atkConfig = state.attacker
    ? renderStatsConfig("atk", state.attackerLevel, state.attackerNature)
    : "";

  const defConfig = state.defender
    ? renderStatsConfig("def", state.defenderLevel, state.defenderNature)
    : "";

  const movesSection = state.attacker
    ? `<div class="build-moves-section">
        <h3 class="build-section-title">Movimientos (máx. 4)</h3>
        <div class="build-slots">${renderMoveSlots(state.attacker.Moves ?? [])}</div>
      </div>`
    : "";

  const dmgSection = renderDamageSection();

  container.innerHTML = `
    <div class="section-header"><h2>Builds & Simulador</h2></div>

    <div class="build-layout">
      <div class="build-col build-col--attacker">
        <h3 class="build-col-title">Atacante</h3>
        <div class="build-search-row">
          <input id="atk-input" class="build-search-input" type="text" placeholder="Nombre del Pokémon..." />
          <button id="atk-btn" class="build-search-btn">Buscar</button>
        </div>
        ${atkCard}
        ${atkConfig}
      </div>

      <div class="build-col build-col--defender">
        <h3 class="build-col-title">Defensor</h3>
        <div class="build-search-row">
          <input id="def-input" class="build-search-input" type="text" placeholder="Nombre del Pokémon..." />
          <button id="def-btn" class="build-search-btn">Buscar</button>
        </div>
        ${defCard}
        ${defConfig}
      </div>
    </div>

    ${movesSection}
    ${dmgSection}
  `;

  bindEvents();

  if (dmgSection) {
    loadDamageTable();
  }
}

// ─── Event binding ────────────────────────────────────────────────────────────

function bindEvents(): void {
  const atkBtn = container.querySelector<HTMLButtonElement>("#atk-btn");
  const atkInput = container.querySelector<HTMLInputElement>("#atk-input");
  atkBtn?.addEventListener("click", () => fetchPokemon("atk", atkInput?.value ?? ""));
  atkInput?.addEventListener("keydown", (e) => {
    if (e.key === "Enter") fetchPokemon("atk", atkInput.value);
  });

  const defBtn = container.querySelector<HTMLButtonElement>("#def-btn");
  const defInput = container.querySelector<HTMLInputElement>("#def-input");
  defBtn?.addEventListener("click", () => fetchPokemon("def", defInput?.value ?? ""));
  defInput?.addEventListener("keydown", (e) => {
    if (e.key === "Enter") fetchPokemon("def", defInput.value);
  });

  if (pokemonNames.length > 0) {
    if (atkInput) createAutocomplete(atkInput, pokemonNames, (name) => fetchPokemon("atk", name));
    if (defInput) createAutocomplete(defInput, pokemonNames, (name) => fetchPokemon("def", name));
  }

  container.querySelectorAll<HTMLSelectElement>(".build-move-select").forEach((sel) => {
    sel.addEventListener("change", () => {
      const idx = parseInt(sel.dataset.slot ?? "0");
      if (sel.value) loadMove(idx, sel.value);
    });
  });

  container.querySelectorAll<HTMLButtonElement>(".build-slot-clear").forEach((btn) => {
    btn.addEventListener("click", () => {
      const idx = parseInt(btn.dataset.slot ?? "0");
      state.slots[idx] = emptySlot();
      buildLayout();
    });
  });

  container.querySelectorAll<HTMLInputElement>(".build-crit-cb").forEach((cb) => {
    cb.addEventListener("change", () => {
      const idx = parseInt(cb.dataset.slot ?? "0");
      state.slots[idx].isCritical = cb.checked;
      loadDamageTable();
    });
  });

  container.querySelectorAll<HTMLInputElement>(".sc-input").forEach((inp) => {
    inp.addEventListener("change", () => handleStatInput(inp));
  });

  container.querySelectorAll<HTMLSelectElement>(".sc-select").forEach((sel) => {
    sel.addEventListener("change", () => handleNatureSelect(sel));
  });

  container.querySelectorAll<HTMLButtonElement>(".build-calc-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const prefix = btn.dataset.prefix as "atk" | "def";
      calcStats(prefix);
    });
  });
}

function handleStatInput(inp: HTMLInputElement): void {
  const prefix = inp.dataset.prefix as "atk" | "def";
  const stat = inp.dataset.stat as keyof core.Stats;
  const field = inp.dataset.field as "iv" | "ev" | "level";
  const val = parseInt(inp.value) || 0;

  if (field === "level") {
    if (prefix === "atk") state.attackerLevel = Math.min(100, Math.max(1, val));
    else state.defenderLevel = Math.min(100, Math.max(1, val));
    return;
  }

  const evs = prefix === "atk" ? state.attackerEVs : state.defenderEVs;
  const ivs = prefix === "atk" ? state.attackerIVs : state.defenderIVs;

  if (field === "ev") {
    const clamped = Math.min(252, Math.max(0, val));
    const without = totalEVs(evs) - (evs[stat] as number);
    evs[stat] = Math.min(clamped, 510 - without) as never;
    inp.value = String(evs[stat]);
  } else {
    ivs[stat] = Math.min(31, Math.max(0, val)) as never;
  }
}

function handleNatureSelect(sel: HTMLSelectElement): void {
  const prefix = sel.dataset.prefix as "atk" | "def";
  if (prefix === "atk") state.attackerNature = sel.value;
  else state.defenderNature = sel.value;
}

// ─── Async actions ────────────────────────────────────────────────────────────

async function fetchPokemon(prefix: "atk" | "def", name: string): Promise<void> {
  if (!name.trim()) return;

  const btn = container.querySelector<HTMLButtonElement>(`#${prefix}-btn`);
  if (btn) btn.disabled = true;

  try {
    const pokemon = await GetPokemon(name.trim().toLowerCase());

    if (prefix === "atk") {
      state.attacker = pokemon;
      state.attackerStats = null;
      state.slots = [emptySlot(), emptySlot(), emptySlot(), emptySlot()];
    } else {
      state.defender = pokemon;
      state.defenderStats = null;
    }

    buildLayout();
    gsap.fromTo(
      container.querySelector(prefix === "atk" ? ".build-col--attacker" : ".build-col--defender"),
      { opacity: 0, y: 10 },
      { opacity: 1, y: 0, duration: 0.3, ease: "power2.out" },
    );
  } catch (err: unknown) {
    alert(`Error: ${String(err)}`);
    if (btn) btn.disabled = false;
  }
}

async function loadMove(slotIdx: number, moveName: string): Promise<void> {
  try {
    const move = await GetMove(moveName);
    state.slots[slotIdx] = { moveName, move, isCritical: false };
    buildLayout();
  } catch (err: unknown) {
    alert(`Error al cargar movimiento: ${String(err)}`);
  }
}

async function calcStats(prefix: "atk" | "def"): Promise<void> {
  const pokemon = prefix === "atk" ? state.attacker : state.defender;
  if (!pokemon) return;

  const level = prefix === "atk" ? state.attackerLevel : state.defenderLevel;
  const nature = prefix === "atk" ? state.attackerNature : state.defenderNature;
  const ivs = prefix === "atk" ? state.attackerIVs : state.defenderIVs;
  const evs = prefix === "atk" ? state.attackerEVs : state.defenderEVs;

  try {
    const stats = await CalculateStats({
      pokemonName: pokemon.Name,
      level,
      natureName: nature,
      ivs,
      evs,
    } as core.StatCalculatorInput);

    if (prefix === "atk") state.attackerStats = stats;
    else state.defenderStats = stats;

    buildLayout();
  } catch (err: unknown) {
    alert(`Error al calcular stats: ${String(err)}`);
  }
}

// ─── Init ──────────────────────────────────────────────────────────────────────

let initialized = false;

export async function initBuilds(): Promise<void> {
  container = document.getElementById("tab-builds") as HTMLElement;
  if (!container) return;

  const tabBtn = document.querySelector<HTMLButtonElement>('[data-tab="builds"]');
  if (!tabBtn) return;

  tabBtn.addEventListener("click", async () => {
    if (initialized) return;
    initialized = true;

    try {
      [natures] = await Promise.all([
        GetNatures(),
        ListPokemon(0, 2000).then((resp) => { pokemonNames = resp.Results.map((r) => r.Name); }),
      ]);
    } catch {
      natures = [];
    }

    buildLayout();
    gsap.fromTo(
      container.querySelector(".section-header"),
      { opacity: 0, y: -8 },
      { opacity: 1, y: 0, duration: 0.3, ease: "power2.out" },
    );
  });
}
