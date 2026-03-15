# Reporte: Migración de Otis a Android

**Fecha**: 2026-03-15
**Proyecto**: Otis — Pokédex y Simulador de Batallas
**Estado actual**: Aplicación de escritorio (Wails v2 — Go + TypeScript/Vite)

---

## 1. Análisis del estado actual

### Stack tecnológico

| Capa | Tecnología | Rol |
|------|-----------|-----|
| Core | Go puro (sin deps externas) | Lógica de negocio: cálculo de daño, batallas, equipos, EVs |
| Shell | Go + goquery | Adaptadores: cliente PokeAPI, scraper PokemonDB, almacenamiento en archivos JSON |
| App | Go + Wails v2 | Entry point, bindings IPC, configuración |
| Frontend | TypeScript + Vite + Bootstrap 5 + ECharts + GSAP | UI completa: Pokédex, explorador, team builder, simulador de batallas |

### Puntos clave de la arquitectura

- **Separación limpia Core/Shell/App**: Core no tiene dependencias externas, lo que facilita reutilizarlo en cualquier plataforma.
- **Frontend web**: La UI ya está hecha en HTML/CSS/JS estándar, lo que la hace portable a WebView.
- **IPC via Wails**: El frontend llama al backend Go a través de `window.go.app.App.Method()` (bindings generados automáticamente por Wails).
- **Persistencia en archivos**: Los equipos se guardan como JSON en `data/teams/`. No hay base de datos.
- **APIs externas**: PokeAPI (REST) y PokemonDB (scraping HTML).
- **Sin autenticación**: No hay cuentas de usuario ni tokens.

### Métricas del proyecto

- ~55 features completadas
- Frontend: ~70KB solo `builds.ts`, ~15 archivos TypeScript
- Backend: ~20 archivos Go entre core, shell y app
- ~40+ métodos expuestos al frontend via Wails bindings

---

## 2. Opciones de migración evaluadas

### Opción A: Capacitor + Go como servidor local (RECOMENDADA)

**Concepto**: Envolver el frontend existente en Capacitor (WebView nativo de Android) y compilar el backend Go como un servidor HTTP local embebido en la app Android usando gomobile.

**Arquitectura**:
```
┌─────────────────────────────────┐
│        Android App (APK)         │
│  ┌───────────────────────────┐  │
│  │   Capacitor WebView       │  │
│  │   (frontend existente)    │  │
│  │   TypeScript + Bootstrap  │  │
│  └──────────┬────────────────┘  │
│             │ HTTP localhost     │
│  ┌──────────▼────────────────┐  │
│  │   Go HTTP Server          │  │
│  │   (gomobile .aar)         │  │
│  │   Core + Shell + REST API │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘
```

**Cambios necesarios**:
1. **App layer**: Reemplazar Wails por un servidor HTTP estándar (`net/http`) que exponga los mismos métodos como endpoints REST.
2. **Frontend**: Reemplazar las llamadas `window.go.app.App.Method()` por `fetch("http://localhost:PORT/api/...")`.
3. **Shell/teams**: Cambiar la ruta de almacenamiento de `data/teams/` al directorio interno de la app Android.
4. **Build**: Usar `gomobile bind` para generar un `.aar` (Android Archive) que inicie el servidor Go al arrancar la app.
5. **Capacitor**: Configurar Capacitor para cargar el frontend desde assets locales y apuntar al servidor Go local.

**Ventajas**:
- Reutiliza >90% del código existente (frontend completo + core + shell)
- El frontend ya funciona en WebView (es HTML/CSS/JS estándar)
- Cambios mínimos en la lógica de negocio (cero cambios en core/)
- Mantiene la misma arquitectura Core/Shell/App
- Puede coexistir con la versión desktop (ambas comparten core/)

**Desventajas**:
- Requiere Android Studio para el proyecto Capacitor
- gomobile tiene limitaciones (no soporta todas las features de Go)
- El servidor HTTP local consume algo más de memoria
- WebView no es tan fluido como UI nativa para animaciones pesadas

**Esfuerzo estimado**: Medio

---

### Opción B: Wails v3 Mobile (EXPERIMENTAL)

