# Fix: no se puede cargar el primer equipo guardado como compañero

**ID**: 0151-fix-companion-first-team-load
**Estado**: todo
**Fecha**: 2026-04-12

---

## Descripción

Al intentar cargar el primer equipo guardado como equipo compañero desde Ajustes > Compañeros, el equipo no se carga. Los demás equipos funcionan correctamente.

**Causa raíz**: En `frontend/src/settings.ts` (líneas 156-183), al hacer click en "Cargar desde equipo", se rellena un `<select>` con `innerHTML` y se escucha el evento `change`. El navegador selecciona automáticamente la primera `<option>` como valor por defecto, pero **no dispara el evento `change`**. Cuando el usuario intenta seleccionar el primer equipo (que ya está seleccionado por defecto), el navegador no emite `change` porque el valor no ha cambiado realmente.

**Flujo del bug**:
1. Usuario hace click en "Cargar desde equipo"
2. `ListTeams()` retorna los equipos, se crean las `<option>` (línea 164-166)
3. El `<select>` muestra el primer equipo seleccionado por defecto (comportamiento HTML estándar)
4. Se registra `addEventListener("change", ..., { once: true })` (línea 168)
5. El usuario selecciona el primer equipo → no hay cambio de valor → **no se dispara `change`** → el equipo no se carga

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — `frontend/src/settings.ts`

## Archivos a crear/modificar

| Archivo | Acción | Descripción |
|---------|--------|-------------|
| `frontend/src/settings.ts` | modificar | Corregir la lógica de selección del equipo compañero para que el primer equipo se pueda cargar |

## Plan de implementación

1. En `frontend/src/settings.ts`, después de rellenar el `<select>` con las opciones (línea 166), insertar una `<option>` placeholder deshabilitada y seleccionada al inicio:
   ```typescript
   teamSelect.innerHTML =
     `<option value="" disabled selected>${t("settings.selectTeamPlaceholder")}</option>` +
     teams.map((t) => `<option value="${t.name}">${t.name}</option>`).join("");
   ```
   Esto garantiza que cualquier equipo (incluido el primero) dispare `change` al seleccionarlo.
2. Verificar que el `change` handler (línea 169) descarte el valor vacío del placeholder.
3. Verificar i18n: el placeholder necesita clave en `es.json` y `en.json`.

## Tests

| Archivo | Qué se testea |
|---------|---------------|
| Manual | Cargar el primer equipo guardado como compañero — debe funcionar |
| Manual | Cargar otros equipos — sin regresión |
| Manual | Verificar con un solo equipo guardado |
| Manual | Verificar con múltiples equipos guardados |

## Criterios de aceptación

- [ ] El primer equipo de la lista se puede cargar como equipo compañero
- [ ] Todos los demás equipos siguen funcionando
- [ ] El placeholder es coherente con el idioma de la UI
- [ ] Funciona en desktop y mobile
