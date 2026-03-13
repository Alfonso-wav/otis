# Vista de exploración: regiones, movimientos y datos detallados

**ID**: 0008-regions-moves-explore-view
**Estado**: todo
**Fecha**: 2026-03-13
**Depende de**: 0006-tab-navigation-system
**Agente**: .claude/agents/frontend

---

## Descripcion

Implementar la pestaña "Explorar" con datos enriquecidos de Pokémon: regiones donde aparecen, movimientos que pueden aprender, habilidades, cadenas evolutivas y más. Requiere múltiples endpoints nuevos de PokeAPI. La vista tendrá subsecciones navegables (regiones, movimientos, habilidades) con visualizaciones interactivas usando ECharts donde tenga sentido.

## Capas afectadas

- **Core**: nuevos tipos de dominio para regiones, movimientos, habilidades, cadenas evolutivas
- **Shell**: nuevos métodos en `PokeAPIClient` para múltiples endpoints de PokeAPI
- **APP**: nuevos bindings Wails para exponer toda la data al frontend
- **Frontend**: implementar vista completa en `pages/explore.ts` con subsecciones

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir structs: `Region`, `Location`, `Move`, `Ability`, `EvolutionChain`, etc. |
| `core/ports.go` | modificar | Añadir métodos: `FetchRegions()`, `FetchRegion(name)`, `FetchMove(name)`, `FetchAbility(name)`, `FetchEvolutionChain(id)` |
| `core/pokemon.go` | modificar | Funciones puras para procesar/filtrar movimientos, agrupar por generación, etc. |
| `shell/pokeapi.go` | modificar | Implementar todos los nuevos fetch contra PokeAPI (`/region`, `/move`, `/ability`, `/evolution-chain`) |
| `app/bindings.go` | modificar | Nuevos bindings Wails para cada endpoint |
| `frontend/src/pages/explore.ts` | modificar | Vista con subsecciones: Regiones, Movimientos, Habilidades |
| `frontend/src/pages/explore/regions.ts` | crear | Subsección de regiones: mapa/lista de regiones con sus locations y Pokémon |
| `frontend/src/pages/explore/moves.ts` | crear | Subsección de movimientos: lista filtrable, detalle de cada movimiento |
| `frontend/src/pages/explore/abilities.ts` | crear | Subsección de habilidades: lista y detalle |
| `frontend/src/charts/type-distribution.ts` | crear | Chart ECharts: distribución de tipos por región |
| `frontend/src/styles/_explore.scss` | crear | Estilos para toda la vista de exploración |
| `frontend/src/styles/main.scss` | modificar | Importar `_explore.scss` |

## Plan de implementacion

### Fase 1: Backend — nuevos tipos y endpoints

1. Definir tipos de dominio en `core/domain.go`:
   - `Region` (nombre, locations, generación principal)
   - `Location` (nombre, región, pokémon que aparecen)
   - `Move` (nombre, tipo, power, accuracy, PP, categoría, descripción)
   - `Ability` (nombre, descripción, pokémon que la tienen)
   - `EvolutionChain` (cadena con stages y condiciones)
2. Extender `PokemonFetcher` en `core/ports.go`
3. Implementar en `shell/pokeapi.go` los fetch a:
   - `GET /region` y `GET /region/{id}`
   - `GET /move/{id}`
   - `GET /ability/{id}`
   - `GET /evolution-chain/{id}`
4. Añadir bindings en `app/bindings.go`
5. Regenerar bindings Wails

### Fase 2: Frontend — vista de Regiones

6. Subsección Regiones: lista de regiones (Kanto, Johto, Hoenn, etc.)
7. Al seleccionar una región: mostrar locations y Pokémon exclusivos/característicos
8. Chart ECharts: distribución de tipos en la región seleccionada (pie chart o bar chart)

### Fase 3: Frontend — vista de Movimientos

9. Lista de movimientos con filtros: por tipo, por categoría (físico/especial/estado)
10. Búsqueda de movimientos por nombre
11. Detalle de movimiento: power, accuracy, PP, tipo, descripción, Pokémon que lo aprenden

### Fase 4: Frontend — vista de Habilidades

12. Lista de habilidades con búsqueda
13. Detalle de habilidad: descripción y Pokémon que la poseen

### Fase 5: Integración

14. Navegación interna entre subsecciones (mini-tabs o sidebar dentro de Explorar)
15. Links cruzados: desde un Pokémon en regiones poder ir a su detalle en Pokédex
16. Animaciones GSAP en todas las transiciones
17. Verificación completa de funcionalidad

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/pokemon_test.go` | Funciones puras de procesamiento de movimientos, habilidades, agrupaciones |
| Manual | Regiones cargan correctamente con sus locations |
| Manual | Movimientos se filtran y buscan correctamente |
| Manual | Habilidades muestran descripción y Pokémon asociados |
| Manual | Charts renderizan datos correctos |
| Manual | Navegación cruzada entre secciones funciona |
| `npm run build` | Build frontend compila |
| `go build ./...` | Build backend compila |

## Criterios de aceptacion

- [ ] La pestaña "Explorar" muestra 3 subsecciones: Regiones, Movimientos, Habilidades
- [ ] Regiones muestra las regiones del mundo Pokémon con sus Pokémon característicos
- [ ] Al menos 1 chart ECharts en la vista de regiones (distribución de tipos)
- [ ] Movimientos tiene búsqueda y filtros por tipo/categoría
- [ ] Detalle de movimiento muestra power, accuracy, PP, tipo y descripción
- [ ] Habilidades tiene lista con búsqueda y detalle
- [ ] Toda la data se carga lazy (no todo de golpe al abrir el tab)
- [ ] Animaciones fluidas en transiciones y aparición de elementos
- [ ] Diseño coherente con la estética "Aura" y consistente con las otras pestañas
- [ ] Responsive en pantallas pequeñas

## Notas

- PokeAPI tiene rate limiting — implementar carga progresiva y considerar cacheo en el futuro.
- Los endpoints de PokeAPI para regiones/movimientos devuelven mucha data — usar paginación y lazy loading agresivo.
- Las descripciones de movimientos/habilidades en PokeAPI vienen en múltiples idiomas — filtrar por `en` o `es` según preferencia.
- Las cadenas evolutivas se pueden visualizar como grafos o timelines — evaluar si ECharts o una visualización custom es más adecuada.
- Esta es la tarea más grande de las tres. Se puede implementar por fases, empezando por Regiones.
