# Crear capa REST HTTP para el backend Go

**ID**: 0057-android-rest-api-layer
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Crear un servidor HTTP estándar en Go que exponga todos los métodos de `app/bindings.go` como endpoints REST JSON. Esta capa reemplaza los bindings IPC de Wails para permitir que el frontend se comunique via HTTP (necesario para la versión Android con Capacitor/WebView).

El servidor coexiste con la versión desktop: `app/desktop/` mantiene Wails, `app/mobile/` contiene el nuevo servidor REST.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Nueva subcapa `app/mobile/` con servidor HTTP, handlers REST y exports para gomobile.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `app/mobile/server.go` | crear | Servidor HTTP con net/http, configuración de CORS, registro de rutas |
| `app/mobile/handlers.go` | crear | Handlers REST que mapean cada método de bindings.go a un endpoint |
| `app/mobile/mobile.go` | crear | Funciones exportadas para gomobile bind (Start/Stop server) |
| `app/main.go` | modificar | Mover a `app/desktop/main.go` para separar entry points |
| `app/bindings.go` | modificar | Mover a `app/desktop/bindings.go` |

## Plan de implementacion

1. Reorganizar `app/` en `app/desktop/` (Wails) y `app/mobile/` (REST)
2. Crear servidor HTTP en `app/mobile/server.go` con rutas para cada binding
3. Implementar handlers REST que reusen la misma lógica de orquestación
4. Configurar CORS para permitir llamadas desde localhost/WebView
5. Crear funciones gomobile-compatible (`Start(port int, dataDir string)`, `Stop()`)
6. Verificar que la versión desktop sigue funcionando sin cambios

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `app/mobile/server_test.go` | Endpoints REST devuelven JSON correcto |
| `app/mobile/handlers_test.go` | Handlers manejan errores correctamente |

## Criterios de aceptacion

- [x] Servidor HTTP arranca y responde en localhost:PORT
- [x] Todos los métodos de bindings.go tienen un endpoint REST equivalente
- [x] CORS configurado correctamente
- [x] La versión desktop (Wails) sigue compilando y funcionando
- [x] Core no se modifica

## Notas

- Usar solo `net/http` de stdlib para maximizar compatibilidad con gomobile.
- Los endpoints deben seguir el patrón: `GET /api/pokemon/{name}`, `POST /api/battle/simulate`, etc.
- Documentar el mapeo bindings → endpoints en un comentario al inicio de handlers.go.
- Este es el paso previo a la compilación con gomobile (tarea 0059).
