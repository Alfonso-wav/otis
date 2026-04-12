# Animación de Heracross en la carga de Bayas

**ID**: 0149-berry-loading-heracross
**Estado**: todo
**Fecha**: 2026-04-12

---

## Descripcion

Añadir una animación de Heracross como indicador de carga mientras se cargan los datos de las bayas en Explorer > Bayas. Actualmente la carga muestra un simple texto `<p class="loading">`. Se debe reemplazar por una animación similar al Diglett de carga pero usando a Heracross.

Heracross es un Pokemon tipo Bug/Fighting conocido por recoger savia/bayas de los árboles (headbutting trees), por lo que encaja temáticamente con la sección de bayas.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — pages/explore/berries.ts, componentes, estilos

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/explore/berries.ts` | modificar | Reemplazar el texto de carga por la animación de Heracross |
| `frontend/src/components/sorting-overlay.ts` | modificar | Añadir función `createInlineHeracross()` o similar |
| `frontend/src/styles/_components.scss` | modificar | Estilos para la animación de Heracross |
| `frontend/src/styles/_dark.scss` | modificar | Soporte dark mode si es necesario |

## Plan de implementacion

1. Definir el sprite de Heracross (artwork oficial de PokeAPI): `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/214.png`
2. Crear una función `createInlineHeracross(container, text?)` en `sorting-overlay.ts` (o un archivo nuevo) que:
   - Muestre el sprite de Heracross con una animación temática (ej. headbutt/shake, como si golpeara un árbol para sacar bayas).
   - Incluya texto de carga debajo.
3. En `berries.ts`, reemplazar el `<p class="loading">` en `initBerries()` por `createInlineHeracross(content, t("berries.loading"))`.
4. Diseñar la animación GSAP: posible secuencia de Heracross inclinándose hacia adelante (headbutt) y volviendo, con loop.
5. Añadir estilos correspondientes.
6. Limpiar la animación al terminar la carga (similar a como se hace con el inline Diglett).

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Ir a Explorer > Bayas y verificar que Heracross aparece durante la carga |
| Manual | Verificar que la animación se reproduce en loop mientras carga |
| Manual | Verificar que la animación desaparece al terminar la carga |
| Manual | Verificar en dark mode |
| Manual | Verificar en mobile |

## Criterios de aceptacion

- [ ] Al cargar las bayas aparece una animación de Heracross en lugar del texto plano
- [ ] La animación es fluida y temática (headbutt o similar)
- [ ] La animación desaparece al completar la carga de datos
- [ ] Funciona en dark mode y light mode
- [ ] Funciona en desktop y mobile
- [ ] No afecta al rendimiento de la carga de bayas
