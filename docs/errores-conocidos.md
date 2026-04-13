# Errores Conocidos y Lecciones Aprendidas

Documento generado a partir del historial de 153 tareas del proyecto Otis.
Contiene los patrones de error recurrentes detectados durante el desarrollo, con ejemplos concretos y reglas para evitarlos.

---

## 1. Race Conditions en Handlers Async

**Frecuencia:** 5+ casos (tasks 0056, 0107, 0108, 0132, 0062)

### Problema

Handlers `async` sin exclusion mutua se ejecutan concurrentemente cuando el usuario hace clicks rapidos o secuenciales. Esto causa estados inconsistentes, overlays bloqueantes que nunca se ocultan, o datos corruptos.

### Casos concretos

| Task | Descripcion |
|------|-------------|
| 0056 | Delete de equipos en secuencia rapida: handlers concurrentes desincronizaban DOM y event listeners |
| 0107 | Clicks rapidos en filtros Legendary/Mythical dejaban `sortingOverlay` visible (z-index 9999), bloqueando toda interaccion |
| 0132 | Keydown handler se ejecutaba antes que el handler de autocomplete, enviando texto parcial al backend |

### Regla de prevencion

```typescript
// SIEMPRE usar guard flag en handlers async que modifican estado compartido
let isProcessing = false;

async function handleAction() {
  if (isProcessing) return;
  isProcessing = true;
  try {
    // ... operacion async
  } finally {
    isProcessing = false;  // SIEMPRE en finally
  }
}
```

- Todo handler async que muestre/oculte overlays DEBE usar `try/finally` para garantizar el cleanup.
- Si hay dependencia de orden entre event listeners (ej. autocomplete vs keydown), verificar el orden de binding y usar `e.defaultPrevented`.

---

## 2. Re-render Destruye Estado del Usuario

**Frecuencia:** 3+ casos (tasks 0075, 0103, 0044)

### Problema

Al re-renderizar secciones completas con `innerHTML`, se pierden valores que el usuario habia introducido (inputs, scroll position, seleccion activa) porque el HTML nuevo tiene valores hardcodeados.

### Casos concretos

| Task | Descripcion |
|------|-------------|
| 0075 | Input de batch simulation se reseteaba a 100 tras cada ejecucion porque el template tenia `value="100"` hardcodeado |
| 0103 | Grid-to-table switch mostraba datos desde offset anterior porque `lastRenderedItems` contenia cache stale |
| 0044 | Boton de batch desaparecia tras reset porque el reset llamaba a `startBattle()` en vez de volver a idle |

### Regla de prevencion

- **Nunca hardcodear valores en templates HTML**. Mantener estado en variables de modulo y usarlas en el template.
- Al cambiar de vista, **resetear explicitamente** caches y offsets (`lastRenderedItems = []`, `offset = 0`).
- Antes de re-render, guardar el estado del usuario que deba persistir.

```typescript
// MAL
`<input type="number" value="100">`

// BIEN
let batchCount = 100;
`<input type="number" value="${batchCount}">`
// Y actualizar batchCount en el handler ANTES de re-render
```

---

## 3. GSAP y Estilos Inline Residuales

**Frecuencia:** 3+ casos (tasks 0108, 0109, 0137)

### Problema

Las animaciones GSAP dejan estilos `transform` y `opacity` inline en los elementos. Estos crean stacking contexts inesperados que rompen la interaccion de elementos con `position: absolute` (dropdowns, overlays).

### Casos concretos

| Task | Descripcion |
|------|-------------|
| 0108 | Dropdowns de filtro dejaban de funcionar tras navegar al detalle y volver. GSAP dejaba `transform` inline que creaba stacking context bloqueando clicks |
| 0109 | Imagen fantasma: el contenido anterior era visible durante la animacion de entrada porque el spinner se asignaba despues de `await showView()` |
| 0137 | Pantalla negra al iniciar: splash dependia de assets procesados por Vite, no disponibles hasta que el bundle cargara |

### Regla de prevencion

- **Siempre** usar `clearProps` en tweens GSAP que afecten `transform` u `opacity`:

```typescript
gsap.fromTo(el, { opacity: 0 }, {
  opacity: 1,
  clearProps: "opacity,transform"  // OBLIGATORIO
});
```

- Antes de iniciar una animacion de transicion, **limpiar el contenido anterior** (ej. poner spinner) para evitar ghost images.
- Assets criticos para splash/primera carga: **inline en HTML** (base64, SVG inline, CSS inline en `<head>`). No depender del bundle de Vite.

---

## 4. Rutas de Assets: Dev vs Build vs APK

**Frecuencia:** 3+ casos (tasks 0093, 0137, 0062)

### Problema

Assets referenciados con rutas `/src/assets/...` funcionan en desarrollo (Vite los resuelve) pero no existen en el build de produccion ni en el APK porque Vite no los copia al `dist/`.

### Casos concretos

