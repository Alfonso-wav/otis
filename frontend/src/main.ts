import "./styles/main.scss";
import { registerPage, initRouter } from "./router";
import { initPokedex } from "./pages/pokedex";
import { initExplore } from "./pages/explore";
import { initBuilds } from "./pages/builds";
import { initSettings } from "./settings";
import { ListGenerations } from "./api";

function dismissSplash(): void {
  const splash = document.getElementById("splash-screen");
  if (!splash) return;
  splash.classList.add("splash-fade-out");
  splash.addEventListener("transitionend", () => splash.remove(), { once: true });
}

// Ping the API to detect when the server is ready, then dismiss splash.
ListGenerations()
  .then(() => dismissSplash())
  .catch(() => dismissSplash());

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

initSettings();
initPokedex();
initExplore();
initBuilds();
