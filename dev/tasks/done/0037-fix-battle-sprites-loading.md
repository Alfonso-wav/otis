# Fix: sprites de batalla no cargan en el simulador

**ID**: 0037-fix-battle-sprites-loading
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Los sprites del simulador de batallas no cargan. El problema tiene dos causas:

1. **El directorio `assets/sprites/` no existe** en el proyecto — los sprites locales (battle-back, battle-front, home-normal) no se encuentran.
2. **Bug en el handler `onerror` de las imágenes**: la comparación `this.src !== '${fallbackSrc}'` siempre es `true` porque `this.src` devuelve la URL absoluta resuelta por el navegador (ej. `http://wails.localhost/assets/sprites/home-normal/pikachu.png`) mientras que `fallbackSrc` es la ruta relativa (`/assets/sprites/home-normal/pikachu.png`). Esto causa un loop infinito entre el primer error y el fallback local, sin nunca llegar al fallback CDN.

## Capas afectadas
- **Frontend**: Corregir la lógica de fallback en `battleSpriteImg()` y en cualquier otro `onerror` de sprites en `builds.ts`.

## Dependencias externas nuevas
Ninguna.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | Corregir `onerror` handler en `battleSpriteImg()` y otros usos de sprites para que el fallback CDN funcione correctamente |

## Plan de implementacion

### Fase 1 — Corregir fallback en `battleSpriteImg()` (líneas 364-372)
1. Cambiar la comparación de `this.src` para que use `this.getAttribute('src')` (devuelve la ruta relativa original) o usar `this.src.endsWith(...)` o un data-attribute para rastrear el intento de fallback.
2. Enfoque recomendado: usar un `data-fallback` attribute counter para evitar loops:
   - Intento 0 (src original falla): cambiar a sprite home-normal local
   - Intento 1 (home-normal falla): cambiar a CDN
   - Intento 2 (CDN falla): ocultar imagen o mostrar placeholder
3. Aplicar el mismo fix a cualquier otro `onerror` handler de sprites en el archivo (buscar todos los `onerror` en el archivo).

### Fase 2 — Verificar otros usos de sprites
4. Revisar todas las funciones que generan `<img>` con sprites en `builds.ts` y asegurar que los fallbacks funcionan:
   - `battleSpriteImg()` (línea 364)
   - Sprites en la sección de selección de Pokémon (buscar otros `onerror` en el archivo)

## Tests
| Archivo | Que se testea |
|---------|---------------|
| Manual | Con `assets/sprites/` vacío o inexistente, los sprites deben cargar desde CDN |
| Manual | Con sprites locales descargados, deben cargar desde local |
| Manual | Sin conexión a internet y sin sprites locales, no debe haber loop infinito |

## Criterios de aceptacion
- [x] Los sprites de batalla cargan correctamente usando el fallback CDN cuando no hay sprites locales
- [x] No hay loop infinito en el handler `onerror`
- [x] Si hay sprites locales descargados, se usan primero
- [x] La cadena de fallback es: battle sprite local → home-normal local → CDN → placeholder/ocultar
- [x] El fix se aplica a todos los `onerror` de sprites en `builds.ts`, no solo a `battleSpriteImg()`
