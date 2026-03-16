# Infinite Scroll en la vista de tarjetas de la Pokédex

**ID**: 0089-pokedex-infinite-scroll
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Reemplazar la paginación actual (botones prev/next con 20 tarjetas por página) por infinite scroll en la vista de grid de la Pokédex. Las tarjetas se cargarán progresivamente conforme el usuario haga scroll hacia abajo, dando una experiencia más fluida y continua.

## Capas afectadas

- **APP (frontend)**: `pokedex.ts` — implementar lógica de infinite scroll con Intersection Observer, acumular tarjetas en el grid en lugar de reemplazarlas por página.
- **Frontend styles**: `_pokemon.scss` — eliminar/ocultar controles de paginación en modo grid, añadir indicador de carga al fondo del grid.
- **Backend/Core/Shell**: sin cambios. Se sigue usando `ListPokemon(offset, limit)` con los mismos parámetros.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Implementar infinite scroll con IntersectionObserver, acumular cards en grid, gestionar estado de carga |
| `frontend/src/styles/_pokemon.scss` | modificar | Añadir estilos para el sentinel/loader al fondo del grid |
| `frontend/index.html` | modificar | Añadir elemento sentinel para el IntersectionObserver al final del grid |
| `frontend/src/animations/transitions.ts` | modificar (si necesario) | Asegurar que `staggerCards` funcione con lotes incrementales (animar solo las cards nuevas, no las ya visibles) |

## Plan de implementacion

### Parte 1: Elemento sentinel e IntersectionObserver

1. Añadir un elemento sentinel (e.g., `<div id="scroll-sentinel"></div>`) después del `#pokemon-grid` en `index.html`.
2. Crear un `IntersectionObserver` que observe el sentinel. Cuando el sentinel entre en el viewport, disparar la carga del siguiente lote.
3. Incluir un indicador de carga (spinner o texto "Cargando...") visible mientras se cargan más tarjetas.

### Parte 2: Acumular tarjetas en lugar de reemplazar

4. Modificar `renderGrid()` para que en modo infinite scroll **añada** (append) las nuevas tarjetas al grid existente en lugar de reemplazar todo el `innerHTML`.
5. Mantener el offset como estado del módulo e incrementarlo con cada lote cargado.
6. Cuando `offset >= totalCount`, desconectar el observer y ocultar el sentinel/loader (se han cargado todos los Pokémon).

### Parte 3: Animación de los nuevos lotes

7. Al añadir un nuevo lote, aplicar `staggerCards` solo a las tarjetas recién añadidas para que tengan la animación de entrada sin re-animar las existentes.
8. Si `staggerCards` actualmente anima todos los hijos del contenedor, modificarlo para aceptar un subset de elementos o un offset de inicio.

### Parte 4: Integración con filtros y vistas

9. Cuando se aplican filtros, resetear el grid (vaciar tarjetas, resetear offset, reconectar observer) y cargar desde el inicio con los filtros activos.
10. `filteredList` ya almacena la lista filtrada completa. En modo filtered, paginar sobre `filteredList` en lugar de llamar a la API.
11. Cuando se cambia a vista tabla, desconectar el observer (la tabla mantiene su paginación actual). Al volver a grid, reconectar.

### Parte 5: Ocultar paginación en modo grid

12. Ocultar los controles de paginación (`#pagination`) cuando se está en vista grid con infinite scroll.
13. Mantener la paginación funcional para la vista tabla.

### Parte 6: Scroll to top

14. Opcionalmente, añadir un botón "volver arriba" que aparezca cuando el usuario haya scrolleado más allá de cierto umbral.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Al hacer scroll hacia abajo se cargan más tarjetas automáticamente |
| (visual) | Se muestra un indicador de carga mientras se obtienen las tarjetas |
| (visual) | Las tarjetas nuevas aparecen con animación stagger |
| (visual) | Al llegar al final de la Pokédex no se intentan más cargas |
| (visual) | Al aplicar un filtro se resetea el grid y se carga desde el inicio |
| (visual) | Al cambiar a vista tabla y volver a grid, el infinite scroll funciona correctamente |
| (visual) | Los controles de paginación no aparecen en vista grid |
| (visual) | Los controles de paginación siguen funcionando en vista tabla |
| (visual) | Funciona en dark mode |
| (visual) | Funciona en mobile (responsive) |

## Criterios de aceptacion

- [ ] Las tarjetas de la Pokédex se cargan progresivamente al hacer scroll (lotes de 20)
- [ ] Se muestra un indicador de carga al fondo del grid mientras se obtienen más tarjetas
- [ ] Las tarjetas nuevas entran con animación stagger
- [ ] Al alcanzar el final de la lista, el loading desaparece y no hay más peticiones
- [ ] Los filtros (generación, tipo, legendario, mítico) funcionan correctamente con infinite scroll
- [ ] Cambiar entre vista grid y tabla no rompe el comportamiento
- [ ] La paginación se oculta en grid y se mantiene en tabla
- [ ] Compatible con dark mode y responsive

## Notas

- La API `ListPokemon(offset, limit)` ya soporta paginación por offset, así que no requiere cambios en backend.
- El `LIMIT = 20` actual se puede reutilizar como tamaño de lote para infinite scroll.
- `IntersectionObserver` tiene soporte universal en navegadores modernos, no requiere polyfill.
- En modo filtrado, los datos ya están en memoria (`filteredList`), así que el scroll infinito pagina sobre el array local sin llamadas extra a la API.
- La vista tabla mantiene paginación clásica porque necesita cargar datos completos (stats) de cada Pokémon y el scroll infinito sería más costoso.
