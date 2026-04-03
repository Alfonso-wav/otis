import { getLocale, setLocale } from "./i18n";

const THEME_KEY = "theme";

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
      const theme = toggle.checked ? "dark" : "light";
      applyTheme(theme);
      localStorage.setItem(THEME_KEY, theme);
    });
  }

  const langSelect = document.getElementById(
    "language-select",
  ) as HTMLSelectElement | null;
  if (langSelect) {
    langSelect.value = getLocale();
    langSelect.addEventListener("change", () => {
      setLocale(langSelect.value);
    });
  }
}

function applyTheme(theme: "dark" | "light"): void {
  document.documentElement.setAttribute("data-bs-theme", theme);
}
