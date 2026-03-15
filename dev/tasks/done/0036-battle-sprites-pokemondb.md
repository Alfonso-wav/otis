# Sprites de batalla Gen 1 desde pokemondb.net

**ID**: 0036-battle-sprites-pokemondb
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Completar el simulador de batallas con sprites estilo Game Boy clásico, scrapeando las imágenes de la **Generación 1** desde las páginas de sprites de pokemondb.net (e.g. `https://pokemondb.net/sprites/charmander`).

**Referencia visual**: `docs/battle.png` — layout clásico de batalla Pokémon:
- **Atacante** (abajo-izquierda): sprite de espaldas ("Back") Gen 1
- **Defensor** (arriba-derecha): sprite frontal ("Normal") Gen 1

**Fuente de sprites**:
- Cada Pokémon tiene una página `/sprites/{nombre}` con tabla de generaciones.
- Se necesitan las imágenes de la **Generation 1** (la más antigua disponible):
  - **Back** (atacante): XPath aprox. `/html/body/main/div[11]/table/tbody/tr[2]/td[3]/span[2]/a/img`
  - **Normal** (defensor): XPath aprox. `/html/body/main/div[11]/table/tbody/tr[2]/td[2]/span[2]/a/img`
- Las URLs de las imágenes apuntan al CDN `img.pokemondb.net/sprites/...`
- Si un Pokémon no tiene sprite Gen 1 (ej. generaciones posteriores), usar la generación más antigua disponible.

## Capas afectadas
- **Core**: Nuevas categorías de sprite (`battle-back`, `battle-front`).
- **Shell**: Scraper que visita `/sprites/{nombre}`, parsea la tabla de generaciones con goquery y extrae las URLs de las imágenes Gen 1 Back y Normal. Descarga las imágenes a disco local.
- **APP**: Exponer binding para descargar sprites de batalla (o extender `DownloadSprites` con las nuevas categorías).
- **Frontend**: Rediseñar `renderBattleSection()` en `builds.ts` para mostrar los sprites en layout clásico de batalla (atacante abajo-izq de espaldas, defensor arriba-der de frente) con barras de HP y nombres al estilo Game Boy.

## Dependencias externas nuevas
Ninguna. Se reutiliza `goquery` y `PokemonDBClient` existente con rate limiting.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Agregar `SpriteCategoryBattleBack` y `SpriteCategoryBattleFront` |
| `core/ports.go` | modificar | Extender interfaz si es necesario |
| `shell/pokemondb_sprites.go` | modificar | Agregar scraping de página `/sprites/{nombre}` para extraer URLs de Gen 1 Back/Normal, y descarga a `assets/sprites/battle-back/` y `assets/sprites/battle-front/` |
| `app/bindings.go` | modificar | Soportar nuevas categorías en `DownloadSprites` |
| `frontend/src/pages/builds.ts` | modificar | Rediseñar sección de batalla con sprites, layout clásico GB |
| `frontend/src/styles/_builds.scss` | modificar | Estilos para layout de batalla clásico (sprites, HP bars, nombres) |

## Plan de implementacion

### Fase 1 — Core: nuevas categorías de sprite
1. En `core/domain.go`, agregar constantes:
   - `SpriteCategoryBattleBack SpriteCategory = "battle-back"`
   - `SpriteCategoryBattleFront SpriteCategory = "battle-front"`

### Fase 2 — Shell: scraper de sprites de batalla
2. En `shell/pokemondb_sprites.go`, agregar función que:
   - Reciba un nombre de Pokémon
   - Haga GET a `/sprites/{nombre}` con rate limiting
   - Parsee la tabla de generaciones con goquery
   - Busque la sección "Generation 1" (o la generación más antigua si no existe Gen 1)
   - Extraiga la URL de la imagen de la columna "Back" (para atacante)
   - Extraiga la URL de la imagen de la columna "Normal" (para defensor)
   - Descargue ambas imágenes a `assets/sprites/battle-back/{nombre}.png` y `assets/sprites/battle-front/{nombre}.png`
3. Integrar en `DownloadAllSprites()` para que soporte las nuevas categorías `battle-back` y `battle-front`, iterando todos los Pokémon del Pokédex.

### Fase 3 — APP: wiring
4. Asegurar que `DownloadSprites` en `app/bindings.go` pasa las nuevas categorías correctamente al shell.

### Fase 4 — Frontend: UI de batalla clásica
5. En `builds.ts`, modificar `renderBattleSection()`:
   - Mostrar sprite del defensor arriba-derecha (frontal, `battle-front`)
   - Mostrar sprite del atacante abajo-izquierda (espaldas, `battle-back`)
   - Barra de HP del defensor arriba con nombre y nivel
   - Barra de HP del atacante abajo con nombre, nivel y HP numérico
   - Fallback a sprite `home-normal` si no existe el sprite de batalla
6. En `_builds.scss`, agregar estilos:
   - Layout de batalla tipo Game Boy (fondo claro/rosado, sprites posicionados)
   - Sprite del atacante: abajo-izq, mayor tamaño, `image-rendering: pixelated`
   - Sprite del defensor: arriba-der, menor tamaño, `image-rendering: pixelated`
   - Cajas de info con nombre, nivel, HP bar estilo retro

## Tests
| Archivo | Que se testea |
|---------|---------------|
| `shell/pokemondb_sprites_test.go` | Parseo correcto de tabla de sprites, extracción de URLs Gen 1 |
| Manual | Verificar que sprites se descargan a `assets/sprites/battle-back/` y `battle-front/` |
| Manual | UI de batalla muestra sprites correctamente posicionados |

## Criterios de aceptacion
- [x] Se scrapean y descargan sprites Gen 1 "Back" a `assets/sprites/battle-back/`
- [x] Se scrapean y descargan sprites Gen 1 "Normal" a `assets/sprites/battle-front/`
- [x] Si no hay Gen 1, se usa la generación más antigua disponible
- [x] Se respeta Crawl-delay de 2 segundos entre requests
- [x] Archivos ya descargados se saltan
- [x] La UI de batalla muestra el atacante abajo-izq de espaldas
- [x] La UI de batalla muestra el defensor arriba-der de frente
- [x] Barras de HP y nombres estilo Game Boy clásico
- [x] Fallback a sprite home-normal si no existe sprite de batalla
- [x] `image-rendering: pixelated` para los sprites Gen 1

## Notas
- Los XPaths proporcionados son aproximados y pueden variar según el Pokémon. El scraper debe ser robusto y buscar por contenido/estructura, no por posición absoluta del DOM.
- Con ~150 Pokémon de Gen 1 y 2s de delay, la descarga tarda ~10 minutos. Para Pokémon de otras generaciones sin sprites Gen 1, se busca la generación más antigua.
- Los sprites Gen 1 son de baja resolución (~56x56 px), lo que da el efecto retro deseado. Usar `image-rendering: pixelated` al escalar.
