import { describe, it, expect } from "vitest";
import { SortCache } from "./sort-cache";

interface Item {
  name: string;
  value: number;
}

const columns = [
  { key: "name", compare: (a: Item, b: Item) => a.name.localeCompare(b.name) },
  { key: "value", compare: (a: Item, b: Item) => a.value - b.value },
];

describe("SortCache", () => {
  it("returns sorted data by column asc", () => {
    const cache = new SortCache<Item>(columns);
    cache.setData([
      { name: "cherry", value: 3 },
      { name: "apple", value: 1 },
      { name: "banana", value: 2 },
    ]);
    const result = cache.get("name", "asc");
    expect(result.map((i) => i.name)).toEqual(["apple", "banana", "cherry"]);
  });

  it("returns sorted data by column desc", () => {
    const cache = new SortCache<Item>(columns);
    cache.setData([
      { name: "cherry", value: 3 },
      { name: "apple", value: 1 },
      { name: "banana", value: 2 },
    ]);
    const result = cache.get("value", "desc");
    expect(result.map((i) => i.value)).toEqual([3, 2, 1]);
  });

  it("invalidates cache and returns empty on get after invalidate", () => {
    const cache = new SortCache<Item>(columns);
    cache.setData([{ name: "a", value: 1 }]);
    cache.invalidate();
    const result = cache.get("name", "asc");
    expect(result).toEqual([]);
  });

  it("returns sourceData for unknown column", () => {
    const cache = new SortCache<Item>(columns);
    const data = [
      { name: "b", value: 2 },
      { name: "a", value: 1 },
    ];
    cache.setData(data);
    const result = cache.get("unknown", "asc");
    expect(result).toEqual(data);
  });

  it("re-computes cache when setData is called with new data", () => {
    const cache = new SortCache<Item>(columns);
    cache.setData([
      { name: "b", value: 2 },
      { name: "a", value: 1 },
    ]);
    expect(cache.get("name", "asc").map((i) => i.name)).toEqual(["a", "b"]);

    cache.setData([
      { name: "z", value: 10 },
      { name: "m", value: 5 },
    ]);
    expect(cache.get("name", "asc").map((i) => i.name)).toEqual(["m", "z"]);
  });
});
