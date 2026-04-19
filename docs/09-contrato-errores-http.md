# 09 - Contrato de Errores y Codigos HTTP

Este documento define como el modulo REST de articulos traduce errores de dominio y errores de entrada a respuestas HTTP.

## Objetivo

- Hacer explicito el contrato de errores para clientes API.
- Reducir ambiguedad al agregar nuevos handlers.
- Mantener consistencia entre dominio y transporte HTTP.

## Errores de dominio

Definidos en domain/errors.go:

- ErrInternalServerError
- ErrNotFound
- ErrConflict
- ErrBadParamInput

## Mapeo actual en handlers

Implementado en internal/rest/article.go mediante getStatusCode(err):

- ErrInternalServerError -> 500 Internal Server Error
- ErrNotFound -> 404 Not Found
- ErrConflict -> 409 Conflict
- Cualquier error no reconocido -> 500 Internal Server Error

## Casos de validacion y parsing (sin getStatusCode)

Ademas del mapeo de dominio, el handler responde codigos directos en errores de entrada:

- Parametro id invalido en ruta -> 404 Not Found
- Body JSON invalido en Store -> 422 Unprocessable Entity
- Body valido en JSON pero invalido por reglas de negocio basicas (validator) -> 400 Bad Request

## Formato de respuesta de error

Hay dos formas actuales de respuesta de error:

1. Estructura JSON estandar:

{
  "message": "detalle"
}

Se usa en errores provenientes de servicio/dominio via getStatusCode.

1. String JSON simple (inconsistente):

"detalle"

Se usa en errores de parsing/parametros en algunos handlers.

## Contrato recomendado para evolucion

Para mejorar consistencia en clientes:

1. Unificar todas las respuestas de error al formato {"message": "..."}.
2. Reservar 404 para recursos no encontrados y evaluar 400 para id invalido.
3. Incluir codigo interno opcional (ejemplo: ARTICLE_NOT_FOUND) en futuras iteraciones.

## Checklist para nuevos endpoints

1. Definir errores de dominio esperados.
2. Mapear errores a HTTP via getStatusCode o helper equivalente.
3. Mantener formato de respuesta consistente.
4. Agregar tests para:
   - camino feliz
   - error de validacion
   - error de dominio
   - error desconocido (fallback 500)
