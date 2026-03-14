# 0022 — Responsive foundation y componentes compartidos

## Descripcion

Establecer una base responsive solida para todo el frontend: crear mixins de breakpoints reutilizables, mejorar la tipografia escalable, y hacer que los componentes compartidos (tabs, filtros, botones, autocomplete, section headers) se adapten correctamente a todos los tamaños de ventana.

## Capas involucradas

- **Frontend (APP)**: archivos SCSS en `frontend/src/styles/` — `_variables.scss`, `_components.scss`, `_tabs.scss`.
- **Core**: No se requieren cambios.
- **Shell**: No se requieren cambios.

## Contexto actual

- Breakpoints existentes: 480px, 768px, 900px, pero definidos inline en cada archivo sin mixins.
- La tipografia usa tamaños fijos (px/rem) sin escalado responsive.
- Los tabs (`_tabs.scss`) tienen ajustes basicos a 480px pero no cubren tablets.
- Los filtros (`_components.scss`) usan flex-wrap pero los selects y pills no escalan bien en mobile.
- Los section headers usan tamaños grandes fijos.
- El autocomplete dropdown no se adapta a pantallas pequeñas.

## Plan de implementacion

### Paso 1 — Crear mixins de breakpoints en _variables.scss

Definir mixins SCSS reutilizables para los breakpoints del proyecto:

```scss
$bp-mobile: 480px;
$bp-tablet: 768px;
$bp-desktop: 1024px;
$bp-wide: 1200px;

@mixin mobile { @media (max-width: $bp-mobile) { @content; } }
@mixin tablet { @media (max-width: $bp-tablet) { @content; } }
@mixin desktop { @media (max-width: $bp-desktop) { @content; } }
@mixin wide { @media (min-width: $bp-wide) { @content; } }
```

- [x] Completado

### Paso 2 — Tipografia responsive

Ajustar los tamaños de fuente principales para que escalen con el viewport:
- Section headers: reducir en tablet y mobile.
- Body text: asegurar legibilidad minima en mobile.
- Usar `clamp()` donde sea apropiado.

- [x] Completado

### Paso 3 — Tabs responsive mejorados

- En tablet: reducir padding y font-size de los tabs.
- En mobile: hacer tabs scrollables horizontalmente si no caben, o reducir a iconos/abreviaciones.
- Asegurar que el indicador animado se ajuste.

- [x] Completado

### Paso 4 — Filtros y pills responsive

- Filter bar: stack vertical en mobile.
- Filter selects: width 100% en mobile.
- Filter pills: tamaño reducido, wrap correcto.
- Botones de accion (clear filters, view toggle): tamaño tactil adecuado (min 44px).

- [x] Completado

### Paso 5 — Section headers y botones responsive

- Section headers: reducir font-size y padding en mobile.
- Botones generales: asegurar min-height tactil, padding adecuado.
- Autocomplete dropdown: ancho 100% en mobile, items mas grandes para touch.

- [x] Completado

## Criterios de aceptacion

- [x] Mixins de breakpoints definidos y usados consistentemente.
- [x] Tipografia legible en todos los tamaños de ventana.
- [x] Tabs navegables y usables en mobile y tablet.
- [x] Filtros y pills correctamente apilados y escalados en mobile.
- [x] Todos los elementos interactivos tienen tamaño tactil adecuado (min 44px).
- [x] No se rompe ningun layout existente en desktop.
- [x] No se introducen cambios en Core ni Shell.

## Dependencias externas nuevas

Ninguna.

## Notas

- Migrar los media queries inline existentes a usar los nuevos mixins.
- Esta tarea es prerequisito de las siguientes (0023, 0024) ya que establece la base.
