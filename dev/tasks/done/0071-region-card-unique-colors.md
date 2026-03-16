# 0071 — Colores únicos para las tarjetas de región

## Descripción

Reemplazar el degradado compartido (`linear-gradient(135deg, #667eea 0%, #764ba2 100%)`) de los headers de las tarjetas de región en Explorar → Regiones por un color sólido diferente para cada región.

## Capas afectadas

- **APP** — frontend (`frontend/src/`)

## Archivos a modificar

1. **`frontend/src/styles/_explore.scss`** (línea 82) — Eliminar el `background: linear-gradient(...)` del selector `.region-card__header` y sustituirlo por un color base o dejarlo sin fondo (se aplicará por región).
2. **`frontend/src/pages/explore/regions.ts`** (línea 127-131) — Añadir un atributo `data-region` ya existente que se puede usar como gancho CSS, o bien asignar el color inline/como clase al header de cada tarjeta.

## Requisitos

1. Eliminar el degradado actual de `.region-card__header`.
2. Asignar un color sólido distinto a cada una de las 10 regiones:
   - `kanto` — Rojo (#E53E3E)
   - `johto` — Dorado/Ámbar (#D69E2E)
   - `hoenn` — Verde esmeralda (#38A169)
   - `sinnoh` — Azul acero (#3182CE)
   - `unova` — Gris oscuro (#4A5568)
   - `kalos` — Azul real (#4299E1)
   - `alola` — Naranja tropical (#ED8936)
   - `galar` — Púrpura (#805AD5)
   - `hisui` — Teal (#319795)
   - `paldea` — Escarlata (#C53030)
   (Los colores exactos son orientativos; el usuario puede ajustarlos.)
3. Mantener el texto blanco y la legibilidad del nombre y el chevron.
4. Respetar dark mode si hay estilos condicionales para `.dark`.

## Enfoque sugerido

- **Opción A (CSS puro):** Usar selectores de atributo `[data-region="kanto"] .region-card__header { background: #E53E3E; }` para cada región en `_explore.scss`.
- **Opción B (inline desde TS):** Crear un mapa `regionName → color` en `regions.ts` y asignar `style="background: ${color}"` al header en el template literal.

La opción A es preferible porque mantiene los estilos en SCSS.

## Tests

- Verificar visualmente que cada tarjeta de región tiene un color sólido diferente.
- Verificar que al expandir/colapsar el header mantiene su color.
- Verificar en dark mode que los colores se ven correctamente.
- Comprobar en móvil que no hay problemas de contraste o legibilidad.
