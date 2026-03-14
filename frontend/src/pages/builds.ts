import gsap from "gsap";
import {
  GetPokemon,
  GetMove,
  GetNatures,
  ListPokemon,
  CalculateStats,
  SimulateDamage,
  InitBattle,
  ExecuteTurn,
  SimulateFullBattle,
  SaveToTeam,
  ListTeams,
  DeleteTeam,
  DeleteTeamMember,
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
  defenderSlots: [BuildSlot, BuildSlot, BuildSlot, BuildSlot];
}

const defaultStats = (): core.Stats => ({
  hp: 0, attack: 0, defense: 0, spAttack: 0, spDefense: 0, speed: 0,
});
const defaultIVs = (): core.Stats => ({
  hp: 31, attack: 31, defense: 31, spAttack: 31, spDefense: 31, speed: 31,
});
const emptySlot = (): BuildSlot => ({ moveName: null, move: null, isCritical: false });

type BattlePhase = "idle" | "attacker-turn" | "over";

interface BattleUIState {
  battleState: core.BattleState | null;
  phase: BattlePhase;
}

let battleUI: BattleUIState = { battleState: null, phase: "idle" };

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
  defenderSlots: [emptySlot(), emptySlot(), emptySlot(), emptySlot()],
};

let natures: core.Nature[] = [];
let pokemonNames: string[] = [];
let container: HTMLElement;
let cachedTeams: core.Team[] = [];

// ─── Helpers ──────────────────────────────────────────────────────────────────

function spriteURL(name: string): string {
  const safeName = name.toLowerCase().replace(/[^a-z0-9-]/g, "");
  return `https://img.pokemondb.net/sprites/home/normal/${safeName}.png`;
}

function typeBadge(type: string): string {
  return `<span class="type-badge type-${type}">${type}</span>`;
}

function typeBadges(types: Array<{ Name: string }>): string {
  return types.map((t) => typeBadge(t.Name)).join(" ");
}

function categoryIcon(cat: string): string {
  const base = "https://img.pokemondb.net/images/icons";
  const map: Record<string, string> = {
    physical: `<img src="${base}/move-physical.png" class="move-cat-icon" alt="Physical" title="Physical">`,
    special:  `<img src="${base}/move-special.png"  class="move-cat-icon" alt="Special"  title="Special">`,
    status:   `<img src="${base}/move-status.png"   class="move-cat-icon" alt="Status"   title="Status">`,
  };
  return map[cat] ?? `<span class="move-cat-unknown">?</span>`;
}

function effectLabel(result: core.DamageResult): string {
  if (result.hasNoEffect) return "Sin efecto";
  const combined = result.multiplier * (result.hasSTAB ? result.stabMultiplier : 1);
  if (result.isSuperEffective) return `¡Super eficaz! ×${combined}`;
  if (result.isNotVeryEffective) return `Poco eficaz ×${combined}`;
  return `×${combined}`;
}

