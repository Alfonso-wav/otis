# Abstraer llamadas frontend: IPC Wails → fetch HTTP

**ID**: 0058-android-frontend-api-abstraction
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Crear un módulo `api.ts` que abstraiga todas las llamadas al backend. Actualmente el frontend usa `window.go.app.App.Method()` (bindings IPC de Wails). El nuevo módulo detectará el entorno (desktop Wails vs. mobile HTTP) y usará el transporte correcto automáticamente.

Todas las páginas del frontend (`pokedex.ts`, `builds.ts`, `explore.ts`, etc.) deben migrar sus llamadas a usar este módulo.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Sin cambios (solo frontend).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/api.ts` | crear | Módulo central con funciones tipadas que abstraen IPC/HTTP |
| `frontend/src/pages/pokedex.ts` | modificar | Reemplazar `window.go.app.App.*` por imports de `api.ts` |
| `frontend/src/pages/builds.ts` | modificar | Reemplazar llamadas IPC por `api.ts` |
| `frontend/src/pages/explore.ts` | modificar | Reemplazar llamadas IPC por `api.ts` |
| `frontend/src/pages/explore/*.ts` | modificar | Reemplazar llamadas IPC en subpáginas |
| `frontend/src/autocomplete.ts` | modificar | Usar `api.ts` para búsqueda |

## Plan de implementacion

1. Crear `api.ts` con función de detección de entorno (`isWails()`)
2. Implementar cada función del API con doble transporte:
   - Wails: `window.go.app.App.Method(args)`
   - HTTP: `fetch("http://localhost:PORT/api/endpoint", { body: JSON.stringify(args) })`
3. Migrar `pokedex.ts` como primera página de prueba
4. Migrar el resto de páginas una por una
5. Verificar que la versión desktop sigue funcionando
6. Eliminar imports directos de `wailsjs/` en las páginas (solo `api.ts` los usa)

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual | Verificar que la app desktop funciona igual tras la migración |

## Criterios de aceptacion

- [x] Todas las llamadas al backend pasan por `api.ts`
- [x] En entorno Wails, usa IPC como antes
- [x] En entorno HTTP (Android), usa fetch a localhost
- [x] No hay imports directos de `wailsjs/` fuera de `api.ts`
- [x] La app desktop funciona sin regresiones

## Notas

- El módulo debe ser totalmente tipado (TypeScript).
- Patrón de detección: `typeof window.go !== 'undefined'` indica Wails.
- La URL del servidor HTTP se configura via constante o variable de entorno Vite.
