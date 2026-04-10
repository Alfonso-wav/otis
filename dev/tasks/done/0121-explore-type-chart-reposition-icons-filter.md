# Task 0121 — Explorar: Tabla de tipos — reposición, iconos con color, filtro arrastrable

## Estado: done

## Goal

Tres mejoras sobre la subpestaña **"Tabla de Tipos"** dentro de Explorar:

1. **Reposicionar la pestaña**: mover "Tabla de Tipos" a la derecha de "Regiones" (nuevo orden: Regiones → Tabla de Tipos → Movimientos → Habilidades → Bayas).

2. **Iconos de tipo con color, sin texto**: los encabezados de fila y columna mostrarán solo el icono SVG del tipo, coloreado con el color de ese tipo. Eliminar el texto abreviado (3 letras).

3. **Filtro de tipos arrastrable**: al hacer clic en un icono de tipo (en el encabezado de fila o de columna), ese tipo se oculta de ambos ejes (fila Y columna desaparecen). El gesto es arrastrable a lo largo de su eje para filtrar múltiples tipos de una pasada. Se añade un botón para restaurar todos los tipos filtrados. La tabla se recompacta visualmente al instante.

## Contexto técnico

### 1. Reposición de la pestaña

**`frontend/src/pages/explore.ts`**:
- Cambiar `TAB_KEYS` de `["typeChart", "regions", "moves", "abilities", "berries"]` a `["regions", "typeChart", "moves", "abilities", "berries"]`.
- Cambiar `activeTab` por defecto de `"typeChart"` a `"regions"` (la primera pestaña del nuevo orden).

### 2. Iconos de tipo con color, sin texto

**`frontend/src/pages/explore/type-chart.ts`**:
- En la función `typeHeader(type)`, eliminar el `<span class="tc-abbr">` con las 3 letras.
- Añadir `data-type="${type}"` al `<th>` del encabezado para que el CSS pueda aplicar el color de fondo.
- Resultado: solo el `<img src="/assets/types/${type}.svg" class="tc-icon" alt="${name}">`.

**`frontend/src/styles/_explore.scss`** (sección type chart):
- Añadir un loop SCSS que genere `background-color` para `.tc-col-header[data-type="fire"]`, `.tc-row-header[data-type="fire"]`, etc., usando `$type-colors` de `_variables.scss`.
- Los iconos SVG de tipo ya están coloreados intrínsecamente, pero el fondo del encabezado les dará contexto visual.
- Ajustar el tamaño del icono si hace falta al quitar el texto (puede crecer un poco, p.ej. 20×20 px).
- Eliminar estilos de `.tc-abbr` (o dejarlos como no-op si se reutilizan en otro sitio — verificar antes de borrar).

**`frontend/src/styles/_dark.scss`**:
- Verificar que las reglas de `.tc-col-header` y `.tc-row-header` en el bloque `[data-bs-theme="dark"]` no sobreescriban el `background-color` del tipo. Si es así, hacer que el color de tipo tenga mayor especificidad o que se aplique con `!important` solo para el hover/filtrado.

### 3. Filtro arrastrable

**`frontend/src/pages/explore/type-chart.ts`** — reescribir `initTypeChart()` para soportar estado de filtro:

#### Estado del módulo
```ts
const filteredTypes = new Set<PokemonType>();
```

#### Renderizado reactivo
- Extraer `renderChart(panel, title)` que genera la tabla solo con los tipos NO filtrados.
- Llamar a `renderChart()` al inicializar y cada vez que cambie `filteredTypes`.
- Tras cada re-render, re-adjuntar los event listeners de filtrado.

#### Lógica de filtrado al hacer clic/arrastrar
- En cada `<th class="tc-col-header" data-type="...">` y `<th class="tc-row-header" data-type="...">`:
  - `mousedown`: marcar el tipo como filtrado, activar modo drag (`isDragging = true`), registrar la dirección del eje (col = horizontal, row = vertical).
  - `mousemove` / `mouseenter` (mientras `isDragging`): filtrar los tipos sobre los que pase el cursor si están en el mismo eje.
  - `mouseup` / `mouseleave` del documento: desactivar modo drag.
- Al filtrar un tipo, se elimina tanto su columna como su fila del render (tabla sin esa fila/columna).

#### Comportamiento de compactación visual
- El re-render elimina las filas y columnas del DOM directamente, no las oculta con CSS.
- Usar una transición CSS breve (fade/scale) en el contenedor para suavizar el cambio de tamaño.
- No se necesita animación frame-by-frame; basta con un `transition: opacity 0.15s` y actualizar el innerHTML.

#### Botón "Restaurar tipos"
- Solo visible cuando `filteredTypes.size > 0`.
- Texto i18n: `typeChart.restoreTypes` (ES: "Restaurar tipos", EN: "Restore types").
- Al hacer clic: `filteredTypes.clear()`, re-render.
- Posición: justo encima de la tabla (zona `.type-chart-labels` o un nuevo `div.tc-filter-bar`).
- Añadir a los locales la clave `typeChart.restoreTypes` en `es.json` y `en.json`.

#### Touch / mobile
- Añadir equivalentes `touchstart`, `touchmove`, `touchend` para soporte en móvil.
- En `touchmove`, usar `document.elementFromPoint()` para detectar el tipo sobre el que se arrastra.

### Archivos afectados

- `frontend/src/pages/explore.ts` — reordenar `TAB_KEYS`, cambiar `activeTab` por defecto
- `frontend/src/pages/explore/type-chart.ts` — quitar abreviaturas, añadir `data-type`, lógica de filtrado
- `frontend/src/styles/_explore.scss` — colores de tipo en encabezados, ajuste tamaño icono, quitar `.tc-abbr` si procede
- `frontend/src/styles/_dark.scss` — verificar que colores de tipo no queden sobrescritos en modo oscuro
- `frontend/src/locales/es.json` — añadir `typeChart.restoreTypes`
- `frontend/src/locales/en.json` — añadir `typeChart.restoreTypes`

## Acceptance criteria

- [ ] La pestaña "Tabla de Tipos" aparece en el orden: Regiones → Tabla de Tipos → Movimientos → Habilidades → Bayas.
- [ ] Al abrir Explorar, la pestaña activa por defecto es "Regiones".
- [ ] Los encabezados de fila y columna muestran solo el icono del tipo, sin texto.
- [ ] Cada encabezado tiene el color de fondo del tipo correspondiente.
- [ ] Hacer clic en un icono de tipo oculta la fila y columna de ese tipo; la tabla se compacta.
- [ ] Arrastrar el cursor/dedo a lo largo del eje del encabezado filtra los tipos por los que pasa.
- [ ] Al ocultar tipos, el botón "Restaurar tipos" aparece.
- [ ] Al pulsar "Restaurar tipos", la tabla vuelve al estado original con los 18 tipos.
- [ ] Las etiquetas se actualizan al cambiar el idioma (locale-changed).
- [ ] No hay errores en consola.
- [ ] La tabla sigue siendo usable en móvil (touch drag funciona).

## Dependencias

Ninguna. Tarea independiente sobre código ya existente en `main`.
