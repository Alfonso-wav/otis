# Feature: radar ofensivo superpuesto al defensivo en tabla de tipos

**ID**: 0154-type-chart-offensive-radar
**Estado**: done
**Fecha**: 2026-04-13

---

## Descripción

En la sección Explorar > Tabla de tipos, al seleccionar un tipo en el heatmap se muestra actualmente un **radar defensivo** (cómo le afecta al tipo seleccionado cada uno de los 18 atacantes). Añadir un segundo radar, el **radar ofensivo**, que muestre la efectividad del tipo seleccionado cuando ataca a cada uno de los 18 tipos.

Ambos radares deben:
- Renderizarse **superpuestos** en el mismo gráfico (dos series en el mismo `RadarChart` de ECharts), compartiendo indicadores y escala.
- Tener controles para **mostrar/ocultar** cada serie independientemente (checkboxes, toggles o leyenda clickable de ECharts).
- Diferenciarse visualmente: por ejemplo, defensivo en color del tipo y ofensivo en un color complementario o con distinto estilo de línea/área.

Los datos ya existen en `OVERRIDES` (`frontend/src/pages/explore/type-chart.ts`) y la función `effectiveness(attacker, defender)` los expone. El radar ofensivo simplemente invierte los argumentos: para el tipo `T` seleccionado, `values[i] = effectiveness(T, ALL_TYPES[i])`.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — modificar chart existente, estilos, i18n

## Archivos a crear/modificar

| Archivo | Acción | Descripción |
|---------|--------|-------------|
| `frontend/src/charts/type-defense-radar.ts` | modificar | Renombrar conceptualmente y añadir serie ofensiva. Exportar función unificada que recibe qué series mostrar (defensivo/ofensivo/ambos) |
| `frontend/src/pages/explore/type-chart.ts` | modificar | Añadir controles de toggle para cada radar, pasar estado al chart, re-render al cambiar |
| `frontend/src/styles/_explore.scss` | modificar | Estilos para los controles de toggle (checkboxes/pills) y leyenda |
| `frontend/src/styles/_dark.scss` | modificar | Dark mode para los nuevos controles |
| `frontend/src/locales/es.json` | modificar | Añadir claves: `offensiveRadar`, `defensiveRadar` (ya existe), `toggleDefensive`, `toggleOffensive` o similar |
| `frontend/src/locales/en.json` | modificar | Mismas claves en inglés |

## Plan de implementación

1. **Extender el chart** (`type-defense-radar.ts`):
   - Renombrar función a `renderTypeRadar` (o mantener compatibilidad con alias) y aceptar opciones `{ showDefense: boolean, showOffense: boolean }`.
   - Calcular dos series:
     - Defensa: `values[i] = effectiveness(ALL_TYPES[i], typeName)` (actual).
     - Ataque: `values[i] = effectiveness(typeName, ALL_TYPES[i])`.
   - Asignar colores distintos y estilo de línea para que se distingan al superponerse (ej. defensivo con `areaStyle` sólido, ofensivo con `areaStyle` rayado o sólo línea).
   - Usar la `legend` de ECharts o los controles externos del paso 2 para toggling.

2. **Controles en `type-chart.ts`**:
   - Añadir un par de checkboxes o pills bajo el título del radar: "Defensivo" / "Ofensivo", ambos marcados por defecto.
   - Guardar estado (`showDefense`, `showOffense`) en variables de módulo (análogo a `selectedRadarType`).
   - Listener `change` → re-render del radar con nuevas opciones (sin recrear instancia innecesariamente si es posible con `setOption`).
   - Al cambiar de tipo seleccionado o idioma, preservar estado de toggles.

3. **Título/hint del radar**: actualizar texto a "Radar de efectividad: {tipo}" (o similar) en lugar de solo "Radar defensivo".

4. **Estilos**:
   - Controles alineados horizontalmente sobre el radar.
   - Touch targets ≥ 44px en mobile.
   - Dark mode: colores de fondo, borde y texto de los pills/checkboxes.

5. **i18n**: añadir claves en `es.json` y `en.json` simultáneamente. Vista escucha `locale-changed` ya (existe listener en `initTypeChart`), verificar que preserva estado de toggles al recargar idioma.

## Tests

| Archivo | Qué se testea |
|---------|---------------|
| Manual | Al seleccionar un tipo aparecen ambos radares superpuestos por defecto |
| Manual | Checkbox "Defensivo" oculta/muestra la serie defensiva |
| Manual | Checkbox "Ofensivo" oculta/muestra la serie ofensiva |
| Manual | Con ambos desactivados el radar queda vacío (o se oculta) sin romper layout |
| Manual | Cambiar de tipo en el heatmap preserva el estado de los toggles |
| Manual | Cambiar de idioma preserva el estado y traduce controles |
| Manual | Dark mode correcto para controles y chart |
| Manual | Responsive en viewport 360px |

## Criterios de aceptación

- [ ] Radar defensivo y ofensivo renderizan superpuestos en el mismo gráfico
- [ ] Controles (checkbox/pill/leyenda) permiten mostrar/ocultar cada serie independientemente
- [ ] Ambas series usan colores/estilos distinguibles
- [ ] Tooltip muestra nombre del tipo y multiplicador para cada serie
- [ ] Textos en `es.json` y `en.json`, vista reacciona a `locale-changed`
- [ ] Animaciones GSAP (si se añaden) con `clearProps`
- [ ] Handlers async (si aplica) con guard flag y `try/finally`
- [ ] Dark mode completo
- [ ] Responsive mobile (360px) con touch targets ≥ 44px
- [ ] Probado en build de producción
