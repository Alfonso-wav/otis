# Migrar frontend a TypeScript + Vite

**ID**: 0003-frontend-typescript-vite
**Estado**: done
**Fecha**: 2026-03-13

---

## Descripcion

Migrar el frontend de vanilla JS a TypeScript con Vite como bundler. Este es el paso fundamental que habilita el resto del stack definido en `lenguajes.md`. Se prioriza compatibilidad con Wails v2 y que la app siga funcionando exactamente igual tras la migración.

## Capas afectadas

- **Core**: sin cambios
- **Shell**: sin cambios
- **APP**: sin cambios en Go. Se actualiza `wails.json` para incluir comandos de build/install de npm

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/package.json` | crear | Dependencias npm: typescript, vite |
| `frontend/tsconfig.json` | crear | Configuración TypeScript strict mode |
| `frontend/vite.config.ts` | crear | Config Vite compatible con Wails v2 |
| `frontend/index.html` | modificar | Apuntar script a `src/main.ts` (entry point Vite) |
| `frontend/src/main.ts` | crear | Migración de `app.js` a TypeScript con tipos |
| `frontend/src/types.ts` | crear | Interfaces TypeScript para Pokemon, Stat, Sprites, etc. |
| `frontend/app.js` | eliminar | Reemplazado por `src/main.ts` |
| `wails.json` | modificar | Añadir `frontend:install`: `npm install`, `frontend:build`: `npm run build` |

## Plan de implementacion

1. Crear `package.json` con dependencias: `typescript`, `vite`
2. Crear `tsconfig.json` con strict mode activado
3. Crear `vite.config.ts` con output compatible Wails (dist → frontend/dist)
4. Crear `frontend/src/types.ts` con interfaces que mapeen los tipos Go del core (Pokemon, PokemonType, Stat, Sprites, PokemonListResult, PokemonListItem)
5. Migrar `app.js` → `frontend/src/main.ts` con tipado estricto
6. Actualizar `index.html` para usar el entry point de Vite (`<script type="module" src="/src/main.ts">`)
7. Actualizar `wails.json` con comandos de install/build
8. Verificar que `wails dev` funciona correctamente
9. Eliminar `frontend/app.js`

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | `wails dev` arranca sin errores, la Pokédex funciona igual que antes |
| `tsc --noEmit` | El código TypeScript compila sin errores de tipo |

## Criterios de aceptacion

- [ ] `npm install` en frontend/ instala dependencias sin errores
- [ ] `npm run build` genera el bundle en frontend/dist
- [ ] `wails dev` arranca y la app funciona idéntica a antes
- [ ] No hay `any` explícito en el código TypeScript
- [ ] Los tipos TypeScript del frontend reflejan los tipos Go del core

## Notas

- Vite es el bundler recomendado por el documento y por Wails v2 para proyectos con build step.
- Los bindings auto-generados en `frontend/wailsjs/` ya tienen `.d.ts` — se aprovechan directamente.
- No se añade ningún framework UI en este paso. Solo TypeScript + Vite.
