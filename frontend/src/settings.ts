import { getLocale, setLocale } from "./i18n";
import { ListPokemon, ListTeams } from "./api";
import { createAutocomplete } from "./autocomplete";

const THEME_KEY = "theme";
const COMPANION_TEAM_KEY = "companion-team";
const COMPANION_KEY = "companion-pokemon"; // old key for migration
const DEFAULT_COMPANIONS = ["diglett", "", "", "", "", ""];

const SPRITE_BASE = "/assets/sprites";
const CDN_FALLBACK = "https://img.pokemondb.net/sprites/home/normal";

function companionSpriteUrl(name: string): string {
  return `${SPRITE_BASE}/${name}.png`;
}

function companionSpriteFallback(name: string): string {
  return `${CDN_FALLBACK}/${name}.png`;
}

interface PendingSettings {
  locale: string | null;
  theme: ("dark" | "light") | null;
  companionTeam: string[] | null;
}

const pending: PendingSettings = { locale: null, theme: null, companionTeam: null };

function currentTheme(): "dark" | "light" {
  return localStorage.getItem(THEME_KEY) === "dark" ? "dark" : "light";
}

export function getCompanionTeam(): string[] {
  const stored = localStorage.getItem(COMPANION_TEAM_KEY);
  if (stored) {
    try {
      const arr = JSON.parse(stored);
      if (Array.isArray(arr)) {
        // Ensure exactly 6 slots
        const team = arr.slice(0, 6).map((s: unknown) => (typeof s === "string" ? s : ""));
        while (team.length < 6) team.push("");
        return team;
      }
    } catch { /* fallthrough */ }
  }
  // Migration: read old single companion
  const old = localStorage.getItem(COMPANION_KEY);
  if (old) {
    const team = [old, "", "", "", "", ""];
    localStorage.setItem(COMPANION_TEAM_KEY, JSON.stringify(team));
    localStorage.removeItem(COMPANION_KEY);
    return team;
  }
  return [...DEFAULT_COMPANIONS];
}

export function renderCompanion(): void {
  const team = getCompanionTeam();
  const container = document.getElementById("header-companion");
  if (!container) return;
  container.innerHTML = team
    .filter((name) => name !== "")
    .map((name) => {
      const url = companionSpriteUrl(name);
      const fallback = companionSpriteFallback(name);
      return `<img class="companion-sprite" src="${url}" alt="${name}" onerror="this.onerror=null;this.src='${fallback}'" />`;
    })
    .join("");
}

function renderCompanionPreview(team: string[]): void {
  const preview = document.getElementById("companion-preview");
  if (!preview) return;
  preview.innerHTML = team
    .filter((name) => name !== "")
    .map((name) => {
      const url = companionSpriteUrl(name);
      const fallback = companionSpriteFallback(name);
      return `<img class="companion-sprite" src="${url}" alt="${name}" onerror="this.onerror=null;this.src='${fallback}'" />`;
    })
    .join("") || '<span style="color:#a0aec0;font-size:0.8rem">—</span>';
}

function hasPendingChanges(): boolean {
  if (pending.locale !== null && pending.locale !== getLocale()) return true;
  if (pending.theme !== null && pending.theme !== currentTheme()) return true;
  if (pending.companionTeam !== null && JSON.stringify(pending.companionTeam) !== JSON.stringify(getCompanionTeam())) return true;
  return false;
}

function updateApplyButton(): void {
  const btn = document.getElementById("settings-apply-btn") as HTMLButtonElement | null;
  if (btn) btn.disabled = !hasPendingChanges();
}

