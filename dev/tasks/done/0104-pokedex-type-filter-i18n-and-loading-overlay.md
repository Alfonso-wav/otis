# Traducir tipos en filtro dropdown y mostrar Mr. Mime en carga de filtros

**ID**: 0104-pokedex-type-filter-i18n-and-loading-overlay
**Estado**: done
**Fecha**: 2026-04-04

---

## Descripcion

Dos problemas en la vista principal de la Pokedex:

1. **Filtro de tipos sin traducir**: Al desplegar el dropdown de tipos, los chips muestran el nombre en ingles (ej. "Fire", "Water") en lugar de usar la traduccion del idioma activo (ej. "Fuego", "Agua"). El resto de la app ya usa `typeName()` de `i18n.ts` para mostrar nombres traducidos, pero el dropdown de filtro usa `tp.Name` directamente.

2. **Mr. Mime no aparece al filtrar por legendarios/miticos**: Al activar el filtro de legendarios o miticos hay un tiempo de carga visible donde se muestra solo texto "Aplicando filtros...". El usuario quiere que el overlay de Mr. Mime (sorting-overlay) aparezca siempre que haya una operacion de carga, no solo al ordenar columnas.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend - pagina Pokedex y componente sorting-overlay

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Linea 971: cambiar `tp.Name.charAt(0).toUpperCase() + tp.Name.slice(1)` por `typeName(tp.Name)` para traducir los chips del filtro de tipos |
| `frontend/src/pages/pokedex.ts` | modificar | Lineas 1227-1243: anadir `showSortingOverlay()` y `hideSortingOverlay()` en los handlers de filtro legendario y mitico |
| `frontend/src/pages/pokedex.ts` | modificar | Funcion `loadFiltered()` (~linea 209) y `loadList()` (~linea 190): considerar mostrar Mr. Mime overlay en lugar del texto plano de carga cuando se invoca desde filtros |

## Plan de implementacion

### Parte 1: Traducir tipos en el dropdown

1. En `pokedex.ts` linea 971, reemplazar:
   ```typescript
   chip.textContent = tp.Name.charAt(0).toUpperCase() + tp.Name.slice(1);
   ```
   por:
   ```typescript
   chip.textContent = typeName(tp.Name);
   ```
   La funcion `typeName` ya esta importada en el archivo (se usa en otras partes).

### Parte 2: Mr. Mime en carga de filtros

1. En los handlers de filtro legendario (linea 1227) y mitico (linea 1236), envolver la llamada a `loadFiltered()`/`loadList()` con `showSortingOverlay()` antes y `hideSortingOverlay()` despues (en el await).
2. Hacer lo mismo en `applyFilters()` (linea 995) que es la funcion llamada al cambiar filtros de tipo y generacion.
3. Verificar que `showSortingOverlay` ya esta importado (linea 17 confirma que si).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Cambiar idioma a espanol, abrir dropdown de tipos: verificar que aparecen "Fuego", "Agua", "Planta", etc. |
| Manual | Cambiar idioma a ingles, abrir dropdown de tipos: verificar que aparecen "Fire", "Water", "Grass", etc. |
| Manual | Activar filtro legendario: verificar que aparece Mr. Mime flotante durante la carga |
| Manual | Activar filtro mitico: verificar que aparece Mr. Mime flotante durante la carga |
| Manual | Activar filtro de tipo/generacion: verificar que aparece Mr. Mime durante la carga |
| Manual | Verificar que Mr. Mime desaparece correctamente al terminar la carga |

## Criterios de aceptacion

- [x] Los chips del dropdown de tipos muestran nombres traducidos segun el idioma activo
- [x] Al filtrar por legendarios aparece el overlay de Mr. Mime durante la carga
- [x] Al filtrar por miticos aparece el overlay de Mr. Mime durante la carga
- [x] Al aplicar filtros de tipo/generacion aparece el overlay de Mr. Mime durante la carga
- [x] El overlay desaparece correctamente al completar la carga
- [x] No hay regresiones en el sorting de columnas (Mr. Mime sigue apareciendo ahi)

## Notas

- `typeName()` se importa desde `../i18n` y ya se usa en lineas 413, 614 y en otros componentes
- `showSortingOverlay(text?)` acepta un texto opcional; usar `t("pokedex.applyingFilters")` como texto
- `loadFiltered()` y `loadList()` son async, por lo que el patron es: show -> await load -> hide
- Los handlers de legendario/mitico (lineas 1227-1243) llaman a `loadFiltered()` o `loadList()` sin await actualmente; hay que anadir await
