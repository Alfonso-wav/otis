# Optimizar UI/UX del frontend para móvil

**ID**: 0061-android-mobile-ux-optimization
**Estado**: todo
**Fecha**: 2026-03-15

---

## Descripcion

Adaptar los estilos y la interacción del frontend para que funcione correctamente en pantallas móviles. Incluye ajustes de viewport, touch targets, navegación, scroll y rendimiento del WebView.

## Capas afectadas

- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Solo frontend (CSS/SCSS y algo de TypeScript).

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Asegurar meta viewport correcto |
| `frontend/src/styles/*.scss` | modificar | Media queries para mobile, touch targets mínimos 48px |
| `frontend/src/pages/pokedex.ts` | modificar | Virtualización de listas para rendimiento |
| `frontend/src/pages/builds.ts` | modificar | Layout adaptado para pantallas pequeñas |
| `frontend/src/router.ts` | modificar | Navegación adaptada (tabs bottom o hamburger menu) |

## Plan de implementacion

1. Auditar todos los componentes en viewport 360x800 (móvil estándar)
2. Ajustar media queries en SCSS para breakpoints mobile
3. Aumentar touch targets (botones, links) a mínimo 48x48px
4. Implementar navegación bottom tabs para mobile
5. Optimizar listas largas (pokédex) con virtualización o paginación
6. Probar en emulador Android y dispositivo real
7. Ajustar animaciones GSAP para rendimiento en WebView

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Test manual | Todas las páginas son usables en pantalla 360x800 |
| Test manual | No hay overflow horizontal |
| Test manual | Los botones son fácilmente pulsables con el dedo |

## Criterios de aceptacion

- [ ] Todas las páginas son usables en pantalla móvil (360px ancho mínimo)
- [ ] No hay scroll horizontal no deseado
- [ ] Touch targets >= 48px
- [ ] Navegación clara y accesible en mobile
- [ ] Rendimiento fluido (>30fps en scroll y animaciones)
- [ ] La versión desktop no se ve afectada (responsive, no destructivo)

## Notas

- Bootstrap 5 ya tiene breakpoints responsive; aprovecharlos al máximo.
- Considerar deshabilitar animaciones pesadas de GSAP en mobile si afectan rendimiento.
- Depende de: 0060 (proyecto Android funcional para testing).
