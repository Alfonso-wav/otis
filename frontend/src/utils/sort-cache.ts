export interface SortDef<T> {
  key: string;
  compare: (a: T, b: T) => number;
}

export class SortCache<T> {
  private cache = new Map<string, T[]>();
  private sourceData: T[] = [];
  private columns: SortDef<T>[];

  constructor(columns: SortDef<T>[]) {
    this.columns = columns;
  }

  setData(data: T[]): void {
    this.sourceData = data;
    this.cache.clear();
    for (const col of this.columns) {
      const asc = [...data].sort(col.compare);
      this.cache.set(`${col.key}:asc`, asc);
      this.cache.set(`${col.key}:desc`, [...asc].reverse());
    }
  }

  get(column: string, direction: "asc" | "desc"): T[] {
    return this.cache.get(`${column}:${direction}`) ?? this.sourceData;
  }

  invalidate(): void {
    this.cache.clear();
    this.sourceData = [];
  }
}
