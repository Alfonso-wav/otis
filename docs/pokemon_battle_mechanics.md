# Mecánicas de Combate en Pokémon — Guía Exhaustiva

> Documento técnico que cubre absolutamente todos los factores que influyen en el cálculo de daño, resistencias y multiplicadores en las batallas Pokémon (Generaciones I–IX).

---

## Índice

1. [La Fórmula de Daño](#1-la-fórmula-de-daño)
2. [Estadísticas Base, IVs y EVs](#2-estadísticas-base-ivs-y-evs)
3. [Naturalezas (Natures)](#3-naturalezas-natures)
4. [Etapas de Estadísticas (Stat Stages)](#4-etapas-de-estadísticas-stat-stages)
5. [STAB — Same Type Attack Bonus](#5-stab--same-type-attack-bonus)
6. [Efectividad de Tipos](#6-efectividad-de-tipos)
7. [Golpes Críticos](#7-golpes-críticos)
8. [Factor Aleatorio (Random Roll)](#8-factor-aleatorio-random-roll)
9. [Condiciones de Estado](#9-condiciones-de-estado)
10. [Clima y Terrenos](#10-clima-y-terrenos)
11. [Habilidades que Modifican Daño](#11-habilidades-que-modifican-daño)
12. [Objetos Equipados](#12-objetos-equipados)
13. [Pantallas Defensivas](#13-pantallas-defensivas)
14. [Movimientos Especiales y Casos Únicos](#14-movimientos-especiales-y-casos-únicos)
15. [Combate Doble y Triple](#15-combate-doble-y-triple)
16. [Mecánicas Generacionales Especiales](#16-mecánicas-generacionales-especiales)
17. [Orden de Turno y Prioridad](#17-orden-de-turno-y-prioridad)
18. [Resumen de Multiplicadores](#18-resumen-de-multiplicadores)

---

## 1. La Fórmula de Daño

La fórmula central de daño (a partir de Gen V en adelante) es:

```
Daño = ((((2 × Nivel / 5 + 2) × Poder × A / D) / 50) + 2)
        × Targets × PB × Weather × GlaiveRush
        × Critical × random × STAB × Type × Burn × other × ZMove × TeraShield
```

Donde cada variable se descompone así:

| Variable      | Descripción |
|---------------|-------------|
| **Nivel**     | Nivel del Pokémon atacante (1–100). |
| **Poder**     | Poder base (Base Power) del movimiento usado. |
| **A**         | Estadística de Ataque o Ataque Especial efectiva del atacante (según la categoría del movimiento). |
| **D**         | Estadística de Defensa o Defensa Especial efectiva del defensor. |
| **Targets**   | 0.75 si el movimiento golpea a múltiples objetivos en dobles/triples; 1.0 en individuales. |
| **PB**        | Parental Bond: 0.25 para el segundo golpe de Parental Bond (Gen VII+). |
| **Weather**   | Multiplicador de clima (1.5, 0.5 o 1.0). |
| **GlaiveRush**| 2.0 si el defensor usó Glaive Rush el turno anterior; 1.0 en caso contrario. |
| **Critical**  | 1.5 en golpe crítico (Gen VI+); 1.0 en caso contrario. |
| **random**    | Factor aleatorio entre 0.85 y 1.00 (ver sección 8). |
| **STAB**      | 1.5 si el tipo del movimiento coincide con un tipo del atacante; 1.0 si no. |
| **Type**      | Producto de la efectividad de tipos (0, 0.25, 0.5, 1, 2, 4). |
| **Burn**      | 0.5 si el atacante está quemado y usa un movimiento físico (excepciones: Fachada/Facade). |
| **other**     | Producto de todos los demás modificadores misceláneos (habilidades, objetos, etc.). |
| **ZMove**     | 0.25 si el defensor se protege de un Z-Move; 1.0 en caso contrario. |
| **TeraShield**| 0.5 en raids Tera si se cumplen ciertas condiciones de escudo; 1.0 en caso contrario. |

**Nota importante:** Todos los truncamientos son hacia abajo (floor) en cada paso intermedio. El daño mínimo siempre es **1** (excepto si la inmunidad de tipo aplica, que da **0**).

---

## 2. Estadísticas Base, IVs y EVs

### 2.1 Las seis estadísticas

Cada Pokémon tiene seis estadísticas (stats):

- **HP** (Hit Points / Puntos de Salud)
- **Attack** (Ataque Físico)
- **Defense** (Defensa Física)
- **Special Attack** (Ataque Especial)
- **Special Defense** (Defensa Especial)
- **Speed** (Velocidad)

### 2.2 Estadísticas Base (Base Stats)

Cada especie tiene valores fijos de base stats que definen su perfil. Por ejemplo, Garchomp tiene 108 de Ataque base, mientras que Blissey tiene 10.

### 2.3 IVs (Individual Values)

Valores individuales heredados o aleatorios al encontrar/criar al Pokémon. Rango: **0–31** por estadística. Representan el "talento genético" del Pokémon.

### 2.4 EVs (Effort Values)

Valores de esfuerzo ganados al derrotar Pokémon o con objetos. Rango: **0–252** por estadística, con un máximo total de **510** entre todas las estadísticas.

### 2.5 Fórmula de cálculo de estadísticas

**Para HP:**
```
HP = ((2 × Base + IV + (EV / 4)) × Nivel / 100) + Nivel + 10
```

**Para cualquier otra estadística:**
```
Stat = (((2 × Base + IV + (EV / 4)) × Nivel / 100) + 5) × NaturalezaMod
```

- **NaturalezaMod** = 1.1 si la naturaleza favorece la stat, 0.9 si la perjudica, 1.0 si es neutra.

**Caso especial — Shedinja:** Su HP siempre es **1**, sin importar IVs, EVs o nivel.

---

## 3. Naturalezas (Natures)

Existen **25 naturalezas**. De ellas, 5 son neutras (no modifican nada) y 20 suben una estadística un 10% (+10%) mientras bajan otra un 10% (-10%). La naturaleza **nunca** afecta a HP.

| Naturaleza | Stat +10% | Stat -10% |
|------------|-----------|-----------|
| Adamant    | Attack    | Sp. Atk   |
| Modest     | Sp. Atk   | Attack    |
| Jolly      | Speed     | Sp. Atk   |
| Timid      | Speed     | Attack    |
| Bold       | Defense   | Attack    |
| Impish     | Defense   | Sp. Atk   |
| Calm       | Sp. Def   | Attack    |
| Careful    | Sp. Def   | Sp. Atk   |
| Brave      | Attack    | Speed     |
| Quiet      | Sp. Atk   | Speed     |
| Relaxed    | Defense   | Speed     |
| Sassy      | Sp. Def   | Speed     |
| Naughty    | Attack    | Sp. Def   |
| Lonely     | Attack    | Defense   |
| Mild       | Sp. Atk   | Defense   |
| Rash       | Sp. Atk   | Sp. Def   |
| Hasty      | Speed     | Defense   |
| Naive      | Speed     | Sp. Def   |
| Lax        | Defense   | Sp. Def   |
| Gentle     | Sp. Def   | Defense   |
| Hardy, Docile, Bashful, Quirky, Serious | — | — |

A partir de Gen VIII se pueden usar **Mint** (Menta) para cambiar la influencia de la naturaleza en stats sin cambiar la naturaleza de nombre.

---

## 4. Etapas de Estadísticas (Stat Stages)

Cada estadística de combate (excepto HP) puede modificarse en combate mediante movimientos como Swords Dance, Calm Mind, Intimidate, etc. Las etapas van de **-6 a +6**.

### 4.1 Multiplicadores de etapas (Atk, Def, SpAtk, SpDef, Spd)

| Etapa | Multiplicador |
|-------|---------------|
| -6    | 2/8 = 0.25    |
| -5    | 2/7 ≈ 0.286   |
| -4    | 2/6 ≈ 0.333   |
| -3    | 2/5 = 0.40    |
| -2    | 2/4 = 0.50    |
| -1    | 2/3 ≈ 0.667   |
|  0    | 2/2 = 1.00    |
| +1    | 3/2 = 1.50    |
| +2    | 4/2 = 2.00    |
| +3    | 5/2 = 2.50    |
| +4    | 6/2 = 3.00    |
| +5    | 7/2 = 3.50    |
| +6    | 8/2 = 4.00    |

Fórmula general: `max(2, 2 + etapa) / max(2, 2 - etapa)`

### 4.2 Multiplicadores de Precisión y Evasión

Para Accuracy y Evasion, los multiplicadores son distintos:

| Etapa | Multiplicador |
|-------|---------------|
| -6    | 3/9 = 0.333   |
| -5    | 3/8 = 0.375   |
| -4    | 3/7 ≈ 0.429   |
| -3    | 3/6 = 0.500   |
| -2    | 3/5 = 0.600   |
| -1    | 3/4 = 0.750   |
|  0    | 3/3 = 1.000   |
| +1    | 4/3 ≈ 1.333   |
| +2    | 5/3 ≈ 1.667   |
| +3    | 6/3 = 2.000   |
| +4    | 7/3 ≈ 2.333   |
| +5    | 8/3 ≈ 2.667   |
| +6    | 9/3 = 3.000   |

### 4.3 Cálculo de probabilidad de acierto

```
P(acierto) = Precisión_base_movimiento × (Etapas_Accuracy_atacante / Etapas_Evasion_defensor) × otros_modificadores
```

Movimientos con precisión "—" (como Aerial Ace, Swift) **nunca fallan** bajo circunstancias normales.

### 4.4 Interacción con golpes críticos

Un golpe crítico **ignora** las bajadas de estadísticas del atacante y las subidas de estadísticas del defensor. Es decir, si tienes -2 Ataque y golpeas un crítico, se usa tu Ataque sin modificar. Si el rival tiene +4 Defensa y recibes un crítico, se ignora esa subida.

---

## 5. STAB — Same Type Attack Bonus

### 5.1 Mecánica básica

Si el tipo del movimiento coincide con **al menos uno de los tipos del Pokémon** atacante, el daño se multiplica por **1.5**.

Ejemplo: Charizard (Fuego/Volador) usa Lanzallamas (Fuego) → STAB aplica (×1.5).

### 5.2 STAB y habilidades

| Habilidad       | Efecto sobre STAB |
|------------------|--------------------|
| **Adaptability** | STAB sube a **2.0** en lugar de 1.5. |
| **Protean / Libero** | El Pokémon cambia al tipo del movimiento que va a usar **antes** de ejecutarlo, obteniendo siempre STAB (Gen IX lo limita a una vez por entrada). |

### 5.3 STAB y Teracristalización (Gen IX)

- Si un Pokémon Teracristaliza a un tipo que ya tenía, el STAB para ese tipo sube a **2.0** (o **2.25** con Adaptability).
- Si Teracristaliza a un tipo nuevo, obtiene STAB de **1.5** para su Tera tipo y mantiene STAB de **1.5** para sus tipos originales.
- Si Teracristaliza a un tipo que ya tenía, los movimientos de sus otros tipos originales mantienen STAB de **1.5**.
- El "Tera STAB" para tipos originales + Tera coincidente es de **2.0** (redondeado; la mecánica exacta usa un cálculo adaptado).

---

## 6. Efectividad de Tipos

### 6.1 Tabla de efectividad

El sistema de tipos es la columna vertebral del combate. Cada tipo de movimiento tiene relaciones de efectividad contra cada tipo de defensor:

| Multiplicador   | Nombre           | Mensaje en juego            |
|------------------|------------------|-----------------------------|
| **×0**          | Inmune           | "No afecta a..."           |
| **×0.25**       | Doblemente resistido | (No hay mensaje especial) |
| **×0.5**        | No muy eficaz    | "No es muy eficaz..."      |
| **×1**          | Daño neutral     | (Sin mensaje)              |
| **×2**          | Supereficaz      | "¡Es supereficaz!"         |
| **×4**          | Doblemente supereficaz | "¡Es supereficaz!" (mismo mensaje) |

### 6.2 Cálculo con doble tipo

El multiplicador final de tipo es el **producto** de la efectividad contra cada tipo del defensor.

Ejemplo: Rayo (Eléctrico) vs Gyarados (Agua/Volador)
- Eléctrico vs Agua = ×2
- Eléctrico vs Volador = ×2
- Total: ×2 × ×2 = **×4**

Ejemplo: Terremoto (Tierra) vs Skarmory (Acero/Volador)
- Tierra vs Acero = ×2
- Tierra vs Volador = ×0
- Total: ×2 × ×0 = **×0** (inmune)

### 6.3 Tabla completa de tipos (18 tipos, Gen VI+)

```
Atacante →       NOR FUE AER VEN TIE ROC BIC FAN ACE FUG AGU PLA ELE PSI HIE DRA SIN HAD
Defensor ↓
Normal            1   2   1   1   1   1   1   0   1   1   1   1   1   1   1   1   1   1
Fuego             1   .5  1   1   2   .5  .5  1   .5  1   2   .5  1   1   .5  1   1   1
Agua              1   1   1   1   1   1   1   1   .5  1   .5  2   2   1   .5  1   1   1
Planta            1   .5  2   2   .5  1   2   1   1   .5  .5  .5  .5  1   2   1   1   1
Eléctrico         1   1   .5  1   2   1   1   1   .5  1   1   1   .5  1   1   1   1   1
Hielo             1   2   1   1   1   2   1   1   2   1   1   1   1   1   .5  1   1   1
Lucha             1   1   2   .5  1   .5  .5  0   1   2   1   1   1   2   1   1   .5  2
Veneno            1   1   1   .5  2   1   .5  .5  1   1   1   .5  1   2   1   1   1   .5
Tierra            1   1   0   .5  1   .5  1   1   1   1   2   2   0   1   2   1   1   1
Volador           1   1   1   1   0   2   .5  1   1   1   1   .5  2   1   2   1   1   1
Psíquico          1   .5  1   1   1   1   2   2   1   1   1   1   1   .5  1   1   2   1
Bicho             1   .5  2   1   .5  2   1   1   1   1   1   .5  1   1   1   1   1   1
Roca              .5  2   .5  .5  2   1   2   1   2   1   2   2   1   1   1   1   1   1
Fantasma          0   1   1   .5  1   1   .5  2   1   1   1   1   1   1   1   1   2   1
Dragón            1   .5  1   1   1   1   1   1   1   1   .5  .5  .5  1   2   2   1   2
Siniestro         1   2   1   1   1   1   2   .5  1   1   1   1   1   0   1   1   .5  2
Acero             .5  2   .5  0   2   .5  .5  .5  .5  .5  1   .5  1   .5  .5  .5  .5  .5
Hada              1   .5  1   2   1   1   .5  1   2   1   1   1   1   1   1   0   .5  1
```

*(Las abreviaturas corresponden a: NOR=Normal, FUE=Fuego, AER=desplazada — se recomienda consultar la tabla oficial completa para precisión pixel-perfect.)*

### 6.4 Inmunidades y cómo romperlas

| Inmunidad                     | Rota por                                      |
|-------------------------------|-----------------------------------------------|
| Normal/Lucha → Fantasma       | Habilidad **Scrappy** o movimiento **Foresight/Odor Sleuth**. |
| Tierra → Volador              | Movimiento **Gravity**, objeto **Iron Ball**, habilidad **Mold Breaker** (no aplica; Levitate es la habilidad, no el tipo). Smack Down/Thousand Arrows eliminan la inmunidad posicional. |
| Psíquico → Siniestro          | No se puede romper por tipo. Miracle Eye permite al Psíquico golpear a Siniestro. |
| Fantasma → Normal             | Habilidad **Scrappy**, Foresight/Odor Sleuth. |
| Dragón → Hada                 | No se puede romper directamente.              |
| Eléctrico → Tierra            | Gravity, Ring Target (objeto).                |
| Veneno → Acero                | Habilidad **Corrosion** (permite envenenar, no daño de tipo). |

### 6.5 Habilidades que anulan inmunidades de tipo

- **Mold Breaker / Teravolt / Turboblaze:** Ignoran habilidades del defensor (como Levitate, Volt Absorb, etc.), pero no inmunidades de tipo puras.
- **Scrappy:** Permite que movimientos Normal y Lucha golpeen a Fantasma.

---

## 7. Golpes Críticos

### 7.1 Probabilidad de golpe crítico

Cada movimiento tiene una etapa de golpe crítico:

| Etapa Crítica | Probabilidad (Gen VI+) |
|---------------|------------------------|
| +0            | 1/24 ≈ 4.17%          |
| +1            | 1/8 = 12.5%           |
| +2            | 1/2 = 50%             |
| +3 o más      | 1/1 = 100%            |

### 7.2 Fuentes de aumento de etapa crítica

| Fuente                    | Etapas |
|---------------------------|--------|
| Movimiento con alto ratio de crítico (Slash, Leaf Blade, etc.) | +1 |
| Habilidad **Super Luck**  | +1     |
| Objeto **Scope Lens / Razor Claw** | +1 |
| Objeto **Leek** (Farfetch'd/Sirfetch'd) | +2 |
| Objeto **Lucky Punch** (Chansey) | +2 |
| Movimiento **Focus Energy** | +2   |
| Efecto **Lansat Berry**   | +2     |

Estos se **suman**, así que Focus Energy (+2) + Super Luck (+1) + movimiento alto ratio (+1) = etapa +4 → 100% de críticos.

### 7.3 Efectos del golpe crítico

- **Multiplicador de daño:** ×1.5 (Gen VI+). En Gen II–V era ×2.
- **Ignora bajadas de estadísticas** del atacante (Ataque o Ataque Especial).
- **Ignora subidas de estadísticas** del defensor (Defensa o Defensa Especial).
- **Ignora Reflect, Light Screen y Aurora Veil** del defensor.
- **No ignora:** quemadura en el atacante (sí lo hacía en Gen II–V, ya no).

### 7.4 Habilidades relacionadas

- **Sniper:** Los críticos hacen ×2.25 en lugar de ×1.5.
- **Battle Armor / Shell Armor:** El oponente no puede recibirte golpes críticos.
- **Anger Point:** Si recibes un crítico, tu Ataque sube a +6.
- **Merciless:** Siempre golpe crítico contra un objetivo envenenado.

---

## 8. Factor Aleatorio (Random Roll)

Cada ataque que hace daño tiene un **multiplicador aleatorio** entre **0.85 y 1.00**, inclusive, distribuido uniformemente. Esto se implementa como un entero aleatorio entre 85 y 100, dividido entre 100, y truncado (floor).

Esto significa que incluso el mismo movimiento, contra el mismo objetivo, con las mismas condiciones, puede hacer daño variable de un turno a otro (+/- 15% respecto al máximo). Cada valor (85, 86, 87... 100) tiene la misma probabilidad: 1/16.

---

## 9. Condiciones de Estado

### 9.1 Condiciones de estado principales (solo una a la vez)

| Estado         | Efecto en combate |
|----------------|-------------------|
| **Quemadura (Burn)** | Daño de movimientos **físicos** reducido al 50% (×0.5). Además, pierde 1/16 de HP máximo por turno (1/8 en Gen I). |
| **Parálisis (Paralysis)** | Velocidad reducida al **50%** (Gen VII+; 25% en gens anteriores). 25% de probabilidad de no poder moverse cada turno. |
| **Envenenamiento (Poison)** | Pierde 1/8 de HP máximo por turno. |
| **Envenenamiento grave (Toxic/Bad Poison)** | Pierde 1/16 × N de HP por turno, donde N incrementa cada turno (1/16, 2/16, 3/16...). |
| **Sueño (Sleep)** | No puede actuar (excepto Ronquido/Sleep Talk). Dura 1–3 turnos (Gen V+). |
| **Congelación (Freeze)** | No puede actuar. 20% de probabilidad de descongelarse cada turno. Ciertos movimientos de Fuego lo curan. |
| **Debilitado (Faint)** | HP = 0. Fuera de combate. |

### 9.2 Condiciones volátiles (pueden apilarse)

| Estado           | Efecto |
|------------------|--------|
| **Confusión**    | 33% de probabilidad de golpearse a sí mismo (Gen VII+). El auto-daño es de tipo físico, poder 40, sin STAB ni efectividad. |
| **Enamoramiento (Attract/Infatuation)** | 50% de probabilidad de no poder atacar. |
| **Flinch (Retroceso)** | Pierde el turno. Solo funciona si el atacante se mueve primero. |
| **Trampa (Bind, Wrap, Whirlpool, etc.)** | No puede huir/ser intercambiado. Pierde 1/8 HP por turno (1/6 con Binding Band). |
| **Leech Seed** | Pierde 1/8 HP por turno, transferido al rival. |
| **Maldición (Curse, de tipo Fantasma)** | Pierde 1/4 HP por turno. |
| **Embargo**     | No puede usar objetos equipados. |
| **Heal Block**  | No puede curarse. |

### 9.3 Impacto de quemadura en el daño (detalle)

La quemadura aplica un multiplicador de **×0.5** al daño de movimientos **físicos**. Excepciones:

- **Facade:** Si el usuario está quemado, el multiplicador de quemadura no se aplica y Facade duplica su poder base (de 70 a 140).
- **Guts:** La habilidad ignora la reducción de daño por quemadura y además aumenta el Ataque ×1.5.

---

## 10. Clima y Terrenos

### 10.1 Climas (Weather)

| Clima           | Efecto en daño | Otros efectos |
|-----------------|----------------|---------------|
| **Sol (Harsh Sunlight)** | Fuego ×1.5, Agua ×0.5. | Solar Beam no necesita carga. Trueno baja a 50% precisión. Synthesis/Moonlight/Morning Sun recuperan 2/3 HP. |
| **Lluvia (Rain)** | Agua ×1.5, Fuego ×0.5. | Trueno y Huracán tienen 100% precisión. Solar Beam baja a 60 de poder. |
| **Tormenta de arena (Sandstorm)** | Roca: Def. Esp. ×1.5 para tipos Roca. | 1/16 HP de daño por turno a todos excepto Roca, Tierra y Acero. |
| **Granizo / Nieve (Hail/Snow)** | Hielo: Defensa ×1.5 para tipos Hielo (Gen IX: Snow). | En Hail: 1/16 HP de daño por turno a todos excepto Hielo. En Snow (Gen IX): sin daño residual, solo +50% Defensa a tipos Hielo. Aurora Veil solo se puede usar bajo Hail/Snow. |
| **Sol extremo (Extremely Harsh Sunlight)** | Fuego ×1.5. Movimientos tipo Agua son completamente anulados. | Solo vía Desolate Land (Primal Groudon). |
| **Lluvia extrema (Heavy Rain)** | Agua ×1.5. Movimientos tipo Fuego son completamente anulados. | Solo vía Primordial Sea (Primal Kyogre). |
| **Corriente de aire (Strong Winds)** | Reduce debilidades de tipo Volador a ×1. | Solo vía Delta Stream (Mega Rayquaza). |

### 10.2 Terrenos (Terrains) — Gen VII+

Los terrenos afectan solo a Pokémon **en contacto con el suelo** (no voladores, ni Levitate, ni bajo Magnet Rise/Telekinesis).

| Terreno            | Efecto en daño | Otros efectos |
|--------------------|----------------|---------------|
| **Eléctrico (Electric Terrain)** | Eléctrico ×1.3 (Gen VIII+; ×1.5 en Gen VII). | Previene Sueño. |
| **Psíquico (Psychic Terrain)** | Psíquico ×1.3 (Gen VIII+; ×1.5 en Gen VII). | Bloquea movimientos de prioridad aumentada contra Pokémon en el suelo. |
| **Hierba (Grassy Terrain)** | Planta ×1.3 (Gen VIII+; ×1.5 en Gen VII). | Recupera 1/16 HP por turno. Terremoto y Bulldoze hacen ×0.5 de daño. |
| **Niebla (Misty Terrain)** | Dragón ×0.5 contra Pokémon en el suelo. | Previene condiciones de estado. |

### 10.3 Habilidades de clima/terreno automático

- **Drought / Drizzle / Sand Stream / Snow Warning:** Activan sol/lluvia/tormenta/nieve al entrar.
- **Electric/Psychic/Grassy/Misty Surge:** Activan el terreno correspondiente al entrar.
- La duración es de **5 turnos** (u 8 con el objeto extensor: Heat Rock, Damp Rock, Smooth Rock, Icy Rock, Terrain Extender).

---

## 11. Habilidades que Modifican Daño

### 11.1 Habilidades que aumentan daño ofensivo

| Habilidad           | Efecto |
|----------------------|--------|
| **Huge Power / Pure Power** | Ataque ×2. |
| **Hustle**           | Ataque ×1.5, pero precisión de movimientos físicos ×0.8. |
| **Guts**             | Ataque ×1.5 cuando tiene un estado (quemadura, parálisis, envenenamiento, etc.). Ignora reducción de quemadura. |
| **Overgrow / Blaze / Torrent / Swarm** | ×1.5 al daño de movimientos del tipo correspondiente cuando HP ≤ 1/3. |
| **Adaptability**     | STAB sube a ×2.0. |
| **Choice Band/Specs efecto en habilidades** | (No es habilidad, sino objeto; ver sección 12). |
| **Technician**       | Movimientos con poder base ≤60 obtienen ×1.5 al poder. |
| **Tough Claws**      | Movimientos de contacto hacen ×1.3 de daño. |
| **Iron Fist**        | Movimientos de puñetazo hacen ×1.2 de daño. |
| **Strong Jaw**       | Movimientos de mordisco hacen ×1.5 de daño. |
| **Mega Launcher**    | Movimientos de aura/pulso hacen ×1.5 de daño. |
| **Sheer Force**      | Movimientos con efectos secundarios hacen ×1.3, pero pierden el efecto secundario. |
| **Reckless**         | Movimientos con retroceso (recoil) hacen ×1.2 de daño. |
| **Analytic**         | ×1.3 si atacas de último en el turno. |
| **Solar Power**      | Ataque Especial ×1.5 bajo sol, pero pierde 1/8 HP por turno. |
| **Sand Force**       | Tierra, Roca y Acero hacen ×1.3 bajo tormenta de arena. |
| **Aerilate / Pixilate / Refrigerate / Galvanize** | Convierte movimientos Normal al tipo correspondiente y añade ×1.2 al poder (Gen VII+). |
| **Neuroforce**       | Movimientos supereficaces hacen ×1.25 de daño adicional. |
| **Tinted Lens**      | Movimientos "no muy eficaces" hacen el doble de daño (la resistencia se reduce a la mitad). |
| **Sniper**           | Golpes críticos hacen ×2.25 en vez de ×1.5. |
| **Stakeout**         | ×2 daño al oponente que acaba de entrar al campo (switch-in). |
| **Supreme Overlord** | +10% daño por cada aliado debilitado (máximo +50%). |

### 11.2 Habilidades que reducen daño recibido

| Habilidad              | Efecto |
|------------------------|--------|
| **Multiscale / Shadow Shield** | Daño recibido ×0.5 cuando HP está al 100%. |
| **Fur Coat**           | Daño físico recibido ×0.5. |
| **Ice Scales**         | Daño especial recibido ×0.5. |
| **Filter / Solid Rock / Prism Armor** | Daño de movimientos supereficaces recibido ×0.75. |
| **Thick Fat**          | Daño de Fuego y Hielo recibido ×0.5. |
| **Heatproof**          | Daño de Fuego recibido ×0.5. |
| **Fluffy**             | Daño de contacto recibido ×0.5, pero Fuego hace ×2. |
| **Levitate**           | Inmune a movimientos tipo Tierra. |
| **Volt Absorb / Water Absorb / Dry Skin** | Absorbe movimientos del tipo correspondiente y recupera HP. |
| **Flash Fire**         | Inmune a Fuego; al recibir un movimiento Fuego, sus propios movimientos Fuego hacen ×1.5. |
| **Storm Drain / Lightning Rod** | Atrae movimientos Agua/Eléctrico (en dobles); inmune y sube Ataque Especial +1. |
| **Motor Drive**        | Inmune a Eléctrico; sube Velocidad +1. |
| **Sap Sipper**         | Inmune a Planta; sube Ataque +1. |
| **Friend Guard**       | Aliados reciben ×0.75 daño (en dobles). |
| **Marvel Scale**       | Defensa ×1.5 cuando tiene un estado alterado. |

### 11.3 Habilidades que ignoran defensas del rival

| Habilidad         | Efecto |
|--------------------|--------|
| **Mold Breaker / Teravolt / Turboblaze** | Ignora habilidades defensivas del rival (Sturdy, Multiscale, Levitate, etc.). |
| **Unaware** | Ignora cambios de estadísticas del oponente (tanto ofensivos como defensivos según el contexto: ignora subidas de Ataque del rival al recibir daño, e ignora subidas de Defensa del rival al atacar). |

---

## 12. Objetos Equipados

### 12.1 Objetos que modifican daño ofensivo

| Objeto                     | Efecto |
|----------------------------|--------|
| **Choice Band**            | Ataque ×1.5, pero solo puede usar un movimiento hasta cambiar. |
| **Choice Specs**           | Ataque Especial ×1.5, misma restricción. |
| **Life Orb**               | Daño ×1.3, pero pierde 1/10 HP por ataque. |
| **Expert Belt**            | Movimientos supereficaces hacen ×1.2. |
| **Metronome (objeto)**     | Cada uso consecutivo del mismo movimiento sube daño en ×0.2 (máx. ×2.0). |
| **Muscle Band**            | Daño de movimientos físicos ×1.1. |
| **Wise Glasses**           | Daño de movimientos especiales ×1.1. |
| **Type-boosting items** (Charcoal, Mystic Water, etc.) | Daño del tipo correspondiente ×1.2. |
| **Plates** (Flame Plate, etc.) | Igual que arriba, ×1.2 al tipo correspondiente. Interactúan con Arceus y Judgment. |
| **Gems** (Fire Gem, etc.) | Se consume; primer movimiento del tipo correspondiente hace ×1.3 (Gen V: ×1.5). |
| **Weather-boosting rocks** (Heat Rock, etc.) | No modifican daño directamente, extienden duración del clima a 8 turnos. |
| **Punching Glove**         | Movimientos de puñetazo hacen ×1.1 y dejan de ser de contacto. |

### 12.2 Objetos que modifican daño defensivo

| Objeto                   | Efecto |
|--------------------------|--------|
| **Eviolite**             | Defensa y Def. Especial ×1.5 si el Pokémon puede evolucionar. |
| **Assault Vest**         | Defensa Especial ×1.5, pero no puede usar movimientos de estado. |
| **Rocky Helmet**         | El atacante que hace contacto pierde 1/6 HP. |
| **Air Balloon**          | Inmune a movimientos tipo Tierra hasta recibir un golpe. |
| **Focus Sash**           | Sobrevive con 1 HP a un golpe que lo noquearía si estaba al 100% HP. Se consume. |
| **Focus Band**           | 10% de probabilidad de sobrevivir con 1 HP a un golpe mortal. |
| **Berries defensivas** (Babiri, Charti, etc.) | ×0.5 daño de un tipo supereficaz específico. Se consume. |

### 12.3 Objetos especiales de velocidad

| Objeto              | Efecto |
|---------------------|--------|
| **Choice Scarf**    | Velocidad ×1.5, pero bloqueado a un movimiento. |
| **Iron Ball**       | Velocidad ×0.5 y el Pokémon pierde inmunidad a Tierra (si es Volador o tiene Levitate). |
| **Lagging Tail / Full Incense** | Siempre se mueve de último en su bracket de prioridad. |
| **Quick Claw**      | 20% de probabilidad de moverse primero en su bracket de prioridad. |

---

## 13. Pantallas Defensivas

### 13.1 Tipos de pantalla

| Movimiento         | Efecto | Duración |
|--------------------|--------|----------|
| **Reflect**        | Daño físico recibido ×0.5 (×0.66 en dobles). | 5 turnos (8 con Light Clay). |
| **Light Screen**   | Daño especial recibido ×0.5 (×0.66 en dobles). | 5 turnos (8 con Light Clay). |
| **Aurora Veil**    | Ambos tipos de daño ×0.5 (×0.66 en dobles). Solo se puede usar bajo Hail/Snow. | 5 turnos (8 con Light Clay). |

### 13.2 Interacciones con pantallas

- Los **golpes críticos ignoran** las pantallas (Gen VI+).
- **Brick Break / Psychic Fangs / Raging Bull** destruyen las pantallas del rival.
- **Defog** elimina las pantallas de ambos lados.
- Las pantallas son un efecto de **campo/equipo**, no del Pokémon individual. Persisten aunque el Pokémon que las puso sea retirado.
- **Infiltrator** ignora las pantallas del rival.

---

## 14. Movimientos Especiales y Casos Únicos

### 14.1 Movimientos con fórmula de daño única

| Movimiento            | Cálculo de daño |
|-----------------------|-----------------|
| **Seismic Toss / Night Shade** | Daño = Nivel del usuario. Ignora tipo del defensor (salvo inmunidades). |
| **Dragon Rage**       | Siempre hace 40 de daño fijo. |
| **Sonic Boom**        | Siempre hace 20 de daño fijo. |
| **Super Fang**        | Reduce HP actual del rival al 50%. |
| **Endeavor**          | Reduce HP del rival al HP actual del usuario. |
| **Final Gambit**      | Daño = HP actual del usuario. El usuario se debilita. |
| **Counter**           | Devuelve ×2 el último daño físico recibido. |
| **Mirror Coat**       | Devuelve ×2 el último daño especial recibido. |
| **Metal Burst**       | Devuelve ×1.5 el último daño recibido (cualquier tipo). |
| **OHKO moves** (Sheer Cold, Fissure, Guillotine, Horn Drill) | Derrota al instante. Precisión = 30% + (Nivel_usuario - Nivel_rival). Falla si el defensor tiene mayor nivel. |
| **Psywave**           | Daño = Nivel × (aleatorio 50–150)%. |
| **Pain Split**        | Ambos Pokémon terminan con el promedio de sus HP actuales. |
| **Bide**              | Almacena daño recibido en 2 turnos, devuelve ×2. |

### 14.2 Movimientos con poder variable

| Movimiento       | Mecánica del poder |
|------------------|-------------------|
| **Eruption / Water Spout** | Poder = 150 × (HP_actual / HP_máx). |
| **Flail / Reversal** | Más poder cuanto menor sea el HP (máximo 200 con HP mínimo). |
| **Gyro Ball**     | Poder = 25 × (Velocidad_rival / Velocidad_propia). Máximo 150. |
| **Electro Ball**  | Poder basado en ratio Velocidad_propia / Velocidad_rival. Máximo 150. |
| **Low Kick / Grass Knot** | Poder basado en peso del rival (20–120). |
| **Heavy Slam / Heat Crash** | Poder basado en ratio de peso propio/rival (40–120). |
| **Acrobatics**    | Poder ×2 (110) si el usuario no tiene objeto. |
| **Facade**        | Poder ×2 (140) si el usuario tiene quemadura, parálisis o envenenamiento. |
| **Hex**           | Poder ×2 (130) si el rival tiene un estado alterado. |
| **Brine**         | Poder ×2 (130) si el rival tiene ≤50% HP. |
| **Venoshock**     | Poder ×2 (130) si el rival está envenenado. |
| **Knock Off**     | Poder ×1.5 (97) si el rival tiene un objeto que se puede quitar. |
| **Stored Power / Power Trip** | Poder = 20 + 20 × (suma de todas las subidas de estadísticas). |
| **Foul Play**     | Usa el Ataque **del defensor** en lugar del atacante para calcular daño. |
| **Body Press**    | Usa la Defensa **del atacante** en lugar de su Ataque para calcular daño. |
| **Photon Geyser / Light That Burns the Sky** | Usa la mayor entre Ataque y Ataque Especial del usuario. |

### 14.3 Movimientos multihit

| Movimiento          | Golpes | Distribución de probabilidad |
|---------------------|--------|------------------------------|
| **Double Hit, Double Kick, etc.** | Siempre 2. | — |
| **Triple Kick/Axel** | 3 (cada golpe sube poder: 10→20→30 o 30→60→90). | — |
| **Population Bomb** | 1–10 (afectado por precisión). | — |
| **Bullet Seed, Icicle Spear, Rock Blast, etc.** | 2–5. | 2 golpes: 35%, 3: 35%, 4: 15%, 5: 15%. |
| **Skill Link (habilidad)** | Siempre 5 golpes para movimientos de 2–5 golpes. | — |
| **Parental Bond (habilidad)** | El usuario golpea dos veces. El segundo golpe hace ×0.25 del daño (Gen VII+). | — |

Cada golpe de un movimiento multihit aplica **independientemente** la probabilidad de crítico, pero los multiplicadores de tipo y STAB son iguales para todos los golpes.

### 14.4 Movimientos de prioridad en protección

| Movimiento     | Efecto |
|----------------|--------|
| **Protect / Detect** | Bloquea cualquier movimiento dirigido al usuario. Probabilidad de fallo aumenta con uso consecutivo (1/1, 1/3, 1/9...). |
| **King's Shield** | Protege y baja Ataque -1 al rival si usó movimiento de contacto. |
| **Baneful Bunker** | Protege y envenena al rival si usó movimiento de contacto. |
| **Spiky Shield** | Protege y causa 1/8 HP de daño si el rival usó contacto. |
| **Silk Trap**   | Protege y baja Velocidad -1 si el rival usó contacto. |
| **Wide Guard**  | Protege al equipo de movimientos que afectan a múltiples objetivos. |
| **Quick Guard** | Protege al equipo de movimientos con prioridad aumentada. |
| **Max Guard**   | Protege de todo, incluyendo Max Moves. Solo durante Dynamax/Gigantamax. |
| **Endure**      | El usuario sobrevive con mínimo 1 HP. |

---

## 15. Combate Doble y Triple

### 15.1 Spread Moves (movimientos de área)

En combate doble, los movimientos que golpean a **múltiples objetivos** sufren una reducción:

- El daño de cada golpe se multiplica por **×0.75**.
- Si solo queda un oponente, no se aplica la reducción.

Ejemplos de spread moves: Earthquake, Surf, Heat Wave, Rock Slide, Dazzling Gleam.

### 15.2 Selección de objetivos

- En dobles, los movimientos de un solo objetivo pueden dirigirse a cualquiera de los Pokémon en campo (incluido el aliado para ciertos movimientos).
- Si el objetivo seleccionado se debilita antes de que el movimiento se ejecute, el movimiento **puede redirigirse** al otro oponente (Gen V+), o fallar si no hay objetivo válido.

### 15.3 Habilidades relevantes en dobles

| Habilidad       | Efecto |
|------------------|--------|
| **Friend Guard** | Aliados reciben ×0.75 daño. |
| **Telepathy**    | Inmune al daño de aliados. |
| **Lightning Rod / Storm Drain** | Atrae todos los movimientos Eléctrico/Agua del campo hacia sí mismo. Sube Sp. Atk +1 e inmune. |
| **Plus / Minus**  | Sp. Atk ×1.5 si un aliado tiene la habilidad complementaria. |
| **Battery**       | Movimientos especiales de aliados hacen ×1.3 de daño. |
| **Power Spot**    | Movimientos de aliados hacen ×1.3 de daño. |
| **Flower Gift** (Cherrim) | Ataque y Def. Especial de aliados ×1.5 bajo sol. |

---

## 16. Mecánicas Generacionales Especiales

### 16.1 Mega Evolución (Gen VI–VII)

- Cambia las stats base del Pokémon, lo cual recalcula todo.
- Puede cambiar tipo(s) y habilidad.
- Solo un Mega por equipo por combate.
- La Mega Evolución ocurre **antes del movimiento**, por lo que la nueva Velocidad se aplica en el mismo turno (desde Gen VII; en Gen VI la velocidad del turno de Mega se calculaba con la velocidad pre-Mega).

### 16.2 Z-Moves (Gen VII)

- Convierten un movimiento del tipo correspondiente en un Z-Move con poder fijo elevado (normalmente 160–200 para físicos/especiales).
- No usan la precisión del movimiento base — **nunca fallan** (excepto si el oponente está en semi-invulnerabilidad como Fly/Dig).
- Si el rival usa Protect, el Z-Move hace **×0.25** del daño en lugar de 0.
- Los movimientos de estado convertidos en Z-Moves otorgan un efecto adicional (como subir una stat) y luego ejecutan el movimiento de estado.

### 16.3 Dynamax / Gigantamax (Gen VIII)

- HP se **duplica** (el HP actual y máximo se multiplican por el Dynamax Level / 10, entre ×1.5 y ×2.0).
- Todos los movimientos se convierten en **Max Moves** con poder fijo según el poder base del movimiento original.
- Los Max Moves nunca fallan.
- Max Moves de estado se convierten en Max Guard.
- Los Max Moves ignoran efectos como Protect (hacen ×0.25 del daño contra Protect, pero Max Guard bloquea completamente).
- Los Max Moves tienen efectos secundarios automáticos (clima, terreno, subida/bajada de stats).
- G-Max Moves (Gigantamax) tienen efectos secundarios únicos.

### 16.4 Teracristalización (Gen IX)

- El Pokémon obtiene un nuevo tipo (Tera Type) que **reemplaza** sus tipos originales para cálculos defensivos.
- Para cálculos ofensivos:
  - Si el Tera tipo coincide con uno original: STAB del movimiento de ese tipo = **×2.0**.
  - Si el Tera tipo es nuevo: obtiene STAB ×1.5 para el nuevo tipo. Mantiene STAB ×1.5 para tipos originales.
- Movimientos **Tera Blast** cambian al Tera tipo del usuario y usan la mayor entre Ataque y Ataque Especial.
- **Stellar Tera tipo:** Cada tipo de movimiento obtiene un boost de STAB ×2.0 una vez (primera vez que se usa cada tipo). No cambia debilidades/resistencias defensivas; en su lugar, mantiene los tipos originales para cálculos defensivos pero con neutralidad forzada a todos los tipos (todo es ×1.0 defensivamente, excepto las debilidades originales que se mantienen según la implementación exacta).

### 16.5 Diferencias entre categoría Física/Especial por generación

- **Gen I–III:** La categoría (física o especial) depende del **tipo** del movimiento, no del movimiento individual. Todos los movimientos Fuego son especiales, todos los movimientos Lucha son físicos, etc.
- **Gen IV+:** Cada movimiento tiene su propia categoría (física, especial o estado) independiente del tipo. Esto fue el "Physical/Special split".

---

## 17. Orden de Turno y Prioridad

### 17.1 Brackets de prioridad

Cada movimiento tiene un valor de prioridad. Se resuelven en orden de mayor a menor prioridad, y dentro del mismo bracket, por Velocidad.

| Prioridad | Movimientos |
|-----------|-------------|
| **+5**    | Helping Hand |
| **+4**    | Protect, Detect, Endure, King's Shield, Baneful Bunker, Spiky Shield, Max Guard, Silk Trap |
| **+3**    | Fake Out, Quick Guard, Wide Guard, Crafty Shield |
| **+2**    | Extreme Speed, First Impression, Feint, Accelerock |
| **+1**    | Aqua Jet, Bullet Punch, Ice Shard, Mach Punch, Quick Attack, Shadow Sneak, Sucker Punch, Water Shuriken, Grassy Glide (en Grassy Terrain), Jet Punch |
| **0**     | La mayoría de movimientos. |
| **-1**    | Vital Throw |
| **-2**    | — |
| **-3**    | Focus Punch, Beak Blast, Shell Trap |
| **-4**    | Avalanche, Revenge |
| **-5**    | Counter, Mirror Coat, Metal Burst |
| **-6**    | Whirlwind, Roar, Circle Throw, Dragon Tail, Teleport |
| **-7**    | Trick Room |

### 17.2 Trick Room

Invierte el orden de velocidad durante 5 turnos. Los Pokémon más lentos actúan primero dentro de cada bracket de prioridad. La prioridad del movimiento **no se invierte**.

### 17.3 Tailwind

Duplica la Velocidad efectiva de todo el equipo durante 4 turnos.

### 17.4 Resolución de empates

Si dos Pokémon tienen la misma Velocidad efectiva en el mismo bracket de prioridad, el orden se decide al **azar** (50/50).

---

## 18. Resumen de Multiplicadores

### 18.1 Pila de multiplicadores ofensivos (se multiplican entre sí)

```
Daño final = Daño_base
    × STAB (1.0 / 1.5 / 2.0 / 2.25)
    × Efectividad de tipo (0 / 0.25 / 0.5 / 1 / 2 / 4)
    × Crítico (1.0 / 1.5 / 2.25 con Sniper)
    × Random roll (0.85–1.00)
    × Clima (0.5 / 1.0 / 1.5)
    × Terreno (1.0 / 1.3)
    × Quemadura (0.5 / 1.0)
    × Objeto ofensivo (1.0 / 1.1 / 1.2 / 1.3 / 1.5)
    × Habilidad ofensiva (varía)
    × Spread en dobles (0.75 / 1.0)
    × Pantallas (0.5 / 0.66 / 1.0)
    × Habilidad defensiva del rival (varía)
    × Objeto defensivo del rival (varía)
    × Otros modificadores
```

### 18.2 Ejemplo de cálculo completo

**Escenario:** Garchomp (Lv. 100, Adamant, 252 Atk EVs, 31 Atk IV, Choice Band) usa Earthquake contra Toxapex (Lv. 100, Bold, 252 HP / 252 Def EVs, 31 Def IV) bajo Grassy Terrain, sin crítico.

1. **Ataque efectivo de Garchomp:**
   - Base Atk: 130
   - Stat Atk = ((2×130 + 31 + 63) × 100/100 + 5) × 1.1 = **394** (Adamant = ×1.1)
   - Con Choice Band: 394 × 1.5 = **591**

2. **Defensa efectiva de Toxapex:**
   - Base Def: 152
   - Stat Def = ((2×152 + 31 + 63) × 100/100 + 5) × 1.1 = **427** (Bold = ×1.1 Def)

3. **Fórmula base:**
   - ((2×100/5 + 2) × 100 × 591 / 427) / 50 + 2 = **~119.2** → truncado a **119**

4. **Multiplicadores:**
   - STAB: Garchomp es Tierra/Dragón, Earthquake es Tierra → ×1.5
   - Efectividad: Tierra vs Veneno = ×2, Tierra vs Agua = ×1 → ×2
   - Crítico: No → ×1.0
   - Random: entre 0.85–1.00
   - Grassy Terrain: Earthquake hace ×0.5 bajo Grassy Terrain
   - No hay quemadura, pantallas, etc.

5. **Cálculo:**
   - 119 × 1.5 × 2 × 1.0 × 0.5 = **178** (con roll máximo)
   - Rango con roll: 178 × 0.85 = **151** a 178 × 1.00 = **178**

6. **HP de Toxapex:** ((2×50 + 31 + 63) × 100/100) + 100 + 10 = **304**
   - Daño: 151–178 / 304 = **49.7%–58.6%** del HP de Toxapex.

---

## Apéndice A — Tabla de Referencia Rápida de Factores

| Factor | Valor mínimo | Valor típico | Valor máximo |
|--------|-------------|-------------|-------------|
| STAB | ×1.0 | ×1.5 | ×2.25 (Tera + Adaptability) |
| Efectividad | ×0 | ×1.0 | ×4 |
| Crítico | ×1.0 | ×1.0 | ×2.25 (Sniper) |
| Random roll | ×0.85 | ×0.925 | ×1.00 |
| Clima | ×0.5 | ×1.0 | ×1.5 |
| Terreno | ×0.5 (Misty vs Dragón) | ×1.0 | ×1.3 |
| Choice Band/Specs | ×1.0 | ×1.0 | ×1.5 |
| Life Orb | ×1.0 | ×1.0 | ×1.3 |
| Quemadura (físico) | ×0.5 | ×1.0 | ×1.0 |
| Pantallas | ×0.5 | ×1.0 | ×1.0 |
| Spread en dobles | ×0.75 | ×1.0 | ×1.0 |
| Stat stages | ×0.25 | ×1.0 | ×4.0 |

## Apéndice B — Precisión: Fórmula Completa

```
P(acierto) = Precisión_movimiento × (Etapas_Accuracy / Etapas_Evasion) × Modificador_habilidad × Modificador_objeto
```

**Habilidades que afectan precisión:**

| Habilidad | Efecto |
|-----------|--------|
| Compound Eyes | Precisión ×1.3 |
| Hustle | Precisión de movimientos físicos ×0.8 |
| Victory Star | Precisión de aliados ×1.1 |
| Sand Veil | Evasión +20% bajo Sandstorm |
| Snow Cloak | Evasión +20% bajo Hail/Snow |
| Tangled Feet | Evasión ×2 mientras está confuso |
| No Guard | Todos los movimientos aciertan (ambos lados) |
| Keen Eye | No puede recibir bajadas de Accuracy |

**Objetos que afectan precisión:**

| Objeto | Efecto |
|--------|--------|
| Wide Lens | Precisión ×1.1 |
| Zoom Lens | Precisión ×1.2 si actúa después del rival |
| Bright Powder / Lax Incense | Evasión del portador ×0.9 contra ataques recibidos |

---

## Apéndice C — Daño Residual por Turno (resumen)

| Fuente | Daño por turno |
|--------|---------------|
| Quemadura | 1/16 HP máx |
| Veneno | 1/8 HP máx |
| Toxic | 1/16 × turno HP máx (acumulativo) |
| Sandstorm / Hail | 1/16 HP máx |
| Leech Seed | 1/8 HP máx (transferido al rival) |
| Curse (fantasma) | 1/4 HP máx |
| Bind/Wrap/etc. | 1/8 HP máx (1/6 con Binding Band) |
| Stealth Rock | 1/8 HP máx × efectividad de Roca contra tipos del Pokémon (1/32 a 1/2) |
| Spikes (1 capa) | 1/8 HP máx |
| Spikes (2 capas) | 1/6 HP máx |
| Spikes (3 capas) | 1/4 HP máx |
| Toxic Spikes (1 capa) | Envenena al entrar |
| Toxic Spikes (2 capas) | Envenena gravemente (Toxic) al entrar |
| Sticky Web | Baja Velocidad -1 al entrar |
| Life Orb | 1/10 HP máx por ataque |
| Black Sludge (no Veneno) | 1/8 HP máx |
| Leftovers | Recupera 1/16 HP máx |
| Black Sludge (Veneno) | Recupera 1/16 HP máx |
| Grassy Terrain | Recupera 1/16 HP máx |

---

> **Nota final:** Este documento cubre las mecánicas principales hasta la Generación IX (Pokémon Escarlata y Púrpura + DLC). Algunas mecánicas menores o edge cases extremadamente raros pueden no estar incluidos. Para cálculos exactos en batalla, se recomienda utilizar herramientas como el **Pokémon Damage Calculator** de Smogon.
