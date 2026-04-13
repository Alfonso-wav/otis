import gsap from "gsap";
import { t } from "../../i18n";

type MapSubSection = "regions" | "world" | "retro";

interface PanZoomState {
  scale: number;
  x: number;
  y: number;
  dragging: boolean;
  startX: number;
  startY: number;
  lastX: number;
  lastY: number;
  initialPinchDist: number;
  initialPinchScale: number;
}

const REGION_FILES: { id: string; ext: string }[] = [
  { id: "kanto", ext: "png" },
  { id: "johto", ext: "png" },
  { id: "hoenn", ext: "png" },
  { id: "sinnoh", ext: "png" },
  { id: "unova", ext: "jpg" },
  { id: "kalos", ext: "jpg" },
  { id: "alola", ext: "png" },
  { id: "galar", ext: "jpg" },
  { id: "paldea", ext: "jpg" },
  { id: "hisui", ext: "png" },
];


const MIN_SCALE = 0.5;
const MAX_SCALE = 5;
const ZOOM_STEP = 0.3;

function createPanZoomState(): PanZoomState {
  return {
    scale: 1, x: 0, y: 0,
    dragging: false, startX: 0, startY: 0, lastX: 0, lastY: 0,
    initialPinchDist: 0, initialPinchScale: 1,
  };
}

function applyTransform(img: HTMLElement, state: PanZoomState): void {
  img.style.transform = `translate(${state.x}px, ${state.y}px) scale(${state.scale})`;
}

function clampScale(scale: number): number {
  return Math.min(MAX_SCALE, Math.max(MIN_SCALE, scale));
}

function pinchDistance(t1: Touch, t2: Touch): number {
  const dx = t1.clientX - t2.clientX;
  const dy = t1.clientY - t2.clientY;
  return Math.sqrt(dx * dx + dy * dy);
}

function attachPanZoom(container: HTMLElement, img: HTMLElement, state: PanZoomState): void {
  // Wheel zoom
  container.addEventListener("wheel", (e) => {
    e.preventDefault();
    const delta = e.deltaY > 0 ? -ZOOM_STEP : ZOOM_STEP;
    state.scale = clampScale(state.scale + delta);
    applyTransform(img, state);
  }, { passive: false });

  // Mouse drag
  container.addEventListener("mousedown", (e) => {
    if (e.button !== 0) return;
    state.dragging = true;
    state.startX = e.clientX - state.x;
    state.startY = e.clientY - state.y;
    container.style.cursor = "grabbing";
  });

  window.addEventListener("mousemove", (e) => {
    if (!state.dragging) return;
    state.x = e.clientX - state.startX;
    state.y = e.clientY - state.startY;
    applyTransform(img, state);
  });

  window.addEventListener("mouseup", () => {
    if (!state.dragging) return;
    state.dragging = false;
    container.style.cursor = "grab";
  });

  // Touch pan & pinch-to-zoom
  container.addEventListener("touchstart", (e) => {
    if (e.touches.length === 1) {
      state.dragging = true;
      state.startX = e.touches[0].clientX - state.x;
      state.startY = e.touches[0].clientY - state.y;
    } else if (e.touches.length === 2) {
      state.dragging = false;
      state.initialPinchDist = pinchDistance(e.touches[0], e.touches[1]);
      state.initialPinchScale = state.scale;
    }
  }, { passive: true });

  container.addEventListener("touchmove", (e) => {
    e.preventDefault();
    if (e.touches.length === 1 && state.dragging) {
      state.x = e.touches[0].clientX - state.startX;
      state.y = e.touches[0].clientY - state.startY;
      applyTransform(img, state);
    } else if (e.touches.length === 2) {
      const dist = pinchDistance(e.touches[0], e.touches[1]);
      const ratio = dist / state.initialPinchDist;
      state.scale = clampScale(state.initialPinchScale * ratio);
      applyTransform(img, state);
    }
  }, { passive: false });

  container.addEventListener("touchend", () => {
    state.dragging = false;
  });
}

function resetView(img: HTMLElement, state: PanZoomState): void {
  state.scale = 1;
  state.x = 0;
  state.y = 0;
  applyTransform(img, state);
}