| Task | Descripcion |
|------|-------------|
| 0093 | 18 iconos SVG de tipos no cargaban en APK: estaban en `src/assets/types/` referenciados como `/src/assets/types/`. Vite no los copio al dist |
| 0137 | Snorlax SVG del splash tardaba en cargar porque dependia del procesamiento del bundle |
| 0062 | Android bloqueaba peticiones HTTP cleartext por defecto (Android 9+) |

### Regla de prevencion

- Assets estaticos que se referencian por URL en runtime **van en `public/`**, no en `src/assets/`.
- Rutas siempre sin prefijo `/src/`: usar `/assets/tipos/fire.svg`, nunca `/src/assets/tipos/fire.svg`.
- **Probar siempre en build de produccion** (`wails build` / APK) antes de dar por terminada una tarea con assets nuevos.
- Para Android: configurar `network_security_config.xml` con `cleartextTrafficPermitted` para `10.0.2.2` y `localhost`.

---

## 5. Fallback Chains con Loops Infinitos

**Frecuencia:** 2 casos (tasks 0037, 0039)

### Problema

En handlers `onerror` de imagenes, comparar `this.src` (URL absoluta resuelta por el navegador) con una ruta relativa de fallback siempre falla, causando un loop infinito entre el src original y el fallback.

### Caso concreto

| Task | Descripcion |
|------|-------------|
| 0037 | `this.src` devuelve `http://wails.localhost/assets/sprites/...` pero `fallbackSrc` es `/assets/sprites/...` — comparacion siempre false, loop infinito |

### Regla de prevencion

```typescript
// NUNCA comparar this.src con rutas relativas
// USAR data attributes como contador de intentos
img.dataset.fallbackAttempt = "0";
img.onerror = function() {
  const attempt = parseInt(this.dataset.fallbackAttempt || "0");
  if (attempt >= MAX_FALLBACKS) {
    this.src = PLACEHOLDER;
    this.onerror = null;  // Cortar el chain
    return;
  }
  this.dataset.fallbackAttempt = String(attempt + 1);
  this.src = FALLBACK_CHAIN[attempt];
};
```

- Siempre incluir un **fallback terminal** (placeholder) y `this.onerror = null` para cortar el chain.
- Nunca comparar URLs absolutas con relativas.

---

## 6. `<select>` y el Primer Option: Change Event No Dispara

**Frecuencia:** 1 caso (task 0151)

### Problema

El navegador auto-selecciona la primera `<option>` de un `<select>`. Si el usuario clickea esa misma opcion, el evento `change` no se dispara porque el valor no cambio.

### Caso concreto

| Task | Descripcion |
|------|-------------|
| 0151 | No se podia cargar el primer equipo como companion porque ya estaba "seleccionado" por defecto |

### Regla de prevencion

- **Siempre** agregar un placeholder disabled como primera opcion:

```html
<select>
  <option value="" disabled selected>-- Seleccionar --</option>
  <option value="team1">Team 1</option>
</select>
```

---

## 7. CSS: Position Absolute vs Fixed en Containers Largos

**Frecuencia:** 2 casos (tasks 0094, 0145)

### Problema

Overlays/spinners centrados con `position: absolute` dentro de containers mas altos que el viewport quedan fuera de la zona visible. El usuario no los ve y cree que la app se colgo.

### Casos concretos

| Task | Descripcion |
|------|-------------|
| 0094 | Mr. Mime sorting overlay invisible en tablas largas: centrado `absolute` dentro de container de 5000px de alto |
| 0145 | Diglett clipped: sprite no centrado horizontalmente en container |

### Regla de prevencion

- Overlays/loading que deben ser siempre visibles: usar `position: fixed` relativo al viewport.
- Sprites en containers con clip: centrar explicitamente con `left: 50%; transform: translateX(-50%)`.
- Probar con contenido largo (200+ filas, scroll activo) antes de dar por bueno un overlay.

---

## 8. Backend No Idempotente

**Frecuencia:** 2 casos (tasks 0055, 0150)

### Problema

Operaciones de backend que fallan cuando se llaman en un estado ya alcanzado (ej. borrar archivo que no existe, endpoint no registrado para plataforma movil).

### Casos concretos

| Task | Descripcion |
|------|-------------|
| 0055 | `DeleteTeam` fallaba con error si el archivo ya no existia. Frontend se quedaba inconsistente |
| 0150 | Endpoint REST de distribucion de tipos faltaba en `handlers.go` (existia en Wails IPC, no en HTTP) |

### Regla de prevencion

- Operaciones de delete: **ignorar `os.ErrNotExist`**. Un delete de algo inexistente es exito, no error.
- **Al agregar un binding Wails, agregar tambien el handler REST** en `app/mobile/handlers.go`. Mantener paridad IPC/HTTP.
- El frontend SIEMPRE debe refrescar UI en `finally`, no solo en el happy path del `try`.

---

## 9. i18n: Traducciones Incompletas y Recarga

**Frecuencia:** 2 casos (tasks 0113, 0117)

### Problema

Textos en ingles aparecen en la version en espanol porque (a) faltan keys en `es.json`, o (b) la vista activa no se re-renderiza cuando el usuario cambia de idioma.

### Casos concretos

