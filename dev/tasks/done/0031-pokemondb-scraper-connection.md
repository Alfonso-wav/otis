# Conexion con pokemondb.net via scraping

**ID**: 0031-pokemondb-scraper-connection
**Estado**: done
**Fecha**: 2026-03-14

---

## Descripcion

Establecer conexion con https://pokemondb.net/ para extraer datos de Pokemon.

**Investigacion previa**:
- pokemondb.net **NO ofrece API publica** (ni gratuita ni de pago).
- El `robots.txt` permite scraping general con `Crawl-delay: 2` segundos.
- Rutas bloqueadas: solo `/pokebase/search?` y `/pokebase/revisions` (foro).
- Bots bloqueados: `wget`, `WebReaper`, `AhrefsBot` (no afecta a un cliente Go).

**Conclusion**: Se implementara scraping HTTP + parsing HTML para extraer datos.

## Capas afectadas
- **Core**: Definir tipos de dominio para los datos extraidos de pokemondb (si difieren de los existentes de PokeAPI) y la interfaz del scraper.
- **Shell**: Implementar el cliente HTTP scraper con rate limiting (2s entre requests), parsing HTML con goquery o similar.
- **APP**: Wiring del nuevo scraper, configuracion de URLs base.

## Dependencias externas nuevas
| Paquete | Uso |
|---------|-----|
| `github.com/PuerkitoBio/goquery` | Parsing HTML / seleccion CSS para scraping |

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/ports.go` | modificar | Agregar interfaz `PokemonDBScraper` con metodos para los datos a extraer |
| `core/domain.go` | modificar | Agregar tipos de dominio si pokemondb ofrece datos que no estan en PokeAPI |
| `shell/pokemondb.go` | crear | Cliente HTTP base: struct `PokemonDBClient`, rate limiter (2s), fetch generico con parsing HTML |
| `shell/pokemondb_pokedex.go` | crear | Scraper de datos de Pokedex: stats, tipos, habilidades, descripciones |
| `app/config.go` | modificar | Agregar config `POKEMONDB_BASE_URL` (default: `https://pokemondb.net`) |
| `app/bindings.go` | modificar | Inyectar el nuevo scraper y exponer metodos al frontend |
| `go.mod` | modificar | Agregar dependencia goquery |

## Plan de implementacion

### Fase 1 — Cliente base y rate limiting
1. Agregar dependencia `goquery` al proyecto (`go get github.com/PuerkitoBio/goquery`)
2. Crear `shell/pokemondb.go` con:
   - Struct `PokemonDBClient` con `http.Client`, `baseURL`, y un rate limiter (ticker o time.Sleep de 2s entre requests)
   - Metodo generico `fetchPage(path string) (*goquery.Document, error)` que respeta el crawl delay
   - User-Agent descriptivo (no imitar navegador, ser transparente)
3. Agregar `POKEMONDB_BASE_URL` en `app/config.go`

### Fase 2 — Scraper de Pokedex basico
4. Analizar la estructura HTML de `https://pokemondb.net/pokedex/all` para mapear selectores CSS
5. Crear `shell/pokemondb_pokedex.go` con metodo `FetchPokedex()` que extraiga la tabla principal
6. Definir interfaz `PokemonDBScraper` en `core/ports.go`
7. Parsear: numero nacional, nombre, tipos, stats base (HP, Atk, Def, SpA, SpD, Spe), total

### Fase 3 — Wiring
8. Inyectar `PokemonDBClient` en APP
9. Exponer al menos un endpoint de prueba al frontend
10. Verificar que el scraping funciona end-to-end

## Tests
| Archivo | Que se testea |
|---------|---------------|
| `shell/pokemondb_test.go` | Rate limiter respeta 2s entre requests |
| `shell/pokemondb_pokedex_test.go` | Parsing correcto de HTML de ejemplo (fixture local, no hit real) |
| Manual | Ejecutar scrape de `/pokedex/all` y verificar datos extraidos |

## Criterios de aceptacion
- [x] `PokemonDBClient` respeta Crawl-delay de 2 segundos entre requests
- [x] Se puede hacer fetch y parse de al menos una pagina de pokemondb.net
- [x] El scraper de Pokedex extrae correctamente: id, nombre, tipos y stats de todos los Pokemon
- [x] La interfaz `PokemonDBScraper` esta definida en Core sin dependencias externas
- [x] El wiring funciona end-to-end desde APP
- [x] Tests unitarios pasan con fixtures HTML locales

## Notas
- Esta tarea establece la **infraestructura base**. Tareas futuras (0032+) pueden agregar scrapers para moves, abilities, items, evolution chains, sprites, movesets competitivos, etc.
- Respetar siempre el `Crawl-delay: 2` del robots.txt. No paralelizar requests.
- Considerar cache local de paginas ya scrapeadas para no repetir requests innecesarios.
- pokemondb.net tiene datos curados que complementan PokeAPI: sprites de alta calidad, movesets competitivos por generacion, tablas de stats mas completas.
