# Toggle de visibilidad de columnas en todas las tablas

**ID**: 0081-table-column-visibility-toggle
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Añadir en todas las tablas de la aplicación un icono de ojo (👁) junto al nombre de cada columna que permita ocultar/mostrar esa columna. La columna de nombre/identificador principal de cada tabla no se puede ocultar. La preferencia del usuario se persiste en `localStorage` para que se mantenga entre sesiones.

## Tablas afectadas

| Tabla | Página | Columna fija (no ocultable) | Columnas ocultables |
|-------|--------|-----------------------------|---------------------|
| Pokémon Stats | Pokédex (table view) | Name | #, Sprite, Types, HP, ATK, DEF, SPA, SPD, Vel, Total |
| Encounters | Pokédex detalle | Location | Game, Method, Chance, Levels, Conditions |
| Moves | Explore > Moves | Name | Type, Category, Power, Accuracy, PP, Priority |
| Abilities | Explore > Abilities | Name | Description, Pokémon Count |
| Stats Config (IV/EV) | Builds | Stat | IV, EV |
| Damage Output (Attacker) | Builds | Move | Type, Category, Min, Max, Effectiveness |
| Damage Output (Defender) | Builds | Move | Type, Category, Min, Max, Effectiveness |

## Capas afectadas

- **APP (frontend)**: Componente reutilizable de toggle, lógica de visibilidad en cada tabla, persistencia en localStorage, estilos.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/components/column-toggle.ts` | crear | Módulo reutilizable: función que recibe config de columnas y tabla-id, renderiza iconos de ojo en los `<th>`, gestiona estado de visibilidad, aplica CSS para ocultar/mostrar columnas, y persiste en localStorage |
| `frontend/src/pages/pokedex.ts` | modificar | Integrar column-toggle en la tabla de stats del Pokédex y en la tabla de encounters |
| `frontend/src/pages/explore/moves.ts` | modificar | Integrar column-toggle en la tabla de moves |
| `frontend/src/pages/explore/abilities.ts` | modificar | Integrar column-toggle en la tabla de abilities |
| `frontend/src/pages/builds.ts` | modificar | Integrar column-toggle en la tabla de stats config (IV/EV) y en las dos tablas de damage output (attacker y defender) |
| `frontend/src/styles/_components.scss` | modificar | Estilos del icono de ojo: tamaño, color, hover, transición, estado activo/inactivo |
| `frontend/src/styles/_dark.scss` | modificar | Estilos del icono de ojo en dark mode |

## Plan de implementacion

1. **Crear `column-toggle.ts`** con una función reutilizable tipo:
   ```ts
   interface ColumnConfig {
     key: string;        // identificador único de la columna
     label: string;      // texto del <th>
     fixed?: boolean;    // true = no se puede ocultar (columna de nombre)
   }

   function initColumnToggle(tableId: string, columns: ColumnConfig[]): void
   ```
   - Renderiza un icono de ojo (SVG o icono existente en el proyecto) dentro de cada `<th>` que no sea `fixed`.
   - Al hacer click en el ojo, alterna la visibilidad de esa columna (th + todos los td correspondientes).
   - Usa clases CSS (`col-hidden`) en las celdas para ocultar/mostrar con `display: none`.
   - El ojo cambia de estilo cuando la columna está oculta (ojo tachado / opacidad reducida).
   - Lee/escribe el estado en `localStorage` con clave `column-visibility-{tableId}`.

2. **Mecanismo de ocultación**: Asignar `data-col="key"` a cada `<th>` y `<td>` correspondiente. Al toggle, añadir/quitar clase `.col-hidden` a todas las celdas con ese `data-col`.

3. **Persistencia en localStorage**:
   - Clave: `column-visibility-{tableId}` (ej: `column-visibility-pokedex-stats`).
   - Valor: JSON con las columnas ocultas: `["sprite", "types"]`.
   - Al inicializar la tabla, leer localStorage y aplicar el estado guardado.

4. **Integrar en cada tabla**:
   - Modificar la generación del HTML de cada tabla para incluir `data-col` en th/td.
   - Llamar a `initColumnToggle(tableId, columns)` después de renderizar.
   - Asegurar que al re-renderizar (por sort, filtro, etc.) se preserva el estado de visibilidad.

5. **Icono de ojo**:
   - Usar SVG inline (ojo abierto = visible, ojo tachado = oculto).
   - Posicionar junto al texto del header, después del sort-indicator si existe.
   - Tamaño pequeño (14-16px), color sutil que no compita con el texto.

6. **Estilos**:
   - `.col-toggle-icon`: cursor pointer, opacity 0.4 por defecto, opacity 1 en hover.
   - `.col-toggle-icon.hidden`: ojo tachado, opacity 0.3.
   - `.col-hidden`: `display: none`.
   - Dark mode: ajustar colores del icono.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Todas las tablas muestran icono de ojo en las cabeceras ocultables |
| (visual) | La columna fija (nombre) NO tiene icono de ojo |
| (visual) | Click en ojo oculta la columna completa (th + todos los td) |
| (visual) | Click en ojo oculto restaura la columna |
| (visual) | El icono cambia de apariencia al ocultar (ojo tachado) |
| (visual) | Recargar la página mantiene las columnas ocultas (localStorage) |
| (visual) | Ordenar/filtrar una tabla preserva las columnas ocultas |
| (visual) | Funciona en dark mode |
| (visual) | Funciona en mobile (responsive) |
| (visual) | La tabla de stats config (IV/EV) permite ocultar IV o EV |
| (visual) | Las dos tablas de damage (attacker/defender) tienen toggles independientes |

## Criterios de aceptacion

- [x] Las 7 tablas tienen iconos de ojo en las columnas ocultables
- [x] La columna de nombre/identificador principal no se puede ocultar en ninguna tabla
- [x] Click en el ojo oculta toda la columna (header + celdas)
- [x] Click en el ojo oculto restaura la columna
- [x] El icono indica visualmente si la columna está oculta o visible
- [x] El estado de visibilidad se persiste en localStorage por tabla
- [x] Al recargar la app se restaura el estado guardado
- [x] Re-renderizados (sort, filtros) no resetean la visibilidad
- [x] Compatible con dark mode
- [x] Responsive: funciona en mobile
- [x] La lógica está centralizada en un módulo reutilizable (`column-toggle.ts`)

## Notas

- Reutilizar el patrón de iconos SVG inline que ya se usa en el proyecto (verificar cómo se implementan los iconos de categoría en moves).
- El `data-col` attribute es la forma más limpia de vincular th con td sin depender del índice de columna, que puede cambiar.
- Para tablas que se re-renderizan completas (como encounters tras sort), hay que re-aplicar `initColumnToggle` o asegurar que el re-render incluya los `data-col` y clases correctas.
