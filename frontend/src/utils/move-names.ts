// move-names.ts — Lazy-loaded cache of move slug → localized name (NameEs).
// Uses GetAllMoves() to build a Map<slug, NameEs> once, then resolves
// localized names without extra network calls.

import { GetAllMoves } from "../api";
import { getLocale } from "../i18n";

let moveNameCache: Map<string, string> | null = null;
let loadPromise: Promise<void> | null = null;

/**
 * Loads the move-name cache (idempotent).
 * Safe to call multiple times — only the first call fetches.
 */
export function loadMoveNames(): Promise<void> {
  if (moveNameCache) return Promise.resolve();
  if (loadPromise) return loadPromise;

  loadPromise = GetAllMoves()
    .then((moves) => {
      moveNameCache = new Map();
      for (const m of moves) {
        if (m.NameEs) {
          moveNameCache.set(m.Name, m.NameEs);
        }
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
  if (!moveNameCache) return fallback;
  return moveNameCache.get(slug) || fallback;
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}
