import gsap from "gsap";

export interface Page {
  id: string;
  container: HTMLElement;
}

const pages: Page[] = [];
let activeId: string | null = null;

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

export function initRouter(defaultId: string): void {
  activeId = defaultId;

  document.querySelectorAll<HTMLButtonElement>(".tab-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const tabId = btn.dataset.tab;
      if (tabId) navigate(tabId);
    });
  });
}
