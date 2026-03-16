# 0070 — Icono de app: Pokeball pixel art

## Descripción

Reemplazar el icono por defecto de la app Android con una Pokeball en estilo pixel art. El icono debe verse nítido en todas las densidades de pantalla y funcionar tanto en formato normal como round.

## Capas afectadas

- **APP** — recursos Android (`android/app/src/main/res/`)

## Requisitos

1. Crear una imagen de Pokeball en estilo pixel art (fondo transparente para el foreground, color sólido para el background).
2. Generar los PNGs en todas las densidades requeridas por Android:
   - `mipmap-mdpi` — 48×48 px (foreground 108×108)
   - `mipmap-hdpi` — 72×72 px (foreground 162×162)
   - `mipmap-xhdpi` — 96×96 px (foreground 216×216)
   - `mipmap-xxhdpi` — 144×144 px (foreground 324×324)
   - `mipmap-xxxhdpi` — 192×192 px (foreground 432×432)
3. Reemplazar los archivos existentes:
   - `ic_launcher.png` en cada carpeta `mipmap-*`
   - `ic_launcher_round.png` en cada carpeta `mipmap-*`
   - `ic_launcher_foreground.png` en cada carpeta `mipmap-*`
4. Actualizar `ic_launcher_background.xml` con un color rojo/blanco apropiado para Pokeball si es necesario.
5. Mantener el sistema de adaptive icons (`mipmap-anydpi-v26/ic_launcher.xml`).

## Diseño del icono

- **Pokeball clásica** en pixel art: mitad superior roja, mitad inferior blanca, línea negra horizontal en el centro, círculo/botón central blanco con borde negro.
- Estilo retro/pixelado fiel a la estética de los juegos clásicos de Pokémon (Game Boy / GBA).
- El pixel art debe escalar bien sin anti-aliasing (nearest neighbor scaling).

## Notas

- Los PNGs se pueden generar programáticamente con un script o con una herramienta de edición de imágenes.
- Considerar usar un canvas HTML o un script en Go/Python para generar los PNGs en todas las resoluciones desde un diseño base.
- No se requieren cambios en el código Go ni en el frontend, solo en los recursos de Android.

## Tests

- Verificar que el icono se muestra correctamente en el launcher del emulador Android.
- Verificar que tanto el icono normal como el round se ven bien.
- Compilar APK con `/build-apk` y confirmar el icono.
