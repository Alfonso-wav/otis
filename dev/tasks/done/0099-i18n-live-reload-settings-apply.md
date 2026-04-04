# 0099 — Recarga en vivo de idioma y botón "Aplicar" en Settings

**Estado**: done

## Descripción

Cuando se cambia el idioma desde Settings, el cambio no se refleja en todas las pestañas hasta recargar la página. Esto se debe a que varias páginas y componentes (builds.ts, abilities.ts, moves.ts, regions.ts, types.ts) no escuchan el evento `locale-changed` y por tanto no re-renderizan su contenido traducido.

Además, se debe añadir un botón "Aplicar" en Settings para que los cambios (idioma, dark mode, etc.) no se apliquen hasta que el usuario pulse dicho botón, en lugar de aplicarse inmediatamente al interactuar con los controles.

## Capas afectadas

- **Core**: ningún cambio.
- **Shell**: ningún cambio.
- **APP (Frontend)**: cambios en settings, builds, explore sub-pages y componentes.

## Cambios realizados

### 1. Traducciones
- Añadida key `settings.apply` en `en.json` ("Apply") y `es.json` ("Aplicar").

### 2. Botón "Aplicar" en Settings
- Añadido botón `#settings-apply-btn` en `index.html`.
- Añadidos estilos en `_settings.scss`.

### 3. Refactorización de `settings.ts`
- Los cambios en idioma y dark mode se almacenan en estado pendiente.
- Solo al pulsar "Aplicar" se persisten y ejecutan.
- El botón se deshabilita cuando no hay cambios pendientes.

### 4. Listener `locale-changed` en `builds.ts`
- Se escucha el evento y se re-renderiza toda la UI de builds llamando a `buildLayout()` + `applyTranslations()`.

### 5. Listeners `locale-changed` en sub-páginas de Explore
- `abilities.ts`: re-invoca `initAbilities()` con datos cacheados.
- `moves.ts`: re-invoca `initMoves()` con datos cacheados.
- `regions.ts`: resetea estado de inicialización y re-invoca `initRegions()`.
- `types.ts`: re-renderiza `renderTypeCards()` con datos cacheados.

### 6. Componentes compartidos
- `column-toggle.ts`: no necesita cambios, se re-inicializa con cada re-render de las tablas.
- `sorting-overlay.ts`: no necesita cambios, se crea fresco cada vez con `t()` actual.
