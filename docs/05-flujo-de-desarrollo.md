# 05 - Flujo de Desarrollo Diario

## Inicio de jornada

1. Sincroniza cambios remotos.
2. Verifica que Docker este activo.
3. Levanta entorno local con make up.

## Desarrollo de una funcionalidad

1. Lee flujo actual en handler -> service -> repository.
2. Implementa cambio en la capa correcta:
   - Regla de negocio: service.
   - Acceso a datos: repository.
   - Transporte HTTP: rest.
3. Si agregas interfaces o firmas, regenera mocks.

## Validaciones antes de commit

1. make lint
2. make tests
3. make build

## Criterios de calidad recomendados

- Evitar mover reglas de negocio a handlers o repositorios.
- Mantener contratos claros por interfaz.
- Propagar context.Context entre capas.
- Mapear errores de dominio a HTTP en delivery.

## Checklist rapido

- Cambio en capa correcta.
- Tests actualizados.
- Linter limpio.
- Build exitoso.
- Endpoint validado manualmente.
