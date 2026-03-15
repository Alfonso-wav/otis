// api.ts — Unified backend abstraction layer.
// Detects Wails (desktop) vs HTTP (mobile) and uses the correct transport.
//
// In Wails mode: calls window['go']['app']['App'][method]() via IPC.
// In HTTP mode: calls fetch() to the REST server.

import type { core } from "../wailsjs/go/models";

// @ts-ignore — Vite injects import.meta.env at build time
const API_BASE: string = (import.meta as any).env?.VITE_API_BASE ?? "http://localhost:8080";

function isWails(): boolean {
  return typeof (window as any)?.go?.app?.App !== "undefined";
}

function wails(method: string, ...args: any[]): Promise<any> {
  return (window as any).go.app.App[method](...args);
}

const RETRY_DELAYS = [500, 1000, 2000];

async function fetchWithRetry(url: string, init?: RequestInit): Promise<Response> {
  let lastError: unknown;
  for (let attempt = 0; attempt <= RETRY_DELAYS.length; attempt++) {
    try {
      const res = await fetch(url, init);
      return res;
    } catch (err) {
      lastError = err;
      if (attempt < RETRY_DELAYS.length) {
        await new Promise((r) => setTimeout(r, RETRY_DELAYS[attempt]));
      }
    }
  }
  throw lastError;
}

async function get<T>(path: string): Promise<T> {
  const res = await fetchWithRetry(`${API_BASE}${path}`);
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error ?? res.statusText);
  }
  return res.json();
}

async function post<T>(path: string, body: unknown): Promise<T> {
  const res = await fetchWithRetry(`${API_BASE}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(data.error ?? res.statusText);
  }
  return res.json();
}

async function put<T>(path: string, body: unknown): Promise<T> {
  const res = await fetchWithRetry(`${API_BASE}${path}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const data = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(data.error ?? res.statusText);
  }
  return res.json();
}

async function del<T>(path: string): Promise<T> {
  const res = await fetchWithRetry(`${API_BASE}${path}`, { method: "DELETE" });
  if (!res.ok) {
    const data = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(data.error ?? res.statusText);
  }
  return res.json();
}

// --- Pokémon base ---

export function ListPokemon(offset: number, limit: number): Promise<core.PokemonListResponse> {
  if (isWails()) return wails("ListPokemon", offset, limit);
  return get(`/api/pokemon?offset=${offset}&limit=${limit}`);
}

export function GetPokemon(name: string): Promise<core.Pokemon> {
  if (isWails()) return wails("GetPokemon", name);
  return get(`/api/pokemon/${encodeURIComponent(name)}`);
}

export function ListTypes(): Promise<core.TypeListResponse> {
  if (isWails()) return wails("ListTypes");
  return get("/api/types");
}

export function GetType(name: string): Promise<core.PokemonTypeDetail> {
  if (isWails()) return wails("GetType", name);
  return get(`/api/types/${encodeURIComponent(name)}`);
}

export function ListRegions(): Promise<core.Region[]> {
  if (isWails()) return wails("ListRegions");
  return get("/api/regions");
}

export function GetRegion(name: string): Promise<core.Region> {
  if (isWails()) return wails("GetRegion", name);
  return get(`/api/regions/${encodeURIComponent(name)}`);
}

export function GetRegionPokemonByType(region: string, typeName: string): Promise<string[]> {
  if (isWails()) return wails("GetRegionPokemonByType", region, typeName);
  return get(`/api/regions/${encodeURIComponent(region)}/pokemon-by-type/${encodeURIComponent(typeName)}`);
}

export function GetMove(name: string): Promise<core.Move> {
  if (isWails()) return wails("GetMove", name);
  return get(`/api/moves/${encodeURIComponent(name)}`);
}

export function GetAbility(name: string): Promise<core.Ability> {
  if (isWails()) return wails("GetAbility", name);
  return get(`/api/abilities/${encodeURIComponent(name)}`);
}

export function GetEvolutionChain(id: number): Promise<core.EvolutionChain> {
  if (isWails()) return wails("GetEvolutionChain", id);
  return get(`/api/evolution-chain/${id}`);
}

export function GetNatures(): Promise<core.Nature[]> {
  if (isWails()) return wails("GetNatures");
  return get("/api/natures");
}

// --- Grupo A: Pokémon extendido ---

export function GetPokemonSpecies(name: string): Promise<core.PokemonSpecies> {
  if (isWails()) return wails("GetPokemonSpecies", name);
  return get(`/api/pokemon-species/${encodeURIComponent(name)}`);
}

export function GetPokemonForm(name: string): Promise<core.PokemonForm> {
  if (isWails()) return wails("GetPokemonForm", name);
  return get(`/api/pokemon-forms/${encodeURIComponent(name)}`);
}

