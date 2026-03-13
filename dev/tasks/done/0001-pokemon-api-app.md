# Aplicacion Pokemon con frontend y conexion a PokéAPI

**ID**: 0001-pokemon-api-app
**Estado**: done
**Fecha**: 2026-03-13

---

## Descripcion

Crear una aplicacion web en Go que consuma la PokéAPI (pokeapi.co/api/v2) y exponga los datos a traves de un frontend interactivo. El backend actua como proxy/orquestador entre el frontend y PokéAPI, siguiendo la arquitectura Core/Shell/APP.

## Capas afectadas

- **Core**: Tipos de dominio (Pokemon, PokemonType, Stat, Sprites), interfaz PokemonFetcher, funciones puras de filtrado y transformacion.
- **Shell**: Cliente HTTP que implementa PokemonFetcher contra PokéAPI.
- **APP**: Servidor HTTP, configuracion, handlers REST, servidor de archivos estaticos para el frontend.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | crear | Tipos de dominio: Pokemon, PokemonType, Stat, Sprites, PokemonListItem, PokemonListResponse |
| `core/ports.go` | crear | Interfaz PokemonFetcher |
| `core/pokemon.go` | crear | Funciones puras: normalizacion de nombres, filtrado por tipo |
| `core/pokemon_test.go` | crear | Tests unitarios de funciones puras |
| `shell/pokeapi.go` | crear | Cliente HTTP a PokéAPI, implementa PokemonFetcher |
| `app/config.go` | crear | Lectura de PORT y POKEAPI_BASE_URL desde env vars |
| `app/handlers/pokemon.go` | crear | Handlers HTTP: ListPokemon, GetPokemon |
| `app/server.go` | crear | Setup del servidor HTTP con rutas |
| `app/main.go` | crear | Entry point, inyeccion de dependencias |
| `frontend/index.html` | crear | Pagina principal del frontend |
| `frontend/style.css` | crear | Estilos del frontend |
| `frontend/app.js` | crear | Logica del frontend: llamadas a la API, renderizado |
| `go.mod` | crear | Modulo Go |

## Plan de implementacion

1. Inicializar modulo Go (`go mod init`)
2. Crear tipos de dominio en `core/domain.go`
3. Definir interfaz PokemonFetcher en `core/ports.go`
4. Implementar funciones puras en `core/pokemon.go` con tests
5. Implementar cliente PokéAPI en `shell/pokeapi.go`
6. Crear configuracion en `app/config.go`
7. Implementar handlers en `app/handlers/pokemon.go`
8. Setup servidor HTTP en `app/server.go`
9. Crear entry point en `app/main.go`
10. Crear frontend basico (HTML/CSS/JS) con lista de Pokemon y vista de detalle
11. Probar la aplicacion end-to-end

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/pokemon_test.go` | Funciones puras: normalizacion de nombres, filtrado por tipo, ordenamiento por stat |

## Criterios de aceptacion

- [ ] El servidor arranca y sirve el frontend en `/`
- [ ] `GET /api/pokemon?offset=0&limit=20` retorna lista paginada de Pokemon en JSON
- [ ] `GET /api/pokemon/{name}` retorna detalle de un Pokemon (nombre, tipos, stats, sprites)
- [ ] El frontend muestra la lista de Pokemon con imagenes (sprites)
- [ ] El frontend permite ver el detalle de un Pokemon al hacer click
- [ ] Core no importa ningun paquete externo (solo stdlib)
- [ ] Shell implementa la interfaz PokemonFetcher definida en Core
- [ ] Tests de Core pasan sin dependencias externas

## Notas

- Se usa solo stdlib de Go (net/http, encoding/json) para mantener cero dependencias externas.
- El frontend es HTML/CSS/JS vanilla, servido como archivos estaticos por el servidor Go.
- PokéAPI no requiere autenticacion. Ver `docs/reporte_conexiones_pokemon_api.md` para detalles.
- Documento de arquitectura completo en `docs/architecture_pokemon_app.md`.
