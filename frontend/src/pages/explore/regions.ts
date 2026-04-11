import gsap from "gsap";
import { ListRegions, GetRegion, GetRegionPokemonByType } from "../../api";
import { renderTypeDistributionChart } from "../../charts/type-distribution";
import { openTypeModal } from "../../components/pokemon-type-modal";
import { openLocationEncounterModal } from "../../components/location-encounter-modal";
import { t, getLocale } from "../../i18n";
import { getLocationNameEs } from "../../data/location-names-es";
import { showDiglettOverlay, hideSortingOverlay } from "../../components/sorting-overlay";

let initialized = false;
let lastContainer: HTMLElement | null = null;

// Track expanded/collapsed state per region
const expandedLocations = new Map<string, boolean>();

const GEN7_PLUS = new Set(["alola", "galar", "hisui", "paldea"]);

function regionLabel(name: string): string {
  const key = `regions.labels.${name}`;
  const label = t(key);
  return label !== key ? label : name.charAt(0).toUpperCase() + name.slice(1);
}

async function loadRegionDetail(
  card: HTMLElement,
  regionName: string,
): Promise<void> {
  const body = card.querySelector<HTMLElement>(".region-card__body")!;
  const isExpanded = card.classList.contains("expanded");

  if (isExpanded) {
    card.classList.remove("expanded");
    gsap.to(body, {
      height: 0,
      opacity: 0,
      duration: 0.25,
      ease: "power2.in",
      onComplete() {
        body.classList.add("hidden");
        body.style.removeProperty("height");
        body.style.removeProperty("opacity");
      },
    });
    return;
  }

  card.classList.add("expanded");

  if (body.dataset.loaded === "true") {
    body.classList.remove("hidden");
    gsap.fromTo(body, { opacity: 0 }, { opacity: 1, duration: 0.2 });
    return;
  }

  body.innerHTML = `<p class="loading">${t("regions.loadingDetail")}</p>`;
  body.classList.remove("hidden");
  showDiglettOverlay(t("regions.loadingDetail"));

  try {
    const region = await GetRegion(regionName);
    const locations = region.Locations ?? [];

    const chartId = `chart-region-${regionName}`;
    const limit = 20;
    const hasMore = locations.length > limit;
    const isExpanded = expandedLocations.get(regionName) ?? false;
    const visibleLocations = isExpanded ? locations : locations.slice(0, limit);

    body.innerHTML = `
      <div class="region-detail">
        <div class="region-locations">
          <h4 class="region-section-title">${t("regions.locations", { count: locations.length })}</h4>
          <div class="region-locations-grid">
            ${visibleLocations
              .map(
                (l) =>
                  `<span class="region-location-tag" data-location="${l.Name}">${l.Names?.[getLocale()] ?? (getLocale() === "es" ? getLocationNameEs(l.Name) : undefined) ?? l.Names?.["en"] ?? l.Name.replace(/-/g, " ")}</span>`,
              )
              .join("")}${
              hasMore
                ? isExpanded
                  ? `<span class="region-location-more region-location-toggle" data-region="${regionName}" data-action="collapse">${t("regions.showLess")}</span>`
                  : `<span class="region-location-more region-location-toggle" data-region="${regionName}" data-action="expand">${t("regions.showMore", { count: locations.length - limit })}</span>`
                : ""
            }
          </div>${GEN7_PLUS.has(regionName) ? `\n          <p class="region-encounter-note">${t("regions.encounterLimited")}</p>` : ""}
        </div>
        <div class="region-chart-container">
          <h4 class="region-section-title">${t("regions.typeDistribution")}</h4>
          <div id="${chartId}" class="region-chart"></div>
        </div>
      </div>`;

    body.dataset.loaded = "true";
    hideSortingOverlay();

    const items = body.querySelectorAll(".region-location-tag");
    gsap.fromTo(
      items,
      { opacity: 0, y: 8 },
      { opacity: 1, y: 0, duration: 0.2, stagger: 0.02, ease: "power2.out" },
    );

    // Click on location tag → open encounter modal (skip Gen VII+ — no data)
    if (GEN7_PLUS.has(regionName)) {
      body.querySelectorAll<HTMLElement>(".region-location-tag").forEach((tag) => {
        tag.classList.add("disabled");
      });
    } else {
      body.querySelectorAll<HTMLElement>(".region-location-tag").forEach((tag) => {
        tag.addEventListener("click", async () => {
          const locName = tag.dataset.location!;
          openLocationEncounterModal(locName);
        });
      });
    }

    // Click on "+X más" / "Mostrar menos" → expand/collapse
    const toggleBtn = body.querySelector<HTMLElement>(".region-location-toggle");
    if (toggleBtn) {
      toggleBtn.addEventListener("click", () => {
        const action = toggleBtn.dataset.action;
        expandedLocations.set(regionName, action === "expand");
        // Re-render by forcing reload
        body.dataset.loaded = "false";
        card.classList.remove("expanded");
        loadRegionDetail(card, regionName);
      });
    }

    renderTypeDistributionChart(chartId, regionName, async (typeName: string) => {
      try {
        const names = await GetRegionPokemonByType(regionName, typeName);
        openTypeModal(regionName, typeName, names ?? []);
      } catch {
        openTypeModal(regionName, typeName, []);
      }
    });
  } catch (err: unknown) {
    hideSortingOverlay();
    body.innerHTML = `<p class="loading" style="color:#e53e3e">${String(err)}</p>`;
  }
}

export async function initRegions(container: HTMLElement): Promise<void> {
  if (initialized) {
    lastContainer = container;
    return;
  }
  initialized = true;
  lastContainer = container;

  container.innerHTML = `<p class="loading">${t("regions.loading")}</p>`;

  try {
    const regions = await ListRegions();

    container.innerHTML = `
      <div class="section-header"><h2>${t("regions.title")}</h2></div>
      <div class="regions-grid" id="regions-grid"></div>`;
    const grid = document.getElementById("regions-grid") as HTMLDivElement;

    grid.innerHTML = regions
      .map(
        (r) => `
      <div class="region-card" data-region="${r.Name}">
        <div class="region-card__header">
          <span class="region-card__name">${regionLabel(r.Name)}</span>
          <span class="region-card__chevron">▼</span>
        </div>
        <div class="region-card__body hidden"></div>
      </div>`,
      )
      .join("");

    const cards = grid.querySelectorAll<HTMLDivElement>(".region-card");
    gsap.fromTo(
      cards,
      { opacity: 0, y: 20 },
      { opacity: 1, y: 0, duration: 0.3, stagger: 0.06, ease: "power2.out" },
    );

    cards.forEach((card) => {
      card.querySelector(".region-card__header")!.addEventListener("click", () =>
        loadRegionDetail(card, card.dataset.region!),
      );
    });
  } catch (err: unknown) {
    container.innerHTML = `<p class="loading" style="color:#e53e3e">${String(err)}</p>`;
  }

  document.addEventListener("locale-changed", () => {
    if (!lastContainer) return;
    initialized = false;
    expandedLocations.clear();
    initRegions(lastContainer);
  });
}
