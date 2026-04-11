import { getLocale, setLocale } from "./i18n";
import { createAutocomplete } from "./autocomplete";
import { ListPokemon } from "./api";

const THEME_KEY = "theme";
const COMPANION_KEY = "companion-pokemon";
const DEFAULT_COMPANION = "diglett";

interface PendingSettings {
  locale: string | null;
  theme: ("dark" | "light") | null;
  companion: string | null;
}

const pending: PendingSettings = { locale: null, theme: null, companion: null };

function currentTheme(): "dark" | "light" {
  return localStorage.getItem(THEME_KEY) === "dark" ? "dark" : "light";
}

export function getCompanion(): string {
  return localStorage.getItem(COMPANION_KEY) || DEFAULT_COMPANION;
}

function companionSpriteUrl(name: string): string {
  const safeName = name.toLowerCase().replace(/[^a-z0-9-]/g, "");
  return `https://img.pokemondb.net/sprites/black-white/anim/normal/${safeName}.gif`;
}

function companionSpriteFallback(name: string): string {
  const safeName = name.toLowerCase().replace(/[^a-z0-9-]/g, "");
  return `https://img.pokemondb.net/sprites/black-white/normal/${safeName}.png`;
}

export function renderCompanion(): void {
  const name = getCompanion();
  const containers = [
    document.getElementById("header-companion"),
    document.getElementById("header-companion-mobile"),
  ];

  for (const container of containers) {
    if (!container) continue;
    const url = companionSpriteUrl(name);
    const fallback = companionSpriteFallback(name);
    container.innerHTML = `<img
      class="companion-sprite"
      src="${url}"
      alt="${name}"
      onerror="this.onerror=null;this.src='${fallback}'"
    />`;
  }
}

function renderCompanionPreview(name: string): void {
  const preview = document.getElementById("companion-preview");
  if (!preview) return;
  const url = companionSpriteUrl(name);
  const fallback = companionSpriteFallback(name);
  preview.innerHTML = `<img
    class="companion-sprite"
    src="${url}"
    alt="${name}"
    onerror="this.onerror=null;this.src='${fallback}'"
  />`;
}

function hasPendingChanges(): boolean {
  if (pending.locale !== null && pending.locale !== getLocale()) return true;
  if (pending.theme !== null && pending.theme !== currentTheme()) return true;
  if (pending.companion !== null && pending.companion !== getCompanion()) return true;
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

  // --- Companion setting ---
  const companionInput = document.getElementById("companion-input") as HTMLInputElement | null;
  if (companionInput) {
    companionInput.value = getCompanion();
    renderCompanionPreview(getCompanion());

    // Load pokemon names for autocomplete
    ListPokemon(0, 2000).then((data) => {
      const names = data.Results.map((r: { Name: string }) => r.Name);
      createAutocomplete(companionInput, names, (name) => {
        companionInput.value = name;
        pending.companion = name;
        renderCompanionPreview(name);
        updateApplyButton();
      });
    });

    companionInput.addEventListener("change", () => {
      const val = companionInput.value.trim().toLowerCase();
      if (val) {
        pending.companion = val;
        renderCompanionPreview(val);
        updateApplyButton();
      }
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
      if (pending.companion !== null && pending.companion !== getCompanion()) {
        localStorage.setItem(COMPANION_KEY, pending.companion);
        renderCompanion();
      }
      pending.locale = null;
      pending.theme = null;
      pending.companion = null;
      updateApplyButton();
    });
  }

  // Render companion in header on init
  renderCompanion();
}

export function cleanupSettings(): void {
  if (pending.theme !== null && pending.theme !== currentTheme()) {
    applyTheme(currentTheme());
  }
  pending.theme = null;
  pending.locale = null;
  pending.companion = null;

  const toggle = document.getElementById(
    "dark-mode-toggle",
  ) as HTMLInputElement | null;
  if (toggle) {
    toggle.checked = currentTheme() === "dark";
  }

  // Reset companion input to saved value
  const companionInput = document.getElementById("companion-input") as HTMLInputElement | null;
  if (companionInput) {
    companionInput.value = getCompanion();
    renderCompanionPreview(getCompanion());
  }

  updateApplyButton();
}

function applyTheme(theme: "dark" | "light"): void {
  document.documentElement.setAttribute("data-bs-theme", theme);
}
