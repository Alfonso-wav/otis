# Fix: scraper de sprites de batalla selecciona sprite B&W en vez de color

**ID**: 0038-fix-battle-sprites-selector
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

El scraper de sprites de batalla (`ScrapeBattleSpriteURLs`) selecciona el sprite incorrecto de pokemondb.net. Cada celda (`td`) de la tabla de sprites contiene **dos spans**:
- `span[1]`: sprite en blanco y negro
- `span[2]`: sprite a color

El scraper actual usa `cells.Eq(idx).Find("img").First()` que pilla el primer `img` (B&W). Debe tomar el segundo (color), que corresponde a `span[2]/a/img`.

**XPaths de referencia (pokemondb.net/sprites/{nombre})**:
- Defensor (Normal/frontal): `/html/body/main/div[11]/table/tbody/tr[2]/td[2]/span[2]/a/img`
- Atacante (Back/espaldas): `/html/body/main/div[11]/table/tbody/tr[2]/td[3]/span[2]/a/img`

Se toma de la fila de **Red & Blue**.

**No hay que re-descargar** los sprites existentes. Solo corregir el scraper para futuras descargas.

## Capas afectadas
- **Shell**: Modificar selector en `ScrapeBattleSpriteURLs`
- **Tests**: Actualizar HTML de test para reflejar la estructura real con dos spans por celda

## Dependencias externas nuevas
Ninguna.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `shell/pokemondb_sprites.go` | modificar | Cambiar selector de `Find("img").First()` a seleccionar el img dentro del segundo span (`span:nth-child(2) img` o equivalente) |
| `shell/pokemondb_sprites_test.go` | modificar | Actualizar HTML de test para incluir dos spans por celda (B&W + color) y verificar que se extrae el sprite a color |

## Plan de implementacion

### Paso 1 — Shell: corregir selector
1. En `shell/pokemondb_sprites.go`, función `ScrapeBattleSpriteURLs` (lineas 106-118):
   - Cambiar `cells.Eq(normalIdx).Find("img").First()` por un selector que tome el img del segundo span: `cells.Eq(normalIdx).Find("span:nth-child(2) a img").First()` (o `span + span img`, o `Find("span").Eq(1).Find("img")`)
   - Aplicar el mismo cambio para la columna Back

### Paso 2 — Tests: actualizar fixtures HTML
2. En `shell/pokemondb_sprites_test.go`, actualizar el HTML de test en `TestScrapeBattleSpriteURLs_Gen1` para:
   - Agregar dos spans por celda: el primero con sprite B&W, el segundo con sprite a color
   - Verificar que el scraper extrae la URL del sprite a color (segundo span)
3. Actualizar igualmente `TestScrapeBattleSpriteURLs_NoGen1_FallbackToOldest`

### Paso 3 — Verificar
4. Ejecutar tests: `go test ./shell/ -run TestScrapeBattleSpriteURLs -v`

## Tests
| Archivo | Que se testea |
|---------|---------------|
| `shell/pokemondb_sprites_test.go` | Que el scraper extrae el sprite a color (span[2]) y no el B&W (span[1]) |

## Criterios de aceptacion
- [ ] El scraper selecciona el img del segundo span (color) en vez del primero (B&W)
- [ ] Tests actualizados con HTML que refleja la estructura real (dos spans por celda)
- [ ] Tests pasan correctamente
- [ ] No se modifican sprites ya descargados
