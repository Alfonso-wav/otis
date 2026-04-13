# Splash: separar "OTIS" del Snorlax y eliminar su respiracion

**ID**: 0160-splash-otis-position-and-static
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

En el splash, el texto `OTIS` aparece demasiado pegado al Snorlax y ademas "respira" (escala y sube/baja) junto con el, porque esta anidado dentro de `.splash-snorlax-wrapper`, que tiene `animation: snorlax-breathe`. Resultado: el label hereda la transformacion del padre.

Se quieren dos cosas:

1. Separarlo un poco mas (posicion un poco mas arriba respecto al Snorlax).
2. Que se quede estatico — no respira, no escala, no sube/baja.

## Capas afectadas

- **Core**: ninguna.
- **Shell**: ninguna.
- **APP**: solo frontend (HTML + CSS del splash).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Mover `<div class="splash-otis">` fuera de `.splash-snorlax-wrapper` para que no herede la animacion de respiracion. Dejarlo como hermano posicionado encima del wrapper. |
| `frontend/index.html` (bloque `<style>` inline) | modificar | Ajustar `.splash-otis` para nueva posicion (fuera del wrapper). Aumentar separacion vertical respecto al Snorlax. |
| `frontend/src/styles/_splash.scss` | modificar | Mismos ajustes que en el CSS inline del `index.html` (mantener paridad entre ambos). Mover `.splash-otis` fuera del contexto de respiracion, ajustar `top`/`bottom` o margen para que quede un poco mas arriba. |
| `frontend/src/main.ts` | verificar | Los selectores `.splash-otis .splash-typed` y `.splash-otis .splash-caret` (lineas 36, 49-50) deben seguir funcionando tras mover el nodo. No deberia requerir cambios si el nodo sigue existiendo como descendiente de `#splash-screen`. |

## Plan de implementacion

1. En `frontend/index.html`, sacar `<div class="splash-otis">...</div>` del interior de `.splash-snorlax-wrapper` y colocarlo como hermano anterior dentro de `#splash-screen` (o envolver ambos en un contenedor estatico que agrupe OTIS + wrapper).
   - Opcion A (minima): OTIS como hermano directo del wrapper, con `position: relative` o similar, sin `absolute` anclado al wrapper.
   - Opcion B: crear contenedor `.splash-stage` con `display: flex; flex-direction: column; align-items: center;` que agrupa OTIS + wrapper; OTIS queda estatico arriba y solo el wrapper respira.
2. Actualizar reglas CSS (tanto en el `<style>` inline del `index.html` como en `frontend/src/styles/_splash.scss`):
   - Quitar `position: absolute; top: -2rem; left: 50%; transform: translateX(-50%);` si pasa a flow normal.
   - Aplicar separacion mayor: p.ej. `margin-bottom: 1.25rem` (o `1.5rem`) bajo OTIS para que quede visiblemente mas arriba del Snorlax.
   - Confirmar que no queda herencia de `animation` ni `transform` del wrapper.
3. Repetir ajuste en el media query `@media (max-width: 400px)` para mobile (reducir separacion proporcional, p.ej. `margin-bottom: 0.9rem`, `font-size: 1.2rem`).
4. Verificar typewriter: el efecto de escritura de `OTIS` lo dispara `main.ts` buscando `.splash-otis .splash-typed`. Al mover el nodo, el selector sigue siendo valido → no requiere cambios JS.
5. Comprobar que el label `POKéDEX` (`.splash-pokedex-label`, debajo del Snorlax) no se ve afectado ni hereda cambios no deseados.

## Tests

| Tipo | Que se verifica |
|------|-----------------|
| Manual desktop | Al abrir la app, el texto `OTIS` queda estatico (no escala ni sube/baja) mientras el Snorlax respira |
| Manual desktop | `OTIS` tiene mas separacion visual respecto al Snorlax que antes |
| Manual desktop | Typewriter de `OTIS` sigue animandose (caracter a caracter) y el caret parpadea |
| Manual mobile 360px | Separacion proporcional correcta, sin tapar el Snorlax ni salirse del viewport |
| Manual APK Android | Mismo resultado que en build web |
| Manual | `POKéDEX` label debajo sigue apareciendo con su typewriter y no se ha movido |

## Criterios de aceptacion

- [ ] `OTIS` no escala ni se mueve con el ciclo de respiracion del Snorlax
- [ ] `OTIS` queda un poco mas arriba (mayor separacion vertical) respecto al Snorlax
- [ ] Typewriter y caret de `OTIS` siguen funcionando
- [ ] `POKéDEX` label inferior no se ve afectado
- [ ] Funciona en desktop, mobile 360px y APK Android
- [ ] Paridad entre el CSS inline de `index.html` y `_splash.scss`

## Notas

- El CSS del splash vive duplicado: una copia inline en `index.html` (para evitar FOUC antes de cargar el bundle) y otra en `_splash.scss`. Ambos deben mantenerse en paridad.
- No tocar la animacion `snorlax-breathe`; solo sacar `OTIS` de su ambito.
- Si la opcion B (contenedor flex) requiere reposicionar el bloque `zzz` flotante, mantener su anclaje al Snorlax (siguen siendo hijos del wrapper).
