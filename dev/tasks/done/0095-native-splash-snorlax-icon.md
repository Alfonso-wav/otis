# 0095 — Mostrar Snorlax en la splash nativa desde el primer instante de arranque

## Descripción

Al abrir la app en el móvil hay varios segundos de pantalla oscura (#1a202c) antes de que la WebView cargue y muestre la splash animada de Snorlax del frontend. El Snorlax aparece solo un instante antes de que la app esté lista, haciendo que toda la espera sea pantalla vacía.

El objetivo es que el Snorlax sea visible desde el primer frame del arranque, usando la API nativa de Android 12 Splash Screen para mostrar el icono de Snorlax centrado sobre el fondo oscuro. Así la transición será: Snorlax nativo → Snorlax frontend animado → app lista.

## Capas afectadas

- **Core**: ningún cambio.
- **Shell**: ningún cambio.
- **APP (Android nativo)**: configurar el tema splash con `windowSplashScreenAnimatedIcon` y `postSplashScreenTheme`.

## Cambios requeridos

### 1. Crear drawable del icono Snorlax para splash

Crear un drawable circular/adaptable para el icono central de la splash screen nativa:

- Crear `android/app/src/main/res/drawable/splash_icon.png` (o SVG → VectorDrawable) con el Snorlax sobre fondo transparente.
- Tamaño recomendado: 288×288 dp (el área visible es un círculo de 216dp según spec de Android 12).
- Generar versiones para distintas densidades si es necesario (`drawable-hdpi/`, `drawable-xhdpi/`, etc.), o usar un único drawable de alta resolución.

### 2. Actualizar `styles.xml` — tema de splash

Modificar `android/app/src/main/res/values/styles.xml`:

```xml
<style name="AppTheme.NoActionBarLaunch" parent="Theme.SplashScreen">
    <item name="windowSplashScreenBackground">#1a202c</item>
    <item name="windowSplashScreenAnimatedIcon">@drawable/splash_icon</item>
    <item name="windowSplashScreenIconBackgroundColor">#1a202c</item>
    <item name="postSplashScreenTheme">@style/AppTheme.NoActionBar</item>
</style>
```

Atributos clave:
- `windowSplashScreenBackground`: fondo oscuro que coincide con el frontend.
- `windowSplashScreenAnimatedIcon`: el icono de Snorlax centrado.
- `windowSplashScreenIconBackgroundColor`: fondo del icono para que no haya círculo blanco.
- `postSplashScreenTheme`: tema que se aplica después de la splash.

### 3. Verificar dependencia de `core-splashscreen`

En `android/app/build.gradle`, asegurar que está la dependencia:

```groovy
implementation 'androidx.core:core-splashscreen:1.0.1'
```

### 4. Verificar `MainActivity` — installSplashScreen

En `MainActivity.java`/`.kt`, verificar que se llama a `installSplashScreen()` antes de `super.onCreate()`. Si usa Capacitor, esto puede estar ya gestionado por el plugin, pero hay que confirmar.

### 5. Opcional — Ajustar duración

Si la splash nativa desaparece muy rápido (antes de que la WebView cargue), considerar usar `setKeepOnScreenCondition` en `MainActivity` para mantener la splash nativa visible hasta que el servidor Go esté listo. Esto evitaría el flash de pantalla vacía entre splash nativa y splash frontend.

## Plan de implementación

1. Extraer/generar PNG de Snorlax sobre fondo transparente desde el SVG existente (`frontend/src/assets/snorlax-splash.svg`).
2. Crear los drawables `splash_icon.png` en las densidades necesarias.
3. Actualizar `styles.xml` con los atributos de splash screen.
4. Verificar/añadir dependencia `core-splashscreen` en `build.gradle`.
5. Verificar que `MainActivity` llama a `installSplashScreen()`.
6. Compilar APK y probar en dispositivo físico.

## Tests

- Instalar APK en dispositivo/emulador y verificar que:
  - Al abrir la app se ve inmediatamente Snorlax centrado sobre fondo oscuro (no pantalla vacía).
  - La transición de splash nativa a splash frontend animada es fluida (mismo fondo, mismo personaje).
  - La splash frontend sigue funcionando con su animación de respiración y zzZ.
  - No hay regresión en el flujo normal de carga.

## Dependencias

- Asset SVG de Snorlax ya disponible en `frontend/src/assets/snorlax-splash.svg`.
- Librería `androidx.core:core-splashscreen` (probablemente ya incluida por Capacitor).

## Notas

- El icono de la splash de Android 12 se muestra dentro de un círculo de 216dp. Asegurar que el Snorlax quepa bien sin recortes.
- Mantener consistencia de color #1a202c en todo el flujo para evitar parpadeos.
- Los drawables `splash.png` existentes en las distintas densidades (landscape/portrait) pueden quedarse como fallback para Android <12, pero el tema `Theme.SplashScreen` de AndroidX gestiona la compatibilidad hacia atrás.
