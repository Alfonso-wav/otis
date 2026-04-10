# Task 0125 — Type Chart: responsivo y sin espacio vacío al filtrar

## Estado: pending

## Goal

Dos mejoras de adaptabilidad sobre la **Tabla de Tipos** en Explorar:

1. **Sin espacio en blanco al filtrar**: cuando el usuario filtra tipos y la tabla se queda más pequeña, el contenedor debe reducirse al tamaño de la tabla en lugar de dejar espacio en blanco a la derecha.
2. **Responsivo para móvil**: añadir media queries para escalar celdas e iconos en pantallas pequeñas, de forma que la tabla sea usable en dispositivos móviles.

---

## Contexto técnico

### Archivos principales

- `frontend/src/styles/_explore.scss` — estilos del type chart (líneas 461–611)
- `frontend/src/pages/explore/type-chart.ts` — HTML del componente (solo lectura)

### Problema 1: Espacio en blanco al filtrar

`.type-chart-scroll` es un elemento de bloque (`display: block` por defecto) y ocupa el 100% del ancho del padre, incluso si la tabla interna es más estrecha. Cuando se filtran tipos (menos columnas), la tabla se encoge pero el contenedor scroll mantiene el ancho completo, dejando espacio vacío.

**Fix**: añadir `width: fit-content; max-width: 100%` al `.type-chart-scroll`. Con esto:
- El contenedor se adapta al ancho de la tabla (sin espacio vacío).
- Nunca excede el padre (`max-width: 100%`), y si la tabla es más ancha, el `overflow-x: auto` activa el scroll.

Las etiquetas `.type-chart-labels` (con `justify-content: space-between`) siguen ocupando el ancho del padre, lo cual es correcto (la etiqueta "Atacante →" queda a la izquierda y "↓ Defensor" a la derecha del contenedor completo, independientemente del tamaño de la tabla).

### Problema 2: Responsivo para móvil

**Dimensiones post-tarea 0124** (base para los breakpoints):
```scss
.tc-col-header  { min-width: 108px; width: 108px; }
.tc-row-header  { min-width: 156px; }
.tc-corner      { min-width: 156px; width: 156px; }
.tc-cell        { width: 108px; min-width: 108px; height: 78px; }
.tc-icon        { width: 60px; height: 60px; }
.type-chart-table { font-size: 1rem; }
```

Con 18 tipos, la tabla tendría un ancho mínimo de: `156px (corner) + 18 × 108px (celdas) = 2100px`. En móvil se necesitan celdas más pequeñas.

**Breakpoints propuestos**:

| Breakpoint | Celda (ancho) | Celda (alto) | Icono | Corner/Row header | Fuente |
|---|---|---|---|---|---|
| ≥ 992px (desktop) | 108px | 78px | 60px | 156px | 1rem |
| 768–991px (tablet) | 72px | 52px | 40px | 104px | 0.85rem |
| 480–767px (móvil grande) | 44px | 32px | 26px | 64px | 0.72rem |
| < 480px (móvil pequeño) | 30px | 22px | 18px | 44px | 0.62rem |

Con < 480px y tabla completa (18 tipos): `44px + 18 × 30px = 584px` → scroll horizontal inevitable, pero es lo esperado para una tabla 18×18. Con filtrado activo (ej. 6 tipos): `44px + 6 × 30px = 224px`, cabe sin scroll.

---

## Cambios requeridos

### 1. Fix espacio en blanco — `_explore.scss`

En el bloque `.type-chart-scroll` (actualmente sin `width`), añadir:

```scss
.type-chart-scroll {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  border-radius: $border-radius-lg;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  // AÑADIR:
  width: fit-content;
  max-width: 100%;
}
```

### 2. Media queries para móvil — `_explore.scss`

Añadir después del bloque `.tc-immune` (tras la línea 611):

```scss
// ─── Type Chart — responsive ─────────────────────────────────────────────────

@media (max-width: 991px) {
  .tc-col-header  { min-width: 72px; width: 72px; padding: 0.6rem 0.3rem; }
  .tc-row-header  { min-width: 104px; padding: 0.5rem 0.8rem; }
  .tc-corner      { min-width: 104px; width: 104px; }
  .tc-cell        { width: 72px; min-width: 72px; height: 52px; font-size: 0.85rem; }
  .tc-icon        { width: 40px; height: 40px; }
  .type-chart-table { font-size: 0.85rem; }
}

@media (max-width: 767px) {
  .tc-col-header  { min-width: 44px; width: 44px; padding: 0.3rem 0.15rem; }
  .tc-row-header  { min-width: 64px; padding: 0.25rem 0.4rem; }
  .tc-corner      { min-width: 64px; width: 64px; }
  .tc-cell        { width: 44px; min-width: 44px; height: 32px; font-size: 0.72rem; }
  .tc-icon        { width: 26px; height: 26px; }
  .type-chart-table { font-size: 0.72rem; }
}

@media (max-width: 479px) {
  .tc-col-header  { min-width: 30px; width: 30px; padding: 0.2rem 0.1rem; }
  .tc-row-header  { min-width: 44px; padding: 0.2rem 0.3rem; }
  .tc-corner      { min-width: 44px; width: 44px; }
  .tc-cell        { width: 30px; min-width: 30px; height: 22px; font-size: 0.62rem; }
  .tc-icon        { width: 18px; height: 18px; }
  .type-chart-table { font-size: 0.62rem; }
}
```

**Nota sobre tarea 0124**: si el icono circular (padding en `.tc-icon`) de 0124 está implementado, añadir override del padding en el breakpoint < 480px:
```scss
@media (max-width: 479px) {
  @each $type, $color in $type-colors {
    .tc-col-header[data-type="#{$type}"] .tc-icon,
    .tc-row-header[data-type="#{$type}"] .tc-icon {
      padding: 4px; // reducido desde 8px del diseño base (tarea 0124)
    }
  }
}
```

---

## Archivos afectados

### Frontend
- `frontend/src/styles/_explore.scss` — fix fit-content + media queries responsive

### Backend
Ninguno.

---

## Acceptance criteria

- [ ] Al filtrar tipos, el contenedor de la tabla se encoge con la tabla (sin espacio en blanco a la derecha).
- [ ] En desktop (≥ 992px): tabla con celdas de tamaño completo (108px si 0124 ya está implementada).
- [ ] En tablet (768–991px): celdas ~72px, tabla usable con pocos tipos filtrados sin scroll horizontal.
- [ ] En móvil grande (480–767px): celdas ~44px, scroll horizontal disponible para tabla completa.
- [ ] En móvil pequeño (< 480px): celdas ~30px, tabla funcional con y sin filtrado.
- [ ] El drag-filter sigue funcionando en todos los breakpoints.
- [ ] El botón "Restaurar tipos" sigue funcionando en todos los breakpoints.
- [ ] Sin regresiones en dark mode.
- [ ] Sin errores en consola.

---

## Dependencias

- **Tarea 0124** (type-chart-layout-size-and-icon-colors): esta tarea asume que 0124 ya fue implementada, ya que los breakpoints están calibrados para las dimensiones post-0124 (celdas de 108px base). Si se implementa antes de 0124, los breakpoints del responsive deberán ajustarse a las dimensiones actuales (36px base).
