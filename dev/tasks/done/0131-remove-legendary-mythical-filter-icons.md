# Task 0131 — Quitar iconos de filtros Legendario y Mítico en Pokédex

## Estado: done

## Goal

Eliminar los emojis ⭐ y ✨ de los botones de filtro "Legendario" y "Mítico" en la vista de la Pokédex, dejando solo el texto.

---

## Contexto técnico

### Archivos principales

- `frontend/index.html` — líneas 56-57, texto fallback de los botones `filter-legendary` y `filter-mythical`
- `frontend/src/locales/es.json` — líneas 29-30, claves `filters.legendary` y `filters.mythical` con emoji
- `frontend/src/locales/en.json` — líneas 29-30, claves `filters.legendary` y `filters.mythical` con emoji

### Estado actual

Los botones de filtro incluyen emojis Unicode como prefijo:
- `⭐ Legendario` / `⭐ Legendary`
- `✨ Mítico` / `✨ Mythical`

Las claves de la sección `pokedex` (líneas 64-65 en ambos JSON) ya están sin emoji y no requieren cambios.

---

## Cambios requeridos

### 1. HTML fallback — `frontend/index.html`

Línea 56: cambiar `⭐ Legendario` → `Legendario`
Línea 57: cambiar `✨ Mítico` → `Mítico`

### 2. Traducciones ES — `frontend/src/locales/es.json`

Línea 29: `"legendary": "⭐ Legendario"` → `"legendary": "Legendario"`
Línea 30: `"mythical": "✨ Mítico"` → `"mythical": "Mítico"`

### 3. Traducciones EN — `frontend/src/locales/en.json`

Línea 29: `"legendary": "⭐ Legendary"` → `"legendary": "Legendary"`
Línea 30: `"mythical": "✨ Mythical"` → `"mythical": "Mythical"`

---

## Archivos afectados

### Frontend
- `frontend/index.html` — quitar emojis del texto fallback de los botones
- `frontend/src/locales/es.json` — quitar emojis de `filters.legendary` y `filters.mythical`
- `frontend/src/locales/en.json` — quitar emojis de `filters.legendary` y `filters.mythical`

### Backend
Ninguno.

### Core
Ninguno.

### Shell
Ninguno.

---

## Acceptance criteria

- [x] El botón de filtro Legendario muestra solo "Legendario" / "Legendary" sin emoji.
- [x] El botón de filtro Mítico muestra solo "Mítico" / "Mythical" sin emoji.
- [x] El cambio aplica en ambos idiomas (ES y EN).
- [x] Los filtros siguen funcionando correctamente (toggle, estilo active, filtrado de lista).
- [x] Sin regresiones en otros filtros de la Pokédex.

---

## Dependencias

Ninguna.
