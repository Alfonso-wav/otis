import "./styles/main.scss";
import gsap from "gsap";
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
  if (!splash || !snorlax) return;

  gsap.to(zzz, { opacity: 0, duration: 0.4 });
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
      // Parar la animación CSS de respiración para que GSAP pueda controlar el transform
      snorlax.style.animation = "none";
      // Snorlax crece y se desvanece lentamente
      gsap.to(snorlax, {
        scale: 3,
        opacity: 0,
        duration: 2.2,
        ease: "power1.in",
      });
      // El fondo oscuro se desvanece muuy lentamente revelando la app
      gsap.to(splash, {
        opacity: 0,
        duration: 3.5,
        ease: "power1.inOut",
        onComplete: () => splash.remove(),
      });
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
