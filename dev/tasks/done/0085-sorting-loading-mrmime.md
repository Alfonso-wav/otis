# Loading overlay con Mr. Mime pensando durante el primer ordenamiento de tablas

**ID**: 0085-sorting-loading-mrmime
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Cuando el usuario hace clic en una columna para ordenar por primera vez, el SortCache necesita pre-computar todos los arrays ordenados, lo que toma un tiempo perceptible. Actualmente se muestra un texto de progreso inline, pero la experiencia es poco pulida.

Se quiere mostrar un **overlay centrado** con una imagen de **Mr. Mime "pensando"** (sprite o ilustracion) mientras se realiza esa primera ordenacion. El overlay debe:
- Aparecer centrado en la pantalla (o en el contenedor principal).
- Mostrar el sprite de Mr. Mime con una animacion sutil (por ejemplo, pulso o bounce con GSAP).
- Opcionalmente mostrar el porcentaje de progreso debajo de la imagen.
- Desaparecer cuando la ordenacion termine.

Esto aplica a todas las tablas que usan SortCache: Pokedex, Moves, Abilities.

## Capas afectadas

- **Core**: No afectada.
- **Shell**: No afectada.
- **APP**: Frontend — nuevo componente de loading overlay, integracion en las paginas con tablas.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/components/sorting-overlay.ts` | crear | Componente del overlay con Mr. Mime pensando |
| `frontend/src/styles/_components.scss` | modificar | Estilos del overlay (centrado, fondo semitransparente, animacion) |
| `frontend/src/pages/pokedex.ts` | modificar | Mostrar overlay durante primer sort en vez del texto de progreso |
| `frontend/src/pages/explore/moves.ts` | modificar | Mostrar overlay durante primer sort |
| `frontend/src/pages/explore/abilities.ts` | modificar | Mostrar overlay durante primer sort |
| `frontend/public/images/mrmime-thinking.png` | crear | Sprite/imagen de Mr. Mime pensando |

## Plan de implementacion

1. Conseguir/crear la imagen de Mr. Mime pensando y colocarla en `frontend/public/images/`.
2. Crear el componente `sorting-overlay.ts` que exponga funciones `show(container, progressText?)` y `hide()`.
3. Agregar estilos CSS para el overlay: fondo semitransparente, centrado absoluto, imagen con animacion (pulse/bounce via GSAP o CSS).
4. Integrar el overlay en `pokedex.ts`: al iniciar el primer sort, llamar `show()`; al terminar, llamar `hide()`. Reemplazar o complementar el texto de progreso actual.
5. Repetir integracion en `moves.ts` y `abilities.ts`.
6. Verificar que el overlay se muestra correctamente en modo claro y oscuro.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar visualmente que el overlay aparece centrado con Mr. Mime durante el primer sort y desaparece al terminar |

## Criterios de aceptacion

- [x] Al hacer el primer sort en cualquier tabla, aparece un overlay centrado con Mr. Mime pensando.
- [x] El overlay tiene una animacion sutil (pulse, bounce o similar).
- [x] El overlay desaparece cuando la ordenacion termina.
- [x] Funciona correctamente en las tres tablas: Pokedex, Moves, Abilities.
- [x] Se ve bien en modo claro y oscuro.
- [x] No aparece en sorts subsecuentes (ya cacheados, son instantaneos).

## Notas

- El proyecto ya usa GSAP para animaciones, se puede reutilizar para la animacion del sprite.
- El SortCache ya emite progreso via callback `onProgress` — se puede aprovechar para mostrar porcentaje debajo de Mr. Mime.
- Para la imagen se puede usar el sprite oficial de Mr. Mime de PokeAPI o un asset personalizado.
