# 0020 — Ordenar columnas de la tabla Pokédex (ascendente / descendente)

## Descripcion

Permitir que el usuario ordene la tabla de la vista principal de la Pokedex haciendo clic en las cabeceras de las columnas. Cada clic alterna entre orden ascendente, descendente y sin orden (estado original). Se debe mostrar un indicador visual en la cabecera activa.

## Capas involucradas

- **Frontend (APP)**: `frontend/src/pages/pokedex.ts`, `frontend/src/styles/_pokemon.scss`
- **Core / Shell**: No se requieren cambios en backend. El ordenamiento se realiza en el frontend sobre los datos ya cacheados en `pokemonDataCache`.

## Contexto actual

- La tabla se renderiza en `renderTable()` dentro de `pokedex.ts` (lineas ~202-267).
- Los datos completos de cada Pokemon ya se cachean en `pokemonDataCache: Map<string, Pokemon>`.
- Las columnas actuales son: `#`, Sprite, Nombre, Tipo, HP, Atk, Def, SpA, SpD, Vel, Total.
- Las cabeceras (`<th>`) son estaticas, sin interaccion.
- Las stats se extraen del array `Pokemon.Stats[]` por indice (0=HP, 1=Atk, 2=Def, 3=SpA, 4=SpD, 5=Vel).

## Plan de implementacion

### Paso 1 — Estado de ordenamiento en `pokedex.ts`

Añadir variables de estado para el sorting:

```ts
type SortDirection = 'asc' | 'desc' | null;
type SortColumn = 'id' | 'name' | 'hp' | 'atk' | 'def' | 'spa' | 'spd' | 'vel' | 'total' | null;

let currentSortColumn: SortColumn = null;
let currentSortDirection: SortDirection = null;
```

- [x] Completado

### Paso 2 — Funcion `sortPokemonData()`

Crear una funcion pura que reciba el array de `Pokemon[]` y retorne una copia ordenada segun la columna y direccion activas:

```ts
function sortPokemonData(data: Pokemon[], column: SortColumn, direction: SortDirection): Pokemon[] {
  if (!column || !direction) return data;
  const sorted = [...data];
  const mult = direction === 'asc' ? 1 : -1;
  sorted.sort((a, b) => {
    switch (column) {
      case 'id': return mult * (a.ID - b.ID);
      case 'name': return mult * a.Name.localeCompare(b.Name);
      case 'hp': return mult * (a.Stats[0].BaseStat - b.Stats[0].BaseStat);
      case 'atk': return mult * (a.Stats[1].BaseStat - b.Stats[1].BaseStat);
      case 'def': return mult * (a.Stats[2].BaseStat - b.Stats[2].BaseStat);
      case 'spa': return mult * (a.Stats[3].BaseStat - b.Stats[3].BaseStat);
      case 'spd': return mult * (a.Stats[4].BaseStat - b.Stats[4].BaseStat);
      case 'vel': return mult * (a.Stats[5].BaseStat - b.Stats[5].BaseStat);
      case 'total': {
        const totalA = a.Stats.reduce((s, st) => s + st.BaseStat, 0);
        const totalB = b.Stats.reduce((s, st) => s + st.BaseStat, 0);
        return mult * (totalA - totalB);
      }
      default: return 0;
    }
  });
  return sorted;
}
```

- [x] Completado

### Paso 3 — Cabeceras interactivas en `renderTable()`

Modificar la generacion de `<th>` para que las columnas ordenables tengan:
- Atributo `data-sort="id|name|hp|atk|def|spa|spd|vel|total"`.
- Clase CSS `sortable` para cursor pointer.
- Un `<span>` para el indicador de direccion (flecha arriba/abajo o ninguna).
- Excluir Sprite y Tipo del sorting (no tiene sentido ordenar por imagen o tipo multi-valor).

```html
<th class="sortable" data-sort="id"># <span class="sort-indicator"></span></th>
<th>Sprite</th>
<th class="sortable" data-sort="name">Nombre <span class="sort-indicator"></span></th>
<th>Tipo</th>
<th class="sortable" data-sort="hp">HP <span class="sort-indicator"></span></th>
...
```

- [x] Completado

### Paso 4 — Event listener en cabeceras

Tras renderizar la tabla, añadir event listeners a los `<th>` con clase `sortable`:

```ts
container.querySelectorAll('th.sortable').forEach(th => {
  th.addEventListener('click', () => {
    const col = th.dataset.sort as SortColumn;
    if (currentSortColumn === col) {
      // Ciclo: asc -> desc -> null
      if (currentSortDirection === 'asc') currentSortDirection = 'desc';
      else if (currentSortDirection === 'desc') { currentSortDirection = null; currentSortColumn = null; }
    } else {
      currentSortColumn = col;
      currentSortDirection = 'asc';
    }
    renderTable(getCurrentPageItems()); // re-renderizar con el orden activo
  });
});
```

- [x] Completado

### Paso 5 — Aplicar sorting en `renderTable()`

Dentro de `renderTable()`, despues de cargar/cachear los datos y antes de generar el HTML de las filas, aplicar el sorting:

```ts
const sortedData = sortPokemonData(pokemonArray, currentSortColumn, currentSortDirection);
// usar sortedData para generar las filas en lugar de pokemonArray
```

- [x] Completado

### Paso 6 — Indicadores visuales en `_pokemon.scss`

Añadir estilos para las cabeceras ordenables y los indicadores:

```scss
.poke-table th.sortable {
  cursor: pointer;
  user-select: none;
  position: relative;

  &:hover {
    background-color: rgba(255, 255, 255, 0.08);
  }
}

.sort-indicator {
  margin-left: 0.25rem;
  font-size: 0.7rem;
  opacity: 0.4;

  &.asc::after { content: '\25B2'; opacity: 1; }   /* triangulo arriba */
  &.desc::after { content: '\25BC'; opacity: 1; }  /* triangulo abajo */
}

th.sortable.active .sort-indicator {
  opacity: 1;
}
```

- [x] Completado

### Paso 7 — Reset de sorting al cambiar pagina o filtros

Cuando el usuario cambia de pagina, aplica filtros, o cambia de vista, resetear el estado de sorting para evitar inconsistencias:

```ts
currentSortColumn = null;
currentSortDirection = null;
```

Esto se aplica en: `prevPage`, `nextPage`, `loadFiltered`, y el toggle de vista.

- [x] Completado

## Criterios de aceptacion

- [x] Clicar en una cabecera ordena la columna en ascendente.
- [x] Clicar de nuevo en la misma cabecera cambia a descendente.
- [x] Un tercer clic vuelve al orden original (sin sorting).
- [x] Se muestra un indicador visual (triangulo) en la cabecera activa con la direccion.
- [x] Las columnas ordenables son: #, Nombre, HP, Atk, Def, SpA, SpD, Vel, Total.
- [x] Sprite y Tipo NO son ordenables.
- [x] El sorting se resetea al cambiar de pagina o aplicar filtros.
- [x] No se introducen cambios en Core ni Shell.
- [x] La animacion GSAP de entrada de filas se ejecuta tras re-ordenar.

## Dependencias externas nuevas

Ninguna.

## Notas

- El sorting opera sobre los datos ya cargados en cache; no requiere llamadas adicionales al backend.
- Al re-renderizar tras un sort, se puede reutilizar la animacion de entrada de filas existente con GSAP stagger.
