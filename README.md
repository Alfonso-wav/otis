# Otis

<p align="center">
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/149.png" width="180" alt="Dragonite" />
</p>

<p align="center">
  <strong>Pokedex, simulador de combates y constructor de equipos.</strong><br/>
  Escrito en Go + TypeScript. Desktop y Android.
</p>

<p align="center">
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png" width="56" />
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/6.png" width="56" />
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/150.png" width="56" />
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/384.png" width="56" />
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/448.png" width="56" />
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/94.png" width="56" />
</p>

---

## Que es Otis

Una aplicacion de escritorio y movil para explorar el mundo Pokemon: buscar Pokemon, consultar movimientos y habilidades, recorrer regiones, montar equipos y simular combates con calculo de dano real.

Todo el backend es Go puro con arquitectura funcional (Core / Shell / App). El frontend es TypeScript con Bootstrap, ECharts y GSAP. Corre como app nativa via Wails (desktop) o Capacitor (Android).

---

## Funcionalidades

<table>
<tr>
<td width="50%">

### Pokedex
<img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/133.png" width="40" />

- Mas de 1000 Pokemon con busqueda y autocompletado
- Filtros por generacion, tipo, legendario/mitico
- Vista en grid o tabla
- Stats, sprites, cadenas evolutivas, movimientos

</td>
<td width="50%">

### Explorar
<img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/249.png" width="40" />

- Regiones y localizaciones
- Catalogo completo de movimientos (poder, precision, PP, prioridad)
- Catalogo de habilidades con efectos
- Tabla de tipos y efectividades

</td>
</tr>
<tr>
<td>

### Constructor de equipos
<img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/445.png" width="40" />

- Equipos de hasta 6 Pokemon
- Nivel, naturaleza, IVs y EVs configurables
- Validacion de limites (510 EVs totales, 252 por stat)
- 4 movimientos por Pokemon
- Guardado persistente en JSON

</td>
<td>

### Simulador de combates
<img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/257.png" width="40" />

- 1v1 con stats y movimientos custom
- Combates turno a turno con seleccion de movimiento
- Auto-simulacion (hasta 200 turnos)
- Batch mode: 1 a 10.000 combates con estadisticas
- Combates de equipo 6v6 con HP arrastrado y cambios

</td>
</tr>
<tr>
<td>

### Calculadora de dano
<img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/376.png" width="40" />

- Formula real de dano Pokemon
- STAB, matchups de tipo, golpes criticos
- Rango de dano (min/max)

</td>
<td>

### Calculadora de EVs y stats
<img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/196.png" width="40" />

- Estimar EVs a partir de stats actuales
- Calcular stats finales con base, IVs, EVs, nivel y naturaleza
- Preview de spreads de EVs

</td>
</tr>
</table>

---

## Arquitectura

```
                    APP (wiring)
                   /           \
                  v             v
              Shell            Core
          (adaptadores)    (logica pura)
              |
              v
         APIs externas
```

| Capa | Responsabilidad | Efectos secundarios |
|------|----------------|---------------------|
| **Core** | Logica de negocio, tipos de dominio, calculo de dano, simulacion de batalla | Ninguno. Funciones puras. |
| **Shell** | Clientes HTTP (PokeAPI), scraping (PokemonDB), almacenamiento de equipos | Todos los I/O van aqui. |
| **App** | Entry point, configuracion, bindings Wails / API REST para movil | Cableado e inyeccion. |

La dependencia siempre fluye hacia adentro: `App -> Shell -> Core`. Core no importa nada del proyecto.

---

## Estructura del proyecto

```
otis/
  core/              # Logica pura: dominio, batallas, dano, EVs, equipos
  shell/             # Adaptadores: PokeAPI, PokemonDB scraper, storage
  app/               # Wails bindings (desktop) + servidor REST (movil)
    mobile/          # HTTP handlers y gomobile exports
  frontend/          # TypeScript + Bootstrap + ECharts + GSAP
    src/
      pages/         # Pokedex, builds, types, explore
      components/    # Modals y componentes reutilizables
      charts/        # Visualizaciones con ECharts
      animations/    # Transiciones con GSAP
  android/           # Proyecto Capacitor para Android
  scripts/           # Build scripts (gomobile)
  data/teams/        # Equipos guardados (JSON)
```

---

## Stack

<table>
<tr>
<td>

**Backend**
- Go 1.25
- Wails v2 (desktop)
- net/http (API movil)
- gomobile (Android .aar)
- goquery (scraping)

</td>
<td>

**Frontend**
- TypeScript
- Vite 6
- Bootstrap 5
- ECharts 6
- GSAP 3
- Sass

</td>
<td>

**Movil**
- Capacitor 8
- Android SDK (API 21+)
- WebView + Go server local

</td>
</tr>
</table>

---

## Como ejecutar

### Desktop

```bash
# Requisitos: Go 1.25+, Node.js 18+, npm

# Instalar dependencias del frontend
npm install --prefix frontend

# Lanzar la app
go run main.go
```

### Android

```bash
# Requisitos: gomobile, Android SDK/NDK, JAVA_HOME

# Compilar backend Go como .aar
./scripts/build-android.sh

# Compilar APK
cd android && ./gradlew assembleDebug
```

El frontend detecta automaticamente el entorno (Wails IPC vs HTTP fetch) a traves de `api.ts`.

---

## APIs externas

| Fuente | Uso |
|--------|-----|
| [PokeAPI](https://pokeapi.co) | Datos de Pokemon, movimientos, habilidades, especies, regiones |
| [PokemonDB](https://pokemondb.net) | Sprites y datos adicionales via scraping |

---

<p align="center">
  <img src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/143.png" width="72" />
  <br/>
  <sub>Otis descansando despues de 10.000 simulaciones de combate.</sub>
</p>
