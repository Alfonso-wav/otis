import gsap from "gsap";
import { initRegions } from "./explore/regions";
import { initMoves } from "./explore/moves";
import { initAbilities } from "./explore/abilities";

type ExploreTab = "regions" | "moves" | "abilities";

let initialized = false;
let activeTab: ExploreTab = "regions";

const ICON_MAP = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M9 6.75V15m6-6v8.25m.503 3.498 4.875-2.437c.381-.19.622-.58.622-1.006V4.82c0-.836-.88-1.38-1.628-1.006l-3.869 1.934c-.317.159-.69.159-1.006 0L9.503 3.252a1.125 1.125 0 0 0-1.006 0L3.622 5.689C3.24 5.88 3 6.27 3 6.695V19.18c0 .836.88 1.38 1.628 1.006l3.869-1.934c.317-.159.69-.159 1.006 0l4.994 2.497c.317.158.69.158 1.006 0Z" /></svg>`;
const ICON_BOLT = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z" /></svg>`;
const ICON_SPARKLES = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09ZM18.259 8.715 18 9.75l-.259-1.035a3.375 3.375 0 0 0-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 0 0 2.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 0 0 2.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 0 0-2.456 2.456ZM16.894 20.567 16.5 21.75l-.394-1.183a2.25 2.25 0 0 0-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 0 0 1.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 0 0 1.423 1.423l1.183.394-1.183.394a2.25 2.25 0 0 0-1.423 1.423Z" /></svg>`;

const tabLabels: Record<ExploreTab, string> = {
  regions: `${ICON_MAP} Regiones`,
  moves: `${ICON_BOLT} Movimientos`,
  abilities: `${ICON_SPARKLES} Habilidades`,
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