**Concepto**: Wails v3 incluye soporte experimental para Android. Migrar de Wails v2 a v3 y usar el target Android.

**Estado**: Wails v3 está en alpha/beta. El soporte móvil es experimental y no está listo para producción.

**Ventajas**:
- Cambio mínimo en el código (misma API de Wails)
- Mantiene los bindings IPC automáticos

**Desventajas**:
- Wails v3 no está estable, especialmente para móvil
- Riesgo de breaking changes
- Documentación limitada para mobile
- Comunidad pequeña para soporte

**Veredicto**: No recomendada hasta que Wails v3 mobile sea estable.

---

### Opción C: PWA (Progressive Web App)

**Concepto**: Convertir el frontend en una PWA que llame directamente a PokeAPI desde el navegador, eliminando el backend Go para la mayoría de funciones.

**Cambios necesarios**:
1. Mover las llamadas a PokeAPI del backend Go al frontend (fetch directo)
2. Reemplazar el almacenamiento de equipos en archivos por IndexedDB/localStorage
3. Reimplementar la lógica de Core (battle, damage, EVs) en TypeScript
4. Agregar manifest.json y service worker para PWA

**Ventajas**:
- Funciona en cualquier dispositivo con navegador
- No necesita publicar en Play Store
- Sin dependencia de Go en mobile

**Desventajas**:
- Requiere reescribir toda la lógica de Core en TypeScript (~1500 líneas de Go)
- Se pierde el scraping de PokemonDB (CORS impide scraping desde el navegador)
- Sin acceso offline robusto (PokeAPI requiere internet)
- No se puede instalar como app nativa "real"

**Esfuerzo estimado**: Alto (reescritura de core)

---

### Opción D: Go + WebView nativo (sin Capacitor)

**Concepto**: Similar a la Opción A pero usando un proyecto Android nativo con WebView directamente, sin Capacitor.

**Ventajas**:
- Control total sobre el ciclo de vida Android
- Sin dependencia de Capacitor/Ionic

**Desventajas**:
- Más código Android manual (Java/Kotlin)
- Sin el ecosistema de plugins de Capacitor
- Más esfuerzo para features como notificaciones, deep links, etc.

**Esfuerzo estimado**: Medio-Alto

---

### Opción E: Reescritura con Flutter/React Native

**Concepto**: Reescribir el frontend completo en Flutter o React Native, manteniendo el backend Go via FFI/gomobile.

**Ventajas**:
- UI 100% nativa, rendimiento óptimo
- Acceso completo a APIs de Android

**Desventajas**:
- Reescritura completa del frontend (~15 archivos TypeScript, ~200KB de código)
- Curva de aprendizaje de Flutter/Dart o React Native
- Se pierde todo el trabajo de UI existente
- Tiempo de desarrollo muy alto

**Esfuerzo estimado**: Muy alto

---

## 3. Comparativa resumida

| Criterio | A: Capacitor+Go | B: Wails v3 | C: PWA | D: WebView nativo | E: Flutter/RN |
|----------|:---:|:---:|:---:|:---:|:---:|
| Reutilización de código | 90% | 95% | 50% | 85% | 30% |
| Estabilidad del approach | Alta | Baja | Alta | Alta | Alta |
| Esfuerzo de migración | Medio | Bajo* | Alto | Medio-Alto | Muy alto |
| Rendimiento en móvil | Bueno | Desconocido | Bueno | Bueno | Excelente |
| Acceso a APIs nativas | Plugins | Limitado | Mínimo | Total | Total |
| UX nativa | Media | Media | Baja | Media | Alta |
| Mantenibilidad | Alta | Riesgo | Media | Media | Alta |

*\*Si Wails v3 estuviera estable*

---

## 4. Recomendación: Opción A — Capacitor + Go HTTP Server

### Justificación

1. **Máxima reutilización**: El frontend TypeScript/Bootstrap/ECharts/GSAP funciona directamente en WebView. Core y Shell de Go se reutilizan intactos.
2. **Cambios acotados**: Solo hay que crear una capa REST en Go (reemplazando los bindings de Wails) y ajustar las llamadas del frontend.
3. **Madurez**: Tanto Capacitor como gomobile son tecnologías probadas en producción.
4. **Coexistencia**: La versión desktop (Wails) y la versión Android pueden compartir el mismo `core/` y `shell/`.

