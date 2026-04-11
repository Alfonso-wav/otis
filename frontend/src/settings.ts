import { getLocale, setLocale } from "./i18n";

const THEME_KEY = "theme";

interface PendingSettings {
  locale: string | null;
  theme: ("dark" | "light") | null;
}

const pending: PendingSettings = { locale: null, theme: null };

function currentTheme(): "dark" | "light" {
  return localStorage.getItem(THEME_KEY) === "dark" ? "dark" : "light";
}

function hasPendingChanges(): boolean {
  if (pending.locale !== null && pending.locale !== getLocale()) return true;
  if (pending.theme !== null && pending.theme !== currentTheme()) return true;
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
      pending.locale = null;
      pending.theme = null;
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

  const toggle = document.getElementById(
    "dark-mode-toggle",
  ) as HTMLInputElement | null;
  if (toggle) {
    toggle.checked = currentTheme() === "dark";
  }
  updateApplyButton();
}

function applyTheme(theme: "dark" | "light"): void {
  document.documentElement.setAttribute("data-bs-theme", theme);
}
