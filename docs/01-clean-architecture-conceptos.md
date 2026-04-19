# 01 - Conceptos de Clean Architecture

## Que es Clean Architecture

Clean Architecture es una forma de organizar software para que la logica de negocio no dependa de frameworks, base de datos o transporte HTTP.

Principios clave:

- Independencia de frameworks.
- Reglas de negocio testeables sin infraestructura real.
- Independencia de UI y base de datos.
- Dependencias apuntando hacia el centro (dominio y casos de uso).

## Capas en este proyecto

Este proyecto usa 4 capas principales:

1. Dominio.
2. Casos de uso (servicios).
3. Repositorios (contratos + implementaciones).
4. Delivery HTTP.

## Como se aplican los conceptos en el repositorio

### Dominio

Define entidades y errores del negocio.

- Entidades: domain/article.go, domain/author.go
- Errores de dominio: domain/errors.go

Estas estructuras no importan Echo, SQL ni Docker.

### Casos de uso

Orquestan reglas de negocio con contratos (interfaces).

- Servicio principal: article/service.go
- Contratos de dependencias: interfaces ArticleRepository y AuthorRepository en article/service.go

Aqui vive la logica de negocio, por ejemplo:

- Verificar conflictos al guardar (Store).
- Completar autores para un lote de articulos (fillAuthorDetails).
- Validar existencia antes de eliminar (Delete).

### Repositorios

Implementan acceso a datos sin mover reglas de negocio al SQL.

- Implementacion MySQL de articulos: internal/repository/mysql/article.go
- Implementacion MySQL de autores: internal/repository/mysql/author.go
- Helpers de cursor de paginacion: internal/repository/helper.go

### Delivery HTTP

Recibe requests, valida entradas, invoca casos de uso y mapea errores a codigos HTTP.

- Handler de articulos: internal/rest/article.go
- Middlewares: internal/rest/middleware/

## Inversion de dependencias en la practica

Las interfaces se declaran del lado consumidor:

- El servicio define que necesita de repositorios.
- El handler define que necesita de servicio.

Luego, en app/main.go, se conectan implementaciones concretas:

- repositorio MySQL -> servicio -> handler REST.

Esto permite cambiar infraestructura sin romper la logica de negocio.
