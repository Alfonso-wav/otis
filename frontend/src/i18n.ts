import es from "./locales/es.json";
import en from "./locales/en.json";

type Translations = Record<string, unknown>;

const locales: Record<string, Translations> = { es, en };

const LOCALE_KEY = "locale";
const DEFAULT_LOCALE = "es";

let current: string = localStorage.getItem(LOCALE_KEY) ?? DEFAULT_LOCALE;

function resolve(key: string, data: Translations): string {
  const parts = key.split(".");
  let node: unknown = data;
  for (const part of parts) {
    if (node === null || typeof node !== "object") return key;
    node = (node as Record<string, unknown>)[part];
  }
  return typeof node === "string" ? node : key;
}

export function t(key: string, params?: Record<string, string | number>): string {
  let value = resolve(key, locales[current] ?? locales[DEFAULT_LOCALE]);
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      value = value.replaceAll(`{${k}}`, String(v));
    }
  }
  return value;
}

export function typeName(apiName: string): string {
  return t(`typeNames.${apiName.toLowerCase()}`);
}

export function statName(apiName: string): string {
  return t(`statNames.${apiName}`);
}

export function getLocale(): string {
  return current;
}

export function setLocale(locale: string): void {
  if (!locales[locale]) return;
  current = locale;
  localStorage.setItem(LOCALE_KEY, locale);
  applyTranslations();
  document.dispatchEvent(new CustomEvent("locale-changed", { detail: locale }));
}

export function applyTranslations(): void {
  document.querySelectorAll<HTMLElement>("[data-i18n]").forEach((el) => {
    const key = el.getAttribute("data-i18n")!;
    el.textContent = t(key);
  });
  document.querySelectorAll<HTMLElement>("[data-i18n-placeholder]").forEach((el) => {
    const key = el.getAttribute("data-i18n-placeholder")!;
    (el as HTMLInputElement).placeholder = t(key);
  });
  document.querySelectorAll<HTMLElement>("[data-i18n-title]").forEach((el) => {
    const key = el.getAttribute("data-i18n-title")!;
    el.title = t(key);
  });
  document.documentElement.lang = current;
}

export function initI18n(): void {
  applyTranslations();
}
