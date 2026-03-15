# Modal con lista de Pokémon por habilidad

**ID**: 0054-ability-pokemon-list-modal
**Estado**: todo
**Fecha**: 2026-03-15
**Depende de**: —

---

## Descripción

En la pestaña Explorar > Habilidades, la tabla muestra en la columna "Pokémon" el número de Pokémon que pueden aprender cada habilidad. Actualmente ese número es solo texto estático.

Se quiere que al hacer clic en el número de Pokémon de cualquier fila, se abra una ventana modal que muestre una tabla con los datos de esos Pokémon (nombre, sprite, tipos, etc.).

## Capas afectadas

- **Core**: No afectada (el struct `Ability` ya contiene `Pokemon []string`).
- **Shell**: No afectada (los datos ya se obtienen de PokeAPI).
- **APP (Frontend)**: Modificar la tabla de habilidades para que el número sea clickeable y abrir un modal con los datos de los Pokémon. Reutilizar el patrón de modal existente en `pokemon-type-modal.ts`.

## Archivos a crear/modificar

| Archivo | Acción | Descripción |
|---------|--------|-------------|
| `frontend/src/pages/explore/abilities.ts` | modificar | 1) Hacer que la celda del conteo de Pokémon sea clickeable (clase CSS + cursor pointer). 2) Al hacer clic, abrir modal con la lista de Pokémon de esa habilidad. |
| `frontend/src/components/ability-pokemon-modal.ts` | crear | Componente modal que recibe un nombre de habilidad y una lista de nombres de Pokémon, y muestra una tabla/grid con sus datos (sprite, nombre, tipos). Seguir el patrón de `pokemon-type-modal.ts`. |
| `frontend/src/styles/_explore.scss` | modificar | Estilos para el modal de Pokémon por habilidad y para la celda clickeable en la tabla de habilidades. Reutilizar los estilos existentes de `.type-modal-*` o extenderlos. |

## Plan de implementación

### Fase 1 — Componente modal

1. Crear `ability-pokemon-modal.ts` siguiendo el patrón de `pokemon-type-modal.ts`:
   - Overlay con fondo semitransparente.
   - Header con nombre de la habilidad y botón de cierre.
   - Body con grid/tabla de Pokémon mostrando: sprite (desde pokemondb local o CDN fallback), nombre capitalizado y tipos con badges de color.
   - Cerrar con clic en overlay, botón X, o tecla Escape.
   - Limpieza de event listeners al cerrar.

### Fase 2 — Integración en tabla de habilidades

1. En `abilities.ts`, cambiar la celda `<td class="num-cell">` del conteo a un elemento clickeable (e.g., `<span class="ability-pokemon-count" data-ability="...">N</span>`).
2. Añadir event listener delegado en la tabla que al clic en `.ability-pokemon-count` abra el modal pasando el nombre de la habilidad y su array de Pokémon.

### Fase 3 — Estilos

1. Añadir estilo `.ability-pokemon-count` con `cursor: pointer`, `text-decoration: underline`, color de enlace, y hover state.
2. Reutilizar o extender los estilos `.type-modal-*` existentes para el nuevo modal.

## Tests

- Verificar que al hacer clic en el número se abre el modal.
- Verificar que el modal muestra todos los Pokémon de la habilidad.
- Verificar que se cierra con X, Escape y clic fuera.
- Verificar que funciona correctamente con habilidades que tienen 0 Pokémon (no abrir modal o mostrar mensaje vacío).

## Criterios de aceptación

- [ ] El número de Pokémon en la tabla de habilidades es clickeable y visualmente distinguible.
- [ ] Al hacer clic se muestra un modal con los datos de los Pokémon que aprenden esa habilidad.
- [ ] El modal muestra sprite, nombre y tipos de cada Pokémon.
- [ ] El modal se cierra con X, Escape o clic fuera del modal.
- [ ] Si la habilidad tiene 0 Pokémon, no se abre el modal o se muestra un mensaje indicativo.
- [ ] El diseño es responsive y coherente con los modales existentes del proyecto.
