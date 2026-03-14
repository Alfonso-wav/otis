# Descarga masiva de sprites de pokemondb.net a assets locales

**ID**: 0033-download-pokemondb-sprites
**Estado**: done
**Fecha**: 2026-03-14

---

## Descripcion

Descargar todas las imagenes de sprites disponibles en pokemondb.net y almacenarlas localmente en `assets/sprites/` para servir desde el frontend sin depender de CDNs externos. Actualmente el frontend referencia sprites de `img.pokemondb.net` y `raw.githubusercontent.com/PokeAPI/sprites` directamente, lo que causa dependencia de red y latencia.

**Estrategia de descarga (la mas eficiente)**:
- **NO scraping HTML**: los sprites estan en URLs predecibles del CDN `img.pokemondb.net`.
- Se usa la lista de Pokemon del pokedex ya scrapeado (`FetchPokedex()`) para obtener los nombres.
- Se construyen las URLs directamente: `https://img.pokemondb.net/sprites/{categoria}/{variante}/{nombre}.png`
- Se respeta el `Crawl-delay: 2` segundos del robots.txt existente.

**Categorias de sprites disponibles en pokemondb.net**:
| Categoria | URL pattern | Resolucion aprox. | Prioridad |
|-----------|-------------|-------------------|-----------|
| Home (normal) | `/sprites/home/normal/{name}.png` | 250x250 | Alta (ya usado en builds/compare) |
| Home (shiny) | `/sprites/home/shiny/{name}.png` | 250x250 | Alta |
| Icons (move types) | `/images/icons/move-physical.png`, `move-special.png`, `move-status.png` | Pequeno | Alta (ya usado) |

## Capas afectadas
- **Core**: Definir tipos para el resultado de descarga (progreso, errores) y la interfaz del downloader.
- **Shell**: Implementar el descargador de imagenes reutilizando `PokemonDBClient` existente con rate limiting.
- **APP**: Exponer binding al frontend para iniciar/monitorear descarga. Configurar Wails para servir assets estaticos.
- **Frontend**: Actualizar URLs de sprites para apuntar a assets locales con fallback a CDN.

## Dependencias externas nuevas
Ninguna. Se reutiliza `net/http` y el `PokemonDBClient` existente.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/ports.go` | modificar | Agregar interfaz `SpriteDownloader` |
| `core/domain.go` | modificar | Agregar tipos `SpriteDownloadResult`, `SpriteCategory` |
| `shell/pokemondb_sprites.go` | crear | Descargador de sprites: itera nombres, construye URLs, descarga con rate limiting, guarda en disco |
| `shell/pokemondb_sprites_test.go` | crear | Tests con HTTP mock server |
| `app/bindings.go` | modificar | Exponer `DownloadSprites()` al frontend |
| `frontend/src/pages/builds.ts` | modificar | Actualizar `spriteURL()` para usar assets locales |
| `frontend/src/pages/compare.ts` | modificar | Usar assets locales con fallback |
| `frontend/src/pages/pokedex.ts` | modificar | Usar assets locales con fallback |
| `frontend/src/pages/types.ts` | modificar | Usar assets locales con fallback |

## Plan de implementacion

### Fase 1 — Core: tipos e interfaz
1. Definir en `core/domain.go` los tipos:
   - `SpriteCategory` (string type: "home-normal", "home-shiny", "icons")
   - `SpriteDownloadResult` con campos: `Total`, `Downloaded`, `Skipped` (ya existentes), `Failed`, `Errors`
2. Definir en `core/ports.go` la interfaz `SpriteDownloader`:
   - `DownloadAllSprites(destDir string, categories []SpriteCategory) (SpriteDownloadResult, error)`

### Fase 2 — Shell: descargador de sprites
3. Crear `shell/pokemondb_sprites.go`:
   - Metodo `DownloadAllSprites()` en `PokemonDBClient`
   - Flujo: llamar `FetchPokedex()` → obtener nombres → para cada nombre y categoria, construir URL → descargar con `fetchPage` rate-limited → guardar como archivo PNG en `assets/sprites/{categoria}/{nombre}.png`
   - Skip si el archivo ya existe (cache local, evita re-descargas)
   - Log de progreso cada 50 Pokemon
4. Crear `shell/pokemondb_sprites_test.go` con httptest mock server

### Fase 3 — APP: wiring y binding
5. Agregar metodo `DownloadSprites()` en `app/bindings.go`
6. Asegurar que el directorio `assets/sprites/` es accesible por Wails como asset estatico

### Fase 4 — Frontend: usar assets locales
7. Crear helper `localSpriteURL(name, category)` que apunte a `assets/sprites/{cat}/{name}.png`
8. Actualizar `builds.ts`, `compare.ts`, `pokedex.ts`, `types.ts` para usar local con fallback a CDN en `onerror`

## Tests
| Archivo | Que se testea |
|---------|---------------|
| `shell/pokemondb_sprites_test.go` | Descarga correcta a disco, skip de existentes, manejo de 404s |
| Manual | Ejecutar descarga completa y verificar archivos en `assets/sprites/` |
| Manual | Frontend muestra sprites locales correctamente |

## Criterios de aceptacion
- [ ] Los sprites home/normal se descargan a `assets/sprites/home-normal/`
- [ ] Los sprites home/shiny se descargan a `assets/sprites/home-shiny/`
- [ ] Los iconos de tipos de movimiento se descargan a `assets/sprites/icons/`
- [ ] Se respeta Crawl-delay de 2 segundos entre requests
- [ ] Archivos ya descargados se saltan (no re-descarga)
- [ ] El frontend usa sprites locales con fallback a CDN
- [ ] `SpriteDownloadResult` reporta total/descargados/saltados/fallidos

## Notas
- Con ~1000 Pokemon y 2s de delay, la descarga completa de UNA categoria tarda ~33 minutos. Con 2 categorias (normal + shiny) ~66 minutos. Considerar descarga en background con reporte de progreso.
- Los archivos descargados deben agregarse a `.gitignore` (son ~250KB cada uno, ~500MB total para 2 categorias).
- El directorio `assets/sprites/` debe crearse automaticamente si no existe.
- Considerar agregar un comando CLI o boton en frontend para iniciar la descarga bajo demanda, no automaticamente al arrancar.
