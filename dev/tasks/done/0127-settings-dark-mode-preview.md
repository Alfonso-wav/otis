# Task 0127 — Settings: previsualización de dark mode en tiempo real

## Estado: pending

## Goal

Actualmente el toggle de dark mode en Settings solo guarda el cambio como "pendiente" y no se aplica hasta pulsar el botón "Aplicar". El usuario debería poder ver el tema oscuro/claro en tiempo real al mover el toggle, sin necesidad de pulsar "Aplicar". El botón "Aplicar" sigue siendo necesario para persistir el cambio (guardarlo en localStorage).

---

## Contexto técnico

### Archivos principales

- `frontend/src/settings.ts` — lógica del toggle de dark mode y botón Aplicar

### Estado actual

En `settings.ts`, el listener del toggle (`dark-mode-toggle`) solo actualiza `pending.theme` y el estado del botón Aplicar. La función `applyTheme()` (que hace `document.documentElement.setAttribute("data-bs-theme", theme)`) solo se ejecuta dentro del handler del botón Aplicar.

### Comportamiento deseado

1. Al mover el toggle → aplicar el tema visualmente de inmediato (`applyTheme()`).
2. El botón "Aplicar" sigue habilitándose cuando hay cambios pendientes vs. lo guardado en localStorage.
3. Al pulsar "Aplicar" → persistir en localStorage.
4. Si el usuario cambia el toggle pero NO pulsa "Aplicar" y navega a otra pestaña → **revertir** al tema guardado en localStorage (el cambio visual era solo preview).
5. Si el usuario vuelve a Settings → el toggle refleja el estado guardado (no el preview).

### Flujo de revert

Escuchar el evento de cambio de pestaña (o el desmontaje de Settings) para llamar `applyTheme(currentTheme())`, que restaura el tema persistido.

---

## Cambios requeridos

### 1. Preview inmediato en el toggle — `settings.ts`

En el listener `change` del toggle, añadir la llamada a `applyTheme()`:

```typescript
toggle.addEventListener("change", () => {
  pending.theme = toggle.checked ? "dark" : "light";
  applyTheme(pending.theme);  // <-- AÑADIR: preview inmediato
  updateApplyButton();
});
```

### 2. Revertir al salir de Settings sin aplicar

Exportar una función `cleanupSettings()` que revierta el tema al guardado:

```typescript
export function cleanupSettings(): void {
  if (pending.theme !== null && pending.theme !== currentTheme()) {
    applyTheme(currentTheme());  // revertir al tema persistido
  }
  pending.theme = null;
  pending.locale = null;
}
```

### 3. Llamar a `cleanupSettings()` al cambiar de pestaña

En el router o en el código que gestiona las pestañas, cuando el usuario sale de la pestaña Settings, invocar `cleanupSettings()`.

### 4. Resetear toggle al entrar en Settings

En `initSettings()`, asegurar que el toggle refleja siempre el estado persistido (ya lo hace con `toggle.checked = saved === "dark"`), pero verificar que se ejecuta cada vez que se entra a la pestaña, no solo en la primera carga.

---

## Archivos afectados

### Frontend
- `frontend/src/settings.ts` — preview inmediato + export `cleanupSettings()`
- `frontend/src/router.ts` (o equivalente) — llamar `cleanupSettings()` al salir de Settings

### Backend
Ninguno.

---

## Acceptance criteria

- [ ] Al mover el toggle de dark mode, el tema cambia visualmente de inmediato (sin pulsar Aplicar).
- [ ] El botón "Aplicar" sigue habilitado/deshabilitado según haya cambios pendientes.
- [ ] Al pulsar "Aplicar", el tema se persiste en localStorage.
- [ ] Si el usuario cambia el toggle pero sale de Settings sin Aplicar, el tema vuelve al guardado.
- [ ] Al volver a Settings, el toggle refleja el estado persistido, no el preview anterior.
- [ ] Sin regresiones en el cambio de idioma.
