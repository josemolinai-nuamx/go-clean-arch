# 02 - Mapa del Proyecto

## Estructura general

- app/: punto de arranque de la aplicacion.
- domain/: entidades y errores de negocio.
- article/: casos de uso del modulo de articulos.
- internal/repository/mysql/: adaptador de persistencia MySQL.
- internal/rest/: adaptador HTTP (Echo).
- internal/rest/middleware/: middlewares transversales.
- docs/: documentacion en espanol.

## Flujo principal de una request

Ejemplo GET /articles:

1. Llega al handler HTTP en internal/rest/article.go.
2. El handler llama al servicio en article/service.go.
3. El servicio usa el contrato de repositorio.
4. La implementacion MySQL ejecuta SQL en internal/repository/mysql/article.go.
5. El resultado vuelve al servicio y luego al handler.

## Bootstrap del sistema

En app/main.go:

1. Carga variables de entorno desde .env.
2. Crea conexion a MySQL.
3. Inicializa Echo y middlewares.
4. Instancia repositorios concretos.
5. Instancia el servicio de articulos.
6. Registra endpoints y arranca servidor.

## Endpoints actuales

- GET /articles
- GET /articles/:id
- POST /articles
- DELETE /articles/:id