export function ListPokemonColors(): Promise<core.NamedResource[]> {
  if (isWails()) return wails("ListPokemonColors");
  return get("/api/pokemon-colors");
}

export function ListPokemonShapes(): Promise<core.NamedResource[]> {
  if (isWails()) return wails("ListPokemonShapes");
  return get("/api/pokemon-shapes");
}

export function ListPokemonHabitats(): Promise<core.NamedResource[]> {
  if (isWails()) return wails("ListPokemonHabitats");
  return get("/api/pokemon-habitats");
}

// --- Grupo B: Naturalezas y cría ---

export function ListNatureNames(): Promise<core.NamedResource[]> {
  if (isWails()) return wails("ListNatureNames");
  return get("/api/nature-names");
}

export function GetNatureDetail(name: string): Promise<core.NatureDetail> {
  if (isWails()) return wails("GetNatureDetail", name);
  return get(`/api/nature-detail/${encodeURIComponent(name)}`);
}

export function GetEggGroup(name: string): Promise<core.EggGroup> {
  if (isWails()) return wails("GetEggGroup", name);
  return get(`/api/egg-groups/${encodeURIComponent(name)}`);
}

export function GetGender(name: string): Promise<core.Gender> {
  if (isWails()) return wails("GetGender", name);
  return get(`/api/genders/${encodeURIComponent(name)}`);
}

export function GetGrowthRate(name: string): Promise<core.GrowthRate> {
  if (isWails()) return wails("GetGrowthRate", name);
  return get(`/api/growth-rates/${encodeURIComponent(name)}`);
}

// --- Grupo C: Movimientos ---

export function GetAllMoves(): Promise<core.Move[]> {
  if (isWails()) return wails("GetAllMoves");
  return get("/api/moves/all");
}

export function ListMoves(offset: number, limit: number): Promise<core.MoveListResponse> {
  if (isWails()) return wails("ListMoves", offset, limit);
  return get(`/api/moves?offset=${offset}&limit=${limit}`);
}

export function GetMoveDamageClass(name: string): Promise<core.MoveDamageClass> {
  if (isWails()) return wails("GetMoveDamageClass", name);
  return get(`/api/move-damage-classes/${encodeURIComponent(name)}`);
}

export function GetMoveAilment(name: string): Promise<core.MoveAilment> {
  if (isWails()) return wails("GetMoveAilment", name);
  return get(`/api/move-ailments/${encodeURIComponent(name)}`);
}

export function GetMoveTarget(name: string): Promise<core.MoveTarget> {
  if (isWails()) return wails("GetMoveTarget", name);
  return get(`/api/move-targets/${encodeURIComponent(name)}`);
}

export function GetMachine(id: number): Promise<core.Machine> {
  if (isWails()) return wails("GetMachine", id);
  return get(`/api/machines/${id}`);
}

// --- Grupo D: Habilidades ---

export function GetAllAbilities(): Promise<core.Ability[]> {
  if (isWails()) return wails("GetAllAbilities");
  return get("/api/abilities/all");
}

export function ListAbilities(offset: number, limit: number): Promise<core.AbilityListResponse> {
  if (isWails()) return wails("ListAbilities", offset, limit);
  return get(`/api/abilities?offset=${offset}&limit=${limit}`);
}

// --- Grupo F: Ubicaciones ---

export function GetLocation(name: string): Promise<core.LocationDetail> {
  if (isWails()) return wails("GetLocation", name);
  return get(`/api/locations/${encodeURIComponent(name)}`);
}

export function GetLocationArea(name: string): Promise<core.LocationArea> {
  if (isWails()) return wails("GetLocationArea", name);
  return get(`/api/location-areas/${encodeURIComponent(name)}`);
}

// --- Grupo G: Stats y generaciones ---

export function GetStatDetail(name: string): Promise<core.StatDetail> {
  if (isWails()) return wails("GetStatDetail", name);
  return get(`/api/stats/${encodeURIComponent(name)}`);
}

export function ListGenerations(): Promise<core.NamedResource[]> {
  if (isWails()) return wails("ListGenerations");
  return get("/api/generations");
}

export function GetGeneration(name: string): Promise<core.Generation> {
  if (isWails()) return wails("GetGeneration", name);
  return get(`/api/generations/${encodeURIComponent(name)}`);
}

export function ListPokedexes(): Promise<core.NamedResource[]> {
  if (isWails()) return wails("ListPokedexes");
  return get("/api/pokedexes");
}

export function GetPokedex(name: string): Promise<core.Pokedex> {
  if (isWails()) return wails("GetPokedex", name);
  return get(`/api/pokedexes/${encodeURIComponent(name)}`);
}

