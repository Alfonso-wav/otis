# Fix: Diglett cortado por la mitad vertical en carga de vista individual

**ID**: 0145-fix-diglett-loading-clipped
**Estado**: todo
**Fecha**: 2026-04-12

---

## Descripcion

Cuando se carga la vista de detalle de un Pokemon individual (Pokedex > click en Pokemon), la animación de Diglett aparece cortada por la mitad vertical. El sprite de Diglett debería verse completo dentro de su contenedor `.diglett-clip`, pero se está recortando lateralmente.

Posibles causas:
- El contenedor `.diglett-clip` o `.diglett-clip--sm` tiene un width insuficiente respecto al sprite.
- El sprite (`.diglett-img` / `.diglett-img--sm`) no está centrado horizontalmente dentro del clip.
- El `position: absolute; bottom: 0` del sprite sin `left` explícito puede causar que el sprite se alinee al borde izquierdo y se corte por la derecha (o viceversa).

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — componente sorting-overlay, estilos _components.scss

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_components.scss` | modificar | Centrar horizontalmente el sprite de Diglett dentro de `.diglett-clip` y `.diglett-clip--sm` |
| `frontend/src/components/sorting-overlay.ts` | investigar | Verificar estructura HTML del inline Diglett |

## Plan de implementacion

1. Reproducir el bug: abrir la vista de detalle de un Pokemon y observar el Diglett de carga.
2. Inspeccionar `.diglett-clip` / `.diglett-clip--sm` y verificar si el sprite tiene centrado horizontal.
3. Añadir `left: 50%; transform: translateX(-50%)` al `.diglett-img` y `.diglett-img--sm` para centrar horizontalmente, o usar `left: 0; right: 0; margin: 0 auto`.
4. Verificar que la animación GSAP (translate Y) no se rompe con el nuevo posicionamiento.
5. Testear tanto la versión fullscreen (`.diglett-clip`) como la inline (`.diglett-clip--sm`).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Abrir detalle de un Pokemon y verificar que Diglett aparece completo |
| Manual | Verificar en regiones (usa inline Diglett) |
| Manual | Verificar animación de entrada/salida del Diglett |
| Manual | Verificar en mobile y desktop |

## Criterios de aceptacion

- [ ] El Diglett de carga aparece completo, sin recorte lateral
- [ ] La animación de pop-up/hide sigue funcionando correctamente
- [ ] Funciona tanto en overlay fullscreen como en inline
- [ ] Funciona en desktop y mobile
