# Task 0163 — Simulaciones: fondo de batalla configurable desde Explore > Mapas

## Estado: done

## Goal
En la visualizacion de batalla del simulador, permitir al usuario **cambiar el fondo** eligiendo entre los mapas disponibles en Explore > Mapas (Kanto, Johto, Hoenn, Sinnoh, Unova, Kalos, Alola, Galar, Paldea, Hisui, world-1..4, retro-blue).

## Contexto tecnico

### Frontend
- Vista batalla: `frontend/src/pages/builds.ts`, funcion `renderBattleSection()` (linea 519+). Arena: `<div class="battle-arena">` con `.battle-arena-top` y `.battle-arena-bottom` (lineas 533-541).
- Actualmente **no tiene fondo**; solo sprites + HP bars. Fondo debe aplicarse a `.battle-arena`.
- Assets disponibles: `frontend/public/assets/maps/` — `{alola,galar,hisui,hoenn,johto,kalos,kanto,paldea,retro-blue,sinnoh,unova,world-1..4}.{png,jpg}`.
- Lista de mapas definida en `frontend/src/pages/explore/maps.ts` (lineas 19-30 regions, 238-243 worlds, 288 retro).

### UX propuesta
- Selector (`<select>`) encima de la arena o en zona de controles. Placeholder disabled: "-- Fondo --".
- Opciones: "Sin fondo" (default) + lista de mapas con nombre traducido (reusar labels de `maps.ts`).
- Aplicar via CSS: `style="background-image:url('/assets/maps/kanto.png'); background-size:cover; background-position:center"` en `.battle-arena`.
- Preferencia persistida en `localStorage` (key `battle-background`) — como otras prefs del proyecto (sort orders, dark mode).
- Recargar estilo en re-render de la vista.

### Consideraciones
- Mapas son grandes (jpg/png full). Verificar que no tapan sprites — sprites deben tener z-index + sombra/contorno si hace falta contraste.
- `<select>` con placeholder disabled (regla CLAUDE.md).
- i18n: labels de regions ya existen en `regions.*`; worlds y retro-blue pueden necesitar keys nuevas.
- Mobile (360px): selector no debe romper fila de botones.

## Acceptance criteria

- [ ] Selector de fondo visible en la vista batalla (idle y durante batalla).
- [ ] Opciones: "Sin fondo" + todos los mapas de Explore > Maps.
- [ ] Cambio de fondo se aplica inmediatamente a `.battle-arena`.
- [ ] Preferencia se guarda en `localStorage` y persiste entre sesiones.
- [ ] Sprites visibles y legibles sobre el fondo (contraste OK).
- [ ] `<select>` tiene placeholder disabled.
- [ ] Dark mode OK.
- [ ] Mobile 360px OK (selector no desborda).
- [ ] i18n EN + ES para nuevas keys.
- [ ] Funciona en build de produccion y APK (assets en `public/`, ruta `/assets/maps/...`).

## Archivos afectados

- `frontend/src/pages/builds.ts` — render selector, handler change, aplicar `background-image` a `.battle-arena`, leer/escribir `localStorage`.
- `frontend/src/styles/_builds.scss` o similar — estilos base `.battle-arena` (background-size, position), selector.
- `frontend/src/locales/{en,es}.json` — keys: `builds.battleBackground`, `builds.noBackground`, y posibles labels de worlds/retro si no existen.

## Notas

- No tocar backend ni core. Solo frontend.
- Reutilizar lista de mapas de `maps.ts` exportandola (evitar duplicacion).
