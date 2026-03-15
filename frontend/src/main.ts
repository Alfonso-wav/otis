import "./styles/main.scss";
import { registerPage, initRouter } from "./router";
import { initPokedex } from "./pages/pokedex";
import { initTypes } from "./pages/types";
import { initExplore } from "./pages/explore";
import { initBuilds } from "./pages/builds";

registerPage({
  id: "pokedex",
  container: document.getElementById("tab-pokedex") as HTMLElement,
});
registerPage({
  id: "types",
  container: document.getElementById("tab-types") as HTMLElement,
});
registerPage({
  id: "explore",
  container: document.getElementById("tab-explore") as HTMLElement,
});
registerPage({
  id: "builds",
  container: document.getElementById("tab-builds") as HTMLElement,
});

initRouter("pokedex");

initPokedex();
initTypes();
initExplore();
initBuilds();
