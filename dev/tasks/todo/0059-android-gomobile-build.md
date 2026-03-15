# Compilar backend Go con gomobile para Android (.aar)

**ID**: 0059-android-gomobile-build
**Estado**: todo
**Fecha**: 2026-03-15

---

## Descripcion

Configurar la compilación del backend Go como una biblioteca Android (.aar) usando gomobile. El archivo .aar contendrá el servidor HTTP REST (tarea 0057) y será invocable desde Java/Kotlin para iniciar y detener el servidor.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Verificar compatibilidad con gomobile (goquery usa net/http, debería funcionar).
- **APP**: `app/mobile/mobile.go` — funciones exportadas para gomobile.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `app/mobile/mobile.go` | modificar | Asegurar que las funciones exportadas cumplen restricciones de gomobile |
| `Makefile` o `scripts/build-android.sh` | crear | Script de build: `gomobile bind -target=android app/mobile` |
| `go.mod` | modificar | Agregar dependencia `golang.org/x/mobile` si es necesario |

## Plan de implementacion

1. Instalar gomobile: `go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init`
2. Verificar que `app/mobile/mobile.go` exporta funciones con tipos compatibles (solo tipos primitivos y strings)
3. Crear script de build: `gomobile bind -target=android -o android/app/libs/otis.aar ./app/mobile`
4. Compilar y resolver errores de compatibilidad
5. Verificar que el .aar se genera correctamente

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Build CI | El .aar se genera sin errores |
| Test manual | El .aar se puede importar en un proyecto Android Studio |

## Criterios de aceptacion

- [ ] `gomobile bind` genera `otis.aar` sin errores
- [ ] El .aar exporta funciones `Start(port int, dataDir string)` y `Stop()`
- [ ] El tamaño del .aar es razonable (<20MB)
- [ ] Documentado el proceso de build en scripts/

## Notas

- gomobile solo soporta tipos exportados simples: int, float, string, bool, []byte, error.
- Si goquery causa problemas con gomobile, considerar mover el scraping a un endpoint REST externo.
- Depende de: 0057 (capa REST).
