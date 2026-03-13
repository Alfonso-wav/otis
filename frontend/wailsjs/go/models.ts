export namespace core {
	
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

