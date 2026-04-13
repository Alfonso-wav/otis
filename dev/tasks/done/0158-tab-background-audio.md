# Musica de fondo por pestaña (intro + pokedex + explore + simulaciones)

**ID**: 0158-tab-background-audio
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

Añadir un sistema de musica de fondo en bucle con cambios por pestaña y fundidos entre tracks. Cuatro mp3 ya existen en `C:\Users\alfon\Documents\otis_pkm_audios`:

- `Intro_Pokemon Red_Blue Opening [JvKoBnv96PM].mp3` → se reproduce al abrir la app (mientras se ve el splash interactivo del Snorlax).
- `Pokedex_Pokemon BlueRed - Cinnabar Island.mp3` → al pulsar el Snorlax y entrar a la pestaña `pokedex`.
- `Explorar_Pokemon BlueRed - Route 1.mp3` → al entrar a la pestaña `explore`.
- `Simulaciones_Pokemon BlueRed - Gym Leader Battle!.mp3` → al entrar a la pestaña `builds` (que se muestra como "Simulaciones" en ES).

Requisitos:
1. Todos los tracks en loop hasta que se cambie de pestaña.
2. Al cambiar de pestaña, fundido de salida (~600ms) del track actual + fundido de entrada (~600ms) del nuevo.
3. La pestaña `settings` NO tiene track propio → mantiene reproduciendose el track de la pestaña anterior sin interrumpirlo.
4. En la pagina de Settings añadir control de volumen (slider 0-100) y toggle mute/activo. Persistir en `localStorage`.

## Capas afectadas

- **Core**: ninguna.
- **Shell**: ninguna (los mp3 son assets estaticos).
- **APP**: frontend completo (nuevo modulo audio + integracion con router + UI en settings + assets publicos).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/public/audio/intro.mp3` | crear | Copiar desde `C:\Users\alfon\Documents\otis_pkm_audios\Intro_*.mp3` |
| `frontend/public/audio/pokedex.mp3` | crear | Copiar desde `Pokedex_*.mp3` |
| `frontend/public/audio/explore.mp3` | crear | Copiar desde `Explorar_*.mp3` |
| `frontend/public/audio/simulations.mp3` | crear | Copiar desde `Simulaciones_*.mp3` |
| `frontend/src/audio.ts` | crear | Modulo de audio: 2 `HTMLAudioElement` pool (actual / siguiente) para crossfade, API `playTrack(id)`, `stop()`, `setVolume(v)`, `setMuted(b)`. Lectura inicial de localStorage. |
| `frontend/src/main.ts` | modificar | Al arrancar, `playTrack("intro")` antes/durante el splash. Cuando el splash dismiss lleva a pokedex, llamar `playTrack("pokedex")`. |
| `frontend/src/router.ts` | modificar | En `navigate(id)`: si `id` ∈ {pokedex, explore, builds} → `playTrack(id)`. Si `id === "settings"` → NO cambiar de track (dejar reproduciendo el anterior). |
| `frontend/src/settings.ts` | modificar | Añadir UI de volumen (slider) y toggle mute. Al cambiar, `setVolume` / `setMuted` inmediatos + persistir. Incluir dentro del flujo de "Apply" o aplicar en vivo (preferible vivo, sin confirmar). |
| `frontend/index.html` | modificar | Añadir seccion en settings con slider de volumen y checkbox de mute. Etiquetas i18n. |
| `frontend/src/locales/en.json` | modificar | Claves: `settings.musicSection`, `settings.musicVolume`, `settings.musicMuted`. |
| `frontend/src/locales/es.json` | modificar | Mismas claves en español. |
| `.gitignore` / git add | verificar | Los mp3 deben commitearse como assets; verificar tamaños (pueden ser pesados). Si >2MB cada uno, considerar compresion. |

## Plan de implementacion

1. **Assets**: copiar los 4 mp3 a `frontend/public/audio/` con nombres normalizados (`intro.mp3`, `pokedex.mp3`, `explore.mp3`, `simulations.mp3`). Recordar que Vite sirve `public/` en la raiz, asi la URL final sera `/audio/intro.mp3`.
2. **Modulo `audio.ts`**:
   ```ts
   type TrackId = "intro" | "pokedex" | "explore" | "simulations";
   const SRC: Record<TrackId, string> = {
     intro: "/audio/intro.mp3",
     pokedex: "/audio/pokedex.mp3",
     explore: "/audio/explore.mp3",
     simulations: "/audio/simulations.mp3",
   };
   const FADE_MS = 600;
   let currentAudio: HTMLAudioElement | null = null;
   let currentTrack: TrackId | null = null;
   let userVolume = readVolume(); // 0..1
   let muted = readMuted();

   export async function playTrack(id: TrackId): Promise<void> {
     if (currentTrack === id && currentAudio && !currentAudio.paused) return;
     const next = new Audio(SRC[id]);
     next.loop = true;
     next.volume = 0;
     await next.play().catch(() => {}); // autoplay policy — may require user gesture
     fadeIn(next, muted ? 0 : userVolume, FADE_MS);
     if (currentAudio) fadeOutAndStop(currentAudio, FADE_MS);
     currentAudio = next;
     currentTrack = id;
   }

   export function setVolume(v: number): void { userVolume = clamp(v, 0, 1); if (currentAudio && !muted) currentAudio.volume = userVolume; persist(); }
   export function setMuted(b: boolean): void { muted = b; if (currentAudio) currentAudio.volume = b ? 0 : userVolume; persist(); }
   ```
3. **Tab mapping** en `router.ts`:
   ```ts
   const TAB_TRACK: Record<string, TrackId> = { pokedex: "pokedex", explore: "explore", builds: "simulations" };
   // dentro de navigate(id):
   const track = TAB_TRACK[id];
   if (track) playTrack(track);
   // id === "settings" → no-op para audio
   ```
4. **Intro**: el primer `playTrack("intro")` va en `main.ts` justo antes del splash. Problema: autoplay policy de navegadores — el audio solo se reproduce tras user gesture. Solucion: intentar `play()` tras el primer click del usuario en el Snorlax; si falla silencioso, arrancar `pokedex.mp3` en su lugar. Alternativamente: el splash requiere click para dismissar → justo despues del click hacemos `playTrack("intro")` y al terminar el dismiss hacemos `playTrack("pokedex")`. **Decision preferida**: arrancar intro dentro del mismo handler de click del splash (cuenta como user gesture) y transicionar a `pokedex` al montar la tab destino.
5. **Settings UI**: nueva seccion en `#tab-settings` con:
   ```html
   <section>
     <h3 data-i18n="settings.musicSection">Musica</h3>
     <label>
       <input type="checkbox" id="settings-music-muted" />
       <span data-i18n="settings.musicMuted">Silenciar musica</span>
     </label>
     <label>
       <span data-i18n="settings.musicVolume">Volumen</span>
       <input type="range" id="settings-music-volume" min="0" max="100" step="1" />
     </label>
   </section>
   ```
   Listeners: `input` en slider → `setVolume(v/100)`. `change` en checkbox → `setMuted(b)`. Aplicar en vivo (no esperar a "Apply").
