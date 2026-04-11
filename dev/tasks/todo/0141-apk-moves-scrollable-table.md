# Tabla de movimientos scrollable en APK (como encuentros)

**ID**: 0141-apk-moves-scrollable-table
**Estado**: todo
**Fecha**: 2026-04-11

---

## Descripcion

En la version APK (Android), la tabla de movimientos en la vista individual del Pokemon debe tener el mismo patron de scroll que la tabla de encuentros: un contenedor con scroll vertical que muestre unos 25 movimientos visibles y scroll lateral a la derecha para ver todas las columnas, sin limite de filas. Actualmente la tabla de movimientos esta limitada a 50 filas y el scroll puede no funcionar bien en movil.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — estilos y logica de la tabla de movimientos

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Eliminar limite de 50 filas en movil, mostrar todos los movimientos |
| `frontend/src/styles/_pokemon.scss` | modificar | Ajustar max-height de .moves-table-wrap para mostrar ~25 filas en movil |

## Plan de implementacion

1. En `pokedex.ts`, detectar si estamos en la version movil/APK y no aplicar el limite de 50 filas (o eliminar el limite por completo y dejar que el scroll maneje la visualizacion).
2. En `_pokemon.scss`, ajustar el `max-height` de `.moves-table-wrap` en mobile para que sea equivalente a ~25 filas visibles (~900px-1000px).
3. Asegurar que `overflow-y: auto`, `overflow-x: auto` y `-webkit-overflow-scrolling: touch` estan presentes.
4. Verificar que el scroll lateral funciona correctamente en Android.
5. Probar con Pokemon que tengan muchos movimientos (ej: Mew, Smeargle).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual (APK) | Verificar scroll vertical fluido con todos los movimientos visibles |
| Manual (APK) | Verificar scroll lateral para columnas que no caben |
| Manual (APK) | Verificar que no hay limite artificial de filas |
| Manual (Desktop) | Verificar que la version desktop no se rompe |

## Criterios de aceptacion

- [ ] En APK, la tabla de movimientos muestra todos los movimientos sin limite
- [ ] El scroll vertical muestra ~25 filas y permite ver el resto haciendo scroll
- [ ] El scroll lateral funciona correctamente en Android
- [ ] La tabla tiene el mismo patron visual que la de encuentros
- [ ] La version desktop sigue funcionando correctamente

## Notas

La tabla de encuentros usa `.encounters-table-wrap` con `max-height: 900px`, `overflow-y: auto`, `overflow-x: auto`, `-webkit-overflow-scrolling: touch`. La tabla de movimientos ya tiene un patron similar pero con `max-height: 1750px` y un limite de 50 filas en JS. Hay que unificar el enfoque: sin limite de filas en JS, y max-height mas reducido en movil (~25 filas visibles).