### Plan de implementación de alto nivel

#### Fase 1: Capa REST en Go
- Crear `app/rest/` con handlers HTTP que expongan los mismos métodos de `bindings.go` como endpoints REST.
- Configurar CORS para `localhost`.
- Compilar con `gomobile bind` generando un `.aar`.

#### Fase 2: Adaptar Frontend
- Crear un módulo `api.ts` que abstraiga las llamadas (actualmente `window.go.app.App.*`) en funciones que hagan `fetch()`.
- Ajustar paths de assets para funcionar dentro de Capacitor.
- Adaptar estilos CSS para mobile (viewport, touch targets, scroll).

#### Fase 3: Proyecto Android con Capacitor
- Inicializar proyecto Capacitor apuntando al frontend Vite.
- Integrar el `.aar` de gomobile como dependencia nativa.
- Crear un plugin Capacitor personalizado que inicie/detenga el servidor Go.
- Configurar el directorio de datos interno para la persistencia de equipos.

#### Fase 4: Optimización mobile
- Adaptar layouts para pantallas pequeñas (responsive).
- Optimizar rendimiento de WebView (lazy loading, virtualización de listas).
- Manejar ciclo de vida Android (pause/resume del servidor Go).
- Splash screen, ícono de app, manifest.

#### Fase 5: Testing y distribución
- Testing en emulador y dispositivos reales.
- Generar APK firmado.
- Preparar para Play Store (si se desea).

### Requisitos técnicos

- **Android Studio** (Arctic Fox o superior)
- **Go 1.24+** con `gomobile` instalado (`go install golang.org/x/mobile/cmd/gomobile@latest`)
- **Node.js 18+** con npm
- **Capacitor CLI** (`npm install @capacitor/cli @capacitor/core @capacitor/android`)
- **Android SDK** (API level 24+ para soporte amplio)

### Riesgos y mitigaciones

| Riesgo | Probabilidad | Mitigación |
|--------|:---:|-----------|
| gomobile no soporta alguna dependencia de Go | Media | goquery usa solo net/http y HTML parsing, que sí están soportados |
| Rendimiento de WebView insuficiente | Baja | Bootstrap ya es ligero; ECharts tiene modo canvas optimizado |
| Almacenamiento de equipos en Android | Baja | Usar `context.getFilesDir()` como base path, inyectado al servidor Go |
| Tamaño del APK | Media | El binario Go (~10-15MB) + frontend (~5MB) = ~20MB total, aceptable |

---

## 5. Estructura propuesta del proyecto post-migración

```
otis/
├── core/                  # SIN CAMBIOS — lógica pura compartida
├── shell/                 # SIN CAMBIOS — adaptadores compartidos
├── app/
│   ├── desktop/           # Entry point Wails (actual main.go)
│   │   ├── main.go
│   │   └── bindings.go
│   └── mobile/            # NUEVO — servidor REST para Android
│       ├── server.go      # HTTP server con endpoints REST
│       ├── handlers.go    # Handlers que llaman a core/shell
│       └── mobile.go      # Exports para gomobile bind
├── frontend/
│   ├── src/
│   │   ├── api.ts         # NUEVO — abstracción fetch/IPC
│   │   └── ...            # Resto del frontend adaptado
│   └── ...
├── android/               # NUEVO — proyecto Capacitor/Android
│   ├── app/
│   │   ├── src/main/java/
│   │   │   └── GoServerPlugin.java  # Plugin que arranca Go
│   │   └── src/main/assets/         # Frontend compilado
│   ├── build.gradle
│   └── capacitor.config.ts
├── go.mod
└── wails.json             # Se mantiene para la versión desktop
```

---

## 6. Conclusión

La arquitectura actual de Otis (Core puro + Shell de adaptadores + frontend web estándar) está **excepcionalmente bien posicionada** para una migración a Android. La separación de capas significa que `core/` y `shell/` se reutilizan sin modificaciones. El frontend, al ser HTML/CSS/JS, funciona directamente en un WebView de Android.

El camino recomendado es **Capacitor + gomobile**: cambios acotados, tecnologías maduras, y máxima reutilización del código existente.