function stabBadge(result: core.DamageResult): string {
  if (result.hasSTAB) return `<span class="stab-badge">STAB</span>`;
  return "";
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
      <img class="build-sprite" src="${spriteURL(pokemon.Name)}" onerror="this.onerror=null;this.src='https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${pokemon.ID}.png'" alt="${pokemon.Name}" />
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

function renderMoveSlots(moves: core.PokemonMoveEntry[], prefix: "atk" | "def"): string {
  const slots = prefix === "atk" ? state.slots : state.defenderSlots;
  const moveOptions = [...moves]
    .sort((a, b) => a.Name.localeCompare(b.Name))
    .map(
      (m) =>
        `<option value="${m.Name}">${m.Name} (${m.Method}${m.Level > 0 ? " lv." + m.Level : ""})</option>`,
    )
    .join("");

  return slots
    .map(
      (slot, i) => `
    <div class="build-slot ${slot.move ? "type-" + slot.move.Type : ""}" data-slot="${i}" data-prefix="${prefix}">
      <div class="build-slot-header">
        <span class="build-slot-num">Slot ${i + 1}</span>
        ${slot.move ? `<button class="build-slot-clear" data-slot="${i}" data-prefix="${prefix}">✕</button>` : ""}
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
          ${prefix === "atk" ? `<label class="build-crit-toggle">
            <input type="checkbox" class="build-crit-cb" data-slot="${i}" ${slot.isCritical ? "checked" : ""} />
            Crítico
          </label>` : ""}
        </div>`
          : `<select class="build-move-select" data-slot="${i}" data-prefix="${prefix}">
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
      <td class="dmg-move">${slot.move!.Name} ${stabBadge(result)}</td>
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

// ─── Battle ───────────────────────────────────────────────────────────────────

function hpBarColor(pct: number): string {
  if (pct > 0.5) return "#48bb78";
  if (pct > 0.25) return "#ed8936";
  return "#e53e3e";
}

function renderHPBar(name: string, current: number, max: number): string {
  const pct = max > 0 ? Math.max(0, current / max) : 0;
  const color = hpBarColor(pct);
  return `
    <div class="battle-pokemon-hp">
      <div class="battle-pokemon-name">${name}</div>
      <div class="hp-bar-track">
        <div class="hp-bar-fill" style="width:${(pct * 100).toFixed(1)}%;background:${color}"></div>
      </div>
      <div class="hp-bar-label">${current} / ${max}</div>
    </div>`;
}

function renderBattleSection(): string {
  const filledSlots = state.slots.filter((s) => s.move !== null);
  if (!state.attacker || !state.defender || filledSlots.length === 0) return "";

  const filledDefSlots = state.defenderSlots.filter((s) => s.move !== null);
  const canAutoSimulate = filledSlots.length > 0 && filledDefSlots.length > 0;

  const bs = battleUI.battleState;
  const phase = battleUI.phase;

  if (phase === "idle" || !bs) {
    return `
      <div class="battle-section" id="battle-section">
        <h3 class="build-section-title">Simulación de batalla</h3>
        <div class="battle-idle-btns">
          <button class="battle-start-btn" id="battle-start-btn">Iniciar batalla turno a turno</button>
          ${canAutoSimulate ? `<button class="battle-auto-btn" id="battle-auto-btn">Simular batalla completa</button>` : ""}
        </div>
      </div>`;
  }

  const atkName = state.attacker.Name;
  const defName = state.defender.Name;

  const winnerName = bs.winner === "attacker" ? atkName : bs.winner === "defender" ? defName : "Empate";
  const winnerBanner = bs.isOver
    ? `<div class="battle-winner-banner">${bs.winner === "draw" ? "¡Empate!" : winnerName + " gana la batalla!"}</div>`
    : "";

  const turnLabel = bs.isOver
    ? ""
    : `<div class="battle-turn-label">Turno ${bs.turnCount + 1} — Elige movimiento de <strong>${atkName}</strong></div>`;

  const moveBtns = state.slots
    .map((slot, i) => {
      if (!slot.move) return "";
      const disabled = bs.isOver ? "disabled" : "";
      const typeClass = `type-${slot.move.Type}`;
      return `<button class="move-btn ${typeClass}" data-slot="${i}" ${disabled}>
        ${slot.move.Name}
        <span class="move-btn-power">${slot.move.Power || "—"}</span>
      </button>`;
    })
    .join("");

  const logEntries = [...(bs.log ?? [])].reverse().slice(0, 20);
  const logHTML = logEntries.map((entry) => `<div class="battle-log-entry">${entry}</div>`).join("");

  return `
    <div class="battle-section" id="battle-section">
      <h3 class="build-section-title">Simulación de batalla — Turno ${bs.turnCount}</h3>
      ${winnerBanner}
      <div class="battle-hp-bars">
        ${renderHPBar(atkName, bs.attackerHP, bs.attackerMaxHP)}
        <div class="battle-vs">VS</div>
        ${renderHPBar(defName, bs.defenderHP, bs.defenderMaxHP)}
      </div>
      ${turnLabel}
      ${!bs.isOver ? `<div class="battle-move-btns" id="battle-move-btns">${moveBtns}</div>` : ""}
      <div class="battle-log" id="battle-log">${logHTML}</div>
      <div class="battle-idle-btns">
        <button class="battle-reset-btn" id="battle-reset-btn">Reiniciar</button>
        ${canAutoSimulate && !bs.isOver ? `<button class="battle-auto-btn" id="battle-auto-btn">Simular batalla completa</button>` : ""}
      </div>
    </div>`;
}

async function startBattle(): Promise<void> {
  if (!state.attacker || !state.defender) return;

  const defenderHasMoves = state.defenderSlots.some((s) => s.move !== null);
  if (!defenderHasMoves) {
    alert("El defensor necesita al menos un movimiento configurado para iniciar la batalla.");
    return;
  }

  const attackerStats = state.attackerStats ?? statsFromPokemon(state.attacker);
  const defenderStats = state.defenderStats ?? statsFromPokemon(state.defender);

  const bs = await InitBattle(attackerStats.hp, defenderStats.hp);
  battleUI = { battleState: bs, phase: "attacker-turn" };
  renderBattleInPlace();
}

async function handleMoveClick(slotIdx: number): Promise<void> {
  let bs = battleUI.battleState;
  if (!bs || bs.isOver) return;
  if (!state.attacker || !state.defender) return;

  const slot = state.slots[slotIdx];
  if (!slot.move) return;

  // ── Attacker turn ──
  const atkStats = state.attackerStats ?? statsFromPokemon(state.attacker);
  const defStats = state.defenderStats ?? statsFromPokemon(state.defender);

  const atkResult = await ExecuteTurn({
    state: bs,
    attackerStats: atkStats,
    defenderStats: defStats,
    attackerTypes: state.attacker.Types,
    defenderTypes: state.defender.Types,
    attackerLevel: state.attackerLevel,
    defenderLevel: state.defenderLevel,
    move: slot.move,
    attackerName: state.attacker.Name,
  } as core.TurnInput);

  bs = atkResult.newState;

  // If attacker KO'd the defender, battle is over
  if (bs.isOver) {
    battleUI = { battleState: bs, phase: "over" };
    renderBattleInPlace();
    return;
  }

  // ── Defender auto-response ──
  const defenderFilledSlots = state.defenderSlots.filter((s) => s.move !== null);
  const defMove = defenderFilledSlots[Math.floor(Math.random() * defenderFilledSlots.length)].move!;

  // Swap HP perspective for defender's turn
  const swappedState = {
    ...bs,
    attackerHP: bs.defenderHP,
    defenderHP: bs.attackerHP,
    attackerMaxHP: bs.defenderMaxHP,
    defenderMaxHP: bs.attackerMaxHP,
  };

  const defResult = await ExecuteTurn({
    state: swappedState,
    attackerStats: defStats,
    defenderStats: atkStats,
    attackerTypes: state.defender.Types,
    defenderTypes: state.attacker.Types,
    attackerLevel: state.defenderLevel,
    defenderLevel: state.attackerLevel,
    move: defMove,
    attackerName: state.defender.Name,
  } as core.TurnInput);

  // Restore HP perspective
  let finalBs = defResult.newState;
  finalBs = {
    ...finalBs,
    attackerHP: finalBs.defenderHP,
    defenderHP: finalBs.attackerHP,
    attackerMaxHP: finalBs.defenderMaxHP,
    defenderMaxHP: finalBs.attackerMaxHP,
    log: [...(bs.log ?? []), ...(finalBs.log ?? []).slice((bs.log ?? []).length)],
  };
  if (finalBs.isOver) {
    finalBs = { ...finalBs, winner: "defender" };
  }

  const nextPhase: BattlePhase = finalBs.isOver ? "over" : "attacker-turn";
  battleUI = { battleState: finalBs, phase: nextPhase };
  renderBattleInPlace();
}

async function simulateFullBattle(): Promise<void> {
  if (!state.attacker || !state.defender) return;

  const atkMoves = state.slots.filter((s) => s.move !== null).map((s) => s.move!);
  const defMoves = state.defenderSlots.filter((s) => s.move !== null).map((s) => s.move!);

  if (atkMoves.length === 0 || defMoves.length === 0) {
    alert("Ambos lados necesitan al menos un movimiento para simular la batalla completa.");
    return;
  }

  const attackerStats = state.attackerStats ?? statsFromPokemon(state.attacker);
  const defenderStats = state.defenderStats ?? statsFromPokemon(state.defender);

  const result = await SimulateFullBattle({
    attackerStats,
    defenderStats,
    attackerTypes: state.attacker.Types,
    defenderTypes: state.defender.Types,
    attackerLevel: state.attackerLevel,
    defenderLevel: state.defenderLevel,
    attackerMoves: atkMoves,
    defenderMoves: defMoves,
    attackerName: state.attacker.Name,
    defenderName: state.defender.Name,
  } as core.FullBattleInput);

  battleUI = { battleState: result, phase: "over" };
  renderBattleInPlace();
}

function renderBattleInPlace(): void {
  const existing = container.querySelector<HTMLElement>("#battle-section");
  if (!existing) {
    buildLayout();
    return;
  }
  existing.outerHTML = renderBattleSection();
  bindBattleEvents();
}

function bindBattleEvents(): void {
  container.querySelector<HTMLButtonElement>("#battle-start-btn")?.addEventListener("click", () => startBattle());
  container.querySelector<HTMLButtonElement>("#battle-reset-btn")?.addEventListener("click", () => startBattle());
  container.querySelector<HTMLButtonElement>("#battle-auto-btn")?.addEventListener("click", () => simulateFullBattle());
  container.querySelectorAll<HTMLButtonElement>(".move-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const idx = parseInt(btn.dataset.slot ?? "0");
      handleMoveClick(idx);
    });
  });
}

// ─── Teams ────────────────────────────────────────────────────────────────────

function buildTeamMember(prefix: "atk" | "def"): core.TeamMember | null {
  const pokemon = prefix === "atk" ? state.attacker : state.defender;
  if (!pokemon) return null;

  const slots = prefix === "atk" ? state.slots : state.defenderSlots;
  const moves = slots.filter((s) => s.move !== null).map((s) => s.move!.Name);

  return {
    pokemonName: pokemon.Name,
    moves,
    level: prefix === "atk" ? state.attackerLevel : state.defenderLevel,
    nature: prefix === "atk" ? state.attackerNature : state.defenderNature,
    ivs: prefix === "atk" ? state.attackerIVs : state.defenderIVs,
    evs: prefix === "atk" ? state.attackerEVs : state.defenderEVs,
  } as core.TeamMember;
}

async function saveToTeam(prefix: "atk" | "def"): Promise<void> {
  const member = buildTeamMember(prefix);
  if (!member) return;

  const teamName = prompt("Nombre del equipo:");
  if (!teamName || !teamName.trim()) return;

  try {
    await SaveToTeam(teamName.trim(), member);
    alert(`${member.pokemonName} guardado en "${teamName.trim()}"`);
    cachedTeams = await ListTeams();
    buildLayout();
  } catch (err: unknown) {
    alert(`Error al guardar: ${String(err)}`);
  }
}

async function importFromTeam(prefix: "atk" | "def", member: core.TeamMember): Promise<void> {
  try {
    const pokemon = await GetPokemon(member.pokemonName);

    if (prefix === "atk") {
      state.attacker = pokemon;
      state.attackerLevel = member.level;
      state.attackerNature = member.nature;
      state.attackerIVs = { ...member.ivs };
      state.attackerEVs = { ...member.evs };
      state.slots = [emptySlot(), emptySlot(), emptySlot(), emptySlot()];
    } else {
      state.defender = pokemon;
      state.defenderLevel = member.level;
      state.defenderNature = member.nature;
      state.defenderIVs = { ...member.ivs };
      state.defenderEVs = { ...member.evs };
      state.defenderSlots = [emptySlot(), emptySlot(), emptySlot(), emptySlot()];
    }

    // Load moves
    const slots = prefix === "atk" ? state.slots : state.defenderSlots;
    const movePromises = member.moves.slice(0, 4).map((name) => GetMove(name));
    const moves = await Promise.all(movePromises);
    moves.forEach((move, i) => {
      slots[i] = { moveName: move.Name, move, isCritical: false };
    });

    // Calculate stats
    const stats = await CalculateStats({
      pokemonName: member.pokemonName,
      level: member.level,
      natureName: member.nature,
      ivs: member.ivs,
      evs: member.evs,
    } as core.StatCalculatorInput);

    if (prefix === "atk") state.attackerStats = stats;
    else state.defenderStats = stats;

    battleUI = { battleState: null, phase: "idle" };
    buildLayout();
  } catch (err: unknown) {
    alert(`Error al importar: ${String(err)}`);
  }
}

async function handleDeleteTeam(name: string): Promise<void> {
  if (!confirm(`Eliminar equipo "${name}"?`)) return;
  try {
    await DeleteTeam(name);
    cachedTeams = await ListTeams();
    buildLayout();
  } catch (err: unknown) {
    alert(`Error: ${String(err)}`);
  }
}

async function handleDeleteTeamMember(teamName: string, index: number): Promise<void> {
  if (!confirm("Eliminar este miembro del equipo?")) return;
  try {
    await DeleteTeamMember(teamName, index);
    cachedTeams = await ListTeams();
    buildLayout();
  } catch (err: unknown) {
    alert(`Error: ${String(err)}`);
  }
}

function renderTeamsSection(): string {
  if (cachedTeams.length === 0) return "";

  const teamsHTML = cachedTeams.map((team) => {
    const membersHTML = team.members
      .map((m, i) => `
        <div class="team-member-row">
          <img class="team-member-sprite" src="${spriteURL(m.pokemonName)}" onerror="this.style.display='none'" alt="${m.pokemonName}" />
          <span class="team-member-name">${m.pokemonName}</span>
          <span class="team-member-detail">Lv.${m.level} ${m.nature}</span>
          <span class="team-member-moves">${m.moves.join(", ")}</span>
          <button class="team-import-btn" data-team="${team.name}" data-idx="${i}" data-prefix="atk" title="Importar como atacante">Atk</button>
          <button class="team-import-btn" data-team="${team.name}" data-idx="${i}" data-prefix="def" title="Importar como defensor">Def</button>
          <button class="team-member-delete-btn" data-team="${team.name}" data-idx="${i}" title="Eliminar miembro">✕</button>
        </div>`)
      .join("");

    return `
      <div class="team-card">
        <div class="team-card-header">
          <span class="team-card-name">${team.name} (${team.members.length}/6)</span>
          <button class="team-delete-btn" data-team="${team.name}">Eliminar equipo</button>
        </div>
        <div class="team-members-list">${membersHTML}</div>
      </div>`;
  }).join("");

  return `
    <div class="teams-section">
      <details class="teams-details">
        <summary class="build-section-title teams-summary">Mis equipos (${cachedTeams.length})</summary>
        <div class="teams-list">${teamsHTML}</div>
      </details>
    </div>`;
}

function bindTeamEvents(): void {
  container.querySelectorAll<HTMLButtonElement>(".team-save-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const prefix = btn.dataset.prefix as "atk" | "def";
      saveToTeam(prefix);
    });
  });

  container.querySelectorAll<HTMLButtonElement>(".team-import-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const teamName = btn.dataset.team!;
      const idx = parseInt(btn.dataset.idx!);
      const prefix = btn.dataset.prefix as "atk" | "def";
      const team = cachedTeams.find((t) => t.name === teamName);
      if (team && team.members[idx]) {
        importFromTeam(prefix, team.members[idx]);
      }
    });
  });

  container.querySelectorAll<HTMLButtonElement>(".team-delete-btn").forEach((btn) => {
    btn.addEventListener("click", () => handleDeleteTeam(btn.dataset.team!));
  });

  container.querySelectorAll<HTMLButtonElement>(".team-member-delete-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      handleDeleteTeamMember(btn.dataset.team!, parseInt(btn.dataset.idx!));
    });
  });
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
        <div class="build-moves-header">
          <h3 class="build-section-title">Movimientos atacante (máx. 4)</h3>
          <button class="battle-random-fill-btn" id="atk-random-fill-btn">Rellenar aleatoriamente</button>
        </div>
        <div class="build-slots">${renderMoveSlots(state.attacker.Moves ?? [], "atk")}</div>
      </div>`
    : "";

  const defenderMovesSection = state.defender
    ? `<div class="build-moves-section defender-moves-section">
        <div class="build-moves-header">
          <h3 class="build-section-title">Movimientos defensor (máx. 4)</h3>
          <button class="battle-random-fill-btn" id="def-random-fill-btn">Rellenar aleatoriamente</button>
        </div>
        <div class="build-slots">${renderMoveSlots(state.defender.Moves ?? [], "def")}</div>
      </div>`
    : "";

  const dmgSection = renderDamageSection();
  const btlSection = renderBattleSection();
  const teamsSection = renderTeamsSection();

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
        ${state.attacker ? `<button class="team-save-btn" data-prefix="atk">Guardar en equipo</button>` : ""}
      </div>

      <div class="build-col build-col--defender">
        <h3 class="build-col-title">Defensor</h3>
        <div class="build-search-row">
          <input id="def-input" class="build-search-input" type="text" placeholder="Nombre del Pokémon..." />
          <button id="def-btn" class="build-search-btn">Buscar</button>
        </div>
        ${defCard}
        ${defConfig}
        ${state.defender ? `<button class="team-save-btn" data-prefix="def">Guardar en equipo</button>` : ""}
      </div>
    </div>

    ${movesSection}
    ${defenderMovesSection}
    ${dmgSection}
    ${btlSection}
    ${teamsSection}
  `;

  bindEvents();
  bindBattleEvents();
  bindTeamEvents();

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

  container.querySelector<HTMLButtonElement>("#atk-random-fill-btn")?.addEventListener("click", () => randomFillSlots("atk"));
  container.querySelector<HTMLButtonElement>("#def-random-fill-btn")?.addEventListener("click", () => randomFillSlots("def"));

  container.querySelectorAll<HTMLSelectElement>(".build-move-select").forEach((sel) => {
    sel.addEventListener("change", () => {
      const idx = parseInt(sel.dataset.slot ?? "0");
      const prefix = (sel.dataset.prefix as "atk" | "def") ?? "atk";
      if (sel.value) loadMove(idx, sel.value, prefix);
    });
  });

  container.querySelectorAll<HTMLButtonElement>(".build-slot-clear").forEach((btn) => {
    btn.addEventListener("click", () => {
      const idx = parseInt(btn.dataset.slot ?? "0");
      const prefix = (btn.dataset.prefix as "atk" | "def") ?? "atk";
      if (prefix === "atk") state.slots[idx] = emptySlot();
      else state.defenderSlots[idx] = emptySlot();
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
      state.defenderSlots = [emptySlot(), emptySlot(), emptySlot(), emptySlot()];
    }
    battleUI = { battleState: null, phase: "idle" };

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

async function loadMove(slotIdx: number, moveName: string, prefix: "atk" | "def" = "atk"): Promise<void> {
  try {
    const move = await GetMove(moveName);
    if (prefix === "atk") {
      state.slots[slotIdx] = { moveName, move, isCritical: false };
    } else {
      state.defenderSlots[slotIdx] = { moveName, move, isCritical: false };
    }
    buildLayout();
  } catch (err: unknown) {
    alert(`Error al cargar movimiento: ${String(err)}`);
  }
}

async function randomFillSlots(prefix: "atk" | "def"): Promise<void> {
  const pokemon = prefix === "atk" ? state.attacker : state.defender;
  if (!pokemon) return;

  const available = pokemon.Moves ?? [];
  if (available.length === 0) return;

  // Pick up to 4 random moves (unique).
  const shuffled = [...available].sort(() => Math.random() - 0.5);
  const picked = shuffled.slice(0, 4);

  try {
    const moves = await Promise.all(picked.map((m) => GetMove(m.Name)));
    const newSlots: [BuildSlot, BuildSlot, BuildSlot, BuildSlot] = [
      emptySlot(), emptySlot(), emptySlot(), emptySlot(),
    ];
    moves.forEach((move, i) => {
      newSlots[i] = { moveName: move.Name, move, isCritical: false };
    });
    if (prefix === "atk") {
      state.slots = newSlots;
    } else {
      state.defenderSlots = newSlots;
    }
    buildLayout();
  } catch (err: unknown) {
    alert(`Error al rellenar movimientos: ${String(err)}`);
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
        ListTeams().then((teams) => { cachedTeams = teams ?? []; }),
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
