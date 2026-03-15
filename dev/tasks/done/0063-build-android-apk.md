# Generar APK instalable para Android

**ID**: 0063-build-android-apk
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Crear un script de build unificado que produzca un `.apk` instalable en un dispositivo Android físico. Actualmente existen todas las piezas (frontend Vite, gomobile .aar, proyecto Capacitor/Android), pero no hay un pipeline que las encadene y genere el APK final. El script debe:

1. Compilar el frontend (`npm run build` → `frontend/dist/`).
2. Compilar el backend Go como `.aar` con gomobile (`scripts/build-android.sh`).
3. Sincronizar Capacitor (`npx cap sync android`) para copiar `dist/` al proyecto Android.
4. Compilar el proyecto Android con Gradle para generar el APK debug (y opcionalmente release).

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Script de build, posibles ajustes en configuración Gradle.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `scripts/build-apk.sh` | crear | Script unificado que ejecuta los 4 pasos del build pipeline |
| `android/app/build.gradle` | modificar | Verificar/ajustar configuración de signing, versionName, minSdk, targetSdk |
| `android/build.gradle` | verificar | Asegurar que las dependencias de Gradle y Android Gradle Plugin son correctas |
| `android/gradle.properties` | verificar | Asegurar configuración de memoria y AndroidX |
| `scripts/build-android.sh` | verificar | Confirmar que genera `android/app/libs/otis.aar` correctamente |

## Plan de implementacion

1. **Verificar prerequisitos**: Confirmar que `gomobile`, `Android SDK`, `NDK`, y `Gradle` están disponibles en el sistema. Documentar versiones requeridas en el script.
2. **Crear `scripts/build-apk.sh`**:
   - Paso 1: `cd frontend && npm run build` — genera `frontend/dist/`.
   - Paso 2: `bash scripts/build-android.sh` — genera `android/app/libs/otis.aar`.
   - Paso 3: `cd frontend && npx cap sync android` — copia dist/ y actualiza plugins en el proyecto Android.
   - Paso 4: `cd android && ./gradlew assembleDebug` — genera el APK en `android/app/build/outputs/apk/debug/app-debug.apk`.
   - Mostrar ruta del APK generado al final.
3. **Revisar `build.gradle`**:
   - `minSdkVersion` al menos 21 (gomobile requiere API 21+).
   - `targetSdkVersion` 34 (o la más reciente disponible).
   - `versionCode` y `versionName` correctos.
   - Dependencia de `otis.aar` como `implementation fileTree(dir: 'libs', include: ['*.aar'])`.
4. **Build debug APK**: Ejecutar el script y verificar que genera el APK.
5. **Copiar APK al usuario**: Indicar ruta final del APK (`android/app/build/outputs/apk/debug/app-debug.apk`) para transferir al dispositivo e instalar.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual | El script `build-apk.sh` completa sin errores |
| Test manual | El APK se genera en la ruta esperada |
| Test manual | El APK se instala en un dispositivo Android físico |
| Test manual | La app abre, carga datos y muestra sprites correctamente |
| Test manual | El servidor Go arranca correctamente dentro del APK |

## Criterios de aceptacion

- [x] Existe un script `scripts/build-apk.sh` que genera el APK en un solo comando
- [x] El APK debug se genera en `android/app/build/outputs/apk/debug/app-debug.apk`
- [x] El APK se puede instalar en un dispositivo Android físico (API 21+)
- [x] La app funciona correctamente: carga Pokémon, muestra sprites, permite batallas
- [x] El script muestra errores claros si falta algún prerequisito (SDK, NDK, gomobile)

## Notas

- Para instalar en un dispositivo físico se necesita habilitar "Instalar desde fuentes desconocidas" en ajustes Android.
- El APK debug no necesita signing con keystore. Para distribución en Play Store se necesitaría un APK release firmado (fuera del alcance de esta tarea).
- Se puede transferir el APK al dispositivo por USB (`adb install app-debug.apk`) o compartir el archivo directamente.
- Prerequisitos del sistema: Android Studio (SDK + NDK), Go 1.25+, gomobile, Node.js 18+.
- Depende de: 0059 (gomobile build), 0060 (proyecto Capacitor), 0061 (UX mobile), 0062 (fix fetch).
