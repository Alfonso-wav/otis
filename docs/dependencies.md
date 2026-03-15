# Dependencias y Estructura del Proyecto

> Última actualización: 2026-03-15

## Módulo
- **Nombre**: `github.com/alfon/pokemon-app`
- **Versión de Go**: 1.24.0

## Estructura del proyecto

```
otis/
├── main.go                          # Entry point — arranca Wails con bindings de app y shell
├── go.mod / go.sum                  # Dependencias Go
├── core/                            # Lógica pura de dominio (sin efectos secundarios)
│   ├── domain.go                    # Tipos de dominio (Pokemon, PokemonType, etc.)
│   ├── ports.go                     # Interfaces/puertos (contratos para Shell)
│   ├── battle.go                    # Lógica de batallas Pokémon
│   ├── battle_test.go               # Tests de batallas
│   ├── damage.go                    # Cálculo de daño
│   ├── damage_test.go               # Tests de daño
│   ├── ev_calc.go                   # Calculadora de EVs
│   ├── pokemon.go                   # Funciones puras de transformación y filtrado
│   ├── pokemon_test.go              # Tests de Pokémon
│   ├── teams.go                     # Lógica de equipos
│   └── teams_test.go                # Tests de equipos
├── shell/                           # Adaptadores al mundo exterior
│   ├── pokeapi.go                   # Cliente HTTP para PokeAPI
│   ├── pokeapi_breeding.go          # Datos de crianza desde PokeAPI
│   ├── pokeapi_locations.go         # Datos de ubicaciones desde PokeAPI
│   ├── pokeapi_meta.go              # Metadatos desde PokeAPI
│   ├── pokeapi_moves.go             # Movimientos desde PokeAPI
│   ├── pokeapi_species.go           # Especies desde PokeAPI
│   ├── pokemondb.go                 # Scraper de PokemonDB
│   ├── pokemondb_pokedex.go         # Pokédex desde PokemonDB
│   ├── pokemondb_pokedex_test.go    # Tests de Pokédex
│   ├── pokemondb_sprites.go         # Sprites desde PokemonDB
│   ├── pokemondb_sprites_test.go    # Tests de sprites
│   ├── pokemondb_test.go            # Tests generales PokemonDB
│   ├── teams.go                     # Persistencia de equipos (JSON)
│   ├── teams_test.go                # Tests de persistencia de equipos
│   └── utils.go                     # Utilidades de shell
├── app/                             # Cableado, bindings y configuración
│   ├── bindings.go                  # Bindings Wails (orquesta Core + Shell)
│   ├── bindings_test.go             # Tests de bindings
│   └── config.go                    # Configuración de la app
├── frontend/                        # Frontend (Vite + TypeScript + SCSS)
│   ├── src/
│   │   ├── animations/              # Animaciones (GSAP)
│   │   ├── charts/                  # Gráficas (ECharts)
│   │   ├── components/              # Componentes UI
│   │   ├── pages/                   # Páginas
│   │   │   └── explore/             # Página de exploración
│   │   └── styles/                  # Estilos (SCSS/Bootstrap)
│   ├── wailsjs/                     # Bindings auto-generados por Wails
│   │   ├── go/app/                  # Bindings JS generados
│   │   └── runtime/                 # Runtime de Wails JS
│   ├── dist/                        # Build de producción
│   └── node_modules/                # Deps JS (bootstrap, echarts, gsap, vite, sass, typescript)
├── data/                            # Datos persistidos
│   └── teams/                       # Equipos guardados (JSON)
├── assets/                          # Assets de la app
├── build/                           # Configuración de build
│   ├── bin/                         # Binarios compilados
│   └── windows/                     # Recursos Windows
├── docs/                            # Documentación
│   ├── patch/                       # Notas de parches
│   └── templates/                   # Plantillas
└── dev/                             # Desarrollo
    └── tasks/
        ├── todo/                    # Tareas pendientes
        └── done/                    # Tareas completadas
```

## Dependencias directas

