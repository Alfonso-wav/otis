// move-names.ts — Lazy-loaded cache of move slug → localized name, power & category.
// Uses GetAllMoves() to build lookup maps once, then resolves data without extra
// network calls.

import { GetAllMoves } from "../api";
import { getLocale } from "../i18n";

interface MoveData {
  nameEs: string;
  power: number;
  category: string; // "physical" | "special" | "status"
}

let moveCache: Map<string, MoveData> | null = null;
let loadPromise: Promise<void> | null = null;

/**
 * Loads the move cache (idempotent).
 * Safe to call multiple times — only the first call fetches.
 */
export function loadMoveNames(): Promise<void> {
  if (moveCache) return Promise.resolve();
  if (loadPromise) return loadPromise;

  loadPromise = GetAllMoves()
    .then((moves) => {
      moveCache = new Map();
      for (const m of moves) {
        moveCache.set(m.Name, {
          nameEs: m.NameEs || "",
          power: m.Power ?? 0,
          category: m.Category || "status",
        });
      }
    })
    .catch(() => {
      // If fetch fails, leave cache null so next call retries
      loadPromise = null;
    });

  return loadPromise;
}

/**
 * Returns the localized display name for a move slug.
 * - locale "es" + NameEs exists → NameEs
 * - otherwise → capitalize(slug with spaces)
 */
export function getLocalizedMoveName(slug: string): string {
  const fallback = capitalize(slug.replace(/-/g, " "));
  if (getLocale() !== "es") return fallback;
  if (!moveCache) return fallback;
  const data = moveCache.get(slug);
  return data?.nameEs || fallback;
}

/**
 * Returns the power of a move, or null if it's a status move / unknown.
 */
export function getMovePower(slug: string): number | null {
  if (!moveCache) return null;
  const data = moveCache.get(slug);
  if (!data) return null;
  return data.power > 0 ? data.power : null;
}

/**
 * Returns the category of a move: "physical", "special", or "status".
 * Falls back to "status" for unknown moves.
 */
export function getMoveCategory(slug: string): "physical" | "special" | "status" {
  if (!moveCache) return "status";
  const data = moveCache.get(slug);
  if (!data) return "status";
  const cat = data.category.toLowerCase();
  if (cat === "physical" || cat === "special") return cat;
  return "status";
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}
