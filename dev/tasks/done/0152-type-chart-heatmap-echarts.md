# Feature: visualización heatmap y gráficos ECharts en tabla de tipos

**ID**: 0152-type-chart-heatmap-echarts
**Estado**: todo
**Fecha**: 2026-04-12

---

## Descripción

Añadir visualizaciones interactivas con ECharts a la sección de la tabla de efectividad de tipos (Explorar > Tabla de tipos), usando los datos de cruces de efectividad que ya existen en el frontend.

Visualizaciones propuestas:
1. **Heatmap** (`echarts/charts` → `HeatmapChart`): matriz 18×18 donde el color indica la efectividad (verde = super eficaz, rojo = resistido, gris oscuro = inmune, neutro = base). Reemplaza o complementa la tabla HTML estática actual.
2. **Radar de defensas** por tipo seleccionado: al hacer click en un tipo del heatmap, muestra un radar con cuántos tipos le son super eficaces, neutros, resistidos e inmunes (resumen defensivo).

Los datos ya existen hardcodeados en `frontend/src/pages/explore/type-chart.ts` (objeto `OVERRIDES`, líneas 12-31). No se necesitan cambios en el backend.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — nuevo chart, estilos, i18n

## Archivos a crear/modificar

| Archivo | Acción | Descripción |
|---------|--------|-------------|
| `frontend/src/charts/type-heatmap.ts` | crear | Componente ECharts heatmap 18×18 con datos de efectividad |
| `frontend/src/charts/type-defense-radar.ts` | crear | Radar defensivo por tipo seleccionado |
| `frontend/src/pages/explore/type-chart.ts` | modificar | Integrar las nuevas visualizaciones (toggle o tabs tabla/heatmap/radar) |
| `frontend/src/styles/_explore.scss` | modificar | Estilos para contenedores de los nuevos charts |
| `frontend/src/styles/_dark.scss` | modificar | Dark mode para los nuevos charts |
| `frontend/src/locales/es.json` | modificar | Etiquetas: "Mapa de calor", "Radar defensivo", etc. |
| `frontend/src/locales/en.json` | modificar | Labels: "Heatmap", "Defense radar", etc. |

## Plan de implementación

1. **Extraer datos**: refactorizar `OVERRIDES` de `type-chart.ts` a un módulo compartido o exportarlo para que los charts puedan consumirlo.
2. **Heatmap** (`type-heatmap.ts`):
   - Usar `HeatmapChart` + `VisualMapComponent` + `GridComponent` de ECharts.
   - Ejes X/Y = los 18 tipos con iconos o nombres abreviados.
   - Escala de color: 0 → gris oscuro, 0.5 → rojo/rosa, 1 → blanco/neutro, 2 → verde.
   - Tooltip con nombre del tipo atacante, defensor y multiplicador.
   - Click handler para seleccionar tipo y mostrar radar.
3. **Radar defensivo** (`type-defense-radar.ts`):
   - Usar `RadarChart` (ya usado en `stats-chart.ts` como referencia).
   - Para un tipo seleccionado: indicadores = los 18 tipos, valor = multiplicador de efectividad recibido.
   - Colores coherentes con la paleta de tipos existente.
4. **Integración en `type-chart.ts`**:
   - Añadir toggle/tabs para alternar entre vista tabla y vista heatmap.
   - El radar aparece debajo del heatmap al seleccionar un tipo.
5. **Estilos y dark mode**: contenedores responsivos, colores adaptados.
6. **i18n**: añadir claves de traducción para las etiquetas nuevas.

## Referencia de charts existentes

- `frontend/src/charts/stats-chart.ts` — Radar chart (imports modulares de ECharts)
- `frontend/src/charts/type-distribution.ts` — Pie/donut chart (import completo de ECharts)
- Paleta de colores de tipos en `type-distribution.ts` línea 3-9

## Tests

| Archivo | Qué se testea |
|---------|---------------|
| Manual | Heatmap se renderiza con los 18×18 cruces correctos |
| Manual | Tooltip muestra atacante, defensor y multiplicador |
| Manual | Click en celda del heatmap muestra radar defensivo |
| Manual | Toggle tabla ↔ heatmap funciona |
| Manual | Dark mode correcto |
| Manual | Responsive en mobile |

## Criterios de aceptación

- [ ] Heatmap 18×18 renderiza correctamente con escala de color por efectividad
- [ ] Radar defensivo se muestra al seleccionar un tipo
- [ ] Toggle entre vista tabla y vista heatmap
- [ ] Tooltip informativo en el heatmap
- [ ] Dark mode completo
- [ ] Responsive (mobile y desktop)
- [ ] Textos i18n en ES y EN
