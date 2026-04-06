# Task 0115 — Stats chart: hexágono y red interior más grandes

## Estado: done

## Goal
El gráfico de estadísticas hexagonal (radar chart de ECharts) resulta demasiado pequeño, especialmente en la vista de comparador donde se superponen varios Pokémon y la red interior es difícil de leer e identificar. Se quiere aumentar considerablemente el tamaño visible del hexágono y la malla interior.

## Contexto técnico
- El contenedor del chart está en `frontend/src/pages/pokedex.ts` línea ~634: `style="width:100%;height:300px;"`
- La configuración del radar está en `frontend/src/charts/stats-chart.ts`
- ECharts permite controlar el tamaño del radar mediante la opción `radius` dentro de `radar` (por defecto ~75% del área)
- También se puede aumentar la altura del contenedor para dar más espacio total al gráfico

## Acceptance criteria
- [ ] El hexágono ocupa notablemente más espacio visual en la vista de detalle
- [ ] La red interior (líneas de guía y áreas) es claramente legible
- [ ] En la vista de comparador (2+ Pokémon) las series superpuestas siguen siendo distinguibles
- [ ] El cambio no rompe el layout del panel de detalle en ningún tamaño de pantalla
- [ ] Las etiquetas de los ejes (HP, Atk, Def…) no quedan cortadas ni solapadas con otros elementos

## Implementación sugerida
1. Aumentar la altura del contenedor `#stats-chart` de `300px` a `420px` (o similar) en `pokedex.ts`
2. En `stats-chart.ts`, añadir/ajustar `radius` en la configuración `radar` para que el polígono ocupe más porcentaje del área (e.g. `radius: '80%'`)
3. Revisar que el `center` del radar siga siendo `['50%', '50%']` para que esté centrado
4. Verificar visualmente en modo comparador con 2-3 Pokémon

## Archivos afectados
- `frontend/src/pages/pokedex.ts` — altura del contenedor `#stats-chart`
- `frontend/src/charts/stats-chart.ts` — opción `radius` (y posiblemente `center`) del radar
