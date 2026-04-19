# 03 - Guia de Primeros Pasos

Esta guia esta pensada para una persona principiante.

## 1. Prerequisitos

Necesitas tener instalado:

- Go (version compatible con go.mod)
- Docker
- Docker Compose
- Make
- Git

Nota importante:

- El Makefile usa el comando docker-compose (CLI clasico con guion).
- Si solo tienes el plugin nuevo docker compose (sin guion), debes instalar compatibilidad o crear un alias para docker-compose.

Comandos utiles para validar:

- go version
- docker --version
- docker-compose --version
- docker compose version
- make --version
- git --version

## 2. Clonar y preparar variables de entorno

1. Clona el repositorio.
2. En la raiz, copia example.env a .env.
3. Revisa los valores de conexion a BD y puerto del servidor.

## 3. Levantar base de datos MySQL

Ejecuta:

make dev-env

Que ocurre:

- Se levanta el servicio mysql definido en compose.yaml.
- Se ejecuta article.sql al iniciar MySQL.
- Se expone el puerto 3306.

## 4. Levantar aplicacion en modo desarrollo

Ejecuta:

make up

Que ocurre:

- Levanta MySQL (si no estaba activo).
- Ejecuta Air para hot reload.
- Inicia API en el puerto configurado (por defecto 9090).

## 5. Probar API rapidamente

En otra terminal:

curl localhost:9090/articles

## 6. Comandos basicos para empezar a contribuir

- make build: compila el binario.
- make lint: valida estilo y problemas estaticos.
- make tests: ejecuta pruebas.
- make go-generate: regenera mocks.

## 7. Flujo recomendado para un primer cambio

1. Crear una rama.
2. Hacer un cambio pequeno.
3. Ejecutar make lint.
4. Ejecutar make tests.
5. Validar manualmente endpoint relacionado.
6. Abrir PR con descripcion clara.
