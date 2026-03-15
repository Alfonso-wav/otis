# Crear proyecto Android con Capacitor e integrar Go server

**ID**: 0060-android-capacitor-project
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Inicializar un proyecto Capacitor con target Android, integrar el .aar de gomobile como dependencia nativa, y crear un plugin Capacitor personalizado que arranque/detenga el servidor Go HTTP al iniciar la app.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Nueva carpeta `android/` con el proyecto Capacitor/Android.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `capacitor.config.ts` | crear | Configuración de Capacitor (appId, appName, webDir) |
| `android/` | crear | Proyecto Android generado por Capacitor |
| `android/app/libs/otis.aar` | copiar | Biblioteca Go compilada (de tarea 0059) |
| `android/app/build.gradle` | modificar | Agregar dependencia al .aar local |
| `android/app/src/main/java/.../GoServerPlugin.java` | crear | Plugin que arranca el servidor Go al iniciar |
| `android/app/src/main/java/.../MainActivity.java` | modificar | Registrar el plugin |
| `frontend/vite.config.ts` | modificar | Ajustar outDir para Capacitor si es necesario |

## Plan de implementacion

1. Instalar Capacitor: `npm install @capacitor/cli @capacitor/core @capacitor/android`
2. Inicializar: `npx cap init otis com.alfon.otis --web-dir frontend/dist`
3. Agregar Android: `npx cap add android`
4. Copiar `otis.aar` a `android/app/libs/`
5. Crear `GoServerPlugin.java` que llame a `Mobile.Start(port, dataDir)` en `onCreate`
6. Registrar plugin en `MainActivity`
7. Configurar permisos en AndroidManifest.xml (INTERNET)
8. Build y test en emulador: `npx cap run android`

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual emulador | La app arranca, el servidor Go responde, el frontend carga |

## Criterios de aceptacion

- [x] La app Android compila sin errores
- [x] El servidor Go arranca automáticamente al abrir la app
- [x] El frontend carga en el WebView y muestra la Pokédex
- [x] Los equipos se guardan en el directorio interno de Android
- [x] El servidor se detiene al cerrar la app

## Notas

- El appId propuesto es `com.alfon.otis`.
- Mínimo API level 24 (Android 7.0) para compatibilidad amplia.
- Usar `context.getFilesDir()` como dataDir para el almacenamiento de equipos.
- Depende de: 0057, 0058, 0059.
