# 08 - Auditoria Tecnica y Plan de Correccion

Este documento resume hallazgos tecnicos detectados y una propuesta de prioridad para corregirlos.

## Estado de implementacion (actualizacion)

Cambios aplicados en esta etapa:

- Corregido: cierre de statements en repositorio de articulos (Store, Delete, Update).
- Corregido: cierre de statement en repositorio de autores (getOne).
- Corregido: simplificacion de concurrencia en fillAuthorDetails para usar un solo Wait de errgroup.

Validacion posterior a cambios:

- Pruebas unitarias ejecutadas: 28 exitosas.
- make build: exitoso.

Pendiente en fase de estabilidad:

- Ninguno.

Resuelto adicionalmente en esta etapa:

- Compatibilidad Docker Compose en Makefile (soporta docker-compose y docker compose).
- Lint operativo en entorno actual mediante ejecucion con go run y perfil compatible de linters.
- Timeout por request ajustado por entorno con default mas seguro.
- Configuracion de pool de conexiones SQL en inicializacion de la aplicacion.
- CORS por entorno con lista de origenes permitidos (allowlist) y fallback controlado para desarrollo.

## Estado verificado de ejecucion

Resultados ejecutados en este repositorio:

- make build: exitoso.
- make tests: exitoso (44 tests, cobertura reportada por paquete).
- make lint: exitoso.
- make dev-env: exitoso tras ajuste de compatibilidad de comando compose.

## Hallazgos principales

## Critico

1. Posible fuga de recursos SQL por statements no cerrados.
   - Ubicacion: internal/repository/mysql/article.go (Store, Delete, Update)
   - Ubicacion: internal/repository/mysql/author.go (getOne)
   - Riesgo: consumo innecesario de recursos y degradacion bajo carga.

## Alto

2. Patrón de concurrencia en fillAuthorDetails con doble Wait en errgroup.
   - Ubicacion: article/service.go
   - Estado: resuelto y validado con suite ejecutada con race detector.

3. Timeout de contexto posiblemente muy bajo para entorno real.
   - Ubicacion: example.env (CONTEXT_TIMEOUT=2)
   - Estado: resuelto (nuevo default recomendado: 15s).

4. Desfase de toolchain para typecheck estricto.
   - Ubicacion: Makefile, go.mod, .golangci.yaml
   - Estado: mitigado con perfil de lint compatible.
   - Riesgo residual: menor cobertura de chequeos semanticos en lint local.

## Medio

5. CORS abierto para cualquier origen.
   - Ubicacion: internal/rest/middleware/cors.go
   - Estado: mitigado con allowlist por entorno.

6. Falta de configuracion de pool de conexiones SQL.
   - Ubicacion: app/main.go
   - Estado: resuelto (SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime, SetConnMaxIdleTime).

7. Cobertura de tests incompleta en algunos caminos.
   - Casos no cubiertos de forma robusta: Update en handler y service, helper de cursor, timeout middleware.

## Plan de correccion priorizado (sin implementacion en esta fase)

## Fase 1 - Estabilidad (inmediato)

1. Cerrar statements con defer Close en repositorios MySQL.
2. Simplificar fillAuthorDetails a un solo punto de Wait seguro.
3. Evaluar retorno progresivo de linters semanticos cuando el entorno use toolchain totalmente alineado.
4. Ajustar timeout por defecto para entorno local y documentar criterio (resuelto).

Estado de Fase 1:

- Items 1, 2 y 4: resueltos.
- Item 3: pendiente estrategico (mejora incremental de lint semantico).

## Fase 2 - Seguridad y operacion

5. Parametrizar CORS por entorno.
6. Configurar pool de conexiones en inicializacion DB (resuelto).
7. Compatibilidad Docker Compose en comandos de Make (resuelto).

## Fase 3 - Calidad y mantenimiento

8. Incrementar cobertura de tests en caminos faltantes.
9. Documentar contrato de errores y codigos HTTP.

Estado de Fase 3:

- Item 8: resuelto (cobertura agregada para helper de cursor, timeout middleware y ramas de error en handlers).
- Item 9: resuelto (documento de contrato de errores/codigos HTTP agregado en docs/09-contrato-errores-http.md).

## Criterios de aceptacion sugeridos

- make lint limpio.
- make tests verde.
- make build exitoso.
- Endpoints principales funcionales con entorno Docker local.
- Documentacion actualizada con cambios tecnicos aplicados.
