# Task 0126 — Type Chart: iconos mucho más grandes

## Estado: done

## Goal

Los iconos de tipo en la tabla de efectividades (Explorar > Tabla de Tipos) son demasiado pequeños. Aumentar significativamente el tamaño de los iconos SVG dentro de las celdas de cabecera para que sean más visibles y reconocibles de un vistazo.

---

## Contexto técnico

### Archivos principales

- `frontend/src/styles/_explore.scss` — estilos del type chart (`.tc-icon`, `.tc-col-header`, `.tc-row-header`, `.tc-cell`)
- `frontend/src/pages/explore/type-chart.ts` — genera el HTML de la tabla (solo lectura, no requiere cambios)

### Estado actual

Los iconos están definidos en `.tc-icon` con `60px × 60px` y un padding de `8px` por el fondo circular de color. Las celdas de cabecera (`.tc-col-header`) tienen `108px × auto` y las de fila (`.tc-row-header`) `156px` de ancho.

Las celdas de datos (`.tc-cell`) son `108px × 78px`, lo cual deja espacio para iconos más grandes en las cabeceras sin romper la proporción de la tabla.

### Propuesta de cambio

Aumentar `.tc-icon` de `60px` a **84px** (40% más grande). Ajustar padding del fondo circular a `10px`. Aumentar `.tc-col-header` a `120px` y la altura de `.tc-cell` a `84px` para acomodar los iconos proporcionalmente. Ajustar `.tc-row-header` a `168px`.

**Nota**: los breakpoints responsive de la tarea 0125 (media queries) deberán re-calibrarse proporcionalmente si ya están implementados.

---

## Cambios requeridos

### 1. Aumentar tamaño de iconos — `_explore.scss`

```scss
.tc-icon {
  width: 84px;   // antes: 60px
  height: 84px;  // antes: 60px
  display: block;
  margin: 0 auto;
}
```

### 2. Ajustar padding del fondo circular

```scss
@each $type, $color in $type-colors {
  .tc-col-header[data-type="#{$type}"] .tc-icon,
  .tc-row-header[data-type="#{$type}"] .tc-icon {
    background-color: $color;
    border-radius: 50%;
    padding: 10px;   // antes: 8px
    box-sizing: border-box;
  }
}
```

### 3. Ajustar dimensiones de celdas

```scss
.tc-col-header {
  min-width: 120px;  // antes: 108px
  width: 120px;      // antes: 108px
}

.tc-row-header {
  min-width: 168px;  // antes: 156px
}

.tc-corner {
  min-width: 168px;  // antes: 156px
  width: 168px;      // antes: 156px
}

.tc-cell {
  width: 120px;      // antes: 108px
  min-width: 120px;  // antes: 108px
  height: 84px;      // antes: 78px
}
```

### 4. Re-calibrar breakpoints responsive (si 0125 ya implementada)

Ajustar proporcionalmente los media queries para mantener coherencia en tablet y mobile.

---

## Archivos afectados

### Frontend
- `frontend/src/styles/_explore.scss` — tamaños de `.tc-icon`, `.tc-col-header`, `.tc-row-header`, `.tc-corner`, `.tc-cell`

### Backend
Ninguno.

---

## Acceptance criteria

- [ ] Los iconos de tipo en la tabla son visiblemente más grandes (~84px con fondo circular).
- [ ] Las celdas de cabecera y datos se ajustan proporcionalmente al nuevo tamaño.
- [ ] La tabla sigue siendo scrollable horizontalmente en pantallas pequeñas.
- [ ] El drag-filter sigue funcionando correctamente.
- [ ] Sin regresiones en dark mode.
- [ ] Si hay breakpoints responsive (tarea 0125), siguen proporcionados.
