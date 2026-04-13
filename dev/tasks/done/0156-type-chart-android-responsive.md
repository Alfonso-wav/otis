# Tabla efectividad tipos: responsive Android

**ID**: 0156-type-chart-android-responsive
**Estado**: todo
**Fecha**: 2026-04-13

---

## Descripcion

En Android, el heatmap de efectividad de tipos (`frontend/src/charts/type-heatmap.ts` + `.tc-heatmap-container`) se ve muy comprimido: las celdas y los iconos/labels de tipos salen minusculos y son dificiles de pulsar. Actualmente el container tiene `height: 400px` en mobile y el heatmap de 18x18 tipos en 360px de ancho deja celdas de ~15-18px — muy por debajo de los 44x44px requeridos para touch targets (ver `.claude/CLAUDE.md` seccion Mobile).

Objetivo: en viewport mobile (<768px), el heatmap debe permitir pulsar cada celda de tipo comodamente. Opciones:
- A) Hacer el heatmap scrolleable horizontalmente con celdas de tamaño minimo fijo (~44px).
- B) Reducir el numero visible: mostrar selector de tipo y una fila/columna de efectividades (vista condensada mobile).
- C) Aumentar altura del container y permitir zoom/pinch.

Preferida: **A** — scroll horizontal con celdas >=44px, labels legibles, mantiene la vista matriz.

## Capas afectadas

- **Core**: ninguna.
- **Shell**: ninguna.
- **APP**: solo frontend.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/charts/type-heatmap.ts` | modificar | En mobile, fijar `grid` con ancho minimo calculado (18 * 44 + labels), aumentar fontSize labels, forzar canvas mayor que container; permitir scroll horizontal |
| `frontend/src/styles/_explore.scss` | modificar | `.tc-heatmap-container` en `@include mobile`: `overflow-x: auto`, `height` auto con `min-height`, ajustar para que el canvas interno pueda ser mas ancho que el viewport |
| `frontend/src/pages/explore/type-chart.ts` | modificar si necesario | Pasar hint de mobile a `renderTypeHeatmap` o detectar dentro del chart |

## Plan de implementacion

1. Definir constante `MIN_CELL_PX = 44` y `LABEL_RESERVE_PX = 96`.
2. En `renderTypeHeatmap`, detectar ancho viewport (`window.innerWidth < 768`).
3. Si mobile: calcular `canvasWidth = LABEL_RESERVE_PX + 18 * MIN_CELL_PX = 888px`. Aplicarlo al container interior (wrapper mas ancho que `.tc-heatmap-container`) y llamar `chartInstance.resize({ width: canvasWidth })`.
4. Envolver `#tc-heatmap` en un wrapper con `overflow-x: auto`, `-webkit-overflow-scrolling: touch`.
5. Subir `fontSize` de axisLabels x/y de 10 → 12-13 en mobile.
6. Aumentar `grid.left` y `grid.top` en mobile para acomodar labels rotados sin cortar.
7. Testear tap: cada celda debe responder; scroll horizontal fluido.
8. Verificar que en tablet/desktop comportamiento no cambia.
9. Añadir `window.addEventListener("resize")` para recalcular al rotar pantalla.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Emulador Android viewport 360x800: celdas >=44px, scroll horizontal funciona |
| Manual | Viewport 768px (tablet): sin scroll, layout normal |
| Manual | Viewport 1200px (desktop): sin cambios visibles |
| Manual | Rotar dispositivo: heatmap se redibuja correctamente |
| Manual | Tap en celda abre tooltip + click selecciona tipo para radar |

## Criterios de aceptacion

- [ ] En 360px de ancho, cada celda del heatmap es >=44x44px
- [ ] Labels de ejes (nombres de tipos) legibles sin zoom en mobile
- [ ] Scroll horizontal fluido en el container del heatmap
- [ ] Click/tap en celda sigue seleccionando tipo para radar
- [ ] Desktop/tablet sin regresion visual
- [ ] Rotacion de pantalla (landscape/portrait) recalcula tamaño
- [ ] Dark mode mantiene estilos correctos

## Notas

- `.tc-heatmap-container` actual: `height: 400px` en mobile (`_explore.scss:619`). Probable que haya que pasar a `height: auto` + `min-height` y poner overflow en wrapper interno, porque ECharts no puede crecer mas que el container.
- Alternativa valida: reemplazar heatmap por vista "selector de tipo + lista de efectividades" solo en mobile. Evaluar UX con usuario antes de implementar.
- Checklist CLAUDE.md: touch targets 44px, `overflow-x: auto` en tablas, probar en 360px.
