# Sistema de navegación por pestañas

**ID**: 0006-tab-navigation-system
**Estado**: done
**Fecha**: 2026-03-13
**Depende de**: 0005-frontend-echarts-gsap
**Agente**: .claude/agents/frontend

---

## Descripcion

Implementar un sistema de pestañas (tabs) en la barra superior de la aplicación que permita navegar entre diferentes vistas. Actualmente la app tiene una sola vista (lista + detalle de Pokémon). Se necesita una barra de navegación con botones/tabs que permitan cambiar entre las distintas secciones de la Pokédex. Este sistema será la base para las vistas que se crearán en las tareas 0007 y 0008.

## Capas afectadas

- **Core**: sin cambios
- **Shell**: sin cambios
- **APP**: sin cambios
- **Frontend**: nueva barra de tabs, refactor de vistas a sistema de páginas

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Añadir barra de tabs en el header, wrappear vistas existentes en contenedores de tab |
| `frontend/src/main.ts` | modificar | Extraer lógica de navegación entre tabs, sistema de routing por tab |
| `frontend/src/pages/pokedex.ts` | crear | Extraer la vista actual (list + detail) a su propio módulo |
| `frontend/src/pages/types.ts` | crear | Placeholder para la vista de tipos (se implementa en 0007) |
| `frontend/src/pages/explore.ts` | crear | Placeholder para la vista de exploración (se implementa en 0008) |
| `frontend/src/router.ts` | crear | Sistema simple de routing/tabs: registrar páginas, cambiar entre ellas |
| `frontend/src/styles/_tabs.scss` | crear | Estilos para la barra de tabs siguiendo la personalidad visual "Aura" |
| `frontend/src/styles/main.scss` | modificar | Importar nuevo partial `_tabs.scss` |

## Plan de implementacion

1. Diseñar la barra de tabs en el header: botones "Pokédex", "Por Tipo", "Explorar"
2. Crear `router.ts`: sistema simple que registra páginas y gestiona la navegación (show/hide + animaciones GSAP)
3. Extraer la vista actual (lista + detalle) a `pages/pokedex.ts` como módulo independiente
4. Crear placeholders para `pages/types.ts` y `pages/explore.ts` (contenido mínimo "Coming soon")
5. Modificar `index.html` para incluir los contenedores de cada tab
6. Integrar el router en `main.ts` como punto de entrada
7. Añadir transiciones GSAP al cambiar de tab (fade out → fade in)
8. Estilos de la barra de tabs: indicador activo, hover, responsive
9. Verificar que la vista actual (Pokédex) sigue funcionando igual tras el refactor

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Navegación entre tabs funciona, tab activo se resalta correctamente |
| Manual | La vista Pokédex (lista, búsqueda, detalle, paginación) sigue funcionando tras el refactor |
| Manual | Transiciones GSAP entre tabs son fluidas |
| `npm run build` | Build compila sin errores |

## Criterios de aceptacion

- [ ] Barra de tabs visible en el header con 3 botones: Pokédex, Por Tipo, Explorar
- [ ] Click en un tab cambia la vista activa con transición animada
- [ ] Tab activo tiene indicador visual claro (color, borde, etc.)
- [ ] La vista Pokédex existente funciona exactamente igual tras el refactor
- [ ] El sistema de routing es extensible para añadir más tabs fácilmente
- [ ] Diseño responsive: tabs se adaptan a pantallas pequeñas
- [ ] Estilos siguen la personalidad visual "Aura" (orgánico, cálido, no genérico)

## Notas

- El router es un sistema simple basado en show/hide de divs + GSAP, no un SPA router completo (no hay URL routing).
- Las páginas de "Por Tipo" y "Explorar" serán placeholders hasta que se implementen en 0007 y 0008.
- La barra de tabs debe integrarse con el header existente (fondo rojo Pokédex) sin romper el diseño actual.
