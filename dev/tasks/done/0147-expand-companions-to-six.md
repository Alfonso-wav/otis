# Ampliar compañeros Pokemon de 1 a 6 en la barra superior

**ID**: 0147-expand-companions-to-six
**Estado**: todo
**Fecha**: 2026-04-12

---

## Descripcion

Actualmente la barra superior (header) muestra 1 Pokemon compañero seleccionado en Settings. Se quiere ampliar a **6 compañeros** (como un equipo completo). Formas de rellenar los slots:

1. **Manualmente**: desde Settings, seleccionar cada Pokemon individualmente para cada uno de los 6 slots.
2. **Desde equipos guardados**: en Simulaciones > Mis equipos, poder usar un equipo guardado para rellenar los 6 compañeros de la barra.

Actualmente:
- `settings.ts` gestiona 1 companion con key `companion-pokemon` en localStorage.
- `index.html` tiene un único `#header-companion`.
- `renderCompanion()` renderiza 1 sprite animado.
- Los equipos se gestionan en `builds.ts` con `ListTeams()`, `SaveToTeam()`, etc. Cada equipo tiene hasta 6 miembros.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna (usa API existente de equipos)
- **APP**: frontend — settings.ts, main.ts, index.html, estilos _pokemon.scss, _settings.scss, builds.ts

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Ampliar `#header-companion` para soportar hasta 6 sprites |
| `frontend/src/settings.ts` | modificar | Ampliar el setting de companion a 6 slots, añadir opción de cargar desde equipo guardado |
| `frontend/src/main.ts` | modificar | Inicializar los 6 companions al arrancar |
| `frontend/src/styles/_pokemon.scss` | modificar | Estilos para los 6 sprites en el header (fila de sprites más pequeños) |
| `frontend/src/styles/_settings.scss` | modificar | Layout del nuevo setting con 6 slots |
| `frontend/src/locales/en.json` | modificar | Etiquetas nuevas para los 6 slots y botón de cargar equipo |
| `frontend/src/locales/es.json` | modificar | Etiquetas nuevas para los 6 slots y botón de cargar equipo |

## Plan de implementacion

1. Cambiar el modelo de datos en localStorage: de `companion-pokemon` (string) a `companion-team` (JSON array de hasta 6 nombres). Mantener retrocompatibilidad leyendo el valor antiguo como primer slot.
2. Modificar `index.html`: el contenedor `#header-companion` albergará hasta 6 sprites en fila.
3. Modificar `renderCompanion()` en `settings.ts` para renderizar los 6 sprites (los slots vacíos no muestran nada o muestran un placeholder sutil).
4. En Settings, crear 6 inputs con autocomplete para cada slot. Añadir un botón "Cargar equipo" que abra un selector con los equipos guardados de `ListTeams()`.
5. Al seleccionar un equipo, rellenar los 6 slots con los nombres de los miembros del equipo.
6. Reducir el tamaño de cada sprite individual para que 6 quepan en el header (ej. 28px en desktop, 22px en mobile).
7. Actualizar traducciones EN/ES.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que aparecen hasta 6 sprites en el header |
| Manual | Verificar que se pueden seleccionar manualmente los 6 Pokemon |
| Manual | Verificar que se puede cargar un equipo guardado |
| Manual | Verificar retrocompatibilidad con el localStorage anterior (1 companion) |
| Manual | Verificar en desktop, tablet y mobile |
| Manual | Verificar en dark mode y light mode |

## Criterios de aceptacion

- [ ] El header muestra hasta 6 sprites de Pokemon compañeros
- [ ] En Settings hay 6 slots editables individualmente con autocomplete
- [ ] Hay un botón para cargar un equipo guardado de Simulaciones
- [ ] Si solo hay 1-5 compañeros, los slots vacíos se manejan limpiamente
- [ ] Retrocompatible: si hay un `companion-pokemon` antiguo, se migra al primer slot
- [ ] Los sprites son proporcionales al header (no desbordan)
- [ ] Funciona en desktop, tablet, mobile
- [ ] Funciona en dark mode y light mode
- [ ] Traducciones EN/ES completas
