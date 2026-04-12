import "./styles/main.scss";
import gsap from "gsap";
import { registerPage, initRouter } from "./router";
import { initPokedex } from "./pages/pokedex";
import { initExplore } from "./pages/explore";
import { initBuilds } from "./pages/builds";
import { initSettings } from "./settings";
import { initI18n } from "./i18n";
import { ListGenerations } from "./api";

// ── Arrow helpers ──────────────────────────────────────────────
const ARROW_POSITIONS = ["n", "ne", "se", "s", "sw", "nw"] as const;
let arrowInterval: ReturnType<typeof setInterval> | null = null;

function createArrowElements(): void {
  const container = document.getElementById("splash-arrows");
  if (!container) return;
  for (const pos of ARROW_POSITIONS) {
    const el = document.createElement("div");
    el.className = `splash-arrow splash-arrow--${pos}`;
    el.dataset.pos = pos;
    container.appendChild(el);
  }
}

/** Pick 1-3 random unique positions from the 8 compass points. */
function pickRandomPositions(): string[] {
  const count = Math.floor(Math.random() * 3) + 1; // 1, 2, or 3
  const shuffled = [...ARROW_POSITIONS].sort(() => Math.random() - 0.5);
  return shuffled.slice(0, count);
}

function startArrowLoop(): void {
  const container = document.getElementById("splash-arrows");
  if (!container) return;
  const arrows = container.querySelectorAll<HTMLElement>(".splash-arrow");

  function tick(): void {
    const active = pickRandomPositions();
    arrows.forEach((el) => {
      const pos = el.dataset.pos ?? "";
      if (active.includes(pos)) {
        gsap.fromTo(el,
          { opacity: 0, scale: 0.4 },
          { opacity: 1, scale: 1, duration: 0.35, ease: "back.out(1.7)" },
        );
      } else {
        gsap.to(el, { opacity: 0, scale: 0.4, duration: 0.25 });
      }
    });
  }

  tick(); // first batch immediately
  arrowInterval = setInterval(tick, 1000);
}

function stopArrows(): void {
  if (arrowInterval !== null) {
    clearInterval(arrowInterval);
    arrowInterval = null;
  }
  const arrows = document.querySelectorAll<HTMLElement>(".splash-arrow");
  arrows.forEach((el) => gsap.to(el, { opacity: 0, duration: 0.2 }));
}

// ── Eyes animation ─────────────────────────────────────────────
function animateEyesOpen(): Promise<void> {
  return new Promise((resolve) => {
    const overlay = document.querySelector<SVGElement>(".splash-eyes-overlay");
    if (!overlay) { resolve(); return; }

    const clipLeft = document.querySelector("#eye-clip-left ellipse") as SVGEllipseElement | null;
    const clipRight = document.querySelector("#eye-clip-right ellipse") as SVGEllipseElement | null;
    if (!clipLeft || !clipRight) { resolve(); return; }

    // All eye shapes (inside the clipped groups) scale proportionally
    const eyeShapes = overlay.querySelectorAll("g[clip-path] ellipse");

    // Show overlay
    gsap.set(overlay, { opacity: 1 });
    gsap.set(eyeShapes, { attr: { transform: "scale(0.5)" }, transformOrigin: "center" });

    const tl = gsap.timeline({ onComplete: resolve });

    // Open clip + scale shapes simultaneously
    tl.to([clipLeft, clipRight], {
      attr: { ry: 10 },
      duration: 0.8,
      ease: "power2.out",
    }, 0);
    tl.to(eyeShapes, {
      attr: { transform: "scale(1)" },
      duration: 0.8,
      ease: "power2.out",
    }, 0);
  });
}

// ── Splash flow ────────────────────────────────────────────────
function showSplashInteractive(): void {
  const splash = document.getElementById("splash-screen");
  const zzz = splash?.querySelector(".splash-zzz");
  const wrapper = splash?.querySelector(".splash-snorlax-wrapper") as HTMLElement | null;
  const snorlax = splash?.querySelector(".splash-snorlax") as HTMLElement | null;
  if (!splash || !snorlax || !wrapper) return;

  if (zzz) gsap.to(zzz, { opacity: 0, duration: 0.4 });
  wrapper.style.cursor = "pointer";

  // Create arrow elements and start loop after 2 seconds
  createArrowElements();
  setTimeout(() => startArrowLoop(), 2000);

  wrapper.addEventListener("click", () => dismissSplashInteractive(splash, wrapper), { once: true });
}

async function dismissSplashInteractive(splash: HTMLElement, wrapper: HTMLElement): Promise<void> {
  wrapper.style.pointerEvents = "none";

  // 1. Stop arrows
  stopArrows();

  // 2. Quick jiggle reaction
  const tl = gsap.timeline();
  tl.to(wrapper, { scale: 1.15, duration: 0.1 })
    .to(wrapper, {
      keyframes: [
        { rotation: -8 }, { rotation: 8 },
        { rotation: -5 }, { rotation: 5 },
        { rotation: 0 },
      ],
      duration: 0.45,
      ease: "power2.out",
    }, "<0.05")
    .to(wrapper, { scale: 1, duration: 0.1 });

  // Wait for jiggle to finish
  await tl.then();

  // 3. Exit animation — grow Snorlax and fade everything
  wrapper.style.animation = "none";
  gsap.to(wrapper, {
    scale: 3,
    opacity: 0,
    duration: 2.2,
    ease: "power1.in",
  });
  gsap.to(splash, {
    opacity: 0,
    duration: 3.5,
    ease: "power1.inOut",
    onComplete: () => splash.remove(),
  });
}

// Ping the API to detect when the server is ready, then show interactive splash.
ListGenerations()
  .then(() => showSplashInteractive())
  .catch(() => showSplashInteractive());

registerPage({
  id: "pokedex",
  container: document.getElementById("tab-pokedex") as HTMLElement,
});
registerPage({
  id: "explore",
  container: document.getElementById("tab-explore") as HTMLElement,
});
registerPage({
  id: "builds",
  container: document.getElementById("tab-builds") as HTMLElement,
});
registerPage({
  id: "settings",
  container: document.getElementById("tab-settings") as HTMLElement,
});

initRouter("pokedex");

initI18n();
initSettings();
initPokedex();
initExplore();
initBuilds();
