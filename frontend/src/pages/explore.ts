import gsap from "gsap";
import { t } from "../i18n";
import { initTypeChart } from "./explore/type-chart";
import { initRegions } from "./explore/regions";
import { initMoves } from "./explore/moves";
import { initAbilities } from "./explore/abilities";
import { initBerries } from "./explore/berries";

type ExploreTab = "typeChart" | "regions" | "moves" | "abilities" | "berries";

let initialized = false;
let activeTab: ExploreTab = "regions";

const ICON_CHART = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M3.375 19.5h17.25m-17.25 0a1.125 1.125 0 0 1-1.125-1.125M3.375 19.5h1.5C5.496 19.5 6 18.996 6 18.375m-3.75 0V5.625m0 12.75v-1.5c0-.621.504-1.125 1.125-1.125m18.375 2.625V5.625m0 12.75c0 .621-.504 1.125-1.125 1.125m1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125m0 3.75h-1.5A1.125 1.125 0 0 1 18 18.375M20.625 4.5H3.375m17.25 0c.621 0 1.125.504 1.125 1.125M20.625 4.5h-1.5C18.504 4.5 18 5.004 18 5.625m3.75 0v1.5c0 .621-.504 1.125-1.125 1.125M3.375 4.5c-.621 0-1.125.504-1.125 1.125M3.375 4.5h1.5C5.496 4.5 6 5.004 6 5.625m-3.75 0v1.5c0 .621.504 1.125 1.125 1.125m0 0h1.5m-1.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125m1.5-3.75C5.496 8.25 6 8.754 6 9.375v1.5m0-5.25v5.25m0-5.25C6 5.004 6.504 4.5 7.125 4.5h9.75c.621 0 1.125.504 1.125 1.125m1.125 2.625h1.5m-1.5 0A1.125 1.125 0 0 1 18 7.125v1.5m1.5-1.5c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125m-17.25 0h1.5m14.25 0h1.5" /></svg>`;
const ICON_MAP = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M9 6.75V15m6-6v8.25m.503 3.498 4.875-2.437c.381-.19.622-.58.622-1.006V4.82c0-.836-.88-1.38-1.628-1.006l-3.869 1.934c-.317.159-.69.159-1.006 0L9.503 3.252a1.125 1.125 0 0 0-1.006 0L3.622 5.689C3.24 5.88 3 6.27 3 6.695V19.18c0 .836.88 1.38 1.628 1.006l3.869-1.934c.317-.159.69-.159 1.006 0l4.994 2.497c.317.158.69.158 1.006 0Z" /></svg>`;
const ICON_BOLT = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="m3.75 13.5 10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75Z" /></svg>`;
const ICON_SPARKLES = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09ZM18.259 8.715 18 9.75l-.259-1.035a3.375 3.375 0 0 0-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 0 0 2.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 0 0 2.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 0 0-2.456 2.456ZM16.894 20.567 16.5 21.75l-.394-1.183a2.25 2.25 0 0 0-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 0 0 1.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 0 0 1.423 1.423l1.183.394-1.183.394a2.25 2.25 0 0 0-1.423 1.423Z" /></svg>`;
const ICON_LEAF = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="ui-icon" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" d="M12.75 3.03v.568c0 .334.148.65.405.864l1.068.89c.442.369.535 1.01.216 1.49l-.51.766a2.25 2.25 0 0 1-1.161.886l-.143.048a1.107 1.107 0 0 0-.57 1.664c.369.555.169 1.307-.427 1.605L9 13.125l.423 1.059a.956.956 0 0 1-1.652.928l-.679-.906a1.125 1.125 0 0 0-1.906.172L4.5 15.75l-.612.153M12.75 3.031a9 9 0 0 0-8.862 12.872M12.75 3.031a9 9 0 0 1 6.69 14.036l-.001.001-1.174 1.174a2.25 2.25 0 0 1-3.182 0l-3.182-3.182a2.25 2.25 0 0 1 0-3.182l3.182-3.182a2.25 2.25 0 0 1 3.182 0l1.174 1.174a9 9 0 0 1-8.19 14.22" /></svg>`;

function tabLabel(tab: ExploreTab): string {
  const icons: Record<ExploreTab, string> = { typeChart: ICON_CHART, regions: ICON_MAP, moves: ICON_BOLT, abilities: ICON_SPARKLES, berries: ICON_LEAF };
  return `${icons[tab]} ${t(`explore.tabs.${tab}`)}`;
}

const TAB_KEYS: ExploreTab[] = ["regions", "typeChart", "moves", "abilities", "berries"];

const tabInited: Record<ExploreTab, boolean> = {
  typeChart: false,
  regions: false,
  moves: false,
  abilities: false,
  berries: false,
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
    case "typeChart":
      initTypeChart(panel);
      break;
    case "regions":
      initRegions(panel);
      break;
    case "moves":
      initMoves(panel);
      break;
    case "abilities":
      initAbilities(panel);
      break;
    case "berries":
      initBerries(panel);
      break;
  }
}

function buildLayout(container: HTMLElement): void {
  container.innerHTML = `
    <div class="explore-nav">
      ${TAB_KEYS
        .map(
          (tab) =>
            `<button class="explore-tab-btn${tab === activeTab ? " active" : ""}" data-explore-tab="${tab}">
              ${tabLabel(tab)}
            </button>`,
        )
        .join("")}
    </div>
    <div class="explore-content">
      ${TAB_KEYS
        .map(
          (tab) =>
            `<div class="explore-panel${tab !== activeTab ? " hidden" : ""}" data-panel="${tab}"></div>`,
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

  document.addEventListener("locale-changed", () => {
    if (!initialized) return;
    container.querySelectorAll<HTMLButtonElement>(".explore-tab-btn").forEach((btn) => {
      const tab = btn.dataset.exploreTab as ExploreTab;
      if (tab) btn.innerHTML = tabLabel(tab);
    });
  });
}
