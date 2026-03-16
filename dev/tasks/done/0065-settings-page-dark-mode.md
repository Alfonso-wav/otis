# Página Settings con Dark Mode

**ID**: 0065-settings-page-dark-mode
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Añadir un icono de engranaje (⚙) en la esquina superior derecha del header que navegue a una nueva vista "Settings". En esta vista, implementar un switch/toggle que active y desactive el modo oscuro (dark mode) de la aplicación. La preferencia debe persistir entre sesiones (localStorage).

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Solo frontend (HTML, TypeScript, SCSS).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Añadir botón de engranaje en `.header-brand` y el contenedor de la página Settings (`#tab-settings`) |
| `frontend/src/styles/_variables.scss` | modificar | Definir variables CSS custom properties para colores que cambien con dark mode (backgrounds, textos, bordes) |
| `frontend/src/styles/_settings.scss` | crear | Estilos de la página Settings (layout, toggle switch) |
| `frontend/src/styles/_dark.scss` | crear | Reglas de dark mode usando `[data-theme="dark"]` en `:root` que sobreescriban las custom properties |
| `frontend/src/styles/main.scss` | modificar | Importar `_settings.scss` y `_dark.scss` |
| `frontend/src/settings.ts` | crear | Módulo que gestiona la página Settings: inicializa el toggle, lee/escribe en localStorage, aplica `data-theme` al `<html>` |
| `frontend/src/main.ts` | modificar | Importar e inicializar el módulo settings; restaurar tema guardado al arrancar |
| `frontend/src/router.ts` | modificar | Registrar la navegación al tab Settings y el botón de vuelta |

## Plan de implementacion

1. **HTML — botón engranaje**:
   - Dentro de `.header-brand`, después de `#search-bar`, añadir un `<button id="settings-btn">` con un icono SVG de engranaje (Bootstrap Icons `bi-gear-fill` o inline SVG).
   - Estilizar para que quede alineado a la derecha con `ms-2`, color blanco, sin borde, cursor pointer.

2. **HTML — página Settings**:
   - Añadir `<div id="tab-settings" class="tab-page hidden">` después de los demás tabs en `<main>`.
   - Contenido: botón de volver (flecha ←), título "Settings", y una fila con label "Dark Mode" + toggle switch.

3. **Router — navegación**:
   - Al hacer clic en `#settings-btn`, ocultar todos los `.tab-page` y `.tab-nav`, mostrar `#tab-settings`.
   - Al hacer clic en el botón de volver, restaurar el tab activo anterior y mostrar `.tab-nav`.

4. **Settings module — dark mode toggle**:
   - Leer `localStorage.getItem('theme')` al inicializar.
   - Si es `"dark"`, poner `document.documentElement.dataset.theme = "dark"` y marcar el toggle como checked.
   - Al cambiar el toggle, alternar `data-theme` y guardar en localStorage.

5. **CSS custom properties**:
   - En `_variables.scss`, definir custom properties en `:root` para los colores principales (fondo, texto, bordes, cards, header).
   - En `_dark.scss`, sobreescribir esas properties bajo `[data-theme="dark"]`.
   - Migrar los estilos existentes más visibles (body background, card backgrounds, textos) para usar las custom properties.

6. **Estilos del toggle**:
   - Switch tipo iOS/Material con CSS puro (input checkbox + label con ::before).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual | El botón de engranaje es visible en header, a la derecha del buscador |
| Test manual | Al pulsar el engranaje se navega a la vista Settings ocultando tabs |
| Test manual | El botón de volver restaura la vista anterior |
| Test manual | El toggle de dark mode cambia los colores de la app inmediatamente |
| Test manual | Al recargar la página, la preferencia de dark mode se mantiene |
| Test manual | Funciona correctamente en móvil (touch target >= 44px) |

## Criterios de aceptacion

- [x] Icono de engranaje visible en el header, alineado a la derecha
- [x] Al pulsar el engranaje se muestra la página Settings y se oculta la navegación de tabs
- [x] Botón de volver funcional que restaura la vista anterior
- [x] Toggle de dark mode cambia el tema de la app en tiempo real
- [x] La preferencia se persiste en localStorage y se restaura al recargar
- [x] El dark mode aplica a: fondo de body, header, cards, tablas, modales, textos y bordes
- [x] No hay regresiones visuales en el modo claro (light mode)
- [x] Responsive: funciona en móvil y desktop

## Notas

- El proyecto usa vanilla TS sin framework, así que la navegación se gestiona mostrando/ocultando divs.
- Bootstrap 5.3 tiene soporte nativo para `data-bs-theme="dark"`, se puede evaluar usar eso en lugar de custom properties manuales si simplifica la implementación.
- El header actual usa `.header-brand` con flex y `#search-bar` con `ms-auto`. El botón de engranaje iría después del search bar.
