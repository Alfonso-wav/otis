import "./styles/main.scss";
import { registerPage, initRouter } from "./router";
import { initPokedex } from "./pages/pokedex";
import { initExplore } from "./pages/explore";
import { initBuilds } from "./pages/builds";
import { initSettings } from "./settings";
import { initI18n } from "./i18n";
import { ListGenerations } from "./api";

function showSplashInteractive(): void {
  const splash = document.getElementById("splash-screen");
  const zzz = splash?.querySelector(".splash-zzz");
  const snorlax = splash?.querySelector(".splash-snorlax") as HTMLElement | null;
  const hint = splash?.querySelector(".splash-hint") as HTMLElement | null;
  if (!splash || !snorlax) return;

  gsap.to(zzz, { opacity: 0, duration: 0.4 });
  gsap.to(hint, { opacity: 1, duration: 0.4, delay: 0.3 });
  snorlax.style.cursor = "pointer";

  snorlax.addEventListener("click", () => dismissSplashInteractive(splash, snorlax), { once: true });
}

function dismissSplashInteractive(splash: HTMLElement, snorlax: HTMLElement): void {
  snorlax.style.pointerEvents = "none";

  gsap.timeline()
    .to(snorlax, { scale: 1.15, duration: 0.1 })
    .to(snorlax, {
      keyframes: [
        { rotation: -8 }, { rotation: 8 },
        { rotation: -5 }, { rotation: 5 },
        { rotation: 0 }
      ],
      duration: 0.45,
      ease: "power2.out"
    }, "<0.05")
    .to(snorlax, { scale: 1, duration: 0.1 })
    .call(() => {
      const overlay = document.createElement("div");
      overlay.className = "splash-exit-overlay";
      splash.appendChild(overlay);
      gsap.fromTo(overlay,
        { clipPath: "circle(0% at 50% 50%)" },
        {
          clipPath: "circle(150% at 50% 50%)",
          duration: 0.55,
          ease: "power2.in",
          onComplete: () => {
            splash.remove();
          }
        }
      );
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
