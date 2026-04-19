package domain

import "errors"

// Domain Errors represent business rules violations and domain-level failures.
// These errors are defined in the domain layer (not in HTTP/REST layer) because:
//   - They represent BUSINESS logic, not presentation details
//   - Multiple interfaces (REST, gRPC, CLI) can reuse the same errors
//   - The HTTP status code mapping (error -> 400/404/500) happens in the REST layer
//
// Example flow:
//   1. Service validates and returns ErrNotFound (domain error)
//   2. REST Handler catches it and converts to HTTP 404 (presentation concern)
//   3. Future gRPC handler could convert same ErrNotFound to grpc.NotFound code

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")
)
