# Mover sección Tipos a subpestaña de Explorar

**ID**: 0079-move-types-to-explore-subtab
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Mover la sección de Tipos de Pokémon que actualmente vive en la pestaña independiente "Por tipo" (`types`) hacia la pestaña "Explorar" (`explore`) como una subpestaña más, al mismo nivel que Regiones, Movimientos y Habilidades. Después, eliminar la pestaña "Por tipo" del tab bar principal.

## Capas afectadas

- **APP (frontend)**: Reestructurar navegación, mover lógica de tipos a subtab de explore, actualizar HTML/SCSS/TS.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/explore/types.ts` | crear | Nuevo módulo subtab de tipos dentro de explore, reutilizando la lógica de `pages/types.ts` (grid de 18 tipos expandibles) |
| `frontend/src/pages/explore.ts` | modificar | Registrar la nueva subtab "Tipos" con botón `data-explore-tab="types"` y panel `data-panel="types"`, añadir import e init del nuevo módulo |
| `frontend/src/pages/types.ts` | eliminar | Ya no se necesita como página independiente |
| `frontend/src/main.ts` | modificar | Eliminar `registerPage` e `initTypes` de la pestaña "types" independiente |
| `frontend/src/router.ts` | modificar | Eliminar registro de la página "types" si está ahí |
| `frontend/src/index.html` | modificar | Eliminar el botón `data-tab="types"` del tab bar principal, eliminar el contenedor `#tab-types`, añadir botón subtab y panel en la sección explore |
| `frontend/src/styles/_explore.scss` | modificar | Añadir estilos para el panel de tipos dentro de explore (reutilizar estilos de `_types.scss`) |
| `frontend/src/styles/_types.scss` | eliminar o vaciar | Los estilos se integran en `_explore.scss` o se mantienen como parcial importado desde explore |

## Plan de implementacion

1. Crear `frontend/src/pages/explore/types.ts` extrayendo la lógica de renderizado de tipos desde `pages/types.ts`. Exportar funciones `initTypesPanel()` y `renderTypesPanel()` siguiendo el patrón de los otros subtabs (regions, moves, abilities).
2. En `explore.ts`, añadir el botón de subtab "Tipos" (con icono apropiado, e.g. shield/puzzle) en el HTML de la nav de explore, y registrar el panel `types` en la lógica de `switchTab()` e `initPanel()`.
3. En `index.html`, eliminar el botón `data-tab="types"` del tab bar principal y el contenedor `#tab-types`.
4. En `main.ts`, eliminar la llamada a `registerPage({ id: "types", ... })` y `initTypes()`.
5. Mover/integrar los estilos de `_types.scss` (cards de tipos, colores por tipo) dentro del contexto de explore o mantener el import pero scoped al panel de explore.
6. Eliminar `pages/types.ts` y limpiar imports huérfanos.
7. Verificar que la subtab "Tipos" aparece como primera o última opción en la nav de explore y que funciona el lazy loading al hacer click.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | La subtab "Tipos" aparece en Explorar junto a Regiones, Movimientos y Habilidades |
| (visual) | Al hacer click en "Tipos", se muestra el grid de 18 tipos expandibles |
| (visual) | Expandir/colapsar cards de tipo funciona correctamente |
| (visual) | Click en un Pokémon dentro de un tipo navega al Pokédex |
| (visual) | La pestaña "Por tipo" ya no aparece en el tab bar principal |
| (visual) | Dark mode se ve correctamente en la nueva subtab |

## Criterios de aceptacion

- [ ] La pestaña "Por tipo" desaparece del tab bar principal
- [ ] En "Explorar" aparece una nueva subtab "Tipos"
- [ ] La subtab muestra el grid de 18 tipos con cards expandibles (misma funcionalidad que antes)
- [ ] Lazy loading: los tipos solo se cargan al hacer click en la subtab
- [ ] Animaciones de expand/collapse funcionan igual
- [ ] Click en Pokémon navega al Pokédex
- [ ] Responsive: se ve bien en móvil y desktop
- [ ] Dark mode funciona correctamente
- [ ] No quedan referencias huérfanas a la antigua pestaña types

## Notas

- Seguir el mismo patrón que `regions.ts`, `moves.ts` y `abilities.ts` para consistencia.
- El orden sugerido de subtabs en explore: Tipos, Regiones, Movimientos, Habilidades (o al final, según preferencia del usuario).
- Los estilos de colores por tipo (`.type-header-${name}`) se reutilizan tal cual.
