# Integrar Bootstrap + SCSS

**ID**: 0004-frontend-bootstrap-scss
**Estado**: todo
**Fecha**: 2026-03-13
**Depende de**: 0003-frontend-typescript-vite

---

## Descripcion

Integrar Bootstrap 5 como sistema de componentes base con personalización via SCSS. Migrar los estilos CSS actuales a SCSS, usando variables Bootstrap customizadas para mantener la identidad visual (colores Pokédex, type badges, etc.). No se elimina funcionalidad visual existente.

## Capas afectadas

- **Core**: sin cambios
- **Shell**: sin cambios
- **APP**: sin cambios

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/package.json` | modificar | Añadir `bootstrap`, `sass` como dependencias |
| `frontend/src/styles/main.scss` | crear | Entry point SCSS: variables custom + import Bootstrap + estilos propios |
| `frontend/src/styles/_variables.scss` | crear | Override de variables Bootstrap (colores, bordes, tipografía) |
| `frontend/src/styles/_pokemon.scss` | crear | Estilos específicos: type badges, stat bars, cards |
| `frontend/src/main.ts` | modificar | Importar `styles/main.scss` |
| `frontend/style.css` | eliminar | Reemplazado por SCSS |
| `frontend/index.html` | modificar | Eliminar link a style.css (Vite inyecta CSS) |

## Plan de implementacion

1. Instalar `bootstrap` y `sass` como dependencias
2. Crear `_variables.scss` con colores custom (rojo Pokédex, colores por tipo)
3. Crear `main.scss` que importa variables → Bootstrap → estilos custom
4. Crear `_pokemon.scss` migrando los estilos actuales de `style.css`
5. Actualizar `main.ts` para importar `main.scss`
6. Actualizar `index.html`: reemplazar grid y layout con clases Bootstrap donde aplique
7. Verificar que la app se ve igual o mejor que antes
8. Eliminar `style.css`

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | La app se ve correctamente, responsive funciona, type badges tienen colores correctos |
| `npm run build` | Build completa sin errores |

## Criterios de aceptacion

- [ ] Bootstrap cargado via SCSS (no CDN)
- [ ] Variables de color personalizadas sobrescriben los defaults de Bootstrap
- [ ] Grid responsive funciona con utilidades Bootstrap
- [ ] Type badges y stat bars mantienen su estilo visual
- [ ] No hay CSS inline ni `style.css` suelto

## Notas

- El documento dice "usa utilidades Bootstrap como base, nunca como producto final". Se personalizan los componentes via SCSS variables, no se usa Bootstrap "de fábrica".
- Mobile-first como indica el documento.
