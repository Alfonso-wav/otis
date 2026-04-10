# Task 0124 — Type Chart: tamaño x3, botón abajo, color en icono

## Estado: pending

## Goal

Cuatro mejoras visuales sobre la **Tabla de Tipos** en Explorar:

1. **Mover el botón "Restaurar tipos"** al final de la tabla (debajo del scroll container), con más espacio vertical, para evitar interferencias con el resize visual.
2. **Hacer la tabla x3 de tamaño**: triplicar las dimensiones de celdas, cabeceras, icono y fuente.
3. **Ajustar el slider de filtro** al nuevo tamaño (principalmente CSS; la lógica de `elementFromPoint` ya funciona con cualquier tamaño).
4. **El color del tipo va en el icono, no en el fondo de la cabecera**: eliminar el `background-color` semitransparente de los `th`, y aplicar el color como fondo circular pequeño al propio `.tc-icon`.

---

## Contexto técnico

### Archivos principales

- `frontend/src/pages/explore/type-chart.ts` — lógica y HTML del componente
- `frontend/src/styles/_explore.scss` — estilos (líneas 461–611)

### Estado actual

**Botón restaurar** (líneas 80-88 de `type-chart.ts`):
```typescript
const filterBar = filteredTypes.size > 0
  ? `<div class="tc-filter-bar"><button class="tc-restore-btn">...</button></div>`
  : "";

panel.innerHTML = `
  <div class="type-chart-wrap">
    <h2 class="type-chart-title">${title}</h2>
    ${filterBar}           ← actualmente ARRIBA, entre title y labels
    <div class="type-chart-labels">...</div>
    <div class="type-chart-scroll">...</div>
  </div>
`;
```

**Dimensiones actuales** (`_explore.scss`):
```scss
.tc-col-header  { min-width: 36px; width: 36px; }
.tc-row-header  { min-width: 52px; }
.tc-corner      { min-width: 52px; width: 52px; }
.tc-cell        { width: 36px; min-width: 36px; height: 26px; }
.tc-icon        { width: 20px; height: 20px; }
.type-chart-table { font-size: 0.68rem; }
```

**Color de tipo en cabeceras** (`_explore.scss` líneas 538-544):
```scss
@each $type, $color in $type-colors {
  .tc-col-header[data-type="#{$type}"],
  .tc-row-header[data-type="#{$type}"] {
    background-color: rgba($color, 0.22);
  }
}
```

Los SVG en `/public/assets/types/*.svg` son **blancos** (`fill="white"`) sobre transparente. Se cargan como `<img>` sin inline SVG.

---

## Cambios requeridos

### 1. Mover botón al fondo — `type-chart.ts`

Mover `${filterBar}` al final del `type-chart-wrap`, después del `type-chart-scroll`:

```typescript
panel.innerHTML = `
  <div class="type-chart-wrap">
    <h2 class="type-chart-title">${title}</h2>
    <div class="type-chart-labels">
      <span class="tc-attacking-label">${attackingLabel}</span>
      <span class="tc-defending-label">${defendingLabel}</span>
    </div>
    <div class="type-chart-scroll">
      <table class="type-chart-table" aria-label="${title}">
        <thead>
          <tr>
            <th class="tc-corner"></th>
            ${headerCells}
          </tr>
        </thead>
        <tbody>${rows}</tbody>
      </table>
    </div>
    ${filterBar}
  </div>
`;
```

### 2. Tabla x3 de tamaño — `_explore.scss`

Triplicar las dimensiones:

```scss
.tc-col-header  { min-width: 108px; width: 108px; padding: 0.9rem 0.45rem; }
.tc-row-header  { min-width: 156px; padding: 0.75rem 1.2rem; }
.tc-corner      { min-width: 156px; width: 156px; }
.tc-cell        { width: 108px; min-width: 108px; height: 78px; font-size: 1.1rem; }
.tc-icon        { width: 60px; height: 60px; }
.type-chart-table { font-size: 1rem; }
```

### 3. Ajuste del botón restaurar — `_explore.scss`

Cambiar `margin-bottom` por `margin-top` para que respire debajo de la tabla:

```scss
.tc-filter-bar {
  margin-top: 1.25rem;
  display: flex;
  justify-content: center;   // centrado bajo la tabla
}

.tc-restore-btn {
  // mantener estilos actuales, ajustar padding al nuevo tamaño
  padding: 0.4rem 1.4rem;
  font-size: 0.9rem;
}
```

### 4. Color del tipo en el icono, no en el fondo — `_explore.scss`

**Eliminar** el bloque actual de backgrounds de cabecera:
```scss
// ELIMINAR esto:
@each $type, $color in $type-colors {
  .tc-col-header[data-type="#{$type}"],
  .tc-row-header[data-type="#{$type}"] {
    background-color: rgba($color, 0.22);
  }
}
```

**Añadir** color como fondo circular en el icono:
```scss
@each $type, $color in $type-colors {
  .tc-col-header[data-type="#{$type}"] .tc-icon,
  .tc-row-header[data-type="#{$type}"] .tc-icon {
    background-color: $color;
    border-radius: 50%;
    padding: 8px;           // espacio interno para que el icono blanco no ocupe todo el círculo
    box-sizing: border-box;
  }
}
```

Con esto: la cabecera tiene fondo neutro (`#f7fafc`), y el icono SVG blanco aparece sobre un círculo de color del tipo. En dark mode, ajustar si es necesario (los colores ya son vivos, pueden funcionar tal cual).

**Nota**: la misma corrección aplica en `_dark.scss` para el bloque equivalente (si existe para dark mode):
```scss
// Buscar y ELIMINAR en _dark.scss el bloque de background-color para tc-col-header/tc-row-header de tipo
```

---

## Archivos afectados

### Frontend
- `frontend/src/pages/explore/type-chart.ts` — mover `${filterBar}` al final del HTML
- `frontend/src/styles/_explore.scss` — dimensiones x3, botón abajo, color en icono
- `frontend/src/styles/_dark.scss` — eliminar override de background de cabecera (si existe)

### Backend
Ninguno.

---

## Acceptance criteria

- [ ] El botón "Restaurar tipos" aparece **debajo** de la tabla (no encima), con margen superior visible.
- [ ] Las celdas de la tabla son aproximadamente x3 del tamaño original (cell ≈108px ancho, ≈78px alto).
- [ ] Los iconos de tipo son ≈60×60px y tienen fondo circular del color de su tipo.
- [ ] Las cabeceras (`th`) NO tienen tinte de color de fondo; solo el círculo del icono lo tiene.
- [ ] El scroll horizontal sigue funcionando con la tabla más grande.
- [ ] El drag-filter (arrastrar cabeceras) sigue funcionando correctamente con el nuevo tamaño.
- [ ] El botón restaurar sigue funcionando (limpia filtros y re-renderiza).
- [ ] Dark mode: los iconos circulares se ven bien en fondo oscuro.
- [ ] Sin errores en consola del frontend.

---

## Dependencias

Ninguna. Tarea visual independiente sobre el componente type-chart existente.
