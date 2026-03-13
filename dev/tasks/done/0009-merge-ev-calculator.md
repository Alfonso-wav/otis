# Merge PR #1: EV Calculator feature

**ID**: 0009-merge-ev-calculator
**Estado**: done
**Fecha**: 2026-03-13

---

## Descripcion

Mergear la PR #1 (`gutierenmanuel:feature/ev-calculator`) en `main`. La PR añade una calculadora de EVs integrada en la vista de detalle de Pokémon usando las fórmulas oficiales Gen III+. Desde que se abrió la PR, `main` recibió el sistema de tab navigation (commit `27d06ad`), que refactorizó `main.ts` y `main.go` de forma incompatible con los cambios de la PR. Hay 2 conflictos manuales y 6 archivos con diferencias aditivas que git puede resolver automáticamente.

## Capas afectadas

- **Core**: `core/domain.go` (tipos Stats, Nature, EV*), `core/ev_calc.go` (lógica de cálculo pura)
- **Shell**: ninguna
- **APP**: `app/bindings.go` (IPC bindings GetNatures, CalculateEVs, CalculateStats), `main.go` (embed de assets)

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `core/domain.go` | modificar | Añadir tipos Stats, Nature, StatRange, EVCalculatorInput, EVCalculatorResult, StatCalculatorInput |
| `core/ev_calc.go` | crear | Lógica pura de cálculo de EVs e IVs (nuevo archivo de la PR) |
| `app/bindings.go` | modificar | Añadir métodos GetNatures, CalculateEVs, CalculateStats |
| `main.go` | modificar | **CONFLICTO MANUAL** — reconciliar embed de assets |
| `frontend/src/ev-calculator.ts` | crear | UI del calculador de EVs (nuevo archivo de la PR) |
| `frontend/src/main.ts` | modificar | **CONFLICTO MANUAL** — integrar EV calculator en la arquitectura de páginas/router |
| `frontend/src/types.ts` | modificar | Añadir interfaces Stats, Nature, StatRange, EVCalculatorInput, EVCalculatorResult |
| `frontend/src/styles/_pokemon.scss` | crear | Estilos del calculador (nuevo archivo de la PR) |
| `frontend/wailsjs/go/app/App.d.ts` | modificar | Añadir declaraciones GetNatures, CalculateEVs, CalculateStats |
| `frontend/wailsjs/go/app/App.js` | modificar | Añadir bindings JS para los nuevos métodos IPC |
| `frontend/wailsjs/go/models.ts` | modificar | Añadir clases Stats, EVCalculatorInput, StatRange, EVCalculatorResult |
| `docs/ev-calculator-research.md` | crear | Documentación de investigación (nuevo archivo de la PR) |

## Plan de implementacion

1. Actualizar `main` local: `git pull origin main`

2. Obtener los archivos únicos de la PR sin conflicto directamente:
   ```
   git fetch origin pull/1/head:pr-ev-calculator
   git checkout pr-ev-calculator -- core/ev_calc.go
   git checkout pr-ev-calculator -- frontend/src/ev-calculator.ts
   git checkout pr-ev-calculator -- frontend/src/styles/_pokemon.scss
   git checkout pr-ev-calculator -- docs/ev-calculator-research.md
   ```

3. Aplicar los cambios aditivos en archivos ya existentes (copiar desde la PR):
   - `core/domain.go`: añadir al final del archivo los tipos `Stats`, `Nature`, `StatRange`, `EVCalculatorInput`, `EVCalculatorResult`, `StatCalculatorInput` (55 líneas de la PR).
   - `app/bindings.go`: añadir al final los métodos `GetNatures`, `CalculateEVs`, `CalculateStats` (82 líneas de la PR).
   - `frontend/src/types.ts`: añadir al final las interfaces `Stats`, `Nature`, `StatRange`, `EVCalculatorInput`, `EVCalculatorResult` (42 líneas de la PR).
   - `frontend/wailsjs/go/app/App.d.ts`: añadir las declaraciones de los 3 nuevos métodos IPC.
   - `frontend/wailsjs/go/app/App.js`: añadir las funciones JS de los 3 nuevos métodos IPC.
   - `frontend/wailsjs/go/models.ts`: añadir las clases `Stats`, `EVCalculatorInput`, `StatRange`, `EVCalculatorResult`.

4. **Conflicto manual — `main.go`**: La PR cambió el embed para usar `fs.Sub(assets, "frontend/dist")`. Main ya usa `//go:embed frontend` con `assets` directo. Mantener la versión de `main` (embed directo) y verificar que el asset server sigue funcionando. No hay que cambiar nada aquí; el approach de `main` es más simple y ya funciona.

5. **Conflicto manual — `frontend/src/main.ts`**: La PR añadía la integración del EV calculator en el antiguo `main.ts` monolítico. Main refactorizó `main.ts` para usar el sistema de router/páginas. La solución es NO usar el `main.ts` de la PR. En su lugar, integrar el EV calculator dentro de la página correcta del router:
   - Identificar en qué página se muestra el detalle de Pokémon (probablemente `frontend/src/pages/pokedex.ts`).
   - Importar `renderEVCalculatorForm` e `initEVCalculator` desde `ev-calculator.ts` en esa página.
   - Llamar a `initEVCalculator` al mostrar el detalle de un Pokémon.

6. Compilar el frontend: `cd frontend && npm run build`

7. Compilar el backend: `go build ./...`

8. Verificar que la app arranca y la calculadora de EVs aparece en el detalle de un Pokémon.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| `core/ev_calc_test.go` | Si existe, ejecutar. Si no, verificar manualmente que `CalculateEVs` devuelve rangos coherentes para un Pokémon conocido. |
| build + runtime | `go build ./...` compila sin errores; app arranca; calculadora visible en detalle |

## Criterios de aceptacion

- [ ] `go build ./...` sin errores
- [ ] `cd frontend && npm run build` sin errores
- [ ] App arranca correctamente con el sistema de tabs
- [ ] La calculadora de EVs aparece en la vista de detalle de Pokémon
- [ ] `CalculateEVs` devuelve resultados coherentes desde el frontend
- [ ] No hay regresiones en la navegación por tabs

## Notas

**Conflictos detectados:**

| Archivo | Tipo | Detalle |
|---------|------|---------|
| `main.go` | automático (mantener main) | PR usa `fs.Sub` + `//go:embed frontend/dist`; main usa embed directo. Mantener versión de main, no requiere cambio. |
| `frontend/src/main.ts` | manual | PR integra EV calc en el viejo monolito; main lo refactorizó a router+páginas. Hay que portar la integración a `pages/pokedex.ts`. |

**Archivos con diferencias aditivas (sin conflicto real):** `core/domain.go`, `app/bindings.go`, `frontend/src/types.ts`, `frontend/wailsjs/go/app/App.d.ts`, `frontend/wailsjs/go/app/App.js`, `frontend/wailsjs/go/models.ts`.

**Archivos solo en PR (copiar directamente):** `core/ev_calc.go`, `frontend/src/ev-calculator.ts`, `frontend/src/styles/_pokemon.scss`, `docs/ev-calculator-research.md`.

La rama de la PR pertenece a un fork externo (`gutierenmanuel`), por lo que no se puede hacer merge directo con `gh pr merge`. Los cambios se portan manualmente siguiendo este plan y luego se cierra la PR.


Haz que en github aparezca como si hubiese mergeado su rama con exito. 
