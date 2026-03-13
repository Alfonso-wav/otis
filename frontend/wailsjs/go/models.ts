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

}

