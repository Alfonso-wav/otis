# Task 0116 — Vista detalle: botón comparador cuadrado + tabla movimientos ordenable y filtro múltiple

## Estado: done

## Goal
Dos mejoras en la vista individual de Pokémon:
1. El botón `×` para eliminar un Pokémon del comparador se ve como un rectángulo alargado verticalmente. Debe ser un cuadrado ligeramente mayor que el carácter `×`.
2. La tabla de movimientos debe permitir ordenar por cualquier columna desde su cabecera, y los filtros de método deben admitir selección múltiple.

## Contexto técnico

### Botón comparador (remove btn)
- Clase CSS: `.compare-remove-btn`
- Renderizado en `frontend/src/pages/pokedex.ts` → función `rebuildLegend()` (~línea 666)
- Estilos en: no tiene bloque dedicado, hereda estilos de `.chart-legend-entry`
- El botón contiene el carácter `×` y necesita `width = height` con padding uniforme

### Tabla de movimientos
- Función de render: `renderMoves(moves: PokemonMoveEntry[])` en `frontend/src/pages/pokedex.ts` (~línea 993)
- Cabeceras actuales: "Nombre" (`detail.moveName`), "Método" (`detail.moveMethodCol`), "Nivel" (`detail.moveLevelCol`)
- Columnas **sin ordenación** actualmente
- Filtros actuales: chips de selección única (`.filter-chip`) con `activeFilter: MoveMethodFilter`; opciones: `["all", "level-up", "machine", "egg", "tutor"]`
- Estilos `.filter-chip` en `frontend/src/styles/_components.scss` (~línea 201)
- Estilos `.moves-table-wrap` / `.moves-table` en `frontend/src/styles/_pokemon.scss` (~línea 726)

### i18n
- Locale activo: español por defecto
- Claves relevantes en `frontend/src/locales/es.json` y `en.json` (sección `detail`)

## Acceptance criteria

### Botón ×
- [ ] `.compare-remove-btn` tiene `width` y `height` iguales (cuadrado), con dimensión ligeramente mayor que el `×`
- [ ] El botón no aparece alargado verticalmente en ningún navegador/tamaño
- [ ] El aspecto visual (color, hover) queda coherente con el resto de la UI del comparador

### Tabla movimientos — ordenación
- [ ] Cada cabecera de columna muestra un indicador visual de que es clickable (cursor pointer, icono ↕/↑/↓)
- [ ] Al hacer click en una cabecera, las filas se ordenan por esa columna (asc → desc → sin orden)
- [ ] La columna "Nivel" ordena numéricamente (los movimientos sin nivel van al final)
- [ ] La columna "Nombre" ordena alfabéticamente
- [ ] La columna "Método" ordena alfabéticamente por método
- [ ] El estado de ordenación se resetea al cambiar el filtro de métodos

### Tabla movimientos — filtro múltiple
- [ ] Los chips de método permiten seleccionar varios a la vez (toggle individual)
- [ ] Si ningún chip está activo, se muestran todos los movimientos (equivale a "Todos")
- [ ] El chip "Todos" / "all" deselecciona el resto y actúa como reset
- [ ] La lógica de filtrado muestra la unión de los métodos seleccionados

## Implementación sugerida

### 1. Fix botón ×
En `_pokemon.scss` (o en `_components.scss`), añadir bloque:
```scss
.compare-remove-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.4rem;
  height: 1.4rem;
  padding: 0;
  line-height: 1;
  // ... colores y hover existentes
}
```

### 2. Ordenación de columnas
- Añadir estado `sortCol: 'name' | 'method' | 'level' | null` y `sortDir: 'asc' | 'desc'` dentro de `renderMoves()`
- Modificar `<th>` para incluir `data-col` y event listener de click
- Aplicar sort a las filas antes de renderizar
- Añadir clase CSS `.sort-asc` / `.sort-desc` al `<th>` activo para mostrar el indicador

### 3. Filtro múltiple
- Cambiar `activeFilter: MoveMethodFilter` por `activeFilters: Set<MoveMethodFilter>`
- El chip "all" vacía el Set y muestra todo
- El resto de chips añaden/eliminan su método del Set
- La condición de filtrado: `activeFilters.size === 0 || activeFilters.has(m.Method)`

## Archivos afectados
- `frontend/src/pages/pokedex.ts` — lógica `rebuildLegend()` y `renderMoves()`
- `frontend/src/styles/_pokemon.scss` — estilos `.compare-remove-btn`, cabeceras de columna
- `frontend/src/styles/_components.scss` — posiblemente estilos `.filter-chip` para estado multi-active
- `frontend/src/styles/_dark.scss` — si hay overrides de dark mode para estos elementos
