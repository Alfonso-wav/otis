# Pokemon companion animado en barra superior desde Settings

**ID**: 0143-settings-pokemon-companion
**Estado**: todo
**Fecha**: 2026-04-11

---

## Descripcion

Nueva funcionalidad en Settings: el usuario puede elegir un Pokemon para tenerlo con una animacion (sprite animado) en la barra superior (header). El sprite se posiciona a la derecha del todo en la barra donde dice "Pokedex". Por defecto sera un Diglett.

**Version navegador/APK**: como la barra superior ya tiene contenido (titulo + search + tabs), se debe crear una segunda fila en el header para el sprite del Pokemon companion, de modo que tenga su espacio propio a la derecha arriba.

**Version desktop (Wails)**: el sprite se coloca directamente en la barra superior a la derecha.

Requiere:
- Nuevo setting para elegir el Pokemon (selector/autocomplete).
- Persistencia de la eleccion en localStorage.
- Sprite animado renderizado en el header.

## Capas afectadas

- **Core**: ninguna
- **Shell**: ninguna
- **APP**: frontend — settings, header, estilos

## Archivos a crear/modificar

| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/index.html` | modificar | Agregar contenedor para el sprite companion en el header |
| `frontend/src/settings.ts` | modificar | Agregar setting de Pokemon companion con selector y persistencia |
| `frontend/src/main.ts` | modificar | Inicializar companion al arrancar la app, leer de localStorage |
| `frontend/src/styles/_pokemon.scss` | modificar | Estilos del sprite companion en el header |
| `frontend/src/styles/_settings.scss` | modificar | Estilos del nuevo setting (selector de Pokemon) |
| `frontend/src/locales/en.json` | modificar | Etiquetas del nuevo setting |
| `frontend/src/locales/es.json` | modificar | Etiquetas del nuevo setting |

## Plan de implementacion

1. Agregar en `index.html` un contenedor `#header-companion` dentro del header, posicionado a la derecha.
2. Para la version navegador/movil, crear una segunda fila (`header-companion-row`) que contenga el sprite.
3. Implementar en `settings.ts` un nuevo control: un campo de texto con autocomplete de nombres de Pokemon (reutilizar la logica de autocomplete existente) + preview del sprite seleccionado.
4. Guardar la eleccion en `localStorage` con key `companion-pokemon` (default: "diglett").
5. En `main.ts`, al inicializar, leer el companion de localStorage y renderizar el sprite animado en el header.
6. Usar sprites de PokemonDB o los sprites locales disponibles con animacion CSS (bounce suave o idle).
7. Agregar traducciones para el nuevo setting.

## Tests

| Archivo | Que se testea |
|---------|---------------|
| Manual | Verificar que Diglett aparece por defecto en el header |
| Manual | Verificar que se puede cambiar el Pokemon desde Settings |
| Manual | Verificar que la eleccion persiste al recargar |
| Manual (APK) | Verificar el layout con segunda fila en movil |
| Manual (Desktop) | Verificar el layout en la barra de Wails |

## Criterios de aceptacion

- [ ] Por defecto, Diglett aparece animado en el header a la derecha
- [ ] En Settings hay un selector para elegir otro Pokemon
- [ ] La eleccion se guarda en localStorage y persiste entre sesiones
- [ ] En version navegador/APK, el sprite esta en una segunda fila del header
- [ ] En version desktop (Wails), el sprite esta directamente en la barra superior derecha
- [ ] El sprite tiene una animacion idle (bounce, respiracion o similar)
- [ ] Las etiquetas del setting estan en EN y ES
- [ ] Funciona en dark mode y light mode

## Notas

El autocomplete de Pokemon ya existe en la app (busqueda del Pokedex). Se puede reutilizar esa logica para el selector en Settings. Los sprites animados pueden ser GIFs de PokemonDB o los sprites estaticos locales con animacion CSS.
