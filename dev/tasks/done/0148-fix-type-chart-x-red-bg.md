# Fix: fondo ovalado rojo detrás de las X en tabla de tipos

**ID**: 0148-fix-type-chart-x-red-bg
**Estado**: todo
**Fecha**: 2026-04-12

---

## Descripcion

En Explorer > Tabla de tipos, los botones "×" que aparecen junto a los iconos de tipo para filtrarlos tienen un fondo ovalado rojo (`background: #e53e3e; border-radius: 50%`). El diseño esperado es que solo se vea la **× en color rojo**, sin fondo circular/ovalado rojo detrás.

El estilo actual en `_explore.scss`:
```scss
.tc-remove-btn {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #e53e3e;
  color: white;
  ...
}
```

Debe cambiarse a fondo transparente con la × en rojo.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — estilos _explore.scss

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_explore.scss` | modificar | Cambiar `.tc-remove-btn` a fondo transparente, × en rojo. Ajustar hover. Actualizar también los breakpoints responsive. |

## Plan de implementacion

1. Cambiar `.tc-remove-btn`:
   - `background: transparent` (o `none`)
   - `color: #e53e3e` (rojo para la ×)
   - Eliminar `border-radius: 50%` (ya no hace falta sin fondo)
   - Ajustar hover: quizás `color: #c53030` y/o `transform: scale(1.15)`
2. Revisar los media queries que también definen `.tc-remove-btn` (991px, 767px, 479px) y asegurarse de que no re-añaden el fondo.
3. Verificar en dark mode que la × roja es visible.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que la × aparece en rojo sin fondo ovalado |
| Manual | Verificar hover: la × se destaca al pasar el cursor |
| Manual | Verificar que al pulsar la × se filtra el tipo correctamente |
| Manual | Verificar en todos los breakpoints (desktop, tablet, mobile) |
| Manual | Verificar en dark mode |

## Criterios de aceptacion

- [ ] La × junto a cada tipo se muestra en rojo sin fondo circular
- [ ] El hover de la × da feedback visual (cambio de color o escala)
- [ ] La funcionalidad de filtrar tipos sigue intacta
- [ ] Se ve correctamente en todos los breakpoints
- [ ] Se ve correctamente en dark mode y light mode
