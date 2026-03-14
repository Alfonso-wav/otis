# STAB multiplicador visible en columna de efectividad

**ID**: 0030-stab-in-effectiveness-column
**Estado**: done
**Fecha**: 2026-03-14

---

## Descripcion
Cuando un Pokemon tiene STAB (su tipo coincide con el tipo del movimiento), el multiplicador x1.5 no se refleja en la columna "Efectividad" de la tabla de dano. Actualmente solo muestra la efectividad de tipo (x1, x2, x0.5, etc.) y el STAB aparece como badge separado. El usuario espera ver el multiplicador combinado: si hay STAB + neutral = x1.5, si hay STAB + super eficaz x2 = x3, etc.

## Capas afectadas
- **Core**: ninguna (el calculo ya es correcto, STAB ya se aplica al dano final)
- **Shell**: ninguna
- **APP**: modificar la funcion `effectLabel` en el frontend para incluir el multiplicador STAB en el texto de efectividad

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | Actualizar `effectLabel()` para multiplicar `result.multiplier` por `result.stabMultiplier` (o por 1.5 si `result.hasSTAB`) al mostrar el texto de efectividad |

## Plan de implementacion
1. En `effectLabel()` (builds.ts:108-113), calcular el multiplicador combinado: `result.multiplier * (result.hasSTAB ? result.stabMultiplier : 1)`
2. Usar ese multiplicador combinado en los textos de efectividad
3. Verificar visualmente que la columna muestra x1.5 para STAB neutral, x3 para STAB + super eficaz, etc.

## Tests
| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que un Pokemon fuego con ataque fuego contra neutral muestra x1.5 |
| Manual | Verificar que STAB + super eficaz muestra x3 |
| Manual | Verificar que sin STAB sigue mostrando x1, x2, etc. normalmente |

## Criterios de aceptacion
- [x] La columna Efectividad muestra x1.5 cuando hay STAB contra tipo neutral
- [x] La columna Efectividad muestra el multiplicador combinado (tipo * STAB) en todos los casos
- [x] Sin STAB, el comportamiento no cambia
- [x] El badge STAB sigue apareciendo junto al nombre del movimiento

## Notas
El campo `Multiplier` en `DamageResult` solo contiene efectividad de tipo. El campo `STABMultiplier` contiene 1.5 o 1.0. La solucion es puramente frontend: combinar ambos valores al renderizar.
