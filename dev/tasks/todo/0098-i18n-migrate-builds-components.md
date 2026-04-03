# 0098 — Migrar strings de Builds y componentes compartidos al sistema i18n

## Descripción

Migrar todos los strings hardcodeados de la página Builds (simulador de batalla, team builder) y los componentes compartidos al sistema de traducciones i18n. Esta es la página más grande del proyecto (~3600 líneas) con la mayor mezcla de idiomas.

## Capas afectadas

- **Core**: ningún cambio.
- **Shell**: ningún cambio.
- **APP (Frontend)**: cambios en builds.ts y componentes compartidos.

## Cambios requeridos

### 1. Ampliar archivos de traducción

Añadir keys en `es.json` y `en.json` para:

**Builds / Team Builder:**
- Slots: "Slot 1", "Slot 2", etc.
- Labels: "Nivel", "Naturaleza", "Habilidad", "Objeto"
- Placeholder: "Nombre del Pokémon..."
- Stats: "HP", "Atk", "Def", "SpA", "SpD", "Vel" (mantener abreviaciones estándar por idioma)
- Labels de stats: "Stat", "Base", "IV", "EV", "Total"
- Acciones de equipo: "Guardar equipo", "Cargar equipo", "Nuevo equipo", "Eliminar"
- Mensajes: "Equipo guardado", "Selecciona un Pokémon"

**Battle Simulator:**
- Efectividad: "¡Super eficaz!", "Poco eficaz", "Sin efecto"
- Labels: "Atacante", "Defensor", "Resultado"
- Botones: "Simular", "Simular en lote"
- Log de batalla: mensajes de daño, turnos, etc.
- STAB, Priority y otros términos técnicos (considerar si traducir o mantener en inglés)

**Componentes compartidos:**
- `column-toggle.ts`: labels de toggle de columnas
- `sorting-overlay.ts`: labels de ordenamiento
- `ability-pokemon-modal.ts`: títulos y labels del modal
- `pokemon-type-modal.ts`: títulos y labels del modal
- `location-encounter-modal.ts`: títulos y labels del modal
- `autocomplete.ts`: placeholders y mensajes

### 2. Modificar `frontend/src/pages/builds.ts`

- Reemplazar todos los strings literales por `t()`.
- Template literals con HTML (innerHTML) necesitan usar `t()` dentro de los `${}`.
- Escuchar `locale-changed` para re-renderizar la UI.
- Tener cuidado con strings que se usan como identificadores vs strings de display.

### 3. Modificar componentes compartidos

- `frontend/src/components/column-toggle.ts`
- `frontend/src/components/sorting-overlay.ts`
- `frontend/src/components/ability-pokemon-modal.ts`
- `frontend/src/components/pokemon-type-modal.ts`
- `frontend/src/components/location-encounter-modal.ts`
- `frontend/src/autocomplete.ts`

### 4. Revisión final de consistencia

- Buscar strings hardcodeados restantes con grep/búsqueda manual.
- Verificar que no quedan mezclas de idiomas en ninguna página.
- Verificar que los archivos `es.json` y `en.json` tienen las mismas keys.

## Plan de implementación

1. Ampliar `es.json` y `en.json` con keys de Builds y componentes.
2. Migrar `builds.ts` — por secciones: team builder primero, luego battle simulator.
3. Migrar componentes compartidos.
4. Revisión exhaustiva de strings restantes con grep.
5. Test completo de cambio de idioma en toda la app.

## Tests

- Verificar que Builds muestra texto en el idioma seleccionado.
- Verificar que el team builder (slots, stats, guardado) está traducido.
- Verificar que el simulador de batalla (log, efectividad, botones) está traducido.
- Verificar que modales y componentes compartidos respetan el idioma.
- Cambiar idioma y verificar que toda la app se actualiza sin recargar.
- Verificar que no quedan strings hardcodeados en ningún rincón de la UI.
- Verificar que no hay regresiones en funcionalidad de batalla, equipos o modales.

## Dependencias

- **0096**: infraestructura i18n.
- **0097**: conviene que esté hecho primero para tener los patrones de migración establecidos, pero no es bloqueante.

## Notas

- `builds.ts` es el archivo más grande (~3600 líneas). La migración debe ser metódica.
- Algunos términos como STAB, IV, EV son universales en la comunidad Pokémon — considerar mantenerlos sin traducir.
- Los nombres de Pokémon, movimientos y habilidades vienen de la API y no se traducen.
- El log de batalla genera strings dinámicos — hay que usar interpolación con `t()`.
