# Fix: donut de distribución de tipos no se muestra en Explorer > Regiones

**ID**: 0144-fix-region-donut-chart
**Estado**: todo
**Fecha**: 2026-04-12

---

## Descripcion

El gráfico donut de distribución de tipos de cada región ha dejado de renderizarse al expandir una tarjeta de región en Explorer > Regiones. El contenedor `#chart-region-{name}` se crea correctamente, pero el chart de ECharts no aparece.

Posibles causas a investigar:
- El contenedor `.region-chart` (height: 220px) puede no tener dimensiones calculadas en el momento de `echarts.init()` si está oculto o recién insertado.
- El `GetRegionTypeDistribution()` puede estar fallando silenciosamente o devolviendo datos vacíos.
- Un cambio reciente en estilos o layout puede haber alterado la visibilidad del contenedor en el momento del render.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna (verificar API si hay regresión en `GetRegionTypeDistribution`)
- **APP**: frontend — charts, pages/explore/regions, estilos

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/explore/regions.ts` | investigar/modificar | Verificar que el contenedor existe y tiene dimensiones cuando se llama a `renderTypeDistributionChart` |
| `frontend/src/charts/type-distribution.ts` | investigar/modificar | Verificar que `echarts.init(el)` funciona con el elemento y que se llama `chart.resize()` tras el render |
| `frontend/src/styles/_explore.scss` | investigar | Verificar que `.region-chart` tiene height correcto y no está oculto |

## Plan de implementacion

1. Reproducir el bug expandiendo una región (ej. Kanto) y verificar si el contenedor `.region-chart` se crea con dimensiones.
2. Verificar en consola si `GetRegionTypeDistribution` devuelve datos o error.
3. Si el contenedor no tiene dimensiones al momento del init de ECharts, añadir un `requestAnimationFrame` o `setTimeout` antes de inicializar el chart para asegurar que el DOM está listo.
4. Considerar llamar `chart.resize()` explícitamente después de que el cuerpo de la card termine su animación de expansión.
5. Verificar en dark mode también.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Expandir cada región y verificar que aparece el donut |
| Manual | Verificar que al hacer click en un tipo del donut se abre el modal |
| Manual | Verificar en dark mode |
| Manual | Verificar en mobile |

## Criterios de aceptacion

- [ ] Al expandir una región, el gráfico donut de tipos se renderiza correctamente
- [ ] El donut es interactivo (click en tipo abre modal)
- [ ] Funciona en light mode y dark mode
- [ ] Funciona en desktop y mobile
