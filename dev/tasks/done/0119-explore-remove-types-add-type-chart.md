# Task 0119 — Explorar: eliminar pestaña Tipos y añadir tabla de fortalezas/debilidades

## Estado: done

## Goal

Dos cambios en la pestaña **Explorar**:

1. **Eliminar la subpestaña "Tipos"**: quitar completamente la subpestaña actual que muestra la lista de tipos con sus Pokémon.

2. **Añadir nueva subpestaña "Tabla de Tipos"**: una tabla interactiva 18×18 que muestre las relaciones de efectividad entre todos los tipos (atacante vs. defensor), con celdas coloreadas según el multiplicador (×2, ×1, ×0.5, ×0).

## Contexto técnico

### Subpestaña "Tipos" a eliminar

- `frontend/src/pages/explore.ts` gestiona las subpestañas de Explorar.
- Actualmente el tipo `ExploreTab = "types" | "regions" | "moves" | "abilities"`.
- `TAB_KEYS = ["types", "regions", "moves", "abilities"]` — eliminar `"types"`.
- El switch en `initPanel()` tiene `case "types": initTypes(panel)` — eliminar.
- El import `import { initTypes } from "./explore/types"` — eliminar.
- El icono `ICON_SHIELD` y la entrada `types: ICON_SHIELD` en `tabLabel` — eliminar.
- El fichero `frontend/src/pages/explore/types.ts` ya no se usará — se puede dejar en disco pero queda huérfano; se eliminará en una limpieza posterior (`/clean`).
- La pestaña activa por defecto era `"types"`. Con su eliminación, la primera pestaña será `"regions"`. Actualizar `activeTab` y `tabInited` acordemente.
- Claves i18n `explore.tabs.types` y toda la sección `"types": {...}` del fichero de locales se pueden dejar (no rompen nada), pero se limpiarán también con `/clean`.

### Nueva subpestaña "Tabla de Tipos"

**Nombre i18n**: `explore.tabs.typeChart` → "Tabla de Tipos" (ES) / "Type Chart" (EN).

**Datos**: La tabla de efectividad de tipos es **estática** (no ha cambiado desde la Gen VI con la adición de Hada). Se hardcodea en el frontend como un objeto TypeScript, sin llamadas al backend.

**Estructura de datos a hardcodear** en `type-chart.ts`:
```ts
// Orden canónico de los 18 tipos (mismo que PokéAPI + tipos de asset)
const TYPES = ["normal","fire","water","electric","grass","ice","fighting",
  "poison","ground","flying","psychic","bug","rock","ghost","dragon",
  "dark","steel","fairy"] as const;

// chart[attacker][defender] = multiplicador (0, 0.5, 1, 2)
const CHART: Record<string, Record<string, number>> = { ... };
```

**Renderizado**:
- Tabla HTML fija con 18 columnas + 1 columna de encabezado de fila.
- Fila de encabezado: iconos de tipo (SVG desde `/assets/types/{name}.svg`) + nombre abreviado.
- Celda vacía (×1) sin color. Celdas ×2 en verde, ×0.5 en rojo/naranja, ×0 en negro/gris.
- Mostrar el multiplicador como texto: "2×", "½", "0".
- La tabla debe tener scroll horizontal en mobile.
- Icono de la subpestaña: usar `ICON_CHART` (un SVG de tabla/cuadrícula de Heroicons).

**Fichero nuevo**: `frontend/src/pages/explore/type-chart.ts`
- Exporta `initTypeChart(panel: HTMLElement): void`
- Sin llamadas a API, sin estado asíncrono — renderizado síncrono inmediato.

**Integración en `explore.ts`**:
- Añadir `"typeChart"` a `ExploreTab`, `TAB_KEYS`, `tabInited`.
- Añadir `import { initTypeChart } from "./explore/type-chart"`.
- Añadir `case "typeChart": initTypeChart(panel)` en `initPanel()`.
- Añadir icono y label en `tabLabel()`.

### Localización

Añadir en `es.json` y `en.json`:
```json
// es.json — dentro de "explore.tabs"
"typeChart": "Tabla de Tipos"

// es.json — nueva sección raíz
"typeChart": {
  "title": "Efectividad de tipos",
  "attackingLabel": "Atacante →",
  "defendingLabel": "↓ Defensor"
}

// en.json análogo en inglés
```

## Acceptance criteria

- [ ] La subpestaña "Tipos" ya no aparece en la navegación de Explorar.
- [ ] La primera pestaña activa al abrir Explorar es "Regiones" (o la primera disponible).
- [ ] La subpestaña "Tabla de Tipos" aparece en la barra de navegación de Explorar con su icono.
- [ ] Al hacer click se muestra la tabla 18×18 de efectividades de tipo.
- [ ] Las celdas de ×2 aparecen en verde, ×0.5 en rojo/naranja, ×0 en gris oscuro.
- [ ] Las celdas de ×1 (neutro) no tienen color especial o son blancas/transparentes.
- [ ] Los encabezados de fila y columna muestran los iconos de tipo desde `/assets/types/{name}.svg`.
- [ ] La tabla tiene scroll horizontal en pantallas pequeñas.
- [ ] Las etiquetas de la subpestaña se actualizan al cambiar el idioma (locale-changed).
- [ ] No hay errores en consola al entrar en la pestaña.

## Archivos afectados

- `frontend/src/pages/explore.ts` — eliminar "types", añadir "typeChart", actualizar activeTab
- `frontend/src/pages/explore/type-chart.ts` — **crear nuevo** con datos hardcodeados y render
- `frontend/src/locales/es.json` — añadir `explore.tabs.typeChart` y sección `typeChart`
- `frontend/src/locales/en.json` — análogo en inglés
