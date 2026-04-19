# 06 - Testing y Calidad

## Objetivo

Garantizar que la logica de negocio y adaptadores se comporten correctamente sin depender solo de pruebas manuales.

## Suites de pruebas existentes

- article/service_test.go
- internal/rest/article_test.go
- internal/repository/mysql/article_test.go
- internal/rest/middleware/cors_test.go
- internal/rest/middleware/timeout_test.go
- internal/repository/helper_test.go

## Estrategia por capa

### Casos de uso

Se testea el servicio con mocks de repositorio para validar reglas de negocio.

### Delivery HTTP

Se testean handlers para validar parsing, codigos HTTP y respuestas.

### Repositorio MySQL

Se testea comportamiento de consultas y persistencia con enfoque unitario.

### Middleware

Se valida comportamiento transversal (ejemplo CORS).

Tambien se valida timeout de request y propagacion de contexto.

### Helper de repositorio

Se valida codificacion y decodificacion de cursor (paginacion por tiempo).

## Comandos de calidad

- make lint
- make tests
- make tests-complete
- make build-race

## Buenas practicas recomendadas

- Agregar tests para paths de error y casos borde.
- Mantener pruebas pequenas y enfocadas.
- No acoplar tests de servicio a detalles de infraestructura.
- Usar mocks solo donde aporten aislamiento real.

## Cobertura en focos de Fase 3

Cobertura verificada con go test -coverprofile:

- internal/rest/article.go: 100%.
- internal/rest/middleware/timeout.go: 100%.
- internal/repository/helper.go: 100%.
