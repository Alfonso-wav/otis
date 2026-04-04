# Pokedex tabla: control de límite de filas y fix de headers en móvil

**ID**: 0101-pokedex-table-row-limit-and-mobile-fix
**Estado**: done
**Fecha**: 2026-04-04

---

## Descripcion

Dos mejoras para la vista de tabla de la Pokédex:

1. **Control de límite de filas**: Añadir un control en la parte inferior de la tabla que permita al usuario elegir cuántas filas mostrar. Por defecto 50, ajustable de 50 en 50 (50, 100, 150, 200...). Reemplaza la paginación actual (LIMIT=20 con botones Prev/Next) por este sistema de límite configurable.

2. **Fix de headers en móvil**: Los títulos de las columnas se solapan y son ilegibles en pantallas pequeñas. La tabla debe garantizar que todos los headers sean legibles. Si no hay espacio suficiente, la tabla debe expandirse hacia la derecha permitiendo scroll horizontal para descubrir el resto de columnas, sin comprimir ni solapar los títulos.

## Capas afectadas

- **Core**: No aplica
- **Shell**: No aplica
- **APP**: Frontend — componente de tabla en `pokedex.ts` y estilos en `_pokemon.scss`

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Cambiar constante LIMIT de 20 a 50, reemplazar paginación Prev/Next por control de límite de filas (+50/-50), cargar filas según el límite seleccionado |
| `frontend/src/styles/_pokemon.scss` | modificar | Asegurar `min-width` en columnas de la tabla para que los headers no se solapen; ajustar estilos del contenedor para scroll horizontal correcto; eliminar sticky columns si interfieren con la legibilidad |

## Plan de implementacion

1. En `pokedex.ts`, cambiar `LIMIT = 20` a un valor inicial de 50 y hacerlo variable (e.g. `let rowLimit = 50`)
2. Reemplazar los controles de paginación (Prev/Next/Page info) por un control de "Límite de filas" con botones para decrementar (-50) e incrementar (+50), mostrando el valor actual
3. Al cambiar el límite, recargar la tabla mostrando las primeras N filas según el nuevo límite
4. En `_pokemon.scss`, añadir `min-width` a las celdas de la tabla (especialmente headers) para evitar compresión
5. Cambiar la tabla a `width: max-content` (o `min-width: 100%`) en tablet/mobile para que se expanda en lugar de comprimir columnas
6. Revisar y ajustar los valores `left` de las sticky columns o eliminar el sticky si causa solapamiento
7. Probar en viewport móvil que los headers sean todos legibles y el scroll horizontal funcione correctamente

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que el control de límite muestra 50 por defecto y permite cambiar de 50 en 50 |
| Manual | Verificar que en móvil los headers no se solapan y el scroll horizontal funciona |

## Criterios de aceptacion

- [x] La tabla muestra 50 filas por defecto al cargar
- [x] Existe un control inferior que muestra el límite actual y permite cambiarlo de 50 en 50
- [x] El botón de -50 se deshabilita cuando el límite es 50 (mínimo)
- [x] Al cambiar el límite la tabla se actualiza mostrando el número correcto de filas
- [x] En pantallas móviles, todos los títulos de columna son legibles sin solapamiento
- [x] La tabla permite scroll horizontal en móvil para ver todas las columnas
- [x] No hay regresiones en la vista de escritorio

## Notas

- La paginación actual usa `LIMIT = 20` con offset y botones Prev/Next (líneas ~20 y ~526-571 de pokedex.ts)
- Los estilos de la tabla están en `_pokemon.scss` líneas 127-206, con sticky columns en tablet/mobile
- Los valores hardcodeados de `left` en sticky columns (0, 2.5rem, 5.5rem) pueden ser parte del problema de solapamiento
- El grid view tiene infinite scroll separado que no debe verse afectado
