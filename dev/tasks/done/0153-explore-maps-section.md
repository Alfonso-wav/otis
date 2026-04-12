# Feature: sección "Mapas" en Explorar con visor de regiones y mapas especiales

**ID**: 0153-explore-maps-section
**Estado**: done
**Fecha**: 2026-04-12

---

## Descripción

Crear una nueva pestaña **"Maps" / "Mapas"** dentro de la página Explorar. La sección tendrá tres sub-secciones internas:

1. **Mapas de regiones** — visor individual de mapas de cada región del mundo Pokémon. El usuario selecciona la región desde un selector (Kanto por defecto). Las imágenes se obtienen de la wiki de Pokémon: https://pokemon.fandom.com/wiki/Region
2. **Mapa del mundo** — imagen del mapa completo del mundo Pokémon (incluyendo Paldea y Lental). Fuente: https://www.reddit.com/r/PokemonMaps/comments/z82a94/pokemon_world_map_now_with_paldea_and_lental/#lightbox
3. **Mapa retro (Game Boy)** — imagen del mapa clásico de Pokémon Blue. Fuente: https://blog.vjeux.com/wp-content/uploads/2023/12/pokemon_blue-1.png

Cada imagen debe ser visualizable con zoom (pinch-to-zoom en mobile, scroll en desktop) y pan para navegar por los detalles del mapa.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — nuevo tab en Explorar, componente visor de imágenes, estilos, i18n

## Obtención de imágenes

Las imágenes de las regiones se deben descargar manualmente y guardar en `frontend/public/assets/maps/`:

| Región | Archivo | Fuente |
|--------|---------|--------|
| Kanto | `kanto.png` | https://pokemon.fandom.com/wiki/Region (sección Kanto) |
| Johto | `johto.png` | https://pokemon.fandom.com/wiki/Region (sección Johto) |
| Hoenn | `hoenn.png` | https://pokemon.fandom.com/wiki/Region (sección Hoenn) |
| Sinnoh | `sinnoh.png` | https://pokemon.fandom.com/wiki/Region (sección Sinnoh) |
| Unova | `unova.png` | https://pokemon.fandom.com/wiki/Region (sección Unova) |
| Kalos | `kalos.png` | https://pokemon.fandom.com/wiki/Region (sección Kalos) |
| Alola | `alola.png` | https://pokemon.fandom.com/wiki/Region (sección Alola) |
| Galar | `galar.png` | https://pokemon.fandom.com/wiki/Region (sección Galar) |
| Paldea | `paldea.png` | https://pokemon.fandom.com/wiki/Region (sección Paldea) |
| Hisui | `hisui.png` | https://pokemon.fandom.com/wiki/Region (sección Hisui) |
| Mundo | `world.png` | Reddit (mapa completo con Paldea y Lental) |
| Retro | `retro-blue.png` | https://blog.vjeux.com/wp-content/uploads/2023/12/pokemon_blue-1.png |

> **Nota**: Las imágenes deben descargarse manualmente y colocarse en la carpeta antes de implementar. No se deben hotlinkear desde las fuentes externas.

## Archivos a crear/modificar

| Archivo | Acción | Descripción |
|---------|--------|-------------|
| `frontend/public/assets/maps/` | crear dir | Directorio para las imágenes de mapas |
| `frontend/src/pages/explore/maps.ts` | crear | Componente principal de la sección Mapas |
| `frontend/src/pages/explore.ts` | modificar | Añadir tab "maps" al tipo `ExploreTab`, `TAB_KEYS`, icono, switch case |
| `frontend/src/locales/en.json` | modificar | Añadir claves `explore.tabs.maps` y sección `maps.*` |
| `frontend/src/locales/es.json` | modificar | Añadir claves `explore.tabs.maps` y sección `maps.*` |
| `frontend/src/styles/_explore.scss` | modificar | Estilos del visor de mapas (contenedor, zoom, selector) |
| `frontend/src/styles/_dark.scss` | modificar | Dark mode para la sección de mapas |

## Plan de implementación

1. **Crear directorio de assets** `frontend/public/assets/maps/` y colocar las imágenes descargadas manualmente.
2. **Crear `maps.ts`**:
   - Definir sub-secciones internas: "Regions", "World", "Retro".
   - Sub-sección **Regions**: selector dropdown con las regiones disponibles (Kanto por defecto). Al cambiar la selección se muestra la imagen correspondiente.
   - Sub-sección **World**: imagen fija del mapa del mundo.
   - Sub-sección **Retro**: imagen fija del mapa Game Boy.
   - Cada imagen dentro de un contenedor con zoom (CSS `transform: scale()` + drag/pan) o usar una librería ligera tipo `panzoom`.
   - Loading overlay mientras carga cada imagen.
3. **Registrar en `explore.ts`**:
   - Añadir `"maps"` al tipo `ExploreTab`.
   - Añadir a `TAB_KEYS`.
   - Crear icono SVG (globo terráqueo o similar).
   - Añadir case en `initPanel`.
   - Añadir entrada en `tabInited`.
4. **i18n**: añadir traducciones en ambos idiomas.
5. **Estilos**: contenedor del visor con aspect ratio, controles de zoom, selector estilizado, dark mode.

## Traducciones necesarias

**EN:**
```json
"explore.tabs.maps": "Maps",
"maps.title": "Maps",
"maps.regions": "Regions",
"maps.world": "World Map",
"maps.retro": "Retro Map",
"maps.selectRegion": "Select a region",
"maps.zoomIn": "Zoom in",
"maps.zoomOut": "Zoom out",
"maps.reset": "Reset"
```

**ES:**
```json
"explore.tabs.maps": "Mapas",
"maps.title": "Mapas",
"maps.regions": "Regiones",
"maps.world": "Mapa del Mundo",
"maps.retro": "Mapa Retro",
"maps.selectRegion": "Selecciona una región",
"maps.zoomIn": "Acercar",
"maps.zoomOut": "Alejar",
"maps.reset": "Restablecer"
```

## Tests

| Archivo | Qué se testea |
|---------|---------------|
| Manual | Tab "Maps" aparece en la barra de Explorar |
| Manual | Selector de regiones muestra todas las regiones y carga la imagen correcta |
| Manual | Kanto se muestra por defecto al entrar |
| Manual | Zoom (scroll wheel y botones) funciona correctamente |
| Manual | Pan/drag funciona para navegar la imagen |
| Manual | Sub-sección "World Map" muestra el mapa del mundo |
| Manual | Sub-sección "Retro Map" muestra el mapa de Game Boy |
| Manual | Dark mode correcto en todas las sub-secciones |
| Manual | Responsive en mobile (pinch-to-zoom) |
| Manual | Cambio de idioma actualiza las etiquetas |

## Criterios de aceptación

- [ ] Nueva pestaña "Maps" / "Mapas" visible en Explorar
- [ ] Sub-sección de regiones con selector dropdown (Kanto por defecto)
- [ ] Imágenes de las 10 regiones principales se cargan correctamente
- [ ] Sub-sección con mapa del mundo completo
- [ ] Sub-sección con mapa retro de Game Boy
- [ ] Zoom y pan funcional en todas las imágenes
- [ ] Dark mode completo
- [ ] Responsive (mobile y desktop)
- [ ] Textos i18n en ES y EN
- [ ] Imágenes servidas localmente desde `assets/maps/`, no hotlinkeadas
