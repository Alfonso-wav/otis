# Task 0132 — Fix simulation autocomplete Enter error

## Estado: done

## Goal

Fix an error popup that appears when selecting a Pokemon with Enter key in the simulation autocomplete (Builds tab). The keydown handler fires before the autocomplete handler, sending the raw partial text to fetchPokemon, which fails.

---

## Contexto tecnico

### Archivos principales

- `frontend/src/pages/builds.ts` — `bindEvents()` function where keydown listeners and autocomplete are registered

### Estado actual

The keydown event listeners on atk/def inputs are registered before the `createAutocomplete()` calls. When the user selects a Pokemon from the autocomplete dropdown with Enter, the keydown handler fires first with the partial text (e.g. "sn" instead of "Snorlax"), triggering `fetchPokemon` with an invalid name and causing an error popup.

---

## Cambios requeridos

### 1. Reorder event binding — `frontend/src/pages/builds.ts`

Move `createAutocomplete()` calls BEFORE the keydown event listeners so the autocomplete's Enter handler (which calls `e.preventDefault()`) fires first.

### 2. Add guard to keydown handlers

Add `&& !e.defaultPrevented` to the Enter key condition in both atk and def keydown handlers, so that when autocomplete consumes the Enter, the builds handler skips its `fetchPokemon` call.

---

## Archivos afectados

### Frontend
- `frontend/src/pages/builds.ts` — reorder event binding and add defaultPrevented guard

### Backend
Ninguno.

### Core
Ninguno.

### Shell
Ninguno.

---

## Acceptance criteria

- [x] Selecting a Pokemon from autocomplete with Enter does not trigger an error popup.
- [x] Typing a full Pokemon name and pressing Enter still works correctly.
- [x] Click-based autocomplete selection still works.
- [x] Both atk and def inputs are fixed.

---

## Dependencias

Ninguna.
