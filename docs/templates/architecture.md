# {nombre_proyecto}

**Fecha**: {fecha}
**Estado**: diseño

---

## Proposito

{que_resuelve_este_proyecto_en_1_2_frases}

## Modulo Go

```
module github.com/{usuario}/{nombre_proyecto}

go {version}
```

---

## Dominio

### Entidades

| Entidad | Descripcion | Campos clave |
|---------|-------------|--------------|
| `{Entidad}` | {descripcion} | `{campo1}`, `{campo2}` |

### Tipos de dominio (Core)

```go
// core/domain.go

type {Entidad} struct {
    {Campo} {tipo}
}

// Enums como constantes
type {TipoEnum} int

const (
    {Valor1} {TipoEnum} = iota
    {Valor2}
)
```

### Reglas de negocio

- {regla_1}
- {regla_2}

---

## Capas

### Core — logica pura

Funciones sin efectos secundarios. Solo transforman datos.

| Archivo | Responsabilidad |
|---------|-----------------|
| `core/domain.go` | Tipos de dominio, structs, enums, constantes |
| `core/{feature}.go` | {funciones_puras_que_contiene} |
| `core/{feature}_test.go` | Tests unitarios de {feature} |

**Interfaces que Core define** (las implementa Shell):

```go
// core/ports.go

type {Repositorio} interface {
    {Metodo}({params}) ({retorno}, error)
}
```

### Shell — adaptadores externos

Todo lo que tiene efectos secundarios.

| Archivo | Tipo | Dependencia externa |
|---------|------|---------------------|
| `shell/{adaptador}.go` | {db/http/storage/queue} | {paquete_o_servicio} |

**Implementa**:

| Interfaz (de Core) | Implementacion (en Shell) |
|---------------------|---------------------------|
| `{Repositorio}` | `shell/{adaptador}.go` |

### APP — cableado y arranque

| Archivo | Responsabilidad |
|---------|-----------------|
| `app/main.go` | Entry point, inyeccion de dependencias |
| `app/config.go` | Lectura de configuracion (env vars, flags) |
| `app/server.go` | Setup del servidor ({http/grpc/cli}) |
| `app/handlers/{handler}.go` | Handler que orquesta Core + Shell |

---

## Estructura de archivos

```
{nombre_proyecto}/
  core/
    domain.go
    ports.go
    {feature}.go
    {feature}_test.go
  shell/
    {adaptador}.go
  app/
    main.go
    config.go
    server.go
    handlers/
      {handler}.go
  go.mod
  go.sum
```

---

## Dependencias externas

| Paquete | Version | Capa | Para que |
|---------|---------|------|----------|
| `{paquete}` | `{version}` | Shell | {proposito} |

---

## Configuracion

| Variable | Tipo | Default | Descripcion |
|----------|------|---------|-------------|
| `{ENV_VAR}` | {tipo} | `{default}` | {descripcion} |

---

## Endpoints / Comandos

| Metodo | Ruta / Comando | Handler | Descripcion |
|--------|----------------|---------|-------------|
| `{GET/POST/CLI}` | `{ruta}` | `{handler}` | {descripcion} |

---

## Flujo de datos

```
{entrada} -> Handler (APP)
                |
                v
            Core.{Funcion}({datos})  // logica pura
                |
                v
            Shell.{Metodo}({resultado})  // efecto secundario
                |
                v
            {salida}
```

---

## Plan de implementacion

1. [ ] Definir tipos de dominio en `core/domain.go`
2. [ ] Definir interfaces (ports) en `core/ports.go`
3. [ ] Implementar logica pura en `core/{feature}.go` + tests
4. [ ] Implementar adaptadores en `shell/`
5. [ ] Cablear todo en `app/` (config, server, handlers)
6. [ ] Tests de integracion

---

## Decisiones de diseno

| Decision | Alternativa descartada | Razon |
|----------|------------------------|-------|
| {decision} | {alternativa} | {por_que} |

---

## Notas

{contexto_adicional}
