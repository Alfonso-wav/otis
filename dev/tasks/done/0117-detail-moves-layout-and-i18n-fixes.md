# Task 0117 — Vista detalle: layout tabla movimientos + i18n Historia y stats hexágono

## Estado: done

## Goal
Tres correcciones en la vista individual de Pokémon:
1. En la tabla de movimientos, reordenar columnas para que "Nivel" vaya antes de "Método", y reducir el espacio excesivo entre ambas columnas.
2. El contenido de la sección "Historia" (flavor text) debe mostrarse en español cuando el locale activo es ES.
3. Las etiquetas del hexágono de stats (hp, attack, defense…) deben mostrarse en español cuando el locale activo es ES.

## Contexto técnico

### 1. Tabla de movimientos — orden y spacing de columnas

- Función de render: `renderMoves()` en `frontend/src/pages/pokedex.ts` (~línea 993)
- Cabeceras actuales: "Nombre" → "Método" → "Nivel"
- Nuevo orden deseado: "Nombre" → "Nivel" → "Método"
- El ancho de las columnas se controla en `frontend/src/styles/_pokemon.scss` (~línea 726, bloque `.moves-table`)
- El espacio excesivo entre Método y Nivel probablemente viene de que Nivel usa `width: auto` o similar y Método ocupa demasiado espacio; ajustar con `width` fijos o `min-width` apropiados en los `<th>/<td>` de cada columna.

### 2. Historia — flavor text en español

- El flavor text se obtiene del backend vía `GetPokemonSpecies(name)` → campo `FlavorText`
- Backend: `shell/pokeapi.go` — función que construye `PokemonSpecies`, actualmente extrae `FlavorText` filtrando solo por versión del juego pero no por idioma → siempre devuelve texto en inglés
- PokéAPI devuelve `flavor_text_entries` con campo `language.name` (`"es"`, `"en"`, etc.)
- Solución: el backend debe filtrar `flavor_text_entries` por idioma. El idioma activo debe pasarse como parámetro desde el frontend al endpoint de la API, o bien el backend puede devolver entradas en ambos idiomas y el frontend elige.
  - Opción preferida (más simple): añadir query param `?lang=es|en` al endpoint de species; el backend filtra por ese idioma y devuelve el flavor text correspondiente.
  - Si no hay entrada en el idioma pedido, hacer fallback a inglés.
- Frontend: al llamar a `GetPokemonSpecies`, pasar el locale actual (`getLocale()`) como parámetro.

### 3. Stats hexágono — etiquetas en español

- Las etiquetas vienen de `pokemon.Stats[i].Name` (strings de PokéAPI: `"hp"`, `"attack"`, `"defense"`, `"special-attack"`, `"special-defense"`, `"speed"`)
- El chart las usa directamente en `frontend/src/charts/stats-chart.ts` (~línea 24): `series[0]?.stats.map((s) => s.Name)`
- Solución frontend: añadir una función de traducción de stat names en `frontend/src/i18n.ts` (similar a `typeName()`) y aplicarla en `stats-chart.ts` antes de pasar las etiquetas al chart.
- Añadir las claves de traducción en `frontend/src/locales/es.json` y `en.json` bajo una nueva sección `statNames`:
  ```json
  "statNames": {
    "hp": "PS",
    "attack": "Ataque",
    "defense": "Defensa",
    "special-attack": "At. Esp.",
    "special-defense": "Def. Esp.",
    "speed": "Velocidad"
  }
  ```
  (en inglés: HP, Attack, Defense, Sp. Atk, Sp. Def, Speed)

### i18n
- Locale activo: español por defecto
- Sistema i18n: `frontend/src/i18n.ts` — función `t(key)`, `typeName()` como referencia para `statName()`
- Locales: `frontend/src/locales/es.json` y `en.json`

## Acceptance criteria

### Tabla movimientos
- [ ] El orden de columnas es: Nombre → Nivel → Método
- [ ] No hay espacio excesivo entre columnas; la tabla se ve compacta y proporcional
- [ ] El orden nuevo se respeta tanto en cabeceras como en filas de datos

### Historia — flavor text
- [ ] En locale ES, el flavor text de "Historia" se muestra en español (cuando PokéAPI tiene entrada en `"es"`)
- [ ] En locale EN, el flavor text se muestra en inglés
- [ ] Si PokéAPI no tiene entrada en el idioma pedido, se muestra el inglés como fallback
- [ ] El cambio de locale recarga o actualiza el flavor text sin necesidad de navegar

### Stats hexágono
- [ ] En locale ES, las etiquetas del hexágono muestran: PS, Ataque, Defensa, At. Esp., Def. Esp., Velocidad
- [ ] En locale EN, las etiquetas muestran: HP, Attack, Defense, Sp. Atk, Sp. Def, Speed
- [ ] El cambio de locale actualiza las etiquetas del chart

## Implementación sugerida

### 1. Reordenar columnas y ajustar spacing

En `renderMoves()` de `pokedex.ts`:
- Cambiar el orden de `<th>` de cabecera: Nombre → `detail.moveLevelCol` → `detail.moveMethodCol`
- Cambiar el orden de `<td>` en cada fila: nombre → nivel → método
- En `_pokemon.scss`, ajustar `width` de las columnas (p.ej. nombre: `auto`, nivel: `3.5rem`, método: `7rem`) para eliminar el hueco.

### 2. Flavor text localizado

**Backend** (`shell/pokeapi.go`):
- En la función que construye `FlavorText`, leer el query param `lang` (o recibirlo como argumento).
- Filtrar `flavor_text_entries` por `language.name == lang`; si no hay resultado, fallback a `"en"`.
- Limpiar saltos de línea y caracteres especiales del texto resultante (ya se hace probablemente).

**API layer** (`app/`):
- Añadir el parámetro `lang` al handler del endpoint de species y pasarlo a la función shell.

**Frontend** (`pokedex.ts`):
- Al llamar `GetPokemonSpecies(name)`, añadir `?lang=${getLocale()}` a la URL.
- Suscribirse al evento `locale-changed` para volver a cargar la sección de Historia con el nuevo locale.

### 3. Traducción de stat names

**`frontend/src/i18n.ts`**:
```typescript
export function statName(name: string): string {
  return t(`statNames.${name}`) ?? name;
}
```

**`frontend/src/locales/es.json`** — añadir sección `statNames`.
**`frontend/src/locales/en.json`** — añadir sección `statNames`.

**`frontend/src/charts/stats-chart.ts`**:
- Importar `statName` de i18n.
- En el mapeo de etiquetas, usar `statName(s.Name)` en lugar de `s.Name`.
- El chart debe re-renderizarse al cambiar el locale (suscribirse a `locale-changed` o recibir etiquetas ya traducidas desde el caller).

## Archivos afectados

- `frontend/src/pages/pokedex.ts` — orden columnas en `renderMoves()`, llamada a species con lang, suscripción a `locale-changed`
- `frontend/src/styles/_pokemon.scss` — widths de columnas de la tabla de movimientos
- `frontend/src/charts/stats-chart.ts` — uso de `statName()` para etiquetas
- `frontend/src/i18n.ts` — nueva función `statName()`
- `frontend/src/locales/es.json` — sección `statNames`
- `frontend/src/locales/en.json` — sección `statNames`
- `shell/pokeapi.go` — filtrado de flavor text por idioma
- `app/` — paso del param `lang` al handler de species