export function ListVersionGroups(): Promise<core.NamedResource[]> {
  if (isWails()) return wails("ListVersionGroups");
  return get("/api/version-groups");
}

export function GetVersionGroup(name: string): Promise<core.VersionGroup> {
  if (isWails()) return wails("GetVersionGroup", name);
  return get(`/api/version-groups/${encodeURIComponent(name)}`);
}

// --- Scraper ---

export function ScrapePokedex(): Promise<core.PokedexDBEntry[]> {
  if (isWails()) return wails("ScrapePokedex");
  return get("/api/scrape/pokedex");
}

// --- Sprites ---

export function DownloadSprites(categories: string[]): Promise<core.SpriteDownloadResult> {
  if (isWails()) return wails("DownloadSprites", categories);
  return post("/api/sprites/download", { categories });
}

// --- Batalla ---

export function SimulateDamage(input: core.DamageInput): Promise<core.DamageResult> {
  if (isWails()) return wails("SimulateDamage", input);
  return post("/api/battle/simulate-damage", input);
}

export function InitBattle(attackerMaxHP: number, defenderMaxHP: number): Promise<core.BattleState> {
  if (isWails()) return wails("InitBattle", attackerMaxHP, defenderMaxHP);
  return post("/api/battle/init", { attackerMaxHP, defenderMaxHP });
}

export function ExecuteTurn(input: core.TurnInput): Promise<core.TurnResult> {
  if (isWails()) return wails("ExecuteTurn", input);
  return post("/api/battle/execute-turn", input);
}

export function SimulateFullBattle(input: core.FullBattleInput): Promise<core.BattleState> {
  if (isWails()) return wails("SimulateFullBattle", input);
  return post("/api/battle/simulate-full", input);
}

export function SimulateMultipleBattles(input: core.FullBattleInput, n: number): Promise<core.BattleReport> {
  if (isWails()) return wails("SimulateMultipleBattles", input, n);
  return post("/api/battle/simulate-multiple", { ...input, n });
}

export function SimulateTeamBattle(team1Name: string, team2Name: string): Promise<core.TeamBattleState> {
  if (isWails()) return wails("SimulateTeamBattle", team1Name, team2Name);
  return post("/api/battle/team-simulate", { team1Name, team2Name });
}

export function SimulateMultipleTeamBattles(team1Name: string, team2Name: string, n: number): Promise<core.TeamBattleReport> {
  if (isWails()) return wails("SimulateMultipleTeamBattles", team1Name, team2Name, n);
  return post("/api/battle/team-simulate-multiple", { team1Name, team2Name, n });
}

// --- Calculadoras ---

export function CalculateEVs(input: core.EVCalculatorInput): Promise<core.EVCalculatorResult> {
  if (isWails()) return wails("CalculateEVs", input);
  return post("/api/calculator/evs", input);
}

export function CalculateStats(input: core.StatCalculatorInput): Promise<core.Stats> {
  if (isWails()) return wails("CalculateStats", input);
  return post("/api/calculator/stats", input);
}

// --- Equipos ---

export function ListTeams(): Promise<core.Team[]> {
  if (isWails()) return wails("ListTeams");
  return get("/api/teams");
}

export function GetTeam(name: string): Promise<core.Team> {
  if (isWails()) return wails("GetTeam", name);
  return get(`/api/teams/${encodeURIComponent(name)}`);
}

export function CreateTeam(name: string): Promise<void> {
  if (isWails()) return wails("CreateTeam", name);
  return post("/api/teams", { name });
}

export function DeleteTeam(name: string): Promise<void> {
  if (isWails()) return wails("DeleteTeam", name);
  return del(`/api/teams/${encodeURIComponent(name)}`);
}

export function SaveToTeam(teamName: string, member: core.TeamMember): Promise<void> {
  if (isWails()) return wails("SaveToTeam", teamName, member);
  return post(`/api/teams/${encodeURIComponent(teamName)}/members`, member);
}

export function UpdateTeamMember(teamName: string, memberIndex: number, member: core.TeamMember): Promise<void> {
  if (isWails()) return wails("UpdateTeamMember", teamName, memberIndex, member);
  return put(`/api/teams/${encodeURIComponent(teamName)}/members/${memberIndex}`, member);
}

export function DeleteTeamMember(teamName: string, memberIndex: number): Promise<void> {
  if (isWails()) return wails("DeleteTeamMember", teamName, memberIndex);
  return del(`/api/teams/${encodeURIComponent(teamName)}/members/${memberIndex}`);
}

export function FillTeamRandom(teamName: string): Promise<core.Team> {
  if (isWails()) return wails("FillTeamRandom", teamName);
  return post(`/api/teams/${encodeURIComponent(teamName)}/fill-random`, {});
}
