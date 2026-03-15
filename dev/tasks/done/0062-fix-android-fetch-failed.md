# Fix "TypeError: Failed to fetch" en Android emulador

**ID**: 0062-fix-android-fetch-failed
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Al visualizar la app en el emulador de Android, no se cargan las imágenes ni las tablas (solo el layout). El error reportado es `TypeError: Failed to fetch`. El frontend hace `fetch()` a `http://localhost:8080` pero las peticiones no llegan al servidor Go. Las causas probables son:

1. **Android bloquea tráfico HTTP cleartext** (Android 9+ lo bloquea por defecto). El `AndroidManifest.xml` no tiene `android:usesCleartextTraffic="true"` ni existe `network_security_config.xml`.
2. **Race condition**: el WebView puede cargar antes de que el servidor Go esté listo en el puerto 8080 (el plugin `GoServerPlugin.load()` arranca el server, pero no hay garantía de que esté escuchando cuando el frontend hace el primer fetch).
3. **Capacitor WebView origin**: con `androidScheme: "http"`, el WebView carga desde `http://localhost` pero las peticiones van a `http://localhost:8080` — puede haber restricciones de mixed content o CORS específicas del WebView Android.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Configuración Android (manifest, network config) y posiblemente frontend (retry logic).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `android/app/src/main/AndroidManifest.xml` | modificar | Añadir `android:usesCleartextTraffic="true"` y referencia a `networkSecurityConfig` |
| `android/app/src/main/res/xml/network_security_config.xml` | crear | Permitir cleartext traffic a `localhost` y `10.0.2.2` (emulador) |
| `frontend/capacitor.config.ts` | modificar | Evaluar si `androidScheme` necesita cambio o si hay que añadir `allowNavigation` para localhost:8080 |
| `frontend/src/api.ts` | modificar | Añadir retry con backoff para el primer fetch (esperar a que el servidor Go esté listo) |
| `android/app/src/main/java/com/alfon/otis/GoServerPlugin.java` | modificar | Añadir health-check loop que confirme que el server está escuchando antes de notificar "listo" |

## Plan de implementacion

1. **Network Security Config**: Crear `network_security_config.xml` que permita cleartext a `localhost` y `10.0.2.2`.
2. **AndroidManifest**: Añadir `android:usesCleartextTraffic="true"` y `android:networkSecurityConfig="@xml/network_security_config"` al tag `<application>`.
3. **Health-check en GoServerPlugin**: Tras llamar a `Mobile.start()`, hacer un loop con HTTP GET a `http://localhost:8080/api/pokemon?offset=0&limit=1` (o un endpoint `/health`) hasta que responda 200, con timeout de 10s.
4. **Retry en frontend**: En `api.ts`, envolver las funciones `get/post/put/del` con un retry (3 intentos, backoff 500ms→1s→2s) para tolerar que el server aún no esté listo.
5. **Verificar Capacitor config**: Confirmar que `androidScheme: "http"` es correcto y que no se necesita `server.url` o `server.allowNavigation`.
6. **Probar en emulador**: Verificar que las tablas e imágenes se cargan correctamente.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual en emulador | Las tablas de Pokémon cargan datos correctamente |
| Test manual en emulador | Las imágenes/sprites se muestran |
| Test manual en emulador | No aparece "TypeError: Failed to fetch" en logcat |
| Test manual en emulador | La app funciona tras cold start (sin race condition) |
| Test en desktop (Wails) | Verificar que los cambios en api.ts no rompen el modo desktop |

## Criterios de aceptacion

- [x] No aparece `TypeError: Failed to fetch` al usar la app en el emulador Android
- [x] Las tablas muestran datos (lista de Pokémon, movimientos, etc.)
- [x] Las imágenes/sprites se cargan (CDN externo o fallback local)
- [x] La app desktop (Wails) sigue funcionando sin regresiones
- [x] El servidor Go está confirmado como "listo" antes de que el frontend haga peticiones

## Notas

- `10.0.2.2` es la IP del host desde el emulador Android (equivalente a `localhost` del host).
- El `androidScheme: "http"` en Capacitor hace que el WebView use `http://localhost` como origin, lo cual puede causar que Android bloquee las peticiones si no está configurado `usesCleartextTraffic`.
- Considerar añadir un endpoint `/api/health` al servidor Go para health-checks rápidos.
- Depende de: 0060 (proyecto Android con Capacitor) y 0058 (abstracción API frontend).
