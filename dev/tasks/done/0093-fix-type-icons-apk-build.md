# Corregir carga de iconos de tipos Pokemon en APK

**ID**: 0093-fix-type-icons-apk-build
**Estado**: done
**Fecha**: 2026-03-17

---

## Descripcion

Al compilar el APK, los iconos SVG de tipos de Pokemon no cargan en ninguna parte de la app. La causa es que los iconos están en `frontend/src/assets/types/` y se referencian con rutas hardcodeadas `/src/assets/types/{type}.svg`. Vite no incluye estos archivos en `dist/`, por lo que al hacer `cap sync` los SVGs no llegan al APK.

## Capas afectadas

- **APP (frontend)**: Mover assets y actualizar rutas en los archivos que renderizan iconos de tipo.

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/public/assets/types/*.svg` | crear (mover) | Mover los 18 SVGs desde `src/assets/types/` a `public/assets/types/` para que Vite los copie a `dist/` |
| `frontend/src/assets/types/` | eliminar | Eliminar directorio original tras mover los SVGs |
| `frontend/src/pages/pokedex.ts` | modificar | Cambiar ruta `/src/assets/types/` → `/assets/types/` (líneas ~404, ~583) |
| `frontend/src/pages/builds.ts` | modificar | Cambiar ruta `/src/assets/types/` → `/assets/types/` (líneas ~172, ~176) |
| `frontend/src/pages/explore/types.ts` | modificar | Cambiar ruta `/src/assets/types/` → `/assets/types/` (línea ~42) |
| `frontend/src/pages/explore/moves.ts` | modificar | Cambiar ruta `/src/assets/types/` → `/assets/types/` (línea ~110) |

## Plan de implementacion

1. Crear directorio `frontend/public/assets/types/`.
2. Mover los 18 archivos SVG de `frontend/src/assets/types/` a `frontend/public/assets/types/`.
3. Eliminar el directorio `frontend/src/assets/types/`.
4. En los 4 archivos TypeScript afectados, reemplazar todas las ocurrencias de `/src/assets/types/` por `/assets/types/`.
5. Verificar que los nombres de tipo usados en las rutas coincidan en case con los archivos SVG (los SVGs son lowercase; si `t.Name` es PascalCase, aplicar `.toLowerCase()`).
6. Ejecutar `npm run build` y verificar que los SVGs aparecen en `dist/assets/types/`.
7. Ejecutar `npx cap sync android` y verificar que los SVGs llegan a los assets del APK.

## Tests

| Tipo | Que se testea |
|------|---------------|
| Visual (dev) | Abrir cada página (Pokedex, Builds, Types, Moves) y verificar que los iconos de tipo se renderizan |
| Visual (APK) | Compilar APK, instalar y verificar iconos en todas las páginas |
| Build | Confirmar que `dist/assets/types/` contiene los 18 SVGs tras `npm run build` |

## Criterios de aceptacion

- [x] Los 18 SVGs están en `frontend/public/assets/types/`
- [x] No quedan referencias a `/src/assets/types/` en el código fuente
- [x] Los iconos de tipo se ven correctamente en modo desarrollo (`npm run dev`)
- [ ] Los iconos de tipo se ven correctamente en el APK compilado
- [x] No hay errores 404 en consola para los SVGs de tipos

## Notas

- Los SVGs son de la librería `duiker101/pokemon-type-svg-icons` (iconos blancos sobre fondo transparente).
- Vite copia el contenido de `public/` directamente a `dist/` sin procesamiento, lo cual es el comportamiento correcto para assets estáticos.
- Revisar si hay problemas de case sensitivity: los archivos son lowercase (`fire.svg`) pero los nombres de tipo podrían venir en PascalCase (`Fire`).
