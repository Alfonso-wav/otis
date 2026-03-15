# Mejoras en gestion de equipos en Builds

**ID**: 0040-team-editor-ux
**Estado**: done
**Fecha**: 2026-03-15

---

## Descripcion

Mejorar la experiencia de gestion de equipos en la pagina de Builds con tres cambios principales:

1. **Equipos arriba**: Mover la seccion "Mis equipos" al inicio de la pagina de Builds (actualmente aparece al final, despues del simulador de batalla).
2. **Guardar en equipo existente**: Al guardar un Pokemon en un equipo, mostrar un selector con los equipos existentes ademas de la opcion de crear uno nuevo (actualmente solo se usa `prompt()` para escribir el nombre).
3. **Edicion completa de equipos**: Permitir editar equipos de forma integral — sacar Pokemon, anadir nuevos miembros a equipos existentes desde la propia seccion de equipos.

## Capas afectadas
- **Core**: Sin cambios (ya existen `AddMemberToTeam`, `RemoveMemberFromTeam`, validaciones).
- **Shell**: Sin cambios (CRUD de equipos ya funciona).
- **APP**: Sin cambios (bindings ya exponen todas las operaciones necesarias).
- **Frontend**: Cambios en `frontend/src/pages/builds.ts` — layout, modal de guardado y controles de edicion.

## Dependencias externas nuevas
Ninguna.

## Archivos a crear/modificar
| Archivo | Accion | Descripcion |
|---------|--------|-------------|
| `frontend/src/pages/builds.ts` | modificar | Reposicionar seccion equipos, nuevo modal de guardado, controles de edicion |

## Plan de implementacion

### Paso 1 — Mover seccion de equipos arriba
1. En `buildLayout()` (linea ~825), mover `${teamsSection}` de su posicion actual (despues de `${btlSection}`) a justo despues del `section-header`, antes del `build-layout`.
2. La seccion de equipos debe aparecer como lo primero visible al entrar en Builds.
- [x] Completado

### Paso 2 — Modal/dropdown de guardado con equipos existentes
3. Reemplazar el `prompt()` actual en `saveToTeam()` (linea ~635) por un modal/dropdown inline que muestre:
   - Lista de equipos existentes con nombre y numero de miembros (ej: "team a (1/6)") — solo los que tengan menos de 6 miembros.
   - Opcion "Crear equipo nuevo" con un campo de texto para el nombre.
   - Boton de confirmar/cancelar.
4. Al seleccionar un equipo existente, llamar a `SaveToTeam()` con ese nombre.
5. Al crear uno nuevo, pedir nombre y llamar a `SaveToTeam()` con el nombre nuevo.
- [x] Completado

### Paso 3 — Edicion de equipos: anadir miembro desde seccion de equipos
6. En cada tarjeta de equipo (si tiene menos de 6 miembros), agregar un boton "Anadir Pokemon" que abra un buscador de Pokemon (reutilizando el autocomplete existente de nombres).
7. Al seleccionar un Pokemon, crear un `TeamMember` con valores por defecto (nivel 50, Hardy, IVs 31, EVs 0, sin moves) y llamar a `SaveToTeam()`.
8. Alternativamente, permitir "Anadir desde build actual" que tome el attacker o defender configurado actualmente.
- [x] Completado

### Paso 4 — Pulir interaccion
9. Asegurar que la seccion de equipos se refresca correctamente despues de cada operacion (ya existe `cachedTeams = await ListTeams(); buildLayout()`).
10. Mantener el `<details>` abierto tras operaciones de edicion (actualmente se cierra al hacer rebuild del layout).
- [x] Completado

## Tests
| Archivo | Que se testea |
|---------|---------------|
| N/A | Cambios puramente de frontend/UI. Core y Shell ya estan testeados. Verificacion manual. |

## Criterios de aceptacion
- [x] La seccion "Mis equipos" aparece al inicio de la pagina de Builds, antes del configurador de attacker/defender
- [x] Al pulsar "Guardar en equipo" se muestra un modal con la lista de equipos existentes y opcion de crear nuevo
- [x] Se puede anadir un Pokemon a un equipo existente desde el modal sin escribir el nombre manualmente
- [x] Se puede anadir un nuevo miembro a un equipo desde la propia tarjeta del equipo
- [x] Se puede eliminar miembros individuales (ya existente, verificar que sigue funcionando)
- [x] El estado del `<details>` (abierto/cerrado) se preserva tras operaciones de edicion
- [x] Equipos llenos (6/6) no muestran opcion de anadir mas miembros

## Notas
- No se necesitan cambios en Go — todo el backend ya soporta las operaciones necesarias.
- El modal de guardado reemplaza al `prompt()` nativo del navegador, mejorando la UX.
- Reutilizar estilos existentes de la pagina (`.team-card`, `.team-member-row`, etc.).
