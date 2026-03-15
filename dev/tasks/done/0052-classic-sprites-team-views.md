# Sprites pixel art clásico para vistas de equipos

**ID**: 0052-classic-sprites-team-views
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Cambiar los sprites de Pokémon en las vistas de **creación de equipos** y **gestión de equipos** (NO en la simulación de combate) para usar sprites pixel art clásico de pokemondb.net en lugar de los sprites 3D de Pokémon Home actuales.

**Cadena de fallback** (de más prioritario a menos):
1. **Black & White** (pixel art clásico, cubre Gen 1-5): `https://img.pokemondb.net/sprites/black-white/normal/{name}.png`
2. **X & Y** (estáticos, cubre Gen 1-6): `https://img.pokemondb.net/sprites/x-y/normal/{name}.png`
3. **Home** (último recurso, cubre todo): `https://img.pokemondb.net/sprites/home/normal/1x/{name}.png`

**Contexto**: actualmente `spriteURL()` apunta a `/assets/sprites/home-normal/{name}.png` (local) con fallback CDN a `home/normal`. Se necesita cambiar la fuente principal y la cadena de fallbacks para las vistas de equipo, manteniendo intactos los sprites de batalla.

**Vistas afectadas** (en `builds.ts`):
- `renderTeamMemberRow()` — sprites 32px en la lista de miembros del equipo
- Cards de Pokémon en el selector de builds (sprite 96px)
- Preview de equipos en la sección de team battle (`.tb-sprite`) — los sprites de preview, NO los de combate en sí
- Cualquier otra referencia a `spriteURL()` que NO sea batalla

**Vistas que NO se tocan**:
- `renderBattleSection()` — usa `battleSpriteURL()` con sprites Gen 1 back/front
- Sprites dentro de la simulación de combate individual y por equipos

## Capas afectadas
- **Frontend**: Modificar funciones de sprite y cadena de fallback en `builds.ts` y estilos en `_builds.scss`.
- **Core**: Sin cambios.
- **Shell**: Sin cambios.
- **APP**: Sin cambios.

## Dependencias externas nuevas
Ninguna. Se usan URLs del CDN de pokemondb.net que ya se usa como fallback en el proyecto.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | Cambiar `spriteURL()` y `spriteFallback()` para apuntar a sprites B&W clásicos con cadena de fallback B&W → X&Y → Home |
| `frontend/src/styles/_builds.scss` | modificar | Asegurar `image-rendering: pixelated` en `.build-sprite` y `.team-member-sprite` para que los sprites pixel art se vean nítidos al escalar |
| `frontend/src/pages/pokedex.ts` | verificar | Si usa `spriteURL` propia, aplicar mismo cambio de sprite clásico |
| `frontend/src/pages/types.ts` | verificar | Si usa `spriteURL` propia, aplicar mismo cambio de sprite clásico |
| `frontend/src/components/pokemon-type-modal.ts` | verificar | Si usa sprites, aplicar mismo cambio |

## Plan de implementacion

### Fase 1 — Actualizar funciones de sprite en builds.ts
1. Modificar `spriteURL(name)` para que apunte a sprites Black & White:
   ```typescript
   function spriteURL(name: string): string {
     const safeName = name.toLowerCase().replace(/[^a-z0-9-]/g, "");
     return `https://img.pokemondb.net/sprites/black-white/normal/${safeName}.png`;
   }
   ```
2. Modificar la cadena de fallback en los `onerror` de las imágenes para usar fallback escalonado:
   - Fallback 1: X & Y → `https://img.pokemondb.net/sprites/x-y/normal/${safeName}.png`
   - Fallback 2: Home → `https://img.pokemondb.net/sprites/home/normal/1x/${safeName}.png`
   - Fallback 3: Local home-normal (si existe) → `/assets/sprites/home-normal/${safeName}.png`
3. NO tocar `battleSpriteURL()` ni `battleSpriteFallbackCDN()` — esas funciones quedan intactas.

### Fase 2 — Asegurar renderizado pixelado
4. En `_builds.scss`, verificar que `.build-sprite` y `.team-member-sprite` tengan `image-rendering: pixelated` para que los sprites B&W se vean nítidos al escalar (no borrosos).

### Fase 3 — Aplicar a otras vistas
5. Revisar `pokedex.ts`, `types.ts` y `pokemon-type-modal.ts` — si tienen funciones `spriteURL` propias, aplicar el mismo cambio de fuente y cadena de fallback.

## Tests
| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que en creación de equipos se ven sprites pixel art B&W |
| Manual | Verificar que en la lista de miembros del equipo se ven sprites pixel art |
| Manual | Verificar que Pokémon de Gen 6+ hacen fallback a X&Y o Home correctamente |
| Manual | Verificar que los sprites de batalla (1v1 y team battle) NO cambiaron |
| Manual | Verificar que sprites se ven nítidos (pixelated) y no borrosos |

## Criterios de aceptacion
- [ ] `spriteURL()` apunta a sprites Black & White como fuente principal
- [ ] Cadena de fallback: B&W → X&Y → Home → local
- [ ] Sprites en creación de equipos muestran pixel art clásico
- [ ] Sprites en lista de miembros del equipo muestran pixel art clásico
- [ ] Sprites de batalla (1v1 y team battle simulation) NO están afectados
- [ ] `image-rendering: pixelated` aplicado a sprites de equipo
- [ ] Pokémon de Gen 6+ muestran fallback correcto (no imagen rota)
- [ ] Otras vistas (pokedex, types, modal) actualizadas si aplica

## Notas
- Los sprites B&W son ~96x96px, lo cual es adecuado para el tamaño actual de `.build-sprite` (96px) y `.team-member-sprite` (32px con escalado).
- La cadena de fallback usa `onerror` con dataset counter, patrón ya usado en el proyecto para sprites de batalla.
- Los sprites B&W cubren 649 Pokémon (Gen 1-5). X&Y añade Gen 6 (~721). Home cubre todas las generaciones.
