# 0041 — Eliminar pestaña Comparar

## Descripción
Eliminar completamente la pestaña **"Comparar"** (comparador de Pokémon) de la aplicación: frontend, backend y bindings.

## Estado
- [x] Done

## Capas afectadas
- **Core** (`core/compare.go`): eliminar archivo completo (tipos `StatComparison`, `PokemonComparison` y función `ComparePokemons`).
- **App** (`app/bindings.go`): eliminar método `ComparePokemons` del struct `App`.
- **Frontend** (`frontend/`): eliminar página, estilos, tab HTML y registros en router/main.

## Plan de implementación

### Paso 1 — Frontend: eliminar tab del HTML
En `frontend/index.html`:
- Eliminar el `<button class="tab-btn" data-tab="compare">Comparar</button>`.
- Eliminar el `<div id="tab-compare" class="tab-page hidden"></div>`.

### Paso 2 — Frontend: eliminar página y estilos
- Eliminar `frontend/src/pages/compare.ts`.
- Eliminar `frontend/src/styles/_compare.scss`.
- Eliminar import de `_compare.scss` en el archivo de estilos principal (si existe).

### Paso 3 — Frontend: limpiar main.ts y router
En `frontend/src/main.ts`:
- Eliminar import de `initCompare`.
- Eliminar registro de la página compare en el router.
- Eliminar llamada a `initCompare()`.

### Paso 4 — Backend: eliminar binding
En `app/bindings.go`:
- Eliminar el método `ComparePokemons(nameA, nameB string)`.

### Paso 5 — Backend: eliminar lógica core
- Eliminar `core/compare.go` completo.

### Paso 6 — Regenerar bindings Wails
- Ejecutar build o `wails generate module` para que los bindings TypeScript reflejen la eliminación.

### Paso 7 — Verificar compilación
- Comprobar que el proyecto compila sin errores (`wails build` o `go build`).
- Verificar que no quedan imports huérfanos.

## Criterios de éxito
- La pestaña "Comparar" no aparece en la UI.
- No queda código muerto relacionado con el comparador.
- El proyecto compila y funciona correctamente.