6. **Persistencia**: `localStorage` keys `music-volume` (string "0".."1") y `music-muted` (string "true"/"false"). Leer al inicializar el modulo.
7. **Fade helpers** con GSAP (ya en el proyecto) o `requestAnimationFrame`:
   ```ts
   function fadeTo(el: HTMLAudioElement, target: number, ms: number, onDone?: () => void) {
     const start = el.volume;
     const t0 = performance.now();
     (function step(now: number) {
       const p = Math.min(1, (now - t0) / ms);
       el.volume = start + (target - start) * p;
       if (p < 1) requestAnimationFrame(step);
       else onDone?.();
     })(t0);
   }
   ```
8. **Cleanup**: al cambiar de pagina, liberar `currentAudio` previo (`pause()`, `src = ""`) tras el fade-out para no mantener buffers.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual desktop | Al cargar la app y clicar Snorlax suena `intro.mp3`, transiciona a `pokedex.mp3` al entrar a la pokedex |
| Manual desktop | Click en tab Explorar → fundido al track de explore |
| Manual desktop | Click en tab Simulaciones → fundido al track de simulations |
| Manual desktop | Click en tab Settings → el track actual sigue reproduciendose sin cambio |
| Manual desktop | Slider volumen cambia volumen en vivo |
| Manual desktop | Toggle mute silencia y desilencia sin reiniciar el track |
| Manual recarga | Volumen y mute persisten tras recargar |
| Manual APK Android | Audios suenan desde `/audio/*.mp3` en el build (verificar que `public/` se copia) |
| Manual APK Android | Al minimizar/restaurar la app el audio pausa/reanuda (comportamiento nativo WebView) |

## Criterios de aceptacion

- [ ] Los 4 mp3 existen en `frontend/public/audio/` y se cargan en dev y en build
- [ ] `intro.mp3` suena al primer gesto del usuario (click en splash Snorlax)
- [ ] Cada una de las 3 pestañas reproduce su track en loop tras la transicion
- [ ] Al cambiar de pestaña hay fundido de salida + fundido de entrada (~600ms)
- [ ] Settings no cambia de track; mantiene el anterior sonando
- [ ] Slider de volumen aplica en vivo
- [ ] Toggle mute aplica en vivo y persiste
- [ ] Funciona en APK Android (assets `public/` incluidos)
- [ ] No hay dos tracks sonando a la vez tras el crossfade

## Notas

- Autoplay policy: en navegadores desktop y Android WebView, `play()` programatico sin gesture falla. El primer `playTrack` debe dispararse desde un handler de evento (click del splash).
- Nombres de archivo: renombrar a `intro.mp3`/`pokedex.mp3`/`explore.mp3`/`simulations.mp3` para evitar rutas con espacios y caracteres raros.
- `public/` es la ruta correcta (ver CLAUDE.md seccion Assets estaticos). NO poner los mp3 en `src/assets/`.
- Consideracion de tamaño: los mp3 de intros suelen ser 2-4 MB. Si el APK crece demasiado, bitrate 128kbps mono es suficiente para musica de fondo. Documentar en `docs/errores-conocidos.md` si se recomprimen.
- Tab interna `builds` se muestra como "Simulaciones" en ES (ver `frontend/src/locales/es.json:15`). El mapeo de audio usa el id tecnico `builds` → `simulations.mp3`.
