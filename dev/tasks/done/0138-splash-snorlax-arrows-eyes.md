# Flechas amarillas y ojos SVG en splash Snorlax

**ID**: 0138-splash-snorlax-arrows-eyes
**Estado**: todo
**Fecha**: 2026-04-11
**Depende de**: 0137-fix-splash-black-screen

---

## Descripcion

Mejorar la interaccion del splash screen de Snorlax con dos nuevas animaciones:

1. **Flechas amarillas**: 2 segundos despues de que aparezca el Snorlax, mostrar flechas amarillas animadas a su alrededor que le senalan. Aleatoriamente se muestran 1, 2 o 3 flechas en diferentes posiciones (arriba, abajo, izquierda, derecha, esquinas). Las flechas aparecen y desaparecen en bucle hasta que el usuario haga click en Snorlax.

2. **Ojos SVG abriendose**: Una vez el usuario clicke en Snorlax, dibujar con SVG unos ojos que se van abriendo progresivamente encima de los ojos cerrados de Snorlax, dando la sensacion de que se despierta. Despues de la animacion de ojos, continuar con la animacion de salida actual.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — splash screen interactivo

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/main.ts` | modificar | Agregar logica de flechas con delay de 2s y animacion de ojos SVG al click |
| `frontend/src/styles/_splash.scss` | modificar | Estilos para flechas amarillas y ojos SVG |
| `frontend/index.html` | modificar | Agregar contenedores SVG para ojos si es necesario |

## Plan de implementacion

1. Crear flechas amarillas como elementos SVG o HTML/CSS (triangulos/chevrons amarillos).
2. Posicionar las flechas alrededor del Snorlax en posiciones predefinidas (8 posiciones: N, NE, E, SE, S, SW, W, NW).
3. Implementar un intervalo que cada ~1s selecciona aleatoriamente 1-3 posiciones y muestra/oculta flechas con animacion fade+traslacion.
4. Iniciar las flechas 2 segundos despues de que el splash sea visible.
5. Al click en Snorlax: detener flechas, mostrar SVG de ojos abriendose (clipPath o animacion de altura/escala del parpado).
6. Los ojos SVG deben posicionarse sobre los ojos cerrados del Snorlax (coordenadas relativas al sprite).
7. Animacion de ojos: apertura gradual ~1s, luego continuar con la animacion de salida existente.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que las flechas aparecen 2s despues, que son aleatorias (1-3), que desaparecen al click |
| Manual | Verificar que los ojos SVG se abren correctamente sobre los ojos de Snorlax |
| Manual | Verificar que la animacion de salida sigue funcionando despues de los ojos |

## Criterios de aceptacion

- [ ] 2 segundos despues del splash, aparecen flechas amarillas senalando al Snorlax
- [ ] Las flechas son aleatorias: entre 1 y 3 simultaneas en posiciones diferentes
- [ ] Las flechas se repiten en bucle hasta el click
- [ ] Al clickar Snorlax, las flechas desaparecen
- [ ] Al clickar, ojos SVG se dibujan sobre los ojos de Snorlax y se abren progresivamente
- [ ] Despues de la animacion de ojos, la animacion de salida actual se ejecuta normalmente
- [ ] Funciona en desktop, web y APK

## Notas

Para los ojos SVG, una tecnica es usar `clipPath` con una elipse cuyo `ry` crece de 0 a su valor final, simulando que el ojo se abre. Los ojos de Snorlax en el SVG actual estan cerrados (lineas curvas). Hay que posicionar los ojos nuevos exactamente encima usando coordenadas relativas al contenedor del sprite.
