# Eliminar flechas este y oeste del splash de Snorlax

**ID**: 0146-remove-snorlax-arrows-east-west
**Estado**: done
**Fecha**: 2026-04-12

---

## Descripcion

Las flechas amarillas que señalan al Snorlax en la pantalla splash antes de entrar a la app deben reducirse de 8 a 6. Se eliminan las flechas en posiciones 3 y 9 del reloj, es decir:
- **E (este / derecha)** — posición 3
- **W (oeste / izquierda)** — posición 9

Las 6 flechas que permanecen son: N, NE, SE, S, SW, NW.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — main.ts (lógica de flechas), estilos _splash.scss

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/main.ts` | modificar | Eliminar "e" y "w" del array `ARROW_POSITIONS` |
| `frontend/src/styles/_splash.scss` | modificar | Eliminar las reglas CSS para `&--e` y `&--w` |

## Plan de implementacion

1. En `main.ts`, cambiar `ARROW_POSITIONS` de `["n", "ne", "e", "se", "s", "sw", "w", "nw"]` a `["n", "ne", "se", "s", "sw", "nw"]`.
2. En `_splash.scss`, eliminar los bloques de estilo para `.splash-arrow--e` y `.splash-arrow--w`.
3. Verificar que la animación de flechas aleatorias sigue funcionando con 6 posiciones.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que solo aparecen 6 flechas (no hay flechas laterales E/W) |
| Manual | Verificar que las flechas siguen animándose aleatoriamente |
| Manual | Verificar que el click en Snorlax sigue funcionando |

## Criterios de aceptacion

- [x] No aparecen flechas en posiciones este (3h) ni oeste (9h)
- [x] Las 6 flechas restantes (N, NE, SE, S, SW, NW) se animan correctamente
- [x] El flujo completo del splash sigue funcionando (click → jiggle → eyes → fade)
