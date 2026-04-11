# Quitar fondo de tarjeta de las habilidades en detalle Pokemon

**ID**: 0142-detail-abilities-no-cards
**Estado**: done
**Fecha**: 2026-04-11

---

## Descripcion

En la vista individual de Pokemon, seccion Habilidades, cada habilidad se muestra dentro de una tarjeta con fondo (#f7fafc en light, #2d3748 en dark), borde izquierdo rojo y border-radius. El usuario quiere eliminar estos fondos de tarjeta para que la informacion se muestre directamente sobre el fondo de la pagina, sin contenedor visual.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — estilos de habilidades

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_pokemon.scss` | modificar | Eliminar background, box-shadow, border-left y border-radius de `.ability-card` |
| `frontend/src/styles/_dark.scss` | modificar | Eliminar override de background dark para `.ability-card` si existe |

## Plan de implementacion

1. En `_pokemon.scss`, localizar los estilos de `.ability-card` (lineas ~873-935).
2. Eliminar: `background`, `border-left`, `border-radius`, `box-shadow`.
3. Mantener: padding (para espaciado), estilos de texto (nombre, descripcion).
4. En `_dark.scss`, eliminar cualquier override de background para `.ability-card`.
5. Verificar visualmente que la informacion se lee bien sobre el fondo de la pagina en ambos modos (light y dark).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que las habilidades se muestran sin fondo de tarjeta en light mode |
| Manual | Verificar que las habilidades se muestran sin fondo de tarjeta en dark mode |
| Manual | Verificar que el texto sigue siendo legible |

## Criterios de aceptacion

- [x] Las habilidades no tienen fondo de tarjeta (sin background)
- [x] No hay borde izquierdo rojo
- [x] No hay border-radius ni box-shadow
- [x] El texto (nombre y descripcion) sigue siendo legible en light y dark mode
- [x] El espaciado entre habilidades se mantiene adecuado

## Notas

Cambio puramente estetico. Solo toca SCSS, no requiere cambios en la logica TypeScript ni en Go.
