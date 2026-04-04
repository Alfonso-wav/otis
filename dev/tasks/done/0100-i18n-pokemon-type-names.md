# 0100 — Traducir nombres de tipos de Pokémon al sistema i18n

## Descripción

Los nombres de los 18 tipos de Pokémon (fire, water, grass, etc.) se muestran siempre en inglés, sin pasar por el sistema de traducciones. Esto afecta principalmente a la pestaña "Tipos" en Explore, pero también a otros lugares donde se muestran badges de tipo (Pokédex, Builds, Movimientos).

Se necesita añadir las traducciones de los nombres de tipo a `es.json` y `en.json`, y modificar los componentes que muestran nombres de tipo para que usen `t()` en lugar del valor crudo de la API.

## Capas afectadas

- **Core**: ningún cambio.
- **Shell**: ningún cambio.
- **APP (Frontend)**: cambios en archivos de traducción y componentes que muestran tipos.

## Cambios requeridos

### 1. Ampliar archivos de traducción

Añadir key `typeNames` en `es.json` y `en.json` con los 18 tipos:

**es.json:**
```json
"typeNames": {
  "normal": "Normal",
  "fire": "Fuego",
  "water": "Agua",
  "grass": "Planta",
  "electric": "Eléctrico",
  "ice": "Hielo",
  "fighting": "Lucha",
  "poison": "Veneno",
  "ground": "Tierra",
  "flying": "Volador",
  "psychic": "Psíquico",
  "bug": "Bicho",
  "rock": "Roca",
  "ghost": "Fantasma",
  "dragon": "Dragón",
  "dark": "Siniestro",
  "steel": "Acero",
  "fairy": "Hada"
}
```

**en.json:**
```json
"typeNames": {
  "normal": "Normal",
  "fire": "Fire",
  "water": "Water",
  "grass": "Grass",
  "electric": "Electric",
  "ice": "Ice",
  "fighting": "Fighting",
  "poison": "Poison",
  "ground": "Ground",
  "flying": "Flying",
  "psychic": "Psychic",
  "bug": "Bug",
  "rock": "Rock",
  "ghost": "Ghost",
  "dragon": "Dragon",
  "dark": "Dark",
  "steel": "Steel",
  "fairy": "Fairy"
}
```

### 2. Crear helper de traducción de tipo

Crear una función helper (en `i18n.ts` o donde sea conveniente):

```typescript
function typeName(apiName: string): string {
  return t(`typeNames.${apiName.toLowerCase()}`);
}
```

### 3. Modificar componentes que muestran nombres de tipo

Archivos a modificar:

- **`frontend/src/pages/explore/types.ts`** (~línea 45): reemplazar `${t.Name}` por `${typeName(t.Name)}` en el span del nombre.
- **`frontend/src/pages/pokedex.ts`** (~líneas 409, 588): badges de tipo en vista detalle y tabla.
- **`frontend/src/pages/builds.ts`** (~líneas 181, 185): badges de tipo en atacante/defensor.
- **`frontend/src/pages/explore/moves.ts`** (~línea 113): columna de tipo en tabla de movimientos.
- **`frontend/src/components/pokemon-type-modal.ts`**: título del modal de tipo.

### 4. Verificar re-render en cambio de idioma

Los listeners de `locale-changed` ya existen en estos componentes (tarea 0099). Solo verificar que al cambiar idioma los nombres de tipo se actualizan correctamente.

## Plan de implementación

1. Añadir traducciones a `es.json` y `en.json`.
2. Crear función helper `typeName()`.
3. Modificar `types.ts` (caso principal).
4. Modificar el resto de componentes que muestran tipos.
5. Test visual de cambio de idioma en todas las vistas afectadas.

## Tests

- Verificar que la pestaña Tipos muestra nombres en español cuando el idioma es español.
- Verificar que la pestaña Tipos muestra nombres en inglés cuando el idioma es inglés.
- Verificar que los badges de tipo en Pokédex, Builds y Movimientos también están traducidos.
- Cambiar idioma en Settings y verificar que los nombres se actualizan sin recargar.
- Verificar que los SVG de iconos siguen cargando correctamente (usan el nombre API, no el traducido).

## Dependencias

- **0096**: infraestructura i18n.
- **0099**: listeners de locale-changed en componentes de Explore.

## Notas

- Los nombres de los SVG de tipo usan el nombre en inglés (`fire.svg`, `water.svg`). La traducción solo afecta al texto visible, NO a las rutas de assets.
- Los nombres de tipo de la API vienen en minúsculas (`fire`, `water`). La key de traducción debe coincidir con este formato.
- Los nombres oficiales en español de los tipos Pokémon son los usados en los juegos de la franquicia.
