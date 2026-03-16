# Tags para columnas ocultas y tipo solo icono

**ID**: 0083-hidden-column-tags-type-icon-only
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Dos mejoras en las tablas del proyecto:

1. **Tags de columnas ocultas**: Cuando se oculta una columna con el icono del ojo, no hay forma de recuperarla. Se debe generar una zona de tags/chips encima de cada tabla que muestre un tag por cada columna oculta. Al pulsar el tag, la columna reaparece y el tag desaparece. Esto da visibilidad al usuario de qué columnas están ocultas y permite restaurarlas fácilmente.

2. **Columna Tipo solo icono**: En la columna "Tipo" de todas las tablas, mostrar solo el icono SVG del tipo (sin texto). Actualmente se muestra `<img> + texto`; se quiere solo el `<img>`.

## Capas afectadas

- **APP (frontend)**: Componente de toggle de columnas, renderizado de type badges en todas las tablas, estilos CSS.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/components/column-toggle.ts` | modificar | Añadir renderizado de tags de columnas ocultas encima de la tabla, con click para restaurar |
| `frontend/src/styles/_components.scss` | modificar | Estilos para los tags de columnas ocultas (`.col-hidden-tags`, `.col-hidden-tag`) |
| `frontend/src/styles/_dark.scss` | modificar | Estilos dark mode para los tags si es necesario |
| `frontend/src/pages/pokedex.ts` | modificar | Quitar texto del type badge, dejar solo icono en columna Tipo |
| `frontend/src/pages/explore/moves.ts` | modificar | Quitar texto del type badge, dejar solo icono |
| `frontend/src/pages/builds.ts` | modificar | Quitar texto del type badge en tablas, dejar solo icono |

## Plan de implementacion

### Parte 1: Tags de columnas ocultas

1. **Modificar `column-toggle.ts`**: En `initColumnToggle`, crear un contenedor `<div class="col-hidden-tags">` e insertarlo justo antes de la tabla (o del wrapper de la tabla). Este contenedor albergará los tags de columnas ocultas.

2. **Función `renderHiddenTags`**: Crear una función que, dado el set de columnas ocultas y el array de `ColumnConfig`, genere un tag `<button class="col-hidden-tag">` por cada columna oculta mostrando el label de la columna. Al hacer click en el tag, se restaura la columna (eliminar del set `hidden`, re-renderizar tags, aplicar visibilidad, actualizar iconos del ojo).

3. **Llamar a `renderHiddenTags`** cada vez que cambie el estado de visibilidad: en el handler del click del ojo y al inicializar (para columnas que ya estaban ocultas desde localStorage).

4. **Estilos CSS** para `.col-hidden-tags` y `.col-hidden-tag`:
   - Contenedor: `display: flex; flex-wrap: wrap; gap: 0.5rem; margin-bottom: 0.5rem;`
   - Tag: aspecto de chip/badge con fondo sutil, texto con el nombre de la columna, icono de "x" o del ojo para indicar que se puede restaurar. Cursor pointer. Hover con efecto visual.
   - Ocultar el contenedor si no hay columnas ocultas.

5. **Verificar** que al re-renderizar tablas (sorting, filtering), los tags se mantengan sincronizados. La función `reapplyColumnVisibility` también debe actualizar los tags.

### Parte 2: Columna Tipo solo icono

6. **En `pokedex.ts`**: Cambiar el renderizado de type badges en la tabla de:
   ```html
   <span class="type-badge type-${t.Name}"><img src="..." class="type-icon">${t.Name}</span>
   ```
   a:
   ```html
   <span class="type-badge type-${t.Name}" title="${t.Name}"><img src="..." alt="${t.Name}" class="type-icon"></span>
   ```
   Añadir `title` para que el nombre aparezca como tooltip al hacer hover.

7. **En `moves.ts`**: Mismo cambio — quitar texto, dejar solo icono con title/alt.

8. **En `builds.ts`**: Mismo cambio en las tablas que muestren tipo.

9. **Ajustar CSS** de `.type-badge` si es necesario: al quitar el texto, puede que el padding necesite ajuste para que el badge se vea bien solo con el icono (menos padding horizontal). Considerar hacer el icono ligeramente más grande (16-18px) para que sea más reconocible sin texto.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Al ocultar una columna, aparece un tag con su nombre encima de la tabla |
| (visual) | Al pulsar el tag, la columna reaparece y el tag desaparece |
| (visual) | Si hay varias columnas ocultas, se muestran varios tags |
| (visual) | Los tags se sincronizan con localStorage (recargar página mantiene tags) |
| (visual) | La zona de tags no aparece si no hay columnas ocultas |
| (visual) | Funciona en todas las tablas: Pokedex, Encounters, Moves, Abilities, Stats Config, Builds |
| (visual) | La columna Tipo muestra solo el icono SVG sin texto |
| (visual) | El tooltip del type badge muestra el nombre del tipo al hacer hover |
| (visual) | Los badges de tipo solo-icono se ven bien en dark mode |
| (visual) | Funciona correctamente en mobile (responsive) |

## Criterios de aceptacion

- [ ] Al ocultar una columna aparece un tag/chip encima de la tabla con el nombre de la columna
- [ ] Al pulsar el tag se restaura la columna y desaparece el tag
- [ ] Los tags se mantienen sincronizados con el estado de visibilidad (localStorage)
- [ ] Si no hay columnas ocultas, no se muestra la zona de tags
- [ ] Funciona en todas las tablas del proyecto que usan column-toggle
- [ ] La columna Tipo en todas las tablas muestra solo el icono SVG (sin texto)
- [ ] El nombre del tipo es accesible via tooltip (title) y alt text
- [ ] Funciona en dark mode
- [ ] Funciona en mobile (responsive)

## Notas

- El componente `column-toggle.ts` es reutilizable y se usa en 6+ tablas distintas. Los tags deben funcionar automáticamente en todas.
- La función `reapplyColumnVisibility` se llama tras re-renderizar tbody (sorting). Debe también actualizar los tags.
- Para los type badges solo-icono, considerar aumentar ligeramente el tamaño del icono y reducir el padding para que el badge sea compacto pero reconocible.
- Los SVGs ya están en `frontend/src/assets/types/` (tarea 0082).
