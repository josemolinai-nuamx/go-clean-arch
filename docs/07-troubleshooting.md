# 07 - Troubleshooting

## Error: no se puede cargar .env

Sintoma:

- Mensaje de error al iniciar app por archivo .env faltante.

Solucion:

1. Copiar example.env a .env.
2. Verificar que ejecutes comandos desde la raiz del repo.

## Error: fallo de conexion a MySQL

Sintoma:

- Mensajes de failed to ping database o connection refused.

Solucion:

1. Ejecutar make dev-env.
2. Esperar healthcheck de MySQL.
3. Verificar variables DATABASE_* en .env.

## Error: puerto ocupado

Sintoma:

- No puede iniciar en 9090 o 3306.

Solucion:

1. Liberar el puerto en uso.
2. Cambiar SERVER_ADDRESS en .env.
3. Ajustar puertos en compose.yaml si aplica.

## Error: herramientas no encontradas

Sintoma:

- make lint, make tests o air fallan por binario faltante.

Solucion:

1. Ejecutar make install-deps.
2. Ejecutar make deps para validar.

## Error: tests intermitentes por timeout

Sintoma:

- Falla por timeout en tests o requests.

Solucion:

1. Revisar CONTEXT_TIMEOUT en .env.
2. Revisar carga de maquina local.
3. Repetir make tests para confirmar reproducibilidad.
