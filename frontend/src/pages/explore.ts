import gsap from "gsap";
import { initRegions } from "./explore/regions";
import { initMoves } from "./explore/moves";
import { initAbilities } from "./explore/abilities";

type ExploreTab = "regions" | "moves" | "abilities";

let initialized = false;
let activeTab: ExploreTab = "regions";

const tabLabels: Record<ExploreTab, string> = {
  regions: "🗺️ Regiones",
  moves: "⚔️ Movimientos",
  abilities: "✨ Habilidades",
};

const tabInited: Record<ExploreTab, boolean> = {
  regions: false,
  moves: false,
  abilities: false,
};

function switchTab(tab: ExploreTab, container: HTMLElement): void {
  if (activeTab === tab) return;
  activeTab = tab;

  container.querySelectorAll<HTMLButtonElement>(".explore-tab-btn").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.exploreTab === tab);
  });

  const panels = container.querySelectorAll<HTMLElement>(".explore-panel");
  panels.forEach((p) => {
    if (p.dataset.panel === tab) {
      p.classList.remove("hidden");
      gsap.fromTo(p, { opacity: 0, y: 10 }, { opacity: 1, y: 0, duration: 0.25 });
      initPanel(tab, p);
    } else {
      p.classList.add("hidden");
    }
  });
}

function initPanel(tab: ExploreTab, panel: HTMLElement): void {
  if (tabInited[tab]) return;
  tabInited[tab] = true;

  switch (tab) {
    case "regions":
      initRegions(panel);
      break;
    case "moves":
      initMoves(panel);
      break;
    case "abilities":
      initAbilities(panel);
      break;
  }
}

function buildLayout(container: HTMLElement): void {
  container.innerHTML = `
    <div class="explore-nav">
      ${(Object.keys(tabLabels) as ExploreTab[])
        .map(
          (t) =>
            `<button class="explore-tab-btn${t === activeTab ? " active" : ""}" data-explore-tab="${t}">
              ${tabLabels[t]}
            </button>`,
        )
        .join("")}
    </div>
    <div class="explore-content">
      ${(Object.keys(tabLabels) as ExploreTab[])
        .map(
          (t) =>
            `<div class="explore-panel${t !== activeTab ? " hidden" : ""}" data-panel="${t}"></div>`,
        )
        .join("")}
    </div>`;

  container.querySelectorAll<HTMLButtonElement>(".explore-tab-btn").forEach((btn) => {
    btn.addEventListener("click", () =>
      switchTab(btn.dataset.exploreTab as ExploreTab, container),
    );
  });

  // Init first tab immediately
  const firstPanel = container.querySelector<HTMLElement>(
    `[data-panel="${activeTab}"]`,
  )!;
  initPanel(activeTab, firstPanel);
}

export function initExplore(): void {
  const container = document.getElementById("tab-explore") as HTMLElement;

  const tabBtn = document.querySelector<HTMLButtonElement>(
    '[data-tab="explore"]',
  );
  if (!tabBtn) return;

  tabBtn.addEventListener("click", () => {
    if (initialized) return;
    initialized = true;
    buildLayout(container);
    gsap.fromTo(
      container.querySelector(".explore-nav"),
      { opacity: 0, y: -10 },
      { opacity: 1, y: 0, duration: 0.3, ease: "power2.out" },
    );
  });
}
