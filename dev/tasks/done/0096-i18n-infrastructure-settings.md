# 0096 — Infraestructura i18n y selector de idioma en Settings

## Descripción

La aplicación tiene strings hardcodeados mezclando español e inglés sin ningún sistema de internacionalización. Esta tarea monta la infraestructura base de i18n y añade el selector de idioma en la página de Settings.

Se necesita:
1. Instalar y configurar una librería i18n ligera (e.g. i18next) compatible con vanilla TypeScript.
2. Crear los archivos de traducción JSON para `es` y `en`.
3. Crear un módulo `i18n.ts` que exponga una función `t(key)` para obtener traducciones.
4. Añadir un selector de idioma (dropdown o toggle ES/EN) en la página de Settings, persistido en `localStorage`.
5. Migrar los strings de `index.html` y `settings.ts` al nuevo sistema como prueba de concepto.

## Capas afectadas

- **Core**: ningún cambio.
- **Shell**: ningún cambio.
- **APP (Frontend)**: nuevo módulo i18n, archivos de traducción, cambios en Settings y HTML.

## Cambios requeridos

### 1. Instalar dependencia i18n

Evaluar opciones ligeras para vanilla TS:
- **i18next** (popular, framework-agnostic) — `npm install i18next`
- **Alternativa casera**: un simple map de traducciones con función `t()` si se quiere evitar dependencias.

### 2. Crear archivos de traducción

```
frontend/src/locales/
  es.json    // traducciones en español
  en.json    // traducciones en inglés
```

Estructura de cada JSON:
```json
{
  "settings": {
    "title": "Configuración",
    "darkMode": "Modo Oscuro",
    "language": "Idioma"
  },
  "tabs": {
    "pokedex": "Pokédex",
    "builds": "Builds",
    "explore": "Explorar",
    "settings": "Configuración"
  },
  "common": {
    "loading": "Cargando...",
    "search": "Buscar",
    "back": "Volver",
    "next": "Siguiente",
    "previous": "Anterior"
  }
}
```

### 3. Crear módulo `frontend/src/i18n.ts`

- Exportar `t(key: string): string` para resolver traducciones.
- Exportar `setLocale(locale: string): void` para cambiar idioma.
- Exportar `getLocale(): string` para leer idioma actual.
- Persistir idioma en `localStorage` con key `"locale"`.
- Default: `"es"` (el idioma mayoritario actual).
- Emitir un evento custom `locale-changed` en `document` para que las páginas puedan re-renderizar.

### 4. Añadir selector de idioma en Settings

En `frontend/index.html`, dentro de `#tab-settings`, añadir una fila para idioma:

```html
<div class="settings-row">
  <span class="settings-label" data-i18n="settings.language">Idioma</span>
  <select id="language-select" class="form-select form-select-sm" style="width: auto;">
    <option value="es">Español</option>
    <option value="en">English</option>
  </select>
</div>
```

En `settings.ts`, conectar el `<select>` con `setLocale()`.

### 5. Migrar strings de index.html y settings

Reemplazar strings hardcodeados en:
- `frontend/index.html`: títulos de tabs, placeholders, labels estáticos.
- `frontend/src/settings.ts`: labels de la página de Settings.

Usar atributos `data-i18n` en HTML estático y llamadas a `t()` en código TS.

## Plan de implementación

1. Instalar i18next (o implementar solución casera).
2. Crear `frontend/src/i18n.ts` con `t()`, `setLocale()`, `getLocale()`.
3. Crear `frontend/src/locales/es.json` y `en.json` con las keys de Settings, tabs y common.
4. Modificar `frontend/index.html` — añadir selector de idioma y `data-i18n` attributes.
5. Modificar `frontend/src/settings.ts` — conectar selector y usar `t()`.
6. Modificar `frontend/src/main.ts` — inicializar i18n al arranque.
7. Verificar que cambiar idioma actualiza la UI de Settings y tabs.

## Tests

- Cambiar idioma en Settings y verificar que labels de Settings cambian.
- Refrescar la página y verificar que el idioma persiste.
- Verificar que el idioma por defecto es español.
- Verificar que no hay regresiones en dark mode toggle.

## Dependencias

- Ninguna tarea previa.

## Notas

- Esta tarea solo migra Settings y HTML base. Las demás páginas se migran en tareas 0097 y 0098.
- Al emitir `locale-changed`, las páginas que ya estén cargadas pueden escuchar y re-renderizar.
- Considerar si los nombres de Pokémon deben traducirse (probablemente no — son nombres propios).
