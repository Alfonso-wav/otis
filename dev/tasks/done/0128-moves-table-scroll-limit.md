# Task 0128 — Pokedex: limitar tabla de movimientos a 50 filas con scroll

## Estado: pending

## Goal

En la vista individual de un Pokemon (Pokedex > detalle), la tabla de movimientos se despliega completamente sin importar cuántas filas tenga. En Pokemon con muchos movimientos (100+), la tabla se hace extremadamente larga y empuja todo el contenido hacia abajo. Limitar la altura visible de la tabla a ~50 filas y añadir scroll vertical para el resto.

---

## Contexto técnico

### Archivos principales

- `frontend/src/pages/pokedex.ts` — función `renderMoves()` (líneas 989–1138), genera la tabla de movimientos
- `frontend/src/styles/_pokemon.scss` — estilos de la tabla de movimientos (`.moves-table-wrap`, `.moves-table`)

### Estado actual

La función `renderTable()` en `pokedex.ts` (línea 1029) genera **todas** las filas sin limitación. El wrapper `.moves-table-wrap` (línea 1069) no tiene `max-height` ni `overflow-y`, así que la tabla crece indefinidamente.

### Solución propuesta

Añadir `max-height` y `overflow-y: auto` al contenedor `.moves-table-wrap` para que la tabla se contenga dentro de un área scrollable. La cabecera (`thead`) debería quedar sticky para que sea visible mientras se hace scroll.

Estimación de altura: cada fila de la tabla tiene ~36-40px. Para ~50 filas: `50 × 38px ≈ 1900px`. Un `max-height` de **1900px** es razonable.

---

## Cambios requeridos

### 1. Añadir scroll al wrapper — `_pokemon.scss`

```scss
.moves-table-wrap {
  max-height: 1900px;
  overflow-y: auto;
}
```

### 2. Hacer thead sticky

```scss
.moves-table thead th {
  position: sticky;
  top: 0;
  z-index: 1;
  background: inherit;  // hereda el fondo de la tabla para no ser transparente
}
```

### 3. Verificar dark mode

Asegurar que el `thead` sticky tenga fondo correcto en dark mode (si `_dark.scss` ya define colores para `.moves-table thead`, el `background: inherit` debería propagarlo).

---

## Archivos afectados

### Frontend
- `frontend/src/styles/_pokemon.scss` — `max-height` + `overflow-y` en `.moves-table-wrap`, sticky thead

### Backend
Ninguno.

---

## Acceptance criteria

- [ ] La tabla de movimientos no excede ~1900px de alto (equivalente a ~50 filas).
- [ ] Si hay más de 50 filas, aparece scroll vertical dentro del contenedor.
- [ ] La cabecera de la tabla (nombre, nivel, método) permanece visible al hacer scroll.
- [ ] Los filtros (All, Level-up, TM, Egg, Tutor) siguen funcionando correctamente.
- [ ] El ordenamiento por columnas sigue funcionando.
- [ ] Correcto en dark mode.
- [ ] Probado con Pokemon que tenga muchos movimientos (ej: Mew, Smeargle).