export function initSettings(): void {
  const saved = localStorage.getItem(THEME_KEY);
  if (saved === "dark") {
    applyTheme("dark");
  }

  // Render companion sprites in the header
  renderCompanion();

  const toggle = document.getElementById(
    "dark-mode-toggle",
  ) as HTMLInputElement | null;
  if (toggle) {
    toggle.checked = saved === "dark";
    toggle.addEventListener("change", () => {
      pending.theme = toggle.checked ? "dark" : "light";
      applyTheme(pending.theme);
      updateApplyButton();
    });
  }

  const langSelect = document.getElementById(
    "language-select",
  ) as HTMLSelectElement | null;
  if (langSelect) {
    langSelect.value = getLocale();
    langSelect.addEventListener("change", () => {
      pending.locale = langSelect.value;
      updateApplyButton();
    });
  }

  // --- Companion team setting (6 slots) ---
  const companionSlots = document.querySelectorAll<HTMLInputElement>(".companion-slot-input");
  const currentTeam = getCompanionTeam();

  if (companionSlots.length > 0) {
    ListPokemon(0, 2000).then((data) => {
      const names = data.Results.map((r: { Name: string }) => r.Name);
      companionSlots.forEach((input, i) => {
        input.value = currentTeam[i] || "";
        createAutocomplete(input, names, (name) => {
          input.value = name;
          if (!pending.companionTeam) pending.companionTeam = [...getCompanionTeam()];
          pending.companionTeam[i] = name;
          renderCompanionPreview(pending.companionTeam);
          updateApplyButton();
        });
        input.addEventListener("change", () => {
          const val = input.value.trim().toLowerCase();
          if (!pending.companionTeam) pending.companionTeam = [...getCompanionTeam()];
          pending.companionTeam[i] = val;
          renderCompanionPreview(pending.companionTeam);
          updateApplyButton();
        });
      });
    });
    renderCompanionPreview(currentTeam);
  }

  // Load from team button
  const loadTeamBtn = document.getElementById("load-team-btn");
  if (loadTeamBtn) {
    loadTeamBtn.addEventListener("click", async () => {
      try {
        const teams = await ListTeams();
        if (teams.length === 0) return;
        const teamSelect = document.getElementById("companion-team-select") as HTMLSelectElement | null;
        if (!teamSelect) return;
        teamSelect.innerHTML = teams.map((t) =>
          `<option value="${t.name}">${t.name}</option>`
        ).join("");
        teamSelect.classList.remove("hidden");
        teamSelect.addEventListener("change", () => {
          const selected = teams.find((t) => t.name === teamSelect.value);
          if (!selected) return;
          const newTeam = selected.members
            .slice(0, 6)
            .map((m) => m.pokemonName || "");
          while (newTeam.length < 6) newTeam.push("");
          pending.companionTeam = newTeam;
          companionSlots.forEach((input, i) => { input.value = newTeam[i] || ""; });
          renderCompanionPreview(newTeam);
          updateApplyButton();
          teamSelect.classList.add("hidden");
        }, { once: true });
      } catch { /* ignore */ }
    });
  }

  const applyBtn = document.getElementById("settings-apply-btn") as HTMLButtonElement | null;
  if (applyBtn) {
    applyBtn.addEventListener("click", () => {
      if (pending.theme !== null && pending.theme !== currentTheme()) {
        applyTheme(pending.theme);
        localStorage.setItem(THEME_KEY, pending.theme);
      }
      if (pending.locale !== null && pending.locale !== getLocale()) {
        setLocale(pending.locale);
      }
      if (pending.companionTeam !== null && JSON.stringify(pending.companionTeam) !== JSON.stringify(getCompanionTeam())) {
        localStorage.setItem(COMPANION_TEAM_KEY, JSON.stringify(pending.companionTeam));
        renderCompanion();
      }
      // reset
      pending.locale = null;
      pending.theme = null;
      pending.companionTeam = null;
      updateApplyButton();
    });
  }
}

export function cleanupSettings(): void {
  if (pending.theme !== null && pending.theme !== currentTheme()) {
    applyTheme(currentTheme());
  }
  pending.theme = null;
  pending.locale = null;
  pending.companionTeam = null;

  const toggle = document.getElementById(
    "dark-mode-toggle",
  ) as HTMLInputElement | null;
  if (toggle) {
    toggle.checked = currentTheme() === "dark";
  }

  // Reset companion inputs to saved values
  const companionSlots = document.querySelectorAll<HTMLInputElement>(".companion-slot-input");
  const saved = getCompanionTeam();
  companionSlots.forEach((input, i) => { input.value = saved[i] || ""; });
  renderCompanionPreview(saved);

  updateApplyButton();
}

function applyTheme(theme: "dark" | "light"): void {
  document.documentElement.setAttribute("data-bs-theme", theme);
}
