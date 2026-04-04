# Mejorar espaciado de botones en simulador de batallas

**ID**: 0105-battle-sim-buttons-spacing
**Estado**: done
**Fecha**: 2026-04-04

---

## Descripcion

Los botones "Simular batalla de equipos" y "Simular N batallas" en la seccion de simulador de batallas estan demasiado pegados a los bordes de sus contenedores. Les falta padding/margen interno para que respiren visualmente y se vean mas esteticos.

Hay dos zonas afectadas:

1. **Batalla de equipos** (`tb-actions`): El contenedor `.tb-actions` (linea 1201 de `builds.ts`) no tiene **ningun estilo CSS definido** en `_builds.scss`. Los botones `.battle-auto-btn` y `.battle-batch-btn` quedan pegados a los bordes del contenedor.

2. **Batalla individual** (`battle-idle-btns` + `battle-batch-row`): Los contenedores `.battle-idle-btns` (linea 483 de `_builds.scss`) y `.battle-batch-row` (linea 1019) tienen margen superior pero carecen de padding lateral suficiente para separarse de los bordes del modulo padre.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend - estilos SCSS de la pagina builds

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/styles/_builds.scss` | modificar | Crear regla para `.tb-actions` con padding y gap adecuados para dar aire a los botones |
| `frontend/src/styles/_builds.scss` | modificar | Revisar `.battle-idle-btns` (linea 483) y `.battle-batch-row` (linea 1019): anadir padding lateral para separar botones de los bordes del modulo |

## Criterios de aceptacion

- [x] Los botones de simulacion tienen espacio visual respecto a los bordes de sus contenedores
- [x] El espaciado es coherente entre la seccion de batalla individual y batalla de equipos
- [x] No se rompe el layout en pantallas moviles
- [x] El aspecto mejora tanto en modo claro como oscuro
