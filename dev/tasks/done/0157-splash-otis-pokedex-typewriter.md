# Splash Snorlax: escritura progresiva "OTIS" y "POKéDEX"

**ID**: 0157-splash-otis-pokedex-typewriter
**Estado**: todo
**Fecha**: 2026-04-13

---

## Descripcion

Añadir al splash interactivo del Snorlax dos textos que se escriben letra a letra (efecto typewriter) de forma lenta y progresiva:

- `OTIS` encima del Snorlax (sobre el SVG, centrado horizontalmente).
- `POKéDEX` debajo del `<p class="splash-zzz-text">` (las Zzz), tambien centrado.

El splash actual (`frontend/index.html`) ya contiene `#splash-screen` con `.splash-snorlax-wrapper` (SVG Snorlax + arrows + ojos overlay) y `.splash-zzz-text` con las 7 Zzz. La logica vive en `frontend/src/main.ts` (`showSplashInteractive` / `dismissSplashInteractive`). Todo el CSS critico del splash esta inlineado en `<head>` del `index.html` para que pinte antes del bundle.

Flujo propuesto:
1. Al montar el splash, `OTIS` empieza a escribirse letra por letra (p. ej. 180-220ms por letra).
2. Cuando termina `OTIS`, comienza `POKéDEX` debajo de las Zzz con la misma cadencia.
3. El usuario puede hacer click en el Snorlax en cualquier momento — si los textos aun no terminaron, se completan al instante y luego se dismissa el splash.

## Capas afectadas

- **Core**: ninguna.
- **Shell**: ninguna.
- **APP**: solo frontend (HTML + TS + CSS inline).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Añadir `<span class="splash-otis"></span>` dentro de `.splash-snorlax-wrapper` (posicionado `absolute` sobre el SVG) y `<span class="splash-pokedex-label"></span>` dentro de `#splash-screen` despues de `.splash-zzz-text`. Añadir CSS inline para tipografia, posicion y cursor de escritura. |
| `frontend/src/main.ts` | modificar | En `showSplashInteractive`, arrancar una secuencia typewriter sobre los dos spans. Exponer `cancelTypewriter()` que se invoca al primer click para completar el texto de inmediato antes de la animacion de dismiss. |
| `frontend/src/styles/_splash.scss` | modificar | Estilos finales no criticos (fonts, sombras, responsive mobile). Fuente coherente con el branding (preferir `monospace` pixel-art estilo Gameboy o la ya usada en header). |

## Plan de implementacion

1. Añadir en `index.html` dentro de `.splash-snorlax-wrapper`:
   ```html
   <div class="splash-otis" aria-hidden="true"><span class="splash-typed"></span><span class="splash-caret">|</span></div>
   ```
   Y despues de `.splash-zzz-text`:
   ```html
   <div class="splash-pokedex-label" aria-hidden="true"><span class="splash-typed"></span><span class="splash-caret">|</span></div>
   ```
2. CSS inline en `<head>`:
   - `.splash-otis`: `position: absolute; top: -0.5rem; left: 50%; transform: translateX(-50%); font-family: 'Press Start 2P', monospace; font-size: 1.8rem; color: #fff; text-shadow: 2px 2px 0 #000;`
   - `.splash-pokedex-label`: `display: block; margin-top: 0.75rem; font-family: same; font-size: 1.4rem; color: #f5c518; text-align: center;`
   - `.splash-caret`: animacion `blink 1s step-end infinite`.
3. En `main.ts`:
   ```ts
   async function typewrite(el: HTMLElement, text: string, perChar = 180): Promise<void> {
     for (const ch of text) {
       if (typewriterCancelled) { el.textContent = text; return; }
       el.textContent = (el.textContent || "") + ch;
       await sleep(perChar);
     }
   }
   ```
   Lanzar: `await typewrite(otisSpan, "OTIS"); await typewrite(pokedexSpan, "POKéDEX");`
4. Al hacer click en el wrapper: set `typewriterCancelled = true`, rellenar ambos spans inmediatamente, ocultar carets, y proceder con `dismissSplashInteractive` (gsap exit animation ya existente).
5. Asegurar `clearProps` en cualquier tween de gsap que toque los spans para no dejar transforms.
6. Verificar en 360px mobile: escalar fuentes (`font-size: 1.2rem` / `1rem`) y posicion de `OTIS` (no solapar con los ojos/arrows del Snorlax).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual desktop | "OTIS" aparece letra a letra sobre Snorlax, termina en ~0.8s, luego "POKéDEX" aparece bajo las Zzz |
| Manual desktop | Click durante typewriter completa el texto y pasa al dismiss |
| Manual mobile 360px | Textos visibles sin solapar ojos/Zzz, escala adecuada |
| Manual APK Android | Splash JS (no el nativo) muestra ambos textos antes del tap |

## Criterios de aceptacion

- [ ] "OTIS" se escribe progresivamente encima del Snorlax
- [ ] "POKéDEX" se escribe progresivamente debajo de las Zzz, despues de "OTIS"
- [ ] Click en el Snorlax con typewriter en curso completa el texto y dismissa
- [ ] El caret parpadea mientras se escribe y desaparece al terminar cada texto
- [ ] Sin regresion visual de arrows, ojos, animacion breathe del Snorlax
- [ ] Mobile 360px: textos legibles, sin overflow
- [ ] Sin transforms residuales de gsap tras dismiss (usar `clearProps`)

## Notas

- Este splash es el splash JS interactivo (post-API ping). El splash nativo Android (`0095-native-splash-snorlax-icon.md`) es independiente y NO se toca.
- No hacer typewriter via `requestAnimationFrame` complejo; el `await sleep(ms)` simple es suficiente y trivial de cancelar.
- El estado `typewriterCancelled` es una bandera de modulo, resetear al montar splash.
- `POKéDEX` lleva la `é` acentuada exactamente como pide el usuario.
