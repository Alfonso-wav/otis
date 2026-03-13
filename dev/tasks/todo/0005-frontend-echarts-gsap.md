# Integrar ECharts + GSAP para visualización y animaciones

**ID**: 0005-frontend-echarts-gsap
**Estado**: todo
**Fecha**: 2026-03-13
**Depende de**: 0003-frontend-typescript-vite, 0004-frontend-bootstrap-scss

---

## Descripcion

Añadir Apache ECharts para gráficas interactivas de stats de Pokémon (reemplazar barras CSS por radar chart o bar chart real) y GSAP para animaciones de UI (transiciones entre vistas, entrada de cards, hover effects avanzados). Se descarta Framer Motion (requiere React, incompatible) y se posponen D3.js y Theatre.js para cuando haya un caso de uso concreto.

## Capas afectadas

- **Core**: sin cambios
- **Shell**: sin cambios
- **APP**: sin cambios

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/package.json` | modificar | Añadir `echarts`, `gsap` |
| `frontend/src/charts/stats-chart.ts` | crear | Componente ECharts para stats de Pokémon (radar o bar chart) |
| `frontend/src/animations/transitions.ts` | crear | Animaciones GSAP: transición list↔detail, entrada de cards |
| `frontend/src/main.ts` | modificar | Integrar chart en detail view, animaciones en navegación |

## Plan de implementacion

1. Instalar `echarts` y `gsap`
2. Crear `stats-chart.ts`: función que recibe stats[] y renderiza un radar chart con ECharts
3. Integrar el chart en la vista de detalle (reemplazar tabla de barras CSS)
4. Crear `transitions.ts`: animaciones GSAP para transición entre vistas y stagger de cards
5. Integrar animaciones en `main.ts` (reemplazar `classList.add/remove('hidden')` por transiciones GSAP)
6. Lazy loading del chart: solo inicializar ECharts cuando se abre el detail view
7. Verificar rendimiento y funcionalidad

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Chart renderiza correctamente con datos reales, animaciones fluidas sin lag |
| `npm run build` | Build completa, tree-shaking de ECharts funciona |

## Criterios de aceptacion

- [ ] Stats de Pokémon se muestran como gráfica interactiva (radar o bar chart)
- [ ] Tooltip en el chart muestra valores al hover
- [ ] Transición list→detail y detail→list es animada (no display:none instantáneo)
- [ ] Cards del grid aparecen con stagger animation al cargar
- [ ] ECharts se carga lazy (no bloquea el render inicial del grid)
- [ ] Bundle size razonable (ECharts con tree-shaking, solo módulos usados)

## Notas

- **Framer Motion descartado**: requiere React. GSAP cubre todas las necesidades de animación sin dependencia de framework.
- **D3.js pospuesto**: no hay caso de uso que justifique su complejidad sobre ECharts para este proyecto.
- **Theatre.js pospuesto**: útil para scrollytelling/narrativas, no aplica al scope actual.
- ECharts config tipada con interfaces TypeScript como indica el documento.
- GSAP timeline para secuencias complejas si se necesitan en el futuro.
