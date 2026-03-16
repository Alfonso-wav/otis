# 0068 — Cambiar color rojo a azul marino en dark mode

## Descripción

En el modo oscuro (dark mode), reemplazar todos los acentos rojos (`$pokedex-red` y sus variantes) por una paleta de azul marino. El modo claro mantiene el rojo original sin cambios.

## Capas afectadas

- **APP (frontend/styles)**: solo cambios CSS/SCSS.

## Cambios requeridos

### 1. Definir variables de color azul marino para dark mode

En `frontend/src/styles/_variables.scss`, añadir las variantes azul marino:

```scss
$navy-blue: #1a365d;
$navy-blue-light: #2a4a7f;
$navy-blue-hover: #2c5282;
$navy-blue-text: #63b3ed;       // texto legible sobre fondos oscuros
$navy-blue-subtle: rgba(26, 54, 93, 0.15); // para sombras y fondos sutiles
```

### 2. Sobreescribir en `_dark.scss`

Dentro del selector `[data-bs-theme="dark"]`, sobreescribir cada aparición de rojo con su equivalente azul marino:

- **Section headers** (`_components.scss` → `.section-header`): `background: $navy-blue;`
- **Focus borders** (inputs, selects, filtros): `border-color: $navy-blue-hover;`
- **Focus shadows**: `box-shadow: 0 0 0 3px $navy-blue-subtle;`
- **Active filter chips/dropdowns**: colores azul marino en lugar de rojo.
- **Reset button hover**: usar `$navy-blue`, `$navy-blue-light`.
- **Pokedex header**: `background: $navy-blue;`
- **Search button**: `color: $navy-blue-text;` y hover `background: $navy-blue-light;`
- **Explore tabs activos**: `background: $navy-blue;`
- **Explore links** (ability count): `color: $navy-blue-text;`
- **Builds**: botones de búsqueda, focus de inputs, battle section border, botones de battle, win bar, labels.
- **Settings toggle**: `background: $navy-blue;` cuando está checked.
- **Back button** en detail view: `border-color` y `color` a azul marino.

### 3. Elementos que MANTIENEN rojo (semántica de peligro/error)

Los siguientes elementos conservan su color rojo en dark mode porque comunican peligro o error:

- **Error text** (`.error-text`, `.error`): mantener `#e53e3e` / `#fc8181`.
- **Delete buttons** (team delete, member delete, modal close hover): mantener `#e53e3e`.
- **Weak effectiveness rows** en `_dark.scss`: mantener `rgba(229, 62, 62, 0.12)` y `#fc8181` (indica debilidad/daño).

### 4. Hardcoded hex values de acento (NO error/peligro)

Buscar `#e53e3e` usado como acento de marca (no error) y sobreescribir en dark mode. Revisar contexto en:
- `_builds.scss` — botones, focus borders, labels de batalla, win bar.
- `_components.scss` — section headers, filtros activos.
- `_pokemon.scss` — pokedex header, search button.

## Notas

- **No tocar el modo claro**: todos los cambios van dentro del selector `[data-bs-theme="dark"]` en `_dark.scss`.
- Verificar contraste WCAG AA sobre fondos oscuros (`#1a202c`, `#2d3748`).
- **Rojo se mantiene para error/peligro/delete** — solo el acento de marca cambia a azul marino.

## Tests

- Alternar dark/light mode y verificar que:
  - Light mode sigue siendo rojo.
  - Dark mode muestra azul marino en todos los componentes listados.
  - No hay texto ilegible por bajo contraste.

## Dependencias

- Ninguna nueva.