| Paquete | Versión actual | Última versión | Estado |
|---------|---------------|----------------|--------|
| `github.com/PuerkitoBio/goquery` | v1.11.0 | v1.12.0 | ⬆️ Actualización disponible |
| `github.com/wailsapp/wails/v2` | v2.11.0 | v2.11.0 | ✅ Al día |

## Dependencias indirectas (go.mod)

| Paquete | Versión actual | Última versión | Estado |
|---------|---------------|----------------|--------|
| `github.com/andybalholm/cascadia` | v1.3.3 | v1.3.3 | ✅ Al día |
| `github.com/bep/debounce` | v1.2.1 | v1.2.1 | ✅ Al día |
| `github.com/go-ole/go-ole` | v1.3.0 | v1.3.0 | ✅ Al día |
| `github.com/godbus/dbus/v5` | v5.1.0 | v5.2.2 | ⬆️ Actualización |
| `github.com/google/uuid` | v1.6.0 | v1.6.0 | ✅ Al día |
| `github.com/gorilla/websocket` | v1.5.3 | v1.5.3 | ✅ Al día |
| `github.com/jchv/go-winloader` | v0.0.0-20210711... | v0.0.0-20250406... | ⬆️ Actualización |
| `github.com/labstack/echo/v4` | v4.13.3 | v4.15.1 | ⬆️ Actualización |
| `github.com/labstack/gommon` | v0.4.2 | v0.4.2 | ✅ Al día |
| `github.com/leaanthony/go-ansi-parser` | v1.6.1 | v1.6.1 | ✅ Al día |
| `github.com/leaanthony/gosod` | v1.0.4 | v1.0.4 | ✅ Al día |
| `github.com/leaanthony/slicer` | v1.6.0 | v1.6.0 | ✅ Al día |
| `github.com/leaanthony/u` | v1.1.1 | v1.1.1 | ✅ Al día |
| `github.com/mattn/go-colorable` | v0.1.13 | v0.1.14 | ⬆️ Actualización |
| `github.com/mattn/go-isatty` | v0.0.20 | v0.0.20 | ✅ Al día |
| `github.com/pkg/browser` | v0.0.0-20240102... | v0.0.0-20240102... | ✅ Al día |
| `github.com/pkg/errors` | v0.9.1 | v0.9.1 | ✅ Al día |
| `github.com/rivo/uniseg` | v0.4.7 | v0.4.7 | ✅ Al día |
| `github.com/samber/lo` | v1.49.1 | v1.53.0 | ⬆️ Actualización |
| `github.com/tkrajina/go-reflector` | v0.5.8 | v0.5.8 | ✅ Al día |
| `github.com/valyala/bytebufferpool` | v1.0.0 | v1.0.0 | ✅ Al día |
| `github.com/valyala/fasttemplate` | v1.2.2 | v1.2.2 | ✅ Al día |
| `github.com/wailsapp/go-webview2` | v1.0.22 | v1.0.23 | ⬆️ Actualización |
| `github.com/wailsapp/mimetype` | v1.4.1 | v1.4.1 | ✅ Al día |
| `golang.org/x/crypto` | v0.44.0 | v0.49.0 | ⬆️ Actualización |
| `golang.org/x/net` | v0.47.0 | v0.52.0 | ⬆️ Actualización |
| `golang.org/x/sys` | v0.38.0 | v0.42.0 | ⬆️ Actualización |
| `golang.org/x/text` | v0.31.0 | v0.35.0 | ⬆️ Actualización |

## Dependencias transitivas adicionales (resueltas por `go list -m all`)

