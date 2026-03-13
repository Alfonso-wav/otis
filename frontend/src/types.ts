export interface Sprites {
  FrontDefault: string;
  FrontShiny: string;
}

export interface Stat {
  Name: string;
  BaseStat: number;
}

export interface PokemonType {
  Name: string;
}

export interface Pokemon {
  ID: number;
  Name: string;
  Types: PokemonType[];
  Stats: Stat[];
  Sprites: Sprites;
  Height: number;
  Weight: number;
}

export interface PokemonListItem {
  Name: string;
  URL: string;
}

export interface PokemonListResponse {
  Count: number;
  Next: string;
  Previous: string;
  Results: PokemonListItem[];
}

// Tipos para la calculadora de EVs
export interface Stats {
  hp: number;
  attack: number;
  defense: number;
  spAttack: number;
  spDefense: number;
  speed: number;
}

export interface Nature {
  name: string;
  increasedStat: string;
  decreasedStat: string;
}

export interface StatRange {
  min: number;
  max: number;
}

export interface EVCalculatorInput {
  pokemonName: string;
  level: number;
  natureName: string;
  currentStats: Stats;
  knownIVs?: Stats;
}

export interface EVCalculatorResult {
  pokemon: string;
  level: number;
  nature: string;
  baseStats: Stats;
  estimatedEVs: Stats;
  evRanges: Record<string, StatRange>;
  totalEVsUsed: number;
  evsRemaining: number;
  maxPossibleStats: Stats;
  usedIVs: Stats;
}
