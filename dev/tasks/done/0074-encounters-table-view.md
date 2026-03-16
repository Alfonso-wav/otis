# 0074 — Convertir la vista de encuentros de tarjetas a tabla unificada

**Estado: done**

## Descripcion

Reemplazar el layout actual de la seccion Encounters en la vista de detalle del Pokemon (tarjetas individuales por location area) por una tabla unica condensada usando el mismo estilo `.poke-table` que las tablas de Moves y Abilities en Explore.

## Capas afectadas

- **APP (frontend)** — refactorizar el renderizado HTML de `loadEncounters()` en `pokedex.ts`
- **Estilos** — reemplazar los estilos de tarjetas por estilos de tabla en `_pokemon.scss` y `_dark.scss`

No se requieren cambios en Core ni Shell. La estructura de datos del backend ya es adecuada.

## Archivos a modificar

### Frontend

1. **`frontend/src/pages/pokedex.ts`** (lineas ~568-641) — Refactorizar `loadEncounters()`:
   - Reemplazar la generacion de HTML basada en `.encounter-location` cards por una tabla `<table class="poke-table encounters-table">`.
   - Columnas sugeridas: **Location**, **Game**, **Method**, **Chance**, **Levels**, **Conditions**.
   - Cada fila representa un detalle de encuentro individual (una combinacion location + version + method).
   - Envolver la tabla en un div `.encounters-table-wrap` con el mismo patron de overflow que `.moves-table-wrap`.

### Estilos

2. **`frontend/src/styles/_pokemon.scss`** (lineas ~494-625) — Reemplazar los estilos actuales de encounters:
   - Eliminar las clases `.encounter-location`, `.encounter-versions`, `.encounter-version`, `.encounter-method`, `.encounter-method-name`, `.encounter-chance`, `.encounter-levels`, `.encounter-conditions`.
   - Anadir `.encounters-table-wrap` siguiendo el patron de `.moves-table-wrap` (overflow-x auto, border-radius, box-shadow, background white).
   - Anadir `.encounters-table` con estilos especificos si se necesitan celdas custom.

3. **`frontend/src/styles/_dark.scss`** (lineas ~593-624) — Actualizar overrides dark mode:
   - Reemplazar los overrides de las clases de tarjeta eliminadas por overrides de `.encounters-table-wrap` siguiendo el patron de `.moves-table-wrap` en dark mode (background `#1a202c`).

## Requisitos

1. La tabla debe usar la clase `.poke-table` como base, igual que Moves y Abilities.
2. Columnas: Location | Game | Method | Chance (%) | Levels | Conditions.
3. Mantener la carga asincrona existente y el mensaje "not found in the wild" cuando no hay datos.
4. El wrapper debe tener scroll horizontal en movil (patron `.encounters-table-wrap` con `overflow-x: auto`).
5. Compatible con dark mode usando los mismos overrides que las tablas existentes.
6. Responsive: legible en desktop y movil.
7. Mantener el spinner de carga existente.

## Referencia de estilo (tablas existentes)

```html
<div class="encounters-table-wrap">
  <table class="poke-table encounters-table">
    <thead>
      <tr>
        <th>Location</th>
        <th>Game</th>
        <th>Method</th>
        <th>Chance</th>
        <th>Levels</th>
        <th>Conditions</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>viridian-forest-area</td>
        <td>red</td>
        <td>walk</td>
        <td>15%</td>
        <td>3-5</td>
        <td>—</td>
      </tr>
    </tbody>
  </table>
</div>
```

## Tests

- Verificar que la tabla se renderiza correctamente con datos de encuentros.
- Verificar que Pokemon sin encuentros siguen mostrando "This Pokemon is not found in the wild."
- Verificar scroll horizontal en pantallas pequenas.
- Verificar dark mode.
- Verificar que la carga asincrona y el spinner funcionan como antes.
