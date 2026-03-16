import gsap from "gsap";

export interface Page {
  id: string;
  container: HTMLElement;
}

const pages: Page[] = [];
let activeId: string | null = null;
let previousId: string | null = null;
const isMobile = window.matchMedia("(max-width: 768px)").matches;

export function registerPage(page: Page): void {
  pages.push(page);
}

export function navigate(id: string): void {
  if (id === activeId) return;

  const next = pages.find((p) => p.id === id);
  if (!next) return;

  const current = pages.find((p) => p.id === activeId);

  document.querySelectorAll<HTMLButtonElement>(".tab-btn").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.tab === id);
  });

  activeId = id;

  if (!current) {
    return;
  }

  if (isMobile) {
    // Skip animations on mobile for snappy tab switching
    current.container.classList.add("hidden");
    next.container.classList.remove("hidden");
  } else {
    gsap.to(current.container, {
      opacity: 0,
      y: -8,
      duration: 0.18,
      ease: "power2.in",
      onComplete() {
        current.container.classList.add("hidden");
        current.container.style.removeProperty("opacity");
        current.container.style.removeProperty("transform");

        next.container.classList.remove("hidden");
        gsap.fromTo(
          next.container,
          { opacity: 0, y: 10 },
          { opacity: 1, y: 0, duration: 0.25, ease: "power2.out" },
        );
      },
    });
  }
}

export function navigateSettings(): void {
  previousId = activeId;
  const settingsPage = document.getElementById("tab-settings");
  const tabNav = document.getElementById("tab-nav");
  if (!settingsPage) return;

  // Hide all tab pages
  pages.forEach((p) => p.container.classList.add("hidden"));
  settingsPage.classList.remove("hidden");

  // Hide tab nav
  if (tabNav) tabNav.classList.add("hidden");
}

export function navigateBack(): void {
  const settingsPage = document.getElementById("tab-settings");
  const tabNav = document.getElementById("tab-nav");

  if (settingsPage) settingsPage.classList.add("hidden");
  if (tabNav) tabNav.classList.remove("hidden");

  const target = previousId ?? "pokedex";
  const page = pages.find((p) => p.id === target);
  if (page) {
    page.container.classList.remove("hidden");
    activeId = target;
  }
}

export function initRouter(defaultId: string): void {
  activeId = defaultId;

  document.querySelectorAll<HTMLButtonElement>(".tab-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const tabId = btn.dataset.tab;
      if (tabId) navigate(tabId);
    });
  });

  document
    .getElementById("settings-btn")
    ?.addEventListener("click", navigateSettings);
  document
    .getElementById("settings-back-btn")
    ?.addEventListener("click", navigateBack);
}
