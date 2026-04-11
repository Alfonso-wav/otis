# Eliminar pantalla negra antes del splash Snorlax

**ID**: 0137-fix-splash-black-screen
**Estado**: todo
**Fecha**: 2026-04-11

---

## Descripcion

Al arrancar la app hay un periodo largo de pantalla completamente negra antes de que aparezca el Snorlax del splash screen. Esto ocurre porque el splash screen (fondo #1a202c + Snorlax SVG) no se muestra hasta que el DOM carga y los assets se resuelven. El objetivo es que el Snorlax aparezca inmediatamente (tiempo 0) al abrir la app. Si la solucion requiere una imagen placeholder previa (por ejemplo un PNG inline o base64), se solicitara al usuario.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — splash screen HTML/CSS/assets

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Optimizar carga del splash para que sea instantaneo (inline SVG o base64, critical CSS) |
| `frontend/src/styles/_splash.scss` | modificar | Asegurar que el splash sea visible sin esperar a que cargue el bundle JS |
| `frontend/src/main.ts` | modificar | Revisar si el flujo de startup retrasa la visibilidad del splash |

## Plan de implementacion

1. Analizar el flujo actual: el splash depende de que Vite procese el SVG en `src="/src/assets/snorlax-splash.svg"` — esto no se resuelve hasta que el bundle se carga.
2. Mover el Snorlax SVG inline en el HTML o usar una version base64/data-uri para que aparezca sin esperar al bundler.
3. Mover los estilos criticos del splash (`#splash-screen`, `.splash-snorlax`, animaciones ZZZ) a un `<style>` inline en el `<head>` del HTML para que se apliquen antes de que cargue el CSS bundle.
4. Verificar que la transicion al estado interactivo sigue funcionando correctamente.
5. Si el SVG es muy grande para inline, pedir al usuario una imagen alternativa (PNG optimizado).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Abrir la app y verificar que Snorlax aparece en < 100ms, sin pantalla negra |

## Criterios de aceptacion

- [ ] Al abrir la app, el Snorlax aparece inmediatamente sin periodo de pantalla negra
- [ ] Las animaciones de ZZZ siguen funcionando correctamente
- [ ] La transicion interactiva (click en Snorlax) sigue funcionando
- [ ] Funciona tanto en la version desktop (Wails) como en la version web/APK

## Notas

El splash actual usa `src="/src/assets/snorlax-splash.svg"` que Vite procesa en build. En dev, el SVG se sirve directamente pero puede haber latencia. La solucion ideal es inline el SVG directamente en el HTML y mover los estilos criticos al head. Si el SVG es demasiado pesado, se pedira al usuario una imagen alternativa.
