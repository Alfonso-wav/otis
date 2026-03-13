export namespace core {
	
	export class Ability {
	    Name: string;
	    Description: string;
	    Pokemon: string[];
	
	    static createFrom(source: any = {}) {
	        return new Ability(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Pokemon = source["Pokemon"];
	    }
	}
	export class NamedResource {
	    Name: string;
	    URL: string;
	
	    static createFrom(source: any = {}) {
	        return new NamedResource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.URL = source["URL"];
	    }
	}
	export class AbilityListResponse {
	    Count: number;
	    Results: NamedResource[];
	
	    static createFrom(source: any = {}) {
	        return new AbilityListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Count = source["Count"];
	        this.Results = this.convertValues(source["Results"], NamedResource);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Stats {
	    hp: number;
	    attack: number;
	    defense: number;
	    spAttack: number;
	    spDefense: number;
	    speed: number;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hp = source["hp"];
	        this.attack = source["attack"];
	        this.defense = source["defense"];
	        this.spAttack = source["spAttack"];
	        this.spDefense = source["spDefense"];
	        this.speed = source["speed"];
	    }
	}
	export class EVCalculatorInput {
	    pokemonName: string;
	    level: number;
	    natureName: string;
	    currentStats: Stats;
	    knownIVs?: Stats;
	
	    static createFrom(source: any = {}) {
	        return new EVCalculatorInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pokemonName = source["pokemonName"];
	        this.level = source["level"];
	        this.natureName = source["natureName"];
	        this.currentStats = this.convertValues(source["currentStats"], Stats);
	        this.knownIVs = this.convertValues(source["knownIVs"], Stats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StatRange {
	    min: number;
	    max: number;
	
	    static createFrom(source: any = {}) {
	        return new StatRange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.min = source["min"];
	        this.max = source["max"];
	    }
	}
	export class EVCalculatorResult {
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
	
	    static createFrom(source: any = {}) {
	        return new EVCalculatorResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pokemon = source["pokemon"];
	        this.level = source["level"];
	        this.nature = source["nature"];
	        this.baseStats = this.convertValues(source["baseStats"], Stats);
	        this.estimatedEVs = this.convertValues(source["estimatedEVs"], Stats);
	        this.evRanges = this.convertValues(source["evRanges"], StatRange, true);
	        this.totalEVsUsed = source["totalEVsUsed"];
	        this.evsRemaining = source["evsRemaining"];
	        this.maxPossibleStats = this.convertValues(source["maxPossibleStats"], Stats);
	        this.usedIVs = this.convertValues(source["usedIVs"], Stats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EggGroup {
	    Name: string;
	    Pokemon: string[];
	
	    static createFrom(source: any = {}) {
	        return new EggGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Pokemon = source["Pokemon"];
	    }
	}
	export class EvolutionStage {
	    Name: string;
	    MinLevel: number;
	    TriggerName: string;
	    EvolvesTo: EvolutionStage[];
	
	    static createFrom(source: any = {}) {
	        return new EvolutionStage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.MinLevel = source["MinLevel"];
	        this.TriggerName = source["TriggerName"];
	        this.EvolvesTo = this.convertValues(source["EvolvesTo"], EvolutionStage);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EvolutionChain {
	    ID: number;
	    Chain: EvolutionStage;
	
	    static createFrom(source: any = {}) {
	        return new EvolutionChain(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Chain = this.convertValues(source["Chain"], EvolutionStage);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Gender {
	    Name: string;
	    Pokemon: string[];
	
	    static createFrom(source: any = {}) {
	        return new Gender(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Pokemon = source["Pokemon"];
	    }
	}
	export class PokemonListItem {
	    Name: string;
	    URL: string;
	
	    static createFrom(source: any = {}) {
	        return new PokemonListItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.URL = source["URL"];
	    }
	}
	export class Generation {
	    Name: string;
	    MainRegion: string;
	    Games: string[];
	    PokemonSpecies: PokemonListItem[];
	    Types: string[];
	    Moves: string[];
	    Abilities: string[];
	
	    static createFrom(source: any = {}) {
	        return new Generation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.MainRegion = source["MainRegion"];
	        this.Games = source["Games"];
	        this.PokemonSpecies = this.convertValues(source["PokemonSpecies"], PokemonListItem);
	        this.Types = source["Types"];
	        this.Moves = source["Moves"];
	        this.Abilities = source["Abilities"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GrowthRateLevel {
	    Level: number;
	    Experience: number;
	
	    static createFrom(source: any = {}) {
	        return new GrowthRateLevel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Level = source["Level"];
	        this.Experience = source["Experience"];
	    }
	}
	export class GrowthRate {
	    Name: string;
	    Formula: string;
	    Levels: GrowthRateLevel[];
	    Pokemon: string[];
	
	    static createFrom(source: any = {}) {
	        return new GrowthRate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Formula = source["Formula"];
	        this.Levels = this.convertValues(source["Levels"], GrowthRateLevel);
	        this.Pokemon = source["Pokemon"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Location {
	    Name: string;
	    Region: string;
	
	    static createFrom(source: any = {}) {
	        return new Location(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Region = source["Region"];
	    }
	}
	export class PokemonEncounter {
	    PokemonName: string;
	    MaxChance: number;
	
	    static createFrom(source: any = {}) {
	        return new PokemonEncounter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PokemonName = source["PokemonName"];
	        this.MaxChance = source["MaxChance"];
	    }
	}
	export class LocationArea {
	    Name: string;
	    Location: string;
	    PokemonEncounters: PokemonEncounter[];
	
	    static createFrom(source: any = {}) {
	        return new LocationArea(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Location = source["Location"];
	        this.PokemonEncounters = this.convertValues(source["PokemonEncounters"], PokemonEncounter);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LocationDetail {
	    Name: string;
	    Region: string;
	    Areas: string[];
	
	    static createFrom(source: any = {}) {
	        return new LocationDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Region = source["Region"];
	        this.Areas = source["Areas"];
	    }
	}
	export class Machine {
	    ID: number;
	    Move: string;
	    Item: string;
	    VersionGroup: string;
	
	    static createFrom(source: any = {}) {
	        return new Machine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Move = source["Move"];
	        this.Item = source["Item"];
	        this.VersionGroup = source["VersionGroup"];
	    }
	}
	export class Move {
	    Name: string;
	    Type: string;
	    Power: number;
	    Accuracy: number;
	    PP: number;
	    Category: string;
	    Description: string;
	
	    static createFrom(source: any = {}) {
	        return new Move(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Power = source["Power"];
	        this.Accuracy = source["Accuracy"];
	        this.PP = source["PP"];
	        this.Category = source["Category"];
	        this.Description = source["Description"];
	    }
	}
	export class MoveAilment {
	    Name: string;
	    Moves: string[];
	
	    static createFrom(source: any = {}) {
	        return new MoveAilment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Moves = source["Moves"];
	    }
	}
	export class MoveDamageClass {
	    Name: string;
	    Description: string;
	    Moves: string[];
	
	    static createFrom(source: any = {}) {
	        return new MoveDamageClass(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Moves = source["Moves"];
	    }
	}
	export class MoveListResponse {
	    Count: number;
	    Results: NamedResource[];
	
	    static createFrom(source: any = {}) {
	        return new MoveListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Count = source["Count"];
	        this.Results = this.convertValues(source["Results"], NamedResource);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MoveTarget {
	    Name: string;
	    Description: string;
	
	    static createFrom(source: any = {}) {
	        return new MoveTarget(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	    }
	}
	
	export class Nature {
	    name: string;
	    increasedStat: string;
	    decreasedStat: string;
	
	    static createFrom(source: any = {}) {
	        return new Nature(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.increasedStat = source["increasedStat"];
	        this.decreasedStat = source["decreasedStat"];
	    }
	}
	export class NatureDetail {
	    Name: string;
	    IncreasedStat: string;
	    DecreasedStat: string;
	    LikesFlavor: string;
	    HatesFlavor: string;
	
	    static createFrom(source: any = {}) {
	        return new NatureDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.IncreasedStat = source["IncreasedStat"];
	        this.DecreasedStat = source["DecreasedStat"];
	        this.LikesFlavor = source["LikesFlavor"];
	        this.HatesFlavor = source["HatesFlavor"];
	    }
	}
	export class PokedexEntry {
	    EntryNumber: number;
	    Pokemon: string;
	
	    static createFrom(source: any = {}) {
	        return new PokedexEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EntryNumber = source["EntryNumber"];
	        this.Pokemon = source["Pokemon"];
	    }
	}
	export class Pokedex {
	    Name: string;
	    IsMainSeries: boolean;
	    Region: string;
	    PokemonEntries: PokedexEntry[];
	
	    static createFrom(source: any = {}) {
	        return new Pokedex(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.IsMainSeries = source["IsMainSeries"];
	        this.Region = source["Region"];
	        this.PokemonEntries = this.convertValues(source["PokemonEntries"], PokedexEntry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Sprites {
	    FrontDefault: string;
	    FrontShiny: string;
	
	    static createFrom(source: any = {}) {
	        return new Sprites(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FrontDefault = source["FrontDefault"];
	        this.FrontShiny = source["FrontShiny"];
	    }
	}
	export class Stat {
	    Name: string;
	    BaseStat: number;
	
	    static createFrom(source: any = {}) {
	        return new Stat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.BaseStat = source["BaseStat"];
	    }
	}
	export class PokemonType {
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new PokemonType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	    }
	}
	export class Pokemon {
	    ID: number;
	    Name: string;
	    Types: PokemonType[];
	    Stats: Stat[];
	    Sprites: Sprites;
	    Height: number;
	    Weight: number;
	
	    static createFrom(source: any = {}) {
	        return new Pokemon(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Types = this.convertValues(source["Types"], PokemonType);
	        this.Stats = this.convertValues(source["Stats"], Stat);
	        this.Sprites = this.convertValues(source["Sprites"], Sprites);
	        this.Height = source["Height"];
	        this.Weight = source["Weight"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class PokemonForm {
	    Name: string;
	    FormName: string;
	    IsMega: boolean;
	    IsBattleOnly: boolean;
	    Types: PokemonType[];
	    Sprites: Sprites;
	
	    static createFrom(source: any = {}) {
	        return new PokemonForm(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.FormName = source["FormName"];
	        this.IsMega = source["IsMega"];
	        this.IsBattleOnly = source["IsBattleOnly"];
	        this.Types = this.convertValues(source["Types"], PokemonType);
	        this.Sprites = this.convertValues(source["Sprites"], Sprites);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class PokemonListResponse {
	    Count: number;
	    Next: string;
	    Previous: string;
	    Results: PokemonListItem[];
	
	    static createFrom(source: any = {}) {
	        return new PokemonListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Count = source["Count"];
	        this.Next = source["Next"];
	        this.Previous = source["Previous"];
	        this.Results = this.convertValues(source["Results"], PokemonListItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PokemonVariety {
	    IsDefault: boolean;
	    Pokemon: string;
	
	    static createFrom(source: any = {}) {
	        return new PokemonVariety(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IsDefault = source["IsDefault"];
	        this.Pokemon = source["Pokemon"];
	    }
	}
	export class PokemonSpecies {
	    Name: string;
	    Order: number;
	    GenderRate: number;
	    CaptureRate: number;
	    BaseHappiness: number;
	    IsBaby: boolean;
	    IsLegendary: boolean;
	    IsMythical: boolean;
	    HatchCounter: number;
	    HasGenderDifferences: boolean;
	    FormsSwitchable: boolean;
	    Genus: string;
	    Color: string;
	    Shape: string;
	    Habitat: string;
	    EggGroups: string[];
	    FlavorText: string;
	    EvolutionChainID: number;
	    Varieties: PokemonVariety[];
	
	    static createFrom(source: any = {}) {
	        return new PokemonSpecies(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Order = source["Order"];
	        this.GenderRate = source["GenderRate"];
	        this.CaptureRate = source["CaptureRate"];
	        this.BaseHappiness = source["BaseHappiness"];
	        this.IsBaby = source["IsBaby"];
	        this.IsLegendary = source["IsLegendary"];
	        this.IsMythical = source["IsMythical"];
	        this.HatchCounter = source["HatchCounter"];
	        this.HasGenderDifferences = source["HasGenderDifferences"];
	        this.FormsSwitchable = source["FormsSwitchable"];
	        this.Genus = source["Genus"];
	        this.Color = source["Color"];
	        this.Shape = source["Shape"];
	        this.Habitat = source["Habitat"];
	        this.EggGroups = source["EggGroups"];
	        this.FlavorText = source["FlavorText"];
	        this.EvolutionChainID = source["EvolutionChainID"];
	        this.Varieties = this.convertValues(source["Varieties"], PokemonVariety);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TypePokemonEntry {
	    Name: string;
	    URL: string;
	
	    static createFrom(source: any = {}) {
	        return new TypePokemonEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.URL = source["URL"];
	    }
	}
	export class PokemonTypeDetail {
	    Name: string;
	    Pokemon: TypePokemonEntry[];
	
	    static createFrom(source: any = {}) {
	        return new PokemonTypeDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Pokemon = this.convertValues(source["Pokemon"], TypePokemonEntry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class Region {
	    Name: string;
	    Locations: Location[];
	
	    static createFrom(source: any = {}) {
	        return new Region(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Locations = this.convertValues(source["Locations"], Location);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class StatCalculatorInput {
	    pokemonName: string;
	    level: number;
	    natureName: string;
	    ivs: Stats;
	    evs: Stats;
	
	    static createFrom(source: any = {}) {
	        return new StatCalculatorInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pokemonName = source["pokemonName"];
	        this.level = source["level"];
	        this.natureName = source["natureName"];
	        this.ivs = this.convertValues(source["ivs"], Stats);
	        this.evs = this.convertValues(source["evs"], Stats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StatDetail {
	    Name: string;
	    IsBattleOnly: boolean;
	    AffectingMovesBuff: string[];
	    AffectingMovesNerf: string[];
	    AffectingNaturesBuff: string[];
	    AffectingNaturesNerf: string[];
	
	    static createFrom(source: any = {}) {
	        return new StatDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.IsBattleOnly = source["IsBattleOnly"];
	        this.AffectingMovesBuff = source["AffectingMovesBuff"];
	        this.AffectingMovesNerf = source["AffectingMovesNerf"];
	        this.AffectingNaturesBuff = source["AffectingNaturesBuff"];
	        this.AffectingNaturesNerf = source["AffectingNaturesNerf"];
	    }
	}
	
	
	export class TypeListResponse {
	    Count: number;
	    Results: PokemonListItem[];
	
	    static createFrom(source: any = {}) {
	        return new TypeListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Count = source["Count"];
	        this.Results = this.convertValues(source["Results"], PokemonListItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class VersionGroup {
	    Name: string;
	    Order: number;
	    Generation: string;
	    Versions: string[];
	    Pokedexes: string[];
	    Regions: string[];
	
	    static createFrom(source: any = {}) {
	        return new VersionGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Order = source["Order"];
	        this.Generation = source["Generation"];
	        this.Versions = source["Versions"];
	        this.Pokedexes = source["Pokedexes"];
	        this.Regions = source["Regions"];
	    }
	}
	export class StatComparison {
	    Name: string;
	    StatA: number;
	    StatB: number;
	    Diff: number;
	    Winner: string;

	    static createFrom(source: any = {}) {
	        return new StatComparison(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.StatA = source["StatA"];
	        this.StatB = source["StatB"];
	        this.Diff = source["Diff"];
	        this.Winner = source["Winner"];
	    }
	}
	export class PokemonComparison {
	    PokemonA: Pokemon;
	    PokemonB: Pokemon;
	    Stats: StatComparison[];
	    TotalA: number;
	    TotalB: number;
	    Winner: string;

	    static createFrom(source: any = {}) {
	        return new PokemonComparison(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PokemonA = this.convertValues(source["PokemonA"], Pokemon);
	        this.PokemonB = this.convertValues(source["PokemonB"], Pokemon);
	        this.Stats = this.convertValues(source["Stats"], StatComparison);
	        this.TotalA = source["TotalA"];
	        this.TotalB = source["TotalB"];
	        this.Winner = source["Winner"];
	    }

	    convertValues(a: any, classs: any, asMap: boolean = false): any {
	        if (!a) { return a; }
	        if (a.slice && a.map) { return (a as any[]).map(elem => this.convertValues(elem, classs)); }
	        else if ("object" === typeof a) { if (asMap) { for (const key of Object.keys(a)) { a[key] = new classs(a[key]); } return a; } return new classs(a); }
	        return a;
	    }
	}

}

