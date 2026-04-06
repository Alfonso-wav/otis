# Fix: imagen del pokemon anterior aparece brevemente al abrir un nuevo detalle

**ID**: 0109-fix-pokemon-detail-ghost-image
**Estado**: done
**Fecha**: 2026-04-06

---

## Descripcion

Al clickar en un pokemon desde la pokedex (tarjeta o fila de tabla), la primera vez carga bien. En los clicks posteriores, durante un instante breve se ve la imagen y el contenido del pokemon seleccionado anteriormente antes de que aparezca el spinner de carga del nuevo pokemon.

## Causa raíz

En `frontend/src/pages/pokedex.ts`, la función `showDetail` tiene el siguiente orden de operaciones:

```ts
async function showDetail(name: string): Promise<void> {
  await showView(detailView, listView);           // 1. Anima la vista (~500ms)
  detailEl.innerHTML = `<p class="loading">…</p>`;  // 2. Limpia el contenido
  const p = await GetPokemon(name);
  renderDetail(p);
}
```

`showView` anima la entrada del `detailView` durante ~500ms (200ms fade-out del list + 300ms fade-in del detail). Durante esa animación, `detailView` ya es visible pero `detailEl` todavía contiene el HTML del pokemon anterior (imagen incluida). El loading spinner solo se asigna **después** de que termina la animación.

**Fix**: mover la asignación del loading spinner a **antes** de `showView`, de modo que cuando el detail view se hace visible ya muestra el estado de carga y no el contenido anterior.

## Capas afectadas

- **Frontend** únicamente: `frontend/src/pages/pokedex.ts`

## Archivos a modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/pokedex.ts` | modificar | Mover `detailEl.innerHTML = loading` antes de `await showView(...)` en `showDetail` |

## Plan de implementacion

En `frontend/src/pages/pokedex.ts`, cambiar la función `showDetail` (líneas ~602-611):

**Antes:**
```ts
async function showDetail(name: string): Promise<void> {
  await showView(detailView, listView);
  detailEl.innerHTML = `<p class="loading">${t("common.loading")}</p>`;
  try {
    const p = await GetPokemon(name);
    renderDetail(p);
  } catch (err: unknown) {
    detailEl.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}
```

**Después:**
```ts
async function showDetail(name: string): Promise<void> {
  detailEl.innerHTML = `<p class="loading">${t("common.loading")}</p>`;
  await showView(detailView, listView);
  try {
    const p = await GetPokemon(name);
    renderDetail(p);
  } catch (err: unknown) {
    detailEl.innerHTML = `<p class="loading error-text">${String(err)}</p>`;
  }
}
```

Con este orden, cuando la animación de entrada del `detailView` comienza, el contenedor ya muestra el spinner de carga en lugar del contenido anterior.

## Tests

| Caso | Verificacion |
|------|--------------|
| Test manual | Click en pokemon A → detalle carga → volver → click en pokemon B → no se ve imagen de A en ningún momento |
| Test manual | Click rápido en varios pokemons consecutivos → nunca aparece imagen fantasma del anterior |
| Test manual | Primera vez que se abre el detalle → sigue funcionando correctamente |
| Test manual | La animación de entrada del detalle sigue siendo fluida |

## Criterios de aceptacion

- [ ] Al navegar al detalle de un segundo pokemon (o siguientes), no aparece la imagen del pokemon anterior en ningún frame visible
- [ ] El spinner de carga aparece correctamente desde el inicio de la transición
- [ ] La animación de entrada/salida de vistas sigue funcionando correctamente
- [ ] No hay regresión en los filtros ni otros comportamientos de la pokedex
