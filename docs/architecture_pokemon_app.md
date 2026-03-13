# Pokemon App

**Fecha**: 2026-03-13
**Estado**: diseño

---

## Proposito

Aplicacion web que consume la PokéAPI para mostrar informacion de Pokémon a traves de un frontend interactivo servido por un backend en Go.

## Modulo Go

```
module github.com/alfon/pokemon-app

go 1.22
```

---

## Dominio

### Entidades

| Entidad | Descripcion | Campos clave |
|---------|-------------|--------------|
| `Pokemon` | Datos basicos de un Pokémon | `ID`, `Name`, `Types`, `Stats`, `Sprites` |
| `PokemonType` | Tipo elemental de un Pokémon | `Name`, `URL` |
| `Stat` | Estadistica base de un Pokémon | `Name`, `BaseStat` |
| `Sprites` | URLs de imagenes/sprites | `FrontDefault`, `FrontShiny` |
| `PokemonListItem` | Elemento de la lista paginada | `Name`, `URL` |
| `PokemonListResponse` | Respuesta paginada de la API | `Count`, `Next`, `Previous`, `Results` |

### Tipos de dominio (Core)

```go
// core/domain.go

type Pokemon struct {
    ID      int
    Name    string
    Types   []PokemonType
    Stats   []Stat
    Sprites Sprites
    Height  int
    Weight  int
}

type PokemonType struct {
    Name string
}

type Stat struct {
    Name     string
    BaseStat int
}

type Sprites struct {
    FrontDefault string
    FrontShiny   string
}

type PokemonListItem struct {
    Name string
    URL  string
}

type PokemonListResponse struct {
    Count    int
    Next     string
    Previous string
    Results  []PokemonListItem
}
```

### Reglas de negocio

- Los nombres de Pokémon se normalizan a minusculas para busquedas.
- Las estadisticas base se usan tal cual vienen de la API (rango 1-255).
- La lista de Pokémon soporta paginacion con offset y limit.

---

## Capas

### Core — logica pura

Funciones sin efectos secundarios. Solo transforman datos.

| Archivo | Responsabilidad |
|---------|-----------------|
| `core/domain.go` | Tipos de dominio: Pokemon, PokemonType, Stat, Sprites, PokemonListItem, PokemonListResponse |
| `core/ports.go` | Interfaces que Shell implementa (PokemonFetcher) |
| `core/pokemon.go` | Funciones puras: filtrado por tipo, busqueda por nombre, ordenamiento por stat |
| `core/pokemon_test.go` | Tests unitarios de funciones puras |

**Interfaces que Core define** (las implementa Shell):

```go
// core/ports.go

type PokemonFetcher interface {
    FetchPokemon(name string) (Pokemon, error)
    FetchPokemonList(offset int, limit int) (PokemonListResponse, error)
}
```

### Shell — adaptadores externos

Todo lo que tiene efectos secundarios.

| Archivo | Tipo | Dependencia externa |
|---------|------|---------------------|
| `shell/pokeapi.go` | http | PokéAPI (pokeapi.co/api/v2) |

**Implementa**:

| Interfaz (de Core) | Implementacion (en Shell) |
|---------------------|---------------------------|
| `PokemonFetcher` | `shell/pokeapi.go` — PokeAPIClient |

### APP — cableado y arranque

| Archivo | Responsabilidad |
|---------|-----------------|
| `app/main.go` | Entry point, inyeccion de dependencias |
| `app/config.go` | Lectura de configuracion (env vars: puerto, base URL de API) |
| `app/server.go` | Setup del servidor HTTP, sirve frontend estatico + API |
| `app/handlers/pokemon.go` | Handler que orquesta Core + Shell para endpoints de Pokemon |

---

## Estructura de archivos

```
pokemon-app/
  core/
    domain.go
    ports.go
    pokemon.go
    pokemon_test.go
  shell/
    pokeapi.go
  app/
    main.go
    config.go
    server.go
    handlers/
      pokemon.go
  frontend/
    index.html
    style.css
    app.js
  go.mod
  go.sum
```

---

## Dependencias externas

| Paquete | Version | Capa | Para que |
|---------|---------|------|----------|
| `net/http` (stdlib) | - | Shell, APP | Cliente HTTP a PokéAPI + servidor HTTP |
| `encoding/json` (stdlib) | - | Shell | Deserializar respuestas JSON de PokéAPI |
| `html/template` (stdlib) | - | APP | Servir frontend (opcional, si se usan templates) |

---

## Configuracion

| Variable | Tipo | Default | Descripcion |
|----------|------|---------|-------------|
| `PORT` | string | `8080` | Puerto del servidor HTTP |
| `POKEAPI_BASE_URL` | string | `https://pokeapi.co/api/v2` | URL base de PokéAPI |

---

## Endpoints / Comandos

| Metodo | Ruta / Comando | Handler | Descripcion |
|--------|----------------|---------|-------------|
| `GET` | `/` | static file server | Sirve el frontend (index.html) |
| `GET` | `/api/pokemon` | `handlers.ListPokemon` | Lista paginada de Pokémon (?offset=0&limit=20) |
| `GET` | `/api/pokemon/{name}` | `handlers.GetPokemon` | Detalle de un Pokémon por nombre o ID |

---

## Flujo de datos

```
GET /api/pokemon/pikachu -> Handler (APP)
                              |
                              v
                          Shell.FetchPokemon("pikachu")  // HTTP GET a PokéAPI
                              |
                              v
                          Core.Pokemon (struct poblado)  // datos del dominio
                              |
                              v
                          JSON response al frontend
```

```
GET /api/pokemon?offset=0&limit=20 -> Handler (APP)
                                        |
                                        v
                                    Shell.FetchPokemonList(0, 20)  // HTTP GET a PokéAPI
                                        |
                                        v
                                    Core.PokemonListResponse  // lista paginada
                                        |
                                        v
                                    JSON response al frontend
```

---

## Plan de implementacion

1. [ ] Definir tipos de dominio en `core/domain.go`
2. [ ] Definir interfaces (ports) en `core/ports.go`
3. [ ] Implementar logica pura en `core/pokemon.go` + tests
4. [ ] Implementar cliente HTTP de PokéAPI en `shell/pokeapi.go`
5. [ ] Cablear config en `app/config.go`
6. [ ] Implementar handlers en `app/handlers/pokemon.go`
7. [ ] Setup servidor HTTP en `app/server.go` + `app/main.go`
8. [ ] Crear frontend basico (HTML/CSS/JS) en `frontend/`
9. [ ] Tests de integracion

---

## Decisiones de diseno

| Decision | Alternativa descartada | Razon |
|----------|------------------------|-------|
| Frontend en HTML/CSS/JS vanilla servido como static files | React/SPA con build step | Simplicidad, sin dependencias de Node.js, el Go server sirve todo |
| Solo stdlib de Go (net/http, encoding/json) | Framework como Gin/Echo | Proyecto pequeno, stdlib es suficiente, menos dependencias |
| PokéAPI como unica fuente de datos | Base de datos local con cache | Primera version simple, sin persistencia; se puede anadir cache despues |

---

## Notas

- PokéAPI no requiere autenticacion y no tiene rate limit oficial (politica de uso justo).
- El frontend hace llamadas al backend Go (no directamente a PokéAPI), asi el backend actua como proxy y se puede anadir cache en el futuro.
- Ver `docs/reporte_conexiones_pokemon_api.md` para analisis detallado de las capacidades y limitaciones de PokéAPI.
