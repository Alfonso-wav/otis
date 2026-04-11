# Task 0136 — Comparador: barra de totales por stat con colores y KPIs

## Estado: done

## Goal
Cuando el usuario compara Pokémon en la vista de detalle, mostrar debajo de la leyenda del gráfico radar:
1. **Barra de stat total (BST)** para cada Pokémon comparado: una barra horizontal segmentada donde cada stat tiene su propio color, representando visualmente la suma total (max ~720–780).
2. **Diferencia de BST** entre el Pokémon base y cada comparado (e.g. "+45" o "−30").
3. **KPIs generales** al comparar: stat más alto de cada uno, stat donde el base gana más y stat donde pierde más.

## Context
- El comparador actual vive en `frontend/src/pages/pokedex.ts` (líneas 647–723).
- Usa `renderStatsChart()` de `frontend/src/charts/stats-chart.ts` con un `ChartSeries[]`.
- Los colores de comparación son `COMPARE_COLORS` (8 colores) — estos son para distinguir Pokémon, NO para stats.
- Cada Pokémon tiene 6 stats: HP, Attack, Defense, Sp. Atk, Sp. Def, Speed.
- Los nombres de stats se traducen con `statName()` de `frontend/src/i18n.ts`.
- Las traducciones están en `frontend/src/locales/en.json` y `es.json`.
- El layout del detalle es: gráfico radar → compare-controls → chart-legend → lore → encounters → abilities → moves.

## Acceptance criteria
- [ ] Debajo de la leyenda, aparece una sección de "stat totals" cuando hay 2+ Pokémon en la comparación
- [ ] Para cada Pokémon comparado se muestra una barra horizontal segmentada con los 6 stats, cada segmento con un color distinto propio del stat (no del Pokémon)
- [ ] Los colores de stats son consistentes: mismo color para HP siempre, mismo color para Attack siempre, etc.
- [ ] Al lado de cada barra se muestra el nombre del Pokémon, el total numérico y la diferencia vs el base (e.g. "+45" verde, "−30" rojo)
- [ ] Se muestran KPIs generales: "Stat más alto: [Pokémon] — [stat] ([valor])", "Mayor ventaja: [stat] (+N)", "Mayor desventaja: [stat] (−N)"
- [ ] Los KPIs solo se muestran cuando hay exactamente 2 Pokémon (base + 1 comparado)
- [ ] Se oculta la sección cuando solo queda el Pokémon base
- [ ] Las etiquetas de UI están en los locales EN y ES
- [ ] La sección respeta dark mode
- [ ] La sección se actualiza al añadir/eliminar Pokémon del comparador

## Layers involved
- **Frontend** (`frontend/src/pages/pokedex.ts`): lógica de renderizado de la nueva sección
- **Frontend** (`frontend/src/styles/_pokemon.scss`): estilos de barras y KPIs
- **Frontend** (`frontend/src/styles/_dark.scss`): variantes dark mode
- **Frontend** (`frontend/src/locales/en.json`, `es.json`): nuevas claves de traducción

## Stat color palette (propuesta)
| Stat        | Color   | Referencia visual      |
|-------------|---------|------------------------|
| HP          | #ff5959 | rojo vida              |
| Attack      | #f5ac78 | naranja ofensivo       |
| Defense     | #fae078 | amarillo escudo        |
| Sp. Atk     | #9db7f5 | azul claro especial    |
| Sp. Def     | #a7db8d | verde defensa especial |
| Speed       | #fa92b2 | rosa velocidad         |

Estos colores son los clásicos de Pokémon usados en wikis y herramientas de stats.

## Notes
- No requiere cambios en Go (Core/Shell) — toda la data de stats ya está disponible en el frontend.
- Las barras segmentadas son HTML/CSS puro (no ECharts) para mantenerlo simple.
- El BST máximo teórico es ~720 (Arceus) pero la barra puede escalar a 800 para dar margen.
