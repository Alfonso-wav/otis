# Type chart: ajustar tabla de tipos al tamano de la ventana

**ID**: 0161-type-chart-fit-viewport
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripcion

La "tabla de tipos" (heatmap 18x18 en `Explore > Type Chart`) se ve mal porque su tamano no se adapta al viewport:

- Desktop: `.tc-heatmap-container` tiene `height: 600px` fijo (tablet 500px), independiente del alto disponible. Queda o demasiado grande o demasiado pequena dependiendo de la ventana.
- Mobile: `applyMobileLayout` en `type-heatmap.ts` fuerza `width = 96 + 18 * 44 = 888px` y `height = 888px`, provocando scroll horizontal obligatorio y celdas demasiado grandes.
- No hay uso de `aspect-ratio` ni calculo basado en el ancho real del contenedor. La tabla es cuadrada (18x18) pero el contenedor la trata como rectangulo libre.

Objetivo: que la heatmap ocupe el ancho disponible del panel y mantenga proporcion cuadrada (o casi), sin scroll horizontal en desktop/tablet y con celdas legibles pero no oversized en mobile.

## Capas afectadas

- **Core**: ninguna.
- **Shell**: ninguna.
- **APP**: solo frontend (CSS + render del chart).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_explore.scss` | modificar | Reemplazar `height: 600px` fijo de `.tc-heatmap-container` por tamano responsivo: `width: 100%`, `aspect-ratio: 1 / 1` (o calc con `max-height: min(70vh, Npx)`), quitar height fijo. Ajustar breakpoints `@include tablet` y `@include mobile` para misma logica. |
| `frontend/src/charts/type-heatmap.ts` | modificar | Eliminar (o relajar) `applyMobileLayout`: no fijar `width`/`height` en pixeles sobre el container; dejar que CSS mande y que `echarts.init` lea el tamano real. Mantener celdas cuadradas ajustando `grid` (top/left/bottom/right) en funcion del ancho actual; en mobile permitir `minWidth` solo si el viewport es estrecho y la legibilidad lo requiere (valor mas bajo que 888px, p.ej. calcular en base al `window.innerWidth` disponible). |
| `frontend/src/charts/type-heatmap.ts` | modificar | En el handler de resize, recalcular layout y llamar `chart.resize()` tras cambiar el container. Asegurar debounce ya existente sigue OK. |
| `frontend/src/pages/explore/type-chart.ts` | verificar | El wrapper `.type-chart-wrap` no impone max-width que limite el ancho disponible; si lo hace, ajustarlo. |

## Plan de implementacion

1. En `_explore.scss`:
   - `.tc-heatmap-container`: `width: 100%; aspect-ratio: 1 / 1; max-height: calc(100vh - 220px); min-height: 320px;`. Quitar `height: 600px`.
   - `@include tablet`: mismo patron, ajustar `max-height` si hace falta (p.ej. `70vh`).
   - `@include mobile`: quitar `overflow-x: auto` si ya no se necesita; mantener `min-height` razonable (p.ej. `280px`). `aspect-ratio: 1 / 1` mantiene cuadrado pegado al ancho del viewport (~360px) sin forzar scroll.
2. En `type-heatmap.ts`:
   - Eliminar el `applyMobileLayout` que fija pixeles; o reducirlo a un fallback cuando `container.clientWidth < X`.
   - Antes de `echarts.init`, asegurar que el container ya tiene dimensiones (esta ya dentro de `requestAnimationFrame`, OK).
   - Calcular `grid.top/left/bottom/right` en funcion del tamano real: si `width < 480px`, reducir `gridLeft` y `gridTop` (labels mas pequenas) para dar mas espacio a las celdas.
   - Ajustar `labelFontSize` y `cellLabelFontSize` proporcional al ancho del container (p.ej. clamp 9..14).
   - Resize handler: recalcular fontSize/grid y `setOption` parcial o `resize()`.
3. Verificar tipografia mobile: con 360px ancho, celdas ~= (360 - 96 - 40) / 18 ≈ 13px. Demasiado pequeno para labels en celda; aceptar que en mobile las celdas no muestren valor (ocultar `label.show` bajo X px) y dejar colores + tooltip como fuente de info.
4. Probar en desktop (1920, 1440, 1280), tablet (~900), mobile (~360) y APK.

## Tests

| Tipo | Que se verifica |
|------|-----------------|
| Manual desktop 1920 | Heatmap cabe en el panel sin exceder el viewport vertical; cuadrada; labels legibles |
| Manual desktop 1280 | Se reescala proporcionalmente, sin scroll horizontal |
| Manual tablet 900 | Sin scroll horizontal, tamano comodo |
| Manual mobile 360 | No scroll horizontal; heatmap ocupa el ancho; leyenda y radar debajo siguen accesibles |
| Manual resize dinamico | Al redimensionar la ventana, la heatmap se ajusta tras el debounce, sin quedarse cortada |
| Manual APK Android | Mismo resultado que mobile web |

## Criterios de aceptacion

- [ ] `.tc-heatmap-container` ya no tiene `height` fijo en pixeles
- [ ] La heatmap se ajusta al ancho del panel manteniendo proporcion cuadrada
- [ ] Sin scroll horizontal en desktop, tablet ni mobile 360px
- [ ] Labels y celdas siguen legibles (con reduccion de label-in-cell en anchos minimos si procede)
- [ ] Resize dinamico funciona (redimensionar la ventana reajusta la heatmap)
- [ ] Radar debajo no se ve afectado

## Notas

- No tocar logica de efectividad (`OVERRIDES`) ni el radar; solo layout de la heatmap.
- Mantener `clearProps` no aplica (no hay GSAP aqui). Sigue aplicando el patron de `requestAnimationFrame` antes de `echarts.init`.
- Si `aspect-ratio` presenta problemas en algun navegador objetivo, alternativa: calcular height en JS leyendo `container.clientWidth` antes de `echarts.init`.
