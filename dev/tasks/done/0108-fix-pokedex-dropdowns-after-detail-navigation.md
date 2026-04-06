# Fix: filtros de generación y tipo dejan de funcionar tras navegar a un pokemon individual y volver

**ID**: 0108-fix-pokedex-dropdowns-after-detail-navigation
**Estado**: done
**Fecha**: 2026-04-06

---

## Descripcion

Después de hacer click en un pokemon de la pokédex (que abre la vista de detalle / tarjeta individual) y luego volver a la vista principal usando el botón de volver, los dropdowns de "Generación" y "Tipo" dejan de responder al click: al pulsar el trigger no se despliegan las opciones.

Este es el mismo síntoma que el bug 0107 (overlay bloqueando clics), pero con un disparador distinto: la navegación detail → list en lugar de los botones legendary/mythical.

## Causa raíz hipotética

La función `showView` en `frontend/src/animations/transitions.ts` utiliza GSAP `fromTo` para animar la entrada de `viewIn` (`listView`) con `{ opacity: 0, y: 20 }` → `{ opacity: 1, y: 0 }`. Al completarse, GSAP deja los estilos inline residuales en el elemento:
- `opacity: 1`
- `transform: translateY(0px)`

La propiedad `transform` en un elemento crea un nuevo **stacking context** y una nueva capa de compositing. En algunos navegadores (especialmente motores WebKit/Blink en entornos Capacitor/mobile), esto puede interferir con la capacidad de los elementos `position: absolute` hijos (como `.filter-dropdown__panel` con `z-index: 100`) de recibir eventos de click, o puede causar que el panel se renderice por debajo de otra capa.

**Posibilidad secundaria**: si en algún flujo existe un overlay (`showSortingOverlay`) que no se eliminó correctamente antes de navegar a la vista de detalle, ese overlay persistiría invisible durante la navegación y quedaría bloqueando la UI al volver.

## Capas afectadas

- **Frontend** únicamente: `frontend/src/animations/transitions.ts` y `frontend/src/pages/pokedex.ts`

## Archivos a modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/animations/transitions.ts` | modificar | Añadir `clearProps: "opacity,transform"` al tween de `viewIn` en `showView` para limpiar los estilos inline de GSAP al completar la animación |
| `frontend/src/pages/pokedex.ts` | modificar | Llamar `hideSortingOverlay()` de forma defensiva en el handler de `backBtn` antes de restaurar `listView` |

## Plan de implementacion

### 1. Limpiar estilos residuales de GSAP en `showView`

En `frontend/src/animations/transitions.ts`, en el `fromTo` de `viewIn`, añadir `clearProps: "opacity,transform"` en el objeto de variables "to":

```ts
tl.fromTo(
  viewIn,
  { opacity: 0, y: 20 },
  {
    opacity: 1,
    y: 0,
    duration: 0.3,
    ease: "power2.out",
    clearProps: "opacity,transform",  // ← AÑADIR
    onStart() {
      viewIn.style.display = "";
      viewIn.style.removeProperty("display");
      viewIn.classList.remove("hidden");
    },
  },
);
```

Esto hace que GSAP elimine sus estilos inline después de completar la animación, dejando el elemento sin `transform` residual y sin `opacity` inline. El elemento queda con sus valores CSS naturales.

Hacer lo mismo para el tween de `viewOut` para que tampoco quede con `opacity: 0` inline (ya tiene `display: none`, pero es más limpio):

```ts
tl.to(viewOut, {
  opacity: 0,
  duration: 0.2,
  ease: "power2.in",
  onComplete() {
    viewOut.style.display = "none";
    // GSAP's clearProps se aplica automáticamente al completar
  },
});
```

En realidad, `clearProps` en `viewOut` es opcional porque el elemento queda oculto con `display: none`. Lo importante es el `clearProps` en `viewIn`.

### 2. Llamar `hideSortingOverlay()` defensivamente en `backBtn`

En `frontend/src/pages/pokedex.ts`, el handler de `backBtn`:

```ts
backBtn.addEventListener("click", async () => {
  hideSortingOverlay();          // ← AÑADIR (defensivo: elimina overlay si por algún motivo quedó activo)
  const { disposeChart } = await import("../charts/stats-chart");
  disposeChart();
  await showView(listView, detailView);
});
```

## Tests

| Caso | Verificacion |
|------|--------------|
| Test manual | Click en pokemon → vista de detalle → botón volver → dropdown generación responde al click |
| Test manual | Click en pokemon → vista de detalle → botón volver → dropdown tipo responde al click |
| Test manual | Click en pokemon → vista de detalle → botón volver → combinación de filtros funciona normalmente |
| Test manual | Ciclo: legendary → mythical → click pokemon → volver → gen/tipo dropdowns abren correctamente |
| Test manual | Click rápido en pokemon antes de que carguen los datos → volver → dropdowns funcionan |

## Criterios de aceptacion

- [ ] Tras navegar a la tarjeta de un pokemon y volver, los dropdowns de generación y tipo abren correctamente
- [ ] El ciclo complete (varios pokémon clicados y vuelta atrás) no degrada los filtros
- [ ] La animación de entrada/salida de vistas sigue funcionando correctamente
- [ ] Los filtros legendary/mythical siguen funcionando (no regresión de 0107)

## Notas

- `clearProps: "opacity,transform"` en GSAP elimina esas propiedades de los estilos inline del elemento al completar la animación. El elemento vuelve a sus valores CSS definidos en las hojas de estilo (o los valores por defecto del navegador).
- En GSAP 3, `clearProps` acepta una string separada por comas con los nombres de propiedades CSS o los aliases de GSAP (`transform`, `opacity`, `x`, `y`, etc.).
- `clearProps: "all"` es más agresivo (limpia todo), pero puede causar parpadeos si el CSS natural del elemento difiere del estado final de la animación. Usar solo las propiedades específicas es más seguro.
- La llamada defensiva a `hideSortingOverlay()` en `backBtn` es de bajo riesgo: si no hay overlay activo, la función simplemente no hace nada.
