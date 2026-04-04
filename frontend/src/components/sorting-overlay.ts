import gsap from "gsap";
import { t } from "../i18n";

const MRMIME_SPRITE = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/122.png";
const MACHAMP_SPRITE = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/68.png";

let overlayEl: HTMLDivElement | null = null;
let tween: gsap.core.Tween | null = null;

function create(sprite: string, altText: string): HTMLDivElement {
  const el = document.createElement("div");
  el.className = "sorting-overlay";
  el.innerHTML = `
    <div class="sorting-overlay__content">
      <img class="sorting-overlay__img" src="${sprite}" alt="${altText}" />
      <p class="sorting-overlay__text">${t("sorting.sorting")}</p>
    </div>`;
  return el;
}

export function showSortingOverlay(text?: string): void {
  hideSortingOverlay();
  overlayEl = create(MRMIME_SPRITE, "Mr. Mime thinking...");
  if (text) {
    overlayEl.querySelector<HTMLParagraphElement>(".sorting-overlay__text")!.textContent = text;
  }
  document.body.appendChild(overlayEl);

  const img = overlayEl.querySelector<HTMLImageElement>(".sorting-overlay__img")!;
  tween = gsap.to(img, {
    y: -10,
    duration: 0.6,
    ease: "sine.inOut",
    yoyo: true,
    repeat: -1,
  });
}

export function updateSortingOverlayText(text: string): void {
  if (!overlayEl) return;
  overlayEl.querySelector<HTMLParagraphElement>(".sorting-overlay__text")!.textContent = text;
}

export function showLoadingOverlay(text?: string): void {
  hideSortingOverlay();
  overlayEl = create(MACHAMP_SPRITE, "Machamp loading...");
  if (text) {
    overlayEl.querySelector<HTMLParagraphElement>(".sorting-overlay__text")!.textContent = text;
  }
  document.body.appendChild(overlayEl);

  const img = overlayEl.querySelector<HTMLImageElement>(".sorting-overlay__img")!;
  tween = gsap.to(img, {
    y: -10,
    duration: 0.6,
    ease: "sine.inOut",
    yoyo: true,
    repeat: -1,
  });
}

export function hideSortingOverlay(): void {
  if (tween) {
    tween.kill();
    tween = null;
  }
  if (overlayEl) {
    overlayEl.remove();
    overlayEl = null;
  }
}
