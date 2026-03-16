# 0072 — Fix dark mode en tablas de Explorar (movimientos y habilidades)

## Descripción

En dark mode, las tablas de Explorar → Movimientos y Explorar → Habilidades se muestran con fondo blanco y texto gris, lo que dificulta la lectura. Hay que aplicar los mismos estilos oscuros que ya tiene la tabla de la pestaña Pokédex.

## Capas afectadas

- **APP** — frontend (`frontend/src/`)

## Archivos a modificar

1. **`frontend/src/styles/_dark.scss`** — Añadir overrides dark mode para `.moves-table-wrap`, `.abilities-table-wrap` y `.ability-desc-cell`.
2. **`frontend/src/styles/_explore.scss`** — (Opcional) Refactorizar los `background: white` hardcodeados si se prefiere usar variables de Bootstrap.

## Problema raíz

- `.moves-table-wrap` (línea 309 de `_explore.scss`) tiene `background: white` sin override en dark mode.
- `.abilities-table-wrap` (línea 374 de `_explore.scss`) tiene `background: white` sin override en dark mode.
- `.ability-desc-cell` (línea 388 de `_explore.scss`) tiene `color: #4a5568` (gris oscuro) que en dark mode es ilegible sobre fondo oscuro.
- La tabla Pokédex NO tiene este problema porque su contenedor `#pokemon-grid` no tiene background hardcodeado — usa el background de la página, que Bootstrap adapta automáticamente en dark mode.

## Requisitos

1. En dark mode, `.moves-table-wrap` y `.abilities-table-wrap` deben tener `background: #1a202c` (consistente con el esquema dark del proyecto).
2. En dark mode, `.ability-desc-cell` debe tener `color: #a0aec0` para ser legible.
3. Las tablas de movimientos y habilidades deben verse visualmente iguales a la tabla de la Pokédex en dark mode.
4. No alterar la apariencia en light mode.

## Enfoque sugerido

Añadir en `_dark.scss` dentro del bloque `[data-bs-theme="dark"]`:

```scss
// ─── Explore table wrappers ────────────────────────────────────────────
.moves-table-wrap,
.abilities-table-wrap {
  background: #1a202c;
}

.abilities-table .ability-desc-cell {
  color: #a0aec0;
}
```

## Tests

- Verificar en dark mode que las tablas de movimientos y habilidades tienen fondo oscuro.
- Verificar que el texto es legible (contraste adecuado) en ambas tablas.
- Verificar que la descripción de habilidades se lee correctamente.
- Comparar visualmente con la tabla de la Pokédex — deben tener apariencia consistente.
- Verificar que en light mode no hay cambios.