| Paquete | Versión | Última versión |
|---------|---------|----------------|
| `atomicgo.dev/cursor` | v0.2.0 | — |
| `atomicgo.dev/keyboard` | v0.2.9 | — |
| `atomicgo.dev/schedule` | v0.1.0 | — |
| `dario.cat/mergo` | v1.0.0 | v1.0.2 |
| `github.com/Masterminds/semver` | v1.5.0 | — |
| `github.com/Microsoft/go-winio` | v0.6.1 | v0.6.2 |
| `github.com/ProtonMail/go-crypto` | v1.1.5 | v1.4.0 |
| `github.com/StackExchange/wmi` | v1.2.1 | — |
| `github.com/acarl005/stripansi` | v0.0.0-20180116... | — |
| `github.com/alecthomas/chroma/v2` | v2.14.0 | v2.23.1 |
| `github.com/aymanbagabas/go-osc52/v2` | v2.0.1 | — |
| `github.com/aymerick/douceur` | v0.2.0 | — |
| `github.com/bitfield/script` | v0.24.0 | v0.24.1 |
| `github.com/charmbracelet/glamour` | v0.8.0 | v1.0.0 |
| `github.com/charmbracelet/lipgloss` | v0.12.1 | v1.1.0 |
| `github.com/charmbracelet/x/ansi` | v0.1.4 | v0.11.6 |
| `github.com/cloudflare/circl` | v1.3.7 | v1.6.3 |
| `github.com/containerd/console` | v1.0.3 | v1.0.5 |
| `github.com/cyphar/filepath-securejoin` | v0.3.6 | v0.6.1 |
| `github.com/davecgh/go-spew` | v1.1.1 | — |
| `github.com/dlclark/regexp2` | v1.11.0 | v1.11.5 |
| `github.com/emirpasic/gods` | v1.18.1 | — |
| `github.com/flytam/filenamify` | v1.2.0 | — |
| `github.com/fsnotify/fsnotify` | v1.9.0 | — |
| `github.com/go-git/gcfg` | v1.5.1-0.20230307... | — |
| `github.com/go-git/go-billy/v5` | v5.6.2 | v5.8.0 |
| `github.com/go-git/go-git/v5` | v5.13.2 | v5.17.0 |
| `github.com/golang/groupcache` | v0.0.0-20210331... | v0.0.0-20241129... |
| `github.com/google/go-cmp` | v0.6.0 | v0.7.0 |
| `github.com/google/shlex` | v0.0.0-20191202... | — |
| `github.com/gookit/color` | v1.5.4 | v1.6.0 |
| `github.com/gorilla/css` | v1.0.1 | — |
| `github.com/itchyny/gojq` | v0.12.13 | v0.12.18 |
| `github.com/itchyny/timefmt-go` | v0.1.5 | v0.1.7 |
| `github.com/jackmordaunt/icns` | v1.0.0 | — |
| `github.com/jaypipes/ghw` | v0.13.0 | v0.23.0 |
| `github.com/jaypipes/pcidb` | v1.0.1 | v1.1.1 |
| `github.com/jbenet/go-context` | v0.0.0-20150711... | — |
| `github.com/kevinburke/ssh_config` | v1.2.0 | v1.6.0 |
| `github.com/kr/text` | v0.2.0 | — |
| `github.com/leaanthony/clir` | v1.3.0 | v1.7.0 |
| `github.com/leaanthony/debme` | v1.2.1 | — |
| `github.com/leaanthony/winicon` | v1.0.0 | — |
| `github.com/lithammer/fuzzysearch` | v1.1.8 | — |
| `github.com/lucasb-eyer/go-colorful` | v1.2.0 | v1.3.0 |
| `github.com/matryer/is` | v1.4.1 | — |
| `github.com/mattn/go-runewidth` | v0.0.16 | v0.0.21 |
| `github.com/microcosm-cc/bluemonday` | v1.0.27 | — |
| `github.com/mitchellh/go-homedir` | v1.1.0 | — |
| `github.com/muesli/reflow` | v0.3.0 | — |
| `github.com/muesli/termenv` | v0.15.3-0.20240618... | v0.16.0 |
| `github.com/nfnt/resize` | v0.0.0-20180221... | — |
| `github.com/pjbgf/sha1cd` | v0.3.2 | v0.5.0 |
| `github.com/pmezard/go-difflib` | v1.0.0 | — |
| `github.com/pterm/pterm` | v0.12.80 | v0.12.83 |
| `github.com/sabhiram/go-gitignore` | v0.0.0-20210923... | — |
| `github.com/sergi/go-diff` | v1.3.2-0.20230802... | v1.4.0 |
| `github.com/skeema/knownhosts` | v1.3.0 | v1.3.2 |
| `github.com/stretchr/testify` | v1.10.0 | v1.11.1 |
| `github.com/tc-hib/winres` | v0.3.1 | — |
| `github.com/tidwall/gjson` | v1.14.2 | v1.18.0 |
| `github.com/tidwall/match` | v1.1.1 | v1.2.0 |
| `github.com/tidwall/pretty` | v1.2.0 | v1.2.1 |
| `github.com/tidwall/sjson` | v1.2.5 | — |
| `github.com/wzshiming/ctc` | v1.2.3 | — |
| `github.com/wzshiming/winseq` | v0.0.0-20200112... | v0.0.0-20200720... |
| `github.com/xanzy/ssh-agent` | v0.3.3 | — |
| `github.com/xo/terminfo` | v0.0.0-20220910... | — |
| `github.com/yuin/goldmark` | v1.7.4 | v1.7.16 |
| `github.com/yuin/goldmark-emoji` | v1.0.3 | v1.0.6 |
| `golang.org/x/image` | v0.12.0 | v0.37.0 |
| `golang.org/x/mod` | v0.29.0 | v0.34.0 |
| `golang.org/x/sync` | v0.18.0 | v0.20.0 |
| `golang.org/x/telemetry` | v0.0.0-20240228... | v0.0.0-20260312... |
| `golang.org/x/term` | v0.37.0 | v0.41.0 |
| `golang.org/x/time` | v0.8.0 | v0.15.0 |
| `golang.org/x/tools` | v0.38.0 | v0.43.0 |
| `golang.org/x/xerrors` | v0.0.0-20190717... | v0.0.0-20240903... |
| `gopkg.in/check.v1` | v1.0.0-20200227... | v1.0.0-20201130... |
| `gopkg.in/warnings.v0` | v0.1.2 | — |
| `gopkg.in/yaml.v3` | v3.0.1 | — |
| `howett.net/plist` | v1.0.0 | v1.0.1 |
| `mvdan.cc/sh/v3` | v3.7.0 | v3.13.0 |

