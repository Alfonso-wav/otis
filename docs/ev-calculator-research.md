# Investigación: Calculadora de EVs para Pokémon

## Índice
1. [Conceptos Fundamentales](#conceptos-fundamentales)
2. [Fórmulas de Cálculo de Stats](#fórmulas-de-cálculo-de-stats)
3. [Effort Values (EVs)](#effort-values-evs)
4. [Individual Values (IVs)](#individual-values-ivs)
5. [Naturalezas](#naturalezas)
6. [Cálculo Inverso de EVs](#cálculo-inverso-de-evs)
7. [Datos Mínimos Requeridos del Usuario](#datos-mínimos-requeridos-del-usuario)
8. [Diseño de la Calculadora](#diseño-de-la-calculadora)
9. [Fuentes](#fuentes)

---

## Conceptos Fundamentales

Las estadísticas finales de un Pokémon se calculan a partir de cinco componentes:

| Componente | Rango | Descripción |
|------------|-------|-------------|
| **Base Stats** | 1-255 | Fijos por especie (ej: Pikachu siempre tiene 90 de Speed base) |
| **IVs** | 0-31 | Valores individuales, fijos al capturar/eclosionar |
| **EVs** | 0-252 | Effort Values, ganados entrenando |
| **Nivel** | 1-100 | Nivel actual del Pokémon |
| **Naturaleza** | ×0.9/1.0/1.1 | Modificador de +10%/-10% a dos stats |

---

## Fórmulas de Cálculo de Stats

### Generación III en adelante (Estándar actual)

#### HP (Puntos de Salud)
```
HP = ⌊(2 × Base + IV + ⌊EV/4⌋) × Nivel / 100⌋ + Nivel + 10
```

**Caso especial:** Shedinja siempre tiene 1 HP.

#### Otras Stats (Ataque, Defensa, Sp.Atk, Sp.Def, Velocidad)
```
Stat = ⌊(⌊(2 × Base + IV + ⌊EV/4⌋) × Nivel / 100⌋ + 5) × Naturaleza⌋
```

Donde:
- `⌊x⌋` = floor (redondear hacia abajo)
- `Base` = stat base de la especie
- `IV` = valor individual (0-31)
- `EV` = effort value (0-252)
- `Nivel` = nivel del Pokémon (1-100)
- `Naturaleza` = 1.1 (aumenta), 1.0 (neutral), 0.9 (disminuye)

### Ejemplo de Cálculo

**Pikachu nivel 50, IV=31, EV=252, Naturaleza Timid (+Speed)**

Stats base de Pikachu: Speed = 90

```
Speed = ⌊(⌊(2 × 90 + 31 + ⌊252/4⌋) × 50 / 100⌋ + 5) × 1.1⌋
Speed = ⌊(⌊(180 + 31 + 63) × 50 / 100⌋ + 5) × 1.1⌋
Speed = ⌊(⌊274 × 0.5⌋ + 5) × 1.1⌋
Speed = ⌊(137 + 5) × 1.1⌋
Speed = ⌊142 × 1.1⌋
Speed = ⌊156.2⌋
Speed = 156
```

---

## Effort Values (EVs)

### Límites

| Generación | Máx. por Stat | Total Máximo |
|------------|---------------|--------------|
| Gen I-II | 65,535 | Sin límite |
| Gen III-V | 255 | 510 |
| Gen VI+ | 252 | 510 |

### Impacto en Stats

- **A nivel 100:** +1 punto de stat por cada 4 EVs
- **Máximo beneficio:** +63 puntos (252 EVs ÷ 4)
- **A nivel 50:** La relación varía debido al floor en la fórmula

### Peculiaridad a Nivel 50 (Competitivo)

En batallas competitivas (nivel 50), con 31 IVs:
- Primeros 4 EVs → +1 stat
- Cada 8 EVs adicionales → +1 stat

Esto explica spreads como 252/4/0/252/0/0 en vez de 255/0/0/255/0/0.

### Fuentes de EVs

| Método | EVs Ganados |
|--------|-------------|
| Combate | 1-3 según especie derrotada |
| Vitaminas | +10 EVs (con límites por gen) |
| Alas/Plumas | +1 EV |
| Pokérus | ×2 EVs de combate |
| Macho Brace | ×2 EVs de combate |
| Power Items | +8 EVs adicionales |

---

## Individual Values (IVs)

### Características

- Rango: 0-31 por stat
- **Fijos permanentemente** al capturar o eclosionar
- Determinan el potencial máximo del Pokémon
- En juegos modernos: Hyper Training permite simular IVs de 31

### Cálculo Inverso de IVs

Dado que la fórmula usa `floor()`, calcular IVs exactos requiere:
1. Pokémon **recién capturado** (EVs = 0)
2. Conocer el **nivel exacto**
3. Conocer la **naturaleza**
4. Usar múltiples niveles para reducir el rango

La fórmula inversa para IVs (asumiendo EVs conocidos):

```
IV = ⌈((Stat / Naturaleza - 5) × 100 / Nivel) - 2 × Base - ⌊EV/4⌋⌉
```

---

## Naturalezas

### Tabla Completa de las 25 Naturalezas

| Naturaleza | Stat +10% | Stat -10% |
|------------|-----------|-----------|
| **Adamant** | Attack | Sp. Attack |
| **Bashful** | — | — (neutral) |
| **Bold** | Defense | Attack |
| **Brave** | Attack | Speed |
| **Calm** | Sp. Defense | Attack |
| **Careful** | Sp. Defense | Sp. Attack |
| **Docile** | — | — (neutral) |
| **Gentle** | Sp. Defense | Defense |
| **Hardy** | — | — (neutral) |
| **Hasty** | Speed | Defense |
| **Impish** | Defense | Sp. Attack |
| **Jolly** | Speed | Sp. Attack |
| **Lax** | Defense | Sp. Defense |
| **Lonely** | Attack | Defense |
| **Mild** | Sp. Attack | Defense |
| **Modest** | Sp. Attack | Attack |
| **Naive** | Speed | Sp. Defense |
| **Naughty** | Attack | Sp. Defense |
| **Quiet** | Sp. Attack | Speed |
| **Quirky** | — | — (neutral) |
| **Rash** | Sp. Attack | Sp. Defense |
| **Relaxed** | Defense | Speed |
| **Sassy** | Sp. Defense | Speed |
| **Serious** | — | — (neutral) |
| **Timid** | Speed | Attack |

### Identificación en el Juego

- **Texto rojo:** stat aumentado (+10%)
- **Texto azul:** stat disminuido (-10%)

---

## Cálculo Inverso de EVs

### Fórmula para Determinar EVs desde Stats

Despejando EV de la fórmula de stats:

#### Para HP:
```
EV = (((HP - Nivel - 10) × 100 / Nivel) - 2 × Base - IV) × 4
```

#### Para otras stats:
```
EV = (((Stat / Naturaleza - 5) × 100 / Nivel) - 2 × Base - IV) × 4
```

### Problema del Floor

Debido a `⌊⌋` en la fórmula original, el cálculo inverso da un **rango de valores posibles**, no un valor exacto.

### EVs Necesarios para +N puntos de stat

```
EVs_necesarios = ⌈(N × 100 / Nivel) / (Naturaleza)⌉ × 4
```

A nivel 100 simplificado:
```
EVs_necesarios = N × 4
```

---

## Datos Mínimos Requeridos del Usuario

Para que la calculadora funcione con **mínimo input del jugador**, necesitamos:

### Datos Esenciales (Mínimos)

| Dato | Cómo obtenerlo | Necesario |
|------|----------------|-----------|
| **Especie** | Ya la tenemos de la Pokédex | Automático |
| **Nivel** | Visible en resumen | Sí |
| **Naturaleza** | Visible en resumen (color rojo/azul) | Sí |
| **Stats actuales** | Visible en resumen | Sí |

### Datos Opcionales (Mejoran precisión)

| Dato | Cómo obtenerlo | Por defecto |
|------|----------------|-------------|
| **IVs** | Juez de IVs o calculadora | Asumir 31 (perfecto) |
| **EVs conocidos** | Si el jugador los entrenó | Asumir 0 |

### Flujo de Usuario Simplificado

1. Usuario selecciona Pokémon en la Pokédex
2. Ingresa: **Nivel**, **Naturaleza**, **Stats actuales** (6 valores)
3. Opcionalmente: IVs conocidos
4. La calculadora muestra:
   - EVs estimados por stat
   - Stats máximos posibles
   - EVs restantes disponibles (510 - usados)
   - Recomendaciones de entrenamiento

---

## Diseño de la Calculadora

### Funcionalidades Principales

#### 1. Calculadora de Stats (Forward)
- **Input:** Base stats, IVs, EVs, Nivel, Naturaleza
- **Output:** Stats finales calculados

#### 2. Calculadora de EVs (Reverse)
- **Input:** Stats actuales, Nivel, Naturaleza, IVs (asumidos o conocidos)
- **Output:** EVs estimados (con rango de error)

#### 3. Optimizador de EVs
- **Input:** Stats deseados, restricciones
- **Output:** Distribución óptima de 510 EVs

### Estructura de Datos Propuesta

```go
type EVCalculatorInput struct {
    PokemonID int      // Para obtener base stats
    Level     int      // 1-100
    Nature    string   // "Adamant", "Jolly", etc.
    CurrentStats Stats // Stats actuales del usuario
    KnownIVs  *Stats   // Opcional, nil = asumir 31
    KnownEVs  *Stats   // Opcional, nil = calcular
}

type Stats struct {
    HP         int
    Attack     int
    Defense    int
    SpAttack   int
    SpDefense  int
    Speed      int
}

type EVCalculatorResult struct {
    EstimatedEVs    Stats
    PossibleEVRange map[string][2]int // min-max por stat
    TotalEVsUsed    int
    EVsRemaining    int
    MaxPossibleStats Stats
    Recommendations []string
}
```

### Arquitectura en el Proyecto

```
core/
├── domain.go      # Añadir tipos Stats, Nature, EVCalcInput/Result
├── ev_calc.go     # Lógica de cálculo de EVs/Stats
└── pokemon.go     # Ya existente

app/
├── bindings.go    # Añadir CalculateEVs(), GetNatures()
└── ...

shell/
├── pokeapi.go     # Ya obtiene base stats
└── ...

frontend/
├── src/
│   ├── ev-calculator.ts  # Lógica UI de la calculadora
│   └── ...
└── ...
```

---

## Fuentes

- [Bulbapedia - Stat](https://bulbapedia.bulbagarden.net/wiki/Stat)
- [Bulbapedia - Effort Values](https://bulbapedia.bulbagarden.net/wiki/Effort_values)
- [Bulbapedia - Nature](https://bulbapedia.bulbagarden.net/wiki/Nature)
- [Pokémon Database - Natures](https://pokemondb.net/mechanics/natures)
- [Pokémon Database - EVs](https://pokemondb.net/ev)
- [Marriland - IV Calculator](https://marriland.com/tools/iv-calculator/)
- [Terresquall - Calculating EVs](https://blog.terresquall.com/2020/07/calculating-evs-needed-to-raise-a-stat-in-pokemon/)
- [PokéAPI Documentation](https://pokeapi.co/docs/v2)
- [Serebii - Effort Values](https://www.serebii.net/games/evs.shtml)
- [VGC Guide - Base Stats](https://www.vgcguide.com/base-stats)
