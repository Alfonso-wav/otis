# Corregir visibilidad de Mr. Mime en overlay de sorting cuando el contenedor es largo

**ID**: 0094-fix-sorting-overlay-mrmime-visibility
**Estado**: done
**Fecha**: 2026-03-17

---

## Descripcion

Cuando el usuario ordena una tabla larga (por ejemplo, la Pokedex con todos los Pokemon), el overlay de Mr. Mime aparece centrado verticalmente dentro del contenedor completo. Si el contenedor es mas alto que el viewport, Mr. Mime queda fuera de la zona visible y el usuario tiene que hacer scroll hacia abajo para verlo.

Mr. Mime debe aparecer siempre visible en la pantalla sin necesidad de scroll.

## Capas afectadas

- **Core**: No afectada.
- **Shell**: No afectada.
- **APP**: Frontend — ajustar CSS del overlay de sorting para que Mr. Mime sea visible en el viewport.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_components.scss` | modificar | Cambiar posicionamiento del `.sorting-overlay` para que Mr. Mime quede visible en el viewport (usar `position: fixed` o `align-items: flex-start` con padding adecuado) |
| `frontend/src/components/sorting-overlay.ts` | modificar | Si es necesario, ajustar la logica de insercion del overlay (por ejemplo, usar `document.body` en vez del contenedor si se cambia a `position: fixed`) |

## Plan de implementacion

1. Cambiar `.sorting-overlay` de `position: absolute` a `position: fixed` para que se centre respecto al viewport, no al contenedor.
2. Si se usa `position: fixed`, ajustar `inset: 0` y el `z-index` para que cubra toda la pantalla correctamente.
3. Verificar que el fondo semitransparente sigue cubriendo el contenedor de la tabla.
4. Alternativa: mantener `position: absolute` pero usar `align-items: flex-start` con `padding-top` calculado para posicionar Mr. Mime en la zona visible del contenedor.
5. Probar en las tres tablas: Pokedex, Moves, Abilities.
6. Probar en modo claro y oscuro.
7. Probar en mobile y desktop.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que Mr. Mime es visible sin scroll al ordenar tablas largas (Pokedex) |
| Manual | Verificar que el overlay sigue cubriendo la tabla correctamente |
| Manual | Verificar en mobile y desktop, modo claro y oscuro |

## Criterios de aceptacion

- [x] Mr. Mime aparece visible en la pantalla sin necesidad de hacer scroll al ordenar cualquier tabla.
- [x] El fondo semitransparente cubre correctamente el area de la tabla.
- [x] Funciona en las tres tablas: Pokedex, Moves, Abilities.
- [x] Se ve correctamente en modo claro y oscuro.
- [x] Se ve correctamente en mobile y desktop.

## Notas

- El overlay actual usa `position: absolute` + `inset: 0` + flexbox centering dentro del contenedor. El problema es que en contenedores mas altos que el viewport, el centro queda fuera de la vista.
- La solucion mas simple es usar `position: fixed` para centrar respecto al viewport. Hay que verificar que el fondo siga cubriendo el area esperada.