function buildViewer(src: string): { wrapper: HTMLElement; state: PanZoomState } {
  const state = createPanZoomState();

  const wrapper = document.createElement("div");
  wrapper.className = "map-viewer";

  const viewport = document.createElement("div");
  viewport.className = "map-viewport";

  const loader = document.createElement("div");
  loader.className = "map-loader";
  loader.innerHTML = `<div class="map-spinner"></div>`;

  const img = document.createElement("img");
  img.className = "map-img";
  img.alt = "Map";
  img.draggable = false;

  img.addEventListener("load", () => {
    loader.classList.add("hidden");
    gsap.fromTo(img, { opacity: 0 }, { opacity: 1, duration: 0.3 });
  });
  img.addEventListener("error", () => {
    loader.innerHTML = `<span class="map-error">Image not available</span>`;
  });
  img.src = src;

  viewport.appendChild(loader);
  viewport.appendChild(img);
  wrapper.appendChild(viewport);

  attachPanZoom(viewport, img, state);

  // Zoom controls
  const controls = document.createElement("div");
  controls.className = "map-controls";
  controls.innerHTML = `
    <button class="map-ctrl-btn map-ctrl-zoomin" title="${t("maps.zoomIn")}">+</button>
    <button class="map-ctrl-btn map-ctrl-zoomout" title="${t("maps.zoomOut")}">-</button>
    <button class="map-ctrl-btn map-ctrl-reset" title="${t("maps.reset")}">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" width="14" height="14"><path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.992 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182" /></svg>
    </button>`;

  controls.querySelector(".map-ctrl-zoomin")!.addEventListener("click", () => {
    state.scale = clampScale(state.scale + ZOOM_STEP);
    applyTransform(img, state);
  });
  controls.querySelector(".map-ctrl-zoomout")!.addEventListener("click", () => {
    state.scale = clampScale(state.scale - ZOOM_STEP);
    applyTransform(img, state);
  });
  controls.querySelector(".map-ctrl-reset")!.addEventListener("click", () => {
    resetView(img, state);
  });

  wrapper.appendChild(controls);

  return { wrapper, state };
}

function buildRegionsSection(): HTMLElement {
  const section = document.createElement("div");
  section.className = "map-section map-section--regions";

  let currentRegion = REGION_FILES[0];

  // Region selector
  const selectorRow = document.createElement("div");
  selectorRow.className = "map-selector-row";

  const label = document.createElement("label");
  label.className = "map-selector-label";
  label.textContent = t("maps.selectRegion");

  const select = document.createElement("select");
  select.className = "map-region-select";
  REGION_FILES.forEach((r) => {
    const opt = document.createElement("option");
    opt.value = r.id;
    opt.textContent = r.id.charAt(0).toUpperCase() + r.id.slice(1);
    if (r.id === currentRegion.id) opt.selected = true;
    select.appendChild(opt);
  });

  selectorRow.appendChild(label);
  selectorRow.appendChild(select);
  section.appendChild(selectorRow);

  // Viewer
  const viewerContainer = document.createElement("div");
  viewerContainer.className = "map-viewer-container";
  section.appendChild(viewerContainer);

  function loadRegion(region: typeof REGION_FILES[number]): void {
    viewerContainer.innerHTML = "";
    const { wrapper } = buildViewer(`/assets/maps/${region.id}.${region.ext}`);
    viewerContainer.appendChild(wrapper);
    gsap.fromTo(wrapper, { opacity: 0, y: 6 }, { opacity: 1, y: 0, duration: 0.25 });
  }

  select.addEventListener("change", () => {
    currentRegion = REGION_FILES.find((r) => r.id === select.value) ?? REGION_FILES[0];
    loadRegion(currentRegion);
  });

  loadRegion(currentRegion);
  return section;
}

const WORLD_MAPS = [
  { file: "world-1.png", label: "1" },
  { file: "world-2.png", label: "2" },
  { file: "world-3.png", label: "3" },
  { file: "world-4.png", label: "4" },
];

function buildWorldSection(): HTMLElement {
  const section = document.createElement("div");
  section.className = "map-section map-section--world";

  // Gallery selector
  const selectorRow = document.createElement("div");
  selectorRow.className = "map-gallery-nav";

  WORLD_MAPS.forEach((m, i) => {
    const btn = document.createElement("button");
    btn.className = `map-gallery-btn${i === 0 ? " active" : ""}`;
    btn.textContent = m.label;
    btn.dataset.idx = String(i);
    selectorRow.appendChild(btn);
  });

  section.appendChild(selectorRow);

  const viewerContainer = document.createElement("div");
  viewerContainer.className = "map-viewer-container";
  section.appendChild(viewerContainer);

  function loadWorld(idx: number): void {
    viewerContainer.innerHTML = "";
    const { wrapper } = buildViewer(`/assets/maps/${WORLD_MAPS[idx].file}`);
    viewerContainer.appendChild(wrapper);
    gsap.fromTo(wrapper, { opacity: 0, y: 6 }, { opacity: 1, y: 0, duration: 0.25 });
    selectorRow.querySelectorAll<HTMLButtonElement>(".map-gallery-btn").forEach((b) => {
      b.classList.toggle("active", Number(b.dataset.idx) === idx);
    });
  }

  selectorRow.querySelectorAll<HTMLButtonElement>(".map-gallery-btn").forEach((btn) => {
    btn.addEventListener("click", () => loadWorld(Number(btn.dataset.idx)));
  });

  loadWorld(0);
  return section;
}

