# go-clean-arch

Implementacion de referencia de Clean Architecture en Go.

## Documentacion en Espanol

La documentacion completa en espanol esta organizada en modulos dentro de [docs/README.md](docs/README.md).

Contenido principal:

1. Conceptos de Clean Architecture.
2. Mapa del proyecto y flujo entre capas.
3. Guia paso a paso para desarrolladores principiantes.
4. Explicacion de herramientas: Docker, Docker Compose, Dockerfile, Make, Go y MySQL.
5. Testing, calidad y troubleshooting.
6. Auditoria tecnica y plan de correccion priorizado.

## Inicio rapido

1. Copia variables de entorno:

```bash
cp example.env .env
```

2. Levanta base de datos local:

```bash
make dev-env
```

3. Inicia la aplicacion en modo desarrollo (hot reload):

```bash
make up
```

4. Prueba endpoint principal:

```bash
curl localhost:9090/articles
```

## Comandos utiles

```bash
make build
make lint
make tests
make go-generate
```

## Nota de alcance de idioma

- Documentacion: espanol.
- Codigo y comentarios de codigo: ingles.
