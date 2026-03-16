# Usar snorlax-seeklogo.svg en la pantalla de carga

**ID**: 0087-splash-snorlax-seeklogo
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Reemplazar el SVG actual del splash screen (`frontend/src/assets/snorlax-splash.svg`, 1.7KB) por el nuevo asset de mayor calidad `assets/snorlax-seeklogo.svg` (25KB). La pantalla de carga ya existe (tarea 0069) con animación de respiración y "zzZ" flotantes; solo hay que cambiar la imagen.

## Capas afectadas

- **APP (frontend)**: Reemplazar el asset SVG y ajustar estilos si el nuevo SVG tiene proporciones o tamaño visual diferente.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/assets/snorlax-splash.svg` | reemplazar | Copiar el contenido de `assets/snorlax-seeklogo.svg` sobre este archivo (o eliminarlo y referenciar el nuevo) |
| `frontend/src/styles/_splash.scss` | modificar | Ajustar tamaño/proporciones de `.splash-snorlax` si el nuevo SVG lo requiere |

## Plan de implementacion

### Parte 1: Reemplazar el asset

1. Copiar `assets/snorlax-seeklogo.svg` a `frontend/src/assets/snorlax-splash.svg`, sobrescribiendo el archivo existente. Así no hay que cambiar la referencia en `index.html`.

### Parte 2: Ajustar estilos si es necesario

2. Comparar visualmente el nuevo SVG con el anterior. Si las proporciones son distintas, ajustar en `_splash.scss` las propiedades de `.splash-snorlax` (`width`, `max-width`, `height`, etc.) para que se vea correctamente centrado y proporcionado.

3. Verificar que la animación de respiración (`scale 1.0 → 1.04`) sigue viéndose bien con el nuevo SVG.

### Parte 3: Verificar dark mode y móvil

4. Comprobar que el SVG se ve bien sobre el fondo oscuro `#1a202c` del splash. Si el SVG tiene fondo blanco o colores que chocan, ajustar el SVG o añadir filtros CSS.

5. Verificar en viewport móvil que el tamaño es adecuado.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | El nuevo SVG se muestra centrado en el splash screen |
| (visual) | La animación de respiración funciona correctamente |
| (visual) | Las "zzZ" flotantes siguen posicionadas correctamente respecto a la imagen |
| (visual) | Se ve bien sobre el fondo oscuro #1a202c |
| (visual) | Proporción correcta en desktop (1024x768) |
| (visual) | Proporción correcta en móvil |
| (visual) | El splash se cierra normalmente tras la carga de la API |

## Criterios de aceptacion

- [x] El splash screen muestra el nuevo SVG `snorlax-seeklogo.svg`
- [x] La imagen está centrada y proporcionada
- [x] La animación de respiración funciona correctamente
- [x] Se ve bien en fondo oscuro
- [x] Responsive: tamaño adecuado en móvil
- [x] El splash se cierra correctamente al terminar la carga

## Notas

- El archivo `assets/snorlax-seeklogo.svg` (25KB) es significativamente más grande que el actual (1.7KB), lo que sugiere mayor detalle/calidad.
- Mantener el nombre `snorlax-splash.svg` en frontend para evitar cambios en `index.html`.
- Si el SVG tiene elementos que no se ven bien a escala pequeña, considerar simplificarlo.
