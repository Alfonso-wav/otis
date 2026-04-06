# Task 0111 — Comparador de Pokémon en la vista de detalle

## Goal
Añadir un botón "Comparar" en la vista de detalle de un Pokémon (Pokédex). Al pulsarlo, el usuario puede elegir un segundo Pokémon. Las estadísticas de ambos se superponen en el hexágono de stats (radar chart), cada uno en un color distinto. Debajo del gráfico aparece el nombre de cada Pokémon en el color de su serie.

## Context
- Vista de detalle: `frontend/src/pages/pokedex.ts`, función `renderDetail(p)` (línea ~613). Genera el HTML del detalle incluyendo `<div id="stats-chart">`.
- Radar chart: `frontend/src/charts/stats-chart.ts`. Función `renderStatsChart(container, stats)`. Actualmente dibuja una sola serie roja. Usa ECharts (radar chart).
- Para elegir el segundo Pokémon se puede reutilizar el sistema de autocomplete existente: `frontend/src/autocomplete.ts` + `createAutocomplete()`.
- API: `GetPokemon(name)` en `frontend/src/api.ts` devuelve `core.Pokemon` con `Stats[]`.
- No se necesitan cambios en backend.

## UI / UX
1. En `renderDetail`, añadir un botón `<button id="compare-btn">Comparar</button>` debajo del chart (o junto al título).
2. Al pulsar, mostrar un campo de búsqueda/autocomplete para elegir el segundo Pokémon (puede ser un pequeño modal inline o un input que aparece in situ).
3. Al seleccionar el segundo Pokémon:
   - Llamar `GetPokemon(name)` para obtener sus stats.
   - Re-renderizar el chart con dos series superpuestas:
     - Serie 1 (Pokémon base): color rojo `#e53e3e` (ya existente).
     - Serie 2 (Pokémon comparado): color azul `#3182ce`.
   - Añadir una leyenda debajo del chart con los nombres de ambos Pokémon en sus colores.
4. Un botón "✕ Limpiar comparación" resetea el chart a una sola serie.

## Scope
- **Frontend only**.
- Capas: APP (pokedex.ts) + charts (stats-chart.ts).

## Steps

### 1. Modificar `stats-chart.ts`
Cambiar la firma de `renderStatsChart` para aceptar una segunda serie opcional:

```ts
export interface ChartSeries {
  label: string;
  stats: Stat[];
  color: string;
}

export function renderStatsChart(container: HTMLElement, primary: ChartSeries, secondary?: ChartSeries): void
```

- Si solo hay `primary`, renderiza como ahora (una serie).
- Si hay `secondary`, añadir una segunda entrada en `series[].data` con su color.
- El `radar.indicator` usa el max(255) para todos los stats (sin cambios).

### 2. Actualizar llamada existente en `pokedex.ts`
Línea ~635:
```ts
renderStatsChart(chartContainer, { label: p.Name, stats: p.Stats, color: "#e53e3e" });
```

### 3. Añadir botón "Comparar" y leyenda en `renderDetail`
Añadir al HTML generado:
```html
<div id="compare-controls" class="compare-controls">
  <button id="compare-btn" class="btn btn-secondary">⚖️ Comparar</button>
</div>
<div id="chart-legend" class="chart-legend hidden"></div>
```

### 4. Lógica de comparación en `pokedex.ts`
Después de `renderStatsChart`, adjuntar listener a `#compare-btn`:
- Mostrar un input de autocomplete (reutilizar `createAutocomplete`) para buscar el segundo Pokémon.
- Al seleccionar: llamar `GetPokemon`, re-renderizar el chart con ambas series, mostrar leyenda.
- Botón de limpiar: vuelve a `renderStatsChart` con solo la serie primaria, oculta leyenda.

### 5. Leyenda debajo del chart
```html
<div class="chart-legend">
  <span style="color: #e53e3e">■ {nombre pokemon 1}</span>
  <span style="color: #3182ce">■ {nombre pokemon 2}</span>
</div>
```

### 6. Estilos
Añadir en el SCSS existente (o inline mínimo):
- `.compare-controls`: `margin-top: 8px; display: flex; gap: 8px; justify-content: center`
- `.chart-legend`: `display: flex; gap: 16px; justify-content: center; margin-top: 8px; font-weight: 600`

## Files to modify
- `frontend/src/charts/stats-chart.ts`
- `frontend/src/pages/pokedex.ts`
- Posiblemente `frontend/src/styles/` (SCSS para legend/compare-controls)

## Acceptance criteria
- [ ] El botón "Comparar" aparece en la vista de detalle de cualquier Pokémon
- [ ] Al pulsar, se puede buscar y elegir un segundo Pokémon
- [ ] El hexágono muestra dos series superpuestas en colores distintos
- [ ] Debajo del hexágono aparecen los nombres de ambos Pokémon en sus respectivos colores
- [ ] El botón "Limpiar" elimina la comparación y vuelve a la vista normal
