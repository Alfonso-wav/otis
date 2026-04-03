# 0097 — Migrar strings de Pokédex y Explore al sistema i18n

## Descripción

Con la infraestructura i18n montada en la tarea 0096, migrar todos los strings hardcodeados de las páginas Pokédex y Explore (types, regions, moves, abilities) al sistema de traducciones.

## Capas afectadas

- **Core**: ningún cambio.
- **Shell**: ningún cambio.
- **APP (Frontend)**: cambios en páginas Pokédex y Explore.

## Cambios requeridos

### 1. Ampliar archivos de traducción

Añadir keys en `es.json` y `en.json` para:

**Pokédex:**
- Labels de columnas: "Nombre", "Tipo", "Total", "HP", etc.
- Filtros: "Generación", "Tipo", "Legendario", "Mítico", "Todos"
- Placeholder de búsqueda: "Buscar por nombre o número..."
- Mensajes: "No se encontraron Pokémon", "Cargando...", "Aplicando filtros..."
- Paginación: "Anterior", "Siguiente"
- Vista: "Grid", "Tabla"

**Explore — Tipos:**
- Labels de pestañas: "Tipos", "Regiones", "Movimientos", "Habilidades"
- Efectividad: "¡Super eficaz!", "Poco eficaz", "Sin efecto", "Normal"
- Labels de tabla de tipos

**Explore — Regiones:**
- Labels de columnas y títulos de región
- Mensajes de carga

**Explore — Movimientos:**
- Labels de columnas: "Nombre", "Tipo", "Categoría", "Poder", "Precisión", "PP"
- Categorías: "Physical" → "Físico"/"Physical", "Special" → "Especial"/"Special", "Status" → "Estado"/"Status"

**Explore — Habilidades:**
- Labels de columnas: "Nombre", "Descripción"

### 2. Modificar `frontend/src/pages/pokedex.ts`

- Reemplazar todos los strings literales por llamadas a `t()`.
- Escuchar evento `locale-changed` para re-renderizar columnas, filtros y mensajes.
- Las configuraciones de columnas (`{ key, label }`) deben usar `t()` al renderizar, no al definir.

### 3. Modificar `frontend/src/pages/explore.ts`

- Reemplazar labels de tabs por `t()`.
- Escuchar `locale-changed`.

### 4. Modificar `frontend/src/pages/explore/types.ts`

- Reemplazar strings de efectividad y labels.

### 5. Modificar `frontend/src/pages/explore/regions.ts`

- Reemplazar labels de columnas y mensajes.

### 6. Modificar `frontend/src/pages/explore/moves.ts`

- Reemplazar labels de columnas y categorías.
- Las categorías "Physical"/"Special"/"Status" pasan a ser traducibles.

### 7. Modificar `frontend/src/pages/explore/abilities.ts`

- Reemplazar labels de columnas.

## Plan de implementación

1. Ampliar `es.json` y `en.json` con todas las keys necesarias.
2. Migrar `pokedex.ts` — es la página más usada, empezar aquí.
3. Migrar `explore.ts` (tabs).
4. Migrar `explore/types.ts`, `explore/moves.ts`, `explore/abilities.ts`, `explore/regions.ts`.
5. Probar cambio de idioma en cada página.

## Tests

- Verificar que todas las páginas muestran texto en el idioma seleccionado.
- Cambiar idioma en Settings y verificar que Pokédex y Explore se actualizan.
- Verificar que filtros, columnas, mensajes de carga y paginación están traducidos.
- Verificar que no hay strings hardcodeados visibles en la UI de estas páginas.
- Verificar que no hay regresiones en sorting, filtros o navegación.

## Dependencias

- **0096**: infraestructura i18n y módulo `t()`.

## Notas

- Los nombres de Pokémon, movimientos y habilidades vienen de la API y no se traducen aquí.
- Solo se traducen labels de UI, mensajes al usuario y placeholders.
- Los componentes compartidos (`column-toggle.ts`, `sorting-overlay.ts`) pueden necesitar ajustes menores si tienen strings propios.
