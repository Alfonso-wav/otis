import gsap from "gsap";
import { ListRegions, GetRegion, GetRegionPokemonByType } from "../../../wailsjs/go/app/App";
import { renderTypeDistributionChart } from "../../charts/type-distribution";
import { openTypeModal } from "../../components/pokemon-type-modal";

let initialized = false;

function regionLabel(name: string): string {
  const labels: Record<string, string> = {
    kanto: "Kanto (Gen I)",
    johto: "Johto (Gen II)",
    hoenn: "Hoenn (Gen III)",
    sinnoh: "Sinnoh (Gen IV)",
    unova: "Unova (Gen V)",
    kalos: "Kalos (Gen VI)",
    alola: "Alola (Gen VII)",
    galar: "Galar (Gen VIII)",
    hisui: "Hisui",
    paldea: "Paldea (Gen IX)",
  };
  return labels[name] ?? name.charAt(0).toUpperCase() + name.slice(1);
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

  body.innerHTML = '<p class="loading">Cargando región...</p>';
  body.classList.remove("hidden");

  try {
    const region = await GetRegion(regionName);
    const locations = region.Locations ?? [];

    const chartId = `chart-region-${regionName}`;
    body.innerHTML = `
      <div class="region-detail">
        <div class="region-locations">
          <h4 class="region-section-title">Localizaciones (${locations.length})</h4>
          <div class="region-locations-grid">
            ${
              locations
                .slice(0, 20)
                .map(
                  (l) =>
                    `<span class="region-location-tag">${l.Name.replace(/-/g, " ")}</span>`,
                )
                .join("") +
              (locations.length > 20
                ? `<span class="region-location-more">+${locations.length - 20} más</span>`
                : "")
            }
          </div>
        </div>
        <div class="region-chart-container">
          <h4 class="region-section-title">Distribución de tipos</h4>
          <div id="${chartId}" class="region-chart"></div>
        </div>
      </div>`;

    body.dataset.loaded = "true";

    const items = body.querySelectorAll(".region-location-tag");
    gsap.fromTo(
      items,
      { opacity: 0, y: 8 },
      { opacity: 1, y: 0, duration: 0.2, stagger: 0.02, ease: "power2.out" },
    );

    renderTypeDistributionChart(chartId, regionName, async (typeName: string) => {
      try {
        const names = await GetRegionPokemonByType(regionName, typeName);
        openTypeModal(regionName, typeName, names ?? []);
      } catch {
        openTypeModal(regionName, typeName, []);
      }
    });
  } catch (err: unknown) {
    body.innerHTML = `<p class="loading" style="color:#e53e3e">${String(err)}</p>`;
  }
}

export async function initRegions(container: HTMLElement): Promise<void> {
  if (initialized) return;
  initialized = true;

  container.innerHTML = '<p class="loading">Cargando regiones...</p>';

  try {
    const regions = await ListRegions();

    container.innerHTML = `
      <div class="section-header"><h2>Regiones</h2></div>
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
}