## Mapa de imports por paquete interno

### `github.com/alfon/pokemon-app` (main.go)
- **Stdlib**: `embed`, `log`, `net/http`
- **Internos**: `app`, `shell`
- **Externos**: `github.com/wailsapp/wails/v2`, `github.com/wailsapp/wails/v2/pkg/options`, `github.com/wailsapp/wails/v2/pkg/options/assetserver`

### `github.com/alfon/pokemon-app/app`
- **Stdlib**: `context`, `fmt`, `math/rand`, `os`
- **Internos**: `core`
- **Externos**: ninguno

### `github.com/alfon/pokemon-app/core`
- **Stdlib**: `errors`, `fmt`, `math`, `sort`, `strings`
- **Internos**: ninguno
- **Externos**: ninguno

### `github.com/alfon/pokemon-app/shell`
- **Stdlib**: `encoding/json`, `fmt`, `io`, `log`, `net/http`, `os`, `path/filepath`, `regexp`, `strconv`, `strings`, `sync`, `time`
- **Internos**: `core`
- **Externos**: `github.com/PuerkitoBio/goquery`

## Resumen

- **Total dependencias directas**: 2
- **Total dependencias indirectas (go.mod)**: 27
- **Total dependencias resueltas (go list -m all)**: 96
- **Dependencias con actualización disponible**: ~50
- **Dependencias externas usadas directamente en código**: 2 (`github.com/wailsapp/wails/v2`, `github.com/PuerkitoBio/goquery`)

### Notas

- Core no tiene dependencias externas, cumpliendo correctamente con la arquitectura de capas definida.
- Shell solo usa `goquery` como dependencia externa directa (para web scraping).
- App no usa dependencias externas — solo stdlib + Core.
- `goquery` tiene actualización menor disponible (v1.11.0 → v1.12.0).
- Las dependencias `golang.org/x/*` tienen actualizaciones disponibles pero son gestionadas transitivamente por Wails.
