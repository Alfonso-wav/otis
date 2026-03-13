import "./styles/main.scss";
import { registerPage, initRouter } from "./router";
import { initPokedex } from "./pages/pokedex";
import { initTypes } from "./pages/types";
import { initExplore } from "./pages/explore";
import { initCompare } from "./pages/compare";

registerPage({
  id: "pokedex",
  container: document.getElementById("tab-pokedex") as HTMLElement,
});
registerPage({
  id: "types",
  container: document.getElementById("tab-types") as HTMLElement,
});
registerPage({
  id: "compare",
  container: document.getElementById("tab-compare") as HTMLElement,
});
registerPage({
  id: "explore",
  container: document.getElementById("tab-explore") as HTMLElement,
});

initRouter("pokedex");

initPokedex();
initTypes();
initCompare();
initExplore();
