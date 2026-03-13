# Migrar aplicacion a Wails con IPC

**ID**: 0002-wails-ipc-migration
**Estado**: done
**Fecha**: 2026-03-13

---

## Descripcion

Convertir la aplicacion web actual (Go HTTP server + frontend estático) a una aplicacion de escritorio usando **Wails v2**. Reemplazar la comunicación HTTP REST entre frontend y backend por **IPC bindings** nativos de Wails, donde el frontend invoca funciones Go directamente a través del bridge de WebView. Esto elimina la latencia HTTP local y mejora la fluidez de la UI.

## Capas afectadas

- **Core**: Sin cambios. La lógica pura y tipos de dominio se mantienen intactos.
- **Shell**: Sin cambios. El cliente HTTP a PokéAPI sigue funcionando igual.
- **APP**: Cambio mayor. Se reemplaza el servidor HTTP por la inicialización de Wails. Los handlers HTTP se convierten en métodos de un struct `App` que se bindea a Wails. El frontend se embebe con `embed.FS`.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `app/main.go` | modificar | Reemplazar `http.ListenAndServe` por `wails.Run()` con configuración de ventana |
| `app/server.go` | eliminar | Ya no se necesita servidor HTTP, Wails maneja el serving |
| `app/bindings.go` | crear | Struct `App` con métodos bindeados: `ListPokemon`, `GetPokemon` |
| `app/handlers/pokemon.go` | eliminar | Los handlers HTTP se reemplazan por métodos IPC en `bindings.go` |
| `app/config.go` | modificar | Eliminar config de PORT (ya no hay servidor HTTP), mantener POKEAPI_BASE_URL |
| `frontend/index.html` | modificar | Agregar script de Wails runtime |
| `frontend/app.js` | modificar | Reemplazar `fetch()` por llamadas IPC: `window.go.main.App.Method()` |
| `go.mod` | modificar | Agregar dependencia de `github.com/wailsapp/wails/v2` |
| `wails.json` | crear | Configuración de proyecto Wails |

## Plan de implementacion

1. Instalar Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`) y verificar con `wails doctor`
2. Agregar dependencia de Wails v2 al `go.mod`
3. Crear `app/bindings.go` con struct `App` que contenga el `PokemonFetcher` y métodos `ListPokemon(offset, limit int) ([]core.Pokemon, error)` y `GetPokemon(name string) (core.Pokemon, error)`
4. Modificar `app/main.go` para usar `wails.Run()` con la config de ventana, embeber frontend con `embed.FS`, y bindear el struct `App`
5. Modificar `app/config.go` para eliminar la configuración de puerto
6. Modificar `frontend/app.js` para reemplazar todas las llamadas `fetch("/api/...")` por `window.go.main.App.ListPokemon()` y `window.go.main.App.GetPokemon()`
7. Modificar `frontend/index.html` para incluir el script de Wails runtime
8. Crear `wails.json` con la configuración del proyecto
9. Eliminar `app/server.go` y `app/handlers/pokemon.go`
10. Ejecutar `wails build` y verificar que la app funciona correctamente

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/pokemon_test.go` | Sin cambios — los tests de Core deben seguir pasando |
| `app/bindings_test.go` | Tests de los métodos del struct App con mock de PokemonFetcher |

## Criterios de aceptacion

- [ ] La app arranca como aplicación de escritorio con `wails dev`
- [ ] El frontend muestra la lista de Pokémon paginada (igual que antes)
- [ ] El detalle de un Pokémon se muestra correctamente al hacer click
- [ ] La búsqueda por nombre funciona
- [ ] La comunicación frontend-backend es por IPC (no HTTP local)
- [ ] Core no tiene cambios ni nuevas dependencias
- [ ] Shell no tiene cambios
- [ ] Los tests existentes de Core pasan sin modificación
- [ ] `wails build` genera un ejecutable funcional
- [ ] La dirección de dependencias se mantiene: APP → Shell → Core

## Notas

- **Wails v2** usa WebView2 en Windows (ya incluido en Windows 11), WebKitGTK en Linux, y WebKit en macOS.
- Los métodos bindeados en Wails deben ser **exportados** (mayúscula) y pueden retornar `(value, error)` — Wails convierte automáticamente el error en un reject de la Promise en JS.
- El frontend se embebe en el binario con `//go:embed` — no se necesita servir archivos estáticos.
- Wails genera automáticamente bindings TypeScript/JS para los métodos Go, lo cual facilita el autocomplete en el frontend.
- La arquitectura Core/Shell/APP se mantiene limpia: Wails es solo un cambio en la capa APP (cómo se expone la funcionalidad), no en la lógica ni en los adaptadores.
