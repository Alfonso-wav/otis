# Animacion Diglett de carga en tarjetas y cargas generales

**ID**: 0139-diglett-loading-cards
**Estado**: todo
**Fecha**: 2026-04-11

---

## Descripcion

Sustituir los estados de carga actuales (textos "Cargando...", spinners genericos) por la animacion de Diglett existente (`showDiglettOverlay`) en todas las cargas de Pokemon: tarjetas del Pokedex, carga de detalles individuales, y cualquier otro punto donde se carguen datos de Pokemon. Ya existe una implementacion de Diglett en `sorting-overlay.ts` que se usa en algunas partes; el objetivo es reutilizarla como loading universal.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — componentes de carga

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Usar Diglett overlay al cargar Pokemon en el grid y en la vista de detalle |
| `frontend/src/pages/explore/regions.ts` | modificar | Verificar que ya usa Diglett (lo usa) y que otros puntos de carga tambien |
| `frontend/src/components/location-encounter-modal.ts` | modificar | Reemplazar spinner CSS por Diglett overlay |
| `frontend/src/components/sorting-overlay.ts` | modificar | Si hace falta, crear variante mas pequena del Diglett para contextos inline |
| `frontend/src/styles/_components.scss` | modificar | Ajustar estilos si se necesita variante inline/compacta |

## Plan de implementacion

1. Identificar todos los puntos de la app donde se muestra un estado de carga (buscar "Cargando", "Loading", spinners CSS, overlays).
2. Para cada punto, evaluar si corresponde usar el Diglett overlay fullscreen (`showDiglettOverlay`) o una version mas compacta.
3. Reemplazar los spinners/textos de carga por la animacion de Diglett.
4. Asegurar que `hideSortingOverlay()` se llama correctamente al terminar cada carga.
5. Probar que no hay solapamiento de overlays si multiples cargas ocurren en paralelo.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que al cargar Pokemon en el Pokedex aparece Diglett |
| Manual | Verificar que al abrir un Pokemon individual aparece Diglett si hay latencia |
| Manual | Verificar que al abrir modales de encuentros aparece Diglett |
| Manual | Verificar que el Diglett desaparece correctamente al terminar la carga |

## Criterios de aceptacion

- [ ] Al cargar Pokemon en el grid del Pokedex, se muestra animacion Diglett
- [ ] Al cargar datos de un Pokemon individual, se muestra Diglett si hay espera
- [ ] Los spinners CSS genericos se reemplazan por Diglett donde aplique
- [ ] La animacion Diglett desaparece correctamente al terminar cada carga
- [ ] No hay conflictos si multiples cargas ocurren simultaneamente
- [ ] Funciona en desktop, web y APK

## Notas

El Diglett overlay actual es fullscreen (z-index 9999). Para cargas parciales (ej: solo el grid de tarjetas) podria hacer falta una variante contenida dentro de un elemento padre en vez de fullscreen. Evaluar caso por caso.