| Task | Descripcion |
|------|-------------|
| 0113 | "Category", "Habitat", "Shape" en ingles dentro de la vista detalle en espanol. Ademas, cambiar idioma no actualizaba la vista detalle |
| 0117 | Flavor text siempre en ingles: backend no filtraba por idioma. Stats labels sin traduccion |

### Regla de prevencion

- Al agregar texto visible al usuario, **agregar la key en AMBOS** `en.json` y `es.json` inmediatamente.
- Toda vista que muestre texto traducido **debe escuchar el evento `locale-changed`** y re-renderizarse.
- Si el backend devuelve texto localizable (flavor text, descripciones), **aceptar parametro `lang`** y filtrar en backend.

---

## 10. ECharts Init con Container Sin Dimensiones

**Frecuencia:** 1 caso (task 0144)

### Problema

`echarts.init(container)` falla silenciosamente si el container tiene dimensiones 0x0 en el momento de la llamada (ej. container oculto, animacion en progreso, tab no visible).

### Caso concreto

| Task | Descripcion |
|------|-------------|
| 0144 | Donut chart de regiones no renderizaba porque el container no tenia dimensiones calculadas al momento del init |

### Regla de prevencion

```typescript
// Diferir init hasta que el container tenga dimensiones
requestAnimationFrame(() => {
  const chart = echarts.init(container);
  chart.setOption(options);
});
```

- Tras expandir/animar un contenedor padre, llamar `chart.resize()`.
- Si el chart esta en un tab oculto, inicializar cuando el tab se haga visible.

---

## 11. Web Scraping: Selector CSS Fragil

**Frecuencia:** 1 caso (task 0038)

### Problema

Scrapers que usan selectores genericos (`Find("img").First()`) seleccionan el elemento equivocado cuando el HTML tiene multiples elementos similares.

### Caso concreto

| Task | Descripcion |
|------|-------------|
| 0038 | Scraper de sprites seleccionaba el sprite B&W (primer `img`) en vez del color (segundo `span > img`) |

### Regla de prevencion

- Usar selectores CSS lo mas especificos posible: `span:nth-child(2) a img` en vez de `img`.
- Incluir HTML de ejemplo en los tests para validar el selector contra la estructura real.
- Los scrapers son fragiles por naturaleza: documentar la estructura HTML esperada junto al selector.

---

## 12. Mobile: Touch Targets y Overflow

**Frecuencia:** 2 casos (tasks 0064, 0101)

### Problema

Layouts con `flex` sin `wrap` o anchos minimos demasiado grandes causan overflow horizontal en pantallas de 360px. Botones demasiado pequenos son dificiles de pulsar en tactil.

### Casos concretos

| Task | Descripcion |
|------|-------------|
| 0064 | Botones de equipo (Atk, Def, delete) cortados en movil: flex sin wrap + min-width 90px en 360px |
| 0101 | Headers de tabla ilegibles en movil: sin min-width, columnas comprimidas hasta ser ilegibles |

### Regla de prevencion

- Touch targets: **minimo 44px x 44px** (directriz Apple/Google).
- Usar `flex-wrap: wrap` por defecto en layouts con multiples items.
- Tablas en movil: `overflow-x: auto` en el container, `min-width` en columnas criticas.
- **Probar siempre en viewport 360px** antes de cerrar una tarea con componentes nuevos.

---

## Resumen de Patrones por Frecuencia

| # | Patron | Casos | Severidad |
|---|--------|-------|-----------|
| 1 | Race conditions en handlers async | 5+ | Alta |
| 2 | Re-render destruye estado del usuario | 3+ | Media |
| 3 | GSAP estilos residuales | 3+ | Alta |
| 4 | Rutas assets dev vs build | 3+ | Alta |
| 5 | Fallback image loops | 2 | Alta |
| 6 | Select primer option sin change | 1 | Baja |
| 7 | Position absolute en containers largos | 2 | Media |
| 8 | Backend no idempotente / paridad IPC-HTTP | 2 | Alta |
| 9 | i18n incompleto / sin recarga | 2 | Media |
| 10 | ECharts init sin dimensiones | 1 | Media |
| 11 | Scraping selector fragil | 1 | Baja |
| 12 | Mobile overflow y touch targets | 2 | Media |

---

## Checklist Pre-Merge

Antes de dar una tarea como completada, verificar:

- [ ] Handlers async tienen guard flags y `try/finally`
- [ ] Templates HTML usan variables de estado, no valores hardcodeados
- [ ] Animaciones GSAP tienen `clearProps`
- [ ] Assets nuevos estan en `public/`, no en `src/assets/`
- [ ] Bindings Wails tienen su handler REST equivalente en `handlers.go`
- [ ] Textos nuevos estan en `en.json` Y `es.json`
- [ ] Vistas con i18n escuchan `locale-changed`
- [ ] Probado en build de produccion (no solo dev)
- [ ] Probado en viewport 360px si hay componentes nuevos
- [ ] Overlays usan `position: fixed` si deben ser siempre visibles
- [ ] `<select>` tienen placeholder disabled
- [ ] Fallback chains de imagenes tienen limite y terminador
