# Fix: endpoint HTTP faltante para distribución de tipos por región

**ID**: 0150-fix-region-donut-missing-endpoint
**Estado**: todo
**Fecha**: 2026-04-12

---

## Descripción

Los gráficos donut de distribución de tipos en Explorar > Regiones muestran "Sin datos de distribución" en la versión mobile/HTTP. La causa raíz es que el endpoint REST `GET /api/regions/{region}/type-distribution` **no fue registrado** en `app/mobile/handlers.go` cuando se añadió `GetRegionTypeDistribution` al backend en el commit d2f4da7.

En desktop (Wails IPC) el binding funciona correctamente, pero en mobile el frontend llama a `/api/regions/{region}/type-distribution`, no hay handler, la petición falla, el `.catch()` en `regions.ts:142` pasa `{}` a `renderTypeDistributionChart`, y la función muestra el mensaje de "Sin datos de distribución" (`type-distribution.ts:19-22`).

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: `app/mobile/handlers.go` — añadir handler HTTP faltante

## Archivos a crear/modificar

| Archivo | Acción | Descripción |
|---------|--------|-------------|
| `app/mobile/handlers.go` | modificar | Añadir handler `GET /api/regions/{region}/type-distribution` que invoque `a.GetRegionTypeDistribution()` |
| `app/mobile/handlers.go` | modificar | Actualizar comentario de documentación (línea 16) para incluir el nuevo endpoint |

## Plan de implementación

1. En `app/mobile/handlers.go`, después del handler de `GetRegionPokemonByType` (línea ~226), añadir:
   ```go
   mux.HandleFunc("GET /api/regions/{region}/type-distribution", func(w http.ResponseWriter, r *http.Request) {
       result, err := a.GetRegionTypeDistribution(r.PathValue("region"))
       if err != nil {
           jsonError(w, err.Error(), http.StatusInternalServerError)
           return
       }
       jsonResponse(w, result)
   })
   ```
2. Actualizar el bloque de documentación al inicio del archivo para incluir:
   ```
   // GetRegionTypeDistribution → GET /api/regions/{region}/type-distribution
   ```
3. Verificar que el donut se renderiza correctamente en mobile.

## Tests

| Archivo | Qué se testea |
|---------|---------------|
| Manual | Expandir cada región en mobile y verificar que aparece el donut |
| Manual | Click en un tipo del donut abre el modal con la lista de Pokémon |
| Manual | Verificar en dark mode |
| Manual | Verificar que desktop (Wails) sigue funcionando sin regresión |

## Criterios de aceptación

- [ ] El endpoint `GET /api/regions/{region}/type-distribution` responde con `map[string]int`
- [ ] Al expandir una región en mobile, el gráfico donut se renderiza correctamente
- [ ] El donut es interactivo (click en tipo abre modal)
- [ ] Funciona en light mode y dark mode
- [ ] Sin regresión en desktop (Wails IPC)
