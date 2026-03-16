# Iconos originales de tipo Pokemon en badges de tablas

**ID**: 0082-pokemon-type-icons
**Estado**: done
**Fecha**: 2026-03-16

---

## Descripcion

Actualmente los badges de tipo Pokemon (acero, veneno, fuego, etc.) en todas las tablas del proyecto muestran solo texto con fondo de color. Se quiere añadir el icono/silueta original del tipo junto al nombre, manteniendo los colores actuales de fondo.

**Fuente de iconos**: Repositorio [duiker101/pokemon-type-svg-icons](https://github.com/duiker101/pokemon-type-svg-icons) — contiene 18 SVGs (uno por tipo) con silueta blanca sobre fondo transparente, ideales para superponer sobre los fondos de color existentes. Los 18 nombres de archivo coinciden exactamente con los nombres de tipo usados en el proyecto (`fire.svg`, `water.svg`, `steel.svg`, etc.).

**URL patrón CDN**: `https://raw.githubusercontent.com/duiker101/pokemon-type-svg-icons/master/icons/{type}.svg`

**Estrategia**: Descargar los 18 SVGs al directorio local de assets para servir offline (consistente con el patrón del proyecto para sprites), con fallback a la URL CDN. Insertar los iconos como `<img>` dentro de los `.type-badge` existentes, antes del texto del tipo.

## Capas afectadas

- **APP (frontend)**: Renderizado de badges de tipo en todas las tablas/vistas, descarga de assets SVG, estilos CSS.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/assets/types/` | crear | Directorio con los 18 SVGs de iconos de tipo descargados |
| `frontend/src/pages/pokedex.ts` | modificar | Actualizar renderizado de type badges en la tabla del Pokedex para incluir icono SVG antes del nombre |
| `frontend/src/pages/explore/moves.ts` | modificar | Actualizar renderizado de type badge en la tabla de movimientos para incluir icono SVG |
| `frontend/src/pages/explore/types.ts` | modificar | Actualizar cabeceras de tarjetas de tipo para incluir icono SVG si aplica |
| `frontend/src/pages/builds.ts` | modificar | Actualizar type badges en las tarjetas de equipo para incluir icono SVG |
| `frontend/src/components/pokemon-type-modal.ts` | modificar | Actualizar type badges en el modal de tipo si los muestra |
| `frontend/src/styles/_components.scss` | modificar | Ajustar estilos de `.type-badge` para acomodar icono + texto (flexbox, gap, tamaño del icono) |
| `frontend/src/styles/_dark.scss` | modificar | Verificar que los iconos blancos se vean bien en dark mode (los fondos de tipo ya son de color, debería funcionar) |

## Plan de implementacion

1. **Descargar los 18 SVGs** desde `https://raw.githubusercontent.com/duiker101/pokemon-type-svg-icons/master/icons/{type}.svg` y guardarlos en `frontend/src/assets/types/{type}.svg`. Los tipos son: `bug`, `dark`, `dragon`, `electric`, `fairy`, `fighting`, `fire`, `flying`, `ghost`, `grass`, `ground`, `ice`, `normal`, `poison`, `psychic`, `rock`, `steel`, `water`.

2. **Crear función utilitaria** `typeIconUrl(typeName: string): string` que retorne la ruta del asset local (importado por Vite) para el tipo dado. Centralizar en un módulo compartido o directamente en cada archivo si es más simple.

3. **Actualizar `.type-badge` en CSS** (`_components.scss`):
   - Hacer el badge `display: inline-flex; align-items: center; gap: 4px;`
   - Definir tamaño del icono: `width: 14px; height: 14px;` (ajustar según se vea)
   - El icono SVG es blanco — para tipos con texto oscuro (`$type-dark-text`), aplicar `filter: brightness(0) saturate(100%)` o similar para que el icono sea oscuro también.

4. **Actualizar renderizado en `pokedex.ts`**: Cambiar el template de type badges de:
   ```html
   <span class="type-badge type-${t.Name}">${t.Name}</span>
   ```
   a:
   ```html
   <span class="type-badge type-${t.Name}"><img src="${typeIconUrl(t.Name)}" alt="" class="type-icon">${t.Name}</span>
   ```

5. **Actualizar renderizado en `moves.ts`**: Mismo patrón — añadir `<img>` del icono dentro del badge de tipo.

6. **Actualizar renderizado en `types.ts`**: Añadir icono en la cabecera de cada tarjeta de tipo.

7. **Actualizar renderizado en `builds.ts`**: Añadir icono en los badges de tipo de las tarjetas de equipo.

8. **Actualizar `pokemon-type-modal.ts`**: Añadir icono si hay badges de tipo visibles.

9. **Verificar dark mode**: Los iconos son blancos sobre fondos de color — debería verse bien en ambos modos. Para los tipos con texto oscuro, asegurar que el icono también sea oscuro.

10. **Verificar responsive**: Comprobar que los badges con icono no rompan el layout en mobile.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| (visual) | Cada type badge en la tabla del Pokedex muestra icono + nombre |
| (visual) | El badge de tipo en la tabla de Moves muestra icono + nombre |
| (visual) | Las tarjetas de tipo en Explore muestran icono en la cabecera |
| (visual) | Los badges de tipo en Builds muestran icono + nombre |
| (visual) | Los colores de fondo se mantienen exactamente igual que antes |
| (visual) | Los iconos son visibles (blanco) sobre fondos oscuros (ghost, dark, fighting, etc.) |
| (visual) | Los iconos se oscurecen sobre fondos claros (grass, electric, ice, fairy, flying, bug) |
| (visual) | Funciona correctamente en dark mode |
| (visual) | Funciona correctamente en mobile (responsive, badges no overflow) |
| (visual) | Los 18 tipos tienen su icono correcto |

## Criterios de aceptacion

- [x] Los 18 SVGs de tipo están descargados en `frontend/src/assets/types/`
- [x] Todo badge `.type-badge` en el proyecto muestra icono SVG + nombre del tipo
- [x] Los colores de fondo actuales no cambian
- [x] Los iconos son visibles y contrastan correctamente sobre todos los colores de fondo
- [x] Los tipos con texto oscuro también tienen icono oscuro
- [x] Funciona en dark mode
- [x] Funciona en mobile (responsive)
- [x] No se rompe ningún layout existente (tablas, tarjetas, modales)

## Notas

- Los SVGs del repo `duiker101/pokemon-type-svg-icons` son siluetas blancas (`fill="white"`) sobre fondo transparente, lo cual es ideal para superponer sobre los fondos de color existentes.
- Para los tipos con texto oscuro (`grass`, `electric`, `ice`, `fairy`, `flying`, `bug`), se necesita un filtro CSS para oscurecer el icono y mantener el contraste.
- Seguir el patrón existente del proyecto de servir assets localmente (como los sprites de Pokemon) para funcionar offline en la app de escritorio y móvil.
- Los nombres de archivo de los SVGs (`fire.svg`, `water.svg`, etc.) coinciden exactamente con los valores de `PokemonType.Name` usados en el frontend.
