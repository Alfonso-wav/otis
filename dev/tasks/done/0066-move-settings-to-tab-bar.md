# Mover botón Settings a la barra de pestañas

**ID**: 0066-move-settings-to-tab-bar
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

El botón de engranaje de Settings está actualmente en la esquina superior derecha del header (dentro de `.header-brand`), una zona inaccesible en móvil en orientación vertical. Mover el botón a la barra de pestañas (`#tab-nav`) junto a Pokédex, Por Tipo, Explorar y Builds, para que sea accesible en todas las resoluciones.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Solo frontend (HTML, TypeScript, SCSS).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Eliminar `#settings-btn` del `.header-brand` y añadir un nuevo botón/tab de Settings dentro de `#tab-nav` |
| `frontend/src/styles/_settings.scss` | modificar | Eliminar estilos del botón `.settings-btn` del header; adaptar el icono de engranaje para que encaje como un tab más en la barra |
| `frontend/src/styles/_tabs.scss` | modificar | Asegurar que el nuevo tab de Settings tenga el mismo estilo que los demás tabs (incluyendo el indicador activo) |
| `frontend/src/router.ts` | modificar | Actualizar la lógica de navegación: el tab Settings ahora es un tab más en `#tab-nav` en lugar de un botón independiente; al seleccionarlo ocultar el contenido de otros tabs y mostrar `#tab-settings`; al seleccionar otro tab, ocultar Settings |
| `frontend/src/settings.ts` | modificar | Ajustar la inicialización si referencia al antiguo `#settings-btn`; asegurar que el botón de volver dentro de la página Settings siga funcionando o eliminarlo si ya no es necesario (ahora se puede volver simplemente pulsando otro tab) |

## Plan de implementacion

1. **HTML — mover botón a tab-nav**:
   - Eliminar el `<button id="settings-btn">` del `.header-brand`.
   - Dentro de `<nav class="tab-nav" id="tab-nav">`, añadir un nuevo `<button class="tab-btn" data-tab="settings">` con el icono SVG de engranaje (reducido) o simplemente el texto "Settings" / icono, coherente con el estilo de los demás tabs.

2. **Router — integrar Settings como tab**:
   - En `initRouter()`, el tab `data-tab="settings"` debe funcionar igual que los demás: al hacer clic, `navigate("settings")` oculta los otros `tab-page` y muestra `#tab-settings`.
   - Eliminar la lógica especial de `navigateSettings()` que ocultaba `#tab-nav`. Ahora `#tab-nav` siempre es visible, incluso en la página Settings.
   - Evaluar si el botón de volver dentro de `#tab-settings` sigue siendo necesario. Si la navegación por tabs es suficiente, eliminarlo para simplificar.

3. **Estilos — coherencia visual**:
   - El nuevo tab de Settings debe verse igual que los demás (misma fuente, padding, indicador activo `::after`).
   - Si se usa un icono SVG en lugar de texto, asegurar que el tamaño sea coherente con el texto de los demás tabs.
   - Eliminar los estilos de `.settings-btn` que ya no aplican (el botón circular del header).

4. **Limpieza**:
   - Eliminar cualquier referencia al antiguo `#settings-btn` en el código.
   - Verificar que no queden estilos huérfanos en `_settings.scss`.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual | El tab Settings aparece en la barra de pestañas junto a los demás |
| Test manual | Al pulsar Settings, se muestra la página de Settings con el toggle de dark mode |
| Test manual | Al pulsar otro tab desde Settings, se vuelve correctamente a ese tab |
| Test manual | El tab-nav permanece visible cuando se está en Settings |
| Test manual | En móvil vertical, el tab Settings es accesible por scroll horizontal en la barra de tabs |
| Test manual | El dark mode toggle sigue funcionando correctamente |
| Test manual | No hay regresiones en la navegación de los demás tabs |

## Criterios de aceptacion

- [ ] El botón de engranaje ya no está en el header
- [ ] Existe un tab "Settings" (con icono o texto) en `#tab-nav`
- [ ] Al pulsar el tab Settings se muestra `#tab-settings` y se marca como activo
- [ ] La navegación entre tabs funciona normalmente incluyendo Settings
- [ ] El `#tab-nav` permanece visible en la página Settings
- [ ] Funciona correctamente en móvil vertical y horizontal
- [ ] El toggle de dark mode sigue funcionando
- [ ] No hay estilos ni código huérfano del antiguo botón del header

## Notas

- El tab-nav usa scroll horizontal en móvil, así que un 5to tab debería caber sin problemas.
- Considerar usar solo un icono de engranaje (sin texto) para el tab de Settings si se quiere ahorrar espacio, o texto corto como un simple icono SVG inline.
- El botón de volver dentro de `#tab-settings` probablemente ya no sea necesario, ya que el usuario puede navegar con las pestañas directamente.
