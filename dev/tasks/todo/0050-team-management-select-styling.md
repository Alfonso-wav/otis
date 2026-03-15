# Corregir estilos de selects en gestión de equipos

**ID**: 0050-team-management-select-styling
**Estado**: todo
**Fecha**: 2026-03-15

---

## Descripcion

Varios cuadros de selección (`<select>`) en la sección de gestión de equipos y batalla por equipos no coinciden con el estilo visual de la aplicación. Se necesita unificar su apariencia con el resto de la app.

Los selects afectados:
- `.team-battle-select` (selectores de equipo 1 y equipo 2 en batalla) — **no tiene estilos CSS definidos**
- `.build-move-select` (selector de movimientos en builds) — tiene estilos pero verificar consistencia
- `.sc-select` / `.sc-nature` (selector de naturaleza) — verificar consistencia

## Capas afectadas

- **Core**: No afectada.
- **Shell**: No afectada.
- **APP**: Solo frontend (CSS/SCSS).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_builds.scss` | modificar | Añadir estilos para `.team-battle-select` y unificar con el resto de selects |
| `frontend/src/pages/builds.ts` | revisar | Verificar que las clases CSS son consistentes en todos los selects |

## Plan de implementacion

1. Revisar todos los `<select>` en `builds.ts` y listar las clases CSS que usan.
2. Definir estilos para `.team-battle-select` que coincidan con el patrón visual existente (`.build-move-select`, `.sc-select`).
3. Verificar que todos los selects comparten: border, border-radius, padding, font-family, colores de fondo y focus.
4. Aplicar los mismos estilos de tema Game Boy si aplica.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| N/A | Verificación visual — los selects deben verse consistentes en toda la sección de builds |

## Criterios de aceptacion

- [ ] `.team-battle-select` tiene estilos definidos que coinciden con el resto de la app
- [ ] Todos los selects en la sección de builds tienen apariencia visual consistente
- [ ] Los selects tienen estados `:focus` y `:hover` coherentes
- [ ] Se ve bien en pantallas pequeñas (responsive)

## Notas

Actualmente `.team-battle-select` no tiene ningún estilo CSS definido en `_builds.scss`. Los otros selects (`.build-move-select`, `.sc-select`) sí tienen estilos. Usar esos como referencia para el nuevo estilo.
