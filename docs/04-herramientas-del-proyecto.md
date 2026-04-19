# 04 - Herramientas del Proyecto

## Go

Lenguaje principal del proyecto.

Uso en este repo:

- Compilacion del binario (make build).
- Ejecucion de pruebas unitarias.
- Gestion de dependencias con go.mod.

## Make

Orquesta comandos repetitivos de desarrollo.

Archivo principal:

- Makefile

Comandos importantes:

- make up
- make down
- make destroy
- make build
- make lint
- make tests
- make install-deps

## Docker

Empaqueta la aplicacion en contenedores.

Archivo:

- Dockerfile

Patron usado:

- Multi-stage build: una etapa para compilar y otra etapa minima para ejecutar.

## Docker Compose

Coordina multiples servicios.

Archivo:

- compose.yaml

Servicios:

- mysql: base de datos.
- web: API del proyecto.

Detalle operativo:

- Los targets de Makefile detectan automaticamente docker-compose o docker compose.
- Esto permite ejecutar make dev-env y make up en entornos con CLI clasico o plugin moderno.

## MySQL

Persistencia principal.

Esquema inicial:

- article.sql

Variables de entorno relacionadas:

- DATABASE_HOST
- DATABASE_PORT
- DATABASE_USER
- DATABASE_PASS
- DATABASE_NAME
- CONTEXT_TIMEOUT
- DB_MAX_OPEN_CONNS
- DB_MAX_IDLE_CONNS
- DB_CONN_MAX_LIFETIME_SEC
- DB_CONN_MAX_IDLE_TIME_SEC
- CORS_ALLOWED_ORIGINS

Valores recomendados para entorno local:

- CONTEXT_TIMEOUT = 15
- DB_MAX_OPEN_CONNS = 25
- DB_MAX_IDLE_CONNS = 10
- DB_CONN_MAX_LIFETIME_SEC = 300
- DB_CONN_MAX_IDLE_TIME_SEC = 120
- CORS_ALLOWED_ORIGINS = "http://localhost:3000,http://localhost:5173"

Nota de CORS:

- En desarrollo, si DEBUG=true y CORS_ALLOWED_ORIGINS esta vacio, el middleware permite *.
- En otros entornos, se recomienda definir explicitamente CORS_ALLOWED_ORIGINS.

## Air

Herramienta de recarga en caliente durante desarrollo.

Uso:

- Se invoca con make up (objetivo dev-air).

## golangci-lint

Analisis estatico del codigo para detectar problemas de calidad.

Uso:

- make lint

## gotestsum y tparse

Mejoran salida y analisis de tests.

Uso:

- make tests
- make tests-complete

## mockery

Genera mocks a partir de interfaces para pruebas unitarias.

Uso:

- make go-generate
