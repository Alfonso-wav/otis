# Fix: filtros de generación y tipo dejan de funcionar tras clickar legendario/mítico

**ID**: 0107-fix-pokedex-filter-dropdowns-after-legendary
**Estado**: done
**Fecha**: 2026-04-06

---

## Descripcion

Después de toquetear los botones de "Legendario" y "Mítico" en la Pokédex (activarlos y desactivarlos varias veces), los dropdowns de "Generación" y "Tipo" dejan de responder: al hacer clic en el trigger no se despliegan las opciones.

## Causa raíz

Hay dos problemas relacionados en `frontend/src/pages/pokedex.ts`:

**1. Condición de carrera entre handlers async concurrentes**

Los handlers de `filterLegendaryBtn` y `filterMythicalBtn` son async y no tienen guardia de exclusión mutua. Si el usuario hace clic en ambos con suficiente rapidez (antes de que el overlay se pinte en pantalla), ambos handlers arrancan concurrentemente. En ese caso, la secuencia de `showSortingOverlay / hideSortingOverlay` queda desincronizada:

- Handler 1 crea el overlay → Handler 2 lo destruye (en su propio `showSortingOverlay`) y crea uno nuevo → Handler 1 termina y llama `hideSortingOverlay`, que destruye el overlay del Handler 2 → Handler 2 termina y llama `hideSortingOverlay` sobre `null`.
- Resultado: el overlay puede quedar visible (z-index 9999, `position: fixed`, cubre toda la pantalla) bloqueando todos los clics, incluidos los de los triggers de los dropdowns.

**2. `hideSortingOverlay` no está garantizado si hay un error inesperado**

Los handlers de legendary/mythical y `applyFilters` no tienen `try/finally`, por lo que si ocurre algún error no capturado entre `showSortingOverlay` y `hideSortingOverlay`, el overlay se queda atascado permanentemente.

## Capas afectadas

- **Frontend** únicamente: `frontend/src/pages/pokedex.ts`

## Archivos a modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Añadir guardia `isFiltering`, cambiar a `try/finally` en handlers y `applyFilters` |

## Plan de implementacion

1. **Añadir variable de guardia** — Declarar `let isFiltering = false` junto al resto de variables de módulo.

2. **Guardia en `filterLegendaryBtn` y `filterMythicalBtn`** — Al inicio del handler comprobar `if (isFiltering) return`. Usar `try/finally` para asegurar que `hideSortingOverlay()` e `isFiltering = false` se llaman siempre:
   ```ts
   filterLegendaryBtn.addEventListener("click", async () => {
     if (isFiltering) return;
     isFiltering = true;
     filter.legendary = !filter.legendary;
     filterLegendaryBtn.classList.toggle("active", filter.legendary);
     resetInfiniteScroll();
     filteredList = [];
     showSortingOverlay(t("pokedex.applyingFilters"));
     try {
       if (hasFilter()) await loadFiltered();
       else await loadList();
     } finally {
       hideSortingOverlay();
       isFiltering = false;
     }
   });
   ```
   Ídem para `filterMythicalBtn`.

3. **Guardia en `applyFilters`** — Mismo patrón: `if (isFiltering) return; isFiltering = true;` con `try/finally`.

4. **Guardia en el handler de `filterResetBtn`** — Añadir `if (isFiltering) return;` para evitar que el reset se ejecute mientras hay una operación en curso.

## Tests

| Caso | Verificacion |
|------|--------------|
| Test manual | Activar legendario → esperar carga → activar mítico → esperar carga → desactivar ambos → abrir dropdown generación: debe desplegarse |
| Test manual | Hacer clic rápido en legendario y mítico antes de que aparezca el overlay → solo una operación debe ejecutarse → overlay debe desaparecer al final |
| Test manual | Filtros de tipo y generación siguen funcionando con combinaciones de legendary/mythical |

## Criterios de aceptacion

- [ ] Tras activar/desactivar legendario y mítico varias veces, los dropdowns de generación y tipo responden correctamente al clic
- [ ] No es posible lanzar dos operaciones de filtrado concurrentes
- [ ] El overlay de carga siempre desaparece al terminar la operación, incluso si ocurre un error
- [ ] El resto de filtros (generación, tipo, reset) siguen funcionando con normalidad

## Notas

- La guardia `isFiltering` debe proteger también `applyFilters()` porque los chips de gen/tipo también muestran el overlay.
- No es necesario deshabilitar visualmente los botones legendary/mythical durante la carga (solo ignorar el clic), ya que el overlay ya comunica visualmente que hay una operación en curso.
- El overlay tiene `position: fixed; inset: 0; z-index: 9999` — si se queda atascado, bloquea absolutamente toda interacción con la UI.
