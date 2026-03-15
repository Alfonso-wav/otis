# Fix: fallback CDN del simulador de batalla usa sprites Home en vez de sprites de batalla

**ID**: 0039-battle-sprites-cdn-fallback
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

El simulador de batalla en `builds.ts` tiene una cadena de fallback para sprites que termina en el CDN de pokemondb.net con sprites `home/normal` (modernos, 3D). Cuando no hay sprites descargados localmente (que es el caso habitual ya que no se versionan), **siempre se muestran los sprites Home** en vez de los sprites de batalla pixelados estilo Game Boy.

**Estado actual del fallback** en `battleSpriteImg()`:
1. `/assets/sprites/battle-front/{name}.png` — local, no existe
2. `/assets/sprites/home-normal/{name}.png` — local, no existe
3. `https://img.pokemondb.net/sprites/home/normal/{name}.png` — **este se muestra** (sprite moderno)

**Estado deseado**:
1. `/assets/sprites/battle-front/{name}.png` — local
2. CDN sprite de batalla pixelado de pokemondb.net (Black & White tiene la mejor cobertura Gen 1-5)
   - Front: `https://img.pokemondb.net/sprites/black-white/anim/normal/{name}.gif`
   - Back: `https://img.pokemondb.net/sprites/black-white/anim/back-normal/{name}.gif`
3. `/assets/sprites/home-normal/{name}.png` — fallback local home
4. `https://img.pokemondb.net/sprites/home/normal/{name}.png` — fallback CDN home (para Pokemon Gen 6+ sin sprites pixelados)

De esta forma, si no hay sprites locales, se cargan los sprites animados de Black & White desde el CDN (que son pixelados y encajan con el estilo Game Boy del simulador).

## Capas afectadas
- **Frontend**: Modificar `battleSpriteImg()` y `battleSpriteURL()` en `builds.ts`

## Dependencias externas nuevas
Ninguna.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | Actualizar `battleSpriteImg()` para añadir fallback CDN de sprites de batalla (Black & White animados) antes de caer a Home sprites. Ajustar la cadena de `onerror` para soportar 4 niveles de fallback. |

## Plan de implementacion

### Paso 1 — Añadir función de URL CDN de batalla
1. En `builds.ts`, añadir una función `battleSpriteFallbackCDN(name, type)` que devuelva:
   - Para `"battle-front"`: `https://img.pokemondb.net/sprites/black-white/anim/normal/{safeName}.gif`
   - Para `"battle-back"`: `https://img.pokemondb.net/sprites/black-white/anim/back-normal/{safeName}.gif`

### Paso 2 — Actualizar cadena de fallback en battleSpriteImg
2. En `battleSpriteImg()`, modificar la cadena de `onerror` para que:
   - Fallback 0→1: de local battle sprite a **CDN battle sprite** (Black & White animado)
   - Fallback 1→2: de CDN battle a **local home-normal**
   - Fallback 2→3: de local home a **CDN home**
   - Fallback 3: ocultar sprite

### Paso 3 — Verificar
3. Ejecutar la app y verificar que el simulador de batalla muestra sprites pixelados animados de Black & White en vez de los sprites Home modernos.

## Tests
| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que el simulador muestra sprites pixelados (BW) cuando no hay sprites locales |
| Manual | Verificar que Pokemon de Gen 6+ caen a sprites Home correctamente |

## Criterios de aceptacion
- [x] `battleSpriteImg()` tiene 4 niveles de fallback (local battle → CDN battle BW → local home → CDN home)
- [x] Los sprites CDN de batalla usan la ruta de Black & White animados de pokemondb.net
- [x] El simulador muestra sprites pixelados por defecto (sin necesidad de descargar sprites localmente)
- [x] Si un Pokemon no tiene sprite BW (Gen 6+), cae correctamente a sprites Home
