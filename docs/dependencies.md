# Dependencias y Estructura del Proyecto

> Última actualización: 2026-03-13

## Módulo
- **Nombre**: `github.com/alfon/pokemon-app`
- **Versión de Go**: 1.18

## Estructura del proyecto

```
otis/
├── main.go                  # Entry point — arranca Wails con bindings de app y shell
├── go.mod
├── go.sum
├── core/                    # Lógica pura de dominio (sin efectos secundarios)
│   ├── domain.go            # Tipos de dominio (Pokemon, PokemonType, etc.)
│   ├── ports.go             # Interfaces/puertos (contratos para Shell)
│   ├── pokemon.go           # Funciones puras de transformación y filtrado
│   └── pokemon_test.go      # Tests unitarios de lógica pura
├── shell/                   # Adaptadores al mundo exterior
│   └── pokeapi.go           # Cliente HTTP para PokéAPI
├── app/                     # Cableado, bindings y configuración
│   ├── bindings.go          # Bindings Wails (orquesta Core + Shell)
│   ├── bindings_test.go     # Tests de bindings
│   └── config.go            # Configuración de la app
├── frontend/                # Assets del frontend (Wails)
│   └── wailsjs/
│       ├── go/app/          # Bindings JS generados
│       └── runtime/         # Runtime de Wails JS
├── dev/
│   └── tasks/
│       ├── todo/            # Tareas pendientes
│       └── done/            # Tareas completadas
└── docs/                    # Documentación del proyecto
```

## Dependencias directas

| Paquete | Versión actual | Última versión | Estado |
|---------|---------------|----------------|--------|
| `github.com/wailsapp/wails/v2` | v2.5.1 | v2.11.0 | ⚠️ Actualización disponible |

## Dependencias indirectas

| Paquete | Versión actual | Última versión | Estado |
|---------|---------------|----------------|--------|
| `github.com/bep/debounce` | v1.2.1 | v1.2.1 | ✅ Al día |
| `github.com/go-ole/go-ole` | v1.2.6 | v1.3.0 | ⚠️ Actualización |
| `github.com/google/uuid` | v1.1.2 | v1.6.0 | ⚠️ Actualización |
| `github.com/jchv/go-winloader` | v0.0.0-20210711... | v0.0.0-20250406... | ⚠️ Actualización |
| `github.com/labstack/echo/v4` | v4.9.0 | v4.15.1 | ⚠️ Actualización |
| `github.com/labstack/gommon` | v0.3.1 | v0.4.2 | ⚠️ Actualización |
| `github.com/leaanthony/go-ansi-parser` | v1.0.1 | v1.6.1 | ⚠️ Actualización |
| `github.com/leaanthony/gosod` | v1.0.3 | v1.0.4 | ⚠️ Actualización |
| `github.com/leaanthony/slicer` | v1.5.0 | v1.6.0 | ⚠️ Actualización |
| `github.com/mattn/go-colorable` | v0.1.11 | v0.1.14 | ⚠️ Actualización |
| `github.com/mattn/go-isatty` | v0.0.14 | v0.0.20 | ⚠️ Actualización |
| `github.com/pkg/browser` | v0.0.0-20210706... | v0.0.0-20240102... | ⚠️ Actualización |
| `github.com/pkg/errors` | v0.9.1 | v0.9.1 | ✅ Al día |
| `github.com/samber/lo` | v1.27.1 | v1.53.0 | ⚠️ Actualización |
| `github.com/tkrajina/go-reflector` | v0.5.5 | v0.5.8 | ⚠️ Actualización |
| `github.com/valyala/bytebufferpool` | v1.0.0 | v1.0.0 | ✅ Al día |
| `github.com/valyala/fasttemplate` | v1.2.1 | v1.2.2 | ⚠️ Actualización |
| `github.com/wailsapp/mimetype` | v1.4.1 | v1.4.1 | ✅ Al día |
| `golang.org/x/crypto` | v0.1.0 | v0.49.0 | ⚠️ Actualización |
| `golang.org/x/exp` | v0.0.0-20220303... | v0.0.0-20260312... | ⚠️ Actualización |
| `golang.org/x/net` | v0.7.0 | v0.52.0 | ⚠️ Actualización |
| `golang.org/x/sys` | v0.5.0 | v0.42.0 | ⚠️ Actualización |
| `golang.org/x/text` | v0.7.0 | v0.35.0 | ⚠️ Actualización |

### Dependencias transitivas adicionales (no en go.mod directo)