// ─── Exposed list for reuse (e.g. battle background selector) ────────────────
export interface BattleBgMap {
  id: string;
  file: string;
  labelKey: string;
}

export const BATTLE_BG_MAPS: BattleBgMap[] = [
  ...REGION_FILES.map<BattleBgMap>((r) => ({
    id: r.id,
    file: `${r.id}.${r.ext}`,
    labelKey: `regions.labels.${r.id}`,
  })),
  ...WORLD_MAPS.map<BattleBgMap>((m, i) => ({
    id: `world-${i + 1}`,
    file: m.file,
    labelKey: `maps.bgLabels.world${i + 1}`,
  })),
  { id: "retro-blue", file: "retro-blue.png", labelKey: "maps.bgLabels.retroBlue" },
];

function buildRetroSection(): HTMLElement {
  const section = document.createElement("div");
  section.className = "map-section map-section--retro";
  const { wrapper } = buildViewer("/assets/maps/retro-blue.png");
  section.appendChild(wrapper);
  return section;
}

export function initMaps(panel: HTMLElement): void {
  let activeSubSection: MapSubSection = "regions";

  const subTabs: { key: MapSubSection; labelKey: string }[] = [
    { key: "regions", labelKey: "maps.regions" },
    { key: "world", labelKey: "maps.world" },
    { key: "retro", labelKey: "maps.retro" },
  ];

  // Sub-navigation
  const subNav = document.createElement("div");
  subNav.className = "map-sub-nav";

  subTabs.forEach(({ key, labelKey }) => {
    const btn = document.createElement("button");
    btn.className = `map-sub-btn${key === activeSubSection ? " active" : ""}`;
    btn.dataset.mapSub = key;
    btn.textContent = t(labelKey);
    subNav.appendChild(btn);
  });

  panel.appendChild(subNav);

  // Content container
  const content = document.createElement("div");
  content.className = "map-content";
  panel.appendChild(content);

  // Build sections lazily
  const sections: Partial<Record<MapSubSection, HTMLElement>> = {};

  function showSubSection(key: MapSubSection): void {
    if (activeSubSection === key && sections[key]) return;
    activeSubSection = key;

    subNav.querySelectorAll<HTMLButtonElement>(".map-sub-btn").forEach((btn) => {
      btn.classList.toggle("active", btn.dataset.mapSub === key);
    });

    if (!sections[key]) {
      switch (key) {
        case "regions": sections[key] = buildRegionsSection(); break;
        case "world": sections[key] = buildWorldSection(); break;
        case "retro": sections[key] = buildRetroSection(); break;
      }
    }

    content.innerHTML = "";
    content.appendChild(sections[key]!);
    gsap.fromTo(content, { opacity: 0, y: 8 }, { opacity: 1, y: 0, duration: 0.25 });
  }

  subNav.querySelectorAll<HTMLButtonElement>(".map-sub-btn").forEach((btn) => {
    btn.addEventListener("click", () => showSubSection(btn.dataset.mapSub as MapSubSection));
  });

  showSubSection("regions");

  // Locale change: update labels
  document.addEventListener("locale-changed", () => {
    subNav.querySelectorAll<HTMLButtonElement>(".map-sub-btn").forEach((btn) => {
      const key = btn.dataset.mapSub as MapSubSection;
      const tab = subTabs.find((st) => st.key === key);
      if (tab) btn.textContent = t(tab.labelKey);
    });

    // Update selector label if regions section exists
    const selectorLabel = panel.querySelector<HTMLElement>(".map-selector-label");
    if (selectorLabel) selectorLabel.textContent = t("maps.selectRegion");

    // Update control tooltips
    panel.querySelectorAll<HTMLButtonElement>(".map-ctrl-zoomin").forEach((b) => b.title = t("maps.zoomIn"));
    panel.querySelectorAll<HTMLButtonElement>(".map-ctrl-zoomout").forEach((b) => b.title = t("maps.zoomOut"));
    panel.querySelectorAll<HTMLButtonElement>(".map-ctrl-reset").forEach((b) => b.title = t("maps.reset"));
  });
}
