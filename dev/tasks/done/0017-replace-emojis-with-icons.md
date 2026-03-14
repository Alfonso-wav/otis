# 0017 — Reemplazar emojis con iconos reales en toda la app

## Estado

done

## Descripción

Sustituir todos los emojis usados como iconos en la interfaz por imágenes reales (SVG o PNG) obtenidas de las fuentes indicadas en `assets/urls.md`:
- **Sprites Pokémon**: usar imágenes de `pokemondb.net` (mejor calidad que las actuales de PokeAPI)
- **Iconos de categoría de movimiento y UI general**: usar iconos SVG de calidad (reemplazar ⚔️, ✨, 🛡️, 🏆, ⚖️, 🗺️, etc.)

## Contexto

Emojis actualmente en uso:

| Emoji | Archivo | Uso |
|-------|---------|-----|
| ⚔️ | `builds.ts:85`, `moves.ts:50`, `moves.ts:157`, `explore.ts:13` | Categoría física |
| ✨ | `builds.ts:86`, `moves.ts:51`, `moves.ts:158`, `explore.ts:14` | Categoría especial / Habilidades |
| 🛡️ | `builds.ts:87` | Categoría estado |
| 💤 | `moves.ts:52` | Categoría estado (filtro) |
| ❓ | `moves.ts:54` | Categoría desconocida |
| 🏆 | `compare.ts:84` | Ganador en comparador |
| ⚖️ | `compare.ts:84` | Empate en comparador |
| 🗺️ | `explore.ts:12` | Pestaña de regiones |

Fuentes disponibles en `assets/urls.md`:
- Sprites: `https://pokemondb.net/sprites/{name}` (galería por Pokémon)
- Iconos pokeball/Pokémon: `https://thenounproject.com/browse/icons/term/pokemon/`

## Capas involucradas

- **Frontend (APP)**: todos los archivos `.ts` con emojis, más estilos SCSS
- **Assets**: posible descarga/inclusión de SVGs locales en `assets/` o `frontend/src/assets/`
- **Core/Shell/Backend**: sin cambios

## Plan de implementación

### Paso 1 — Sprites Pokémon de mejor calidad

Cambiar la función `spriteURL()` en `builds.ts`, `compare.ts`, `types.ts` y `pokedex.ts` para usar sprites de Pokémon HOME de pokemondb.net:

```typescript
function spriteURL(nameOrId: string | number): string {
  // pokemondb usa nombres en minúsculas, sin guiones para formas especiales
  const name = typeof nameOrId === "number"
    ? nameOrId.toString()  // fallback con ID si no hay nombre
    : nameOrId.toLowerCase().replace(/[^a-z0-9-]/g, "");
  return `https://img.pokemondb.net/sprites/home/normal/${name}.png`;
}
```

Ajustar las llamadas existentes para pasar el nombre del Pokémon en lugar del ID donde sea posible (en builds.ts y compare.ts ya se tiene `pokemon.Name`).

**Fallback**: si la imagen de pokemondb falla (`onerror`), caer al sprite de PokeAPI actual.

```html
<img src="${spriteURL(pokemon.Name)}"
     onerror="this.src='https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${pokemon.ID}.png'"
     alt="${pokemon.Name}" />
```

### Paso 2 — Iconos de categoría de movimiento

Los iconos de categoría de movimiento en la Pokédex oficial usan imágenes específicas. Usar las imágenes de categoría de Bulbapedia/pokemondb:

```
Physical: https://img.pokemondb.net/images/icons/move-physical.png
Special:  https://img.pokemondb.net/images/icons/move-special.png
Status:   https://img.pokemondb.net/images/icons/move-status.png
```

Crear función `categoryIcon()` que devuelva un `<img>` tag en lugar de emoji:

```typescript
function categoryIcon(cat: string): string {
  const base = "https://img.pokemondb.net/images/icons";
  const map: Record<string, string> = {
    physical: `<img src="${base}/move-physical.png" class="move-cat-icon" alt="Physical" title="Physical">`,
    special:  `<img src="${base}/move-special.png"  class="move-cat-icon" alt="Special"  title="Special">`,
    status:   `<img src="${base}/move-status.png"   class="move-cat-icon" alt="Status"   title="Status">`,
  };
  return map[cat] ?? `<span class="move-cat-unknown">?</span>`;
}
```

Aplicar en: `builds.ts`, `moves.ts`, `explore.ts`.

### Paso 3 — Iconos de tipo Pokémon

Si en algún archivo se usan texto plano para los tipos (como "Fire", "Water"), explorar si pokemondb expone iconos de tipo:
```
https://img.pokemondb.net/images/icons/type/{type-lowercase}.png
```
Verificar disponibilidad real antes de usar. Si no están disponibles, mantener las badges de CSS actuales.

### Paso 4 — Iconos de UI general (trofeo, balanza, mapa)

Para `compare.ts` (🏆 ganador, ⚖️ empate) y `explore.ts` (🗺️ regiones), usar SVG inline o iconos de una CDN libre como [Heroicons](https://heroicons.com/) o [Phosphor Icons](https://phosphoricons.com/) que tienen licencia MIT:

- 🏆 ganador → icono `trophy` de Heroicons SVG inline
- ⚖️ empate → icono `scale` de Heroicons SVG inline
- 🗺️ regiones → icono `map` de Heroicons SVG inline
- ✨ habilidades → icono `sparkles` de Heroicons SVG inline

Descargar los SVGs necesarios a `frontend/src/assets/icons/` e importarlos como strings (Vite permite `?raw` imports de SVG).

```typescript
import trophyIcon from "../assets/icons/trophy.svg?raw";
import scaleIcon from "../assets/icons/scale.svg?raw";
```

### Paso 5 — Estilos SCSS

Añadir clases para los nuevos iconos:
```scss
.move-cat-icon {
  width: 20px;
  height: 14px;
  vertical-align: middle;
  image-rendering: pixelated; // para sprites pixel art
}

.ui-icon {
  width: 1em;
  height: 1em;
  vertical-align: middle;
  display: inline-block;

  svg {
    width: 100%;
    height: 100%;
    fill: currentColor;
  }
}
```

### Paso 6 — Verificación y limpieza

- Buscar cualquier emoji restante en todos los `.ts` con `grep -r "[^\x00-\x7F]" frontend/src/`
- Asegurarse de que los `alt` y `title` de todas las imágenes son descriptivos
- Verificar que los fallbacks funcionan en caso de error de carga de imagen

## Criterios de aceptación

- [ ] Los sprites de Pokémon cargan desde pokemondb.net (mayor resolución) con fallback a PokeAPI
- [ ] Los iconos de categoría de movimiento (físico/especial/estado) son imágenes reales, no emojis
- [ ] Los iconos de UI (trofeo, empate, mapa, habilidades) son SVGs, no emojis
- [ ] Todos los `<img>` tienen atributo `alt` descriptivo
- [ ] El fallback de sprites funciona correctamente si pokemondb no responde
- [ ] No quedan emojis usados como iconos en ningún archivo `.ts`
- [ ] Los estilos de los iconos son consistentes con el diseño existente

## Notas

- No tocar el backend ni core
- Las URLs de pokemondb.net son públicas; no requieren API key
- Si alguna URL de imagen no existe en pokemondb (formas alternativas, regionales), el `onerror` asegura el fallback
- Heroicons tiene licencia MIT: se pueden copiar los SVGs directamente al proyecto
- Los SVGs de heroicons están en: https://github.com/tailwindlabs/heroicons/tree/master/src/24/outline
