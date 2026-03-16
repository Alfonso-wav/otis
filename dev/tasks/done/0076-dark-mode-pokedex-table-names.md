# 0076 — Fix: nombres de Pokemon apenas visibles en la tabla Pokedex en modo oscuro

**Estado: done**

## Descripcion

En la vista de tabla de la Pokedex (`.poke-table`), los nombres de los Pokemon (`.poke-table__name`) usan el color `#2d3748` (gris oscuro) definido en `_pokemon.scss`. En modo oscuro, el fondo de la tabla es `#1a202c` o `#2d3748`, lo que hace que los nombres sean practicamente invisibles por falta de contraste.

**Comportamiento esperado**: en modo oscuro, los nombres de los Pokemon deben mostrarse en color claro (blanco o similar) para ser legibles.

## Capas afectadas

- **APP (frontend)** — agregar override de color en los estilos dark mode.

No se requieren cambios en Core ni Shell.

## Archivos a modificar

### Frontend

1. **`frontend/src/styles/_dark.scss`**

- **Linea ~27-41**: dentro del bloque `.poke-table` del dark mode, agregar un override para `.poke-table__name`:
  ```scss
  .poke-table__name {
    color: #ffffff;
  }
  ```
- Esto sigue el patron ya existente en `_dark.scss` donde se sobreescriben colores de texto para dark mode (ej. `th { color: #a0aec0; }`).

## Requisitos

1. Los nombres de Pokemon en la tabla Pokedex deben ser claramente legibles en modo oscuro.
2. El color en modo claro (`#2d3748`) no debe verse afectado.
3. El cambio debe aplicarse solo dentro del selector dark mode (`[data-bs-theme="dark"]`).

## Tests

- Activar modo oscuro en Settings, ir a Pokedex, cambiar a vista de tabla y verificar que los nombres se leen con claridad.
- Verificar que en modo claro los nombres siguen mostrando el color original.
