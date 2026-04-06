# Task 0114 — Comparador: sin emoji, comparaciones apilables y colores únicos

## Estado: done

## Goal
El comparador de estadísticas en la vista de detalle tiene tres problemas:
1. **Emoji prohibido**: el botón de comparar usa ⚖️, en contra de la norma del proyecto (solo sprites o nada).
2. **Solo compara dos Pokémon**: al añadir un segundo comparado se reemplaza el anterior. El usuario quiere poder apilar comparaciones ad infinitum.
3. **Color único por comparado**: cada Pokémon añadido debe tener un color diferente y persistente en la leyenda y en la gráfica.

## Acceptance criteria
- [x] El botón "Comparar" no tiene emoji
- [x] Se pueden añadir 3 o más Pokémon a la comparación
- [x] Cada Pokémon comparado tiene un color único y diferente al del base
- [x] La leyenda muestra todos los comparados con su color
- [x] Se puede eliminar un comparado individualmente con su "×" en la leyenda
- [x] Al eliminar un comparado, la gráfica se actualiza correctamente
- [x] El Pokémon base nunca se puede eliminar del comparador
- [x] `renderStatsChart` acepta un array de series (refactor compatible)
