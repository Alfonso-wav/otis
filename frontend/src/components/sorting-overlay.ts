import gsap from "gsap";
import { t } from "../i18n";

const MRMIME_SPRITE = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/122.png";
const MACHAMP_SPRITE = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/68.png";
const DIGLETT_SPRITE = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/50.png";

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

export function showDiglettOverlay(text?: string): void {
  hideSortingOverlay();
  const el = document.createElement("div");
  el.className = "sorting-overlay";
  el.innerHTML = `
    <div class="sorting-overlay__content">
      <div class="diglett-hole">
        <div class="diglett-clip">
          <img class="sorting-overlay__img diglett-img" src="${DIGLETT_SPRITE}" alt="Diglett digging..." />
        </div>
        <div class="diglett-ground"></div>
      </div>
      <p class="sorting-overlay__text">${text ?? t("sorting.sorting")}</p>
    </div>`;
  overlayEl = el;
  document.body.appendChild(el);

  const img = el.querySelector<HTMLImageElement>(".diglett-img")!;
  const tl = gsap.timeline({ repeat: -1 });
  tl.fromTo(img, { y: 60 }, { y: 0, duration: 0.4, ease: "back.out(1.7)" })
    .to(img, { y: -6, duration: 0.3, ease: "sine.inOut", yoyo: true, repeat: 3 })
    .to(img, { y: 60, duration: 0.3, ease: "power2.in" })
    .to(img, { y: 60, duration: 0.3 });
  tween = tl as unknown as gsap.core.Tween;
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

// ─── Inline Diglett (contained within a parent element) ─────────────────────

const inlineTweens = new Map<HTMLElement, gsap.core.Tween>();

/**
 * Creates an inline Diglett loading indicator inside the given container.
 * Unlike the fullscreen overlay, this is a contained element that respects
 * the parent's size. Returns the wrapper element for later removal.
 */
export function createInlineDiglett(container: HTMLElement, text?: string): HTMLDivElement {
  // Clean up any orphaned inline tweens (elements removed from DOM by innerHTML replacement)
  for (const [el, tw] of inlineTweens) {
    if (!el.isConnected) {
      tw.kill();
      inlineTweens.delete(el);
    }
  }

  const wrapper = document.createElement("div");
  wrapper.className = "diglett-inline";
  wrapper.innerHTML = `
    <div class="diglett-inline__content">
      <div class="diglett-hole diglett-hole--sm">
        <div class="diglett-clip diglett-clip--sm">
          <img class="sorting-overlay__img diglett-img diglett-img--sm" src="${DIGLETT_SPRITE}" alt="Diglett loading..." />
        </div>
        <div class="diglett-ground diglett-ground--sm"></div>
      </div>
      <p class="diglett-inline__text">${text ?? t("common.loading")}</p>
    </div>`;
  container.innerHTML = "";
  container.appendChild(wrapper);

  const img = wrapper.querySelector<HTMLImageElement>(".diglett-img")!;
  const tl = gsap.timeline({ repeat: -1 });
  tl.fromTo(img, { y: 40 }, { y: 0, duration: 0.4, ease: "back.out(1.7)" })
    .to(img, { y: -4, duration: 0.3, ease: "sine.inOut", yoyo: true, repeat: 3 })
    .to(img, { y: 40, duration: 0.3, ease: "power2.in" })
    .to(img, { y: 40, duration: 0.3 });
  inlineTweens.set(wrapper, tl as unknown as gsap.core.Tween);

  return wrapper;
}

/**
 * Removes an inline Diglett element and kills its animation.
 */
export function removeInlineDiglett(wrapper: HTMLDivElement | null): void {
  if (!wrapper) return;
  const tw = inlineTweens.get(wrapper);
  if (tw) {
    tw.kill();
    inlineTweens.delete(wrapper);
  }
  wrapper.remove();
}