| Paquete | Versión |
|---------|---------|
| `atomicgo.dev/cursor` | v0.1.1 |
| `atomicgo.dev/keyboard` | v0.2.8 |
| `bitbucket.org/creachadair/shell` | v0.0.7 |
| `github.com/Masterminds/semver` | v1.5.0 |
| `github.com/Microsoft/go-winio` | v0.4.16 |
| `github.com/acarl005/stripansi` | v0.0.0-20180116... |
| `github.com/alecthomas/chroma` | v0.10.0 |
| `github.com/aymerick/douceur` | v0.2.0 |
| `github.com/bitfield/script` | v0.19.0 |
| `github.com/charmbracelet/glamour` | v0.5.0 |
| `github.com/containerd/console` | v1.0.3 |
| `github.com/davecgh/go-spew` | v1.1.1 |
| `github.com/dlclark/regexp2` | v1.4.0 |
| `github.com/emirpasic/gods` | v1.12.0 |
| `github.com/flytam/filenamify` | v1.0.0 |
| `github.com/fsnotify/fsnotify` | v1.4.9 |
| `github.com/go-git/gcfg` | v1.5.0 |
| `github.com/go-git/go-billy/v5` | v5.2.0 |
| `github.com/go-git/go-git/v5` | v5.3.0 |
| `github.com/golang-jwt/jwt` | v3.2.2+incompatible |
| `github.com/google/shlex` | v0.0.0-20191202... |
| `github.com/gookit/color` | v1.5.2 |
| `github.com/gorilla/css` | v1.0.0 |
| `github.com/imdario/mergo` | v0.3.12 |
| `github.com/jackmordaunt/icns` | v1.0.0 |
| `github.com/jbenet/go-context` | v0.0.0-20150711... |
| `github.com/kevinburke/ssh_config` | v0.0.0-20201106... |
| `github.com/kr/pretty` | v0.3.0 |
| `github.com/kr/text` | v0.2.0 |
| `github.com/leaanthony/clir` | v1.3.0 |
| `github.com/leaanthony/debme` | v1.2.1 |
| `github.com/leaanthony/winicon` | v1.0.0 |
| `github.com/lithammer/fuzzysearch` | v1.1.5 |
| `github.com/lucasb-eyer/go-colorful` | v1.2.0 |
| `github.com/matryer/is` | v1.4.0 |
| `github.com/mattn/go-runewidth` | v0.0.13 |
| `github.com/microcosm-cc/bluemonday` | v1.0.17 (retracted) |
| `github.com/mitchellh/go-homedir` | v1.1.0 |
| `github.com/muesli/reflow` | v0.3.0 |
| `github.com/muesli/termenv` | v0.9.0 |
| `github.com/nfnt/resize` | v0.0.0-20180221... |
| `github.com/olekukonko/tablewriter` | v0.0.5 |
| `github.com/pterm/pterm` | v0.12.49 |
| `github.com/rivo/uniseg` | v0.2.0 |
| `github.com/sabhiram/go-gitignore` | v0.0.0-20210923... |
| `github.com/sergi/go-diff` | v1.2.0 |
| `github.com/stretchr/objx` | v0.1.0 |
| `github.com/stretchr/testify` | v1.8.0 |
| `github.com/tc-hib/winres` | v0.1.5 |
| `github.com/thoas/go-funk` | v0.9.1 |
| `github.com/tidwall/gjson` | v1.9.3 |
| `github.com/tidwall/match` | v1.1.1 |
| `github.com/tidwall/pretty` | v1.2.0 |
| `github.com/tidwall/sjson` | v1.1.7 |
| `github.com/wzshiming/ctc` | v1.2.3 |
| `github.com/wzshiming/winseq` | v0.0.0-20200112... |
| `github.com/xanzy/ssh-agent` | v0.3.0 |
| `github.com/xo/terminfo` | v0.0.0-20210125... |
| `github.com/yuin/goldmark` | v1.4.13 |
| `github.com/yuin/goldmark-emoji` | v1.0.1 |
| `golang.org/x/image` | v0.5.0 |
| `golang.org/x/mod` | v0.6.0-dev... |
| `golang.org/x/term` | v0.5.0 |
| `golang.org/x/time` | v0.0.0-20201208... |
| `golang.org/x/tools` | v0.1.12 |
| `golang.org/x/xerrors` | v0.0.0-20200804... |
| `gopkg.in/check.v1` | v1.0.0-20200227... |
| `gopkg.in/warnings.v0` | v0.1.2 |
| `gopkg.in/yaml.v3` | v3.0.1 |

## Mapa de imports por paquete interno

### `github.com/alfon/pokemon-app` (main.go)
- **Stdlib**: `embed`, `log`
- **Internos**: `github.com/alfon/pokemon-app/app`, `github.com/alfon/pokemon-app/shell`
- **Externos**: `github.com/wailsapp/wails/v2`, `github.com/wailsapp/wails/v2/pkg/options`, `github.com/wailsapp/wails/v2/pkg/options/assetserver`

### `github.com/alfon/pokemon-app/app` (bindings.go, config.go)
- **Stdlib**: `context`, `os`
- **Internos**: `github.com/alfon/pokemon-app/core`
- **Externos**: ninguno

### `github.com/alfon/pokemon-app/core` (domain.go, ports.go, pokemon.go)
- **Stdlib**: `strings`
- **Internos**: ninguno
- **Externos**: ninguno

### `github.com/alfon/pokemon-app/shell` (pokeapi.go)
- **Stdlib**: `encoding/json`, `fmt`, `net/http`, `time`
- **Internos**: `github.com/alfon/pokemon-app/core`
- **Externos**: ninguno

## Resumen

- **Total dependencias directas**: 1
- **Total dependencias indirectas (go.mod)**: 22
- **Total dependencias transitivas**: 63 adicionales (resueltas por `go list -m all`)
- **Total dependencias resueltas**: 86
- **Dependencias con actualización disponible**: 62
- **Dependencias retractadas**: 1 (`github.com/microcosm-cc/bluemonday` v1.0.17)

### Notas

- La dependencia directa `wails/v2` tiene una actualización mayor disponible (v2.5.1 → v2.11.0).
- `golang.org/x/crypto` y `golang.org/x/net` están muy desactualizadas y pueden contener vulnerabilidades de seguridad conocidas.
- `bluemonday` v1.0.17 está **retractada** por sus autores — se recomienda actualizar a v1.0.27.
- Core no tiene dependencias externas, cumpliendo con la arquitectura de capas definida.
- Shell solo usa stdlib + Core, sin dependencias externas directas.
