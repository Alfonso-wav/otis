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

export interface TypePokemonEntry {
  Name: string;
  URL: string;
}

export interface PokemonTypeDetail {
  Name: string;
  Pokemon: TypePokemonEntry[];
}

export interface TypeListResponse {
  Count: number;
  Results: PokemonListItem[];
}
