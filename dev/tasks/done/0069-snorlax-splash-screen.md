# 0069 — Splash screen animado con Snorlax al iniciar la app

## Descripción

Al abrir el APK en el móvil hay una pantalla negra durante varios segundos mientras el servidor Go arranca y la WebView carga. Reemplazar esa pantalla negra por una splash screen con una animación de Snorlax (estilo dormido/despertando) que se muestre durante toda la fase de carga y transicione suavemente al contenido de la app.

La implementación tiene dos partes:
1. **Splash nativa Android**: reemplazar el drawable `splash.png` actual por una imagen de Snorlax para eliminar el pantallazo negro inmediato.
2. **Splash animada en frontend**: añadir una pantalla de carga HTML/CSS/GSAP con animación de Snorlax que se muestre mientras el servidor Go inicia y la app se hidrata. Al estar disponible la API, la splash se desvanece y muestra la app.

## Capas afectadas

- **Core**: ningún cambio.
- **Shell**: ningún cambio.
- **APP (Android nativo)**: reemplazar los drawables `splash.png` en todas las densidades.
- **APP (frontend)**: nueva pantalla de carga animada con lógica de espera al servidor.

## Cambios requeridos

### 1. Splash nativa Android — drawables

Reemplazar `android/app/src/main/res/drawable*/splash.png` en todas las densidades con una imagen estática de Snorlax sobre fondo oscuro (#1a202c o similar) que coincida con el tema de la app:

| Directorio | Resolución aproximada |
|---|---|
| `drawable/` | 480×320 |
| `drawable-land-hdpi/` | según actual |
| `drawable-land-mdpi/` | según actual |
| `drawable-land-xhdpi/` | según actual |
| `drawable-land-xxhdpi/` | según actual |
| `drawable-land-xxxhdpi/` | según actual |
| `drawable-port-hdpi/` | según actual |
| `drawable-port-mdpi/` | según actual |
| `drawable-port-xhdpi/` | según actual |
| `drawable-port-xxhdpi/` | según actual |
| `drawable-port-xxxhdpi/` | según actual |

Actualizar `styles.xml` si es necesario para que el fondo del tema `AppTheme.NoActionBarLaunch` use el color de fondo correcto.

### 2. Splash animada en frontend

Crear la animación en el frontend que se muestra sobre la WebView mientras espera al servidor:

**a) HTML (`frontend/index.html`)**:
- Añadir un `<div id="splash-screen">` antes del `<div id="app">` con:
  - Imagen/SVG de Snorlax centrado.
  - Texto sutil tipo "Cargando..." o "Snorlax está despertando...".
  - Fondo que coincida con el splash nativo (#1a202c).

**b) CSS (`frontend/src/styles/`)**:
- Crear `_splash.scss` con estilos del splash:
  - Pantalla completa, centrado, z-index alto.
  - Animación CSS de Snorlax: respiración idle (scale sutil), zzZ flotantes, o similar.
  - Transición de salida (fade-out).

**c) JS/TS (`frontend/src/`)**:
- En el flujo de inicialización existente (`api.ts` o `main.ts`), cuando la API responda exitosamente:
  - Añadir clase de fade-out al splash.
  - Tras la transición, eliminar el elemento del DOM.
- Si ya existe lógica de retry/espera al servidor en `api.ts`, aprovecharla para coordinar la desaparición del splash.

### 3. Asset de Snorlax

- Crear o conseguir un SVG/PNG de Snorlax con licencia adecuada para uso en la app.
- Ubicar en `frontend/src/assets/snorlax-splash.svg` (o `.png`).
- El mismo asset (en formato PNG) se usa para los drawables nativos de Android.

## Plan de implementación

1. Diseñar/obtener el asset de Snorlax (SVG para frontend, PNG para Android).
2. Generar las imágenes `splash.png` en todas las densidades de Android.
3. Reemplazar los drawables existentes.
4. Ajustar `styles.xml` si el fondo necesita cambio.
5. Añadir el markup del splash en `index.html`.
6. Crear `_splash.scss` con las animaciones.
7. Integrar la lógica de ocultación del splash con el flujo de inicialización existente.
8. Probar en emulador y dispositivo físico.

## Tests

- Instalar el APK en un dispositivo/emulador y verificar que:
  - Al abrir la app se ve inmediatamente la imagen de Snorlax (no pantalla negra).
  - La animación de Snorlax se reproduce mientras carga.
  - Cuando la app está lista, la splash desaparece con fade-out suave.
  - No hay parpadeo ni salto entre la splash nativa y la splash del frontend.
- Verificar en modo claro y oscuro.
- Verificar en orientación portrait (principal).

## Dependencias

- Asset gráfico de Snorlax (SVG/PNG).
- GSAP ya está disponible en el proyecto si se quiere usar para la animación (aunque CSS puro puede ser suficiente).

## Notas

- Coordinar el color de fondo de la splash nativa con el de la splash del frontend para que la transición sea imperceptible.
- Mantener la animación ligera (CSS preferred) para no retrasar la carga.
- El texto "Snorlax está despertando..." es una sugerencia; ajustar al gusto del usuario.
