# Cambiar botones de simulación de morado a rojo

**ID**: 0078-simulate-buttons-red-color
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Los botones de simulación en la pestaña Build ("Simular batalla completa", "Simular batalla de equipos" y "Simular N batallas") usan color morado (#6b46c1) que no encaja con el estilo visual de la app. Deben usar el rojo corporativo ($pokedex-red: #e53e3e) como el resto de botones.

## Capas afectadas

- **APP (frontend)**: Solo cambios de estilo CSS/SCSS.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_builds.scss` | modificar | Reemplazar `#6b46c1` por `$pokedex-red` en `.battle-auto-btn` (líneas ~491-506), `.battle-batch-btn` (líneas ~1048-1063) y `.battle-batch-input:focus` (línea ~1045) |

## Plan de implementacion

1. En `_builds.scss`, clase `.battle-auto-btn`: cambiar `background: #6b46c1` → `background: $pokedex-red` y `darken(#6b46c1, 8%)` → `darken($pokedex-red, 8%)` en hover.
2. En `_builds.scss`, clase `.battle-batch-btn`: mismo cambio de background y hover.
3. En `_builds.scss`, clase `.battle-batch-input:focus`: cambiar `border-color: #6b46c1` → `border-color: $pokedex-red`.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Verificar que los 4 botones de simulación (batalla, equipos, batch x2) se muestran en rojo |

## Criterios de aceptacion

- [x] Los botones `.battle-auto-btn` se muestran en rojo (#e53e3e)
- [x] Los botones `.battle-batch-btn` se muestran en rojo (#e53e3e)
- [x] El hover oscurece el rojo correctamente
- [x] El input de batch tiene borde rojo al hacer focus
- [x] No hay restos de color morado (#6b46c1) en la pestaña Build

## Notas

- El color morado `#6b46c1` está hardcodeado 5 veces en `_builds.scss`. Reemplazar todas las ocurrencias por la variable `$pokedex-red` para mantener consistencia.
- Cambio puramente visual, sin impacto en lógica.
