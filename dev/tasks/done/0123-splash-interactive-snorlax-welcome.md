# Task 0123 — Pantalla de bienvenida interactiva con Snorlax

## Estado: pending

## Goal

Añadir una **pantalla de bienvenida interactiva** antes de entrar a la aplicación. El flujo completo es:

1. App arranca → pantalla de carga existente (Snorlax dormido + ZZZ flotantes mientras la API inicializa).
2. API lista → splash transiciona al estado "interactivo": los ZZZ desaparecen, aparece un texto de "toca a Snorlax" y el cursor cambia a pointer sobre el Snorlax.
3. Usuario hace click en el Snorlax → Snorlax se agita (shake GSAP) + feedback visual.
4. Tras el shake → transición de salida con **degradado** (overlay radial o sweep que cubre toda la pantalla) y la app aparece.

### Diseño del estado interactivo (libertad creativa)

- Los ZZZ dejan de flotar (se desvanecen o paran).
- El Snorlax mantiene su respiración suave (`snorlax-breathe`).
- Aparece un texto sutil bajo el Snorlax: *"¡Toca a Snorlax!"* / *"Tap Snorlax!"* (i18n).
- Cursor `pointer` sobre `.splash-snorlax`.
- Fondo puede tener un brillo/glow suave pulsando detrás del Snorlax para invitar al click.

### Animaciones de click

- **Shake**: `gsap.to(".splash-snorlax", { keyframes: [{rotation: -8}, {rotation: 8}, {rotation: -5}, {rotation: 5}, {rotation: 0}], duration: 0.5, ease: "power2.out" })`
- **Scale pop**: breve scale-up al inicio del shake (`scale: 1.15` → vuelve a 1) para que "reaccione".
- **Transición de salida**: un `div.splash-overlay-out` (posición fija, inset 0, fondo `#1a202c`) que aparece con un clip-path o radial-gradient que se expande desde el centro hacia los bordes, luego el overlay completo hace `opacity: 0` para revelar la app. Duración total ~0.8s.

---

## Contexto técnico

### Estado actual

El splash screen (`#splash-screen`) en `frontend/index.html` se auto-descarta via `dismissSplash()` en `main.ts` en cuanto la API responde. No hay interacción de usuario.

```ts
// main.ts — actual
function dismissSplash(): void {
  const splash = document.getElementById("splash-screen");
  if (!splash) return;
  splash.classList.add("splash-fade-out");
  splash.addEventListener("transitionend", () => splash.remove(), { once: true });
}

ListGenerations()
  .then(() => dismissSplash())
  .catch(() => dismissSplash());
```

### Cambios en `main.ts`

Dividir en dos fases:

```ts
function showSplashInteractive(): void {
  const splash = document.getElementById("splash-screen");
  const zzz = splash?.querySelector(".splash-zzz");
  const snorlax = splash?.querySelector(".splash-snorlax") as HTMLElement | null;
  const hint = splash?.querySelector(".splash-hint") as HTMLElement | null;
  if (!splash || !snorlax) return;

  // Fase 1: apagar ZZZ y mostrar hint
  gsap.to(zzz, { opacity: 0, duration: 0.4 });
  gsap.to(hint, { opacity: 1, duration: 0.4, delay: 0.3 });
  snorlax.style.cursor = "pointer";

  // Fase 2: click en Snorlax
  snorlax.addEventListener("click", () => dismissSplashInteractive(splash, snorlax), { once: true });
}

function dismissSplashInteractive(splash: HTMLElement, snorlax: HTMLElement): void {
  snorlax.style.pointerEvents = "none";

  // Shake + pop
  gsap.timeline()
    .to(snorlax, { scale: 1.15, duration: 0.1 })
    .to(snorlax, {
      keyframes: [
        { rotation: -8 }, { rotation: 8 },
        { rotation: -5 }, { rotation: 5 },
        { rotation: 0 }
      ],
      duration: 0.45,
      ease: "power2.out"
    }, "<0.05")
    .to(snorlax, { scale: 1, duration: 0.1 })
    .call(() => {
      // Crear y animar overlay de salida
      const overlay = document.createElement("div");
      overlay.className = "splash-exit-overlay";
      splash.appendChild(overlay);
      gsap.fromTo(overlay,
        { clipPath: "circle(0% at 50% 50%)" },
        {
          clipPath: "circle(150% at 50% 50%)",
          duration: 0.55,
          ease: "power2.in",
          onComplete: () => {
            splash.remove();
          }
        }
      );
    });
}

ListGenerations()
  .then(() => showSplashInteractive())
  .catch(() => showSplashInteractive());
```

### Cambios en `frontend/index.html`

Añadir dentro de `#splash-screen` el elemento hint (oculto inicialmente):

```html
<p class="splash-hint" data-i18n="splash.hint" style="opacity:0">¡Toca a Snorlax!</p>
```

### Cambios en `frontend/src/styles/_splash.scss`

- Añadir `.splash-hint` (texto de invitación, similar a `.splash-text` pero con color más vivo y cursor pointer hint).
- Añadir `.splash-exit-overlay` (position fixed, inset 0, background `#1a202c`, z-index alto, sin pointer events tras crear).
- Añadir `@keyframes splash-glow` para el fondo pulsante detrás del Snorlax (optional: `box-shadow` o `radial-gradient` que pulsa en el `#splash-screen`).

```scss
.splash-hint {
  margin-top: 1rem;
  color: #e2e8f0;
  font-family: "Exo 2", sans-serif;
  font-size: 0.9rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  opacity: 0; // GSAP lo anima a 1
}

.splash-exit-overlay {
  position: fixed;
  inset: 0;
  z-index: 99998;
  background: #1a202c;
  pointer-events: none;
}
```

### Localización

**`frontend/src/locales/es.json`** y **`en.json`** — añadir en `"splash"`:

```json
"hint": "¡Toca a Snorlax!"   // es.json
"hint": "Tap Snorlax!"        // en.json
```

---

## Archivos afectados

### Frontend
- `frontend/index.html` — añadir `.splash-hint` dentro de `#splash-screen`
- `frontend/src/main.ts` — refactorizar `dismissSplash()` en `showSplashInteractive()` + `dismissSplashInteractive()`
- `frontend/src/styles/_splash.scss` — añadir `.splash-hint`, `.splash-exit-overlay` y animación de glow
- `frontend/src/locales/es.json` — añadir `splash.hint`
- `frontend/src/locales/en.json` — añadir `splash.hint`

### Backend
Ninguno.

---

## Acceptance criteria

- [ ] Al arrancar la app, el splash existente (Snorlax + ZZZ + texto cargando) se muestra mientras la API inicializa.
- [ ] Cuando la API responde, los ZZZ desaparecen y aparece el texto de invitación (i18n).
- [ ] El cursor cambia a `pointer` al pasar sobre el Snorlax en estado interactivo.
- [ ] Al hacer click en el Snorlax, se reproduce la animación de shake + pop.
- [ ] Tras el shake, un overlay de transición con degradado radial cubre la pantalla.
- [ ] El overlay desaparece revelando la aplicación (el `#splash-screen` se elimina del DOM).
- [ ] El texto "¡Toca a Snorlax!" / "Tap Snorlax!" respeta el idioma activo en el momento del arranque.
- [ ] No hay errores en consola del frontend.
- [ ] El comportamiento de fallback (API error) también lleva al estado interactivo (el usuario puede entrar igualmente).

---

## Dependencias

Ninguna. Tarea independiente sobre el splash screen existente.
